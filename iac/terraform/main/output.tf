output "interpreter_pubsub_topic_id" {
    value = google_pubsub_topic.interpreter.id
}

output "interpreter_grpc_cloud_run_uri" {
    value = google_cloud_run_v2_service.interpreter_grpc.uri
}

# This is only populated when creating a new key.
# Is this necessary? The ID token should be made in interpretRequestWithAuth().
output "interpreter_grpc_sa_key" {
    value = google_service_account_key.interpeter_grpc.private_key
    sensitive = true
}
