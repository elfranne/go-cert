
# Generate CA:
openssl req -newkey rsa:4096 -keyout ca.key -new -x509  -out ca.crt -days 3650 -config config_ca.cnf -nodes

# Generate CSR:
openssl req -newkey rsa:4096 -keyout go-cert.key -new -sha256 -config config_cert.cnf -out go-cert.csr -nodes

# Required files for CA:
echo 00 > ca.srl
touch index.txt

# Generate valid self signed certificate:
openssl ca -config config_ca_sign.cnf -out go-cert.crt -in go-cert.csr

# combine
cat ca.crt go-cert.crt go-cert.key > go-cert-combined.crt

# Generate expired self signed certificate:
openssl ca -config config_ca_sign.cnf -out go-cert-expired.crt -in go-cert.csr -startdate 20000101000000Z -enddate 20000201000000Z

# combine
cat ca.crt go-cert-expired.crt go-cert.key > go-cert-expired-combined.crt

# Generate not yet valid self signed certificate:
openssl ca -config config_ca_sign.cnf -out go-cert-not-yet.crt -in go-cert.csr -startdate  20500101000000Z -enddate 20500201000000Z

# combine
cat ca.crt go-cert-not-yet.crt go-cert.key > go-cert-not-yet-combined.crt





