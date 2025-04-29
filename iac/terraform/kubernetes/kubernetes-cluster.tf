#
# Creating the Kubernetes cluster
#

resource "google_service_account" "default" {
  account_id   = "gke-service-account"
  display_name = "Service Account"
}

resource "google_container_cluster" "primary" {
  name     = "primary-gke-cluster"
  location = "us-central1-a" # TODO variableize

  remove_default_node_pool = true
  initial_node_count       = 1
  enable_autopilot         = true

  deletion_protection = false

  #network    = google_compute_network.default.id
  #subnetwork = google_compute_subnetwork.default.id
}

resource "google_container_node_pool" "primary" {
  name       = "primary-node-pool"
  location   = "us-central1-a"
  cluster    = google_container_cluster.primary.name
  node_count = 1

  node_config {
    spot         = true
    machine_type = "e2-micro" # TODO create custom machine type?

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}
