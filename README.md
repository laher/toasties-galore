# toasties-galore 

Presentation about Continuous Delivery for Gophercon AU

## Approaches detailed

 * Table-driven integration tests; fixtures, assertions, chaining together calls:
  - [ ] Using httptest or go-mocks - example provided 
 * Managing database definitions with ‘migrations’:
  - [X] go-migrate. See [chillybin](./chillybin/main.go)
 * Bundling resources with your app (Docker or go-bindata):
  - [ ] docker in this case
 * Mitigating risk of changes:
  - [ ] Feature Flags implemented via config, shared config, and db
  - [ ] Versioning your APIs 
 * Zero-downtime deploys:
  - [X] graceful restart. See [sandwich-press](./sandwich-press/main.go)
  - [ ] healthchecks
 * Metrics and alerting
 * Tooling for your build system: 
  - [ ] Containerising build steps 
  - Deployment and confirmation (ECS/K8S) 
  - [X] `go list -deps` for granular version checking. See [last_commit.sh](./scripts/last_commit.sh)
