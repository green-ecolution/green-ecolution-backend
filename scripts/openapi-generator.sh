#!/bin/env bash

API_DOCS_DIR="./pkg/client"

generate_client_docker() {
  docker run --rm -v "${PWD}":/local -u "$(id -u):$(id -g)" openapitools/openapi-generator-cli generate \
    -i /local/docs/swagger.yaml \
    -o /local/$API_DOCS_DIR \
    -g go \
    --package-name client \
    --git-host github.com \
    --git-user-id green-ecolution \
    --git-repo-id green-ecolution-backend \
    --additional-properties=isGoSubmodule=true 

  exit_on_error "Could not create the client"
}

cleanup() {
  rm -rf $API_DOCS_DIR
}

exit_on_error() {
  if [ $? -ne 0 ]; then
    echo "❌ $1"
    exit 1
  fi
}

cleanup
generate_client_docker
