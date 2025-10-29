# Thread-Safe Configuration Access with `sync.RWMutex`

This project now zentralizes all configuration access through a `sync.RWMutex` to keep the shared `config.AppConfig` safe when multiple gRPC handlers execute concurrently.

## Why a RWMutex?

The gRPC service implementations (`ioctrl`, `network`, `system`, `video`, `watermark`, `lux`, and `fota`) all read and update the global configuration. These handlers are invoked on independent goroutines, which means simultaneous map writes or write/read combinations can trigger data races or the dreaded `concurrent map read and map write` panic. Using a read/write mutex gives:

- **Mutual exclusion for writers:** only one update can mutate the config at a time.
- **Concurrent reads:** multiple handlers can inspect the config in parallel when no writes are happening.
- **Minimal refactoring cost:** service code now calls helper functions exposed by `config/manager.go`, so the synchronization logic stays in one place.

## Key Helpers

```go
// config/manager.go
var (
    configLock sync.RWMutex
    AppConfig  Config
)

func UpdateConfig(update func(*Config)) error {
    configLock.Lock()
    defer configLock.Unlock()
    update(&AppConfig)
    return SaveAppConfigDefault()
}

func ReadConfig(read func(Config)) {
    configLock.RLock()
    defer configLock.RUnlock()
    read(AppConfig)
}
```

- `UpdateConfig` must wrap every mutation. It takes a callback, executes it while holding the write lock, and persists the result.
- `ReadConfig` hands callers an immutable snapshot (passed by value) while holding a read lock. Callers avoid copying large maps manually and remain shielded from concurrent writes.

## Calling Patterns

### Reading

```go
var ipv4 string
cfg.ReadConfig(func(current cfg.Config) {
    ipv4 = current.Network.IPv4
})
```

### Writing

```go
err := cfg.UpdateConfig(func(c *cfg.Config) {
    if c.LEDs == nil {
        c.LEDs = make(map[string]cfg.LEDConfig)
    }
    c.LEDs[channelKey] = cfg.LEDConfig{
        StatusLed: status,
        RecLedOn:  in.RecLedOn,
    }
})
```

Every gRPC handler now uses this pattern, so there are no lingering direct references to `cfg.AppConfig` in the service structs themselves.

## Testing

For consistency with existing CI runs, use the amd64 build target when executing unit tests locally:

```bash
cd grpcd-copilot/grpcd-study/grpcd-copilot
GOARCH=amd64 go test ./...
```

This enforces a deterministic build environment and verifies the RWMutex guarantees do not break existing behavior.

## Next Steps

- If additional packages need safe config access, route them through these helpers instead of reaching for `AppConfig` directly.
- Consider wrapping commonly accessed config subsets in dedicated accessor functions for even clearer intent.
- Keep any future configuration mutations within `UpdateConfig` to preserve the threading guarantees.
