---
name: svelte
description: >
  Svelte 5 runes patterns. Trigger: When writing Svelte components, SvelteKit routes, or Svelte stores.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Writing Svelte components (.svelte)
- Using Svelte 5 runes ($state, $derived, $effect)
- Building SvelteKit routes
- Managing component state

## Critical Patterns

### Svelte 5 runes (REQUIRED)
```svelte
<script>
  let count = $state(0);
  let doubled = $derived(count * 2);

  $effect(() => {
    console.log(`Count is ${count}`);
  });
</script>

<button onclick={() => count++}>
  {count} × 2 = {doubled}
</button>
```

### Props with $props rune
```svelte
<script>
  let { name, age = 18 } = $props();
</script>

<p>{name} is {age} years old</p>
```

### SvelteKit data loading
```typescript
// +page.server.ts
export async function load({ params }) {
  const post = await db.post.findUnique({
    where: { slug: params.slug }
  });
  return { post };
}
```

## Anti-Patterns
### Don't: Use Svelte 4 store syntax in Svelte 5
```svelte
<!-- ❌ Svelte 4 -->
<script>
  import { writable } from 'svelte/store';
  const count = writable(0);
</script>

<!-- ✅ Svelte 5 -->
<script>
  let count = $state(0);
</script>
```

## Quick Reference
| Task | Command |
|------|---------|
| Create | `npx sv create my-app` |
| Dev | `npm run dev` |
| Build | `npm run build` |
| Test | `npm run test` |
