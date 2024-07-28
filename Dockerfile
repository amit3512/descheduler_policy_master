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
###########################################################
# FROM golang:1.22.4

# WORKDIR /go/src/sigs.k8s.io/descheduler
# COPY . .
# ARG ARCH
# ARG VERSION
# RUN VERSION=${VERSION} make build.$ARCH

# FROM scratch

# MAINTAINER Kubernetes SIG Scheduling <kubernetes-sig-scheduling@googlegroups.com>

# USER 1000

# COPY --from=0 /go/src/sigs.k8s.io/descheduler/_output/bin/descheduler /bin/descheduler

# CMD ["/bin/descheduler", "--help"]
################################################################################
# Use a base image with Go installed
FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /go/src/github.com/kubernetes-sigs/descheduler

# Copy the local code to the container
COPY . .

# Build the descheduler binary
RUN make build

# Use a minimal image for the final stage
FROM gcr.io/distroless/static:latest

# Set the working directory inside the container
WORKDIR /

# Copy the descheduler binary from the builder
COPY --from=builder /go/src/github.com/kubernetes-sigs/descheduler/_output/bin/descheduler /descheduler

# Command to run the descheduler
ENTRYPOINT ["/descheduler"]

