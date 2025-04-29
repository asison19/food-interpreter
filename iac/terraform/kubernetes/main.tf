#
# Creating the Kubernetes workload
#

data "google_project" "project" {
  project_id = var.GCP_PROJECT_ID
}

data "google_client_config" "default" {}

data "google_container_cluster" "primary" {
  name     = "primary-gke-cluster" # TODO variableize and use from kubernetes-cluster project.
  location = "us-central1-a" # TODO same as above.
}

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
          # TODO variableize this
          # TODO semantic versioning only happens on master branch. Make it such that other branches have branch name in it for the version to differentiate and have it possible for other branches to push to push to GCP.
          image = "us-central1-docker.pkg.dev/food-interpreter/food-interpreter-repository/food-interpreter:${ vars.food_interpreter_image_version }"
          name  = "food-interpreter-container"
          port {
            container_port = 8080
            name           = "food-interpreter-svc"
          }
          #security_context {
          #  allow_privilege_escalation = false
          #  privileged                 = false
          #  read_only_root_filesystem  = false

          #  capabilities {
          #    add  = []
          #    drop = ["NET_RAW"]
          #  }
          #}
          #liveness_probe {
          #  # TODO
          #}
        }
      }
    }
  }
}

resource "kubernetes_service" "food-interpreter" {
  metadata {
    name = "food-interpreter-loadbalancer"
    annotations = {
      "networking.gke.io/load-balancer-type" = "Internal" # Remove to create an external loadbalancer
    }
  }

  spec {
    selector = {
      app = kubernetes_deployment.food-interpreter.spec[0].selector[0].match_labels.app
    }

    ip_family_policy = "RequireDualStack"

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

#output "endpoint" {
#  value = data.google_container_cluster.my_cluster.endpoint
#}
#
#output "instance_group_urls" {
#  value = data.google_container_cluster.my_cluster.node_pool[0].instance_group_urls
#}
#
#output "node_config" {
#  value = data.google_container_cluster.my_cluster.node_config
#}
#
#output "node_pools" {
#  value = data.google_container_cluster.my_cluster.node_pool
#}