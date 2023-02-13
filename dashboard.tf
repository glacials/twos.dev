resource "kubernetes_deployment" "dashboard" {
  metadata {
    name = "dashboard"
    labels = {
      App = "Dashboard"
    }
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        App = "Dashboard"
      }
    }
    template {
      metadata {
        labels = {
          App = "Dashboard"
        }
      }
      spec {
        container {
          image = "kubernetesui/dashboard"
          name  = "dsahboard"
          port {
            container_port = 80
          }
          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "dashboard" {
  metadata {
    name = "dashboard"
  }
  spec {
    selector = {
      App = kubernetes_deployment.dashboard.spec.0.template.0.metadata[0].labels.App
    }
    port {
      port        = 80
      target_port = 80
    }
    type = "LoadBalancer"
  }
}

resource "cloudflare_record" "dashboard" {
  zone_id = var.zone_id
  name    = "dashboard"
  value   = kubernetes_service.dashboard.status.0.load_balancer.0.ingress.0.hostname
  type    = "CNAME"
  proxied = true
}
