To integrate with Google Home, we need a cloud service.


## Command farm

Server

```shell
bazel run --run_under="cd $PWD &&" //cmd/queueserver -- --alsologtostderr --alsologtostderr --server-cert secrets/server-cert.pem --server-private-key secrets/server-key.pem  --client-ca-cert secrets/client-signer-cert-authority-cert.pem
```

Client

```shell
bazel run --run_under="cd $PWD &&" //cmd/cloud-listener -- --alsologtostderr --client-cert secrets/client-cert.pem --client-private-key secrets/client-key.pem --command-queue-addr "0.0.0.0:8083" --server-ca-cert secrets/server-ca-cert.pem
```

Debug: Send a message

```shell
gcloud pubsub topics publish iot-commands --message "hello there again" --attribute user-id=redshouse
```
