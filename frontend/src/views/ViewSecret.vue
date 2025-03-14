<template>
  <div class="min-h-[calc(100vh-8rem)] flex flex-col items-center justify-center w-full">
    <!-- Title and Description -->
    <div class="text-center mb-8 w-full">
      <h2 class="text-3xl sm:text-4xl font-extrabold text-gray-900 dark:text-white">查看密钥内容</h2>
      <p class="mt-2 text-sm sm:text-base text-gray-600 dark:text-gray-400">
        请输入提取码查看密钥内容，注意内容只能查看一次
      </p>
    </div>

    <!-- Main Card -->
    <div class="bg-white dark:bg-gray-800 shadow-xl rounded-lg overflow-hidden hover:shadow-2xl transition-shadow duration-300 max-w-3xl w-full">
      <!-- Password Input Form -->
      <div v-if="!showContent && !error" class="p-6 sm:p-8">
        <div class="space-y-4">
          <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
            提取码
          </label>
          <div class="relative">
            <input 
              type="text"
              id="password"
              v-model="password"
              @keyup.enter="viewSecret"
              class="w-full px-3 py-2 text-base border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white font-mono text-lg tracking-wider text-center"
              placeholder="请输入提取码"
            />
          </div>
          <button 
            @click="viewSecret"
            :disabled="!password || loading"
            class="w-full flex items-center justify-center px-4 py-3 border border-transparent text-base font-medium rounded-lg text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <template v-if="!loading">
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
              </svg>
              查看内容
            </template>
            <template v-else>
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              正在解密...
            </template>
          </button>
        </div>
      </div>

      <!-- Error State -->
      <div v-if="error" class="p-6 space-y-6">
        <div class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800 dark:text-red-300">
                无法查看密钥
              </h3>
              <div class="mt-2 text-sm text-red-700 dark:text-red-400">
                <p>{{ error }}</p>
              </div>
            </div>
          </div>
        </div>

        <button 
          @click="goHome" 
          class="w-full flex items-center justify-center px-4 py-3 border border-gray-300 dark:border-gray-600 shadow-sm text-base font-medium rounded-lg text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
          </svg>
          返回首页
        </button>
      </div>

      <!-- Show Content -->
      <div v-if="showContent" class="p-6 space-y-6">
        <!-- Success Message -->
        <div class="rounded-md bg-green-50 dark:bg-green-900/20 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-green-800 dark:text-green-300">
                解密成功
              </h3>
              <div class="mt-2 text-sm text-green-700 dark:text-green-400">
                <p>此内容已被解密，并将在您关闭页面后永久销毁</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Content Box -->
        <div class="bg-gray-50 dark:bg-gray-900/50 rounded-lg p-4">
          <div class="flex justify-between items-center mb-2">
            <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300">密钥内容</h3>
            <button 
              @click="copyContent" 
              class="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 rounded-md text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path>
              </svg>
              复制内容
            </button>
          </div>
          <div class="min-h-[10rem] whitespace-pre-wrap break-words text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-800 p-4 rounded-md border border-gray-200 dark:border-gray-700">
            {{ content }}
          </div>
        </div>

        <!-- Warning -->
        <div class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm text-red-700 dark:text-red-400">
                注意: 此内容已被销毁，无法再次查看
              </p>
            </div>
          </div>
        </div>

        <!-- Back Button -->
        <button 
          @click="goHome" 
          class="w-full flex items-center justify-center px-4 py-3 border border-gray-300 dark:border-gray-600 shadow-sm text-base font-medium rounded-lg text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
          </svg>
          返回首页
        </button>
      </div>
    </div>

    <!-- Copy Notification -->
    <div 
      v-if="showCopyNotification" 
      class="fixed bottom-4 right-4 bg-gray-800 text-white px-4 py-2 rounded-lg shadow-lg transition-all duration-300 transform"
      :class="{'translate-y-0 opacity-100': showCopyNotification, 'translate-y-4 opacity-0': !showCopyNotification}"
    >
      已复制到剪贴板
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()
const password = ref('')
const content = ref('')
const showContent = ref(false)
const error = ref('')
const loading = ref(false)
const showCopyNotification = ref(false)

const viewSecret = async () => {
  if (!password.value || loading.value) return
  
  const key = route.params.key as string
  if (!key) {
    error.value = '无效的密钥链接'
    return
  }

  loading.value = true
  error.value = ''
  
  try {
    const response = await fetch(`/api/secret/get?key=${key}&password=${encodeURIComponent(password.value)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('密钥不存在或已被查看')
      } else if (response.status === 401) {
        throw new Error('提取码错误')
      } else {
        throw new Error('服务器错误，请稍后重试')
      }
    }

    const data = await response.json()
    content.value = data.content
    showContent.value = true
  } catch (err) {
    if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = '未知错误，请稍后重试'
    }
  } finally {
    loading.value = false
  }
}

const copyContent = () => {
  if (!content.value) return
  
  navigator.clipboard.writeText(content.value)
    .then(() => {
      showCopyNotification.value = true
      setTimeout(() => {
        showCopyNotification.value = false
      }, 2000)
    })
    .catch(() => alert('复制失败，请手动复制'))
}

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
/* 可以移除这些样式，因为我们已经使用了Tailwind类 */
</style>