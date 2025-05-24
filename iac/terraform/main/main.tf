data "google_project" "project" {
  project_id = var.GCP_PROJECT_ID
}

resource "google_artifact_registry_repository" "food-interpreter-repository" {
  location      = var.GCP_PROJECT_REGION
  repository_id = "food-interpreter-repository"
  description   = "Food Interpreter docker repository"
  format        = "DOCKER"

  cleanup_policies {
    id = "delete-old"
    action = "DELETE"
    condition {
      older_than = "30d"
    }
  }
  cleanup_policies {
    id     = "keep-amount"
    action = "KEEP"
    most_recent_versions {
      keep_count = 5
    }
  }
}

# TODO move this out, same as gateway
resource "google_cloud_run_v2_service" "interpreter" {
  name     = "interpreter-cloud-run"
  location = var.GCP_PROJECT_REGION
  deletion_protection = false

  template {
    timeout = "10s"
    containers {
      image = "${ var.GCP_PROJECT_REGION }-docker.pkg.dev/${ var.GCP_PROJECT_ID }/${ google_artifact_registry_repository.food-interpreter-repository.name }/food-interpreter-interpreter:latest"
      resources {
        limits = {
          memory = "1024Mi"
          cpu    = "1000m"
        }
      }
    }
  }
}

resource "google_service_account" "interpreter_cloud_run" {
  account_id   = "interpreter-cloud-run"
  display_name = "Interpreter Cloud Run Service Account"
}

resource "google_cloud_run_service_iam_binding" "interpreter_servicesinvoker" {
  location = google_cloud_run_v2_service.interpreter.location
  service  = google_cloud_run_v2_service.interpreter.name
  role     = "roles/run.servicesInvoker"
  members  = ["serviceAccount:${ google_service_account.interpreter_cloud_run.email }"]
}

resource "google_project_iam_binding" "interpreter_logwriter" {
  project = var.GCP_PROJECT_ID
  role    = "roles/logging.logWriter"
  members = ["serviceAccount:${ google_service_account.interpreter_cloud_run.email }"]
}

resource "google_project_service_identity" "pubsub_agent" {
  provider = google-beta
  project  = data.google_project.project.project_id
  service  = "pubsub.googleapis.com"
}

resource "google_project_iam_binding" "project_token_creator" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  members = ["serviceAccount:${ google_project_service_identity.pubsub_agent.email }"]
}

# TODO schema?
resource "google_pubsub_topic" "interpreter" {
  name = "interpreter-topic"

  labels = {
    service = "interpreter"
  }

  message_retention_duration = "86000s"
}

resource "google_pubsub_subscription" "interpreter" {
  name                  = "interpreter-subscription"
  topic                 = google_pubsub_topic.interpreter.id
  ack_deadline_seconds  = 20
  retain_acked_messages = false

  labels = {
    service = "interpreter"
  }

  push_config {
    push_endpoint = "${ google_cloud_run_v2_service.interpreter.uri }/interpret"
    oidc_token {
      service_account_email = google_service_account.interpreter_cloud_run.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }

  dead_letter_policy {
    dead_letter_topic = google_pubsub_topic.interpreter-dlq.id
    max_delivery_attempts = 5
  }

  depends_on = [ google_cloud_run_v2_service.interpreter ]
}

# TODO 404 messages aren't going to DLQ. Check Dead Lettering issues on console.
# TODO 404s after renaming everything when placing message straight to pubsub topic
resource "google_pubsub_topic" "interpreter-dlq" {
  name = "interpreter-topic-dlq"

  labels = {
    service = "interpreter"
  }

  message_retention_duration = "604800s" # 7 days
}

resource "google_pubsub_subscription" "interpreter-dlq" {
  name                 = "interpreter-subscription-dlq"
  topic                = google_pubsub_topic.interpreter-dlq.id
  ack_deadline_seconds = 20

  labels = {
    service = "interpreter"
  }
}

resource "google_project_iam_member" "subscriber_role" {
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:service-${var.GCP_PROJECT_ID}@gcp-sa-pubsub.iam.gserviceaccount.com"
  project = var.GCP_PROJECT_ID
}

resource "google_project_iam_member" "editor_role" {
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:service-${var.GCP_PROJECT_ID}@gcp-sa-pubsub.iam.gserviceaccount.com"
  project = var.GCP_PROJECT_ID
}
