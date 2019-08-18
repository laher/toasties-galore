# toasties-galore 

Presentation about Continuous Delivery for Gophercon AU

## Approaches detailed

 * [ ] table-driven integration tests; fixtures, assertions, chaining together calls:
  - [ ] Using httptest or go-mocks - example provided 
 * [ ] Managing database definitions with ‘migrations’:
  - [ ] go-migrate
 * [ ] Bundling resources with your app (Docker or go-bindata):
  - [ ] docker in this case
 * [ ] Mitigating risk of changes:
  - [ ] Feature Flags implemented via config, shared config, and db
  - [ ] Versioning your APIs 
 * Zero-downtime deploys:
  - [ ] graceful restart
  - [ ] healthchecks
 * Metrics and alerting
 * Tooling for your build system: 
  - [ ] Containerising build steps 
  - Deployment and confirmation (ECS/K8S) 
  - [ ] `go list -deps` for granular version checking
