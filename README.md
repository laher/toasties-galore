# toasties-galore 

Code and resources for a talk about Continuous Deployment for Gophercon AU [![Build Status](https://travis-ci.org/laher/toasties-galore.svg?branch=master)](https://travis-ci.org/laher/toasties-galore)

## Monorepo

This is a small monorepo with a naive implementation of some HTTP APIs, and some build tooling.

For the purposes of the talk, there are also 3 Pull Requests open for discussion.

## Services

The jafflr is an imaginary machine which assembles and cooks toasted sandwiches.

It attaches onto a chilly bin ('Chilly bin' is NZ English parlance - aka Esky, in Australian).

There are 2 services, `jafflr` and the `chillybin`. Crucially, `jafflr` depends on `chillybin`.

## Endpoints

 * chillybin:
  - [X] `/pick?name=cheese&quantity=1` - pick stock from chillybin
  - [X] `/` - show current stock
  - [X] `/restock` - restock
  - [X] `/health` - health check
 * jafflr:
  - [X] `/toastie?i=cheese&i=vegemite` - make a toastie 
   * this invokes chillybin 
  - [X] `/` retrieves status (toasting/available)
  - [X] `/health` - health check

## Approaches detailed

 * Table-driven integration tests; fixtures, assertions, chaining together calls:
  - [ ] Using httptest - example provided 
 * Managing database definitions with ‘migrations’:
  - [X] go-migrate. See [chillybin](./chillybin/main.go)
  - With mongodb this can be as simple as "EnsureIndex()":w
 * Bundling resources with your app (Docker or go-bindata):
  - [X] docker in this case
 * Mitigating risk of changes:
  - [ ] Feature Flags implemented via ENV or db
    - 3 steps:
      1. Release something cross-compatible.
      2. Switch client over to the new API
      3. Delete old endpoint (for 3rd parties this gets more complicated)
  - [ ] Versioning your APIs 
      1. Could be a path `/v2/`
      2. Could be version headers
      3. GraphQL supports `deprecated fields`
      4. gRPC, Twirp [Protocol Buffers] support deprecations and field renaming
 * Zero-downtime deploys:
  - [X] HTTP - graceful restart. See [jafflr](./jafflr/main.go)
  - [X] healthchecks
 * Metrics and alerting
 * Tooling for your build system: 
  - [X] Containerising build steps 
  - Deployment and confirmation (ECS/K8S) 
  - [X] `go list -deps` for granular version checking. See [last_commit.sh](./scripts/last_commit.sh)
