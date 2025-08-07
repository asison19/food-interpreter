resource "google_cloud_run_v2_service" "gateway" {
  name     = "gateway-cloud-run"
  location = var.GCP_PROJECT_REGION
  deletion_protection = false

  # TODO, still without latest tag it will have a change in terraform plan on all runs.
  template {
    timeout = "10s"
    containers {
      image = "${ var.GCP_PROJECT_REGION }-docker.pkg.dev/${ var.GCP_PROJECT_ID }/${ google_artifact_registry_repository.food-interpreter-repository.name }/food-interpreter-gateway"
      resources {
        limits = {
          memory = "1024Mi"
          cpu    = "1000m"
        }
      }
    }
    service_account = google_service_account.gateway_cloud_run.email
  }
  depends_on = [
    google_project_iam_member.gateway_act_as
  ]
}

resource "google_service_account" "gateway_cloud_run" {
  account_id   = "gateway-cloud-run"
  display_name = "Gateway Cloud Run Service Account"
}

resource "google_cloud_run_service_iam_binding" "gateway_servicesinvoker" {
  location = google_cloud_run_v2_service.gateway.location
  service  = google_cloud_run_v2_service.gateway.name
  role     = "roles/run.invoker"
  members  = ["serviceAccount:${ google_service_account.gateway_cloud_run.email }"]
}

resource "google_pubsub_topic_iam_binding" "binding" {
  project = google_pubsub_topic.interpreter.project
  topic   = google_pubsub_topic.interpreter.id
  role    = "roles/pubsub.publisher"
  members = ["serviceAccount:${ google_service_account.gateway_cloud_run.email }"]
}

resource "google_project_iam_member" "gateway_act_as" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${ google_service_account.gateway_cloud_run.email }"
}
