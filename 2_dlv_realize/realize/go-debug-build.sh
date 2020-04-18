#!/usr/bin/env ash

# To debug with Delve, disable optimization.
# see: https://github.com/derekparker/delve/blob/master/Documentation/usage/dlv_exec.md
CGO_ENABLED=0 go build -gcflags "all=-N -l" "$@"
