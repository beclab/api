#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)

# Resolve the code-generator package location. Prefer (in order):
#   1. an explicit CODEGEN_PKG override
#   2. a vendored copy at ./vendor/k8s.io/code-generator
#   3. the module cache entry reported by `go list`
#   4. a sibling ../code-generator checkout (legacy layout)
if [[ -z "${CODEGEN_PKG:-}" ]]; then
    if [[ -d "${SCRIPT_ROOT}/vendor/k8s.io/code-generator" ]]; then
        CODEGEN_PKG="${SCRIPT_ROOT}/vendor/k8s.io/code-generator"
    else
        # Ensure the module is present in the local cache before resolving its path.
        (cd "${SCRIPT_ROOT}" && go mod download k8s.io/code-generator) >/dev/null
        CODEGEN_PKG=$(cd "${SCRIPT_ROOT}" && go list -m -f '{{.Dir}}' k8s.io/code-generator 2>/dev/null || true)
        if [[ -z "${CODEGEN_PKG}" ]]; then
            CODEGEN_PKG="${SCRIPT_ROOT}/../code-generator"
        fi
    fi
fi

if [[ ! -f "${CODEGEN_PKG}/kube_codegen.sh" ]]; then
    echo "error: cannot find kube_codegen.sh in CODEGEN_PKG=${CODEGEN_PKG}" >&2
    echo "hint: run 'go mod download k8s.io/code-generator' or set CODEGEN_PKG to an existing checkout" >&2
    exit 1
fi

source "${CODEGEN_PKG}/kube_codegen.sh"

THIS_PKG="github.com/beclab/api"

kube::codegen::gen_helpers \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
    "${SCRIPT_ROOT}/api"

kube::codegen::gen_client \
    --with-watch \
    --output-dir "${SCRIPT_ROOT}/pkg/generated" \
    --output-pkg "${THIS_PKG}/pkg/generated" \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
    --plural-exceptions "Terminus:Terminus" \
    "${SCRIPT_ROOT}/api"

