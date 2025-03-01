// TODO GKE

data "google_project" "project" {
  id   = var.GCP_PROJECT_ID
}

# TODO look into automatically deleting images.
resource "google_artifact_registry_repository" "food-interpreter-repository" {
  location      = "us-central1"
  repository_id = "food-interpreter-repository"
  description   = "Food Interpreter docker repository"
  format        = "DOCKER"

  cleanup_policies {
    id = "keep-amount"
    action = "KEEP"
    most_recent_versions {
      keep_count = 5
    }
  }
}

resource "google_cloud_run_v2_service" "lexer" {
  name     = "lexer-cloud-run"
  location = "us-central1"
  deletion_protection = false

  template {
    timeout = "10s"
    containers {
      image = "us-central1-docker.pkg.dev/${ var.GCP_PROJECT_ID }/${ google_artifact_registry_repository.food-interpreter-repository.name }/food-interpreter:latest"
      resources {
        limits = {
          memory = "1024Mi"
        }
      }
    }
  }
}

resource "google_service_account" "lexer_cloud_run" {
  account_id   = "lexer-cloud-run"
  display_name = "Lexer Cloud Run Service Account"
}

resource "google_cloud_run_service_iam_binding" "binding" {
  location = google_cloud_run_v2_service.lexer.location
  service  = google_cloud_run_v2_service.lexer.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${google_service_account.lexer_cloud_run.email}"]
}

resource "google_project_service_identity" "pubsub_agent" {
  provider = google-beta
  project  = data.google_project.project.project_id
  service  = "pubsub.googleapis.com"
}

resource "google_project_iam_binding" "project_token_creator" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  members = ["serviceAccount:${google_project_service_identity.pubsub_agent.email}"]
}

# TODO do I still need this?
data "google_iam_policy" "lexer" {
  binding {
    role = "roles/viewer" # TODO give it the roles it needs later on.
    members = [
      "serviceAccount:${ google_service_account.lexer_cloud_run.email }",
    ]
  }
}

resource "google_cloud_run_v2_service_iam_policy" "lexer" {
  project     = google_cloud_run_v2_service.lexer.project
  location    = google_cloud_run_v2_service.lexer.location
  name        = google_cloud_run_v2_service.lexer.name
  policy_data = data.google_iam_policy.lexer.policy_data
}

# TODO schema?
resource "google_pubsub_topic" "lexer" {
  name = "lexer-topic"

  labels = {
    service = "lexer"
  }

  message_retention_duration = "86000s"
}

resource "google_pubsub_subscription" "example" {
  name  = "lexer-subscription"
  topic = google_pubsub_topic.lexer.id

  ack_deadline_seconds = 20

  labels = {
    service = "lexer"
  }

  push_config {
    push_endpoint = google_cloud_run_v2_service.lexer.uri

    oidc_token {
      service_account_email = google_service_account.lexer_cloud_run.email
    }

    attributes = {
      x-goog-version = "v1"
    }
  }
  depends_on = [ google_cloud_run_v2_service.lexer ]
}
