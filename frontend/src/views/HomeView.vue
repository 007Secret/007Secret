<script setup lang="ts">
import { computed, ref } from 'vue'

const content = ref('')
const showResult = ref(false)
const shareLink = ref('')
const password = ref('')
const loading = ref(false)
const createError = ref('')
const copyMessage = ref('')

const canSubmit = computed(() => content.value.trim().length > 0 && !loading.value)

async function createSecret() {
  if (!canSubmit.value) return
  loading.value = true
  createError.value = ''

  try {
    const response = await fetch('/api/secret', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: content.value }),
    })
    if (!response.ok) throw new Error('Failed to create secret')

    const data = await response.json()
    const url = new URL(`/s/${data.key}`, window.location.origin)
    url.searchParams.set('password', data.password)
    shareLink.value = url.toString()
    password.value = data.password
    showResult.value = true
  } catch {
    createError.value = '创建分享失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

async function copy(text: string, label: string) {
  try {
    await navigator.clipboard.writeText(text)
    copyMessage.value = `${label}已复制`
    window.setTimeout(() => (copyMessage.value = ''), 2000)
  } catch {
    copyMessage.value = '复制失败，请手动复制'
  }
}

function reset() {
  content.value = ''
  shareLink.value = ''
  password.value = ''
  createError.value = ''
  showResult.value = false
}
</script>

<template>
  <div class="mx-auto flex w-full max-w-7xl items-center px-5 py-12 sm:px-7 sm:py-16 lg:px-8">
    <div class="grid w-full items-center gap-12 xl:grid-cols-[minmax(0,1fr)_minmax(500px,1.03fr)] xl:gap-20">
      <section class="text-center xl:text-left">
        <div class="mb-5 inline-flex items-center gap-3 text-sm font-bold text-[var(--primary)]">
          <span class="h-2 w-2 rounded-full bg-[var(--primary)] shadow-[0_0_0_6px_rgba(37,99,235,.1)]"></span>
          一次性秘密分享
        </div>
        <h1 class="text-[clamp(2.7rem,5vw,4.6rem)] font-black leading-[1.08] tracking-[-0.05em]">
          让秘密，只被看见<span class="text-[var(--primary)]">一次</span>
        </h1>
        <p class="mx-auto mt-6 max-w-xl text-base leading-8 text-[var(--muted)] sm:text-lg xl:mx-0">
          安全分享密码、密钥与敏感信息。无需注册，读取后自动销毁。
        </p>

        <div class="mx-auto mt-10 grid max-w-xl gap-4 text-left sm:grid-cols-3 xl:mx-0">
          <div class="border-t border-[var(--line)] pt-4"><strong class="block">一次查看</strong><span class="text-xs text-[var(--muted)]">接收者只能读取一次</span></div>
          <div class="border-t border-[var(--line)] pt-4"><strong class="block">自动销毁</strong><span class="text-xs text-[var(--muted)]">读取后立即永久删除</span></div>
          <div class="border-t border-[var(--line)] pt-4"><strong class="block">无需注册</strong><span class="text-xs text-[var(--muted)]">打开即可安全分享</span></div>
        </div>
      </section>

      <section class="app-surface overflow-hidden rounded-[24px] p-6 sm:p-8">
        <form v-if="!showResult" class="space-y-5" @submit.prevent="createSecret">
          <div>
            <div class="mb-1 flex items-start justify-between gap-4">
              <div><h2 class="text-2xl font-extrabold tracking-tight sm:text-3xl">创建一次性秘密</h2><p class="mt-1 text-sm text-[var(--muted)]">内容将在被读取后永久销毁</p></div>
              <span class="mt-1 rounded-full bg-[var(--primary-soft)] px-3 py-1 text-xs font-bold text-[var(--primary)]">加密</span>
            </div>
          </div>

          <div>
            <div class="mb-2 flex items-center justify-between"><label for="secret-content" class="text-sm font-bold">秘密内容</label><span class="text-xs text-[var(--muted)]">{{ content.length }} / 10,000</span></div>
            <textarea
              id="secret-content"
              v-model="content"
              aria-label="秘密内容"
              maxlength="10000"
              rows="8"
              placeholder="输入你想安全分享的内容…"
              class="min-h-48 w-full resize-none rounded-2xl border border-[var(--line)] bg-[var(--surface-soft)] px-4 py-3 text-[var(--text)] outline-none transition placeholder:text-slate-400 focus:border-[var(--primary)] focus:ring-4 focus:ring-blue-500/10"
            ></textarea>
          </div>

          <p v-if="createError" role="alert" class="rounded-xl border border-red-300/60 bg-red-500/10 px-4 py-3 text-sm font-medium text-[var(--danger)]">{{ createError }}</p>

          <div class="flex items-start gap-2.5 text-xs leading-5 text-[var(--muted)]">
            <svg class="mt-0.5 h-4 w-4 flex-none text-[var(--primary)]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m12 3 8 3v5c0 5-3.4 8.4-8 10-4.6-1.6-8-5-8-10V6l8-3Z" stroke="currentColor" stroke-width="1.8"/><path d="m9 12 2 2 4-4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/></svg>
            内容将被加密存储，我们无法读取
          </div>

          <button type="submit" :disabled="!canSubmit" class="flex h-14 w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-b from-blue-500 to-blue-700 font-bold text-white shadow-lg shadow-blue-600/20 transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-45 disabled:hover:translate-y-0">
            <svg v-if="loading" class="h-5 w-5 animate-spin" viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="9" stroke="currentColor" stroke-opacity=".3" stroke-width="3"/><path d="M21 12a9 9 0 0 0-9-9" stroke="currentColor" stroke-width="3" stroke-linecap="round"/></svg>
            <svg v-else class="h-5 w-5" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M7 10V7a5 5 0 0 1 10 0v3m-9 0h8a2 2 0 0 1 2 2v7H6v-7a2 2 0 0 1 2-2Z" stroke="currentColor" stroke-width="1.8"/></svg>
            {{ loading ? '正在生成…' : '生成安全链接' }}
          </button>
        </form>

        <div v-else class="space-y-5">
          <div class="flex items-start gap-4"><span class="grid h-11 w-11 flex-none place-items-center rounded-full bg-emerald-500/12 text-xl text-[var(--success)]">✓</span><div><p class="text-sm font-bold text-[var(--success)]">创建成功</p><h2 class="mt-1 text-2xl font-extrabold sm:text-3xl">安全链接已生成</h2><p class="mt-1 text-sm text-[var(--muted)]">请将链接与提取码分别发送给接收者</p></div></div>

          <div class="border-t border-[var(--line)] pt-5"><label for="share-link" class="mb-2 block text-sm font-bold">分享链接</label><div class="flex flex-col gap-2 sm:flex-row"><input id="share-link" :value="shareLink" readonly class="h-13 min-w-0 flex-1 rounded-xl border border-[var(--line)] bg-[var(--surface-soft)] px-4 font-mono text-sm text-[var(--primary)]"><button class="h-13 rounded-xl bg-[var(--primary)] px-5 font-bold text-white" @click="copy(shareLink, '链接')">复制链接</button></div></div>

          <div><label for="extract-code" class="mb-2 block text-sm font-bold">提取码</label><div class="flex flex-col gap-2 sm:flex-row"><input id="extract-code" :value="password" readonly class="h-14 min-w-0 flex-1 rounded-xl border border-[var(--line)] bg-[var(--surface-soft)] px-4 text-center font-mono text-2xl font-black tracking-[.3em]"><button class="h-14 rounded-xl bg-[var(--primary)] px-5 font-bold text-white" @click="copy(password, '提取码')">复制提取码</button></div></div>

          <div class="flex gap-3 rounded-xl border border-amber-400/50 bg-amber-500/10 p-4 text-sm text-[var(--warning)]"><span aria-hidden="true">⚠</span><div><strong class="block">仅可查看一次</strong><span>内容被读取后会立即销毁，且无法恢复</span></div></div>
          <button class="h-13 w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] font-bold text-[var(--primary)] transition hover:bg-[var(--primary-soft)]" @click="reset">创建新的分享</button>
        </div>
      </section>
    </div>

    <Transition name="toast"><div v-if="copyMessage" role="status" class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 rounded-xl bg-slate-950 px-4 py-2.5 text-sm font-semibold text-white shadow-xl">{{ copyMessage }}</div></Transition>
  </div>
</template>

<style scoped>
.h-13 { height: 3.25rem; }
.toast-enter-active, .toast-leave-active { transition: opacity .2s, transform .2s; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, .5rem); }
</style>
