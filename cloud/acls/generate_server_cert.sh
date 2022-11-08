#!/bin/bash
# From https://dev.to/techschoolguru/how-to-secure-grpc-connection-with-ssl-tls-in-go-4ph

set -e # fail on non-zero exit

REPO_ROOT=`git rev-parse --show-toplevel`

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 3650 -nodes -keyout server-ca-key.pem -out server-ca-cert.pem -subj "/C=US/ST=California/L=San Mateo/O=Gonzojive Heatpump Open Source Project/OU=IT/CN=*.daly.red/emailAddress=reddaly\+heatpump@gmail.com"

echo "CA's self-signed certificate is at server-ca-cert.pem"

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=US/ST=California/L=San Mateo/O=Gonzojive Heatpump Open Source Project/OU=IT/CN=*.daly.red/emailAddress=reddaly\+heatpump@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA server-ca-cert.pem -CAkey server-ca-key.pem -CAcreateserial -out server-cert.pem -extfile "$REPO_ROOT/cloud/acls/server-x509-extfile.cnf"

echo "Server's signed certificate is at server-cert.pem"