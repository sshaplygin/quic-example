package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/lucas-clemente/quic-go/http3"
)

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs: getRootCA(currentPath),
		},
	}
	defer roundTripper.Close()

	client := &http.Client{
		Transport: roundTripper,
	}

	addr := "https://localhost:8080"
	rsp, err := client.Get(addr)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("Body length: %d bytes \n", body.Len())
	log.Printf("Response body %s \n", body.Bytes())
}

func getRootCA(certPath string) *x509.CertPool {
	caCertPath := path.Join(certPath, "ca.pem")
	caCertRaw, err := os.ReadFile(caCertPath)
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
