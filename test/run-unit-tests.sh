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

set -x
set -o errexit
set -o nounset

# This just run unit-tests. Ignoring the current directory so as to avoid running e2e tests.
PRJ_PREFIX="github.com/amit3512/descheduler_policy_master"
go test $(go list ${PRJ_PREFIX}/... | grep -v ${PRJ_PREFIX}/vendor/| grep -v ${PRJ_PREFIX}/test/)
