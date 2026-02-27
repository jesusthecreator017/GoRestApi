# Plan: Full-Stack Issue Tracker (Go API Handlers + React Frontend)

## Context

The Go API has a complete storage layer for issues (Create, GetByID, List + sqlc-generated UpdateStatus, Delete) but **no HTTP handlers or routes** — only `GET /v1/health` exists. The `web/` directory is empty. We need to:
1. Wire up CRUD issue handlers in Go with server-side validation
2. Build a React Kanban board frontend with shadcn/ui, Tailwind CSS, and Zod validation
3. Dockerize the frontend and add it to docker-compose

---

## Phase 1: Go API Backend

### 1.1 JSON helpers — New file: `cmd/api/json.go`
- `writeJSON(w, status, envelope)` — writes JSON response with Content-Type header
- `readJSON(r, dst)` — decodes JSON body with `DisallowUnknownFields()`
- `errorJSON(w, status, message)` — error response
- `validationErrorJSON(w, errors map[string]string)` — 422 with field errors
- `envelope` type alias: `map[string]any`

### 1.2 CORS middleware — New file: `cmd/api/middleware/cors.go`
- Follows existing `Middleware` type from `middleware.go`v
- `CORS(allowedOrigin string) Middleware` — closure pattern like `Timeout`
- Sets `Access-Control-Allow-Origin`, `-Methods` (GET, POST, PATCH, DELETE, OPTIONS), `-Headers` (Content-Type)
- Returns 204 early for OPTIONS preflight requests
- Must be **first** in the middleware stack so preflight doesn't hit Timeout

### 1.3 Extend Storage interface — Edit: `internal/store/storage.go`
Add two methods to the `Issues` interface:
```go
UpdateStatus(context.Context, string, StatusType) (*Issue, error)
Delete(context.Context, string) error
```

### 1.4 Implement UpdateStatus & Delete — Edit: `internal/store/issues.go`
- `UpdateStatus` — parses UUID, calls `queries.UpdateIssueStatus(ctx, dbsqlc.UpdateIssueStatusParams{ID, Status})`, returns domain issue
- `Delete` — parses UUID, calls `queries.DeleteIssue(ctx, uid)`
- Both use existing sqlc-generated code in `dbsqlc/issues.sql.go`

### 1.5 Issue handlers — New file: `cmd/api/issues.go`
Five handlers as methods on `*application`:

| Route | Handler | Validation |
|-------|---------|------------|
| `POST /v1/issues` | `createIssueHandler` | title required & ≤255 chars, author_id required & valid UUID |
| `GET /v1/issues` | `listIssuesHandler` | none |
| `GET /v1/issues/{id}` | `getIssueHandler` | id from path (`r.PathValue("id")`) |
| `PATCH /v1/issues/{id}/status` | `updateIssueStatusHandler` | status must be Incomplete/In-Progress/Complete |
| `DELETE /v1/issues/{id}` | `deleteIssueHandler` | id from path |

- Use `pgx.ErrNoRows` check in GetByID/UpdateStatus/Delete to return 404 vs 500

### 1.6 Config & routing — Edit: `cmd/api/api.go` + `cmd/api/main.go`
- Add `corsOrigin string` to `config` struct
- Read `CORS_ORIGIN` env var in `main.go` (default `"http://localhost:5173"`)
- Add `middleware.CORS(app.config.corsOrigin)` as first item in the stack
- Register all 5 issue routes in `mount()`

---

## Phase 2: React Frontend

### 2.1 Scaffold — in `web/`
```bash
npm create vite@latest web -- --template react-ts
cd web && npm install
npx shadcn@latest init    # sets up Tailwind CSS v4 + shadcn
npm install zod @tanstack/react-query
npx shadcn@latest add button card dialog input label select badge textarea alert-dialog
```

### 2.2 File structure
```
web/src/
  schemas/issue.ts          — Zod schemas (Issue, CreateIssueInput, StatusType)
  api/client.ts             — fetch wrapper with VITE_API_URL base
  api/issues.ts             — typed CRUD functions
  hooks/useIssues.ts        — React Query hooks
  components/kanban/
    KanbanBoard.tsx          — 3-column layout, calls useIssues(), groups by status
    KanbanColumn.tsx         — single column with header + issue cards
    IssueCard.tsx            — shadcn Card with title, description, status Badge, delete button
  components/issues/
    CreateIssueDialog.tsx    — shadcn Dialog with form, Zod validation
    DeleteIssueAlert.tsx     — shadcn AlertDialog for confirmation
  App.tsx                   — QueryClientProvider + header + KanbanBoard
```

### 2.3 Zod schemas — `web/src/schemas/issue.ts`

```typescript
import { z } from "zod";

export const StatusType = z.enum(["Incomplete", "In-Progress", "Complete"]);
export type StatusType = z.infer<typeof StatusType>;

export const IssueSchema = z.object({
  id: z.string().uuid(),
  author_id: z.string().uuid(),
  title: z.string().min(1, "Title is required").max(255, "Title too long"),
  description: z.string(),
  status: StatusType,
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});
export type Issue = z.infer<typeof IssueSchema>;

export const CreateIssueSchema = z.object({
  author_id: z.string().uuid("Must be a valid UUID"),
  title: z.string().min(1, "Title is required").max(255, "Title too long"),
  description: z.string().default(""),
});
export type CreateIssueInput = z.infer<typeof CreateIssueSchema>;

export const UpdateStatusSchema = z.object({
  status: StatusType,
});
export type UpdateStatusInput = z.infer<typeof UpdateStatusSchema>;
```

### 2.4 API layer — `web/src/api/`

**`client.ts`**
```typescript
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json", ...init?.headers },
    ...init,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error ?? `HTTP ${res.status}`);
  }
  return res.json();
}
```

**`issues.ts`**
```typescript
import { apiFetch } from "./client";
import type { Issue, CreateIssueInput, StatusType } from "../schemas/issue";

export const issuesApi = {
  list: () =>
    apiFetch<{ issues: Issue[] }>("/v1/issues").then((r) => r.issues),

  getById: (id: string) =>
    apiFetch<{ issue: Issue }>(`/v1/issues/${id}`).then((r) => r.issue),

  create: (data: CreateIssueInput) =>
    apiFetch<{ issue: Issue }>("/v1/issues", {
      method: "POST",
      body: JSON.stringify(data),
    }).then((r) => r.issue),

  updateStatus: (id: string, status: StatusType) =>
    apiFetch<{ issue: Issue }>(`/v1/issues/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    }).then((r) => r.issue),

  delete: (id: string) =>
    apiFetch<{ message: string }>(`/v1/issues/${id}`, {
      method: "DELETE",
    }),
};
```

### 2.5 React Query hooks — `web/src/hooks/useIssues.ts`

```typescript
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { issuesApi } from "../api/issues";
import type { CreateIssueInput, StatusType } from "../schemas/issue";

export function useIssues() {
  return useQuery({ queryKey: ["issues"], queryFn: issuesApi.list });
}

export function useCreateIssue() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateIssueInput) => issuesApi.create(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["issues"] }),
  });
}

export function useUpdateIssueStatus() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: StatusType }) =>
      issuesApi.updateStatus(id, status),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["issues"] }),
  });
}

export function useDeleteIssue() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => issuesApi.delete(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["issues"] }),
  });
}
```

### 2.6 Components

- **KanbanBoard**: Fetches issues via `useIssues()`, filters into 3 columns by status, renders `KanbanColumn` x3 in a flex row. Has "New Issue" button → opens `CreateIssueDialog`.
- **KanbanColumn**: Takes `status`, `label`, `issues[]`. Renders header with count badge, maps issues to `IssueCard`.
- **IssueCard**: shadcn `Card` with title, truncated description, colored `Badge` for status, inline `Select` for status change, trash icon → opens `DeleteIssueAlert`.
- **CreateIssueDialog**: shadcn `Dialog` with form inputs. Validates with `CreateIssueSchema.safeParse()`, shows inline errors, calls `useCreateIssue` on submit.
- **DeleteIssueAlert**: shadcn `AlertDialog` with confirm/cancel. Calls `useDeleteIssue` on confirm.

---

## Phase 3: Docker Integration

### 3.1 Frontend Dockerfile — New file: `web/Dockerfile`

```dockerfile
# Build
FROM node:22-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG VITE_API_URL=http://localhost:8080
RUN npm run build

# Serve
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

### 3.2 Nginx config — New file: `web/nginx.conf`

```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

### 3.3 Update `docker-compose.yml`

Add `frontend` service and `CORS_ORIGIN` to `api`:

```yaml
services:
  db:
    # ... unchanged ...

  api:
    build: .
    environment:
      DATABASE_URL: "postgres://postgres:postgres@db:5432/fswithgo?sslmode=disable"
      ADDR: ":8080"
      CORS_ORIGIN: "http://localhost:3000"
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy

  frontend:
    build:
      context: ./web
      args:
        VITE_API_URL: "http://localhost:8080"
    ports:
      - "3000:80"
    depends_on:
      - api

volumes:
  pgdata:
```

---

## Files Summary

| Action | File | Purpose |
|--------|------|---------|
| **Create** | `cmd/api/json.go` | JSON response/request helpers |
| **Create** | `cmd/api/middleware/cors.go` | CORS middleware |
| **Create** | `cmd/api/issues.go` | 5 CRUD handlers with validation |
| **Edit** | `internal/store/storage.go` | Add UpdateStatus, Delete to interface |
| **Edit** | `internal/store/issues.go` | Implement UpdateStatus, Delete |
| **Edit** | `cmd/api/api.go` | Add routes + CORS to stack + config |
| **Edit** | `cmd/api/main.go` | Add corsOrigin to config |
| **Create** | `web/` (scaffold) | Vite + React + TS project |
| **Create** | `web/src/schemas/issue.ts` | Zod schemas |
| **Create** | `web/src/api/client.ts` | Fetch wrapper |
| **Create** | `web/src/api/issues.ts` | Typed API functions |
| **Create** | `web/src/hooks/useIssues.ts` | React Query hooks |
| **Create** | `web/src/components/kanban/*.tsx` | Board, Column, Card |
| **Create** | `web/src/components/issues/*.tsx` | CreateDialog, DeleteAlert |
| **Create** | `web/src/App.tsx` | Root component |
| **Create** | `web/Dockerfile` | Multi-stage frontend build |
| **Create** | `web/nginx.conf` | SPA nginx config |
| **Edit** | `docker-compose.yml` | Add frontend service + CORS env |

---

## Verification

```bash
# 1. Rebuild and start everything
docker compose up --build

# 2. Test Go API with curl
curl http://localhost:8080/v1/issues
curl http://localhost:8080/v1/issues -X POST \
  -H "Content-Type: application/json" \
  -d '{"author_id":"550e8400-e29b-41d4-a716-446655440000","title":"Test issue"}'

# 3. Open frontend
open http://localhost:3000

# 4. Test CRUD in the UI
# - Create issue via dialog
# - Change status via dropdown
# - Delete issue via trash icon + confirm
```
