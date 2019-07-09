package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/lucas-clemente/quic-go/http3"
)

func GetRootCA(certPath string) *x509.CertPool {
	caCertPath := path.Join(certPath, "ca.pem")
	caCertRaw, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}
	p, _ := pem.Decode(caCertRaw)
	if p.Type != "CERTIFICATE" {
		panic("expected a certificate")
	}
	caCert, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		panic(err)
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(caCert)
	return certPool
}

func main() {
	currentPath, err := os.Getwd()
	// currentPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(currentPath)
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs: GetRootCA(currentPath),
		},
	}
	defer roundTripper.Close()
	client := &http.Client{
		Transport: roundTripper,
	}
	addr := "https://localhost:8080"
	rsp, err := client.Get(addr)
	if err != nil {
		log.Fatal(err)
	}

	defer rsp.Body.Close()

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("Body length: %d bytes", body.Len())
	log.Printf("Response body %s", body.Bytes())
}
