// Program reverse-proxy forwards all http traffic from this process to
// a remote url.
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/golang/glog"
)

const cloudRunPortEnvVar = "PORT"

func main() {
	rpURL, err := url.Parse("http://daly.red:42598")
	if err != nil {
		glog.Fatal(err)
	}

	http.Handle("/", httputil.NewSingleHostReverseProxy(rpURL))

	port := "8080"

	// PORT env is from Cloud Run.
	if got := os.Getenv(cloudRunPortEnvVar); got != "" {
		port = got
	}

	glog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
