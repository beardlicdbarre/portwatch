# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/yourname/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and log file:

```bash
portwatch start --interval 60 --log /var/log/portwatch.log
```

Run a one-time snapshot of currently open ports:

```bash
portwatch scan
```

When an unexpected port opens or closes, `portwatch` will print an alert to stdout (and optionally to a log file):

```
[ALERT] 2024-06-10 14:32:01 - New port opened: TCP 8080
[ALERT] 2024-06-10 14:35:44 - Port closed: TCP 3000
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `30` | Scan interval in seconds |
| `--log` | `""` | Path to log file (optional) |
| `--baseline` | `""` | Path to a known-good ports snapshot |

---

## License

MIT © 2024 yourname