
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

test-restock: ## Run tests
	CHILLYBIN_ADDR=http://localhost:7001 JAFFLR_ADDR=http://localhost:7000 go test -v ./integration-tests -run Restock

test-burnt: ## Run tests
	CHILLYBIN_ADDR=http://localhost:7001 JAFFLR_ADDR=http://localhost:7000 DONENESS=burnt go test -v ./integration-tests -run Burnt

test-medium-done: ## Run tests
	CHILLYBIN_ADDR=http://localhost:7001 JAFFLR_ADDR=http://localhost:7000 DONENESS=medium go test -v ./integration-tests -run Burnt

dot: ## Generate dotfile image
	#unflatten -l 2 "toasties.dot" | dot -Tpng -o "toasties.png"
	cat "diagrams/toasties.dot" | dot -s144 -Tsvg -o "diagrams/toasties.svg"
	cat "diagrams/toasties-changes.dot" | dot -Tsvg -o "diagrams/toasties-changes.svg"
	cat "diagrams/pipeline.dot" | dot -Tsvg -o "diagrams/pipeline.svg"

help:
	@grep -E -h '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
