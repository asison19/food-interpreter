resource "google_pubsub_topic" "example" {
  name = "food-interpreter-lexer"

  labels = {
    application = "food-interpreter"
  }

  message_retention_duration = "86600s"
}
