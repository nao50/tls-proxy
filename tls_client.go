package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
)

func main() {
	// Read CA certificate
	caFile, _ := filepath.Abs("2_tls_server/certificate/certs/ca.pem")
	rootPem, err := ioutil.ReadFile(caFile)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}

	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(rootPem); !ok {
		fmt.Printf("Err: failed to parse root certificate")
		return
	}
	//

	cert, err := tls.LoadX509KeyPair("2_tls_server/certificate/certs/client.pem", "2_tls_server/certificate/certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
		// InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()

	log.Println("client: connected to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println("Client: Server public key is:")
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	message := "Hello\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)
	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("client: exiting")
}
