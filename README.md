# 🍔 bite

**Let an AI agent order your food and groceries.**

`bite` gives your agent hands. Delivery apps here have no public API and no real
web checkout — so an LLM can't order for you. `bite` speaks their private mobile
APIs directly and exposes clean, JSON-in / JSON-out commands any agent can call:
search, browse, build a cart, (soon) check out. Point Claude, your own agent, or a
cron job at it and stop opening the app.

Started as "I'm too lazy to open Talabat." Turned into a tool layer my agent uses to
restock the fridge. Written in Go.

## Why it's built for agents

- **Every command returns JSON** → pipe straight into an LLM or a script, no scraping.
- **Stateless & composable** → `search` → `cart add` → `checkout`, one call each.
- **Reads are safe, paying is gated** → an agent can browse and build carts freely;
  spending money is always a separate, explicit, confirmed step.
- **Multi-provider** → one interface, many apps.

## Providers

| Provider | Account | Food | Grocery | Notes |
|----------|:-------:|:----:|:-------:|-------|
| **Talabat** | ✅ | ✅ restaurants + menu | ✅ talabat mart | open API |
| **Rabbit**  | ✅ | ✅ restaurants | ✅ Supermarket+ | open API, no anti-fraud |
| **Breadfast** | — | — | — | blocked: native anti-tamper SDK that self-destructs on rooted/emulated devices |

## What an agent can do today

```bash
# "find me milk and add the cheapest to my Rabbit cart"
bite rabbit store                        # which store serves the address
bite rabbit categories --store-id 27 --store-name EGY010SOD
bite rabbit cart add --product <id> --store-id 27 --lat 30.04 --lon 30.98

# "what's good on Talabat near me?"
bite talabat food restaurants --lat 30.04 --lon 30.98 --area 8058
bite talabat food menu <branchId>
bite talabat mart add --vendor <uuid> --branch <id> --chain <id> --product <uuid>

bite talabat profile        # account info, saved addresses, order history
bite rabbit profile
```

Each returns JSON the agent reads to decide the next call. Placing the paid order
is the one thing that stays behind a human "yes."

## Setup

```bash
go build -o bin/bite ./cmd/bite
bite talabat session import ./session.json   # one-time, per provider
bite rabbit  session import ./session.json
```

`bite` replays the headers the real app sends. You capture them once (device +
mitmproxy), drop them in a JSON file, import. Details in
[`docs/refresh.md`](./docs/refresh.md); shape in the `*.example.json` files. The
`authorization` token expires every few hours — re-import to refresh (this is the
next thing to automate so agents never stall).

## Credentials

Tokens never touch the repo. They live in `~/.config/bite/` (chmod 600), loaded at
runtime. The repo ships only example files with `<placeholders>`.

## Roadmap

- [x] Talabat + Rabbit read/search/cart
- [ ] **Place orders** — behind an agent-friendly confirm gate
- [ ] Auto token refresh (so agents run unattended)
- [ ] A single agent-facing interface (MCP server / tool schema)
- [ ] Breadfast (needs anti-tamper bypass on a real device)

## Heads up

Personal project for automating **your own** account. Unofficial, not affiliated
with Talabat / Delivery Hero / Rabbit. Automating an account may break their ToS —
that's on you. MIT.
