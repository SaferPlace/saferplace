# Contributing to Safer.Place

## Running Locally

Make sure you have the latest version of Go installed

```sh
go run ./cmd/saferplace
```

This will start the binary on the default port `8080` (configurable using the
`PORT` env var).

## Commit Messages

Please make sure they are prefixed with the package or packages they are
touching. If the change is in the `internal` folder, the `internal` prefix can
be omitted.