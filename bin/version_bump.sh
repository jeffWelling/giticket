#!/bin/bash
set -x

# Run tests
bin/run_tests.sh
if [ $? -ne 0 ]; then
    echo "Tests failed, please fix"
    exit 1
fi

#Get version number
VERSION=$(grep Version pkg/common/version.go | cut -d '=' -f 2 | sed 's/ //g' | sed 's/"//g')

MAJOR_VERSION=$(echo $VERSION | cut -d '.' -f 1)
MINOR_VERSION=$(echo $VERSION | cut -d '.' -f 2)
PATCH_VERSION=$(echo $VERSION | cut -d '.' -f 3)

# Increment PATCH_VERSION
BUMPED_PATCH_VERSION=$((PATCH_VERSION + 1))

# Write out new version file
VERSION_FILE_CONTENTS="package common"

echo $VERSION_FILE_CONTENTS > pkg/common/version.go
echo "const Version = \"$MAJOR_VERSION.$MINOR_VERSION.$BUMPED_PATCH_VERSION\"" >> pkg/common/version.go

git add pkg/common/version.go
git commit -m "Bump version to $MAJOR_VERSION.$MINOR_VERSION.$BUMPED_PATCH_VERSION and tag"

# Syntax check
gofmt -e pkg/common/version.go > /dev/null
if [ $? -ne 0 ]; then
    echo "gofmt failed"
    exit 1
fi

git tag v"$MAJOR_VERSION"."$MINOR_VERSION"."$BUMPED_PATCH_VERSION" && \
git push origin v"$MAJOR_VERSION"."$MINOR_VERSION"."$BUMPED_PATCH_VERSION"
