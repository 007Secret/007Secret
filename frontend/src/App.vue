<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div class="min-h-screen h-screen flex flex-col bg-gray-50 dark:bg-gray-900">
    <!-- Header -->
    <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
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
            <nav class="ml-8 flex space-x-8">
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
          <div class="flex items-center">
            <button class="bg-gray-200 dark:bg-gray-700 p-2 rounded-full text-gray-500 dark:text-gray-300 hover:text-gray-600 dark:hover:text-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
              <span class="sr-only">切换主题</span>
              <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col overflow-auto">
      <RouterView />
    </main>

    <!-- Footer -->
    <footer class="flex-shrink-0 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col items-center justify-center">
          <p class="text-sm text-gray-500 dark:text-gray-400">
            © {{ new Date().getFullYear() }} 007Secret. 保障您的信息安全
          </p>
          <div class="mt-2 flex space-x-6">
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
