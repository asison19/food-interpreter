data "google_project" "project" {
  project_id = var.GCP_PROJECT_ID
}

data "google_client_config" "default" {}

#
# Creating the Kubernetes cluster
#
# TODO what was I doing with this?
resource "google_service_account" "default" {
  account_id   = "gke-service-account"
  display_name = "Service Account"
}

resource "google_container_cluster" "primary" {
  name     = "food-interpreter-cluster"
  location = var.GCP_PROJECT_REGION

  #remove_default_node_pool = true
  initial_node_count       = 1
  enable_autopilot         = true # TODO keep autopilot?

  deletion_protection = false

  #network    = google_compute_network.default.id
  #subnetwork = google_compute_subnetwork.default.id
}
