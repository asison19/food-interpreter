output "kubernetes_cluster_id" {
    value       = google_container_cluster.primary.id
    description = "The ID of the GKE Kubernetes cluster."
}

output "deployment_id" {
    value       = kubernetes_deployment.food_interpreter.id
    description = "The ID of the Kubernetes deployment."
}

output "service_id" {
    value       = kubernetes_service.food_interpreter.id
    description = "The ID of the Kubernetes service."
}
