locals {
  gcp_location = "us-west4"

  # We keep a map of (shorthand image name => pinned image reference) in a file
  # in this directory.
  #
  # To regenerate the file based on image-versions-config.json, run
  /*
   bazel run --run_under="cd $PWD &&" //cmd/cloud/update-image-versions -- --alsologtostderr --input "cloud/deployment/terraform/environments/infra-dev/image-versions.json" --output "cloud/deployment/terraform/environments/infra-dev/image-versions.json"
  */
  image_versions = jsondecode(file("${path.module}/image-versions.json"))["images"]
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

resource "google_project_service" "iam" {
  project = var.project
  service = "iam.googleapis.com"

  disable_dependent_services = true
}

# Create service accounts (robots) and establish their permissions.
resource "google_service_account" "google_actions_http_job" {
  project      = var.project
  account_id   = "google-actions-http-job"
  display_name = "Cloud Run user for Google Actions HTTP responder"
}

# App Engine Admin API.
resource "google_project_service" "appengine" {
  project = var.project
  service = "appengine.googleapis.com"

  disable_dependent_services = true
}

# Firestore database creation
# https://cloud.google.com/firestore/docs/solutions/automate-database-create#firestoretf
resource "google_app_engine_application" "app" {
  project       = var.project
  location_id   = local.gcp_location
  database_type = "CLOUD_FIRESTORE"
}

# Artifact storage of container images.

resource "google_artifact_registry_repository" "my-repo" {
  location      = local.gcp_location
  repository_id = "project-images"
  description   = "Images pushed through bazel rules."
  format        = "DOCKER"
}

resource "google_cloud_run_service" "google_actions_http_endpoint" {
  name     = "google-actions-http-endpoint"
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
        # To debug:
        # bazel run //cmd/reverse-proxy:push-image
        #
        # To start the real server:
        # bazel run //cmd/cloud/http-endpoint
        image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/http-endpoint@sha256:9bf8231e968aa47adf49a7f1f9b3a7b22024db8709496510dfd746ee7c122f9e"
        args = [
          "--alsologtostderr",
          "--state-service-addr", "state-service-6mavwogfvq-wn.a.run.app:443",
        ]
      }
      service_account_name = google_service_account.google_actions_http_job.email
    }
  }
}

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
        image = local.image_versions["queue-service"].resolved
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
        image = "us-west4-docker.pkg.dev/heatpump-dev/project-images/stateservice-image:main"
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

resource "google_cloud_run_service_iam_policy" "noauth_google_actions_http_endpoint" {
  location = google_cloud_run_service.google_actions_http_endpoint.location
  project  = google_cloud_run_service.google_actions_http_endpoint.project
  service  = google_cloud_run_service.google_actions_http_endpoint.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

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
