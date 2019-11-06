
Designing Go services for Continuous Deployment
Gophercon AU
1 Nov 2019
Tags: go, cd, tools, gophercon

Am Laher
Lead Engineer, Vend 
@amfastic

---

@title[Begin with a Premise]

Continuous Deployment can improve the design of your code in subtle, but remarkable ways.

Note:

: Hmm, sounds good, right.
: Maybe a stretch? Let's test the premise.
: I'm going to touch on a lot of topics, but don't worry if it's too rushed to capture the details - I'm providing working samples, which you can explore in your time.
: But first, ...

---

@title[What is this, "Continuous Deployment"?]

_For_the_purposes_of_this_talk_..._

Continuous Deployment is a release strategy: any mainline change that passes testing, is automatically released into production.

Note:

: Contiuous Deployment will automatically put your working software in front of users.
: Yes, that means a push or merge to master.
: Wonderful.
: Just think about applying some this approach to your codebase - a calm and peaceful thought, right.

---

@title[Calm & Peaceful, right?]

.image https://media.giphy.com/media/5b5OU7aUekfdSAER5I/giphy.gif

OK, OK , it does sound a little bit super-risky.
Let's explore...

Note:

: Yup, there's some things to get right:
: - there's work to be done up-front,  some work culture implications, and some caveats.
: Let's explore ...

---

@title[Let's make up a story ...]

Note:

: I'm going to try to share my own learnings, with a story about gophers.

---

@title[Toasties Galore]

Meet our heroine, Gita the Gopher.

.image ./diagrams/toasties-galore-bw.jpg _ 600

_Gita_is_an_innovator,_with_a_passion_for_technology_and_melted_cheese._

Note:

: This is her food stall, Toasties Galore.
: She's starting a business, and she's going to beat the competition with advanced cheese technology.

---

@title[State of the art Techmology ]

Gita builds her own machinery, it's software-driven. 

.image ./diagrams/jafflr-chillybin-bw.jpg _ 800

We - *you*and*I* - build the software.

Note:

: Gita is a hardware freak.
: She's also tech savvy and she advised us to use CD. She knows her cheese.
: - Chillybin is where we store the ingredients. A chilly bin (NZ English), also known as an Esky in Australia
: - Jafflr is a state of the art grilling machine.
:   - It harnesses laser technology to melt the cheese just right
:   - and a mechanical arm to fetch ingredients from the chillybin

---

@title[Toasties Galore software]

An overview of the codebase:

- A monorepo.
- 2 HTTP services, jafflr+chillybin, which interact.
- Tests.
- Docker things.
- Scripts and tooling.
- A CI/CD pipeline.

See [[https://github.com/laher/toasties-galore]]

Note:

: For this presentation, let's look at the software supporting our imaginary food stall, "Toasties Galore".
: Monorepo - multiple services in one code repository.
: There's a lot in there, more than I could hope to communicate in 30 minutes. Take a look. Toasties galore.

---

@title[The Toasties Galore software architecture]

Invoking Jafflr: `POST /toastie`

.image diagrams/toasties.dot.svg _ 600

This simple deployment powers a 24-hour, high-throughput toasting service, with some agressive expectations. _Gophers_are_demanding._

Note:

: Gita sends an HTTP request via a basic UI, ... 
: The main thing to note here is that there are contracts here between Jafflr and Chillybin (and Chillybin & Postgres)

---

@title[The CI/CD Pipeline]

.image ./diagrams/pipeline.dot.svg

Note:

: Here's the pipeline as a flow chart.
: We're using github, travis and dockerhub, but you can substitute any number of alternatives.
: Changes come into github via a merge or push, and github notifies Travis via a webhook.
: The top box represents *any change* on any branch. Maybe a feature branch, where you develop something experimental.
: The bottom box represents the deployment bit - master only. If everything builds successfully, then it'll go directly to prod.

---

@title[The pipeline]

In this case, we're using [[https://travis-ci.org/laher/toasties-galore/branches][Travis]]

.code ./.travis.yml

Note:

: Note how we have 'stages'. The main 'scripts', plus a 'deploy' section
: If you fork and PR my code, it'll trigger a build. 
: Try it out. Hell, set up a travis on your own fork.
: Click on travis to see some builds.

---

@title[Eek, auto-deploy is scary]

Note:

: What are you thinking?
: You liked it until there was a 'deploy' section? Fair.
: We're gonna need some things for this to work

---

@title[The Baseline]

Let's consider some Important Things in our code:

- Observability
- Seamless Deploys
- Integration Tests

Note:

: What's the bare minimum to start feeling more comfortable?
: Observability: What's going on in production?
: Seamless deploys: If I can't gatekeep _when_ we release, I at least want to know that releases won't interrupt our service.
: Testing: We already have tests in our services, right. But now I NEED to know they talk nicely to one another, without staging and manual tests.
: So, let's look inside Jafflr

---

@title[Jafflr schematic]

.image ./diagrams/jafflr-schematic-diagram.jpg _ 700

Note:

: OMG Gita. That's one of her drawings she wanted you to see.

---

@title[Jafflr]

.play ./jafflr/main.go /^func main/,/^}/

Note:

: So, here's the main func in jafflr. Things to notice:
: It's small
: some config, some wiring and start some things.
: In bold, some setup relating to the 3 things we were talking about

---

@title[Observability: Instrumentation]

Use metrics to know what's happening in production. e.g. in a middleware ...

.code -edit ./tpi/http.go /^func TracingMiddleware/,/^}/

Note:

: The bulk of your observability should be provided by your platform, but you'll probably want to add some instrumentation itno your code.
: Middlewares are a typical place you might add some instrumentation.
: Also, jobs, startup, other timings.
: Middleware is the name given to a wrapper around a http handler or group of handlers.
: _There's_also_ [[https://github.com/laher/toasties-galore/blob/059053a/tpi/http.go#L65:L73][another middleware]] _for_logging._

---

@title[Seamless Deploys: The Rolling Update]

Stop listening, keep handling ... e.g.

.code -edit ./tpi/http.go /^func GracefulShutdown/,/^}/

Continues processing in-flight requests, while making way for a newer deployment.

Note:

: Your platform (e.g. ECS/K8s/...) will hopefully support some kind of load-balancing and rolling update process.
: When you receive an interrupt from a 'kill' or similar:
: Exit as quickly as you can without breaking unfinished requests.
: i.e. Stop listening on the port, while continuing to fulfil existing connections. 
: Exit as soon as your requests are completed OR once the context is done (timed out).
: Allows other [new] instances to listen promptly, keeping the release snappy.
: Note: similar approach is possible for non-http services - not shown here.

---

@title[The Arm - An interaction]

.image ./diagrams/jafflr-arm.jpg _ 600

Note:

: Next: let's look at that interaction

---

@title[Integration Tests]

Few tests - verify the interactions.

.code -edit ./integration-tests/integration_test.go  /^func TestHappyPath/,/^}/

Note:

: Run tests against a 'prod-like' environment (could be docker-compose, k8s, ...)
: Things to note:
: Invoke outer service (jafflr), assert response and inner service state (chillybin)
: Doesn't assume state. Reset at the start, Check state before & after method-under-test
: Test is independent from environment. With appropriate 'demo customer' constraints, I could equally run this on prod, local, staging ...
: Needs to be non-racy, non-flaky. Preferably, fast
: Few tests - "Happy Path and some error handling". Permutations are for service/unit tests
: Yay, now we know that our services will play nicely. Halfway there.

---

@title[Success]

Toasties Galore is now a highly successful operation. Gophers *love* the recipes. Jafflr is busy 24-7.

.image ./diagrams/toasties-galore-popular.jpg _ 600

Note:

: See how much happier Gita looks. She's so excited.

---

@title[But then ...]

Some gophers aren't happy ...

- We're super busy, and some orders are going to the wrong gophers. 
- Gwenda got pineapple instead of vegemite.
- Gary received cheese instead of Sheese (soy cheese).

Not good.

---

@title[Design Problems - What to do?]

Decisions ...

*Gita:* Track orders with customers' names.

*Us:* alter the Chillybin API to support richer request data.

Note:

: We need to add new business logic.
: We got some design flaws, we want to fix them too.

---

@title[Planned changes to /pick]

.image diagrams/before.dot.svg _ 600

.image diagrams/after.dot.svg _ 600

Note:

: It's a POST instead of a GET. Imagine there's a lot more metadata; also, it really shouldn't be a GET. It's removing ingredients.
: We're inserting a new entity, orders.
: Sweeet.

---

@title[Oh noes, Breaking Changes]

- API signature needs to change.
- DB structure changes.

.image https://media.giphy.com/media/WrK9dwj8TNPr2/giphy.gif

Note:

: I'm scared, Rick.

---

@title[Us: OMG]

*Gita*: s'all good. Let's dance.

---

@title[The Foxtrot]

So, let's dance.

.image https://media.giphy.com/media/BiA0154sSljkA/giphy.gif

A 3-step dance: 

- Add new implementation.
- Migrate comfortably.
- Delete old implementation.

Note:

: I chose the foxtrot. It could have been the waltz or something. I dunno, I just googled 3-step dance. Foxtrot.

---

@title[The 3-step release dance]

Three steps, or _stages_ 

.image ./diagrams/pulls.png _ 600

With CD, it's easy, you'll consider it more often.

Note:

: We just lined up 3 changes, bottom to top. 
: Each stage may take some time before we're ready to move on to the next.
: Or not. The main thing is that we make them distinct, and follow the sequence.

---

@title[Step 1 - Add new feature implementation]

- Serve the new approach under a different URL, /v2/pick
- Now there are 2 versions of `pick`, side-by-side.
- v2 can be tested manually, but isn't invoked by any other service yet.
- No changes to calling service (jafflr). Just chillybin.

See [[https://github.com/laher/toasties-galore/pull/1/files]]

Note:

: Note the changes:
: refactored the 'update' to an 'extracted method'
: new route, /v2/pick
: an additional integration test

---

@title[Step 2 - Invoke new implementation]

- Use Feature Flags?

	func HasFeature(customer string, feature string) bool {
		switch feature {
		case "pick.v2":
			if customer == "gita" {
				return true
			}
			return false
		default:
			panic("unknown feature")
		}
	}

See [[https://github.com/laher/toasties-galore/pull/2/files]]

- HasFeature could be backed by a service, a db table, or by configuration.
- Alternatives: Canary Deploys, Shadowing, Blue-Green deploys

Note:

: We have a new PR
: Now, jafflr *conditionally* invokes Chillybin
: Switching this on and off is super chill.

---

@title[Step 3 - Delete old implementation]

- Do it soon.
- Line it up in advance. Plan for it, otherwise you'll forget.

See [[https://github.com/laher/toasties-galore/pull/3/files]]

It's the easy part, but it's also easy to forget it, or fear it. 

---

@title[Foxtrot - Conclusions]

When deployment is a non-event:

- you can make your releases safer.
- you code remains cleaner.

When you have more viable options within easy reach, you're empowered to make better decisions.

---

@title[Tips and Tricks]

---

@title[On Configuration ]

Deprecate extensive config files...

- Configurable-everything is a popular pattern
- Makes it easier to reconfigure while an on-call emergency?
- Knobs and twiddles are nice.

---

@title[But, with CD, we can simplify ]

- Use $ENV variables for DSNs and such. [see 12-factor Apps]
- Use constants more often.

Deployment is fast. Just edit code/artifact

Note:

: - Where possible, compile it in - simplify your code.
: - In other cases, it can live in the artifact (file in docker image).

---

@title[On Monorepos]

Monorepos offer some challenges and some advantages

- Benefit: easier to maintain CI/CD pipelines centrally.
- Benefit: in-house dependencies are all local.
- Benefit: testing component combinations is easier.
- Challenge: identifying what to build, test, deploy on each commit.
- Challenge: speed.
- Challenge: 3p deps: _See toasties-galore for top-level go-mod w Docker_

Note:

: Monorepos are source code repositories containing multiple services

---

@title[Deploying services in a monorepo]

Q: Do we really redeploy everything on every change?
A: Depends on language and tooling available. 

Note:

: Care needed to ensure we trigger builds exactly when we should.

---

@title[Publish / deploy monorepo service - Go]

See [[./scripts/publish.sh][publish.sh]] and [[./scripts/deploy.sh][deploy.sh]]

.code -edit ./scripts/publish.sh
 
---

@title[Last relevant commit - Go]

Calculate changes affecting a service [in a monorepo]

Note:

: In a monorepo, it's real useful to know which services are affected by a commit.
: Or put another way, which services have been affected since current production.
: * Current version
: Check 'current version of a service'.

.code ./scripts/last_commit.sh

Note:

: - calls `go list -deps` and eventually `| git rev-list`
: - Allows us to run regression tests at the right time.
: - Report that the deploy was really successful.

---

@title[Why last_commit.sh ?]

- Deploy only affected services.
- Decide which integration tests to run, post-deploy

Note:

: I've previously scripted up version checks against k8s, to determine when to run tests. 
: To me this kind of thing is crucial for regression testing.

---

@title[On DB Migrations]

Migrations will always be challenging 

- Altering big tables - time-consuming, performance concerns.
- Rollbacks are awkward.

Where possible:

- Use transactions (Postgres has transactional DDL) to *test* them in advance.
- Make them idempotent (`if not exists`) to *run* them in advance.
- Split them up - 'create table' is cheap. Isolate the tricky stuff.

Note:

: So, let's take a look at chillybin

---

@title[Chillybin]

.image ./diagrams/chillybin.jpg _ 500

Note:

: This is the original patent for the chillybin

---

@title[Chilly bin - migrations]

.code ./chillybin/main.go /^func runMigrationsSource/,/^}/

Note:

: This is using go-migrate
: Can run this from func main() but I can understand why you wouldn't
: Alternatively, run it manually in a 'job'
: Deploy the migration scripts as part of the container

---

@title[On DIY Tooling]

At Vend we have a slackbot, @overlord, which offers a number of CD-related utilities. e.g.

- canary deploys
- rollbacks
- feature flagging
- service status

Note:

: these are the kinds of tools which you may end up writing the code for.

---

@title[On CI/CD Software]

They all work, and they're all different. I've tried heaps. 

_My_examples_are_only_chosen_to_communicate_some_ideas._

Just choose one.

---

@title[Conclusions]

---

@title[How will it change my code?]

_More_than_anything_else,_CD_affects_ *how* _you_ *change* _your_code._

- CD forces you to deal with Important Things, up front. Most of these have their own payoff.
- CD encourages good practices.
- CD enables safe and easy changes.

Impact:

- Simpler: fewer contraptions in the code.
- Tidyier: fewer cobwebby corners.
- More dynamic: change is simpler -> more, better options.

---

@title[Advice]

It can be hard to retrofit CD. Start early or add it incrementally.

For research and business speak, see [[https://continuousdelivery.com]].

Note:

: Touted Benefits
: Continuous Deployment will empower you to make better choices, and encourage you to develop good practices early on.
: CD means [[https://continuousdelivery.com][reduced risk, faster time-to-market, higher quality, better products, *happier*teams*]]...

---

@title[Advice: Play with Toasties Galore]

See Toasties Galore on [[https://github.com/laher/toasties-galore][github]] or [[https://travis-ci.org/laher/toasties-galore][travis]] to try this stuff out - maybe it'll help you.

Note:

: Thanks
: I hope that helped encourage you along this journey.
: If you don't have it early on, you may have a lot more work to do, but each of these efforts has a payoff.
: I've just tried to show you some go-specific hints about the CI/CD toolchain.
: I skipped whole areas - CI systems, rollbacks, cultural obstacles, ...
: I didn't even try to recommend a CI system
: It can take some time to redesign your software for Continuous Deployment - after seeing this talk, I hope you’ll begin or complete that journey.
