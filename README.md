# gofer

Go toolkit that automates delivery apps (Talabat first; Breadfast & Rabbit next)
over their private mobile APIs. Named for the Go gopher — and for what it does:
runs your errands.

## Design principles

- **Credentials never live in the repo.** Session tokens are stored in
  `~/.config/gofer/<provider>.json` with `0600` perms, loaded at runtime.
  The repo ships only `*.example` templates.
- **Providers are pluggable.** Each app implements a small common interface;
  vertical-specific calls (food vs. grocery) live in the provider package.
- **Paying is gated.** Read/search/cart are automated. Placing a paid order is a
  deliberate, separate, confirmed step — never automatic.

## Layout

```
cmd/gofer            CLI entrypoint
internal/secrets     secure session store (~/.config/gofer, 0600)
internal/providers   Provider interface + per-app clients
  talabat            account, food (restaurants), grocery (talabat mart)
```

## Quick start

```bash
go build -o bin/gofer ./cmd/gofer

# Import a captured session into the secure store (one-time / on refresh)
bin/gofer talabat session import /path/to/session_headers.json

# Reads
bin/gofer talabat profile
bin/gofer talabat addresses
bin/gofer talabat food restaurants --lat 30.0414 --lon 30.9809
bin/gofer talabat food menu <branchId>

# Cart (no payment)
bin/gofer talabat mart cart --vendor <uuid> --branch <id> --chain <id>
bin/gofer talabat mart add  --vendor <uuid> --branch <id> --chain <id> --product <uuid> --qty 1
```

## Token freshness

The `authorization` JWT expires after a few hours. `session import` refreshes it.
Tokens are minted by the real app; a companion capture step (rooted Android
emulator + mitmproxy) re-extracts them when needed. See `docs/refresh.md`.
