package routing

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/andidroid/testgo/pkg/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var client *http.Client

var ROUTING_SERVCICE_URL string
var FLEET_SERVCICE_URL string
var SEARCH_SERVCICE_URL string
var MESSAGING_SERVCICE_URL string

func init() {

	ROUTING_SERVCICE_URL = util.LookupEnv("ROUTING_SERVCICE_URL", func() string {
		util.LookupEnv("ROUTING_SERVCICE_HOST", "localhost")
		util.LookupEnv("ROUTING_SERVCICE_PORT", "80")
		util.LookupEnv("ROUTING_SERVCICE_PROTOCOLL", "http")
		return os.ExpandEnv("$ROUTING_SERVCICE_PROTOCOLL://$ROUTING_SERVCICE_HOST:$ROUTING_SERVCICE_PORT")
	}())
	fmt.Println("reading ROUTING_SERVCICE_URL: ", ROUTING_SERVCICE_URL)

	FLEET_SERVCICE_URL = util.LookupEnv("FLEET_SERVCICE_URL", func() string {
		util.LookupEnv("FLEET_SERVCICE_HOST", "localhost")
		util.LookupEnv("FLEET_SERVCICE_PORT", "80")
		util.LookupEnv("FLEET_SERVCICE_PROTOCOLL", "http")
		return os.ExpandEnv("$FLEET_SERVCICE_PROTOCOLL://$FLEET_SERVCICE_HOST:$FLEET_SERVCICE_PORT")
	}())
	fmt.Println("reading FLEET_SERVCICE_URL: ", FLEET_SERVCICE_URL)

	SEARCH_SERVCICE_URL = util.LookupEnv("SEARCH_SERVCICE_URL", func() string {
		util.LookupEnv("SEARCH_SERVCICE_HOST", "localhost")
		util.LookupEnv("SEARCH_SERVCICE_PORT", "80")
		util.LookupEnv("SEARCH_SERVCICE_PROTOCOLL", "http")
		return os.ExpandEnv("$SEARCH_SERVCICE_PROTOCOLL://$SEARCH_SERVCICE_HOST:$SEARCH_SERVCICE_PORT")
	}())
	fmt.Println("reading SEARCH_SERVCICE_URL: ", SEARCH_SERVCICE_URL)

	MESSAGING_SERVCICE_URL = util.LookupEnv("MESSAGING_SERVCICE_URL", func() string {
		util.LookupEnv("MESSAGING_SERVCICE_HOST", "localhost")
		util.LookupEnv("MESSAGING_SERVCICE_PORT", "80")
		util.LookupEnv("MESSAGING_SERVCICE_PROTOCOLL", "http")
		return os.ExpandEnv("$MESSAGING_SERVCICE_PROTOCOLL://$MESSAGING_SERVCICE_HOST:$MESSAGING_SERVCICE_PORT")
	}())
	fmt.Println("reading MESSAGING_SERVCICE_URL: ", MESSAGING_SERVCICE_URL)

	client = CreateClient()

}

func GetClient() *http.Client {

	// http.DefaultClient().
	return client
}

func CreateClient() *http.Client {

	// certPool := x509.NewCertPool()
	// certPool.AppendCertsFromPEM(caCert)
	// transport:= &http.Transport{
	// 	// TLSClientConfig: &tls.Config{
	// 	//   RootCAs:      certPool,
	// 	//   Certificates: []tls.Certificate{ClientCert},
	// 	// },
	// };

	transport := http.DefaultTransport
	client := &http.Client{
		Timeout:   time.Second,
		Transport: otelhttp.NewTransport(transport),
	}

	// http.DefaultClient().
	return client
}
