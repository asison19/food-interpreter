#
# Creating the Kubernetes workload
#

data "google_project" "project" {
  project_id = var.GCP_PROJECT_ID
}

data "google_client_config" "default" {}

# TODO naming, hyphen, move the resource.
resource "kubernetes_deployment" "food-interpreter" {
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
          image = "us-central1-docker.pkg.dev/food-interpreter/food-interpreter-repository/food-interpreter:${ var.FOOD_INTERPRETER_IMAGE_VERSION }"
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

# TODO naming, hyphen, move the resource.
resource "kubernetes_service" "food-interpreter" {
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
