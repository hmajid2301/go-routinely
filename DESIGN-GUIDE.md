# Goroutinely Design Guide

This document defines the design system for Goroutinely, based on Bento.io's warm, approachable aesthetic.

## Design Philosophy

- **Warm & Approachable**: Uses warm cream backgrounds with soft pink accents
- **Clean & Minimal**: Focuses on content, minimal decoration
- **Playful but Professional**: Friendly without being childish
- **GitHub-inspired**: Contribution graphs and streak tracking for habit visualization

---

## Color Palette

### Primary Colors (Bento Pink)

```
Primary:          #FFBCBA   rgb(255, 188, 186)
Primary Hover:    #FFA7A5   rgb(255, 167, 165)
Primary Active:   #E89896   rgb(232, 152, 150)
```

**Tailwind Usage:**
```html
<!-- Solid backgrounds -->
<div class="bg-[#FFBCBA] hover:bg-[#FFA7A5] active:bg-[#E89896]">

<!-- Borders -->
<div class="border-2 border-[#FFBCBA] hover:border-[#FFA7A5]">

<!-- Text -->
<span class="text-[#E89896]">
```

### Background Colors

```
Warm Background:  #FFF4E9   rgb(255, 244, 233)   - Main page background
Card Background:  #FFFFFF   rgb(255, 255, 255)   - Card/container background
Surface:          #FDE5D8   rgb(253, 229, 216)   - Hover surfaces
Mantle:           #FCE9D6   rgb(252, 233, 214)   - Code/input backgrounds
```

**Tailwind Usage:**
```html
<body class="bg-[#FFF4E9]">
  <div class="bg-white rounded-2xl shadow-sm">
    <input class="bg-[#FCE9D6]">
  </div>
</body>
```

### Text Colors

```
Primary Text:     #553630   rgb(85, 54, 48)      - Headings, important text
Base Text:        #333333   rgb(51, 51, 51)      - Body text
Dark Text:        #3a2623   rgb(58, 38, 35)      - Extra emphasis
Neutral:          #8b7b76   rgb(139, 123, 118)   - Secondary text, captions
```

**Tailwind Usage:**
```html
<h1 class="text-[#553630]">Heading</h1>
<p class="text-[#333333]">Body text</p>
<small class="text-[#8b7b76]">Caption</small>
```

### Alpha Variants (Translucent Pink)

```
Light Pink BG:    rgba(255, 188, 186, 0.04)    - Header/nav backgrounds
Border Light:     rgba(255, 188, 186, 0.15)    - Subtle borders
Border Medium:    rgba(255, 188, 186, 0.3)     - Visible borders
Background Light: rgba(255, 188, 186, 0.08)    - Subtle backgrounds
Background Med:   rgba(255, 188, 186, 0.204)   - Active/selected states
Semi-Transparent: #ffbcba9e                    - Buttons with transparency
```

**Tailwind Usage:**
```html
<header class="bg-[rgba(255,188,186,0.04)] border-b border-[rgba(255,188,186,0.15)]">

<div class="bg-[rgba(255,188,186,0.08)] border-2 border-[rgba(255,188,186,0.3)]">

<!-- Active navigation item -->
<a class="bg-[rgba(255,188,186,0.204)]">

<!-- Semi-transparent button -->
<button class="bg-[#ffbcba9e] border-2 border-[#FFBCBA]">
```

### Mood Colors (for mood tracking pixels)

```
Mood 1 (Worst):   #ffb3ba   rgb(255, 179, 186)   - Pastel Rose
Mood 2 (Low):     #ffbe76   rgb(255, 190, 118)   - Mango Orange
Mood 3 (Neutral): #95e1d3   rgb(149, 225, 211)   - Sea Foam
Mood 4 (Good):    #bae1ff   rgb(186, 225, 255)   - Pastel Sky
Mood 5 (Best):    #e0bbff   rgb(224, 187, 255)   - Pastel Lilac
```

**Tailwind Usage:**
```html
<div class="bg-[#ffb3ba]">Worst mood</div>
<div class="bg-[#ffbe76]">Low mood</div>
<div class="bg-[#95e1d3]">Neutral</div>
<div class="bg-[#bae1ff]">Good mood</div>
<div class="bg-[#e0bbff]">Best mood</div>
```

---

## Typography

### Font Family

```
Primary Font:  'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif
```

**Tailwind Setup:**
```js
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'system-ui', 'sans-serif'],
      },
    },
  },
}
```

### Type Scale

```
3XL:    30px / 1.875rem     - Page titles
2XL:    24px / 1.5rem       - Section headings
XL:     20px / 1.25rem      - Card titles
Base:   16px / 1rem         - Body text (default)
SM:     14px / 0.875rem     - Captions, small text
XS:     12px / 0.75rem      - Labels, metadata
```

**Tailwind Usage:**
```html
<h1 class="text-3xl font-bold">Page Title</h1>
<h2 class="text-2xl font-bold">Section Heading</h2>
<h3 class="text-xl font-semibold">Card Title</h3>
<p class="text-base">Body text</p>
<small class="text-sm text-[#8b7b76]">Caption</small>
<span class="text-xs text-[#8b7b76]">Metadata</span>
```

### Font Weights

```
Regular:  400    - Body text
Medium:   500    - Subtle emphasis
Semibold: 600    - Buttons, labels
Bold:     700    - Headings
Extrabold: 800   - Large numbers, stats
```

**Tailwind:**
```html
<p class="font-normal">Regular</p>
<p class="font-medium">Medium</p>
<p class="font-semibold">Semibold</p>
<p class="font-bold">Bold</p>
<p class="font-extrabold">Extra Bold</p>
```

---

## Spacing Scale

```
XS:    4px  / 0.25rem   - Tight gaps
SM:    8px  / 0.5rem    - Small spacing
MD:    16px / 1rem      - Default spacing
LG:    24px / 1.5rem    - Card padding
XL:    32px / 2rem      - Section spacing
2XL:   48px / 3rem      - Page margins
```

**Tailwind Usage:**
```html
<div class="p-1">    <!-- 4px padding -->
<div class="p-2">    <!-- 8px padding -->
<div class="p-4">    <!-- 16px padding -->
<div class="p-6">    <!-- 24px padding -->
<div class="p-8">    <!-- 32px padding -->
<div class="p-12">   <!-- 48px padding -->

<div class="space-y-4">  <!-- 16px vertical gap -->
<div class="gap-6">      <!-- 24px grid gap -->
```

---

## Border Radius

```
SM:   8px  / 0.5rem     - Small elements, inputs
MD:   16px / 1rem       - Cards, buttons
LG:   20px / 1.25rem    - Large cards
XL:   28px / 1.75rem    - Hero sections
Full: 9999px            - Pills, circular buttons
```

**Tailwind:**
```html
<button class="rounded-lg">     <!-- 8px -->
<div class="rounded-2xl">       <!-- 16px -->
<div class="rounded-3xl">       <!-- 24px (closest to 20px) -->
<button class="rounded-full">   <!-- 9999px -->
```

---

## Shadows

```
SM:  0 1px 3px rgba(0, 0, 0, 0.04)      - Subtle elevation
MD:  0 2px 8px rgba(0, 0, 0, 0.06)      - Card hover
LG:  0 4px 16px rgba(0, 0, 0, 0.08)     - Modal, dropdown
Pink SM: 0 2px 8px rgba(255, 167, 165, 0.15)   - Pink accent shadow
Pink MD: 0 4px 12px rgba(255, 167, 165, 0.3)   - Pink hover shadow
```

**Tailwind:**
```html
<div class="shadow-sm">   <!-- Default subtle shadow -->
<div class="shadow-md">   <!-- Card hover -->
<div class="shadow-lg">   <!-- Modal -->

<!-- Custom pink shadows -->
<div class="shadow-[0_2px_8px_rgba(255,167,165,0.15)]">
<div class="shadow-[0_4px_12px_rgba(255,167,165,0.3)]">
```

---

## Components

### Buttons

#### Primary Button
```html
<button class="
  bg-[#ffbcba9e]
  border-2 border-[#FFBCBA]
  shadow-[0_6px_0_#E8908E]
  text-[#333333]
  font-semibold
  rounded-md
  px-8 py-3
  transition-all duration-300
  hover:bg-[#FFA7A5]
  hover:border-[#FFA7A5]
  hover:shadow-[0_8px_0_#D67573]
  hover:-translate-y-0.5
  active:translate-y-1
  active:shadow-[0_2px_0_#D67573]
">
  Primary Action
</button>
```

#### Secondary Button
```html
<button class="
  bg-[rgba(255,188,186,0.08)]
  border-2 border-[rgba(255,188,186,0.3)]
  shadow-[0_6px_0_rgba(255,188,186,0.2)]
  text-[#553630]
  font-medium
  rounded-md
  px-8 py-3
  transition-all duration-300
  hover:bg-[rgba(255,167,165,0.15)]
  hover:border-[rgba(255,167,165,0.5)]
  hover:shadow-[0_8px_0_rgba(255,167,165,0.3)]
  hover:-translate-y-0.5
  active:translate-y-1
">
  Secondary Action
</button>
```

### Cards

```html
<div class="
  bg-white
  rounded-2xl
  p-6
  shadow-sm
  border border-[rgba(255,188,186,0.1)]
  transition-all duration-200
  hover:shadow-[0_2px_8px_rgba(255,167,165,0.12)]
  hover:-translate-y-0.5
  hover:border-[rgba(255,188,186,0.2)]
">
  Card content
</div>
```

### Header

```html
<header class="
  flex items-center justify-between
  px-6 py-4
  bg-[rgba(255,188,186,0.04)]
  border-b border-[rgba(255,188,186,0.15)]
  shadow-[0_1px_3px_rgba(255,167,165,0.08)]
  sticky top-0 z-50
  backdrop-blur-lg
">
  <!-- Header content -->
</header>
```

### Navigation Items

```html
<!-- Active nav item -->
<a class="
  flex flex-col items-center justify-center
  px-4 py-2 rounded-xl
  bg-[rgba(255,188,186,0.204)]
  text-[#944c4b]
  shadow-[0_2px_8px_rgba(255,167,165,0.15)]
  transition-all duration-150
">
  <span class="text-2xl">📊</span>
  <span class="text-xs font-semibold">Overview</span>
</a>

<!-- Inactive nav item -->
<a class="
  flex flex-col items-center justify-center
  px-4 py-2 rounded-xl
  text-[#8b7b76]
  transition-all duration-150
  hover:bg-[rgba(255,167,165,0.1)]
  hover:-translate-y-0.5
">
  <span class="text-2xl">📅</span>
  <span class="text-xs font-semibold">Month</span>
</a>
```

### Stats Cards

```html
<div class="
  bg-[rgba(255,188,186,0.15)]
  border-2 border-[rgba(255,188,186,0.3)]
  rounded-2xl
  p-6
  text-center
">
  <div class="text-4xl font-extrabold text-[#E89896]">21</div>
  <small class="text-sm text-[#944c4b]">🔥 Current Streak</small>
</div>
```

### GitHub-Style Contribution Grid

```html
<div class="flex gap-0.5">
  <!-- Week column -->
  <div class="flex flex-col gap-0.5">
    <!-- Varying intensity using different pink shades -->
    <div class="w-3 h-3 bg-[#E89896] rounded-sm"></div>  <!-- High activity -->
    <div class="w-3 h-3 bg-[#FFA7A5] rounded-sm"></div>  <!-- Medium-high -->
    <div class="w-3 h-3 bg-[#FFBCBA] rounded-sm"></div>  <!-- Medium -->
    <div class="w-3 h-3 bg-[rgba(255,188,186,0.3)] rounded-sm"></div>  <!-- Low -->
    <div class="w-3 h-3 bg-[rgba(255,188,186,0.1)] rounded-sm"></div>  <!-- Very low -->
  </div>
</div>
```

---

## Responsive Breakpoints

```
sm:  640px   - Small tablets
md:  768px   - Tablets
lg:  1024px  - Desktop
xl:  1280px  - Large desktop
2xl: 1536px  - Extra large
```

**Tailwind Usage:**
```html
<div class="
  grid grid-cols-1
  sm:grid-cols-2
  md:grid-cols-3
  lg:grid-cols-4
  gap-4
">
```

---

## HTMX Integration

### Basic HTMX Patterns

```html
<!-- Get request on click -->
<button
  hx-get="/api/habits"
  hx-target="#habit-list"
  hx-swap="innerHTML"
  class="btn-primary"
>
  Load Habits
</button>

<!-- Post form with swap -->
<form
  hx-post="/api/habits/create"
  hx-target="#habit-list"
  hx-swap="beforeend"
  class="space-y-4"
>
  <input type="text" name="name" class="input-field">
  <button type="submit" class="btn-primary">Create</button>
</form>

<!-- Trigger on input change -->
<input
  hx-get="/api/search"
  hx-trigger="keyup changed delay:500ms"
  hx-target="#search-results"
  class="input-field"
>
```

### Loading States

```html
<button
  hx-get="/api/data"
  hx-indicator="#spinner"
  class="btn-primary"
>
  Load Data
  <span id="spinner" class="htmx-indicator ml-2">
    <!-- Spinner icon -->
  </span>
</button>
```

---

## Dark Mode

Dark theme colors (for future implementation):

```
BG Warm:      #1e1e2e
BG Card:      #181825
BG Surface:   #313244
Text:         #cdd6f4
Neutral:      #a6adc8
```

---

## Usage Examples

### Complete Page Layout

```html
<body class="bg-[#FFF4E9]">
  <!-- Header -->
  <header class="
    flex items-center justify-between
    px-6 py-4
    bg-[rgba(255,188,186,0.04)]
    border-b border-[rgba(255,188,186,0.15)]
    shadow-[0_1px_3px_rgba(255,167,165,0.08)]
    sticky top-0 z-50
  ">
    <a href="/" class="flex items-center gap-2 text-xl font-bold text-[#553630]">
      <img src="/logo.png" alt="Logo" class="w-8 h-8">
      Goroutinely
    </a>
    <button class="btn-icon">🌙</button>
  </header>

  <!-- Main Content -->
  <main class="container mx-auto px-4 py-8 pb-24">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold text-[#553630]">Habits</h1>
        <p class="text-sm text-[#8b7b76] mt-1">Track your daily routines</p>
      </div>
    </div>

    <!-- Habit Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div class="
        bg-white rounded-2xl p-6 shadow-sm
        border border-[rgba(255,188,186,0.1)]
        hover:shadow-[0_2px_8px_rgba(255,167,165,0.12)]
        hover:-translate-y-0.5
        transition-all duration-200
      ">
        <h3 class="text-xl font-semibold text-[#553630]">Morning Run</h3>
        <p class="text-sm text-[#8b7b76] mt-2">5km daily</p>
      </div>
    </div>
  </main>

  <!-- Bottom Navigation -->
  <nav class="
    fixed bottom-0 left-0 right-0
    flex justify-around
    px-4 py-2
    bg-white/95 backdrop-blur-lg
    border-t border-[rgba(255,188,186,0.15)]
    shadow-[0_-1px_3px_rgba(255,167,165,0.08)]
    z-50 h-18
  ">
    <a class="flex flex-col items-center justify-center px-4 py-2 rounded-xl bg-[rgba(255,188,186,0.204)]">
      <span class="text-2xl">📊</span>
      <span class="text-xs font-semibold text-[#944c4b]">Overview</span>
    </a>
    <!-- More nav items... -->
  </nav>
</body>
```

---

## References

- Mockups Directory: `/mockups/` - Visual reference for all components
- Bento.io Website: https://warpstreamlabs.github.io/bento/ - Inspiration source
- Color Palette: Based on Bento's warm, approachable aesthetic
- HTMX Documentation: https://htmx.org/
- Tailwind CSS Documentation: https://tailwindcss.com/
