# fileinspect

A command-line tool for inspecting files and detecting tampering using
SHA-256 hashes. Point it at a file (or a whole directory) and it records
a fingerprint you can verify against later — if the file changes, the
hash won't match.

## Usage

```bash
# Inspect a single file: prints metadata, computes SHA-256, saves record.json
fileinspect <file>

# Hash every file in a directory recursively, save batch_record.json
fileinspect --batch <directory>

# Verify a file against its previously saved record.json
fileinspect --verify <file>
```

### Example

```bash
$ fileinspect main.go
Name: main.go
Size: 804
Modified: 2026-07-26 09:12:03
SHA-256: dbe9e57b00d20adfb68c145ce434c31da5ff7e60bfe4b1dd2e30c1da90db9ffc
Record saved to record.json

$ fileinspect --verify main.go
Recorded hash: dbe9e57b00d20adfb68c145ce434c31da5ff7e60bfe4b1dd2e30c1da90db9ffc
Current hash:  dbe9e57b00d20adfb68c145ce434c31da5ff7e60bfe4b1dd2e30c1da90db9ffc
VERIFIED: File has not been tampered with.
```

If the file changes after being recorded, `--verify` reports the mismatch
and **exits with a non-zero status code**, so it can be used in scripts,
pre-commit hooks, or CI pipelines to fail fast on unexpected file changes:

```bash
$ fileinspect --verify main.go
Recorded hash: dbe9e57b00d20adfb68c145ce434c31da5ff7e60bfe4b1dd2e30c1da90db9ffc
Current hash:  9f2c1a...
ALERT: File has been modified since recording.
$ echo $?
1
```

See [`examples/`](examples/) for sample output from single-file and batch
runs.

## Known limitations

- `record.json` and `batch_record.json` are overwritten on every run —
  there's no per-file naming yet, so recording a second file discards the
  first file's record. Fine for single-file audit workflows; a limitation
  to be aware of if inspecting many files one at a time.
- `--verify` currently reads only from `record.json` in the working
  directory, so it verifies against the most recently recorded single
  file, not against `batch_record.json` entries.

## Build & test

```bash
go build .
go vet ./...
go test ./... -v
```

## Structure

- `main.go` — CLI entry point, flag parsing, command routing
- `file_operations.go` — hashing, single-file and batch processing
- `file_record.go` — `FileRecord` struct
- `verify.go` — verification against a saved record
- `file_operations_test.go` — tests covering hashing, batch mode, and
  tamper detection
- `testevidence/` — sample fixture files used in manual/local testing
- `examples/` — sample recorded output

## Part of

`fileinspect` is a standalone tool and a component of the Aegis Anchor
evidence preservation system.

## License

[MIT](LICENSE)
