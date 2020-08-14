package depscloud

import (
	"crypto/tls"
	"net/url"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	VariableBaseURL = "DEPSCLOUD_BASE_URL"

	DefaultBaseURL = "https://api.deps.cloud"
)

func translateBaseURL(baseURL string) (bool, string) {
	tls := false
	uri, _ := url.Parse(baseURL)

	if uri.Scheme == "https" {
		tls = true
	}

	host := uri.Host
	if !strings.Contains(host, ":") {
		if tls {
			host = host + ":443"
		} else {
			host = host + ":80"
		}
	}

	return tls, host
}

func Connect(options ...grpc.DialOption) (*grpc.ClientConn, error) {
	baseURL := DefaultBaseURL
	if val := os.Getenv(VariableBaseURL); val != "" {
		baseURL = val
	}

	isSecure, target := translateBaseURL(baseURL)

	if isSecure {
		creds := credentials.NewTLS(&tls.Config{})
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithInsecure())
	}

	return grpc.Dial(target, options...)
}
