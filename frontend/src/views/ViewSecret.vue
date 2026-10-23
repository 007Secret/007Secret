<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { receiveAttachment } from '../lib/attachmentTransfer'

const router = useRouter()
const route = useRoute()
const password = ref('')
const content = ref('')
const showContent = ref(false)
const showTransfer = ref(false)
const error = ref('')
const loading = ref(false)
const copyMessage = ref('')
const attachmentStatus = ref('')
let stopReceivingAttachment: (() => void) | undefined

onMounted(() => {
  const queryPassword = route.query.password
  password.value = Array.isArray(queryPassword) ? queryPassword[0] ?? '' : queryPassword ?? ''
})

async function viewSecret() {
  if (!password.value.trim() || loading.value) return
  const key = route.params.key as string
  if (!key) {
    error.value = '无效的密钥链接'
    return
  }

  loading.value = true
  error.value = ''
  showTransfer.value = false
  stopReceivingAttachment?.()
  stopReceivingAttachment = undefined
  try {
    const metadataResponse = await fetch(`/api/secret/attachment?key=${key}&password=${encodeURIComponent(password.value)}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    })
    if (!metadataResponse.ok) {
      if (metadataResponse.status === 404) throw new Error('密钥不存在或已被查看')
      if (metadataResponse.status === 401) throw new Error('提取码错误')
      throw new Error('服务器错误，请稍后重试')
    }
    const metadata = await metadataResponse.json()
    if (!metadata.attachment) {
      await consumeSecret(key)
      return
    }

    showTransfer.value = true
    attachmentStatus.value = '正在等待发送方建立 WebRTC 点对点连接…'
    stopReceivingAttachment = receiveAttachment({ key, password: password.value }, {
      onStatus: (message) => { attachmentStatus.value = message },
      onProgress: (transferred, total) => {
        attachmentStatus.value = `正在接收附件：${Math.round(transferred / total * 100)}%`
      },
      onComplete: async (file) => {
        const url = URL.createObjectURL(file)
        const link = document.createElement('a')
        link.href = url
        link.download = file.name
        link.click()
        URL.revokeObjectURL(url)
        attachmentStatus.value = '附件已接收并开始下载'
        try {
          await consumeSecret(key)
        } catch (err) {
          attachmentStatus.value = `${attachmentStatus.value}；秘密内容无法读取：${err instanceof Error ? err.message : '请稍后重试'}`
        }
      },
      onError: (message) => {
        attachmentStatus.value = message
      },
    })
  } catch (err) {
    error.value = err instanceof Error ? err.message : '未知错误，请稍后重试'
  } finally {
    loading.value = false
  }
}

async function consumeSecret(key: string) {
  const response = await fetch(`/api/secret/get?key=${key}&password=${encodeURIComponent(password.value)}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
    })
  if (!response.ok) {
    if (response.status === 404) throw new Error('密钥不存在或已被查看')
    if (response.status === 401) throw new Error('提取码错误')
    throw new Error('服务器错误，请稍后重试')
  }
  const data = await response.json()
  content.value = data.content
  showTransfer.value = false
  showContent.value = true
}

async function copyContent() {
  try {
    await navigator.clipboard.writeText(content.value)
    copyMessage.value = '内容已复制'
    window.setTimeout(() => (copyMessage.value = ''), 2000)
  } catch {
    copyMessage.value = '复制失败，请手动复制'
  }
}

function goHome() {
  router.push('/')
}
</script>

<template>
  <div class="mx-auto flex w-full max-w-5xl items-center px-5 py-12 sm:px-7 sm:py-16 lg:px-8">
    <div class="w-full text-center">
      <template v-if="!showContent && !showTransfer && !error">
        <div class="mb-4 inline-flex items-center gap-3 text-sm font-bold text-[var(--primary)]"><span class="h-2 w-2 rounded-full bg-[var(--primary)] shadow-[0_0_0_6px_rgba(37,99,235,.1)]"></span>一次性秘密</div>
        <h1 class="text-[clamp(2.5rem,6vw,4rem)] font-black leading-tight tracking-[-0.045em]">查看秘密内容</h1>
        <p class="mx-auto mt-4 max-w-2xl text-base text-[var(--muted)] sm:text-lg">输入提取码后，秘密将被读取并立即销毁</p>
        <div class="mx-auto my-7 flex max-w-xl items-center justify-center gap-3 text-xs text-[var(--muted)] sm:gap-6 sm:text-sm"><strong class="text-[var(--primary)]">① 验证提取码</strong><span>→</span><span>② 查看内容</span><span>→</span><span>③ 自动销毁</span></div>

        <form class="app-surface mx-auto max-w-2xl space-y-5 rounded-[24px] p-6 text-left sm:p-8" @submit.prevent="viewSecret">
          <div><label for="password" class="mb-2 block text-sm font-bold">提取码</label><input id="password" v-model="password" aria-label="提取码" type="text" autocomplete="one-time-code" spellcheck="false" placeholder="请输入提取码" class="h-16 w-full rounded-2xl border border-[var(--line)] bg-[var(--surface-soft)] px-5 text-center font-mono text-2xl font-black uppercase tracking-[.3em] text-[var(--text)] outline-none transition placeholder:text-base placeholder:font-normal placeholder:tracking-normal focus:border-[var(--primary)] focus:ring-4 focus:ring-blue-500/10"></div>
          <div class="flex items-center gap-2 text-xs text-[var(--muted)]"><svg class="h-4 w-4 text-[var(--primary)]" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m12 3 8 3v5c0 5-3.4 8.4-8 10-4.6-1.6-8-5-8-10V6l8-3Z" stroke="currentColor" stroke-width="1.8"/></svg>提取码仅用于本次解密</div>
          <button type="submit" :disabled="!password.trim() || loading" class="flex h-14 w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-b from-blue-500 to-blue-700 font-bold text-white shadow-lg shadow-blue-600/20 transition hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-45 disabled:hover:translate-y-0"><svg v-if="loading" class="h-5 w-5 animate-spin" viewBox="0 0 24 24" fill="none"><circle cx="12" cy="12" r="9" stroke="currentColor" stroke-opacity=".3" stroke-width="3"/><path d="M21 12a9 9 0 0 0-9-9" stroke="currentColor" stroke-width="3" stroke-linecap="round"/></svg><span>{{ loading ? '正在解密…' : '解密并查看' }}</span></button>
          <div class="flex gap-3 rounded-xl border border-amber-400/50 bg-amber-500/10 p-4 text-sm text-[var(--warning)]"><span aria-hidden="true">⚠</span><div><strong class="block">打开前请确认</strong><span>成功读取后，此链接将永久失效</span></div></div>
        </form>
      </template>

      <template v-else-if="showTransfer">
        <div class="app-surface mx-auto max-w-xl space-y-5 rounded-[24px] p-6 sm:p-9">
          <div class="mx-auto grid h-14 w-14 place-items-center rounded-full bg-blue-500/10 text-2xl text-[var(--primary)]">↔</div>
          <div><p class="text-sm font-bold text-[var(--primary)]">WebRTC 点对点传输</p><h1 class="mt-1 text-3xl font-black tracking-tight sm:text-4xl">正在等待发送方</h1><p class="mt-2 text-sm text-[var(--muted)]">附件不会经由服务器中转或暂存。</p></div>
          <p class="rounded-xl border border-blue-300/60 bg-blue-500/10 p-4 text-left text-sm text-[var(--primary)]">{{ attachmentStatus }}</p>
          <div class="flex gap-3 rounded-xl border border-amber-400/50 bg-amber-500/10 p-4 text-left text-sm text-[var(--warning)]"><span aria-hidden="true">⚠</span><div><strong class="block">双方必须同时在线</strong><span>发送方需保持创建分享的原始页面打开；接收方保持当前分享链接页面打开。若连接失败，请重新尝试连接。</span></div></div>
          <button class="h-13 w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] font-bold text-[var(--primary)] hover:bg-[var(--primary-soft)]" @click="viewSecret">重新尝试连接</button>
        </div>
      </template>

      <template v-else-if="showContent">
        <div class="app-surface mx-auto max-w-3xl space-y-5 rounded-[24px] p-6 sm:p-9">
          <div class="mx-auto grid h-14 w-14 place-items-center rounded-full bg-emerald-500/10 text-2xl text-[var(--success)]">✓</div>
          <div><p class="text-sm font-bold text-[var(--success)]">解密成功</p><h1 class="mt-1 text-3xl font-black tracking-tight sm:text-4xl">秘密内容</h1><p class="mt-2 text-sm text-[var(--muted)]">内容已读取并从服务器永久销毁</p></div>
          <div class="border-t border-[var(--line)] pt-5 text-left"><div class="mb-2 flex items-center justify-between gap-3"><span class="text-sm font-bold text-[var(--primary)]">仅此一次</span><button class="rounded-lg bg-[var(--primary)] px-4 py-2 text-sm font-bold text-white" @click="copyContent">复制内容</button></div><div class="min-h-40 whitespace-pre-wrap break-words rounded-xl border border-[var(--line)] bg-[var(--surface-soft)] p-5 font-mono leading-8">{{ content }}</div></div>
          <p v-if="attachmentStatus" class="rounded-xl border border-blue-300/60 bg-blue-500/10 p-4 text-left text-sm text-[var(--primary)]">{{ attachmentStatus }}</p>
          <div class="flex gap-3 rounded-xl border border-amber-400/50 bg-amber-500/10 p-4 text-left text-sm text-[var(--warning)]"><span aria-hidden="true">⚠</span><div><strong class="block">此秘密已销毁</strong><span>关闭页面后将无法再次查看，请妥善保存需要的信息</span></div></div>
          <button data-test="home-button" class="h-13 w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] font-bold text-[var(--primary)] hover:bg-[var(--primary-soft)]" @click="goHome">返回首页</button>
        </div>
      </template>

      <template v-else>
        <div class="app-surface mx-auto max-w-xl space-y-5 rounded-[24px] p-6 sm:p-9">
          <div class="mx-auto grid h-14 w-14 place-items-center rounded-full bg-red-500/10 text-2xl font-bold text-[var(--danger)]">!</div>
          <div><p class="text-sm font-bold text-[var(--muted)]">无法读取</p><h1 class="mt-1 text-3xl font-black tracking-tight sm:text-4xl">这个秘密已不存在</h1><p class="mt-2 text-sm text-[var(--muted)]">链接可能已被查看、过期，或输入有误</p></div>
          <div role="alert" class="rounded-xl border border-red-400/40 bg-red-500/10 p-4 text-left text-sm text-[var(--danger)]"><strong class="block">链接已失效</strong>{{ error }}</div>
          <div class="space-y-3 rounded-xl border border-[var(--line)] bg-[var(--surface-soft)] p-4 text-left text-sm text-[var(--muted)]"><p>✓ 内容已被读取并销毁</p><p>✓ 分享链接已过期</p><p>✓ 链接地址不完整</p></div>
          <button data-test="home-button" class="h-13 w-full rounded-xl bg-[var(--primary)] font-bold text-white" @click="goHome">返回首页</button>
          <p class="text-xs text-[var(--muted)]">一次性秘密在读取后无法恢复</p>
        </div>
      </template>
    </div>

    <Transition name="toast"><div v-if="copyMessage" role="status" class="fixed bottom-6 left-1/2 z-50 -translate-x-1/2 rounded-xl bg-slate-950 px-4 py-2.5 text-sm font-semibold text-white shadow-xl">{{ copyMessage }}</div></Transition>
  </div>
</template>

<style scoped>
.h-13 { height: 3.25rem; }
.toast-enter-active, .toast-leave-active { transition: opacity .2s, transform .2s; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translate(-50%, .5rem); }
</style>
