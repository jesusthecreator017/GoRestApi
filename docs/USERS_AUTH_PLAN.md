# Add Users & Basic Authentication

## Context

The app currently has issues and authors but no concept of user accounts or authentication. Anyone can create, update, or delete issues without identifying themselves. This plan adds a `users` table with email/password auth using bcrypt + JWT tokens, following the exact patterns already established in the codebase (repository interfaces, sqlc, anti-corruption layer, helpers, middleware).

This is a **learning guide** — no code is provided, only the steps, file paths, and concepts you need to research and implement yourself.

---

## Phase 1: Users Table & Store Layer

### 1.1 — Create the migrationv

**Files to create:**
- `cmd/migrate/migrations/000003_create_users.up.sql`
- `cmd/migrate/migrations/000003_create_users.down.sql`

**What the `users` table needs:**
- `id` — BIGSERIAL primary key (same pattern as `issues`)
- `email` — TEXT, UNIQUE, NOT NULL
- `password_hash` — TEXT, NOT NULL (you will never store plain passwords)
- `name` — TEXT, NOT NULL
- `created_at` — TIMESTAMPTZ, DEFAULT NOW()
- `updated_at` — TIMESTAMPTZ, DEFAULT NOW()

**Research:** Look at `cmd/migrate/migrations/000001_create_issues.up.sql` and `000002_create_authors.up.sql` to match the exact style.

**Run:** `DATABASE_URL="pgx5://..." go run ./cmd/migrate -direction up`

### 1.2 — Write sqlc queries

**File to create:** `internal/store/queries/users.sql`

**Queries you need (use sqlc annotations like `:one`, `:exec`):**
- `CreateUser` — INSERT returning all fields
- `GetUserByEmail` — SELECT by email (used during login)
- `GetUserByID` — SELECT by id (used by auth middleware)

**Research:** Look at `internal/store/queries/issues.sql` for the annotation format (`-- name: FunctionName :one`).

**Run:** `sqlc generate` — this creates Go code in `internal/store/dbsqlc/`

### 1.3 — Create the domain struct & store wrapper

**File to create:** `internal/store/users.go`

**What to build:**
- A `User` domain struct with JSON tags (like `store.Issue` in `internal/store/issues.go`)
- A `UserStore` struct that wraps `*dbsqlc.Queries` (same pattern as `IssueStore`)
- Methods: `Create(ctx, params) (User, error)`, `GetByEmail(ctx, email) (User, error)`, `GetByID(ctx, id) (User, error)`
- A converter function (like `toDomainIssue`) to translate sqlc types → domain types
- Return `store.ErrNotFound` when `pgx.ErrNoRows` occurs

### 1.4 — Wire into Storage

**File to modify:** `internal/store/storage.go`

- Add a `Users` interface field to the `Storage` struct (same pattern as `Issues` and `Authors`)
- In `NewStorage()`, create a `UserStore` and assign it

---

## Phase 2: Password Hashing

### 2.1 — Add bcrypt dependency

**Run:** `go get golang.org/x/crypto/bcrypt`

### 2.2 — Hash and verify passwords

**Research topics:**
- `bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)` — returns a hash
- `bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))` — returns nil if match, error if not

**Where to use:**
- Hash the password in the **register handler** before calling `store.Users.Create()`
- Compare in the **login handler** after fetching the user by email

**Important:** The store layer should never see plain-text passwords. Hashing happens in the handler (or a thin service layer). The store only receives and stores the hash.

---

## Phase 3: JWT Tokens

### 3.1 — Add JWT dependency

**Run:** `go get github.com/golang-jwt/jwt/v5`

### 3.2 — Create a JWT helper

**File to create:** `internal/auth/jwt.go` (new package)

**What to build:**
- `GenerateToken(userID int64, secret string, expiry time.Duration) (string, error)` — creates a signed JWT with a `sub` claim (user ID) and `exp` claim
- `ValidateToken(tokenString string, secret string) (int64, error)` — parses and validates, returns the user ID

**Research topics:**
- `jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`
- `jwt.RegisteredClaims` struct — use `Subject` (user ID as string) and `ExpiresAt`
- `jwt.Parse(tokenString, keyFunc)` with `jwt.WithValidMethods([]string{"HS256"})`

### 3.3 — Add JWT_SECRET env var

**File to modify:** `cmd/api/main.go`

- Add `jwtSecret` to the `config` struct in `cmd/api/api.go`
- Load it with `env.GetString("JWT_SECRET", "")` in `main.go`
- Consider requiring it (don't allow empty in production)

---

## Phase 4: Auth Handlers

### 4.1 — Register handler

**File to create:** `cmd/api/users.go`

**`POST /v1/users/register`:**
1. Parse JSON body: `{ "email", "password", "name" }`
2. Validate fields into `errs map[string]string` (email format, password length >= 8, name not empty)
3. Hash password with bcrypt
4. Call `app.store.Users.Create()`
5. Handle unique constraint error (email already taken) → 409 Conflict
6. Generate JWT token
7. Return `{ "user": {...}, "token": "..." }`

### 4.2 — Login handler

**`POST /v1/users/login`:**
1. Parse JSON body: `{ "email", "password" }`
2. Call `app.store.Users.GetByEmail()`
3. If not found → 401 Unauthorized (don't reveal whether email exists)
4. Compare password with bcrypt
5. If mismatch → 401 Unauthorized
6. Generate JWT token
7. Return `{ "user": {...}, "token": "..." }`

### 4.3 — Register routes

**File to modify:** `cmd/api/api.go` (the `mount()` function)

Add:
```
POST /v1/users/register
POST /v1/users/login
```

---

## Phase 5: Auth Middleware

### 5.1 — Create the middleware

**File to create:** `cmd/api/middleware/auth.go`

**What it does:**
1. Read the `Authorization` header → expect `Bearer <token>`
2. If missing/malformed → call `next` but don't set user in context (allow unauthenticated access to public routes)
3. Validate the JWT using your `auth.ValidateToken()`
4. Fetch the user from `store.Users.GetByID()`
5. Store the user in the request context using `context.WithValue`

**Research topics:**
- `strings.TrimPrefix(header, "Bearer ")` to extract the token
- Custom context key type to avoid collisions: `type contextKey string`
- `r.WithContext(ctx)` and `next.ServeHTTP(w, r.WithContext(ctx))`

### 5.2 — Create a "require auth" middleware

This is a second, stricter middleware you wrap around protected routes:

1. Check if user exists in context
2. If not → 401 response
3. If yes → call `next`

### 5.3 — Wire middleware

**File to modify:** `cmd/api/api.go`

**Two approaches to consider:**
- **Option A:** Add the auth middleware to the global stack (it reads the token but doesn't reject — just sets context). Then wrap specific routes with the "require auth" middleware.
- **Option B:** Create a separate mux/handler group for protected routes.

Research Go 1.22 `http.ServeMux` patterns for grouping routes with different middleware.

---

## Phase 6: Protect Existing Routes

Once auth works, decide which routes need protection:

| Route | Auth required? |
|---|---|
| `POST /v1/users/register` | No |
| `POST /v1/users/login` | No |
| `GET /v1/health` | No |
| `GET /v1/issues` | No (public read) |
| `GET /v1/issues/{id}` | No (public read) |
| `POST /v1/issues` | Yes |
| `PATCH /v1/issues/{id}/status` | Yes |
| `DELETE /v1/issues/{id}` | Yes |

You can also associate the logged-in user as the issue author instead of accepting `author` in the request body.

---

## File Summary

| Action | File |
|---|---|
| Create | `cmd/migrate/migrations/000003_create_users.up.sql` |
| Create | `cmd/migrate/migrations/000003_create_users.down.sql` |
| Create | `internal/store/queries/users.sql` |
| Create | `internal/store/users.go` |
| Create | `internal/auth/jwt.go` |
| Create | `cmd/api/users.go` |
| Create | `cmd/api/middleware/auth.go` |
| Modify | `internal/store/storage.go` — add Users interface |
| Modify | `cmd/api/api.go` — add routes + config |
| Modify | `cmd/api/main.go` — load JWT_SECRET |
| Run | `go get golang.org/x/crypto/bcrypt` |
| Run | `go get github.com/golang-jwt/jwt/v5` |
| Run | `sqlc generate` after writing queries |
| Run | Migrate up after writing migration SQL |

---

## Verification

1. **Migrate:** Run migration up, check `users` table exists in Postgres
2. **Register:** `curl -X POST localhost:8080/v1/users/register -d '{"email":"test@test.com","password":"password123","name":"Test"}'` → should return user + token
3. **Login:** `curl -X POST localhost:8080/v1/users/login -d '{"email":"test@test.com","password":"password123"}'` → should return user + token
4. **Protected route without token:** `curl -X POST localhost:8080/v1/issues -d '...'` → should return 401
5. **Protected route with token:** `curl -H "Authorization: Bearer <token>" -X POST localhost:8080/v1/issues -d '...'` → should succeed
6. **Duplicate email:** Register same email again → should return 409

---

## Suggested Learning Order

1. Start with Phase 1 — get the database table and store working first
2. Test it manually by writing a quick handler that creates a user (plain text password temporarily)
3. Then add bcrypt (Phase 2) — replace the plain text with hashing
4. Then add JWT (Phase 3) — this is independent of the database
5. Build the real handlers (Phase 4) — combining store + bcrypt + JWT
6. Add middleware last (Phase 5) — once you can issue and validate tokens
7. Protect routes (Phase 6) — the final wiring
