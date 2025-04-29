#
# Creating the Kubernetes cluster
#

# TODO can I remove this?
resource "google_service_account" "default" {
  account_id   = "gke-service-account"
  display_name = "Service Account"
}

resource "google_compute_network" "default" {
  name = "primary-network"

  auto_create_subnetworks  = false
  enable_ula_internal_ipv6 = true
}

resource "google_compute_subnetwork" "default" {
  name = "primary-subnetwork"

  ip_cidr_range = "10.0.0.0/16"
  region        = "us-central1"

  stack_type       = "IPV4_IPV6"
  ipv6_access_type = "INTERNAL" # Change to "EXTERNAL" if creating an external loadbalancer

  network = google_compute_network.default.id
  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "192.168.0.0/24"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "192.168.1.0/24"
  }
}


# TODO rename and move these resources? Can't move, immutable name
resource "google_container_cluster" "primary" {
  #name     = "food-interpreter" # TODO destroy and rename? name is immutable
  name     = "primary-gke-cluster"
  location = var.GCP_PROJECT_REGION

  #remove_default_node_pool = true
  initial_node_count       = 1
  enable_autopilot         = true # TODO keep autopilot?
  enable_l4_ilb_subsetting = true

  deletion_protection = false

  network    = google_compute_network.default.id
  subnetwork = google_compute_subnetwork.default.id

  ip_allocation_policy {
    stack_type                    = "IPV4_IPV6"
    services_secondary_range_name = google_compute_subnetwork.default.secondary_ip_range[0].range_name
    cluster_secondary_range_name  = google_compute_subnetwork.default.secondary_ip_range[1].range_name
  }
}
