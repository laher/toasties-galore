#!/bin/bash

echo "BUILD/PUBLISH: begin"

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

function docker_tag_exists() {
    echo "CHECK: does $1 docker tag exist for commit $2 ?"
    curl --silent -f -lSL https://index.docker.io/v1/repositories/$1/tags/$2 > /dev/null
}

for SVC in jafflr chillybin
do
    export SVC_COMMIT=$(./scripts/last_commit.sh "$SVC")
    if docker_tag_exists "$DOCKER_USERNAME/$SVC" "$SVC_COMMIT"; then
        echo "SKIP: Docker tag already exists for $DOCKER_USERNAME/$SVC: $SVC_COMMIT"
    else
        echo "PUBLISH: publishing image for $DOCKER_USERNAME/$SVC: $SVC_COMMIT"
        docker build -t $SVC -f ./$SVC/Dockerfile . # rebuild just in case
        docker tag $SVC $DOCKER_USERNAME/$SVC
        docker tag $SVC $DOCKER_USERNAME/$SVC:$SVC_COMMIT
        docker push $DOCKER_USERNAME/$SVC
    fi
done

echo "BUILD/PUBLISH: done"
