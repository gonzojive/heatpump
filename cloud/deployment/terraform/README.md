Code in here is largely based on
https://cloud.google.com/architecture/managing-infrastructure-as-codeand
https://cloud.google.com/run/docs/deploying#terraform


One time setup required to get cloud build working is mostly based on
https://cloud.google.com/architecture/managing-infrastructure-as-codeand.

```shell
PROJECT_ID=$(gcloud config get-value project)
CLOUDBUILD_SA="$(gcloud projects describe $PROJECT_ID     --format 'value(projectNumber)')@cloudbuild.gserviceaccount.com"

gsutil mb gs://${PROJECT_ID}-tfstate
gsutil versioning set on gs://${PROJECT_ID}-tfstate

gcloud services enable cloudbuild.googleapis.com compute.googleapis.com

# I had to enable one other API manually for Terraform to work. I forget which.

gcloud projects add-iam-policy-binding $PROJECT_ID     --member serviceAccount:$CLOUDBUILD_SA --role roles/editor --role roles/iam.securityAdmin
```

Then push to the `infra-dev` branch for the changes to take effect.
