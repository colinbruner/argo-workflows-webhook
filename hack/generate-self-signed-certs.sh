#!/bin/bash

openssl req \
  -x509 \
  -newkey rsa:2048 \
  -keyout server.key \
  -out server.crt \
  -sha256 \
  -days 3650 \
  -nodes \
  -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname"
