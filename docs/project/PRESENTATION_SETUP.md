# Isolated FSMO presentation setup

**Current owner acceptance and publication authority, 2026-10-11 — DEC-085:** Final Staff/Admin desktop and Student/Faculty mobile visual review passed; the owner reported no remaining UI issues and authorized auditing, committing and pushing the completed Phases 8–14 application and shared UI to `release/elabtrack-v2-presentation-ready`. Main and production deployment remain untouched. Normal schema 000008 and existing records are preserved. Official FSMO terms/current acceptance, production Student domains, verified activation/Brevo and operational deployment readiness remain separate gates. Earlier pending-review/no-commit statements below are historical task snapshots. [Release audit and publication handoff](RELEASE_PUBLICATION_REPORT.md).

This is fictional local rehearsal data under DEC-082. Final owner acceptance of Phases8–14 remains pending. The accepted Phase7 demo is retained separately. Do not point the new API at the normal database, which remains migration000008 until a separately authorized rollout.

|Environment|PostgreSQL|API|Browser|
|---|---|---|---|
|Normal development, preserved|5434|8080|5173|
|Accepted Phase7 demo, preserved|54835|18086|15176|
|Full-system presentation|54836|18087|15177|
|Disposable engineering verification|54832|18085|15175|

From the repository root in WSL Ubuntu, with Docker Desktop WSL integration or a local Docker daemon running:

```bash
python3 integration/presentation-demo.py setup
python3 integration/presentation-demo.py status
python3 integration/presentation/verify.py
python3 integration/presentation-demo.py start
```

Open **http://127.0.0.1:15177/login** on the same machine. The permanent bottom banner identifies the fictional presentation and port. API readiness: http://127.0.0.1:18087/api/v1/ready. Loopback bindings do not expose this dataset to campus networks. Never confuse this with normal5173 or accepted demo15176.

Setup requires Python3.12+, Go matching backend/go.mod (1.27.1), Node24 (or another manifest-supported release), npm, Docker/Compose with `!override` support, available ports54836/18087/15177, and dependencies installed. On a new developer machine, install the toolchains and run `npm ci` in frontend and `go mod download` in backend while internet is available; fetch `postgres:18.6-alpine` before disconnecting. Tool locations come from PATH, optionally ELABTRACK_GO/ELABTRACK_NODE. APP_BIND_HOST=127.0.0.1 makes the presentation API loopback-only; the normal API default is unchanged. Startup refuses occupied18087/15177 rather than stopping another process. A fresh working copy refuses an already existing presentation container and directs you to its owning working copy. No credentials or hidden owner environment file is required: setup generates private separate database/JWT/account secrets and initializes the fresh labeled volume. Runtime receives only its limited database credentials; the migrator is separate. Source is copied from an explicit public allowlist into ignored backend/tmp/presentation-demo/snapshot, with hashes. Historical migrations1–8 must match acceptedc2741c5; setup then applies paired1–11. Source changes require a deliberate build and restart, never reseeding an existing target.

Credentials are generated, hashed with existing bcrypt, and stored only in ignored0600 private files. Retrieve them **in your own interactive terminal**:

```bash
python3 integration/presentation-demo.py credentials
```

Use `admin`, `staff1`, `staff2`, `student01`…`student20`, or `faculty01`…`faculty05`. Do not paste passwords into chat, shell arguments, reports or source. No credential list is included in this guide. Every account uses the real email/password login and refresh system. Seeded test-ready activation is confined to this explicit disposable target; no mail delivery or production activation verification is claimed.

The database starts with40 equipment,20 Students,5 Faculty,2 Staff,1 Admin,40 protected catalog illustrations,28 fictional profile avatars,35 borrowing histories and24 explicit **fictional historical** demonstration terms acceptances. Student20 has no seeded acceptance and demonstrates the real first-use checkbox/acceptance workflow. All terms are titled **DEMONSTRATION TERMS — NOT OFFICIAL FSMO POLICY**. They never enter normal data. Artwork is original local vector illustration rendered to PNG, with source/provenance/hash manifest; it is not a photograph of real stock or students. Inventory and dashboard values derive from application transactions, not fixed UI metrics. Due-soon24h is demo configuration only; production lead/repeat remains OPEN-025.

### Rebuild the shared UI without reseeding

The DEC-083 desktop/mobile reconstruction is the same feature implementation in normal and presentation builds. The presentation adds only its permanent fictional-data banner and corresponding bottom spacing. PDF generation uses the bundled licensed font and local libraries, so no CDN is required.

After source changes on this existing presentation, retain all current data:

```bash
python3 integration/presentation-demo.py stop
python3 integration/presentation-demo.py build
python3 integration/presentation-demo.py start
python3 integration/presentation/verify.py
```

Do not run setup/reset as a shortcut to rebuilding. Prior acceptance and UI tests create legitimate synthetic borrowing/audit records and durable notification/terms state; the original seed counts describe a fresh setup. Student20 may already have accepted the demo terms after a rehearsal. Never delete acceptance/history to reproduce first use. A deliberate guarded reset is a separate operator action described below, and is not performed by this reconstruction task. The normal database still needs a separately authorized migration from 8 to the verified compatible schema 11 before full-system integration.

Stop/restart without changing data:

```bash
python3 integration/presentation-demo.py stop
python3 integration/presentation-demo.py start
```

Reset deliberately restores all scenarios and the fresh Student20 terms interaction. It verifies config/database name/environment, Compose ownership, exact volume and loopback port, and the database's seeded marker. It stops only its owned processes, retains the old volume and private config in a timestamped directory, and creates a **new** isolated volume. No volume is deleted.

```bash
python3 integration/presentation-demo.py reset --confirm elabtrack_v2_presentation_demo
python3 integration/presentation-demo.py start
python3 integration/presentation/verify.py
```

Run `setup` again after an interrupted initial attempt: it resumes only a verified empty target, and never seeds a nonempty one. Seed initialization is one transaction; fixture failures roll back users/stock/terms/histories together, while migrations and the unseeded identity marker remain for safe diagnosis. Generated credentials are written before seeding. Diagnostic logs are private; do not publish them. An unverified, partially seeded or foreign target is refused. Retained volumes/configuration provide recovery options; no automatic restoration into normal data exists.

If startup fails, check Docker, free ports, private file permissions and `status`; use private api/frontend/last-error logs locally. Do not run `docker compose down -v`, clear normal data, or weaken origin/session policy. Mismatching configuration requires correction to the recorded presentation identity rather than guessing a target. Rebuild application source only while the demo is stopped:

```bash
python3 integration/presentation-demo.py stop
python3 integration/presentation-demo.py build
python3 integration/presentation-demo.py start
```

Production remains a separate review: official FSMO terms publication/current acceptance, approved Student domains, verified Brevo activation/recovery, HTTPS/cookie/origin configuration, backup/restore procedures and normal migration authorization. Fresh production installation must not import demo users, synthetic consent or fictional stock. No V1 migration is required. Account creation without a verified sender must report pending/unavailable delivery honestly; an offline seed is not a live activation or password-recovery workflow. Use the generated demonstration accounts for the walkthrough.


Before disconnecting, prepare the pinned dependencies and Docker image. A clean machine needs downloads once; runtime does not fetch them. To make a separate offline npm cache from prepared public packages, with no copied npmrc/credentials:

```bash
node integration/presentation/prepare-npm-cache.mjs frontend/package-lock.json /path/to/prepared/npm-cache backend/tmp/presentation-offline-cache
```

In a clean working copy, run `npm ci --offline --cache ../backend/tmp/presentation-offline-cache --no-audit --no-fund` from frontend. The export reports missing platform-optional **or unprepared** tarballs; missing required ones must be prepared while online, then the offline install repeated. Never treat the export count alone as proof. The verified rehearsal installed373 packages from a separate cache into new node_modules, with distinct empty user/global npm configuration files; no owner npmrc, env, compiled binary or installed project dependencies were copied. Go's pinned dependencies were vendored before disconnecting, then API/seed built with `GOFLAGS=-mod=vendor`, `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, an empty module cache and a fresh build cache. Optional ELABTRACK_BUILD_CACHE selects that cache. The source snapshot includes a prepared backend/vendor when present. Installed toolchains and cached postgres:18.6-alpine remain documented prerequisites.

For a consistent local backup and actual restore rehearsal, stop this presentation first:

```bash
python3 integration/presentation-demo.py stop
python3 integration/presentation/recovery.py
python3 integration/presentation-demo.py start
```

Recovery verifies the marked source, creates a0600 custom PostgreSQL archive in an ignored timestamped directory and restores it into a **new** labeled54837 database/volume using its bootstrap owner. All table digests, sequence states and11 migration records must match. The restored container is stopped to release its port; both recovered volume/configuration and backup remain private and retained. This command never overwrites a database or restores normal data. Failure logs remain private. A filesystem/archive listing is insufficient proof of a backup.

An offline rehearsal blocks non-loopback browser traffic from before initial login; it does not alter Windows/WSL networking or firewall settings. Windows-host browser forwarding, physical devices, Safari/Firefox and institutional deployment are separate manual checks. If Windows cannot reach127.0.0.1:15177, first test readiness/browser URLs from WSL and verify WSL localhost forwarding; do not bind the private demo to a campus-facing interface as a workaround.

The [complete29-step walkthrough](PRESENTATION_WALKTHROUGH.md) includes a shorter route. [Production gaps](PRODUCTION_READINESS_GAPS.md) separates approved rules from remaining approvals and technical dependencies.
