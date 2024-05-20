package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var caCert x509.Certificate
var leafCert x509.Certificate
var leafKey rsa.PrivateKey

func main() {
	errors := 0
	cert := flag.String("cert", "/etc/ssl/cert.crt", "Path to the certificate")
	day := flag.Int("day", 30, "Minimum days left on the cert, default 30")
	flag.Parse()

	// load and parse cert/CA/key
	certLoad, err := ioutil.ReadFile(*cert)
	if err != nil {
		fmt.Printf("failed to read cert file: %s\n", err)
		os.Exit(1)
	}
	for block, rest := pem.Decode(certLoad); block != nil; block, rest = pem.Decode(rest) {
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				panic(err)
			}
			if cert.IsCA {
				caCert = *cert
			} else {
				leafCert = *cert
			}
		case "PRIVATE KEY":

			if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
				leafKey = *key
			}
			if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
				leafKey = *key.(*rsa.PrivateKey)
			}
		default:
			fmt.Printf("you should not land here, the certificate is most likely invalid")
			os.Exit(1)
		}
	}

	// Check cert/key/CA are present
	if len(leafCert.Subject.SerialNumber) == 0 {
		errors += 1
		fmt.Printf("Cert is missing\n")
	}
	if len(caCert.Subject.SerialNumber) == 0 {
		errors += 1
		fmt.Printf("CA is missing\n")
	}
	if len(leafKey.PublicKey.N.Bytes()) == 0 {
		errors += 1
		fmt.Printf("CA is missing\n")
	}

	// check leaf age
	leafDaysLeft := int(time.Until(leafCert.NotAfter) / (24 * time.Hour))
	if leafDaysLeft < *day {
		errors += 1
	}
	fmt.Printf("cert has %d days left (minumum %d)\n", leafDaysLeft, *day)

	// check leaf if NotBefore
	if time.Until(leafCert.NotBefore) > 0 {
		errors += 1
		fmt.Printf("cert is not yet valid: %s\n", leafCert.NotBefore)
	}
	// check CA expiry
	caDaysLeft := int(time.Until(caCert.NotAfter) / (24 * time.Hour))
	if caDaysLeft < *day {
		errors += 1
		fmt.Printf("CA is expiring in: %d days\n", caDaysLeft)
	}

	// check CA if NotBefore
	if time.Until(caCert.NotBefore) > 0 {
		errors += 1
		fmt.Printf("CA is not yet valid: %s\n", caCert.NotBefore)
	}

	// check CA match leaf
	roots := x509.NewCertPool()
	roots.AddCert(&caCert)
	verifyCA := x509.VerifyOptions{Roots: roots}
	if _, err := leafCert.Verify(verifyCA); err != nil {
		errors += 1
		fmt.Printf("Cert failed to validate against CA: %s\n", err)
	}

	// Check key/cert matching
	leafCertPub := leafCert.PublicKey
	leafCertPubKey := *leafCertPub.(*rsa.PublicKey)
	if !bytes.Equal(leafCertPubKey.N.Bytes(), leafKey.N.Bytes()) {
		errors += 1
		fmt.Printf("Cert and Key do not match!\n")
	}

	// exit
	if errors != 0 {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
