#
# Creating the Kubernetes cluster
#

resource "google_service_account" "default" {
  account_id   = "gke-service-account"
  display_name = "Service Account"
}

# TODO rename and move these resources? Can't move, immutable name
resource "google_container_cluster" "primary" {
  #name     = "food-interpreter" # TODO destroy and rename? name is immutable
  name     = "primary-gke-cluster"
  location = var.GCP_PROJECT_REGION

  #remove_default_node_pool = true
  initial_node_count       = 1
  enable_autopilot         = true # TODO keep autopilot?

  deletion_protection = false

  #network    = google_compute_network.default.id
  #subnetwork = google_compute_subnetwork.default.id
}
