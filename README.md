# GateKeeper

![GateKeeper logo](logo.png)

GateKeeper bridges an **IPGuard** intercom system to **Home Assistant**. It talks to the
IPGuard cloud over ZeroMQ, pairs with your account via SMS, and exposes your doors and
intercoms to Home Assistant as [MQTT Discovery](https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery)
entities:

- a **door-open button** for each relay/door
- an **incoming-call event** entity
- a **call snapshot camera** entity

It ships both as a Home Assistant Supervisor add-on and as a standalone Go service/Docker
image.

## Legal Notice

### Independent Project

GateKeeper is an independent, community-developed open-source project.

It is not affiliated with, associated with, authorised by, endorsed by, or supported by IP Door Entry Ltd, Noralsy, BATICONNECT, IPGUARD, or any of their affiliates.

This project is not an official IPGUARD client or an officially supported integration.

### Interoperability

GateKeeper is independently developed software intended to provide interoperability with IPGUARD-compatible intercom systems.

The project does not include or distribute proprietary source code, binaries, or application assets from the official IPGUARD applications.

All third-party trademarks, product names and service names are the property of their respective owners and are used solely to identify compatible systems and services.

### Authorised Use Only

GateKeeper is intended to be used only with IPGUARD accounts, devices and resources that the user is authorised to access.

Users are responsible for ensuring that their use of GateKeeper complies with applicable laws, access permissions and contractual obligations.

The software must not be used to bypass access controls, obtain unauthorised access to accounts or devices, or interfere with the operation of third-party infrastructure.

### Physical Access and Security

GateKeeper can issue commands that may result in physical doors or access-controlled entry points being unlocked.

Users are solely responsible for the configuration, deployment and security of their installations, including MQTT access controls, authentication, network security and the protection of credentials.

The software should not be configured to grant access automatically without appropriate safeguards.

### Disclaimer of Warranty and Liability

GateKeeper is provided "AS IS", without warranties or guarantees of any kind, express or implied.

The authors and contributors make no guarantees regarding availability, reliability, compatibility, security or correct operation of the software.

To the maximum extent permitted by applicable law, the authors and contributors shall not be liable for any damages or losses arising from the use or inability to use the software, including unintended door unlocking, unauthorised physical access, service interruptions or loss of data.

Nothing in this notice excludes or limits liability where such exclusion or limitation is prohibited by applicable law.

## How it works

1. GateKeeper connects to the IPGuard auth/main servers over ZeroMQ (`goczmq`).
2. On first run you pair the app with your IPGuard account through a small web UI: enter
   your name and phone number, receive an SMS code, and GateKeeper exchanges it for a
   long-lived key.
3. Once paired, GateKeeper publishes your doors/intercoms to Home Assistant via MQTT
   Discovery, and relays incoming-call events (with a snapshot) and door-open commands
   between Home Assistant and IPGuard.

Pairing state (phone/code/key) is persisted to disk so it survives restarts and updates.

## Running as a Home Assistant add-on

1. Install the official **Mosquitto broker** add-on first (Settings → Add-ons → Add-on
   Store) if you don't already have an MQTT broker. GateKeeper auto-detects its
   credentials from the Supervisor — no manual entry needed.
2. In Home Assistant, go to **Settings → Add-ons → Add-on Store → ⋮ → Repositories** and
   add:

   ```
   https://github.com/avallac/GateKeeper
   ```

   Then find **GateKeeper** in the store and install it.
3. Start the add-on, then open its panel from the Home Assistant sidebar.
4. Enter your name and phone number; you'll receive an SMS code to pair the app with
   your IPGuard account.
5. Once paired, your doors/intercoms appear under **Settings → Devices & Services →
   MQTT** as a device with a door-open button, a call event, and a call snapshot
   camera.

### Add-on configuration

Open the add-on's **Configuration** tab to set `auth_address` / `main_address` if your
IPGuard provider uses non-default servers. See [`build.yaml`](build.yaml) and
[`repository.yaml`](repository.yaml) for add-on metadata.

## Running standalone

### Requirements

- Go 1.25+
- A C toolchain plus `libzmq3-dev` and `libczmq-dev` (or equivalents) for the CGO-based
  `goczmq` dependency
- An MQTT broker reachable from the host

### Configuration

GateKeeper reads configuration in layers, each overriding the previous:

1. [`configs/app.yaml`](configs/app.yaml) — base config (`authAddress`, `mainAddress`,
   `addr`)
2. [`configs/app.local.yaml`](configs/app.local.yaml) — local overrides, e.g. pointing
   `authAddress`/`mainAddress` at a local [`cmd/fake-server`](cmd/fake-server) instance
   via `fakeServer: true`
3. When running under the Home Assistant Supervisor (`SUPERVISOR_TOKEN` is set):
   environment variables (`MQTT_HOST`, `MQTT_PORT`, `MQTT_USERNAME`, `MQTT_PASSWORD`,
   `MQTT_SSL`), `/data/options.json` (add-on options), and the Supervisor's MQTT service
   discovery

Pairing state is stored in [`configs/user.yaml`](configs/user.yaml) (or `/data/config.yaml`
under the Supervisor).

### Build & run

```sh
make build   # builds ./bin/gatekeeper (requires CGO)
make run     # go run ./cmd/app
make test    # go test ./...
make fmt     # go fmt ./...
make vet     # go vet ./...
```

The web UI is served on `:8080` by default (`addr` in `configs/app.yaml`).

### Using the fake server

[`cmd/fake-server`](cmd/fake-server) stands in for the real IPGuard auth/main servers, so
you can run and test GateKeeper without a real intercom. To switch the client to it:

1. Start the fake server:

   ```sh
   go run ./cmd/fake-server
   ```

   It binds the auth/main ZeroMQ endpoints on `tcp://*:5579` and `tcp://*:5577`, and
   serves a small control UI on `:8081`.

2. Point GateKeeper at it via `configs/app.local.yaml`:

   ```yaml
   authAddress: "tcp://localhost:5579"
   mainAddress: "tcp://localhost:5577"
   fakeServer: true
   ```

3. Run GateKeeper as usual (`make run`). With `fakeServer: true`, any previously stored
   pairing key is discarded on startup, so the app re-pairs against the fake server. This
   re-pair is needed to obtain a new encryption key for the main channel, since the fake
   server generates a fresh key every time it starts.

### Docker

```sh
make docker-build
make docker-run
```

`docker-run` mounts `configs/user.yaml` read-write (for persisted pairing state) and
`configs/app.local.yaml` read-only into the container.

## Development

- [`cmd/app`](cmd/app) — the GateKeeper service entry point
- [`cmd/fake-server`](cmd/fake-server) — a mock of the IPGuard auth/main servers for
  development and testing
- [`internal/service/ipguard`](internal/service/ipguard) — IPGuard protocol client and
  message types
- [`internal/service/mqtt`](internal/service/mqtt) — MQTT client and Home Assistant
  discovery payloads
- [`internal/handlers`](internal/handlers) — HTTP handlers for the pairing/dashboard web
  UI
- [`web`](web) — HTML templates and static assets for the web UI

Run the unit tests with `make test` (or `go test ./...`).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
