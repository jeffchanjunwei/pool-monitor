package main

import (
	//"os"
	"log"
	"time"
	"flag"
	"strings"
	"net/http"
	"crypto/tls"
        "crypto/x509"
        "io/ioutil"
        "net"
        "errors"
        "fmt"


	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func main() {

	var etcdAddr []string
	var cert string
	var key string
	var caCert string
        var etcdEndPoints string

        flag.StringVar(&etcdEndPoints, "etcdEndPoints", "http://127.0.0.1:2379", "The etcd address used for ip pool management")
        flag.StringVar(&caCert, "caCert", "", "The CA certificate when using https")
        flag.StringVar(&key, "key", "", "The private key when using https")
        flag.StringVar(&cert, "cert", "", "The certificate when ")

        flag.Parse()

        etcdAddr = strings.Split(etcdEndPoints,",")

	var c client.Client
        var err error
        var transport = &http.Transport{
                Proxy: http.ProxyFromEnvironment,
                Dial: (&net.Dialer{
                        Timeout: 30 * time.Second,
                        KeepAlive: 30 * time.Second,
                }).Dial,
                TLSHandshakeTimeout: 10 * time.Second,
        }

        tlsConfig := &tls.Config{
                InsecureSkipVerify: false,
        }

        cfg := client.Config{
                Endpoints: etcdAddr,
                HeaderTimeoutPerRequest: time.Duration(3) * time.Second,
        }

        if caCert != "" {
                certBytes, err := ioutil.ReadFile(caCert)
                if err != nil {
                        fmt.Println(errors.New("CA file load failed"))
                }

                fmt.Println(certBytes)
                fmt.Println("test-1")

                caCertPool := x509.NewCertPool()
                ok := caCertPool.AppendCertsFromPEM(certBytes)

                if ok {
                        tlsConfig.RootCAs = caCertPool
                }
        }

	if cert != "" && key != "" {
                tlsCert, err := tls.LoadX509KeyPair(cert, key)
                if err != nil {
                        fmt.Println(err)
                        fmt.Println(errors.New("Cert file or key file load failed"))
                }
                tlsConfig.Certificates = []tls.Certificate{tlsCert}
        }

        transport.TLSClientConfig = tlsConfig
        cfg.Transport = transport

        c, err = client.New(cfg)
        if err != nil {
                fmt.Println(err)
                fmt.Println(errors.New("Error when creating etcd client"))
        }

        kapi := client.NewKeysAPI(c)

/*
	etcdAddr := os.Getenv("ETCD_ENDPOINTS")

	cfg := client.Config{
		Endpoints: []string{etcdAddr},
		Transport: client.DefaultTransport,

		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)
*/
	// watch the change of "/calico/v1/ipam/v4/pool"
	log.Println("watch the change of directory /calico/v1/ipam/v4/pool")

	 watcher := kapi.Watcher("/calico/v1/ipam/v4/pool", &client.WatcherOptions{
		Recursive: true,
	 })

	for true{
		resp, _ := watcher.Next(context.Background())
		//fmt.Println(resp)
		//fmt.Println(err)

		if resp != nil {
			Refresh()
		}
	}

}

