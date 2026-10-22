import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { maxAttachmentSize, startSendingAttachment } from '../attachmentTransfer'

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
