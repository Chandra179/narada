# Narada

Single-page React + Vite app for Bazi / elemental personality analysis. All client-side, no backend.

## Commands

| Action | Command |
|--------|---------|
| Dev server | `npm run dev` |
| Build | `npm run build` (runs `tsc -b` then `vite build`) |
| Lint | `npm run lint` (oxlint) |
| Preview | `npm run preview` |

## Architecture

- Entrypoint: `src/main.tsx` → renders `<App/>`
- No router — screens are a manual state machine (`'onboard' | 'deep' | 'dashboard'`) in `App.tsx`
- State persists to `localStorage` under key `wu-xing-state-v1` via `src/hooks/usePersistentState.ts`
- All calculations are deterministic, pure functions in `src/lib/calculations.ts`
- City dataset is hardcoded in `src/lib/constants.ts`
- `esoteric-analysis-app.html` is the pre-React prototype — not the active source

## TypeScript gotchas

- `erasableSyntaxOnly: true` — no `enum`, no `namespace`, no parameter properties
- `verbatimModuleSyntax: true` — always use `import type` for type-only imports
- `noUnusedLocals` / `noUnusedParameters` are errors

## Tailwind v4

Uses `@import "tailwindcss"` and `@theme` directive in `src/index.css`. Custom design tokens are declared in the `@theme` block, NOT in `tailwind.config.*`. No v3 compatibility.

## No tests

No test framework configured. No test files exist.

## Go server (dev dependency)

- `server/main.go` + `server/calculations.go` — `POST /api/profile` takes `{"birthdate":"YYYY-MM-DD"}`, returns computed profile JSON
- Start: `go run .` from the `server/` directory (listens on `:8080`)
- Vite proxies `/api` → `localhost:8080` in dev via `vite.config.ts`
- Build does NOT compile the Go binary — server must be started separately

## Styling conventions

- Fonts: Fraunces (display), Manrope (body), IBM Plex Mono (mono) — loaded from Google Fonts in `index.html`
- CSS custom properties mirror Tailwind theme tokens
- Unlock-CTA button uses `.unlock-cta` class in global CSS, not a Tailwind utility class
