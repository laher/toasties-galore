#!/bin/bash

echo "BUILD/DEPLOY: begin"

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

function docker_tag_exists() {
    echo "CHECK: does $1 docker tag exist for commit $2 ?"
    curl --silent -f -lSL https://index.docker.io/v1/repositories/$1/tags/$2 > /dev/null
}

function is_deployed_version() {
    # check k8s/ecs/... here
    # e.g. aws blah check
    # e.g. kubectl check thing
    echo "CHECK: is $1 deployed with version $2? checking K8S/ECS/... for currently-running version"
    return 1
}

for SVC in jafflr chillybin
do
    export SVC_COMMIT=$(./scripts/last_commit.sh "$SVC")
    if docker_tag_exists "$DOCKER_USERNAME/$SVC" "$SVC_COMMIT"; then
        echo "SKIP: Docker tag already exists for $DOCKER_USERNAME/$SVC: $SVC_COMMIT"
    else
        echo "BUILD: building image for $DOCKER_USERNAME/$SVC: $SVC_COMMIT"
        docker build -t $SVC -f ./$SVC/Dockerfile .
        docker tag $SVC $DOCKER_USERNAME/$SVC
        docker tag $SVC $DOCKER_USERNAME/$SVC:$SVC_COMMIT
        docker push $DOCKER_USERNAME/$SVC
    fi

    if is_deployed_version "$SVC" "$SVC_COMMIT"; then
        echo "SKIP: Service $SVC version is already deployed to version: $SVC_COMMIT"
    else
        echo "DEPLOY: Service $SVC deploy to version '$SVC_COMMIT' goes here"
    fi
done

echo "BUILD/DEPLOY: done"
