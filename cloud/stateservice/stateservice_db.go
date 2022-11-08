package stateservice

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/golang/protobuf/proto"
	"github.com/gonzojive/heatpump/cloud/acls"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/proto/controller"
)

type Store interface {
	StoreDeviceState(ctx context.Context, user *acls.Identity, s *controller.DeviceState) error
	GetDeviceState(ctx context.Context, user *acls.Identity, name string) (*controller.DeviceState, error)
	ListDevices(ctx context.Context, user *acls.Identity) ([]string, error)
}

const (
	devicesCollectionName = "iot-devices"
	reportedStateProtoCol = "reported-device-state-proto"
)

func userDeviceDoc(client *firestore.Client, user *acls.Identity, deviceID string) *firestore.DocumentRef {
	return client.Collection(devicesCollectionName).
		Doc(fmt.Sprintf("user-%s-device-%s", user.ID(), deviceID))
}

type firestoreStore struct {
	client *firestore.Client
}

func newFirestoreBackedStore(ctx context.Context, params *cloudconfig.Params) (*firestoreStore, error) {
	fsClient, err := createFirestoreClient(ctx, params.GCPProject)
	if err != nil {
		return nil, err
	}
	return &firestoreStore{fsClient}, nil
}

func (st *firestoreStore) StoreDeviceState(ctx context.Context, user *acls.Identity, s *controller.DeviceState) error {
	encoded, err := proto.Marshal(s)
	if err != nil {
		return err
	}
	// TODO(reddaly): Transactionalize.
	if _, err := st.client.Collection(devicesCollectionName).
		Doc(fmt.Sprintf("user-%s-device-%s", user.ID(), s.GetName())).
		Set(ctx, map[string]interface{}{
			reportedStateProtoCol: encoded,
		}, firestore.MergeAll); err != nil {
		return fmt.Errorf("error adding state to firestore: %w", err)
	}
	return nil
}

func (st *firestoreStore) GetDeviceState(ctx context.Context, user *acls.Identity, name string) (*controller.DeviceState, error) {
	snapshot, err := userDeviceDoc(st.client, user, name).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting state from firestore: %w", err)
	}
	state := &controller.DeviceState{}
	if err := dataAtPathToProto(snapshot, reportedStateProtoCol, state); err != nil {
		return nil, fmt.Errorf("error getting state from firestore: %w", err)
	}
	return state, nil
}

func (st *firestoreStore) ListDevices(ctx context.Context, user *acls.Identity) ([]string, error) {
	return nil, fmt.Errorf("unimplemented")
}

func createFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	// Sets your Google Cloud Platform project ID.
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create FireStore client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client, nil
}

func dataAtPathToProto(d *firestore.DocumentSnapshot, path string, msg proto.Message) error {
	val, err := d.DataAt(path)
	if err != nil {
		return fmt.Errorf("error getting data at path %q: %w", path, err)
	}
	wireBytes, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("error getting data at path %q: not a bytes field", path)
	}
	return proto.Unmarshal(wireBytes, msg)
}
