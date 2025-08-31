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
      older_than = "3d"
    }
  }
  cleanup_policies {
    id     = "keep-amount"
    action = "KEEP"
    most_recent_versions {
      keep_count = 2
    }
  }
}

resource "google_project_iam_binding" "logwriter" {
  project = var.GCP_PROJECT_ID
  role    = "roles/logging.logWriter"
  members = [
    "serviceAccount:${ google_service_account.interpreter_cloud_run.email}",
    "serviceAccount:${ google_service_account.gateway_cloud_run.email}",
    "serviceAccount:${ google_service_account.interpreter_grpc_cloud_run.email }"
  ]
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
