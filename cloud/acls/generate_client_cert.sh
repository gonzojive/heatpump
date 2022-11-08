#!/bin/bash
# From https://dev.to/techschoolguru/how-to-secure-grpc-connection-with-ssl-tls-in-go-4ph

set -e # fail on non-zero exit

REPO_ROOT=`git rev-parse --show-toplevel`

# Generate client's private key and certificate signing request (CSR)
echo "Generating client's private key"
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-signing-req.pem -subj "/C=US/ST=California/L=San Mateo/O=Gonzojive Heatpump Open Source Project/OU=IT/CN=*.daly.red/emailAddress=reddaly\+heatpump@gmail.com"

echo "Signing client's CSR"
# Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-signing-req.pem -days 60 -CA client-signer-cert-authority-cert.pem -CAkey client-signer-cert-authority-key.pem -CAcreateserial -out client-cert.pem -extfile "$REPO_ROOT/cloud/acls/client-x509-extfile.cnf"

#echo "Client's signed certificate"
#openssl x509 -in client-cert.pem -noout -text
