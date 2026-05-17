#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

CODE_GEN_DIR=$(go list -mod=mod -m -f '{{.Dir}}' k8s.io/code-generator)
source "${CODE_GEN_DIR}/kube_codegen.sh"

rm -f "$(go env GOPATH)"/bin/*-gen

CURRENT_DIR=$(dirname $0)
PROJECT_ROOT="${CURRENT_DIR}"/..

kube::codegen::gen_helpers \
--boilerplate "${GARDENER_HACK_DIR}/LICENSE_BOILERPLATE.txt" \
"${PROJECT_ROOT}/pkg/apis/cloudflare"

kube::codegen::gen_helpers \
--boilerplate "${GARDENER_HACK_DIR}/LICENSE_BOILERPLATE.txt" \
"${PROJECT_ROOT}/pkg/apis/config"
