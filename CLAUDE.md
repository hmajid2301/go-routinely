### 0 — Purpose
These rules ensure maintainability, safety, and developer velocity.
**MUST** rules are enforced by CI/review; **SHOULD** rules are strong recommendations; **CAN** rules are allowed without extra approval.

---

### 1 — Before Coding
- **BP-1 (MUST)** Ask clarifying questions for ambiguous requirements.
- **BP-2 (MUST)** Draft and confirm an approach (API shape, data flow, failure modes) before writing code.
- **BP-3 (SHOULD)** When >2 approaches exist, list pros/cons and rationale.
- **BP-4 (SHOULD)** Define testing strategy (unit/integration) and observability signals up front.

---

### 2 — Modules & Dependencies
- **MD-1 (SHOULD)** Prefer stdlib; introduce deps only with clear payoff; track transitive size and licenses.
- **MD-2 (CAN)** Use `govulncheck` for updates.

---

### 3 — Code Style
- **CS-1 (MUST)** Enforce `gofmt`, `go vet`
- **CS-2 (MUST)** Avoid stutter in names: `package kv; type Store` (not `KVStore` in `kv`).
- **CS-3 (SHOULD)** Small interfaces near consumers; prefer composition over inheritance.
- **CS-4 (SHOULD)** Avoid reflection on hot paths; prefer generics when it clarifies and speeds.
- **CS-5 (MUST)** Use input structs for function receiving more than 2 arguments. Input contexts should not get in the input struct.
- **CS-6 (SHOULD)** Declare function input structs before the function consuming them.

---

### 4 — Errors
- **ERR-1 (MUST)** Wrap with `%w` and context: `fmt.Errorf("open %s: %w", p, err)`.
- **ERR-2 (MUST)** Use `errors.Is`/`errors.As` for control flow; no string matching.
- **ERR-3 (SHOULD)** Define sentinel errors in the package; document behavior.
- **ERR-4 (CAN)** Use `context.WithCancelCause` and `context.Cause` for propagating error causes.

---

### 5 — Concurrency
- **CC-1 (MUST)** The **sender** closes channels; receivers never close.
- **CC-2 (MUST)** Tie goroutine lifetime to a `context.Context`; prevent leaks.
- **CC-3 (MUST)** Protect shared state with `sync.Mutex`/`atomic`; no "probably safe" races.
- **CC-4 (SHOULD)** Use `errgroup` for fan‑out work; cancel on first error.
- **CC-5 (CAN)** Prefer buffered channels only with rationale (throughput/back‑pressure).

---

### 6 — Contexts
- **CTX-1 (MUST)** If a function takes in `ctx context.Context` it must be the first parameter; never store ctx in structs.
- **CTX-2 (MUST)** Propagate non‑nil `ctx`; honor `Done`/deadlines/timeouts.
- **CTX-3 (CAN)** Expose `WithX(ctx)` helpers that derive deadlines from config.

---

### 7 — Testing
- **T-1 (MUST)** Table‑driven tests; deterministic and hermetic by default.
- **T-2 (MUST)** Run `-race` in CI; add `t.Cleanup` for teardown.
- **T-3 (SHOULD)** Mark safe tests with `t.Parallel()`.

---

### 8 — Logging & Observability
- **OBS-1 (MUST)** Structured logging (`slog`) with levels and consistent fields.
- **OBS-2 (SHOULD)** Correlate logs/metrics/traces via request IDs from context.
- **OBS-3 (CAN)** Provide debug/pprof endpoints guarded by auth or local‑only access.

---

### 9 — Performance
- **PERF-1 (MUST)** Measure before optimizing: `pprof`, `go test -bench`, `benchstat`.
- **PERF-2 (SHOULD)** Avoid allocations on hot paths; reuse buffers with care; prefer `bytes`/`strings` APIs.
- **PERF-3 (CAN)** Add microbenchmarks for critical functions and track regressions in CI.

---

### 10 — Configuration
- **CFG-1 (MUST)** Config via env/flags; validate on startup; fail fast.
- **CFG-2 (MUST)** Treat config as immutable after init; pass explicitly (not via globals).
- **CFG-3 (SHOULD)** Provide sane defaults and clear docs.
- **CFG-4 (CAN)** Support config hot‑reload only when correctness is preserved and tested.

---

### 11 — APIs & Boundaries
- **API-1 (MUST)** Document exported items: `// Foo does …`; keep exported surface minimal.
- **API-2 (MUST)** Accept interfaces where variation is needed; **return concrete types** unless abstraction is required.
- **API-3 (SHOULD)** Keep functions small, orthogonal, and composable.
- **API-5 (CAN)** Use constructor options pattern for extensibility.

---

### 12 — Security
- **SEC-1 (MUST)** Validate inputs; set explicit I/O timeouts; prefer TLS everywhere.
- **SEC-2 (MUST)** Never log secrets; manage secrets outside code (env/secret manager).
- **SEC-3 (SHOULD)** Limit filesystem/network access by default; principle of least privilege.
- **SEC-4 (CAN)** Add fuzz tests for untrusted inputs.

---

### 13 — CI/CD
- **CI-1 (MUST)** Lint, vet, test (`-race`), and build on every PR; cache modules/builds.
- **CI-2 (MUST)** Reproducible builds with `-trimpath`; embed version via `-ldflags "-X main.version=$TAG"`.
- **CI-3 (SHOULD)** Require review sign‑off for rules labeled (MUST).
- **CI-4 (CAN)** Publish SBOM and run `govulncheck`/license checks in CI.

---

### 14 — Tooling
- **TL-1 (CAN)** Linters: `golangci-lint`, `staticcheck`, `gofumpt`.
- **TL-2 (CAN)** Security: `govulncheck`, dependency scanners.
- **TL-3 (CAN)** Testing: `gotestsum`, `mockgen`/`counterfeiter`.

---

### 16 — Tooling Gates (examples)
- **G-1 (MUST)** `task lint` passes.
- **G-2 (MUST)** `task unit:test` and `task integration:test` passes with project config.

---

### Appendix — Rationale for IDs
Stable IDs (e.g., **BP-1**, **ERR-2**) enable precise code‑review comments, changelogs, and automated policy checks (e.g., "violates **CC-2**"). Keep IDs stable; deprecate with notes instead of renumbering.

---

## Writing Functions Best Practices
1. Can you read the function and HONESTLY easily follow what it's doing? If yes, then stop here.
2. Does the function have very high cyclomatic complexity? (number of independent paths, or, in a lot of cases, number of nesting if if-else as a proxy). If it does, then it's probably sketchy.
3. Are there any common data structures and algorithms that would make this function much easier to follow and more robust? Parsers, trees, stacks / queues, etc.
4. Does it have any hidden untested dependencies or any values that can be factored out into the arguments instead? Only care about non-trivial dependencies that can actually change or affect the function.
5. Brainstorm 3 better function names and see if the current name is the best, consistent with rest of codebase.

---

## 17 — Frontend & Templates

### HTML Templates
- **FE-1 (MUST)** Use HTMX for frontend interactivity; return HTML fragments from handlers.
- **FE-2 (MUST)** Use Tailwind CSS for styling; no custom CSS files unless absolutely necessary.
- **FE-3 (MUST)** Follow Bento.io design system principles (see DESIGN-GUIDE.md for details).
- **FE-4 (SHOULD)** Return semantic HTML; use proper elements (`<button>`, `<nav>`, `<main>`, etc.).
- **FE-5 (SHOULD)** Components should be self-contained and reusable.
- **FE-6 (CAN)** Use Alpine.js for client-side state management when needed.

### HTMX Usage
- **HTMX-1 (MUST)** Use `hx-get`, `hx-post`, etc. for AJAX requests.
- **HTMX-2 (MUST)** Return HTML fragments that swap into the DOM; no JSON responses for HTMX requests.
- **HTMX-3 (SHOULD)** Use `hx-target` and `hx-swap` to control DOM updates.
- **HTMX-4 (SHOULD)** Use `hx-trigger` for custom event handling.
- **HTMX-5 (CAN)** Use HTMX extensions when they simplify interactions.

### Tailwind CSS
- **TW-1 (MUST)** Use Tailwind utility classes for all styling.
- **TW-2 (MUST)** Follow the Bento color palette defined in DESIGN-GUIDE.md:
  - Primary: `bg-[#FFBCBA]`, `border-[#FFBCBA]`
  - Primary hover: `hover:bg-[#FFA7A5]`
  - Background: `bg-[#FFF4E9]` (warm cream)
  - Card background: `bg-white`
- **TW-3 (SHOULD)** Use consistent spacing scale: `space-{xs,sm,md,lg,xl,2xl}` maps to `{1,2,4,6,8,12}`.
- **TW-4 (SHOULD)** Use Tailwind's responsive prefixes (`sm:`, `md:`, `lg:`) for breakpoints.
- **TW-5 (CAN)** Extend Tailwind config for custom design tokens if needed.

### Component Structure
- **COMP-1 (MUST)** Handler functions should return `templ.Component` from `a-h/templ`.
- **COMP-2 (MUST)** Keep templates in `internal/templates/` directory.
- **COMP-3 (SHOULD)** Break large pages into smaller component files.
- **COMP-4 (SHOULD)** Use templ's component composition for reusable UI elements.
- **COMP-5 (CAN)** Pass data to templates via Go structs; keep template logic minimal.

### Design Guidelines Reference
- **DG-1 (MUST)** Refer to `mockups/DESIGN-GUIDE.md` for:
  - Exact color values and usage patterns
  - Typography scale and font choices
  - Spacing and layout principles
  - Component design patterns
  - Interactive state styles (hover, active, disabled)
- **DG-2 (SHOULD)** Match mockup designs in `mockups/` directory for consistency.
- **DG-3 (SHOULD)** Use pink accent colors (`rgba(255, 188, 186, 0.204)`, `#ffbcba9e`, `#FFA7A5`) throughout for buttons, borders, and highlights.

### Accessibility
- **A11Y-1 (MUST)** Use semantic HTML and ARIA labels where appropriate.
- **A11Y-2 (SHOULD)** Ensure keyboard navigation works for all interactive elements.
- **A11Y-3 (SHOULD)** Provide sufficient color contrast ratios.
- **A11Y-4 (CAN)** Test with screen readers and keyboard-only navigation.

---

Do not generate any markdown files without me explicitly saying.
