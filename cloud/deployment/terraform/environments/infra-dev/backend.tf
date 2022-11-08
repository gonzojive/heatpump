terraform {
  backend "gcs" {
    bucket = "heatpump-dev-tfstate"
    prefix = "env/dev"
  }
}
