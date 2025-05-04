terraform {
  backend "gcs" {
    bucket  = "z8nfq4kvf1x2d7pr-terraform-remote-backend"
    prefix  = "terraform/state"
  }
}
