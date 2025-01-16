#!/bin/bash

# https://kubernetes.io/docs/tasks/administer-cluster/certificates/#openssl

SCRIPT_DIR=$(dirname $0)
SECRETS_DIR="testdata"
DEPLOY_DIR="deploy"
mkdir -p $SECRETS_DIR
mkdir -p $DEPLOY_DIR

# Generate CA
if [[ ! -f $SECRETS_DIR/ca.key ]]; then
  openssl genrsa -out $SECRETS_DIR/ca.key 2048
fi

# Generate CA Certificate
if [[ ! -f $SECRETS_DIR/ca.crt ]]; then
  openssl req -x509 -new -nodes \
    -key $SECRETS_DIR/ca.key \
    -subj "/CN=argo-webhook" \
    -days 10000 -out $SECRETS_DIR/ca.crt
fi

# Generate Server Key
if [[ ! -f $SECRETS_DIR/server.key ]]; then
  openssl genrsa -out $SECRETS_DIR/server.key 2048
fi

# Generate CSR
if [[ ! -f $SECRETS_DIR/server.csr ]]; then
  openssl req -new \
    -key $SECRETS_DIR/server.key \
    -out $SECRETS_DIR/server.csr \
    -config $SCRIPT_DIR/csr.conf
fi

# Generate Server Certificate
if [[ ! -f $SECRETS_DIR/server.crt ]]; then
  openssl x509 -req \
    -in $SECRETS_DIR/server.csr \
    -CA $SECRETS_DIR/ca.crt -CAkey $SECRETS_DIR/ca.key \
    -CAcreateserial -out $SECRETS_DIR/server.crt -days 10000 \
    -extensions v3_ext -extfile $SCRIPT_DIR/csr.conf -sha256

  CA_CERT_BASE64=$(cat $SECRETS_DIR/ca.crt | base64)
  #yq -i ".webhooks[0].clientConfig.caBundle = $CA_CERT_BASE64" deploy/resources/mutatingwebhook.yaml

  TLS_CRT_BASE64=$(cat $SECRETS_DIR/server.crt | base64)
  TLS_KEY_BASE64=$(cat $SECRETS_DIR/server.key | base64)

  # Regenerate secret-tls.yaml
  cat <<EOF > $DEPLOY_DIR/resources/secret-tls.yaml
apiVersion: v1
kind: Secret
metadata:
  name: argo-webhook-tls
type: kubernetes.io/tls
data:
  tls.crt:
    ${TLS_CRT_BASE64}
  tls.key:
    ${TLS_KEY_BASE64}
EOF

fi

echo "[INFO]: Add Base64 caCertificate to deploy/resources/mutatingwebhook.yaml"
echo $CA_CERT_BASE64

echo "[INFO]: Validating Server Certificate with CA"
openssl verify -verbose -CAfile $SECRETS_DIR/ca.crt $SECRETS_DIR/server.crt