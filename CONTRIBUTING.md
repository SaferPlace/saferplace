# Contributing to SaferPlace

## Local Development

> **NOTE**
> This documentation is incomplete and needs to be expanded, but should be enough to get you
> started.

### Workflow

1. Find a task from the [Project Backlog]
2. Make a PR and reference the issue its closing
3. Ensure the PRs are passing

### Dependencies

- Go
- node and pnpm
- Docker

For now these instructions are not well optimized

### Running

Open up 3 tabs in your terminal

#### Docker Compose

```sh
# ~/workdir/saferplace
$ docker compose up
```

Enable anonymous readonly access on minio

- Navigate to http://localhost:9001/buckets/images/admin/prefix using the credentials from `.envrc`.
- There click on `Add Access Rule` and add an anonymous access rule
  - `Prefix` - `/`
  - `Access` - `readonly`

This will ensure that you need authentication to write to the bucket, but not to read.

#### Go Backend

The Go Backend can be ran as individual components, or as a single binary. For development its easiest to run it as a monolithic binary. By default, it will be running on `http://localhost:8001`

```sh
# ~/workdir/saferplace
$ go run ./cmd/saferplace
```

##### Backend Configuration

You can run the application as a series of individual components by setting the `<component>` or `<component1>,<component1>` arg.

```sh
# ~/workdir/saferplace
$ PORT=8001 go run ./cmd/saferplace report
# In another terminal
$ PORT=8002 go run ./cmd/saferplace viewer
```

Since each binary spins up a webserver, you need to use separate ports.

The binary can be either configured using environment variables, or using a configuration file.
The configuration file can be specified using `-config` flag.

```sh
# ~/workdir/saferplace
$ PORT=8001 go run ./cmd/saferplace -config=config.example.yaml
```

While you can use the config file for local testing, its not necessary as we use `direnv`. Sane
defaults are set for developing locally, and they can be enabled using `direnv allow`. If there
are any secret values, or values that you don't want submitted, use the `.secret.envrc` file which
is loaded automatically if provided.

If you want to view all configuration options you can take a look at the `config.go` file which
contains all the options and their default values.

Tracing is disabled by default but can be enabled using `SAFERPLACE_TRACING_ENABLED=true`, and
setting the endpoint to the `otel-collector` running in Docker Compose with
`SAFERPLACE_TRACING_ENDPOINT=localhost:4317`.

#### PWA Frontend

```sh
# ~/workdir/saferplace/packages/pwa
$ pnpm run dev --host
```

#### Reviewer App

```sh
# ~/workdir/saferplace/packages/review-ui
$ pnpm run dev --host
```

Then proceed to http://localhost:5173/login, and change the backend to `http://localhost:8001`.

[Project Backlog]: https://github.com/orgs/SaferPlace/projects/2/
[Project Milestones]: https://github.com/SaferPlace/saferplace/milestones
