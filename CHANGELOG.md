# Changelog

## 1.0.1

- Added `README.md` with setup, configuration, fake-server, and legal-notice sections.
- Limited the in-memory event log to the last 50 entries; older events are now evicted
  as new ones arrive.

## 1.0.0

- Initial Home Assistant Supervisor add-on release.
- MQTT Discovery for door-open button, incoming-call event, and call snapshot camera.
- Ingress support for the pairing/dashboard web UI.
- Configurable `auth_address`/`main_address` via add-on options.
