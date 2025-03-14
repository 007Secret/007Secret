<template>
  <div class="about-container flex flex-col items-center justify-center w-full">
    <div class="card max-w-3xl w-full">
      <div class="card-header">
        <h1 class="card-title">关于 007Secret</h1>
        <p class="card-description">安全、简单的一次性信息分享服务</p>
      </div>
      
      <div class="card-body">
        <section class="about-section">
          <h2 class="section-title">我们的使命</h2>
          <p class="section-content">
            007Secret 致力于提供一个安全、简单的方式来分享敏感信息。在当今数字时代，信息安全变得越来越重要，
            我们的平台确保您的敏感信息只能被查看一次，并在查看后立即销毁，不留下任何痕迹。
          </p>
        </section>
        
        <section class="about-section">
          <h2 class="section-title">如何工作</h2>
          <div class="steps-container">
            <div class="step">
              <div class="step-number">1</div>
              <div class="step-content">
                <h3 class="step-title">创建密钥</h3>
                <p class="step-description">
                  在首页输入您想要分享的敏感信息，点击"生成分享链接"按钮。
                </p>
              </div>
            </div>
            
            <div class="step">
              <div class="step-number">2</div>
              <div class="step-content">
                <h3 class="step-title">分享链接</h3>
                <p class="step-description">
                  系统会生成一个唯一的链接和提取码。将这两者分别通过不同的渠道分享给接收者。
                </p>
              </div>
            </div>
            
            <div class="step">
              <div class="step-number">3</div>
              <div class="step-content">
                <h3 class="step-title">一次性查看</h3>
                <p class="step-description">
                  接收者使用链接和提取码查看信息。一旦信息被查看，它将立即从我们的服务器上删除。
                </p>
              </div>
            </div>
          </div>
        </section>
        
        <section class="about-section">
          <h2 class="section-title">安全保障</h2>
          <ul class="security-features">
            <li class="feature-item">
              <svg class="feature-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
              </svg>
              <span>端到端加密传输</span>
            </li>
            <li class="feature-item">
              <svg class="feature-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
              </svg>
              <span>一次性查看后自动销毁</span>
            </li>
            <li class="feature-item">
              <svg class="feature-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span>信息存储时间限制</span>
            </li>
            <li class="feature-item">
              <svg class="feature-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
              <span>不保存访问日志</span>
            </li>
          </ul>
        </section>
        
        <section class="about-section">
          <h2 class="section-title">适用场景</h2>
          <div class="use-cases">
            <div class="use-case">
              <h3 class="use-case-title">密码分享</h3>
              <p class="use-case-description">
                安全地分享账户密码给团队成员或家人，而不必担心密码被第三方截获。
              </p>
            </div>
            
            <div class="use-case">
              <h3 class="use-case-title">敏感文档</h3>
              <p class="use-case-description">
                分享包含个人信息或商业机密的文档内容，确保信息不会被泄露。
              </p>
            </div>
            
            <div class="use-case">
              <h3 class="use-case-title">私密消息</h3>
              <p class="use-case-description">
                发送私密消息给他人，确保消息在被阅读后不会留下任何痕迹。
              </p>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 无需额外的脚本逻辑
</script>

<style scoped>
:root {
  --primary-color: #3b82f6;
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --text-primary: #1f2937;
  --text-secondary: #4b5563;
  --border-color: #e5e7eb;
  --radius: 0.5rem;
}

.dark {
  --primary-color: #3b82f6;
  --bg-primary: #1f2937;
  --bg-secondary: #111827;
  --text-primary: #f9fafb;
  --text-secondary: #d1d5db;
  --border-color: #374151;
}

.about-container {
  padding: 2rem 0;
  width: 100%;
}

.card {
  background-color: var(--bg-primary);
  border-radius: var(--radius);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.card-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  text-align: center;
}

.card-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

@media (min-width: 640px) {
  .card-title {
    font-size: 1.875rem;
  }
}

.card-description {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

.card-body {
  padding: 1.5rem;
}

@media (min-width: 640px) {
  .card-body {
    padding: 2rem;
  }
}

.about-section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 1rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border-color);
}

@media (min-width: 640px) {
  .section-title {
    font-size: 1.5rem;
  }
}

.section-content {
  font-size: 0.875rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
}

@media (min-width: 640px) {
  .section-content {
    font-size: 1rem;
  }
}

.steps-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

@media (min-width: 768px) {
  .steps-container {
    flex-direction: row;
  }
}

.step {
  display: flex;
  gap: 1rem;
  align-items: flex-start;
  flex: 1;
}

.step-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  background-color: var(--primary-color);
  color: white;
  border-radius: 9999px;
  font-weight: 600;
  flex-shrink: 0;
}

.step-content {
  flex: 1;
}

.step-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.25rem 0;
}

@media (min-width: 640px) {
  .step-title {
    font-size: 1.125rem;
  }
}

.step-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

@media (min-width: 640px) {
  .step-description {
    font-size: 1rem;
  }
}

.security-features {
  list-style: none;
  padding: 0;
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

@media (min-width: 640px) {
  .security-features {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .security-features {
    grid-template-columns: repeat(4, 1fr);
  }
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

@media (min-width: 640px) {
  .feature-item {
    font-size: 1rem;
  }
}

.feature-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: var(--primary-color);
}

.use-cases {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

@media (min-width: 640px) {
  .use-cases {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .use-cases {
    grid-template-columns: repeat(3, 1fr);
  }
}

.use-case {
  background-color: var(--bg-secondary);
  border-radius: var(--radius);
  padding: 1.25rem;
}

.use-case-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

@media (min-width: 640px) {
  .use-case-title {
    font-size: 1.125rem;
  }
}

.use-case-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
}

@media (min-width: 640px) {
  .use-case-description {
    font-size: 1rem;
  }
}
</style> 