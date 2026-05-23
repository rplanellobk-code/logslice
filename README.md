# logslice

A fast log file splitter and time-range extractor for large structured log archives.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

---

## Usage

Extract log entries within a specific time range:

```bash
logslice --input /var/log/app.log --from "2024-01-15T08:00:00" --to "2024-01-15T12:00:00" --output slice.log
```

Split a large log file into chunks by time interval:

```bash
logslice split --input /var/log/app.log --interval 1h --output-dir ./chunks/
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Path to the source log file | required |
| `--from` | Start timestamp (RFC3339) | — |
| `--to` | End timestamp (RFC3339) | — |
| `--interval` | Chunk duration for split mode | `1h` |
| `--output` | Output file path | stdout |
| `--output-dir` | Output directory for split mode | `./` |
| `--format` | Log timestamp format | `rfc3339` |

---

## Features

- Handles multi-gigabyte log files efficiently with streaming I/O
- Supports common structured log formats (JSON, logfmt, Apache, nginx)
- Binary search indexing for fast time-range seeks
- Gzip input/output support

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)