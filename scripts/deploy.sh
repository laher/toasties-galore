#!/bin/bash

export ENVIRONMENT=${1:-prod}
echo "BUILD/DEPLOY: begin deploy to $ENVIRONMENT"

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

function is_deployed() {
    echo "TODO CHECK: is $1 deployed to $3 with version $2? checking K8S/ECS/..."
    return 1
}

function deploy() {
    echo "TODO DEPLOY: Service $1 deploy to $3 with version '$2' ... K8S/ECS/... "
    return 0
}

for SVC in jafflr chillybin
do
    export SVC_COMMIT=$(./scripts/last_commit.sh "$SVC")
    if is_deployed_version "$SVC" "$SVC_COMMIT" "$ENVIRONMENT"; then
        echo "SKIP: Service $SVC version is already deployed to version: $SVC_COMMIT"
    else
        echo "DEPLOY: Tag $DOCKER_USERNAME/$SVC as $ENVIRONMENT release"
        if docker inspect --type=image $DOCKER_USERNAME/$SVC:$SVC_COMMIT
        then
            echo "SKIP: image already built"
        else
            echo "BUILD: need to rebuild image"
            docker build -t $SVC -f ./$SVC/Dockerfile .
            docker tag $SVC "$DOCKER_USERNAME/$SVC:$SVC_COMMIT"
        fi
        docker tag "$DOCKER_USERNAME/$SVC:$SVC_COMMIT" "$DOCKER_USERNAME/$SVC:$ENVIRONMENT"
        docker push "$DOCKER_USERNAME/$SVC:$ENVIRONMENT" # <- in lieu of actual deploy
        deploy "$SVC" "$SVC_COMMIT" "$ENVIRONMENT"
    fi
done

echo "BUILD/DEPLOY: done deploy to $ENVIRONMENT"
