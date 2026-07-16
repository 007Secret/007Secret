import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import HomeView from '../../views/HomeView.vue'

describe('HomeView', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('shows a secure creation form with live character count', async () => {
    const wrapper = mount(HomeView)
    const textarea = wrapper.get('textarea[aria-label="秘密内容"]')

    expect(wrapper.text()).toContain('让秘密，只被看见一次')
    expect(wrapper.text()).toContain('内容将被加密存储，我们无法读取')
    expect(wrapper.text()).toContain('0 / 10,000')
    expect(wrapper.get('button[type="submit"]').attributes('disabled')).toBeDefined()

    await textarea.setValue('hello')
    expect(wrapper.text()).toContain('5 / 10,000')
    expect(wrapper.get('button[type="submit"]').attributes('disabled')).toBeUndefined()
  })

  it('renders the generated link and extraction code after success', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ key: 'abc123', password: '8F4K2P' }),
    }))
    const wrapper = mount(HomeView)

    await wrapper.get('textarea').setValue('top secret')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(fetch).toHaveBeenCalledWith('/api/secret', expect.objectContaining({ method: 'POST' }))
    expect(wrapper.text()).toContain('安全链接已生成')
    expect(wrapper.get('#extract-code').element).toHaveProperty('value', '8F4K2P')
    expect(wrapper.text()).toContain('仅可查看一次')
    expect(wrapper.text()).toContain('创建新的分享')
  })

  it('shows a request failure inline', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false }))
    const wrapper = mount(HomeView)

    await wrapper.get('textarea').setValue('top secret')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(wrapper.get('[role="alert"]').text()).toContain('创建分享失败，请稍后重试')
  })
})
