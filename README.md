# toasties-galore 

Presentation about Continuous Delivery for Gophercon AU

## Services

The Jafflotron is an imaginary machine which assembles and cooks toasted sandwiches.

It attaches onto a chilly bin ('Chilly bin' is NZ English parlance - aka Esky, in Australian).

There are 2 services, `jafflotron` and the `chillybin`. Crucially, `jafflotron` depends on `chillybin`.

## Endpoints

 * chillybin:
  - [X] `/pick?name=cheese&quantity=1` - pick stock from chillybin
  - [X] `/` - show current stock
  - [ ] `/restock` - restock
  - [X] `/health` - health check
 * jafflotron:
  - [ ] `/toastie?filling[]=cheese&filling[]=vegemite` - make a toastie 
   * this invokes chillybin 
  - [ ] `/` retrieves status (toasting/available)
  - [X] `/health` - health check

## Approaches detailed

 * Table-driven integration tests; fixtures, assertions, chaining together calls:
  - [ ] Using httptest - example provided 
 * Managing database definitions with ‘migrations’:
  - [X] go-migrate. See [chillybin](./chillybin/main.go)
 * Bundling resources with your app (Docker or go-bindata):
  - [X] docker in this case
 * Mitigating risk of changes:
  - [ ] Feature Flags implemented via config, shared config, and db
  - [ ] Versioning your APIs 
 * Zero-downtime deploys:
  - [X] graceful restart. See [sandwich-press](./sandwich-press/main.go)
  - [ ] healthchecks
 * Metrics and alerting
 * Tooling for your build system: 
  - [X] Containerising build steps 
  - Deployment and confirmation (ECS/K8S) 
  - [X] `go list -deps` for granular version checking. See [last_commit.sh](./scripts/last_commit.sh)
