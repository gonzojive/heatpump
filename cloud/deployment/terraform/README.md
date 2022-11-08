Example based on https://cloud.google.com/run/docs/deploying#terraform

To run the example on Cloud Run:

```shell
bazel run //cmd/heatpump-oauth-server:push_go_image
```

```shell
cd cloud/google/config
terraform init
terraform plan
terraform apply
```