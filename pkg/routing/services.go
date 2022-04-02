package routing

import (
	"net/http"
	"time"
)

var client *http.Client

func init() {
	client = CreateClient()
}

func GetClient() *http.Client {

	// http.DefaultClient().
	return client
}

func CreateClient() *http.Client {

	// certPool := x509.NewCertPool()
	// certPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
        Timeout: time.Second,
		Transport: &http.Transport{
			// TLSClientConfig: &tls.Config{
			//   RootCAs:      certPool,
			//   Certificates: []tls.Certificate{ClientCert},
			// },
		},
	}

	// http.DefaultClient().
	return client
}
