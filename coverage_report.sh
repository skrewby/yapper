#!/usr/bin/env bash

t="/tmp/go-cover.$$.tmp"
go test -coverprofile=$t ./... && go tool cover -html=$t && unlink $t
