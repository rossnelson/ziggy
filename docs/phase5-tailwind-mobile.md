# Phase 5: Tailwind Migration + Mobile Support

## Overview

Migrate from scoped Svelte CSS to Tailwind CSS and add responsive mobile support with a bottom drawer chat interface.

## Goals

1. Replace inline `<style>` blocks with Tailwind utility classes
2. Add responsive breakpoints for mobile/tablet/desktop
3. Implement bottom drawer chat for mobile view
4. Maintain current desktop layout at larger breakpoints

## Tailwind Setup

### Installation

```bash
cd web
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

### Configuration

**tailwind.config.js**
```js
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        ziggy: {
          green: '#4ade80',
          purple: '#a855f7',
          dark: '#0a0a12',
          panel: 'rgba(26, 26, 46, 0.95)',
        },
      },
      fontFamily: {
        mono: ['monospace'],
      },
    },
  },
  plugins: [],
}
```

**src/app.css**
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

## Breakpoints Strategy

| Breakpoint | Width | Layout |
|------------|-------|--------|
| Mobile | < 640px | Stacked, bottom drawer chat |
| Tablet | 640px - 1024px | Side-by-side, compact |
| Desktop | > 1024px | Current layout |

## Component Migration Plan

### 1. Game.svelte

**Current:** Flex row with controls, canvas, chat side-by-side

**Mobile (<640px):**
- Full-width game canvas centered
- Controls as floating buttons or bottom bar
- Chat as bottom drawer (separate component handles this)

**Desktop (>=640px):**
- Keep current side-by-side layout
- Use `sm:` prefix for desktop styles

### 2. Chat.svelte (Desktop)

Keep as side panel for `sm:` and up breakpoints.

### 3. ChatDrawer.svelte (New - Mobile)

Bottom drawer component for mobile:

**States:**
- `collapsed`: Only grab bar visible (48px)
- `peek`: Input + last message (~120px)
- `half`: 50% viewport height
- `full`: 90% viewport height

**Features:**
- Grab bar with pill indicator
- Touch drag with momentum
- Snap to nearest state on release
- Click grab bar to toggle collapsed/half
- Smooth spring animation

**Structure:**
```svelte
<div class="fixed inset-x-0 bottom-0 z-50 sm:hidden">
  <!-- Grab bar -->
  <div class="h-6 flex justify-center items-center cursor-grab">
    <div class="w-10 h-1 bg-ziggy-green/50 rounded-full" />
  </div>

  <!-- Chat content -->
  <div class="bg-ziggy-panel h-full overflow-hidden">
    <!-- Messages -->
    <!-- Input -->
  </div>
</div>
```

### 4. Controls.svelte

**Mobile:**
- Horizontal row below game canvas
- Or floating action buttons

**Desktop:**
- Keep vertical layout

### 5. Stats.svelte

**Mobile:**
- Compact horizontal bars
- Smaller text

**Desktop:**
- Keep current layout

## Migration Order

1. **Setup Tailwind** - Install, configure, verify working
2. **Create base styles** - Define color palette, common utilities
3. **Migrate Game.svelte** - Container layout, responsive grid
4. **Migrate Stats.svelte** - Simple component, good starting point
5. **Migrate Controls.svelte** - Button styles, responsive layout
6. **Migrate Message.svelte** - Typography, speech bubble
7. **Migrate Ziggy.svelte** - Animations (keep keyframes in CSS)
8. **Migrate Background.svelte** - Gradient backgrounds
9. **Create ChatDrawer.svelte** - New mobile-only component
10. **Migrate Chat.svelte** - Desktop chat, hide on mobile
11. **Final polish** - Test all breakpoints, animations

## File Changes Summary

### New Files
- `tailwind.config.js`
- `postcss.config.js`
- `src/app.css`
- `src/lib/ChatDrawer.svelte`

### Modified Files
- `src/App.svelte` - Import app.css
- `src/lib/Game.svelte` - Tailwind classes, responsive layout
- `src/lib/Chat.svelte` - Tailwind classes, hide on mobile
- `src/lib/Controls.svelte` - Tailwind classes, responsive
- `src/lib/Stats.svelte` - Tailwind classes
- `src/lib/Message.svelte` - Tailwind classes
- `src/lib/Ziggy.svelte` - Tailwind classes (keep animations)
- `src/lib/Background.svelte` - Tailwind classes

## Bottom Drawer Implementation Details

### Touch Handling

```typescript
let startY = 0;
let currentY = 0;
let drawerHeight = 48; // collapsed

function onTouchStart(e: TouchEvent) {
  startY = e.touches[0].clientY;
}

function onTouchMove(e: TouchEvent) {
  const deltaY = startY - e.touches[0].clientY;
  drawerHeight = Math.max(48, Math.min(window.innerHeight * 0.9, drawerHeight + deltaY));
  startY = e.touches[0].clientY;
}

function onTouchEnd() {
  // Snap to nearest state
  const vh = window.innerHeight;
  if (drawerHeight < 80) snapTo('collapsed');
  else if (drawerHeight < vh * 0.3) snapTo('peek');
  else if (drawerHeight < vh * 0.7) snapTo('half');
  else snapTo('full');
}
```

### Snap Heights

```typescript
const snapHeights = {
  collapsed: 48,
  peek: 120,
  half: window.innerHeight * 0.5,
  full: window.innerHeight * 0.9,
};
```

### Animation

Use CSS transition for smooth snapping:
```css
.drawer {
  transition: height 0.3s cubic-bezier(0.32, 0.72, 0, 1);
}
.drawer.dragging {
  transition: none;
}
```

## Testing Checklist

- [ ] Tailwind builds correctly
- [ ] Desktop layout unchanged
- [ ] Mobile layout stacks properly
- [ ] Bottom drawer opens/closes
- [ ] Touch drag works smoothly
- [ ] Snap points feel natural
- [ ] Chat input accessible in all states
- [ ] Messages scroll properly
- [ ] Mystery dropdown works on mobile
- [ ] Controls responsive on all sizes
- [ ] Stats readable on mobile
- [ ] Ziggy animations preserved
- [ ] No flash of unstyled content


## Future Enhancement: 100 HP Celebration

When user reaches 100 HP:
- Confetti animation burst
- Special congratulatory message from Ziggy
- Maybe a subtle glow effect on the HP bar
- Sound effect (optional)

Consider using canvas-confetti library or CSS-based particle effect.

