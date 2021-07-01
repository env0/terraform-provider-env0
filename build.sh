#!/usr/bin/env sh

set -e
go fmt
(cd env0 && go fmt)
(cd client && go fmt)
CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o terraform-provider-env0
