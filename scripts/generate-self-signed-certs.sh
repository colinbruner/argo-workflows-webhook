#!/bin/bash

FIXTURES_DIR="fixtures"
mkdir -p $FIXTURES_DIR

if [[ ! -f $FIXTURES_DIR/server.key ]] || [[ ! -f $FIXTURES_DIR/server.crt ]]; then
  openssl req -nodes -x509 -sha256 \
    -newkey rsa:2048 \
    -keyout $FIXTURES_DIR/server.key \
    -out $FIXTURES_DIR/server.crt \
    -days 3650 \
    -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname"
fi