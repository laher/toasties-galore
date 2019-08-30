SVC_PATH = $(shell pwd)
BASE_PATH = $(dir $(shell pwd))
TAG=$$(../scripts/last_commit.sh .)
IMAGE="toastiehub.org/$(SVC)"

print-tag: ## print current version (commit hash)
	@echo "TAG: $(TAG)"

docker-build: ## build docker image
docker-build: print-tag
		${SUDO} docker build -t "$(IMAGE):$(TAG)" -f $(SVC_PATH)/Dockerfile $(BASE_PATH)

help:
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
