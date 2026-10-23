import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { maxAttachmentSize, receiveAttachment, startSendingAttachment } from '../attachmentTransfer'

class FakeChannel {
  bufferedAmount = 0
  onopen: (() => Promise<void> | void) | null = null
  onmessage: ((event: { data: string }) => Promise<void> | void) | null = null
  send = vi.fn()
  close = vi.fn()
}

class FakePeer {
  static instances: FakePeer[] = []
  onicecandidate: ((event: { candidate: unknown }) => void) | null = null
  channel = new FakeChannel()
  createDataChannel = vi.fn(() => this.channel)
  createOffer = vi.fn(async () => ({ type: 'offer', sdp: 'fake-offer' }))
  createAnswer = vi.fn(async () => ({ type: 'answer', sdp: 'fake-answer' }))
  iceGatheringState = 'complete'
  getStats = vi.fn(async () => new Map<string, { type: string }>())
  setLocalDescription = vi.fn(async () => undefined)
  setRemoteDescription = vi.fn(async () => undefined)
  addIceCandidate = vi.fn(async () => undefined)
  close = vi.fn()

  constructor() {
    FakePeer.instances.push(this)
  }
}

class FakeSocket {
  static instances: FakeSocket[] = []
  onmessage: ((event: { data: string }) => Promise<void> | void) | null = null
  onerror: (() => void) | null = null
  onopen: (() => void) | null = null
  onclose: (() => void) | null = null
  send = vi.fn()
  close = vi.fn()

  constructor() {
    FakeSocket.instances.push(this)
  }
}

describe('attachment transfer', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    FakePeer.instances = []
    FakeSocket.instances = []
    vi.stubGlobal('WebSocket', FakeSocket)
    vi.stubGlobal('RTCPeerConnection', FakePeer)
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
  })

  it('limits attachments to 10 MB', () => {
    expect(maxAttachmentSize).toBe(10 * 1024 * 1024)
  })

  it('reports a sender signaling disconnect instead of silently waiting', () => {
    const onError = vi.fn()
    startSendingAttachment(new File([], 'report.txt'), { key: 'key', password: 'pass' }, { onError })
    FakeSocket.instances[0].onclose?.()
    expect(onError).toHaveBeenCalledWith('信令连接已断开，请重新创建分享并保持发送页面打开')
  })

  it('reports receiver signaling disconnect immediately and clears its timeout', async () => {
    const onError = vi.fn()
    receiveAttachment({ key: 'key', password: 'pass' }, { onError })
    FakeSocket.instances[0].onclose?.()
    expect(onError).toHaveBeenCalledWith('信令连接已断开，请重新尝试连接')
    await vi.advanceTimersByTimeAsync(30_000)
    expect(onError).toHaveBeenCalledTimes(1)
  })

  it('identifies unresolved local-network candidates after the sender has answered', async () => {
    vi.stubGlobal('navigator', { userAgent: 'Mozilla/5.0 (Macintosh)' })
    const onError = vi.fn()
    const onStatus = vi.fn()
    const stop = receiveAttachment({ key: 'key', password: 'pass' }, { onError, onStatus })
    const socket = FakeSocket.instances[0]
    const peer = FakePeer.instances[0]
    peer.onicecandidate?.({ candidate: { candidate: 'candidate:1 1 udp 1 local-host.local 12345 typ host', address: 'local-host.local' } })
    await socket.onmessage?.({ data: JSON.stringify({ type: 'offer', payload: { type: 'offer', sdp: 'test' } }) })
    await socket.onmessage?.({ data: JSON.stringify({ type: 'candidate', payload: { candidate: 'candidate:2 1 udp 1 remote-host.local 12346 typ host' } }) })

    expect(onStatus).toHaveBeenCalledWith('已联系到发送方，正在建立 WebRTC 点对点连接…')
    await vi.advanceTimersByTimeAsync(30_000)
    expect(onError).toHaveBeenCalledWith('WebRTC 本地网络连接失败。请检查 macOS「系统设置 → 隐私与安全 → 本地网络」中的 Chrome 权限，然后重新连接。')
    expect(peer.close).toHaveBeenCalled()
    stop()
  })

  it('keeps the sender available until receipt is acknowledged and supports a receiver retry', async () => {
    const onError = vi.fn()
    const stop = startSendingAttachment(new File([], 'report.txt'), { key: 'key', password: 'pass' }, { onError })
    const socket = FakeSocket.instances[0]

    await socket.onmessage?.({ data: JSON.stringify({ type: 'ready' }) })
    expect(FakePeer.instances).toHaveLength(1)
    expect(FakePeer.instances[0].createOffer).toHaveBeenCalledOnce()

    await vi.advanceTimersByTimeAsync(60_000)
    expect(socket.close).not.toHaveBeenCalled()
    expect(onError).not.toHaveBeenCalled()

    await socket.onmessage?.({ data: JSON.stringify({ type: 'ready' }) })
    expect(FakePeer.instances).toHaveLength(2)
    expect(FakePeer.instances[0].close).toHaveBeenCalledOnce()
    expect(FakePeer.instances[1].createOffer).toHaveBeenCalledOnce()

    await FakePeer.instances[1].channel.onopen?.()
    await socket.onmessage?.({ data: JSON.stringify({ type: 'ready' }) })
    expect(FakePeer.instances).toHaveLength(3)

    await FakePeer.instances[2].channel.onmessage?.({ data: JSON.stringify({ kind: 'received' }) })
    expect(socket.close).toHaveBeenCalledOnce()
    await socket.onmessage?.({ data: JSON.stringify({ type: 'ready' }) })
    expect(FakePeer.instances).toHaveLength(3)

    stop()
  })
})
