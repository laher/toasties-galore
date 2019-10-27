SVC_PATH = $(shell pwd)
BASE_PATH = $(dir $(shell pwd))
TAG=$$(../scripts/last_commit.sh .)
IMAGE="$(SVC)"

print-tag: ## print current version (commit hash)
	@echo "TAG: $(TAG)"

docker-build: ## build docker image
docker-build: print-tag
		${SUDO} docker build -t "$(IMAGE)" -f $(SVC_PATH)/Dockerfile $(BASE_PATH)
		${SUDO} docker tag "$(IMAGE):latest" "$(IMAGE):$(TAG)"

build-linux: ## Build (linux) binary
	GO111MODULES=on CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a .

run: ## Run locally
	GO111MODULES=on go run -mod=vendor .

reflex: ## Run (file watcher)
	reflex -s -g '*.go' -- sh -c "$(MAKE) run"

reflex-nohup: ## Run (file watcher)
	nohup $(MAKE) reflex &

help:
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
