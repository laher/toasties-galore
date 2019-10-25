#!/bin/bash

echo "deploy"

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

function docker_tag_exists() {
    curl --silent -f -lSL https://index.docker.io/v1/repositories/$1/tags/$2 > /dev/null
}

export JAFFLR_COMMIT=$(./scripts/last_commit.sh jafflr)
echo "checking jafflr tag $JAFFLR_COMMIT"
if docker_tag_exists "$DOCKER_USERNAME/jafflr" "$JAFFLR_COMMIT"; then
    echo "Tag exists for $DOCKER_USERNAME/jafflr: $JAFFLR_COMMIT"
else
    docker build -t jafflr -f ./jafflr/Dockerfile .
    docker tag jafflr $DOCKER_USERNAME/jafflr
    docker tag jafflr $DOCKER_USERNAME/jafflr:$JAFFLR_COMMIT
    docker push $DOCKER_USERNAME/jafflr
    echo "chillybin deploy goes here"
fi

export CHILLYBIN_COMMIT=$(./scripts/last_commit.sh chillybin)
echo "checking chillybin tag $JAFFLR_COMMIT"
if docker_tag_exists "$DOCKER_USERNAME/chillybin" "$CHILLYBIN_COMMIT"; then
    echo "Tag exists for $DOCKER_USERNAME/chillybin: $CHILLYBIN_COMMIT"
else
    docker build -t chillybin -f ./chillybin/Dockerfile .
    docker tag chillybin $DOCKER_USERNAME/chillybin
    docker tag chillybin $DOCKER_USERNAME/chillybin:$CHILLYBIN_COMMIT
    docker push $DOCKER_USERNAME/chillybin
    echo "chillybin deploy goes here"
fi
