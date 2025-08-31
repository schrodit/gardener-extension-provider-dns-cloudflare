
VERSION=$(cat ./VERSION)
if [ -z "${VERSION}" ]; then
    goreleaser build --snapshot --clean
    jq -r '.version' ./dist/metadata.json > ./VERSION
    VERSION=$(cat ./VERSION)
fi
echo "$VERSION"
