<script setup>
import { RouterView } from 'vue-router'
import { ref, onMounted, watch } from 'vue'

// 创建一个响应式变量来跟踪暗黑模式状态
const isDarkMode = ref(false)

// 在组件挂载时检查本地存储或系统偏好
onMounted(() => {
  // 检查本地存储中的偏好
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDarkMode.value = savedTheme === 'dark'
  } else {
    // 如果没有保存的偏好，则检查系统偏好
    isDarkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  applyTheme()
})

// 监听isDarkMode的变化并应用主题
watch(isDarkMode, () => {
  applyTheme()
  // 保存用户偏好到本地存储
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light')
})

// 应用主题到HTML元素
const applyTheme = () => {
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark')
    document.body.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
    document.body.classList.remove('dark')
  }
}

// 切换暗黑模式
const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value
}
</script>

<template>
  <div class="w-full h-screen flex flex-col bg-gray-50 dark:bg-gray-900 mx-auto overflow-hidden">
    <!-- Header -->
    <header class="w-full bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="flex-shrink-0 flex items-center">
              <div class="w-8 h-8 bg-blue-600 rounded-md flex items-center justify-center">
                <svg class="h-5 w-5 text-white" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 15V17M6 9V7C6 4.79086 7.79086 3 10 3H14C16.2091 3 18 4.79086 18 7V9M6 9C3.79086 9 2 10.7909 2 13V17C2 19.2091 3.79086 21 6 21H18C20.2091 21 22 19.2091 22 17V13C22 10.7909 20.2091 9 18 9M6 9H18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
              <span class="ml-2 text-xl font-bold text-gray-900 dark:text-white">007Secret</span>
            </router-link>
            <nav class="hidden md:ml-8 md:flex md:space-x-8">
              <router-link to="/" class="border-blue-500 text-gray-900 dark:text-white inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                首页
              </router-link>
              <a href="#" class="border-transparent text-gray-500 dark:text-gray-300 hover:border-gray-300 hover:text-gray-700 dark:hover:text-gray-200 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                关于我们
              </a>
              <a href="#" class="border-transparent text-gray-500 dark:text-gray-300 hover:border-gray-300 hover:text-gray-700 dark:hover:text-gray-200 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium">
                使用帮助
              </a>
            </nav>
          </div>
          <div class="flex items-center space-x-4">
            <button @click="toggleDarkMode" class="bg-gray-200 dark:bg-gray-700 p-2 rounded-full text-gray-500 dark:text-gray-300 hover:text-gray-600 dark:hover:text-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
              <span class="sr-only">切换主题</span>
              <svg v-if="!isDarkMode" class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
              <svg v-else class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </button>
            <!-- 移动端菜单按钮 -->
            <button type="button" class="md:hidden inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500">
              <span class="sr-only">打开菜单</span>
              <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col justify-start overflow-hidden w-full">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8 pt-1 pb-7">
        <RouterView />
      </div>
    </main>

    <!-- Footer -->
    <footer class="w-full flex-shrink-0 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex flex-col items-center">
          <p class="text-sm text-gray-500 dark:text-gray-400 text-center mb-2">
            © {{ new Date().getFullYear() }} 007Secret. 保障您的信息安全
          </p>
          <div class="flex justify-center space-x-6">
            <a href="#" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
              <span class="sr-only">隐私政策</span>
              隐私政策
            </a>
            <a href="#" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
              <span class="sr-only">使用条款</span>
              使用条款
            </a>
            <a href="#" class="text-gray-400 hover:text-gray-500 dark:hover:text-gray-300">
              <span class="sr-only">联系我们</span>
              联系我们
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<style>
@import '@/assets/tailwind.css';

html, body {
  width: 100%;
  max-width: 100%;
  height: 100%;
  overflow: hidden;
  margin: 0;
  padding: 0;
  position: fixed;
}

#app {
  width: 100%;
  min-width: 100%;
  height: 100%;
  overflow: hidden;
  display: flex;
  justify-content: center;
}

/* 禁用所有滚动 */
main, .container, .min-h-screen {
  overflow: hidden;
}

/* 使用Tailwind的container类的默认样式 */
.container {
  width: 100%;
  margin-left: auto;
  margin-right: auto;
}

/* 自定义动画 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}

/* 确保文本区域在暗模式下有正确的颜色 */
textarea::placeholder, input::placeholder {
  color: #9ca3af;
}
.dark textarea::placeholder, .dark input::placeholder {
  color: #6b7280;
}

/* 平滑过渡效果 */
.transition-colors {
  transition-property: background-color, border-color, color, fill, stroke;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}
</style>
