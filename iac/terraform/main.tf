// TODO GKE

resource "google_pubsub_topic" "lexer" {
  name = "food-interpreter-lexer"

  labels = {
    application = "food-interpreter"
  }

  message_retention_duration = "86600s"
}

resource "google_artifact_registry_repository" "food-interpreter-repository" {
  location      = "us-central1"
  repository_id = "food-interpreter-repository"
  description   = "Food Interpreter docker repository"
  format        = "DOCKER"
}
