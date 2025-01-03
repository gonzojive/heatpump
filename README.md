# Chiltrix CX34 dashboard software

This project is used to communicate with a hydronic heating system that uses the
[Chiltrix CX34](https://www.chiltrix.com/small-chiller-home.html) air-to-water
heat pump. The components of the system are as follows:

1. A Raspberry Pi.
2. The [Chiltrix CX34](https://www.chiltrix.com/small-chiller-home.html) heat
   pump chiller. The Pi communicates with the chiller using using [serial
   communication](https://en.wikipedia.org/wiki/RS-485), which is used natively
   by the CX34 to talk with its controller per the [wiring
   diagram](https://www.chiltrix.com/documents/CX34-2-wiring-diagram-HIGH-RES.pdf)
   and [ProtoAir FPA-W44
   manual](https://www.chiltrix.com/control-options/Remote-Gateway-BACnet-Guide-rev2.pdf).

3. A USB-to-RS484 dongle that connects the Raspberry Pi to the Chiltrix. This
   device should show up as `/dev/ttyUSB0`.

4. Optional [DS18B20](https://www.adafruit.com/product/381) temperature sensors,
   which can now be bought for about $2 each. Several temperature sensors can be
   connected to single GPIO port on the Pi. I followed [this guide](https://www.circuitbasics.com/raspberry-pi-ds18b20-temperature-sensor-tutorial/#:~:text=The%20DS18B20%20temperature%20sensor%20is,accurate%20and%20take%20measurements%20quickly.) to get up and running and
   then wrote the `github.com/gonzojive/heatpump/tempsensor` Go package to read
   the sensor values.

![screenshot of dashboard as of 2020/12/24](https://raw.githubusercontent.com/gonzojive/heatpump/main/docs/screenshot-2020-12-24.png "screenshot of dashboard as of 2020/12/24")

There is also [work in progress](https://github.com/gonzojive/heatpump/issues/5)
Google Home integration, which allows setting the temperature of the fan coil
units from the Google Home app. You can also say commands like:

* "Set the temperature to 72 degrees"
* "Set the fan speed to low"

![screenshot of dashboard as of
2020/12/24](https://raw.githubusercontent.com/gonzojive/heatpump/main/docs/screenshot-google-home.png
"screenshot of Google Home app for controlling fan coil units")

## Implementation notes

This project is written in Go.

## Setup

### Software installation

Install go on the Raspberry Pi according to [the official
instructions](https://golang.org/doc/install).

Install the collector binary as of the latest version:

```shell
go install github.com/gonzojive/heatpump/cmd/cx34collector@latest github.com/gonzojive/heatpump/cmd/cx34install@latest
```

Install the systemd service so the service runs at boot.

```shell
sudo `which cx34install`
```

Inspect the logs with

```shell
journalctl -u cx34collector.service
```

#### Developer notes

`./tools/install-waterpi.sh` is an example script that will build and copy
binaries to a remote Raspberry pi.

`grpcui -plaintext 192.168.86.24:8084` will run a local
[`grpcui`](https://github.com/fullstorydev/grpcui) frontend for querying the
heatpump after `grpcui` is installed.

### CX34 control

## General development tips

### Wiring

The Chiltrix uses RS485 to communicate with its controller. I have a
USB-to-RS484 with a ch340T chip and pins labelbed "A / D+" and "B / D-".

TODO: more info

### ssh access with password login disabled

My `~/.ssh/config` has a section like

```
Host waterpi
    HostName 192.168.86.22
    User pi
```

and the Pi's `~/.ssh/authorized_hosts` and `/etc/ssh/sshd_config` has a section with

```
# To disable tunneled clear text passwords, change to no here!
PasswordAuthentication no
```

### sshfs

Mounting the Raspberry Pi's OS can be helpful for developing on a workstation. I followed [this guide](https://www.digitalocean.com/community/tutorials/how-to-use-sshfs-to-mount-remote-file-systems-over-ssh) and ran these commands:

``` shell
sudo mkdir /mnt/waterpi
sudo chown <YOUR_USERNAME> /mnt/waterpi
sudo sshfs -o "allow_other,default_permissions,IdentityFile=/home/<YOUR_USERNAME>/.ssh/id_rsa" pi@192.168.86.22:/ /mnt/waterpi
```

### Fan coil unit control and scheduling

There is a gRPC server that can be started for getting and setting fan coil
state. Run this command on the fan coil unit:

```shell
go run fancoil/cmd/fancoil_status/fancoil_status.go --alsologtostderr --start-server --grpc-port 8083
```

You can then run the scheduler on the device:

```shell
go run fancoil/cmd/fancoil_schedule/fancoil_schedule.go --alsologtostderr --fancoil-service localhost:8083
```


# Cloud services

For the Google Home app integration to work, some cloud-hosted services are
needed. These can be deployed as follows:

## 1) Build and push docker images

(optional): Build the docker images used to run all of the services. This is
also done automatically by a Cloud Build hook that is triggered by pushes to
`main`, but it can be useful for pushing local images.

```shell
./cloud/deployment/cloudbuild/image-pusher/push-all-images.sh
```

## 2) Update the images referenced by the Terraform config

Terraform is used to deploy the Cloud Run images based on the contents of
`cloud/deployment/terraform/environments/infra-dev/image-versions.json`, which
is pulled into the main terraform config
`cloud/deployment/terraform/environments/infra-dev/main.tf`.

To update the pinned version of the images, run the following:

```shell
bazel run --run_under="cd $PWD &&" //cmd/cloud/update-image-versions -- --alsologtostderr --input "cloud/deployment/terraform/environments/infra-dev/image-versions.json" --output "cloud/deployment/terraform/environments/infra-dev/image-versions.json"
```

## 3) Configure the Actions Console settings

The following linkes might need to be updated if the Cloud Run deployment
changes significantly:

https://console.actions.google.com/u/0/project/hydronics-9f50d/actions/smarthome/
https://console.actions.google.com/u/0/project/hydronics-9f50d/smarthomeaccountlinking/
