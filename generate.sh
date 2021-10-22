#!/bin/sh

set -o errexit

rm -rf gen && mkdir gen cd gen
swagger generate cli -f swagger.yaml
cd -

rm gen/cli/models_runner_group_model.go




