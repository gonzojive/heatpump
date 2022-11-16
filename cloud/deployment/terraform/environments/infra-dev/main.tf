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


resource "google_project_iam_binding" "firebase_readwrite" {
  project = var.project
  # Yes, datastore.owner, not firestore.
  #
  # See https://cloud.google.com/firestore/docs/security/iam
  role    = "roles/datastore.owner"

  members = [
    google_service_account.state_service_job.member,
    google_service_account.google_actions_http_job.member,
  ]
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

resource "google_service_account" "auth_service_job" {
  project      = var.project
  account_id   = "auth-service-job"
  display_name = "Cloud Run user for AuthService gRPC service that issues tokens to IoT clients."
}

resource "google_service_account" "state_service_job" {
  project      = var.project
  account_id   = "state-service-job"
  display_name = "Cloud Run user for StateService gRPC service that stores device state."
}

resource "google_service_account" "command_queue_service_job" {
  project      = var.project
  account_id   = "command-queue-service-job"
  display_name = "Cloud Run user for StateService gRPC service that stores device state."
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
      service_account_name = google_service_account.google_actions_http_job.email
      containers {
        # To debug:
        # bazel run //cmd/reverse-proxy:push-image
        #
        # To start the real server:
        # bazel run //cmd/cloud/http-endpoint
        image = local.image_versions["google-actions-http-endpoint"].resolved
        args = [
          "--alsologtostderr",
          "--state-service-addr", "state-service-6mavwogfvq-wn.a.run.app:443",
        ]
      }
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
      service_account_name = google_service_account.command_queue_service_job.email

      timeout_seconds = 3600

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
      service_account_name = google_service_account.state_service_job.email
      containers {
        # bazel run //cmd/stateservice:push-image
        image = local.image_versions["state-service"].resolved
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

resource "google_cloud_run_service" "auth_service" {
  name     = "auth-service"
  location = local.gcp_location

  metadata {
    annotations = {
      "run.googleapis.com/client-name" = "terraform"
      "run.googleapis.com/ingress"     = "all"
    }
  }

  template {
    spec {
      service_account_name = google_service_account.auth_service_job.email
      containers {
        # bazel run //cmd/authservice:push-image
        image = local.image_versions["auth-service"].resolved
        # Enable HTTP/2 so gRPC works.
        # https://cloud.google.com/run/docs/configuring/http2
        ports {
          name           = "h2c"
          container_port = 8092
        }
        args = [
          "--grpc-port", 8092,
        ]
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
      }
    }
  }
}

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

resource "google_cloud_run_service_iam_policy" "noauth_auth_service" {
  location = google_cloud_run_service.auth_service.location
  project  = google_cloud_run_service.auth_service.project
  service  = google_cloud_run_service.auth_service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

resource "google_cloud_run_service_iam_policy" "noauth_command_queue_service" {
  location = google_cloud_run_service.command_queue_service.location
  project  = google_cloud_run_service.command_queue_service.project
  service  = google_cloud_run_service.command_queue_service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

# Secret names
#
# Do not store secret information in terraform. Use the command line or the GCP
# UI to store secret material. However, we can use Terraform to create secret
# resources and specify who has access to them. The contents of the secrets are
# stored in "secret versions" in GCP parlance and are specified elsewhere.
resource "google_secret_manager_secret" "device_token_signer_private_rsa" {
  project = var.project
  secret_id = "device-token-signer-private-rsa"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_iam_binding" "device_token_signer_private_rsa" {
  project = var.project
  secret_id = google_secret_manager_secret.device_token_signer_private_rsa.secret_id

  role = "roles/secretmanager.secretAccessor"
  
  members = [
    google_service_account.auth_service_job.member,
  ]
}

resource "google_secret_manager_secret" "google_actions_oauth_client_secret" {
  project = var.project
  secret_id = "google-actions-oauth-client-secret"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_iam_binding" "google_actions_oauth_client_secret" {
  project = var.project
  secret_id = google_secret_manager_secret.google_actions_oauth_client_secret.secret_id

  role = "roles/secretmanager.secretAccessor"
  
  members = [
    google_service_account.google_actions_http_job.member,
  ]
}

resource "google_pubsub_topic" "iot_commands" {
  project = var.project
  name = "iot-commands"
  message_retention_duration = "432000s" // 5 days
}

resource "google_pubsub_subscription" "iot_commands_queue_server" {
  project = var.project
  topic = google_pubsub_topic.iot_commands.name
  name  = "queueserver_pull_1"

  labels = {}

  enable_exactly_once_delivery = true
  enable_message_ordering    = false

  # 20 minutes
  retain_acked_messages      = false
  ack_deadline_seconds = 60
  message_retention_duration = google_pubsub_topic.iot_commands.message_retention_duration

  expiration_policy {
    ttl = "2678400s"
  }
  retry_policy {
    maximum_backoff = "600s"
    minimum_backoff = "10s"
  }
}

resource "google_pubsub_topic_iam_binding" "iot_commands_publisher" {
  project = google_pubsub_topic.iot_commands.project
  topic = google_pubsub_topic.iot_commands.name
  role = "roles/pubsub.publisher"
  members = [
    google_service_account.google_actions_http_job.member,
  ]
}

resource "google_pubsub_subscription_iam_binding" "iot_commands_queue_server" {
  project = google_pubsub_topic.iot_commands.project
  subscription = google_pubsub_subscription.iot_commands_queue_server.name
  role = "roles/pubsub.subscriber"
  members = [
    google_service_account.command_queue_service_job.member,
  ]
}
