import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

const routerPush = vi.fn()
const route = { params: { key: 'abc123' }, query: { password: '' } as Record<string, string> }

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: routerPush }),
  useRoute: () => route,
}))

import ViewSecret from '../../views/ViewSecret.vue'

describe('ViewSecret', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
    routerPush.mockReset()
    route.params.key = 'abc123'
    route.query = { password: '' }
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
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ content: 'server password: top-secret' }),
    }))
    const wrapper = mount(ViewSecret)

    await wrapper.get('input').setValue('8F4K2P')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(wrapper.text()).toContain('解密成功')
    expect(wrapper.text()).toContain('server password: top-secret')
    expect(wrapper.text()).toContain('此秘密已销毁')
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
