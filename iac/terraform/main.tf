// TODO GKE

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

resource "google_cloud_run_v2_job" "lexer" {
  name     = "lexer-cloud-run"
  location = "us-central1"
  deletion_protection = false

  template {
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
}

resource "google_service_account" "lexer_cloud_run" {
  account_id   = "lexer-cloud-run"
  display_name = "Lexer Cloud Run Service Account"
}

data "google_iam_policy" "lexer" {
  binding {
    role = "roles/viewer" # TODO give it the roles it needs later on.
    members = [
      google_service_account.lexer_cloud_run.email,
    ]
  }
}

resource "google_cloud_run_v2_job_iam_policy" "lexer" {
  project     = google_cloud_run_v2_job.lexer.project
  location    = google_cloud_run_v2_job.lexer.location
  name        = google_cloud_run_v2_job.lexer.name
  policy_data = data.google_iam_policy.lexer.policy_data
}
