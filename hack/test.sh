
export KUBEBUILDER_ASSETS="$(setup-envtest use --use-env -p path ${ENVTEST_K8S_VERSION})"
GARDENER_HACK_DIR="$(go list -m -f "{{.Dir}}" github.com/gardener/gardener)/hack"
bash "${GARDENER_HACK_DIR}/test.sh" ./cmd/... ./pkg/...
