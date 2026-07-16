import { beforeEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import App from '../../App.vue'

describe('App shell', () => {
  beforeEach(() => {
    vi.stubGlobal('matchMedia', vi.fn(() => ({ matches: false })))
    localStorage.clear()
  })

  it('renders the secure product shell without placeholder navigation', () => {
    const wrapper = mount(App, {
      global: {
        stubs: {
          RouterLink: { template: '<a><slot /></a>' },
          RouterView: { template: '<div data-test="router-view" />' },
        },
      },
    })

    expect(wrapper.get('a').text()).toContain('007Secret')
    expect(wrapper.text()).toContain('端到端加密')
    expect(wrapper.text()).toContain('阅后即焚')
    expect(wrapper.find('button[aria-label*="主题"]').exists()).toBe(true)
    expect(wrapper.text()).not.toContain('关于我们')
    expect(wrapper.text()).not.toContain('使用帮助')
  })
})
