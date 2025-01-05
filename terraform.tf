terraform {
  backend "gcs" {
    bucket  = var.GCP_STATE_BUCKET
    prefix  = "terraform/state"
  }
}
