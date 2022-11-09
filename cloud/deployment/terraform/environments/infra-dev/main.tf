locals {
  gcp_location = "us-west4"
}

provider "google" {
  project = var.project
}

# Services (APIs) to enable
resource "google_project_service" "artifactregistry" {
  project = var.project
  service = "artifactregistry.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

resource "google_project_service" "cloudrun" {
  project = var.project
  service = "run.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

# Artifact storage of container images.

resource "google_artifact_registry_repository" "my-repo" {
  location      = local.gcp_location
  repository_id = "project-images"
  description   = "Images pushed through bazel rules."
  format        = "DOCKER"
}

# resource "google_cloud_run_service" "google_actions_http_endpoint" {
#   name     = "google-actions-http-endpoint"
#   location = locals.gcp_location

#   metadata {
#     annotations = {
#       "run.googleapis.com/client-name" = "terraform"
#       "run.googleapis.com/ingress"     = "all"
#     }
#   }

#   template {
#     spec {
#       containers {
#         image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/reverse-proxy-image@sha256:981653a39aa265d670c04693f0660a6646068f32933ad26c0885637120a87cfc"
#       }
#     }
#   }
# }

# resource "google_cloud_run_service" "command_queue_service" {
#   name     = "command-queue-service"
#   location = locals.gcp_location

#   metadata {
#     annotations = {
#       "run.googleapis.com/client-name" = "terraform"
#       "run.googleapis.com/ingress"     = "all"
#     }
#   }

#   template {
#     spec {
#       containers {
#         image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/command-queue-service-image@sha256:50311f03fdc452c56055793e871a2e2ff96cf793fedb148043884e6858899ae2"
#         # Enable HTTP/2 so gRPC works.
#         # https://cloud.google.com/run/docs/configuring/http2
#         ports {
#           name           = "h2c"
#           container_port = 8083
#         }
#       }
#     }

#     metadata {
#       annotations = {
#         "autoscaling.knative.dev/maxScale" = "1"
#       }
#     }
#   }
# }

# resource "google_cloud_run_service" "state_service" {
#   name     = "state-service"
#   location = locals.gcp_location

#   metadata {
#     annotations = {
#       "run.googleapis.com/client-name" = "terraform"
#       "run.googleapis.com/ingress"     = "all"
#     }
#   }

#   template {
#     spec {
#       containers {
#         image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/stateservice-image@sha256:3b80c5732a7f350ef938b88129e1cceb45745add2cbae73bc817d35288209df4"
#         # Enable HTTP/2 so gRPC works.
#         # https://cloud.google.com/run/docs/configuring/http2
#         ports {
#           name           = "h2c"
#           container_port = 8089
#         }
#       }
#     }

#     metadata {
#       annotations = {
#         "autoscaling.knative.dev/maxScale" = "1"
#       }
#     }
#   }
# }

data "google_iam_policy" "noauth" {
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

# resource "google_cloud_run_service_iam_policy" "noauth" {
#   location = google_cloud_run_service.google_actions_http_endpoint.location
#   project  = google_cloud_run_service.google_actions_http_endpoint.project
#   service  = google_cloud_run_service.google_actions_http_endpoint.name

#   policy_data = data.google_iam_policy.noauth.policy_data
# }

# resource "google_cloud_run_service_iam_policy" "noauth_state_service" {
#   location = google_cloud_run_service.state_service.location
#   project  = google_cloud_run_service.state_service.project
#   service  = google_cloud_run_service.state_service.name

#   policy_data = data.google_iam_policy.noauth.policy_data
# }

# resource "google_cloud_run_service_iam_policy" "noauth_command_queue_service" {
#   location = google_cloud_run_service.command_queue_service.location
#   project  = google_cloud_run_service.command_queue_service.project
#   service  = google_cloud_run_service.command_queue_service.name

#   policy_data = data.google_iam_policy.noauth.policy_data
# }

module "pubsub_iot_commands" {
  source  = "terraform-google-modules/pubsub/google"
  version = "~> 4.0.1"

  topic                            = "iot-commands"
  project_id                       = var.project
  topic_message_retention_duration = "432000s" // 5 days
  push_subscriptions               = []
  pull_subscriptions = [
    {
      name                         = "queueserver_pull_1" // required
      ack_deadline_seconds         = 60                   // optional
      max_delivery_attempts        = 5                    // optional
      maximum_backoff              = "600s"               // optional
      minimum_backoff              = "10s"                // optional
      enable_exactly_once_delivery = true                 // optional
    }
  ]
}
