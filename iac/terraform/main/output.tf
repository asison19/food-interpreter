output "interpreter_pubsub_topic_id" {
    value = google_pubsub_topic.interpreter.id
}

output "interpreter_grpc_cloud_run_uri" {
    value = google_cloud_run_v2_service.interpreter_grpc.urls[0]
}
