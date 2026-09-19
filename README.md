<h1 align="center">🍔 bite</h1>

<p align="center">
  <em>Automate your food and groceries. Order without touching your phone.</em>
</p>

<p align="center">
  <a href="#"><img src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white" alt="Go 1.25"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License"></a>
  <a href="#status"><img src="https://img.shields.io/badge/status-alpha-orange.svg" alt="Alpha"></a>
  <img src="https://img.shields.io/badge/creds-never%20in%20repo-brightgreen.svg" alt="Credentials never in repo">
</p>

---

**bite** is a Go toolkit that drives delivery apps through their *private* mobile
APIs — the same endpoints the phone app calls, no browser automation, no screen
scraping. Point it at a store, build a cart, and (soon) place the order, all from
your terminal or your own agent.

It exists because the good delivery apps in the region ship no public API and no
usable web checkout — but their mobile backends are just JSON. bite speaks that
JSON directly.

> **Status:** alpha. Talabat (Egypt) is implemented for account, food and grocery
> reads plus cart writes. Order placement is intentionally not automated yet — see
> [Roadmap](#roadmap).

## Why

- 📱 **No phone.** Search, compare and cart from the CLI (or a script, or an agent).
- ⚡ **Fast.** Direct API calls return in milliseconds — no emulator in the hot path.
- 🧩 **Pluggable.** Providers implement a small shared surface; Talabat first,
  Breadfast & Rabbit next.
- 🔒 **Safe by design.** Your tokens never enter the repo. Paying is never automatic.

## Providers

| Provider | Account | Food (restaurants) | Grocery (q-commerce) | Place order |
|----------|:-------:|:------------------:|:--------------------:|:-----------:|
| **Talabat** (🇪🇬) | ✅ profile, addresses | ✅ list, menu | ✅ catalog, cart | 🚧 planned |
| Breadfast | — | — | — | — |
| Rabbit | — | — | — | — |

## Install

```bash
go install github.com/0xSherlokMo/bite/cmd/bite@latest
# or from source:
git clone https://github.com/0xSherlokMo/bite && cd bite
go build -o bin/bite ./cmd/bite
```

## Quick start

bite authenticates by replaying the exact headers the mobile app sends. You supply
those once via a session file (see [Authentication](#authentication)), then:

```bash
# import your captured session (one-time / when tokens expire)
bite talabat session import ./session.json

# account
bite talabat profile
bite talabat addresses

# food
bite talabat food restaurants --lat 30.04 --lon 30.98 --area 8058
bite talabat food menu 631087

# grocery / talabat mart  (no payment involved)
bite talabat mart cart --vendor <uuid> --branch <id> --chain <id> --lat 30.04 --lon 30.98
bite talabat mart add  --vendor <uuid> --branch <id> --chain <id> \
                       --product <uuid> --qty 2 --lat 30.04 --lon 30.98
```

Every command prints JSON (or a table for lists), so it pipes cleanly into `jq`,
scripts, or an LLM agent.

## How it works

```
              ┌────────────────────────────────────────────┐
   your CLI   │  cmd/bite         command routing / output  │
   or agent ──▶  internal/providers/talabat                 │
              │    account · food · grocery   (typed calls) │
              │  internal/secrets  session store (0600)     │
              └───────────────┬────────────────────────────┘
                              │ replays app headers (JWT + device tokens)
                              ▼
                        api.talabat.com  (private mobile API)
```

The mobile backends gate requests on an `authorization` JWT **plus** device and
anti-fraud headers (`x-api-key`, `x-device-id`, `incognia-token`, `x-segmenttoken`,
…). bite stores that header set and attaches it to each request. Different
microservices require different subsets — the grocery service, for example, needs
`x-api-key` and `x-customer-id` that the account service doesn't — so bite persists
the **union** of headers seen across services.

## Authentication

Tokens are minted by the real app, so you capture them once:

1. Run the app on a device/emulator behind an HTTPS proxy you control
   (e.g. mitmproxy with its CA trusted).
2. Trigger a few authenticated requests (open the app, pull to refresh).
3. Export the union of request headers to a JSON object: `{ "header": "value", … }`.
4. `bite talabat session import ./session.json`.

At minimum the file needs `authorization` and `x-device-id`; grocery endpoints also
need `x-api-key` and `x-customer-id`. See [`docs/refresh.md`](./docs/refresh.md) for
the full flow, and [`internal/secrets/talabat.example.json`](./internal/secrets/talabat.example.json)
for the shape.

The `authorization` JWT expires after a few hours; re-run `session import` to refresh.

## Security

- **Credentials never live in the repository.** Sessions are written to
  `~/.config/bite/<provider>.json` with `0600` permissions and are loaded at
  runtime. The repo ships only `*.example.json` templates with placeholders.
- `.gitignore` blocks config, `*.session.json`, `*.flows`, and `.env`.
- **Paying is never automatic.** Reads and cart edits are automated; placing a paid
  order will always be a deliberate, separately confirmed action.

Found a security issue? Please open a private advisory rather than a public issue.

## Roadmap

- [ ] **Order placement** — behind an explicit confirmation gate
- [ ] **Search** — "find shawarma near me"
- [ ] **Automated token refresh** — headless re-mint when the JWT expires
- [ ] **Breadfast** and **Rabbit** providers
- [ ] Typed menu / catalog models (currently raw JSON)

## Legal

bite is an independent, unofficial project for **personal use with your own
account**. It is not affiliated with, endorsed by, or connected to Talabat,
Delivery Hero, or any other company. Automating an account may conflict with a
provider's Terms of Service — you are responsible for how you use it. No
warranty; see [LICENSE](./LICENSE).

## License

MIT © [0xSherlokMo](https://github.com/0xSherlokMo)
