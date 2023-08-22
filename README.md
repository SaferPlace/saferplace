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

Add the following to your env, `direnv` recommended, but do not commit the `.envrc`.

```sh
export MINIO_ENDPOINT=localhost:9000
export MINIO_ACCESS_KEY=saferplace
export MINIO_SECRET_KEY=supersecret
```

### Running

Open up 3 tabs in your terminal

#### Docker Compose

```sh
# ~/workdir/realtime
$ docker compose up
```

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
