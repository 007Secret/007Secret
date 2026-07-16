<script setup lang="ts">
import { RouterView } from 'vue-router'
import { computed, onMounted, ref, watch } from 'vue'

const isDarkMode = ref(false)
const themeLabel = computed(() => (isDarkMode.value ? '切换为浅色主题' : '切换为深色主题'))

function applyTheme() {
  document.documentElement.classList.toggle('dark', isDarkMode.value)
}

onMounted(() => {
  const savedTheme = localStorage.getItem('theme')
  isDarkMode.value = savedTheme
    ? savedTheme === 'dark'
    : window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
  applyTheme()
})

watch(isDarkMode, () => {
  applyTheme()
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light')
})

function toggleDarkMode() {
  isDarkMode.value = !isDarkMode.value
}
</script>

<template>
  <div class="relative flex min-h-screen flex-col overflow-hidden bg-[var(--page)] text-[var(--text)]">
    <div aria-hidden="true" class="app-grid pointer-events-none fixed inset-0 opacity-80"></div>
    <div aria-hidden="true" class="pointer-events-none fixed left-1/2 top-[-12rem] h-[34rem] w-[54rem] -translate-x-1/2 rounded-full bg-blue-500/10 blur-3xl"></div>

    <header class="sticky top-0 z-20 border-b border-[var(--line)] bg-[color-mix(in_srgb,var(--surface)_88%,transparent)] backdrop-blur-xl">
      <div class="mx-auto flex h-[72px] w-full max-w-7xl items-center justify-between px-5 sm:px-7 lg:px-8">
        <router-link to="/" class="group flex items-center gap-3" aria-label="007Secret 首页">
          <span class="grid h-10 w-10 place-items-center rounded-xl bg-gradient-to-br from-blue-500 to-blue-700 text-white shadow-lg shadow-blue-600/20 transition-transform group-hover:-translate-y-0.5">
            <svg class="h-5 w-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="M7 10V7a5 5 0 0 1 10 0v3m-9 0h8a2 2 0 0 1 2 2v7H6v-7a2 2 0 0 1 2-2Z" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
              <path d="M12 14v2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </span>
          <span class="text-xl font-extrabold tracking-tight">007Secret</span>
        </router-link>

        <div class="flex items-center gap-3">
          <div class="hidden items-center gap-2 rounded-full border border-blue-200 bg-[var(--primary-soft)] px-3.5 py-2 text-xs font-semibold text-[var(--primary)] sm:flex dark:border-blue-800">
            <span class="h-2 w-2 rounded-full bg-emerald-500 shadow-[0_0_0_4px_rgba(16,185,129,.12)]"></span>
            端到端加密 · 阅后即焚
          </div>
          <button
            type="button"
            :aria-label="themeLabel"
            :title="themeLabel"
            class="grid h-10 w-10 place-items-center rounded-xl border border-[var(--line)] bg-[var(--surface)] text-[var(--muted)] transition hover:-translate-y-0.5 hover:text-[var(--primary)]"
            @click="toggleDarkMode"
          >
            <svg v-if="!isDarkMode" class="h-5 w-5" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M20.5 15.2A8.7 8.7 0 0 1 8.8 3.5 8.7 8.7 0 1 0 20.5 15.2Z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" /></svg>
            <svg v-else class="h-5 w-5" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="4" stroke="currentColor" stroke-width="1.8"/><path d="M12 2v2m0 16v2M4.9 4.9l1.4 1.4m11.4 11.4 1.4 1.4M2 12h2m16 0h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
          </button>
        </div>
      </div>
    </header>

    <main class="relative z-10 flex flex-1">
      <RouterView />
    </main>

    <footer class="relative z-10 px-5 py-5 text-center text-xs text-[var(--muted)]">
      © {{ new Date().getFullYear() }} 007Secret · 为敏感信息而生
    </footer>
  </div>
</template>
