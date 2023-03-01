#!/bin/bash
VERSION=v2.2.8
COMMIT_HASH=$(git rev-parse --short HEAD)
ENVIRONMENT=${VARIABLE:-"dev"} 
BUILDSTRING="$VERSION-$COMMIT_HASH-$ENVIRONMENT"
echo Using the "$ENVIRONMENT" environment
echo Building uberbot "$VERSION-$COMMIT_HASH"
if [[ "$ENVIRONMENT" == "dev" ]]; then
  go build -gcflags="all=-N -l" -ldflags "-X github.com/ubergeek77/uberbot/v2/core.VERSION=${BUILDSTRING} -X github.com/ubergeek77/uberbot/v2/core.ENVIRONMENT=${ENVIRONMENT}";
else
  go build -ldflags "-X github.com/ubergeek77/uberbot/v2/core.VERSION=${BUILDSTRING} -X github.com/ubergeek77/uberbot/v2/core.ENVIRONMENT=${ENVIRONMENT}";
fi
