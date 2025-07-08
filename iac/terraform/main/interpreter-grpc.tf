resource "google_cloud_run_v2_service" "interpreter_grpc" {
  name     = "interpreter-grpc-cloud-run"
  location = var.GCP_PROJECT_REGION
  deletion_protection = false

  template {
    timeout = "10s"
    containers {
      image = "${ var.GCP_PROJECT_REGION }-docker.pkg.dev/${ var.GCP_PROJECT_ID }/${ google_artifact_registry_repository.food-interpreter-repository.name }/food-interpreter-interpreter-grpc:latest"
      ports {
        name           = "h2c"
        container_port = 8080
      }
      resources {
        limits = {
          memory = "1024Mi"
          cpu    = "1000m"
        }
      }
    }
  }
}

resource "google_service_account" "interpreter_grpc_cloud_run" {
  account_id   = "interpreter-grpc-cloud-run"
  display_name = "Interpreter gRPC Cloud Run Service Account"
}

resource "google_cloud_run_service_iam_binding" "interpreter_grpc_servicesinvoker" {
  location = google_cloud_run_v2_service.interpreter_grpc.location
  service  = google_cloud_run_v2_service.interpreter_grpc.name
  role     = "roles/run.invoker"
  members  = [
    "serviceAccount:${ google_service_account.interpreter_grpc_cloud_run.email }",
    "serviceAccount:${ google_service_account.gateway_cloud_run.email }"
  ]
}

#
# Service account key for gateway to interpreter-grpc
#

# note this requires the terraform to be run regularly
resource "time_rotating" "interpreter_grpc" {
  rotation_days = 30
}

resource "google_service_account_key" "interpeter_grpc" {
  service_account_id = google_service_account.interpreter_grpc_cloud_run.name

  keepers = {
    rotation_time = time_rotating.interpreter_grpc.rotation_rfc3339
  }
}
