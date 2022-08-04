package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"crypto/tls"

	echo "github.com/jpillora/go-echo-server/handler"
	"github.com/jpillora/go-echo-server/udp"
	"github.com/jpillora/opts"
)

var version = "0.0.0-src"

var (
	tlscert, tlskey string
	port            int
	secure          bool
)

func main() {
	certpem := getEnv("CERT_FILE", "certs/cert.pem")
	keypem := getEnv("KEY_FILE", "certs/key.pem")

	flag.StringVar(&tlscert, "tlsCertFile", certpem, "File contaains the X509 Certificate for HTTPS")
	flag.StringVar(&tlskey, "tlsKeyFile", keypem, "File containing the X509 private key")
	flag.IntVar(&port, "port", 443, "http port")
	flag.BoolVar(&secure, "secure", true, "use tls")

	flag.Parse()

	c := struct {
		Port        int  `help:"Port" env:"PORT"`
		UDP         bool `help:"UDP mode"`
		echo.Config `type:"embedded"`
	}{
		Port: port,
	}
	opts.New(&c).
		Name("go-echo-server").
		Version(version).
		Repo("github.com/jpillora/go-echo-server").
		Parse()
	//udp mode?
	if c.UDP {
		log.Fatal(udp.Start(c.Port))
	}
	//http mode
	server := http.Server{
		Addr: fmt.Sprintf(":%v", c.Port),
	}

	if secure {
		certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
		if err != nil {
			log.Fatalf("Failed to load key pair: %v", err)
			// log.Errorf("Failed to load key pair: %v", err)
		}
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{certs},
		}
	}
	h := echo.New(c.Config)
	log.Printf("Listening for http requests on %d...", c.Port)
	server.Handler = h
	if secure {
		log.Fatal(server.ListenAndServeTLS("", ""))

	} else {
		log.Fatal(server.ListenAndServe())
	}

	// log.Fatal(http.ListenAndServe(":"+strconv.Itoa(c.Port), h))
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
