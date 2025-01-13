#!/bin/bash

# Requires: https://github.com/anderseknert/kube-review
# Usage: ./get-admission-review-json.sh <resource-file>
# Usage: ./scripts/get-admission-review-json.sh testdata/argo/cronworkflow.yaml

# NOTE: can switch out create for different operations, e.g. update, delete, connect, etc.
kube-review create $1