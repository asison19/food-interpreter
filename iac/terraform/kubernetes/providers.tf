provider "google" {
  project = var.GCP_PROJECT_ID
  region  = var.GCP_PROJECT_REGION
}

provider "kubernetes" {
  host = "https://${ google_container_cluster.primary.endpoint }"
  token = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth[0].cluster_ca_certificate)

  ignore_annotations = [
    "^autopilot\\.gke\\.io\\/.*",
    "^cloud\\.google\\.com\\/.*"
  ]
}
