# ovh-ddns-update

A small command-line tool that keeps an [OVH DynHost](https://help.ovhcloud.com/csm/en-dns-dynhost) record in sync with your current public IPv4 address. Run it periodically (e.g. from cron or a systemd timer) on a machine behind a dynamic IP.

## How it works

1. Detects the current public IPv4 address using [go-external-ip](https://github.com/glendc/go-external-ip), which queries several public services and takes a consensus.
2. Reads the last IP it pushed from a small state file (the "database", see `--database`).
3. If the two differ, it calls the OVH DynDNS endpoint (`https://www.ovh.com/nic/update`) with your DynHost credentials sent as HTTP basic auth.
4. If OVH answers `good` or `nochg`, the new IP is saved to the state file. Any other answer (bad credentials, unknown zone, HTTP error, ...) is logged and the program exits with a non-zero status **without** saving the IP, so the update is retried on the next run.

If the IP hasn't changed, no request is sent to OVH.

## Usage

```
Usage: ovh-ddns-update [--database DATABASE] DOMAIN OVHID OVHPASSWORD

Positional arguments:
  DOMAIN                 OVH DynHost domain
  OVHID                  OVH DynHost identifier
  OVHPASSWORD            OVH DynHost identifier password

Options:
  --database DATABASE    database file [default: .ovh-ddns-update.db]
  --help, -h             display this help and exit
```

`OVHID` and `OVHPASSWORD` are the DynHost identifier and its password, which you create in the OVH control panel (Domains → your zone → DynHost → Manage access). They are not your OVH account credentials.

Example crontab entry, running every 5 minutes:

```
*/5 * * * * /usr/local/bin/ovh-ddns-update --database /var/lib/ovh-ddns-update.db home.example.com myid mypassword
```

The default database path is relative to the working directory, so pass an absolute `--database` path when running from cron or a service.

## Building

Requires [Go](https://go.dev/dl/) 1.26 or newer (see `go.mod`).

```sh
git clone https://github.com/fabiocody/ovh-ddns-update.git
cd ovh-ddns-update
go build -o ovh-ddns-update .
```

For a smaller, statically linked binary, optionally cross-compiled:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ovh-ddns-update .
```

Common targets: `linux/amd64`, `linux/386`, `linux/arm` (set `GOARM=5`, `6` or `7`), `linux/arm64`, `darwin/amd64`, `darwin/arm64`.

### Releases

Release binaries are built with [GoReleaser](https://goreleaser.com/) (configuration in `.goreleaser.yml`). Pushing a tag matching `v*` triggers the `goreleaser_release` GitHub Actions workflow, which publishes a binary per OS/architecture to the GitHub release. Pull requests run a snapshot build (`goreleaser build --snapshot`) to check that everything still compiles.

To try a release build locally:

```sh
goreleaser build --snapshot --clean
```

## Code overview

The whole program is in `main.go`:

| Function | Purpose |
| --- | --- |
| `main` | Parses arguments, compares current and previous IP, triggers the update. |
| `GetCurrentIP` | Detects the public IPv4 address through a consensus of external services. |
| `GetPreviousIP` | Reads the last pushed IP from the database file (`not_found` if it doesn't exist yet). |
| `UpdateDDNS` | Sends the update request to OVH and validates the response. |
| `SaveCurrentIP` | Writes the IP to the database file. |
| `HandleError` | Logs the error with caller info and exits. |

Dependencies: [go-arg](https://github.com/alexflint/go-arg) for argument parsing, [logrus](https://github.com/sirupsen/logrus) for logging, [go-external-ip](https://github.com/glendc/go-external-ip) for IP detection.

## License

See [LICENSE](LICENSE).
