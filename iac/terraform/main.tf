// TODO GKE

resource "google_pubsub_topic" "lexer" {
  name = "food-interpreter-lexer"

  labels = {
    application = "food-interpreter"
  }

  message_retention_duration = "86600s"
}
