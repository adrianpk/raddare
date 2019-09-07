#!/bin/sh
# Build
# ./scripts/build.sh

# Free ports
killall -9 raddare

# Set environment variables
REV=$(eval git rev-parse HEAD)
# Service
export RDR_SVC_NAME="raddare"
export RDR_SVC_REVISION=$REV
# OSRM
export RDR_OSRM_HOST="router.project-osrm.org"
export RDR_OSRM_API_VER="v1"
export RDR_OSRM_REQ_TIMEOUT_SECS="5"

go build -o ./bin/raddare ./cmd/raddare.go
./bin/raddare
# go run -race cmd/raddare.go
