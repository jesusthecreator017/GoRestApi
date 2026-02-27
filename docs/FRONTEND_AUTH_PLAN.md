# Frontend Authentication — Landing Page, Login, Register, Dashboard

## Context

The Go backend already has working JWT auth (`POST /v1/users/register`, `POST /v1/users/login`) and protected routes. The frontend has zero auth code — no context, no token storage, no login/register pages. We need to wire the frontend to the backend auth system so that:

- Unauthenticated users see a **landing page** at `/`
- Authenticated users see the **dashboard** at `/`
- Login (`/login`) and register (`/register`) pages exist with forms
- JWT token is stored in `localStorage` and auto-injected into API calls

---

## Step 1: Create user schemas

**Create:** `web/src/schemas/user.ts`

Follow the pattern in `web/src/schemas/issue.ts`. Define:

- `UserSchema` — `{ id: string (uuid), email: string, name: string, created_at: string, updated_at: string }` → export `User` type
- `LoginInputSchema` — `{ email: string (min 1), password: string (min 8) }` → export `LoginInput` type
- `RegisterInputSchema` — `{ email: string (min 1), name: string (min 1), password: string (min 8) }` → export `RegisterInput` type
- `AuthResponseSchema` — `{ user: UserSchema, token: string }` → export `AuthResponse` type

---

## Step 2: Create auth API module

**Create:** `web/src/api/auth.ts`

Follow the pattern in `web/src/api/issues.ts`. Export an `authApi` object:

- `login(data: LoginInput)` → `apiFetch<AuthResponse>("/v1/users/login", { method: "POST", body: JSON.stringify(data) })`
- `register(data: RegisterInput)` → `apiFetch<AuthResponse>("/v1/users/register", { method: "POST", body: JSON.stringify(data) })`

---

## Step 3: Add token injection to `apiFetch`

**Modify:** `web/src/api/client.ts`

Add a `getToken()` helper that reads `localStorage.getItem("token")` (guarded with `typeof window !== "undefined"` for SSR safety). In `apiFetch`, if a token exists, merge `Authorization: Bearer <token>` into the headers.

---

## Step 4: Create AuthContext

**Create:** `web/src/contexts/AuthContext.tsx` (mark `"use client"`)

This is the core auth state manager:

- **State:** `user: User | null`, `token: string | null`, `isLoading: boolean` (starts `true`)
- **On mount** (`useEffect`): read `token` and `user` (JSON) from `localStorage`, hydrate state, set `isLoading = false`
- **`login(token, user)`**: save both to `localStorage` and state
- **`logout()`**: remove both from `localStorage`, clear state, call `queryClient.clear()` (via `useQueryClient`) to wipe cached data
- Export a **`useAuth()`** hook that calls `useContext` and throws if outside provider

---

## Step 5: Wire AuthProvider into the app

**Modify:** `web/src/app/providers.tsx`

Nest `<AuthProvider>` inside `<QueryClientProvider>` (AuthProvider needs `useQueryClient` for logout):

```tsx
<QueryClientProvider>
  <AuthProvider>
    {children}
  </AuthProvider>
</QueryClientProvider>
```

---

## Step 6: Create auth mutation hooks

**Create:** `web/src/hooks/useAuth.ts`

Follow the pattern in `web/src/hooks/useIssues.ts`. Create TanStack Query mutations:

- `useLogin()` — `useMutation` calling `authApi.login`. On success: `auth.login(data.token, data.user)`, then `router.push("/")`
- `useRegister()` — `useMutation` calling `authApi.register`. Same onSuccess behavior.

Both use `useAuth()` from context and `useRouter()` from `next/navigation`.

---

## Step 7: Create login and register pages

**Create:** `web/src/app/login/page.tsx` (`"use client"`)

- Centered card with email + password inputs (shadcn `Card`, `Input`, `Label`, `Button`)
- Local `useState` for each field
- Submit calls `useLogin().mutate()`
- Show error message from mutation error
- Link to `/register` ("Don't have an account? Register")
- Redirect to `/` if already authenticated (`useAuth().user` is not null)

**Create:** `web/src/app/register/page.tsx` (`"use client"`)

- Same structure but with three fields: name, email, password
- Uses `useRegister()` hook
- Link to `/login` ("Already have an account? Login")
- Redirect to `/` if already authenticated

---

## Step 8: Conditional home page — landing vs dashboard

**Modify:** `web/src/app/page.tsx`

- Import `useAuth` from `@/contexts/AuthContext`
- If `isLoading`: render nothing (or a minimal loader)
- If `user` is null: render a **landing page** — centered heading ("Track your issues"), subtitle, and two `Button`s linking to `/login` and `/register`
- If `user` exists: render `<Dashboard />` as it does now

**Modify:** `web/src/app/Dashboard.tsx`

- Import `useAuth`, destructure `user` and `logout`
- Display user's name somewhere visible
- Add a "Logout" button that calls `logout()` then `router.push("/")`

---

## File Summary

| Order | Action | File |
|-------|--------|------|
| 1 | Create | `web/src/schemas/user.ts` |
| 2 | Create | `web/src/api/auth.ts` |
| 3 | Modify | `web/src/api/client.ts` — inject Authorization header |
| 4 | Create | `web/src/contexts/AuthContext.tsx` |
| 5 | Modify | `web/src/app/providers.tsx` — add AuthProvider |
| 6 | Create | `web/src/hooks/useAuth.ts` |
| 7a | Create | `web/src/app/login/page.tsx` |
| 7b | Create | `web/src/app/register/page.tsx` |
| 8a | Modify | `web/src/app/page.tsx` — conditional landing/dashboard |
| 8b | Modify | `web/src/app/Dashboard.tsx` — show user name + logout |

---

## Gotchas

- **localStorage + SSR**: Every `localStorage` access must be guarded with `typeof window !== "undefined"`. The `AuthProvider` defers reads to `useEffect`; `getToken()` in `client.ts` needs the guard too.
- **Hydration mismatch**: Auth state is unknown during SSR. The `isLoading` flag prevents a flash of wrong UI — the page waits until the client checks localStorage.
- **Query cache on logout**: `queryClient.clear()` in `logout()` prevents one user's data from leaking to the next session.

---

## Verification

1. `pnpm dev` — app loads at `http://localhost:3000`
2. Visit `/` unauthenticated → see landing page with login/register buttons
3. Click "Register" → `/register` form, fill in name/email/password, submit → redirected to `/` showing dashboard with user's name
4. Click "Logout" → back to landing page
5. Visit `/login` → log in with same credentials → dashboard again
6. Create an issue → should work (token sent automatically)
7. Hard refresh → still logged in (token persisted in localStorage)
