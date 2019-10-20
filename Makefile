
install:
	go get github.com/cespare/reflex

run-chillybin: ## Run chillybin
	cd chillybin && $(MAKE) run

run-jafflr: ## Run jafflr
	cd jafflr && $(MAKE) run

reflex-chillybin: ## Run chillybin (file watcher)
	cd chillybin && $(MAKE) reflex

reflex-jafflr: ## Run jafflr (file watcher)
	cd jafflr && $(MAKE) reflex

start-postgres: ## start postgres
	docker-compose start postgres

run-all: ## Run all
run-all: start-postgres
	cd chillybin && $(MAKE) reflex-nohup
	cd jafflr && $(MAKE) reflex-nohup

test: ## Run tests
	CHILLYBIN_ADDR=http://localhost:7001 JAFFLR_ADDR=http://localhost:7000 go test -v ./integration-tests 

help:
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
