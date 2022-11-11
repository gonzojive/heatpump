/*
Copyright (c) 2019 Tadej Slamic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package oauthstore

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/util/bazelrunfiles"
	portpicker "github.com/johnsiilver/golib/development/portpicker"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/iterator"
)

const firebaseEmulatorBazelPath = "cloud/oauthstore/firestore-emulator"

/*
func TestMain(m *testing.M) {
	project, ok := os.LookupEnv("PROJECT_ID")
	if !ok {
		log.Fatalln("PROJECT_ID env variable is missing")
	}
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: project}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	c, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(func() int {
		defer c.Close()
		return m.Run()
	}())
}
*/

func emulatorClient(ctx context.Context, t *testing.T) (*firestore.Client, func()) {
	binPath, err := bazelrunfiles.Runfile(firebaseEmulatorBazelPath)
	if err != nil {
		t.Fatalf("failed to find firestore emulator path: %v", err)
	}

	version, err := exec.CommandContext(ctx, binPath, "--version").Output()
	if err != nil {
		t.Fatalf("error running %s --version: %v", binPath, err)
	}

	glog.Infof("firebase emulator version: %q", strings.TrimSpace(string(version)))

	port, err := portpicker.TCP(portpicker.Local4())

	if err != nil {
		t.Fatalf("failed to get port to use for firestore emulator: %v", err)
	}

	if err := os.Setenv("FIRESTORE_EMULATOR_HOST", fmt.Sprintf("localhost:%d", port)); err != nil {
		t.Fatalf("failed to set FIRESTORE_EMULATOR_HOST env var: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	cmd := exec.CommandContext(ctx, binPath, "--port", strconv.Itoa(port))

	var cmdOutput strings.Builder
	cmd.Stdout = &cmdOutput
	cmd.Stderr = &cmdOutput

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start the emulator: %v", err)
	}
	glog.Infof("starting emulator on localhost:%d", port)

	eg.Go(func() error {
		glog.Infof("waiting for emulator to finish...")
		err := cmd.Wait()
		if err, isExitErr := err.(*exec.ExitError); isExitErr {
			if err.ExitCode() == 130 { // terminated by owner (SIGTERM)
				return nil
			}
			glog.Errorf("emulator failed with error code %d: %v", err.ExitCode(), err)
		}
		if err != nil {
			glog.Infof("output from emulator: %q; err = %v; ctx.Err() = %v", cmdOutput.String(), err, ctx.Err())
			return err
		}
		return nil
	})
	cleanup := func() {
		glog.Infof("cleaning up emulator...")
		// Equivalent of Ctrl-c the emulator and wait for it to exit.
		if cmd.Process != nil {
			glog.Infof("sending SIGTERM to emulator...")
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				t.Fatalf("failed to send SIGTERM to emulator: %v", err)
			}
		}
		if err := eg.Wait(); err != nil {
			t.Fatalf("cleanup finished with error: %v", err)
		}
	}
	glog.Infof("starting NewClient()")
	c, err := firestore.NewClient(ctx, "test")
	if err != nil {
		t.Fatalf("error creating firestore client: %v", err)
		defer cleanup()
	}

	t.Run("emulatorstartup", func(t *testing.T) {
		testFirstoreIsSane(ctx, t, c)
	})

	glog.Infof("started emulator, continuing to test function with %v...", c)
	return c, cleanup
}

func TestFirstoreIsSane(t *testing.T) {
	ctx := context.Background()

	c, cancel := emulatorClient(ctx, t)
	defer cancel()
	testFirstoreIsSane(ctx, t, c)
}

func testFirstoreIsSane(ctx context.Context, t *testing.T, c *firestore.Client) {
	if _, err := c.Collection("things").
		Doc("thing2").
		Set(ctx, map[string]interface{}{
			"property": "forty two",
		}, firestore.MergeAll); err != nil {
		t.Fatalf("error adding state to firestore: %v", err)
	}

	doc, err := c.Collection("things").Doc("thing2").Get(ctx)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	got, err := doc.DataAt("property")
	if err != nil {
		t.Fatal(err)
	}
	if got.(string) != "forty two" {
		t.Fatalf("got %q, want %q", got.(string), "forty two")
	}
}

func TestStoreClient(t *testing.T) {
	ctx := context.Background()

	c, cancel := emulatorClient(ctx, t)
	defer cancel()
	client := NewTokenStorage(c, "tests")
	type holder struct {
		key string
		get func(context.Context, string) (oauth2.TokenInfo, error)
		del func(context.Context, string) error
	}
	tokens := map[*models.Token]holder{
		{Access: "access"}:   {key: "access", get: client.GetByAccess, del: client.RemoveByAccess},
		{Code: "code"}:       {key: "code", get: client.GetByCode, del: client.RemoveByCode},
		{Refresh: "refresh"}: {key: "refresh", get: client.GetByRefresh, del: client.RemoveByRefresh},
	}
	for i, h := range tokens {
		ctx := context.Background()
		err := client.Create(ctx, i)
		assert.Nil(t, err)

		tok, err := h.get(ctx, h.key)
		assert.Nil(t, err)
		assert.Equal(t, i, tok)

		err = h.del(ctx, h.key)
		assert.Nil(t, err)

		_, err = h.get(ctx, h.key)
		assert.NotNil(t, err)

		err = h.del(ctx, h.key)
		assert.Nil(t, err)
	}
}

func TestNoDocument(t *testing.T) {
	ctx := context.Background()
	c, cancel := emulatorClient(ctx, t)
	defer cancel()

	store := NewTokenStorage(c, "tests")
	info, err := store.GetByRefresh(ctx, "whoops")
	assert.Nil(t, info)
	if err != iterator.Done {
		t.Fatalf("expected iterator.Done, got err = %v", err)
	}
}

func TestIsNilOrZero(t *testing.T) {
	tokens := map[oauth2.TokenInfo]bool{
		nil:                               true,
		&models.Token{}:                   true,
		&models.Token{Access: "access"}:   false,
		&models.Token{Code: "code"}:       false,
		&models.Token{Refresh: "refresh"}: false,
	}
	for tok, expected := range tokens {
		result := isNilOrZero(tok)
		assert.Equal(t, expected, result)
	}
}
