package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/adrg/xdg"
	"github.com/golang/glog"
	"github.com/mtraver/iotcore"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	if err := run(); err != nil {
		glog.Exitf("error: %v", err)
	}
}

func run() error {
	flag.Parse()

	c, err := NewClient(context.Background())
	if err != nil {
		return err
	}

	glog.Infof("got client %v", c)
	return nil
}

const (
	configDirName = "heatpump"
)

type Client struct {
	device *iotcore.Device
	c      mqtt.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	dev, err := loadIOTCoreDevice()
	if err != nil {
		return nil, err
	}

	rootsPath, err := xdg.ConfigFile(path.Join(configDirName, "roots.pem"))
	if err != nil {
		return nil, fmt.Errorf("could not determine roots file path %w", err)
	}
	rootsCertBytes, err := ioutil.ReadFile(rootsPath)
	if err != nil {
		return nil, fmt.Errorf("could not read roots file %q: %w", rootsPath, err)
	}

	client, err := dev.NewClient(iotcore.DefaultBroker, bytes.NewReader(rootsCertBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to make MQTT client: %w", err)
	}

	if err := waitToken(ctx, client.Connect()); err != nil {
		return nil, fmt.Errorf("failed to connect to device IoT server: %w", ctx.Err())
	}

	if err := waitToken(ctx, client.Publish(dev.TelemetryTopic(), 1, true, []byte("{\"temp\": 18.0}"))); err != nil {
		return nil, fmt.Errorf("failed to publish: %w", err)
	}

	if err := waitToken(ctx, client.Publish(dev.StateTopic(), 1, true, []byte("{\"old\": true}"))); err != nil {
		return nil, fmt.Errorf("failed to publish state: %w", err)
	}

	if err := waitToken(ctx, client.Subscribe(dev.ConfigTopic(), 1, func(_ mqtt.Client, msg mqtt.Message) {
		glog.Infof("got config message %q", string(msg.Payload()))
	})); err != nil {
		return nil, fmt.Errorf("failed to publish state: %w", err)
	}

	time.Sleep(time.Minute)

	return &Client{
		dev,
		client,
	}, nil
}

func waitToken(ctx context.Context, tok mqtt.Token) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-tok.Done():
		return tok.Error()
	}

}

func loadIOTCoreDevice() (*iotcore.Device, error) {
	relPath := path.Join(configDirName, "devices/grant-cx34/ec_private.pem")
	privateKeyPath, err := xdg.ConfigFile(relPath)
	if err != nil {
		return nil, fmt.Errorf("could not find private key config file %q: %w", relPath, err)
	}

	return &iotcore.Device{
		ProjectID:   "redapps",
		RegistryID:  "heatpumps",
		DeviceID:    "grant-cx34",
		PrivKeyPath: privateKeyPath,
		Region:      "us-central1",
	}, nil
}
