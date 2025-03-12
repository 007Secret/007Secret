<template>
  <div class="min-h-[calc(100vh-8rem)] flex items-center justify-center py-16 px-8">
    <div class="max-w-4xl w-full">
      <!-- Title and Description -->
      <div class="text-center mb-8">
        <h2 class="text-3xl font-extrabold text-gray-900 dark:text-white">安全分享您的秘密信息</h2>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          一次性查看，自动销毁，保障您的信息安全
        </p>
      </div>

      <!-- Main Card -->
      <div class="bg-white dark:bg-gray-800 shadow-xl rounded-lg overflow-hidden hover:shadow-2xl transition-shadow duration-300">
        <!-- Create Secret Form -->
        <div v-if="!showResult" class="p-6">
          <div class="space-y-4">
            <label for="secret-content" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              密钥内容
            </label>
            <div class="relative">
              <textarea 
                id="secret-content"
                v-model="content"
                rows="8"
                class="w-full px-3 py-2 text-base border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white resize-none"
                placeholder="在这里输入你想要分享的内容..."
              ></textarea>
            </div>
            <button 
              @click="createSecret"
              :disabled="!content"
              class="w-full flex items-center justify-center px-4 py-3 border border-transparent text-base font-medium rounded-lg text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
              </svg>
              生成分享链接
            </button>
          </div>
        </div>

        <!-- Show Result -->
        <div v-if="showResult" class="p-6 space-y-6">
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
                  密钥创建成功
                </h3>
                <div class="mt-2 text-sm text-green-700 dark:text-green-400">
                  <p>您可以复制以下链接和密钥分享给他人</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Share Link -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              分享链接
            </label>
            <div class="flex rounded-lg shadow-sm">
              <input 
                type="text" 
                readonly 
                :value="shareLink" 
                class="flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-l-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white sm:text-sm"
              />
              <button 
                @click="copyLink" 
                class="inline-flex items-center px-4 py-2 border border-l-0 border-gray-300 dark:border-gray-600 rounded-r-lg bg-blue-600 text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path>
                </svg>
              </button>
            </div>
          </div>

          <!-- Password -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              提取码
            </label>
            <div class="flex rounded-lg shadow-sm">
              <input 
                type="text" 
                readonly 
                :value="password" 
                class="flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-l-lg border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white font-mono text-lg tracking-wider text-center"
              />
              <button 
                @click="copyPassword" 
                class="inline-flex items-center px-4 py-2 border border-l-0 border-gray-300 dark:border-gray-600 rounded-r-lg bg-blue-600 text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"></path>
                </svg>
              </button>
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
                  注意: 此链接只能被查看一次，查看后将自动销毁
                </p>
              </div>
            </div>
          </div>

          <!-- Create New Button -->
          <button 
            @click="reset" 
            class="w-full flex items-center justify-center px-4 py-3 border border-gray-300 dark:border-gray-600 shadow-sm text-base font-medium rounded-lg text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
            创建新的分享
          </button>
        </div>
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

const content = ref('')
const showResult = ref(false)
const shareLink = ref('')
const password = ref('')
const showCopyNotification = ref(false)

async function createSecret() {
  if (!content.value.trim()) return;
  
  try {
    const response = await fetch('/api/secret', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ content: content.value }),
    })

    if (!response.ok) {
      throw new Error('Failed to create secret')
    }

    const data = await response.json()
    shareLink.value = `${window.location.origin}/s/${data.key}`
    password.value = data.password
    showResult.value = true
  } catch (error) {
    alert('创建分享失败，请稍后重试')
    console.error('Error:', error)
  }
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text)
    .then(() => {
      showCopyNotification.value = true
      setTimeout(() => {
        showCopyNotification.value = false
      }, 2000)
    })
    .catch(() => alert('复制失败，请手动复制'))
}

function copyLink() {
  copyToClipboard(shareLink.value)
}

function copyPassword() {
  copyToClipboard(password.value)
}

function reset() {
  content.value = ''
  showResult.value = false
  shareLink.value = ''
  password.value = ''
}
</script>