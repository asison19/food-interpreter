terraform {
  backend "gcs" {
    bucket  = "95fa2a619dd6f929-terraform-remote-backend"
    prefix  = "terraform/state"
  }
}
