# Refreshing the Talabat session

The `authorization` JWT expires after a few hours; other headers (`x-api-key`,
`incognia-token`, `x-device-id`, `x-segmenttoken`) are longer-lived but also come
from the app. When calls start returning 401, re-mint a session:

1. Boot the rooted Android emulator snapshot that is already logged in:
   ```
   emulator -avd talabat -writable-system -no-snapshot-load \
     -http-proxy http://127.0.0.1:8080 & 
   # (or load the `logged_in` snapshot)
   ```
2. Run mitmproxy on :8080 with the CA trusted in the emulator's system store.
3. Open Talabat / pull to refresh so the app issues authenticated requests.
4. Export the **union** of headers seen across `api.talabat.com` +
   `userlocation.talabat.com` requests to a JSON object `{header: value}`.
5. Import it:
   ```
   bite talabat session import session_full.json
   ```

The header set must include at minimum `authorization` and `x-device-id`, and for
grocery endpoints also `x-api-key` and `x-customer-id`.

> Automating steps 1–4 (headless token minting) is the next milestone.
