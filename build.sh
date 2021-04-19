#!/bin/sh
set -e
go fmt
(cd env0tfprovider && go fmt)
(cd env0apiclient && go fmt)
CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o terraform-provider-env0
