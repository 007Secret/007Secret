import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

const routerPush = vi.fn()
const route = { params: { key: 'abc123' }, query: { password: '' } as Record<string, string> }

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: routerPush }),
  useRoute: () => route,
}))

vi.mock('../../lib/attachmentTransfer', () => ({ receiveAttachment: vi.fn() }))

import { receiveAttachment, type TransferCallbacks } from '../../lib/attachmentTransfer'
import ViewSecret from '../../views/ViewSecret.vue'

describe('ViewSecret', () => {
  const stopReceiving = vi.fn()
  const attachment = { name: 'report.pdf', size: 3, type: 'application/pdf' }
  let callbacks: TransferCallbacks

  beforeEach(() => {
    vi.restoreAllMocks()
    routerPush.mockReset()
    route.params.key = 'abc123'
    route.query = { password: '' }
    stopReceiving.mockReset()
    vi.mocked(receiveAttachment).mockReset().mockImplementation((_credentials, handlers = {}) => {
      callbacks = handlers
      return stopReceiving
    })
    vi.stubGlobal('URL', class extends URL {
      static createObjectURL = vi.fn(() => 'blob:received-report')
      static revokeObjectURL = vi.fn()
    })
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('prefills a password from the link and presents the one-time flow', async () => {
    route.query = { password: '8F4K2P' }
    const wrapper = mount(ViewSecret)
    await flushPromises()

    expect(wrapper.get('input[aria-label="提取码"]').element).toHaveProperty('value', '8F4K2P')
    expect(wrapper.text()).toContain('输入提取码后，秘密将被读取并立即销毁')
    expect(wrapper.text()).toContain('打开前请确认')
  })

  it('renders decrypted content and the destruction notice after success', async () => {
    vi.stubGlobal('fetch', vi.fn()
      .mockResolvedValueOnce({ ok: true, json: async () => ({ content: 'server password: top-secret' }) }))
    const wrapper = mount(ViewSecret)

    await wrapper.get('input').setValue('8F4K2P')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(wrapper.text()).toContain('解密成功')
    expect(wrapper.text()).toContain('server password: top-secret')
    expect(wrapper.text()).toContain('此秘密已销毁')
  })

  it('shows text immediately while an attachment is still connecting', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ content: 'text available without WebRTC', attachment }),
    })
    vi.stubGlobal('fetch', fetchMock)
    const wrapper = mount(ViewSecret)

    await wrapper.get('input').setValue('8F4K2P')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(wrapper.text()).toContain('text available without WebRTC')
    expect(wrapper.text()).toContain('双方必须同时在线')
    expect(wrapper.text()).toContain('附件不会经由服务器中转或暂存')
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(String(fetchMock.mock.calls[0][0])).toContain('/api/secret/get?')
    expect(wrapper.find('a[download]').exists()).toBe(false)
  })

  it('keeps text visible after a failed WebRTC connection and retries only the attachment', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ content: 'keep this text', attachment }),
    })
    vi.stubGlobal('fetch', fetchMock)
    const wrapper = mount(ViewSecret)
    await wrapper.get('input').setValue('8F4K2P')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    const previousCallbacks = callbacks
    callbacks.onError?.('发送方不在线')
    await flushPromises()

    expect(wrapper.text()).toContain('keep this text')
    expect(wrapper.text()).toContain('WebRTC 未连接')
    expect(wrapper.text()).toContain('发送方不在线')
    expect(wrapper.find('a[download]').exists()).toBe(false)
    const retry = wrapper.findAll('button').find(button => button.text() === '重新尝试连接')!
    await retry.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('keep this text')
    expect(wrapper.text()).toContain('正在等待发送方')
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(stopReceiving).toHaveBeenCalledTimes(1)
    previousCallbacks.onComplete?.(new File(['old'], 'old.pdf'))
    await flushPromises()
    expect(wrapper.find('a[download]').exists()).toBe(false)
    expect(URL.createObjectURL).not.toHaveBeenCalled()
  })

  it('offers a manual download after receiving a file without automatically clicking a link', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ content: 'file notes', attachment }),
    }))
    const click = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => undefined)
    const wrapper = mount(ViewSecret)
    await wrapper.get('input').setValue('8F4K2P')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    callbacks.onComplete?.(new File(['pdf'], 'report.pdf', { type: 'application/pdf' }))
    await flushPromises()

    expect(wrapper.text()).toContain('file notes')
    expect(click).not.toHaveBeenCalled()
    const download = wrapper.get('a[download]')
    expect(download.text()).toBe('下载文件')
    expect(download.attributes('href')).toBe('blob:received-report')
    expect(download.attributes('download')).toBe('report.pdf')
    expect(URL.revokeObjectURL).not.toHaveBeenCalled()
    expect(fetch).toHaveBeenCalledTimes(1)

    wrapper.unmount()
    expect(URL.revokeObjectURL).toHaveBeenCalledWith('blob:received-report')
    expect(stopReceiving).toHaveBeenCalledTimes(1)
    callbacks.onComplete?.(new File(['late'], 'late.pdf'))
    expect(URL.createObjectURL).toHaveBeenCalledTimes(1)
  })

  it('preserves retrieved text if the browser cannot start a WebRTC receiver', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ content: 'still readable', attachment }),
    }))
    vi.mocked(receiveAttachment).mockImplementation(() => { throw new Error('unavailable') })
    const wrapper = mount(ViewSecret)
    await wrapper.get('input').setValue('8F4K2P')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(wrapper.text()).toContain('still readable')
    expect(wrapper.text()).toContain('WebRTC 未连接')
  })

  it('renders a consumed-link error and returns home', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 404 }))
    const wrapper = mount(ViewSecret)

    await wrapper.get('input').setValue('8F4K2P')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(wrapper.text()).toContain('这个秘密已不存在')
    expect(wrapper.text()).toContain('密钥不存在或已被查看')
    await wrapper.get('button[data-test="home-button"]').trigger('click')
    expect(routerPush).toHaveBeenCalledWith('/')
  })
})
