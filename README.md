# SaferPlace

Trying to make the world a little bit safer

> **warning**
> SaferPlace is not associated with any official institutions.

## What is SaferPlace

SaferPlace is a platform created to report incidents which happen in the area. So often there are
crimes and violations in your surrounding which are not reported. Originally designed to tell you
which areas are safe to live in for college based on the Gardai crime report data, we are now
expanding the idea by letting you report safety infractions.

The project is currently in **heavy** development, but if you would like to see what we are
currently working on, take a look at our [Project Milestones] and [Project Backlog]

The current plan is to have the site fully functional by the end of 2023 with the website being
available at [https://safer.place]

## Old Code

- If you would like to see the original implementation of saferplace, take a look at
  [Original Saferplace].
- There was an attempt at rewriting the original codebase, but that attempt was then replaced by
  "realtime saferplace" (community sourced, this repository). If you would like to see it, the
  efforts are kept on the [saferplace-v1 Branch]

---

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
$ pnpm run dev --host
```

#### Reviewer App

```sh
# ~/workdir/realtime/packages/review-ui
$ pnpm run dev --host
```

Then proceed to http://localhost:5173/login, and change the backend to `http://localhost:8001`.

[Project Backlog]: https://github.com/orgs/SaferPlace/projects/2/
[Project Milestones]: https://github.com/SaferPlace/saferplace/milestones
[https://safer.place]: https://safer.place
[Original Saferplace]: https://github.com/saferplace/original
[saferplace-v1 Branch]: https://github.com/SaferPlace/saferplace/tree/saferplace-v1
