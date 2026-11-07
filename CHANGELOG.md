# Changelog

## 1.1.0

- Split the former combined "ring" event into two distinct MQTT Discovery events per
  board: a **Call** (`ring`) event, now fired on the incoming SDP offer, and a **Call
  Photo** (`photo`) event, fired on the call-log message with the snapshot.
- Door no longer auto-opens on an incoming SDP offer; it now only triggers the `ring`
  event — opening the door is a separate, explicit action.
- Renamed the generic MQTT topics `.../event` to `.../ring_event` and `.../photo_event`
  to match the new events.

## 1.0.1

- Added `README.md` with setup, configuration, fake-server, and legal-notice sections.
- Limited the in-memory event log to the last 50 entries; older events are now evicted
  as new ones arrive.

## 1.0.0

- Initial Home Assistant Supervisor add-on release.
- MQTT Discovery for door-open button, incoming-call event, and call snapshot camera.
- Ingress support for the pairing/dashboard web UI.
- Configurable `auth_address`/`main_address` via add-on options.
