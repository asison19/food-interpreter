output "interpreter_pubsub_topic_id" {
    value = google_pubsub_topic.interpreter.id
}

output "interpreter_grpc_cloud_run_uri" {
    value = google_cloud_run_v2_service.interpreter_grpc.uri
}

output "gateway_cloud_run_uri" {
    value = google_cloud_run_v2_service.gateway.uri
}

output "gateway_service_account" {
    value = google_service_account.gateway_cloud_run.email
}
