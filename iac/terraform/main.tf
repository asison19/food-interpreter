// TODO GKE

# TODO look into automatically deleting images.
resource "google_artifact_registry_repository" "food-interpreter-repository" {
  location      = "us-central1"
  repository_id = "food-interpreter-repository"
  description   = "Food Interpreter docker repository"
  format        = "DOCKER"
}
