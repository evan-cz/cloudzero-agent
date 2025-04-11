#!/usr/bin/env bash

set -euo pipefail

FAILED=false

# Check Go version consistency
function check_go_version() {
    DESIRED_GO_VERSION="$(git grep -hPo '^go [0-9]+\.[0-9]+(\.[0-9]+)?$' go.mod | awk '{print $2}')"
    if [ -z "$DESIRED_GO_VERSION" ]; then
        echo "Error: No Go version found in go.mod" >&2
        exit 1
    fi
    DESIRED_GO_VERSION_NO_MICRO="$(echo $DESIRED_GO_VERSION | awk -F. '{print $1"."$2}')"

    if [ "$DESIRED_GO_VERSION_NO_MICRO" = "$DESIRED_GO_VERSION" ]; then
        echo "Error: Go version in go.mod does not include a micro version (${DESIRED_GO_VERSION} -> ${DESIRED_GO_VERSION}.0?)" >&2
        FAILED=TRUE
    fi

    # Dockerfiles
    find . -type f -iname 'Dockerfile*' -print0 | while IFS= read -r -d '' DOCKERFILE; do
        git grep -q " golang:${DESIRED_GO_VERSION}" ${DOCKERFILE} || {
            echo "${DOCKERFILE} does not have the desired Go version (${DESIRED_GO_VERSION})" >&2
            FAILED=true
        }
    done

    # go.mod
    for GO_MOD in \
            go.mod \
            tests/integration/test_server/go.mod; do
        git grep -q "^go ${DESIRED_GO_VERSION}\$" ${GO_MOD} || {
            echo "${GO_MOD} does not have the desired Go version (${DESIRED_GO_VERSION})" >&2
            FAILED=true
        }
    done
}
check_go_version

# Check that Helm chart contains expected pattern(s) for version bump
function check_helm_chart_version_bump() {
    # This is the pattern we use in .github/workflows/release-to-main.yml to
    # update the version of the CloudZero Agent container in the Helm chart.
    # This just exists to make sure we don't accidentally break it by tweaking
    # the comment or something.
    if ! grep -Eq '^( +tag): +[^ ]+  (# <- Software release corresponding to this chart version.)$' helm/values.yaml; then
        echo "Helm chart does not contain expected pattern for version bump" >&2
        FAILED=true
    fi
}
check_helm_chart_version_bump
if [ "$FAILED" = true ]; then
    exit 1
fi
