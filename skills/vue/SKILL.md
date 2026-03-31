---
name: vue
description: >
  Vue 3 Composition API patterns. Trigger: When writing Vue components, composables, or Vue/Nuxt applications.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Writing Vue 3 components (.vue)
- Creating composables (useXxx)
- Building Vue/Nuxt applications
- Using Pinia for state management

## Critical Patterns

### Composition API with script setup (REQUIRED)
```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const count = ref(0)
const doubled = computed(() => count.value * 2)

onMounted(() => {
  console.log('Component mounted')
})
</script>

<template>
  <button @click="count++">{{ count }} × 2 = {{ doubled }}</button>
</template>
```

### Composable pattern (REQUIRED)
```typescript
// useCounter.ts
export function useCounter(initial = 0) {
  const count = ref(initial)
  const increment = () => count.value++
  const decrement = () => count.value--
  const reset = () => count.value = initial
  return { count, increment, decrement, reset }
}
```

### Pinia store
```typescript
// stores/user.ts
export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isLoggedIn = computed(() => !!user.value)
  async function login(creds: Credentials) { /* ... */ }
  return { user, isLoggedIn, login }
})
```

## Anti-Patterns
### Don't: Use Options API in new projects
```vue
<!-- ❌ Options API -->
<script>
export default {
  data() { return { count: 0 } },
  methods: { increment() { this.count++ } }
}
</script>
```

## Quick Reference
| Task | Command |
|------|---------|
| Create | `npm create vue@latest` |
| Dev | `npm run dev` |
| Build | `npm run build` |
| Test | `vitest` |
