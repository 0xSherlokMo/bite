# 🍔 bite

Order food and groceries from the terminal instead of doomscrolling the app.

Delivery apps here have no public API and no real web checkout — but their phone
apps just talk JSON. `bite` talks that JSON directly. Started as a "I'm too lazy to
open Talabat" side project; ended up being kind of useful, so here it is.

Written in Go. Talabat works today. Breadfast and Rabbit are on the list.

## What it does

```bash
bite talabat profile
bite talabat addresses
bite talabat food restaurants --lat 30.04 --lon 30.98 --area 8058
bite talabat food menu 631087
bite talabat mart add --vendor <uuid> --branch <id> --chain <id> --product <uuid> --qty 2 --lat 30.04 --lon 30.98
```

Reads and cart edits work. Actually placing an order (i.e. spending money) isn't
automated yet — on purpose. That'll be a separate, confirm-first command.

## Setup

```bash
go build -o bin/bite ./cmd/bite
bite talabat session import ./session.json
```

`bite` just replays the headers the real app sends. You grab them once (device +
mitmproxy), drop them in a JSON file, and import. Details in
[`docs/refresh.md`](./docs/refresh.md); shape in
[`talabat.example.json`](./internal/secrets/talabat.example.json). The `authorization`
token expires every few hours — re-import to refresh.

## Creds

Your tokens never touch the repo. They live in `~/.config/bite/` (chmod 600),
loaded at runtime. The repo only ships example files with `<placeholders>`.

## Roadmap

- [ ] Place orders (with a confirm gate)
- [ ] Search
- [ ] Auto token refresh
- [ ] Breadfast + Rabbit

## Heads up

Personal project, personal account, unofficial. Not affiliated with Talabat /
Delivery Hero. Automating your account might break their ToS — that's on you. MIT.
