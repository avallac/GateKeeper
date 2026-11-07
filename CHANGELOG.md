# Changelog

## 1.1.3

- Fixed the same kind of decode error on `brd_sp_call_log`'s `timestamp` (`cannot
  unmarshal string into ... int64`): the IPGuard server sometimes sends it as an
  RFC3339 date string instead of a Unix-millisecond number. The field is now parsed
  via a new tolerant `FlexInt64` type.

## 1.1.2

- Fixed the same kind of decode error on `responder_type` (`cannot unmarshal string
  into ... int`): the field is now parsed via a new tolerant `FlexInt` type.

## 1.1.1

- Fixed a decode error on `brd_sp_call_log` (`cannot unmarshal string into ... bool`):
  the IPGuard server inconsistently sends `open_status` as either a JSON boolean or a
  stringified boolean. This field is now parsed via a new tolerant `FlexBool` type, so
  the call-log message (and its photo/snapshot) is no longer dropped.

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
