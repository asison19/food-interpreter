resource "google_cloud_run_v2_service" "interpreter_grpc" {
  name     = "interpreter-grpc-cloud-run"
  location = var.GCP_PROJECT_REGION
  deletion_protection = false

  template {
    timeout = "10s"
    containers {
      image = "${ var.GCP_PROJECT_REGION }-docker.pkg.dev/${ var.GCP_PROJECT_ID }/${ google_artifact_registry_repository.food-interpreter-repository.name }/food-interpreter-interpreter-grpc"
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
