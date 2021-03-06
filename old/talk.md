# Premise:

Continuous Delivery can drive the design of your code for Great Good. I'll share some of the tricks and strategies I've used to deploy Go services safely, smoothly, and without ceremony. I'll cover integration testing, db migrations, zero-downtime deploys, feature flags, and some useful go tooling.


Character: Gita and Toasties Galore

Conflict: breaking changes

Resolution: the 3-step release

With CD, this becomes very natural. The more ceremony you have in your release process, the harder this becomes.

Other concerns:

 * Another premise: manual deployment steps will keep increasing over time as your deployment grows and changes. Nip it in the bud.
 * Frontloading: integration testing, managing breaking changes, downtimeless deploys, observability
 * DB migrations
 * Config or not - not! ENV, in many cases. See 12-factor apps.
 * Tooling - e.g. monorepos - the build decision
 * Tooling - canary tasks
 * Observability - put version into metrics (where cardinality allows)
