#!/bin/bash

imageName="http-injest"
repoName="us-east1-docker.pkg.dev/shadowtacticalautomation/sct"

# Do Go Stuff
go get -u
go mod tidy

# Determine the current version from git tags
version=$(git describe --tags --abbrev=0)
major=$(echo "$version" | cut -d. -f1)
minor=$(echo "$version" | cut -d. -f2)
patch=$(echo "$version" | cut -d. -f3)
echo "Major: $major, Minor: $minor, Patch: $patch"

np=$(($patch+1))
newTag="$major.$minor.$np"
branch=$(git rev-parse --abbrev-ref HEAD)

# Write the new version to version.txt
oldVersion=$(cat version.txt)
echo $newTag > version.txt
git commit -am "Bump version to $newTag"
git push origin $branch

# Tag and build the Docker image
git tag $newTag
git push origin $newTag
GOOS=linux GOARCH=amd64 go build
docker build -t $repoName/$imageName:$newTag .
docker push $repoName/$imageName:$newTag
rm $imageName

# # Redeploy the service
# sed -i '' "s/$oldVersion/$newTag/g" ops/prod/deployment.yaml
# cd ops/prod && ./redeploy.sh
# git commit -am "Deployed version $newTag"
# git push origin $branch
