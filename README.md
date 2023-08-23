# SaferPlace Realtime Monorepo

This monorepo is designed to temporarily house saferplace realtime application
for the ease of development

## Local Development

> **NOTE**
> This documentation is incomplete and needs to be expanded, but should be enough to get you
> started.

### Workflow

1. Find a task from the [Github Project Backlog]
2. Make a PR and reference the issue its closing
3. Ensure the PRs are passing

### Dependencies

- Go
- node and pnpm
- Docker

For now these instructions are not well optimized

To load the recommended Env Vars automatically use `direnv`. If there are any secret variables
eg. OAuth credentials, use the `.secret.envrc`.

### Running

Open up 3 tabs in your terminal

#### Docker Compose

```sh
# ~/workdir/realtime
$ docker compose up
```

Enable anonymous readonly access on minio

- Navigate to http://localhost:9001/buckets/images/admin/prefix using the credentials from `.envrc`.
- There click on `Add Access Rule` and add an anonymous access rule
  - `Prefix` - `/`
  - `Access` - `readonly`

This will ensure that you need authentication to write to the bucket, but not to read.

#### Go Backend

Go backend is currently running as a single app. It will run on `http://localhost:5173`

```sh
# ~/workdir/realtime
$ go run ./cmd/realtime
```

#### PWA Frontend

```sh
# ~/workdir/realtime/packages/pwa
$ pnpm run dev
```

Then proceed to http://localhost:5173/login, and change the backend to `http://localhost:8001`.

[Github Project Backlog]: https://github.com/orgs/SaferPlace/projects/2/
