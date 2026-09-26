# mycodes

Source of [mycod.es](https://mycod.es): paste code, get a shareable link with syntax
highlighting and line anchors (`#L5`, `#L5-9`). Pastes expire after 7 days.

- `frontend/`: SvelteKit app (editor, viewer, `/:id/raw` plain-text route)
- `server/`: Go (Fiber + GORM + SQLite) API under `/api/v1/code`

## Getting Started

1. Create an empty `live.db` in `server`. Docker Compose mounts this file; if it is missing, Docker creates a directory instead and SQLite fails to open it.

```shell
# On Windows:
fsutil file createNew server/live.db 0
# On Mac, or Linux:
touch server/live.db
```

2. Copy `server/.env.sample` to `server/.env` and fill in the values:

```shell
cp server/.env.sample server/.env
```

3. Launch the server and the frontend:

```shell
cd server && go run app.go
cd frontend && pnpm install && API_URL=http://localhost:4000 pnpm dev
```

`API_URL` tells the dev server where the API is. Without it, the Docker service
names (`server`, `mycodes-server`) are used.

4. Or, use Docker Compose to launch both:

```shell
make dev
```

## Tests

```shell
cd server && go test ./...
cd frontend && pnpm check && pnpm build
```
