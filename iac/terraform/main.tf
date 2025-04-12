// TODO GKE

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

resource "google_cloud_run_v2_service" "lexer" {
  name     = "lexer-cloud-run"
  location = var.GCP_PROJECT_REGION
  deletion_protection = false

  template {
    timeout = "10s"
    containers {
      image = "${ var.GCP_PROJECT_REGION }-docker.pkg.dev/${ var.GCP_PROJECT_ID }/${ google_artifact_registry_repository.food-interpreter-repository.name }/food-interpreter:latest"
      resources {
        limits = {
          memory = "1024Mi"
          cpu    = "1000m"
        }
      }
    }
  }
}

resource "google_service_account" "lexer_cloud_run" {
  account_id   = "lexer-cloud-run"
  display_name = "Lexer Cloud Run Service Account"
}

resource "google_cloud_run_service_iam_binding" "lexer_servicesinvoker" {
  location = google_cloud_run_v2_service.lexer.location
  service  = google_cloud_run_v2_service.lexer.name
  role     = "roles/run.servicesInvoker"
  members  = ["serviceAccount:${ google_service_account.lexer_cloud_run.email }"]
}

resource "google_project_iam_binding" "lexer_logwriter" {
  project = var.GCP_PROJECT_ID
  role    = "roles/logging.logWriter"
  members = ["serviceAccount:${ google_service_account.lexer_cloud_run.email }"]
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
resource "google_pubsub_topic" "lexer" {
  name = "lexer-topic"

  labels = {
    service = "lexer"
  }

  message_retention_duration = "86000s"
}

resource "google_pubsub_subscription" "lexer" {
  name                 = "lexer-subscription"
  topic                = google_pubsub_topic.lexer.id
  ack_deadline_seconds = 20

  labels = {
    service = "lexer"
  }

  push_config {
    push_endpoint = "${ google_cloud_run_v2_service.lexer.uri }/lexer"
    oidc_token {
      service_account_email = google_service_account.lexer_cloud_run.email
    }
    attributes = {
      x-goog-version = "v1"
    }
  }

  dead_letter_policy {
    dead_letter_topic = google_pubsub_topic.lexer-dlq.id
    max_delivery_attempts = 5
  }

  depends_on = [ google_cloud_run_v2_service.lexer ]
}

resource "google_pubsub_topic" "lexer-dlq" {
  name = "lexer-topic-dlq"

  labels = {
    service = "lexer"
  }

  message_retention_duration = "604800s" # 7 days
}

resource "google_pubsub_subscription" "lexer-dlq" {
  name                 = "lexer-subscription-dlq"
  topic                = google_pubsub_topic.lexer.id
  ack_deadline_seconds = 20

  labels = {
    service = "lexer"
  }
}
