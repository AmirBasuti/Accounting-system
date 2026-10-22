# Accounting Repository Layer in Go

A Go and PostgreSQL project exploring persistence for ledgers, vouchers, and voucher items. The implementation uses GORM repositories, transactions around voucher changes, and validation of selected record relationships.

The current entry point connects to the database and attempts schema migration. The repository contains persistence code and integration tests; it does not currently expose an HTTP API or a complete user-facing accounting application.

## Implemented components

| Component | Responsibility |
|---|---|
| `SL` | Sub-ledger code, title, detail flag, and version. |
| `DL` | Detail-ledger code, title, and version. |
| `Voucher` | Voucher number, version, and associated items. |
| `VoucherItem` | Ledger references and integer debit/credit fields. |

Repositories provide create, read, update, and delete operations. Voucher creation, updates, and deletion use GORM transactions; voucher reads preload their items. Validation checks ledger existence, whether an item uses debit or credit, and whether a voucher contains between 2 and 500 items. Several operations compare record versions before changing data.

These rules describe the implemented behavior. They do not establish complete accounting correctness or concurrency-safe locking.

## Requirements and build

- Go 1.22.5 or later, as declared in `go.mod`.
- PostgreSQL for database execution and integration tests.

```sh
git clone https://github.com/AmirBasuti/Accounting-system.git
cd Accounting-system
go mod download
go build -mod=readonly ./...
```

The source compiled with Go 1.22.5 on Windows/amd64 during the 22 October 2026 review. Database execution and integration tests were not run during that review.

## Database setup status

Portable setup needs work before this repository can offer a verified quickstart:

- `internal/database/config.go` loads a `.env` file from a fixed development-machine path. A `.env` in the repository root is not used automatically by this revision.
- The loader expects `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, and `DB_PORT`.
- Both GORM AutoMigrate and SQL migration files exist. Their table-naming conventions need reconciliation before selecting one authoritative setup procedure.
- `.env` is not currently covered by the repository's ignore rules. Keep local credentials outside tracked files.

The `cmd/main.go` entry point attempts database connection and migration; it is not a server. Database configuration and schema initialization should be made portable before publishing execution instructions as a working installation procedure.

## Tests

The `tests/` directory contains PostgreSQL integration test source covering ledger CRUD, version mismatches, voucher operations, and ledger references. The suite connects through the same configuration loader and creates or deletes database records.

It is not a self-contained unit-test suite: its setup depends on an existing database and schema, and the test bootstrap currently ignores the connection error. Use a disposable database when validating it after the setup issues are resolved. No passing test result is claimed here.

## Project layout

```text
cmd/main.go                Database connection and migration entry point
internal/database/         Configuration, connection, and GORM migration
internal/models/           Ledger, voucher, and voucher-item data models
internal/repository/       Persistence operations and validation
migrations/                SQL migration scripts
tests/                     PostgreSQL integration tests
```

## Current limitations

- Whole-voucher debit and credit totals are not checked for equality. The current validation is per item.
- Version comparisons use separate reads and writes; concurrent-update safety has not been demonstrated.
- Detail-ledger deletion checks require review against the actual column names.
- Integer amount fields have no documented currency or unit convention.
- There is no claim of production use, financial compliance, or suitability for handling real financial records.

## License

No license is specified in this repository at the reviewed revision.
