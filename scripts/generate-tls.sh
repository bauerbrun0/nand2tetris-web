#!/bin/bash

HOSTS="localhost,127.0.0.1"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --local-ip)
      if [[ -n "$2" ]]; then
        HOSTS="$HOSTS,$2"
        shift
      else
        echo "Error: --local-ip requires an argument"
        exit 1
      fi
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
  shift
done

go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host "$HOSTS"
mkdir -p tls
mv cert.pem tls/
mv key.pem tls/
