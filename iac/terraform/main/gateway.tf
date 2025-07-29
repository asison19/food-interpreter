resource "google_cloud_run_v2_service" "gateway" {
  name     = "gateway-cloud-run"
  location = var.GCP_PROJECT_REGION
  deletion_protection = false

  template {
    timeout = "10s"
    containers {
      image = "${ var.GCP_PROJECT_REGION }-docker.pkg.dev/${ var.GCP_PROJECT_ID }/${ google_artifact_registry_repository.food-interpreter-repository.name }/food-interpreter-gateway:latest"
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

#
# Used for GitHub Actions API Testing
#
resource "google_project_iam_member" "gateway_token_creator" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:${ google_service_account.gateway_cloud_run.email }"
}

resource "google_service_account_key" "gateway_cloud_run" {
  service_account_id = google_service_account.gateway_cloud_run.name
  public_key_type    = "TYPE_X509_PEM_FILE"

  keepers = {
    rotation_time = time_rotating.gateway_cloud_run.rotation_rfc3339
  }
}

resource "time_rotating" "gateway_cloud_run" {
  rotation_days = 7
}
