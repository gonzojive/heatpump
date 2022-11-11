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

resource "google_project_service" "secretmanager" {
  project = var.project
  service = "secretmanager.googleapis.com"

  timeouts {
    create = "30m"
    update = "40m"
  }

  disable_dependent_services = true
}

resource "google_project_service" "firestore" {
  project = var.project
  service = "firestore.googleapis.com"

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
#   location = local.gcp_location

#   metadata {
#     annotations = {
#       "run.googleapis.com/client-name" = "terraform"
#       "run.googleapis.com/ingress"     = "all"
#     }
#   }

#   template {
#     spec {
#       containers {
#         # bazel run //cmd/queueserver:push-image
#         image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/reverse-proxy-image@sha256:981653a39aa265d670c04693f0660a6646068f32933ad26c0885637120a87cfc"
#       }
#     }
#   }
# }

resource "google_cloud_run_service" "command_queue_service" {
  name     = "command-queue-service"
  location = local.gcp_location

  metadata {
    annotations = {
      "run.googleapis.com/client-name" = "terraform"
      "run.googleapis.com/ingress"     = "all"
    }
  }

  template {
    spec {
      containers {
        # bazel run //cmd/queueserver:push-image
        image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/command-queue-service-image@sha256:507394c79022e2cb9e31f22ec1bbf2d0c9506e6d4ba74e8dc9c15d0e64d45302"
        # Enable HTTP/2 so gRPC works.
        # https://cloud.google.com/run/docs/configuring/http2
        ports {
          name           = "h2c"
          container_port = 8083
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
      }
    }
  }
}

resource "google_cloud_run_service" "state_service" {
  name     = "state-service"
  location = local.gcp_location

  metadata {
    annotations = {
      "run.googleapis.com/client-name" = "terraform"
      "run.googleapis.com/ingress"     = "all"
    }
  }

  template {
    spec {
      containers {
        # bazel run //cmd/stateservice:push-image
        image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/stateservice-image@sha256:0a480e08794147437170dd0c0af7aa92e953bec99d852e943ceb45b3cfd3b1a9"
        # Enable HTTP/2 so gRPC works.
        # https://cloud.google.com/run/docs/configuring/http2
        ports {
          name           = "h2c"
          container_port = 8089
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
      }
    }
  }
}

# resource "google_cloud_run_service" "auth_service" {
#   name     = "auth-service"
#   location = local.gcp_location

#   metadata {
#     annotations = {
#       "run.googleapis.com/client-name" = "terraform"
#       "run.googleapis.com/ingress"     = "all"
#     }
#   }

#   template {
#     spec {
#       containers {
#         # bazel run //cmd/authservice:push-image
#         image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/authservice-image@sha256:7b8df3707487ca9f7fafd7d2b9f1514ac8fd71b55b33e311cc5dd44ae283068b"
#         # Enable HTTP/2 so gRPC works.
#         # https://cloud.google.com/run/docs/configuring/http2
#         ports {
#           name           = "h2c"
#           container_port = 8092
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

resource "google_cloud_run_service_iam_policy" "noauth_state_service" {
  location = google_cloud_run_service.state_service.location
  project  = google_cloud_run_service.state_service.project
  service  = google_cloud_run_service.state_service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

# resource "google_cloud_run_service_iam_policy" "noauth_auth_service" {
#   location = google_cloud_run_service.auth_service.location
#   project  = google_cloud_run_service.auth_service.project
#   service  = google_cloud_run_service.auth_service.name

#   policy_data = data.google_iam_policy.noauth.policy_data
# }

resource "google_cloud_run_service_iam_policy" "noauth_command_queue_service" {
  location = google_cloud_run_service.command_queue_service.location
  project  = google_cloud_run_service.command_queue_service.project
  service  = google_cloud_run_service.command_queue_service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

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
