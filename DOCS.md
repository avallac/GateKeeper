# GateKeeper

Bridges an IPGuard intercom to Home Assistant: opens doors and surfaces incoming-call
notifications (with a snapshot) as MQTT-discovered entities.

## First-time setup

1. Install the official **Mosquitto broker** add-on first (Settings → Add-ons → Add-on
   Store), if you don't already have an MQTT broker. GateKeeper auto-detects its
   credentials — no manual entry needed.
2. Start this add-on, then open its panel from the Home Assistant sidebar.
3. Enter your name and phone number; you'll receive an SMS code to pair the app with
   your IPGuard account.
4. Once paired, your doors/intercoms appear under **Settings → Devices & Services →
   MQTT** as a device with a door-open button, a call event, and a call snapshot
   camera.

## Configuration

Open the add-on's **Configuration** tab to set `auth_address`/`main_address` if your
IPGuard provider uses non-default servers.

Pairing state (phone/code/key) is stored under the add-on's persistent `/data` volume,
so it survives add-on updates and restarts.
