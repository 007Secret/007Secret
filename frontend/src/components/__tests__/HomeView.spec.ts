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

  it('fills an empty secret with the selected attachment name without replacing typed text', async () => {
    const wrapper = mount(HomeView)
    const fileInput = wrapper.get('input[type="file"]')
    const secret = wrapper.get('textarea[aria-label="秘密内容"]')
    const file = new File(['report'], 'report.pdf', { type: 'application/pdf' })

    Object.defineProperty(fileInput.element, 'files', { value: [file], configurable: true })
    await fileInput.trigger('change')
    expect((secret.element as HTMLTextAreaElement).value).toBe('report.pdf')

    await secret.setValue('keep this secret')
    Object.defineProperty(fileInput.element, 'files', { value: [new File(['new'], 'new.pdf')], configurable: true })
    await fileInput.trigger('change')
    expect((secret.element as HTMLTextAreaElement).value).toBe('keep this secret')
  })

  it('rejects attachments larger than 10 MB before creating a secret', async () => {
    const wrapper = mount(HomeView)
    const oversized = new File([new Uint8Array(10 * 1024 * 1024 + 1)], 'large.bin')

    const fileInput = wrapper.get('input[type="file"]')
    Object.defineProperty(fileInput.element, 'files', { value: [oversized], configurable: true })
    await fileInput.trigger('change')

    expect(wrapper.get('[role="alert"]').text()).toContain('文件不能超过 10 MB')
    expect(wrapper.get('button[type="submit"]').attributes('disabled')).toBeDefined()
  })

  it('includes selected attachment metadata in the create request', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ key: 'abc123', password: '8F4K2P' }),
    }))
    const wrapper = mount(HomeView)
    const file = new File(['report'], 'report.pdf', { type: 'application/pdf' })

    const fileInput = wrapper.get('input[type="file"]')
    Object.defineProperty(fileInput.element, 'files', { value: [file], configurable: true })
    await fileInput.trigger('change')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    const [, options] = vi.mocked(fetch).mock.calls[0]
    expect(JSON.parse(String(options?.body))).toEqual({
      content: 'report.pdf',
      attachment: { name: 'report.pdf', size: file.size, type: 'application/pdf' },
    })
  })
})
