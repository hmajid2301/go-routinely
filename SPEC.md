# Goroutinely - Gamified Habit Tracker

## Core Intent

A self-hosted, gamified habit tracker app aimed at the /r/selfhosted and /r/homelab community. The app centers around a "year in pixels" view for mood tracking, with detailed habit tracking as a secondary layer. Users own their data, deploy easily via Nix or Docker, and access through a PWA with offline support.

## Target Audience

Self-hosting enthusiasts who want:
- Full data ownership and privacy
- Easy deployment via Nix module or Docker
- Modern, cozy pixel art aesthetic with Catppuccin color scheme
- Mobile-first PWA experience with offline capability
- OAuth/OIDC integration with their existing auth infrastructure (Authentik, Authelia, etc.)

## Core Features

### 1. Year in Pixels (Primary Feature)
- **Daily mood tracking**: Users rate each day on a 1-5 scale
- **Visual calendar**: Full year view with color-coded pixels representing mood
- **Responsive layout**:
  - **Mobile**: Show current month, swipe left/right to navigate between months
  - **Desktop**: Show full year (all 12 months) in scrollable view
- **Pixel design**: Rounded squares with visible spacing (modern, cozy feel)
- **Color scheme**: Catppuccin pastel colors (gradient from low to high mood)
- **Interactive**: Click/tap any day to see detailed habit completions in bottom sheet modal
- **Sparse storage**: Only store rated days, generate missing days on-demand in UI
- **Performance target**: < 200ms load time for year view (optimized for 5-20 users)
- **Year view rendering** (Round 14 decision):
  - Server renders full HTML with all 372 pixels (~15-25 KB uncompressed)
  - Rely on gzip/brotli compression (typically reduces to ~3-5 KB over wire)
  - No lazy-loading or virtualization for MVP (simpler architecture)
  - Works without JavaScript (progressive enhancement)
  - Trade-off: Slightly larger payload for simpler implementation and better compatibility
- **Year selector**: Dropdown for 2026, 2025, 2024, etc. in year-overview page
- **Historical data access** (Round 15 decision):
  - Users can select any past year (2024, 2023, etc.) and freely edit/add mood ratings
  - No restrictions on editing old data - full edit access to all years
  - Allows backfilling missed days or correcting past entries
  - No read-only mode for historical years
- **Editable ratings**: Users can go back and edit past mood ratings (any year, any day)
- **Three visual states**:
  - Future dates (grayed out/muted appearance, not interactive)
  - Today (highlighted with distinct border/glow)
  - Past unrated dates (empty state with subtle border, distinct from future)

### 2. Habit Tracking (Secondary Feature)
Support multiple habit types:
- **Binary habits**: Yes/no completion (e.g., "went to gym", "brushed teeth")
- **Numeric habits**: Count-based tracking with optional goals (e.g., "50 pushups")
- **Duration habits**: Time-based tracking (e.g., "30 min workout")
- **Negative habits**: Tracking avoidance (e.g., "didn't smoke")

Habits tracked independently from mood ratings - users can track habits without rating their day, and vice versa.

**Input Validation (Round 14 decisions)**:
- Store habit values as `NUMERIC(10,2)` in Postgres (supports decimals)
- **Count habits**: Enforce integer-only validation (0-10,000 max), reject decimals
- **Duration habits**: Allow decimals (e.g., 15.5 minutes), max 1,440 minutes (24 hours)
- **Client validation**: HTML5 input attributes (type="number", min="0", max="10000", step="1" or step="0.1")
- **Backend validation**: Always re-validate on server (defense in depth), return 400 Bad Request with error message
- **Database constraint**: `CHECK (value >= 0 AND value <= 10000)` on habit_entries table

### 3. Habit Entry UX
- **Real-time quick entry**: Tap/click habit to mark complete immediately (optimized for mobile)
- **Swipe gestures**: Swipe right to complete, left to undo (mobile-native feel)
- **Modal for numeric/duration**: Binary habits one-tap, numeric/duration habits open modal for value entry
- **Optimistic UI**: Show changes immediately, sync in background
- **Touch-optimized**: Large tap targets (min 44x44px), spacious layout
- **Celebration animations**: Special feedback when completing habits (confetti, glow effect)

### 4. Habit Views
- **List view**: All active habits displayed as minimal list items with:
  - Colored circle indicator (habit color)
  - Optional emoji icon (user can pick emoji for each habit)
  - Today's status (completed/not completed)
  - For numeric habits: show percentage for today's entry only (not rolling average)
  - Swipeable for quick actions
- **Individual habit calendar**: See completion history for a specific habit
- **Goal tracking**:
  - Show intensity/percentage for numeric habits (e.g., 25/50 pushups = 50%)
  - Display raw value (e.g., "75 pushups") in day details bottom sheet
  - **Goal exceeding**: Visual indicator (gold star, glow effect) when goal exceeded (>100%)
  - Goals do not affect mood pixel colors (strict separation of concerns)
- **Streaks**: Track consecutive days (deferred to v0.2+)
  - **Streak calculation** (Round 13 decision):
    - Streaks start from habit creation date (not before)
    - Archiving a habit resets streak permanently (no recovery even if recreated)
    - New habit with same name as archived habit = completely independent streak
    - All day boundaries use user's configured timezone (not server UTC)
    - Example: User in PST completes habit at 11:59pm PST → counts for that PST day, not next UTC day
  - Timezone changes recalculate streaks based on new timezone
- **Archived habits**: Hidden from active list but visible in historical day views

### 5. Authentication & Authorization
- **OAuth/OIDC integration**: Support Authentik, Authelia, and other OIDC providers
- **Session management**:
  - Implement OAuth refresh tokens for silent re-authentication
  - Queue offline data across sessions (persist in IndexedDB until synced)
  - If session expires with pending offline data, restore after re-auth
- **Admin approval workflow**:
  - New OAuth users require manual admin approval (configurable via env var)
  - No automatic approval rules (email domain whitelisting deferred to future)
  - Rejected users are permanently blocked based on OAuth sub
- **Role-based access**:
  - Admin role: Approve users, access admin dashboard, view system health, revoke user access
  - User role: Standard access to personal data only
- **In-app admin dashboard**: Admins see pending user approvals and can approve/reject
- **User revocation** (Round 13 decision):
  - Admin can revoke active user's access via "Revoke Access" button
  - Graceful 5-minute sync window: set `revoke_grace_until` timestamp, block after grace period
  - During grace: allow sync API requests only (read/write habit/mood entries), block UI access
  - After grace: return 403 Forbidden on all requests, invalidate all sessions
  - Cancel scheduled push notifications on revocation
  - Preserve user data in database (soft revocation, not deletion)
- **Rate limiting**: Hard rate limits (429 errors) to prevent abuse
  - Per-user limits: e.g., 100 habit entries/hour, 50 mood entries/hour
  - Apply to all users including admins (no exemptions for MVP)
- **User data deletion** (Round 14 decisions):
  - User clicks "Delete All Data" button in profile page
  - Show confirmation modal requiring password re-entry (prevent accidental deletion)
  - Auto-export all user data as JSON file (browser download prompt before deletion)
  - Soft-delete: set `status='deleted'` and `deleted_at` timestamp
  - 30-day recovery window: data remains in database but user cannot login
  - After 30 days: background job permanently purges deleted user records
  - Admin can manually restore deleted user within 30-day window (set status back to 'approved')

### 6. PWA Features
- **Offline-first architecture**: Users can track habits and rate days offline
- **Background sync**: Data syncs automatically when connection restored
- **Conflict resolution**: Last-write-wins (timestamp-based) for multi-device conflicts
- **Offline storage**:
  - Separate IndexedDB stores for mood entries and habit entries
  - Auto-delete oldest offline entries when quota exceeded (FIFO)
  - Mobile Safari quota: ~50MB limit consideration
- **Push notifications**: Daily reminders to track habits/rate day (user opt-in)
- **Install to home screen**: Full PWA manifest for native-like experience
- **Mobile-first design**: Optimized for phone/tablet, works on desktop

### 7. Navigation & Layout
- **Hybrid mobile patterns**:
  - Bottom navigation bar for main sections (Year, Habits, Profile)
  - Floating action button (FAB) for quick actions (add habit, rate day)
  - Bottom sheet for day details (swipeable, maintains context)
- **Desktop patterns**:
  - Top navigation bar or sidebar for main sections
  - Modal overlays for day details
- **HTMX navigation**:
  - SPA-style htmx swaps for main sections (fast, no full reloads)
  - URL updates via hx-boost for browser history
  - Bottom sheet/modal for day details (no page navigation)
- **HTMX response strategy**:
  - Return minimal HTML (single element) for updates
  - Swipe gesture completion: return updated habit list item only
  - Mood rating: return updated pixel element only
  - Optimistic UI updates to ensure responsiveness

### 8. Timezone Handling
- **Per-user timezone setting**: Each user configures their timezone in profile settings
- **UTC storage**: All timestamps stored as UTC in PostgreSQL
- **Client-side display**: Convert to user's timezone for display
- **OAuth timezone inference**: Attempt to get timezone from OAuth claims, fallback to UTC, user can override

## Technology Stack

### Backend
- **Language**: Go
- **Database**: PostgreSQL with sqlc for type-safe queries
- **Migrations**: Goose migrations (separate container in docker-compose, manual for Nix)
- **Auth**: OAuth/OIDC library (go-oidc or similar)
- **Server**: Standard library or lightweight router (chi, echo, fiber)
- **Database pooling**: pgx/v5 with dynamic connection pool (2x CPU cores)

### Frontend
- **Templating**: HTMX for dynamic updates
- **Styling**: TailwindCSS v4+ (built via Nix, no Node.js needed)
- **UI Components**: DaisyUI with custom Catppuccin theme
- **PWA**: Service Worker for offline support, Web Push API for notifications
- **Assets**: Embedded in Go binary via embed.FS
- **Icons**: Optional emoji support for habits (native emoji, no icon library needed for MVP)
- **Build**: Tailwind CLI from nixpkgs, no npm/node required

### Deployment
- **Single binary**: Embed frontend assets and static files (built via buildGoApplication from gomod2nix)
- **Migrations**: Separate from binary (run via goose CLI or docker container)
- **Nix module**: Declare as NixOS service with PostgreSQL integration
- **Docker/Docker Compose**: Official images with easy setup
- **Configuration**: Environment variables for all settings
- **Build tool**: gomod2nix for reproducible Go builds (no vendor directory needed)

## Data Model

### Core Entities

```sql
-- Users (from OAuth)
users:
  - id (uuid, pk)
  - email (string, unique)
  - name (string)
  - oauth_provider (string)
  - oauth_sub (string, unique)
  - role (enum: admin, user)
  - status (enum: pending, approved, rejected, deleted)
  - timezone (string, e.g., "America/New_York")
  - theme_preference (enum: light, dark, auto) -- for Catppuccin Latte/Mocha
  - deleted_at (timestamptz, nullable) -- soft delete with 30-day recovery
  - revoked_at (timestamptz, nullable) -- admin revocation timestamp
  - revoke_grace_until (timestamptz, nullable) -- 5-min sync grace period
  - created_at (timestamptz)
  - updated_at (timestamptz)

-- Mood tracking (year in pixels)
mood_entries:
  - id (uuid, pk)
  - user_id (uuid, fk -> users)
  - date (date) -- stored as DATE type, interpreted in user's timezone
  - rating (int 1-5, not null)
  - notes (text, nullable)
  - created_at (timestamptz)
  - updated_at (timestamptz)
  - unique(user_id, date)

-- Habit definitions
habits:
  - id (uuid, pk)
  - user_id (uuid, fk -> users)
  - name (string, not null)
  - description (text, nullable)
  - type (enum: binary, numeric, duration, negative)
  - goal_value (numeric, nullable) -- for numeric/duration habits
  - goal_unit (string, nullable) -- "pushups", "minutes", etc.
  - color (string, hex color for UI, nullable) -- Catppuccin palette
  - emoji (string, nullable) -- optional emoji icon (single Unicode emoji)
  - archived_at (timestamptz, nullable) -- soft delete
  - created_at (timestamptz)
  - updated_at (timestamptz)

-- Habit completions
habit_entries:
  - id (uuid, pk)
  - habit_id (uuid, fk -> habits)
  - user_id (uuid, fk -> users) -- denormalized for query performance
  - date (date) -- stored as DATE type, interpreted in user's timezone
  - completed (boolean, not null, default false)
  - value (numeric(10,2), nullable) -- for numeric/duration habits (max 10k count, 1440 min)
  - notes (text, nullable)
  - synced_at (timestamptz, nullable) -- for offline sync tracking
  - retry_count (int, default 0) -- optimistic UI retry attempts
  - last_retry_at (timestamptz, nullable) -- last retry timestamp
  - created_at (timestamptz)
  - updated_at (timestamptz)
  - unique(habit_id, date)
  - check(value >= 0 AND value <= 10000) -- validation: count habits max 10k, duration max 1440min
```

### Key Design Decisions
- **Dual independent entities**: Mood entries and habit entries are separate, no hard linkage
- **Sparse storage**: Only store explicit entries, generate missing days in application layer
- **Goal-based intensity**: Habits with numeric goals calculate percentage completion (value/goal)
- **Denormalized user_id**: In habit_entries for faster per-user queries
- **Soft delete**: Use archived_at for habits instead of hard delete
- **Separate habits**: New habit with same name as archived one is completely independent (no linking)
- **DATE type for dates**: Dates stored without timezone, interpreted in user's timezone context
- **Archived habit visibility**: Archived habits hidden from active list but shown in historical day views

### Indexes
```sql
-- Performance-critical indexes
CREATE INDEX idx_mood_entries_user_date ON mood_entries(user_id, date DESC);
CREATE INDEX idx_habit_entries_user_date ON habit_entries(user_id, date DESC);
CREATE INDEX idx_habit_entries_habit_date ON habit_entries(habit_id, date DESC);
CREATE INDEX idx_habits_user_active ON habits(user_id) WHERE archived_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE status = 'pending';
```

## Design System

### Visual Aesthetic
- **Theme**: Warm, cozy, Bento-inspired with custom pastel palette
- **Primary inspiration**: [Bento.io](https://warpstreamlabs.github.io/bento/) professional cozy aesthetic
- **Logo**: Pixel art gopher mascot (Go language mascot) - branding only
- **Pixel art scope**: Logo and mood pixels use pixel art style, UI components are modern
- **UI components**: Bento-style squared buttons (6px radius), soft shadows, spacious layout
- **Density**: Spacious and calm (generous padding, large touch targets, room to breathe)
- **Gopher emotional states**: Hardcoded per context - neutral (login/default), celebrating (habit completion), encouraging (empty states), thinking (modals)
- **Mockups**: Complete HTML/CSS mockups in /mockups/ directory with interactive demos

### Color Palette: Bento-Inspired Custom Pastels

#### Theme Support
- **Light theme**: Warm cream backgrounds (#FFF4E9) with pastel accents
- **Dark theme**: Catppuccin Mocha (dark backgrounds, muted pastels)
- **User preference**: Stored in database, toggle in settings
- **System preference**: Respect `prefers-color-scheme` media query as default

#### Mood Colors (5-level gradient - distinct hues for pattern recognition)
Research-backed colors from year-in-pixels apps (easily distinguishable at a glance):

**All themes (consistent colors across light/dark):**
- Rating 1 (worst): `#ffb3ba` (Pastel Rose) - gentle sadness
- Rating 2 (low): `#ffbe76` (Mango Orange) - muted optimism
- Rating 3 (neutral): `#95e1d3` (Sea Foam) - balanced calm
- Rating 4 (good): `#bae1ff` (Pastel Sky) - peaceful contentment
- Rating 5 (best): `#e0bbff` (Pastel Lilac) - dreamy joy

**Design principle**: Different hues (not just saturations) for instant pattern recognition in year view

#### UI Colors (Light Theme - Bento-inspired)

**Background Hierarchy:**
- Base: `#FFF4E9` (Warm cream) - main background
- Card: `#ffffff` (Pure white) - elevated cards
- Surface: `#FFEFE0` (Peachy cream) - secondary surfaces
- Mantle: `#FCE9D6` (Deeper peachy) - tertiary surfaces

**Action Colors:**
- Primary: `#d32f2f` (Bento red) - important actions, following Bento.io
- Secondary: `#FFD6AF` (Peach) - secondary actions
- Accent: `#EB8788` (Coral red) - highlights and FAB

**Text Colors:**
- Text: `#3a2623` (Dark brown) - primary text
- Neutral: `#654641` (Medium brown) - secondary text

#### UI Colors (Dark Theme - Catppuccin Mocha)

**Background Hierarchy:**
- Base: `#1e1e2e` (Mocha base)
- Card: `#181825` (Slightly darker)
- Surface: `#11111b` (Darkest)

**Action Colors:**
- Primary: `#f5c2e7` (Pink)
- Secondary: `#cba6f7` (Mauve)
- Accent: `#fab387` (Peach)

**Text Colors:**
- Text: `#cdd6f4` (Light text)
- Neutral: `#7c8090` (Gray text)

### Typography

**Font Stack**:
```css
font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
```
- Modern, readable sans-serif (Inter or system default)
- No pixel fonts (better readability, especially on mobile)
- Optional: Pixel art gopher logo uses pixel font in SVG

**Type Scale** (spacious for cozy feel):
- Heading 1: `text-3xl` (30px, for page titles)
- Heading 2: `text-2xl` (24px, for section headers)
- Heading 3: `text-xl` (20px, for cards/modals)
- Body: `text-base` (16px, default)
- Small: `text-sm` (14px, for labels/meta)
- Tiny: `text-xs` (12px, for timestamps)

**Line Height**: Relaxed (`leading-relaxed` = 1.625) for spacious feel

### Spacing & Layout

**Padding/Margin Scale** (generous for spacious feel):
- xs: 4px (`p-1`)
- sm: 8px (`p-2`)
- md: 16px (`p-4`) - default for cards
- lg: 24px (`p-6`) - sections
- xl: 32px (`p-8`) - page padding

**Touch Targets**:
- Minimum: 44x44px (Apple/Google guidelines)
- Preferred: 48x48px for primary actions
- Spacing between tappable elements: ≥ 8px

**Border Radius**:
- Small: `rounded` (4px) - list items, mood pixels
- Buttons: 6px - Bento-style squared buttons (override DaisyUI via tailwind.config.js)
- Medium: `rounded-lg` (8px) - cards
- Large: `rounded-xl` (16px) - modals
- Full: `rounded-full` - FAB, avatars

**DaisyUI Integration Strategy**:
- Use DaisyUI components for structure and utilities
- Override button border-radius in tailwind.config.js: `theme.extend.borderRadius['.btn'] = '6px'`
- Maintain Bento aesthetic while leveraging DaisyUI's component classes

### Pixel Design Specification

**Mood Pixels**:
- Shape: Rounded squares (4px border radius for subtle softness)
- Size (mobile): 40x40px (comfortable tap target)
- Size (desktop): 32x32px (year-overview), 40x40px (month-view)
- Spacing: 4px gap between pixels (`gap-1`)
- Grid layout: CSS Grid for consistent spacing
- States:
  - Empty (past unrated): Light gray border, hollow center, clickable
  - Empty (future): Muted background with dashed outline, not interactive
  - Today: Thick border with Accent color (#EB8788), subtle glow effect
  - Rated (1-5): Filled with corresponding mood color (pastel rose to lilac)

**Calendar Layout**:
- Mobile (portrait): 7 pixels per row (week), month view with swipe navigation
- Desktop (landscape): Month grid view, all 12 months visible
- **Empty months**: Show full calendar grid with all empty pixels (no data placeholders)
  - Maintains consistent layout across all months
  - Invites interaction even for empty periods
  - Empty past dates distinct from future dates visually

### Animation & Motion

**Principles**:
- **Smooth transitions**: 200-300ms for major UI changes (page swaps, modals)
- **Micro-interactions**: 50-100ms for button presses, checkbox checks
- **Respect `prefers-reduced-motion`**: Disable/reduce animations for accessibility
- **Celebration animations**: Special effects for positive actions (habit completion, high mood rating)

**Specific Animations**:
- Modal/bottom sheet open: `slide-up` with `fade-in` (250ms, ease-out)
- Page transitions (htmx): `fade` (200ms) or `slide` (250ms)
- Swipe gestures: Follow finger, spring animation on release (300ms, spring)
- Habit completion: Scale + color change (150ms, ease-out) + confetti burst
- Mood rating: Pixel grows and fills with color (200ms, ease-out)
- FAB press: Scale down (100ms) then up (150ms, bounce)
- Loading states: Skeleton shimmer animation (1500ms loop)

**CSS Transitions**:
```css
/* Default transition */
transition: all 200ms cubic-bezier(0.4, 0, 0.2, 1);

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  transition: none;
  animation: none;
}
```

### Component Specifications

#### Bottom Navigation Bar (Mobile)
- Height: 64px (includes safe area inset)
- 3 items: Year (home icon), Habits (checklist icon), Profile (user icon)
- Active state: Accent color with label
- Inactive state: Neutral color, icon only or with muted label
- Position: `fixed bottom-0`, `z-50`

#### Floating Action Button (FAB)
- Size: 56x56px (standard FAB size)
- Position: Bottom-right, 16px from edge (above bottom nav on mobile)
- Color: Accent (Catppuccin Peach)
- Icon: Plus (+) for add actions, or context-aware (rate day, add habit)
- Shadow: Elevated (`shadow-lg`)
- Animation: Rotate 45° when active, scale on press

#### Bottom Sheet (Day Details Modal)
- Mobile: Slides up from bottom, covers 70-80% of screen
- Swipe handle: Centered drag indicator (40x4px rounded bar)
- Backdrop: Semi-transparent overlay (`bg-black/50`)
- Content: Scrollable, padded (24px)
- Close: Swipe down or tap backdrop

#### Habit List Items
- Layout: Horizontal flex, space between
- Left: Colored circle indicator (16px) + optional emoji + habit name
- Right: Completion state (checkmark or value) + swipe area
- Height: 56px (comfortable tap target)
- Hover/active: Subtle background color change
- Swipe right: Reveal "complete" action (green background)
- Swipe left: Reveal "undo" action (red background)

#### Cards (Settings, Habit Details)
- Background: Surface color (Catppuccin Surface)
- Border: None or subtle (1px, Overlay0)
- Shadow: `shadow-md` for elevation
- Padding: 24px (`p-6`)
- Border radius: `rounded-lg` (8px)

## Scale & Performance Requirements

### Target Scale
- **Primary**: Small community instances (5-20 users)
- **Secondary**: Personal/family use (1-5 users)
- **Performance targets**:
  - Year view load: < 200ms
  - Month view load (mobile): < 100ms
  - Habit entry submission: < 100ms
  - Bottom sheet open: < 50ms (instant feel)
  - PWA offline→online sync: < 5s for typical dataset (365 mood entries + ~20 habits × 365 days)

### Database Optimization
- Indexes on (user_id, date) for both mood_entries and habit_entries
- Partial index for active (non-archived) habits
- **Connection pooling**:
  - Dynamic pool size: 2x CPU cores (scales with hardware)
  - Queue requests with 5s timeout when pool exhausted
  - Better UX than immediate rejection, surfaces capacity issues gradually
- Query optimization via sqlc generated code
- Consider PostgreSQL DATE type optimization (no timezone overhead)
- **Migration strategy**:
  - Transaction-based migrations with automatic rollback on failure
  - App refuses to start if migration fails (prevents data corruption)
  - Requires manual intervention to fix and retry

### Frontend Performance
- Service worker caching for static assets (Tailwind CSS, HTMX, JS)
- Optimistic UI updates (no waiting for server confirmation)
- Lazy loading for year view (render current month first, then adjacent months)
- Debounced sync for offline mode (batch updates every 5s)
- Image optimization: Pixel art gopher logo as optimized SVG (< 10KB)

### Lighthouse Targets
- **Performance**: 95+
- **Accessibility**: 100
- **Best Practices**: 95+
- **SEO**: 90+
- **PWA**: 100

## Deployment Architecture

### Configuration (Environment Variables)
```bash
# Database
DATABASE_URL=postgresql://user:pass@localhost/goroutinely

# OAuth/OIDC
OAUTH_PROVIDER=authentik  # or authelia, generic
OAUTH_CLIENT_ID=goroutinely
OAUTH_CLIENT_SECRET=supersecret
OAUTH_ISSUER_URL=https://auth.example.com/application/o/goroutinely/
OAUTH_REDIRECT_URL=https://habits.example.com/auth/callback
OAUTH_SCOPES=openid,profile,email

# Admin approval
REQUIRE_ADMIN_APPROVAL=true  # default: true
INITIAL_ADMIN_EMAIL=admin@example.com  # first user with this email becomes admin

# Server
PORT=8080
BASE_URL=https://habits.example.com
ENVIRONMENT=production  # development, production

# PWA (optional, for push notifications)
ENABLE_PUSH_NOTIFICATIONS=false
VAPID_PUBLIC_KEY=...
VAPID_PRIVATE_KEY=...
VAPID_SUBJECT=mailto:admin@example.com

# Session
SESSION_SECRET=random-secret-key-change-me
SESSION_DURATION=720h  # 30 days

# Theme (optional, default: auto)
DEFAULT_THEME=auto  # light, dark, auto (follow system preference)
```

### Deployment Options

#### Nix Module
```nix
# Example NixOS configuration
services.goroutinely = {
  enable = true;
  settings = {
    port = 8080;
    oauth = {
      provider = "authentik";
      clientId = "goroutinely";
      issuerUrl = "https://auth.example.com";
    };
    database.createLocally = true;
  };
};
```

**Nix Module Decisions**:
- Single instance per host (one goroutinely service per NixOS system)
- Assumes reverse proxy (Caddy, Nginx) for TLS termination
- Module serves HTTP only on configured port
- Standard self-hosted deployment pattern

#### Docker / Docker Compose
```yaml
# docker-compose.yml example
version: "3.9"
services:
  goroutinely:
    image: ghcr.io/hmajid2301/goroutinely:latest
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://goroutinely:password@db/goroutinely
      OAUTH_CLIENT_ID: ${OAUTH_CLIENT_ID}
      OAUTH_CLIENT_SECRET: ${OAUTH_CLIENT_SECRET}
      OAUTH_ISSUER_URL: ${OAUTH_ISSUER_URL}
      OAUTH_REDIRECT_URL: ${OAUTH_REDIRECT_URL}
      BASE_URL: https://habits.example.com
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    volumes:
      - goroutinely-data:/var/lib/postgresql/data
    environment:
      POSTGRES_DB: goroutinely
      POSTGRES_USER: goroutinely
      POSTGRES_PASSWORD: password
    restart: unless-stopped

volumes:
  goroutinely-data:
```

**Docker Decisions**:
- Built via Nix (`pkgs.dockerTools.buildImage`) - no Dockerfile needed
- Layered image with Nix store paths (efficient caching)
- Runtime image: ~40MB (binary + cacerts for OAuth/HTTPS)
- No base image needed (Nix creates minimal container from scratch)
- Includes only necessary dependencies (binary + cacert package)
- Security: Nix ensures reproducible, minimal attack surface

#### Single Binary
```bash
# Build
go build -o goroutinely ./cmd/server

# Run (requires PostgreSQL running separately)
export DATABASE_URL=postgresql://...
./goroutinely
```

## MVP Scope (v0.1.0)

### Must-Have Features
1. **Auth**: OAuth/OIDC with admin approval workflow
2. **Year in pixels**: Interactive calendar with mood tracking (1-5), adaptive layout (mobile/desktop)
3. **Habit management**: Add/edit/archive habits (all types: binary, numeric, duration, negative)
4. **Habit tracking**: Real-time quick entry with swipe gestures and modals
5. **Day details bottom sheet**: Click pixel to see mood + habits for that day
6. **User profile**: Set timezone, theme preference (light/dark), view account info
7. **Admin dashboard**: Approve/reject pending users
8. **PWA basics**: Manifest, service worker, offline mode with background sync
9. **Responsive UI**: Mobile-first, works on desktop, Catppuccin theme
10. **Single binary deployment**: Embedded assets, auto-migrations

### MVP UI Components
- Bottom navigation bar (mobile) / top nav (desktop)
- Floating action button (FAB) for quick actions
- Bottom sheet for day details (mobile) / modal (desktop)
- Habit list with swipe gestures
- Mood rating selector (1-5 buttons or slider)
- Settings page with theme toggle
- Admin dashboard with pending users table

### Out of Scope for MVP (Future Versions)
- Push notifications (v0.2)
- Habit streaks and statistics (v0.2)
- Social features (compare with friends) (v2.0)
- Habit templates/presets (v0.3)
- Export data (CSV, JSON) (v0.3)
- Advanced theme customization (v0.3)
- Habit categories/tags (v0.4)
- Reminder scheduling (custom times) (v0.3)
- Multi-day view options (week, month detailed) (v0.3)
- Habit insights/analytics (v0.4)
- Email domain whitelist for auto-approval (v0.2)
- Re-request access for rejected users (v0.3)
- User-customizable mood color palette (v0.3)
- Bulk habit entry workflow (v0.2)
- Habit card layout (v0.2, MVP uses list items)

## Technical Implementation Details

This section provides concrete implementation strategies based on the **banterbus** reference project and proven patterns from the self-hosted community. All tooling comes from nixpkgs where possible, eliminating Node.js/npm dependencies.

### Nix Build Strategy

**Build Tool**: gomod2nix for reproducible Go builds
- No need for vendor directory
- Generates `gomod2nix.toml` from go.mod/go.sum
- Reproducible builds across machines
- Follows banterbus/nix-go-htmx-tailwind-template patterns

**Flake Structure** (`flake.nix`):
```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };

        # Development tools from nixpkgs (no npm/node)
        devPackages = with pkgs; [
          go_1_25
          goose          # DB migrations
          air            # Live reload
          golangci-lint  # Linting
          sqlc           # Type-safe SQL
          tailwindcss    # CSS framework (CLI from nixpkgs)
          templ          # Go templating (optional, or use html/template)
          watchman       # File watching
        ];
      in
      {
        # Main package (binary)
        packages.default = pkgs.buildGoApplication {
          pname = "goroutinely";
          version = "0.1.0";
          src = ./.;
          modules = ./gomod2nix.toml;

          # Embed static assets
          preBuild = ''
            ${pkgs.tailwindcss}/bin/tailwindcss -i ./web/static/input.css -o ./web/static/output.css --minify
          '';
        };

        # Docker container
        packages.container = pkgs.dockerTools.buildImage {
          name = "goroutinely";
          tag = "latest";
          created = "now";
          copyToRoot = pkgs.buildEnv {
            name = "image-root";
            paths = [
              self.packages.${system}.default
              pkgs.cacert  # For HTTPS/OAuth
            ];
            pathsToLink = [ "/bin" ];
          };
          config = {
            Cmd = [ "${self.packages.${system}.default}/bin/goroutinely" ];
            ExposedPorts = { "8080/tcp" = {}; };
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
          };
        };

        # Dev shell (no npm/node, everything from nixpkgs)
        devShells.default = pkgs.mkShell {
          packages = devPackages ++ [
            pkgs.gomod2nix.packages.${system}.default
          ];
          shellHook = ''
            export GOOSE_DRIVER=postgres
            export GOOSE_DBSTRING="$DATABASE_URL"
          '';
        };
      }
    );
}
```

**Project Structure**:
```
goroutinely/
├── flake.nix                 # Nix flake config
├── flake.lock               # Locked dependencies
├── gomod2nix.toml           # Generated by gomod2nix
├── go.mod                   # Go dependencies
├── go.sum                   # Go checksums
├── main.go                  # Entry point
├── internal/
│   ├── config/             # Config loading
│   ├── service/            # Business logic
│   ├── store/
│   │   └── db/
│   │       ├── migrations/ # Goose SQL migrations
│   │       ├── query.sql   # sqlc queries
│   │       └── *.go        # Generated by sqlc
│   └── transport/
│       └── http/           # HTTP handlers
├── web/
│   ├── static/
│   │   ├── input.css      # Tailwind source
│   │   ├── output.css     # Generated (gitignored)
│   │   └── js/            # PWA service worker
│   └── templates/          # Go html/template
├── sqlc.yaml               # sqlc config
├── tailwind.config.js      # Tailwind config
└── docker-compose.yml      # Dev environment
```

**Key Benefits**:
- **No Node.js/npm**: Tailwind CSS from nixpkgs (v4+)
- **Reproducible**: gomod2nix locks Go dependencies
- **Single command**: `nix build` produces binary
- **Fast CI**: Nix cache reuses builds
- **Dev consistency**: Same tools across all developers

**Workflow**:
1. Update Go dependencies: `go get ...` → `gomod2nix`
2. Build binary: `nix build`
3. Build container: `nix build .#container`
4. Dev mode: `nix develop` → `air` (live reload)
5. Run migrations: `goose up` (from dev shell)

### Rate Limiting Strategy
- **Per-user limits** (enforced via middleware):
  - Habit entries: 100 per hour
  - Mood entries: 50 per hour
  - Habit CRUD operations: 50 per hour
  - Authentication attempts: 10 per 15 minutes
- **Response**: HTTP 429 Too Many Requests with Retry-After header
- **Storage**: In-memory sliding window counter (reset hourly)
- **MVP**: No Redis dependency, acceptable for 5-20 users
- **Future**: Redis-backed rate limiting for larger instances (v0.3+)

### OAuth Refresh Token Flow
```
1. Initial auth: User logs in, receive access token + refresh token
2. Store refresh token securely (HttpOnly cookie or encrypted localStorage)
3. Access token expires (typically 1 hour)
4. Client detects 401 Unauthorized on API call
5. Automatically attempt refresh using refresh token
6. If refresh succeeds: retry original request with new access token
7. If refresh fails: redirect to login, preserve offline data in IndexedDB
8. After successful re-auth: sync pending offline data
```

**TLS/HTTPS Handling (Round 15 Decision)**:
- Allow `http://` in BASE_URL for any environment (development, production, localhost)
- Log security warnings to stderr if BASE_URL is http:// (not https://)
- Don't block application startup for http:// - admin responsible for TLS via reverse proxy
- OAuth providers may reject http:// redirects (except localhost) - this is provider's choice
- **Rationale**: Permissive for dev/testing, assumes production uses reverse proxy (Caddy, Nginx) for TLS termination

**Health Check Endpoint (Round 15 Decision)**:
- `GET /health` returns simple `200 OK` with body `{"status": "ok"}`
- No dependency checks (no DB ping, no OAuth provider check)
- Just verifies Go process is alive and HTTP server responding
- **Rationale**: Sufficient for MVP scale (5-20 users), simple, fast (<1ms response)
- Future enhancement (v0.3+): Add `/health/detailed` with DB/OAuth checks for monitoring

### Offline Data Persistence Strategy
**IndexedDB Stores**:
- `mood_entries`: {id, date, rating, notes, synced, created_at}
- `habit_entries`: {id, habit_id, date, completed, value, notes, synced, created_at}
- `pending_sync`: Queue of unsynced operations with retry metadata

**Quota Management**:
- Monitor available quota via `navigator.storage.estimate()`
- Warn user at 80% capacity
- Auto-delete oldest unsynced entries (FIFO) at 95% capacity
- Mobile Safari limit: ~50MB (roughly 10,000 habit entries)

**Sync Process**:
1. Queue operations in IndexedDB with `synced: false`
2. Display optimistically in UI immediately (keep visible even on failure)
3. Background sync attempts POST to server
4. On success: mark `synced: true`, keep for 7 days then purge
5. On failure: auto-retry with exponential backoff (1s, 2s, 4s, 8s, max 64s), track retry_count in habit_entries
6. On exhausted retries (5+ attempts): show error toast to user, keep entry in pending state
7. On session expiry: pause sync, resume after re-auth

**Key Decision (Round 14)**: Keep optimistic UI state visible during failures. Auto-retry silently. Only show error if all retries exhausted.

### HTMX Partial Update Strategy
**Single Element Updates** (minimal HTML):
- Habit completion: `<div id="habit-{id}" hx-swap-oob="true">...</div>`
- Mood pixel: `<div id="pixel-{date}" hx-swap-oob="true">...</div>`
- Swipe action: Return just the updated list item with new state

**Out-of-Band Swaps** (hx-swap-oob) - Single Response Multi-Element Update:
- **Primary use case**: Complete habit → update habit item + update stats card + update pixel (if visible)
- Server returns **single HTTP response** with multiple hx-swap-oob elements:
```html
<!-- Primary response: updated habit item -->
<div id="habit-123" class="habit-item">...</div>

<!-- OOB swap: update stats card -->
<div id="weekly-stats" hx-swap-oob="true">
  <div>21 completed</div>
</div>

<!-- OOB swap: update mood pixel if on same page -->
<div id="pixel-2026-01-08" hx-swap-oob="true" class="pixel pixel-rated-4">8</div>
```
- **Benefit**: Atomic multi-element updates in one HTTP request, no race conditions

**Error Handling**:
- 4xx/5xx responses: Show error toast ONLY if retries exhausted (not immediately)
- Network failure: Queue for retry, show "syncing" indicator
- Timeout (5s): Continue retrying in background silently

**Key Decision (Round 15)**: Use hx-swap-oob for atomic multi-element updates, not separate requests or triggers.

### Database Migration Strategy (Goose)

**Migration Tool**: Goose (https://github.com/pressly/goose)
- Industry-standard migration tool
- Simple up/down migrations
- Timestamp-based versioning
- SQL-only (no Go code in migrations for simplicity)

**Migration File Format**:
```sql
-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
```

**Migration Location**: `internal/store/db/migrations/`

**Deployment Approaches**:

1. **Docker Compose** (Development/Testing):
```yaml
migrate:
  image: gomicro/goose:3.25.0
  depends_on:
    postgres:
      condition: service_healthy
  environment:
    POSTGRES_HOST: postgres
    POSTGRES_PORT: 5432
    POSTGRES_USER: postgres
    POSTGRES_PASSWORD: postgres
    POSTGRES_DB: goroutinely
  command: [
    "goose", "-dir", "/migrations", "postgres",
    "host=postgres port=5432 user=postgres password=postgres dbname=goroutinely sslmode=disable",
    "up"
  ]
  volumes:
    - ./internal/store/db/migrations:/migrations:ro
  restart: "no"
```

2. **NixOS Module** (Production):
```nix
# Migrations run via systemd oneshot service before main service starts
systemd.services.goroutinely-migrate = {
  description = "Run Goroutinely database migrations";
  wantedBy = [ "goroutinely.service" ];
  before = [ "goroutinely.service" ];
  serviceConfig = {
    Type = "oneshot";
    RemainAfterExit = true;
    ExecStart = "${pkgs.goose}/bin/goose -dir ${migrations-path} postgres $DATABASE_URL up";
  };
};
```

3. **Manual** (Local development):
```bash
goose -dir internal/store/db/migrations postgres $DATABASE_URL up
goose -dir internal/store/db/migrations postgres $DATABASE_URL status
goose -dir internal/store/db/migrations postgres $DATABASE_URL down
```

**Why NOT Embedded Migrations**:
- Migrations should be explicit and auditable
- Separating migrations from app binary allows independent versioning
- Easier rollback (just run goose down)
- No need to rebuild/redeploy app for schema changes
- Follows banterbus pattern (proven in production)
- Simpler error handling (migration failures don't crash app startup)

**Migration Failure Handling**:
- Goose wraps each migration in a transaction automatically
- On failure: transaction rolls back, migration marked as failed
- Admin must fix issue and re-run `goose up`
- Goose tracks applied migrations in `goose_db_version` table
- No automatic retries (prevents cascading failures)

### Connection Pool Configuration
**Dynamic Sizing**:
```go
numCPU := runtime.NumCPU()
maxConns := numCPU * 2
minConns := max(2, numCPU / 2)

pool, err := pgxpool.New(ctx, connString)
pool.Config().MaxConns = maxConns
pool.Config().MinConns = minConns
pool.Config().MaxConnIdleTime = 5 * time.Minute
pool.Config().HealthCheckPeriod = 1 * time.Minute
```

**Request Queueing**:
- When all connections busy: queue request in memory
- Timeout after 5 seconds: return 503 Service Unavailable
- Queue size limit: 100 pending requests (prevent memory exhaustion)
- Metrics: Log queue depth and timeout events for monitoring

## Key Concerns & Risks

### Scope Creep
- **Risk**: Adding features beyond MVP delays initial release
- **Mitigation**: Strict adherence to MVP feature list, maintain backlog for future versions
- **Decision point**: Before adding any feature, ask "Is this required for initial user feedback?"

### PWA Offline Sync Complexity
- **Risk**: Conflict resolution when multiple devices sync (e.g., phone logs 50 pushups, tablet logs 40)
- **Mitigation**: Last-write-wins strategy based on updated_at timestamp
- **Assumption**: Most users will use single device, multi-device conflicts rare in MVP
- **Future**: Consider operational transformation or manual conflict resolution UI (v0.3+)

### Performance with Sparse Storage
- **Risk**: Generating year view on-demand may be slow with empty state logic
- **Mitigation**:
  - Single query to fetch all user entries for year range
  - Generate missing days in Go (fast, < 1ms for 365 days)
  - Frontend renders incrementally (current month first, then adjacent months)
  - Lazy load past months on scroll (desktop year view)
- **Fallback**: Add Redis caching if query performance degrades (unlikely for 5-20 users)

### OAuth Setup Complexity for Users
- **Risk**: Self-hosters may struggle with OIDC configuration (client ID, secrets, redirect URLs)
- **Mitigation**:
  - Detailed documentation with step-by-step Authentik/Authelia examples
  - Clear error messages when OAuth misconfigured
  - Health check endpoint showing OAuth status (`/health`)
  - Sample `.env` file with all required variables
- **Future**: Consider optional local auth fallback (v2.0+)

### UI/UX Feeling Sluggish
- **Risk**: HTMX requests feel slow, offline mode not smooth, modals laggy, animations janky
- **Mitigation**:
  - Optimistic UI updates (show changes immediately, sync in background)
  - Loading states and skeleton screens for slower operations (< 1% of requests)
  - Aggressive service worker caching (static assets never hit network after first load)
  - Fast server response times (< 100ms target for mutations)
  - Smooth CSS transitions with hardware acceleration (`transform`, `opacity`)
  - Respect `prefers-reduced-motion` for accessibility
  - Test on real low-end devices (not just fast dev machines)
- **Performance budget**:
  - Total JS: < 50KB (HTMX + service worker + minimal app code)
  - Total CSS: < 100KB (Tailwind + DaisyUI purged)
  - First Contentful Paint: < 1.5s on 3G

### Timezone Edge Cases
- **Risk**: Date boundaries incorrect (user in UTC+10 vs UTC-5), habits logged on wrong day
- **Mitigation**:
  - Enforce user timezone setting on first login (detect from browser, confirm with user)
  - Show "Today" date prominently in UI with timezone indicator (e.g., "Today, Jan 5 (PST)")
  - Use PostgreSQL DATE type (no timezone in storage, interpreted client-side)
  - Handle daylight saving time transitions (Go time package handles this automatically)
  - Server always sends UTC, client converts to user's timezone for display
- **Testing**: Test with users in different timezones, especially edge cases (UTC+14, UTC-12, DST transitions)

### Catppuccin Integration with DaisyUI
- **Risk**: DaisyUI theme system may not support Catppuccin colors easily, or conflict
- **Mitigation**:
  - Research DaisyUI custom theme creation before implementation
  - Create separate Catppuccin Latte and Mocha themes in `tailwind.config.js`
  - Test theme toggle extensively (light/dark switching without page reload)
  - Fallback: If DaisyUI incompatible, use raw Tailwind with Catppuccin variables
- **Reference**: https://github.com/catppuccin/catppuccin (official Catppuccin palette)

### Rate Limiting False Positives
- **Risk**: Hard rate limits (429 errors) may block legitimate use cases (bulk import, scripts)
- **Mitigation**:
  - Set limits high enough for normal usage (100 entries/hour = ~1.6/min)
  - Clear error messages explaining limits and reset time
  - Admin can monitor rate limit hits via logs
  - Document limits in API documentation
- **Future**: Add admin override or temporary limit increase (v0.3+)

### IndexedDB Quota Exhaustion
- **Risk**: Auto-deleting oldest entries may cause unexpected data loss if user doesn't sync regularly
- **Mitigation**:
  - Show persistent warning banner at 80% quota
  - Block new entries at 95%, force user to sync or delete manually
  - Document offline storage limits in user guide
  - Test thoroughly on iOS Safari (most restrictive quota)
- **Edge case**: User offline for weeks accumulating data - may hit quota before reconnecting

### HTMX Single-Element Updates
- **Risk**: Returning minimal HTML may cause UI desync if client state differs from server
- **Mitigation**:
  - Include version/hash in element IDs to detect stale updates
  - Fallback to full section re-render if update fails
  - Optimistic UI rolls back on server error
  - Periodic full refresh (e.g., on page focus) to resync state
- **Testing**: Test with slow/flaky network, race conditions

### Migration Failure in Production
- **Risk**: App refuses to start, production instance down until manual fix
- **Mitigation**:
  - Thoroughly test migrations in staging with production-like data
  - Backup database before deploying new version (automated in Nix module)
  - Document recovery procedures for common failures
  - Health check endpoint fails gracefully, shows "maintenance" message
- **Monitoring**: Alert admins immediately on migration failure

## Success Criteria

### Technical
- Year view loads in < 200ms for 5-20 user instance
- PWA scores 90+ on Lighthouse (Performance, PWA, Accessibility, Best Practices)
- Offline mode works reliably: habit entries sync successfully on reconnect
- Single binary deployment under 50MB (including embedded assets)
- Zero-downtime migrations (auto-migrate on startup works without data loss)
- Test coverage > 70% for business logic (habit tracking, mood entries, auth)

### User Experience
- New user can create first habit and track it within 2 minutes (post-OAuth setup)
- Mobile UI feels responsive and native-like (< 100ms tap response)
- Catppuccin aesthetic is cohesive and appealing (positive feedback from beta testers)
- Admin can approve new users in < 30 seconds
- Swipe gestures feel natural on mobile (no accidental triggers)
- Bottom sheet animations are smooth (60 fps on mid-range devices)
- Theme toggle works seamlessly without page reload

### Accessibility
- Keyboard navigation works for all actions (tab, enter, space, escape)
- Screen reader compatible (ARIA labels, semantic HTML)
- Color contrast ratios meet WCAG AA standards (Catppuccin palette is designed for this)
- Touch targets ≥ 44x44px
- `prefers-reduced-motion` respected

### Community & Adoption
- Positive reception on /r/selfhosted (upvoted post, constructive feedback)
- Clear documentation: README with quickstart, deployment guides for Nix/Docker/binary
- Active issue triage: Respond to issues within 48 hours
- At least 3 deployment methods working: Nix module, Docker Compose, binary
- 10+ GitHub stars within first month
- 2+ community contributions (PRs, issues, documentation)

## Decisions Made (Interview Rounds)

### Round 1: Data Model & Auth Foundations
1. **Data modeling**: Dual independent entities (mood_entries and habit_entries separate, no forced linkage between mood and habits)
2. **Missing data handling**: Generate on-demand/sparse storage (only store completed entries, generate missing days in app layer)
3. **Auth strategy**: OAuth with admin approval, configurable via env variables at startup
4. **Pixel view UX**: Year in pixels shows daily mood (1-5 rating), click day to see habits; color gradient based on mood rating

### Round 2: Goals & Entity Relationships
1. **Core entity model**: Dual independent entities - mood and habits fully decoupled, users can track either independently
2. **Goals & targets**: Goal-based intensity/percentage calculation (habits with numeric goals show completion as percentage: 25/50 = 50%)
3. **Admin approval flow**: In-app admin dashboard for approvals + role-based access (admin vs user roles)
4. **Day vs habit independence**: Fully independent - users can track habits without rating days and vice versa

### Round 3: PWA, Design & Deployment
1. **PWA features**: Offline-first with background sync, push notifications for reminders (user opt-in), full PWA functionality (not native app)
2. **Design system**: Custom warm color palette (Catppuccin pastels), modern sans-serif body text, rounded corners and soft shadows, DaisyUI UI library
3. **Scale requirements**: Small community (5-20 users) with < 200ms year view performance target
4. **Deployment**: Single binary with embedded assets, auto-migrate on startup for simplicity

### Round 4: Sync, Colors & Entry UX
1. **Sync conflicts**: Last-write-wins (timestamp-based) for multi-device conflicts (simple, assumes single-device for MVP)
2. **Mood colors**: Fixed color gradient (Catppuccin palette), editable mood ratings (users can change past ratings)
3. **Habit entry UX**: Real-time quick entry (tap to complete), swipe gestures for completion, modal for numeric/duration habits
4. **Admin approval rules**: Manual approval only (no automatic domain whitelist), rejected users permanently blocked

### Round 5: Archival, Timezones & Navigation
1. **Habit archival**: Archived habits show in historical views only (past day details), separate habits (no linking/versioning)
2. **Empty states**: Three visual states - past/unrated, today (highlighted), future (grayed out/distinct)
3. **Timezone handling**: Store UTC with per-user timezone setting, convert on display
4. **Navigation pattern**: Hybrid approach - SPA-style htmx swaps for main sections, bottom sheet (mobile) / modal (desktop) for day details

### Round 6: UI/UX & Design System Details
1. **Grid layout**: Adaptive - mobile shows current month (swipe navigation), desktop shows full year
2. **Pixel art scope**: Pixel art in branding only (gopher logo + mood pixels), UI components are modern
3. **Color palette**: Catppuccin (define now) with Latte (light) and Mocha (dark) theme support
4. **Mobile patterns**: Hybrid - bottom navigation bar, FAB, top nav (support both mobile and desktop patterns)

### Round 7: Design System Refinement
1. **Pixel design**: Rounded squares with spaced grid (modern, cozy feel)
2. **Mood colors**: Catppuccin pastel spectrum (Red → Yellow → Green → Teal → Mauve for 1-5 ratings)
3. **UI colors**: Neutral with warm accents (let mood pixels stand out, use Catppuccin Peach/Pink for primary actions)
4. **Density**: Spacious and calm (generous padding, large touch targets, relaxed line height)

### Round 8: Component Details
1. **Theme support**: Light (Catppuccin Latte) + Dark (Catppuccin Mocha) with user toggle
2. **Habit display**: List items (minimal) with colored circle indicators and optional emoji support
3. **Month navigation**: Swipe left/right (native mobile feel)
4. **Motion design**: Full animation suite (smooth transitions, micro-interactions, celebration animations, respect `prefers-reduced-motion`)

### Round 9: Technical Edge Cases & Implementation
1. **Streak logic**: Streaks always reset on archive, timezone changes recalculate streaks (deferred to v0.2+)
2. **Offline storage**: Separate IndexedDB stores (moods, habits), auto-delete oldest entries when quota exceeded (FIFO)
3. **HTMX granularity**: Return minimal HTML (single element updates for habits/pixels with optimistic UI)
4. **Goal exceeding**: Distinct visual indicator for exceeding goals (gold star, glow effect), goals don't affect mood pixels

### Round 10: Error Handling & Data Integrity
1. **Session expiry**: Implement OAuth refresh tokens, queue offline data across sessions (restore after re-auth)
2. **Empty months**: Show full calendar grid with all empty pixels (consistent layout, no placeholders)
3. **Migration failure**: Transaction-based rollback, refuse to start on failure (requires manual fix)
4. **Habit metrics**: Show only today's value (no rolling averages), display raw value + percentage in day details

### Round 11: Deployment & Scaling
1. **Nix module**: Single instance per host, assume reverse proxy for TLS, serve HTTP only
2. **Connection pooling**: Dynamic pool size (2x CPU cores), queue requests with 5s timeout when exhausted
3. **Docker build**: Built via Nix dockerTools (~40MB with cacerts), no Dockerfile needed
4. **Rate limiting**: Hard rate limits with 429 errors (100 entries/hour), apply to all users including admins

### Round 12: Nix Build & Migrations (Based on Banterbus)
1. **Build system**: gomod2nix for reproducible Go builds, no vendor directory needed
2. **Migrations**: Goose (separate from binary), run via systemd oneshot (NixOS) or docker container (compose)
3. **Frontend build**: Tailwind CSS CLI from nixpkgs (v4+), no Node.js/npm required
4. **Container images**: Nix dockerTools.buildImage (layered, efficient caching, minimal)

### Round 13: Mockup Review & Implementation Details (2026-01-08)
1. **Streak calculations**: Streaks start from habit creation date, archiving resets streak permanently, all day boundaries use user's configured timezone (not server UTC)
2. **Gopher emotional states**: Hardcoded per template/context (login=neutral, completion=celebrating, empty-state=encouraging, modal=thinking) - no dynamic backend logic
3. **DaisyUI integration**: Use DaisyUI components but override border-radius via tailwind.config.js to achieve 6px Bento-style buttons while keeping DaisyUI utilities
4. **User revocation**: Graceful 5-minute sync window before full lockout + cancel all scheduled notifications on revoke (preserve user's pending offline data)

### Round 14: Offline Sync & Data Edge Cases (2026-01-08)
1. **Optimistic UI failures**: Keep optimistic state visible, auto-retry with exponential backoff (1s, 2s, 4s, 8s, max 64s), only show error toast if all retries exhausted
2. **Year view rendering**: Server renders full HTML with all 372 pixels (~15-25 KB), rely on gzip compression, no lazy-loading complexity for MVP (simpler architecture, works without JS)
3. **Delete all data**: Soft-delete with 30-day recovery + require password confirmation modal + auto-export JSON before deletion (triple safety: confirmation, export, recovery window)
4. **Numeric validation**: Store as NUMERIC(10,2), enforce integer-only for count habits, allow decimals for duration, max limits (10k count, 1440 min), client HTML5 + backend re-validation

### Round 15: HTMX Patterns & Deployment Edge Cases (2026-01-08)
1. **HTMX out-of-band swaps**: Single HTTP response with multiple hx-swap-oob elements (habit item + stats card + pixel if needed) for atomic multi-element updates
2. **Historical year access**: Allow full edit access to all past years (2024, 2023, etc.) - users can backfill or correct old data anytime with no restrictions
3. **TLS/OAuth deployment**: Allow http:// in any environment but log security warnings - don't block startup (admin responsible for TLS via reverse proxy, permissive for dev/testing)
4. **Health check endpoint**: Simple 200 OK without dependency checks (dumb health check) - just verifies Go process is alive, sufficient for MVP scale (5-20 users)

## Reference Projects & Inspiration

### Technical References
- **Similar stack**: https://gitlab.com/hmajid2301/banterbus (Go + HTMX + Tailwind + PostgreSQL + sqlc)
- **Template project**: /home/haseeb/projects/nix-go-htmx-tailwind-template (project structure reference)

### Design Inspiration
- **Year in pixels**:
  - https://lifeinpixels.org/
  - https://play.google.com/store/apps/details?id=com.pixa.app&hl=sv
  - https://play.google.com/store/apps/details?id=com.apsterix.apppixel&hl=en-US
  - https://atoms.jamesclear.com/ (James Clear's habit tracker)
- **Color scheme**: https://github.com/catppuccin/catppuccin (official Catppuccin palette)
- **UI components**: https://daisyui.com/ (TailwindCSS component library)

## Next Steps

### Phase 1: Project Setup (Week 1)
1. Initialize Go module with project structure (internal/, web/, migrations/)
2. Set up flake.nix with gomod2nix, dev tools from nixpkgs (no Node.js)
3. Generate gomod2nix.toml from go.mod
4. Configure sqlc for type-safe queries (query.sql → generated Go code)
5. Create initial Goose migrations (users, mood_entries, habits, habit_entries)
6. Set up docker-compose with postgres + goose migrate service
7. Configure Tailwind CSS build in Nix (preBuild step)

### Phase 2: Backend Core (Week 2-3)
1. Implement OAuth/OIDC authentication flow (go-oidc library)
2. Build admin approval workflow (pending users, approve/reject)
3. Create user CRUD operations
4. Implement mood entry API (create, read, update by date range)
5. Implement habit CRUD (create, read, update, archive)
6. Implement habit entry tracking (create, update, fetch by date range)
7. Add timezone handling (per-user setting, UTC storage, conversion utilities)

### Phase 3: Frontend Foundation (Week 3-4)
1. Set up HTMX + Tailwind + DaisyUI build pipeline
2. Create Catppuccin theme configuration (Latte + Mocha)
3. Build base layout templates (header, nav, footer)
4. Implement authentication pages (login, callback, pending approval)
5. Create theme toggle functionality (light/dark, persist to backend)

### Phase 4: Year in Pixels View (Week 4-5)
1. Build year view template (adaptive: full year on desktop, month on mobile)
2. Implement pixel rendering with Catppuccin mood colors
3. Add swipe navigation for mobile month switching
4. Create bottom sheet/modal for day details
5. Implement mood rating selector and submission
6. Add empty states (past/unrated, today highlight, future grayed out)

### Phase 5: Habit Tracking UI (Week 5-6)
1. Build habit list view (minimal list items with colored circles + emojis)
2. Implement swipe gestures for completion (touch events + htmx)
3. Create habit add/edit modal with type selection
4. Build quick completion flow (one-tap for binary, modal for numeric)
5. Add celebration animations (confetti, glow effects)
6. Implement FAB for quick actions

### Phase 6: PWA & Offline Support (Week 6-7)
1. Create service worker with caching strategy
2. Build PWA manifest with Catppuccin theme colors
3. Implement offline storage (IndexedDB for habit entries)
4. Add background sync for offline-to-online data transfer
5. Create sync status indicator in UI
6. Test offline mode thoroughly (airplane mode, slow 3G)

### Phase 7: Admin & Settings (Week 7)
1. Build admin dashboard (pending users table, approve/reject actions)
2. Create user profile page (timezone setting, theme preference, account info)
3. Implement settings page (theme toggle, timezone picker)
4. Add role-based access control (admin vs user routes)

### Phase 8: Polish & Testing (Week 8)
1. Optimize database queries (add indexes, test with 20 users × 365 days data)
2. Run Lighthouse audits, fix performance issues
3. Test across devices (iOS Safari, Android Chrome, desktop browsers)
4. Add loading states and skeleton screens
5. Write unit tests for business logic (target 70% coverage)
6. Fix accessibility issues (keyboard nav, screen reader, contrast)

### Phase 9: Deployment & Documentation (Week 9)
1. Create Nix module for NixOS deployment
2. Build Docker image and docker-compose.yml
3. Write deployment documentation (README, guides for each method)
4. Create OAuth setup guides (Authentik, Authelia step-by-step)
5. Add health check endpoint and monitoring
6. Set up CI/CD (GitHub Actions for builds and tests)

### Phase 10: Beta Release (Week 10)
1. Deploy beta instance for testing
2. Invite beta testers (friends, /r/selfhosted community)
3. Gather feedback and iterate
4. Fix critical bugs
5. Prepare v0.1.0 release announcement
6. Publish to GitHub with releases and changelog

---

## Specification Summary

This specification was developed through **15 rounds of structured interviews** covering:

**Rounds 1-8: Core Architecture & Design**
- Data modeling, authentication, PWA features
- Catppuccin color system, Tailwind + DaisyUI integration
- Mobile-first UX with swipe gestures, bottom sheets, FAB
- Adaptive layouts (mobile month view, desktop year view)

**Rounds 9-12: Technical Implementation & Build System**
- Edge cases: streaks, offline storage, goal exceeding
- Error handling: session expiry, empty data, migration failures
- Deployment: Nix modules, Docker via Nix, connection pooling
- Performance: rate limiting, HTMX partial updates, IndexedDB quota management
- Build system: gomod2nix, Goose migrations, Tailwind from nixpkgs

**Rounds 13-15: Mockup Review & Implementation Refinement (2026-01-08)**
- Streak calculation details (creation date start, archive reset, user timezone boundaries)
- Gopher emotional states (hardcoded per context, no backend logic)
- DaisyUI styling override strategy (6px border-radius via Tailwind config)
- User revocation with graceful sync period (5-min window + notification cleanup)
- Optimistic UI retry strategy (exponential backoff, keep visible state)
- Year view rendering approach (full HTML 372 pixels, gzip-compressed)
- Data deletion safety (soft-delete 30-day, password confirm, auto-export)
- Input validation rules (NUMERIC(10,2), integer count, decimal duration, max limits)
- HTMX multi-element updates (single response with hx-swap-oob)
- Historical data editing (full access to all past years)
- TLS permissiveness (allow http:// with warnings, don't block)
- Health check simplicity (200 OK only, no dependency checks)

**Key Technical Decisions**:
- **Nix-first build system**: gomod2nix + dockerTools (no Node.js/npm, everything from nixpkgs)
- **Goose migrations**: Separate from binary, explicit deployment steps
- **Offline-first PWA** with IndexedDB (separate stores, FIFO deletion)
- **OAuth refresh tokens** with cross-session data persistence
- **HTMX minimal updates** with optimistic UI and out-of-band swaps
- **Dynamic connection pooling** (2x CPU cores, request queueing)
- **Hard rate limiting** (429 errors, no admin exemption)
- **Nix Docker images** via dockerTools (~40MB, layered, reproducible)
- **Single Nix instance** per host, assuming reverse proxy

**Deployment Targets**:
- Small community instances (5-20 users)
- Year view < 200ms, habit entry < 100ms
- PWA Lighthouse score 90+
- Single binary < 30MB (Go + embedded assets)
- Docker runtime image ~40MB (includes cacerts)

**MVP Scope (v0.1.0)**:
10 core features with deferred enhancements to v0.2-v0.4. Focus on getting initial user feedback without scope creep.

This specification is implementation-ready with clear technical requirements, error handling strategies, and deployment patterns.
