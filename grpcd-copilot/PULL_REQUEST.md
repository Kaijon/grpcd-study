# Refactor config usage into cfg package and update MQTT handlers

## Summary
- remove the duplicate `config.go` and add `config_bridge.go` so legacy globals delegate to `cfg`
- update MQTT handlers and command helper to read/write config through `cfg`
- align the copied UnitTest tree to import `grpcd/config` instead of removed globals

## Testing
- `GOARCH=amd64 go test ./...`
