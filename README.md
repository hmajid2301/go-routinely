# Goroutinely

A cozy, self-hosted habit tracker with a "year in pixels" mood view. Built for the /r/selfhosted community.

## Tech Stack

- **Backend**: Go + pgx + sqlc
- **Frontend**: HTMX + TailwindCSS + DaisyUI (Catppuccin theme)
- **Database**: PostgreSQL + Goose migrations
- **Build**: Nix (gomod2nix) - no Node.js needed!

## Getting Started

```bash
git clone https://gitlab.com/hmajid2301/goroutinely
cd goroutinely

nix develop
# or if you have direnv
direnv allow

task dev
```

The application will be available at `http://localhost:8381`.
