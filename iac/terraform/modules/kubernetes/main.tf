data "google_project" "project" {
  project_id = var.GCP_PROJECT_ID
}

data "google_client_config" "default" {}

#
# Creating the Kubernetes cluster
#

# TODO what was I doing with this?
resource "google_service_account" "default" {
  account_id   = "gke-service-account"
  display_name = "Service Account"
}

# TODO rename and move these resources? Can't move, immutable name
resource "google_container_cluster" "primary" {
  name     = "food-interpreter-cluster"
  location = var.GCP_PROJECT_REGION

  #remove_default_node_pool = true
  initial_node_count       = 1
  enable_autopilot         = true # TODO keep autopilot?

  deletion_protection = false

  #network    = google_compute_network.default.id
  #subnetwork = google_compute_subnetwork.default.id
}

#
# Creating the Kubernetes workload
#

resource "kubernetes_deployment" "food_interpreter" {
  metadata {
    name = "food-interpreter"
  }

  spec {
    selector {
      match_labels =  {
        app = "food-interpreter"
      }
    }

    template {
      metadata {
        labels = {
          app = "food-interpreter"
        }
      }
      spec {
        container {
          # TODO this will have the latest at the time. How to ensure it has the latest all the time? Try similar way as cloud run? What about rollbacks? ArgoCD?
          # TODO is there a way to change the deployment without going through terraform if the image in the GAR gets updated?
          image = var.deployment_image
          name  = "food-interpreter-container"
          port {
            container_port = 8080
            name           = "food-int"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "food_interpreter" {
  metadata {
    name = "food-interpreter-loadbalancer"
    annotations = {
      "networking.gke.io/load-balancer-type" = "Internal"
    }
  }

  spec {
    selector = {
      app = kubernetes_deployment.food-interpreter.spec[0].selector[0].match_labels.app
    }

    port {
      port        = 80
      target_port = kubernetes_deployment.food-interpreter.spec[0].template[0].spec[0].container[0].port[0].name
    }

    type = "LoadBalancer"
  }

  depends_on = [time_sleep.wait_service_cleanup]
}

# Provide time for Service cleanup
resource "time_sleep" "wait_service_cleanup" {
  depends_on = [ google_container_cluster.primary ]

  destroy_duration = "180s"
}
