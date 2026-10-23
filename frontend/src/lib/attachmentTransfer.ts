export const maxAttachmentSize = 10 * 1024 * 1024

export type TransferCredentials = { key: string; password: string }
export type TransferCallbacks = {
  onStatus?: (message: string) => void
  onProgress?: (transferred: number, total: number) => void
  onComplete?: (file: File) => void
  onError?: (message: string) => void
}

const chunkSize = 64 * 1024
const connectionTimeout = 30_000
const offlinePeerMessage = '无法建立 WebRTC 点对点连接。请确认发送方与接收方同时在线并保持页面打开，然后重试。'

function signalingURL(credentials: TransferCredentials, role: 'sender' | 'receiver') {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const query = new URLSearchParams({ ...credentials, role })
  return `${protocol}//${window.location.host}/api/transfer/signal?${query}`
}

function sendSignal(socket: WebSocket, type: 'offer' | 'answer' | 'candidate' | 'ready', payload?: unknown) {
  socket.send(JSON.stringify({ type, payload }))
}

function setupIce(peer: RTCPeerConnection, socket: WebSocket, onCandidate?: (candidate: RTCIceCandidate | null) => void) {
  peer.onicecandidate = (event) => {
    onCandidate?.(event.candidate)
    if (event.candidate) sendSignal(socket, 'candidate', event.candidate)
  }
}

export function startSendingAttachment(file: File, credentials: TransferCredentials, callbacks: TransferCallbacks = {}) {
  if (typeof WebSocket === 'undefined' || typeof RTCPeerConnection === 'undefined') {
    callbacks.onError?.('当前浏览器不支持 WebRTC 点对点传输')
    return () => undefined
  }
  const socket = new WebSocket(signalingURL(credentials, 'sender'))
  let closed = false
  let peer: RTCPeerConnection | undefined
  let channel: RTCDataChannel | undefined
  let attempt = 0

  const close = () => {
    if (closed) return
    closed = true
    channel?.close()
    peer?.close()
    socket.close()
  }
  const fail = (message: string) => {
    if (closed) return
    callbacks.onError?.(message)
    close()
  }

  const beginAttempt = async () => {
    if (closed) return
    attempt += 1
    const currentAttempt = attempt
    channel?.close()
    peer?.close()
    const currentPeer = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] })
    const currentChannel = currentPeer.createDataChannel('file-transfer', { ordered: true })
    peer = currentPeer
    channel = currentChannel
    setupIce(currentPeer, socket)
    currentChannel.onmessage = ({ data }) => {
      if (currentAttempt !== attempt || typeof data !== 'string') return
      try {
        if (JSON.parse(data).kind === 'received') close()
      } catch {
        // Ignore malformed control data from the peer.
      }
    }
    currentChannel.onopen = async () => {
      if (closed || currentAttempt !== attempt) return
      try {
        currentChannel.send(JSON.stringify({ kind: 'metadata', name: file.name, type: file.type, size: file.size }))
        let transferred = 0
        for (let offset = 0; offset < file.size; offset += chunkSize) {
          if (currentAttempt !== attempt) return
          while (currentChannel.bufferedAmount > chunkSize * 8) await new Promise(resolve => setTimeout(resolve, 20))
          if (currentAttempt !== attempt) return
          currentChannel.send(await file.slice(offset, offset + chunkSize).arrayBuffer())
          transferred += Math.min(chunkSize, file.size - offset)
          callbacks.onProgress?.(transferred, file.size)
        }
        currentChannel.send(JSON.stringify({ kind: 'complete' }))
      } catch {
        if (currentAttempt === attempt) callbacks.onError?.('附件发送失败，请让接收方重新尝试连接')
      }
    }
    const offer = await currentPeer.createOffer()
    if (closed || currentAttempt !== attempt) return
    await currentPeer.setLocalDescription(offer)
    if (closed || currentAttempt !== attempt) return
    sendSignal(socket, 'offer', offer)
  }
  socket.onmessage = async ({ data }) => {
    try {
      const signal = JSON.parse(String(data))
      if (signal.type === 'ready') {
        callbacks.onStatus?.('接收方已上线，正在建立 WebRTC 点对点连接…')
        await beginAttempt()
      }
      if (signal.type === 'answer' && peer) await peer.setRemoteDescription(signal.payload)
      if (signal.type === 'candidate' && peer) await peer.addIceCandidate(signal.payload)
    } catch {
      fail('无法处理点对点连接信令')
    }
  }
  socket.onerror = () => fail('无法连接信令服务')
  socket.onclose = () => {
    if (channel?.readyState !== 'open') fail('信令连接已断开，请重新创建分享并保持发送页面打开')
  }

  return close
}

export function receiveAttachment(credentials: TransferCredentials, callbacks: TransferCallbacks = {}) {
  if (typeof WebSocket === 'undefined' || typeof RTCPeerConnection === 'undefined') {
    callbacks.onError?.('当前浏览器不支持 WebRTC 点对点传输')
    return () => undefined
  }
  const socket = new WebSocket(signalingURL(credentials, 'receiver'))
  const peer = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] })
  let received = 0
  let metadata: { name: string; type: string; size: number } | null = null
  const chunks: BlobPart[] = []
  let closed = false
  let connected = false
  let receivedOffer = false
  let hasLocalMdnsCandidate = false
  let hasRemoteMdnsCandidate = false
  let timeout: number | undefined

  const close = () => {
    if (closed) return
    closed = true
    if (timeout) window.clearTimeout(timeout)
    peer.close()
    socket.close()
  }
  const fail = (message: string) => {
    if (closed) return
    callbacks.onError?.(message)
    close()
  }

  timeout = window.setTimeout(async () => {
    if (closed || connected) return
    let message = receivedOffer
      ? '双方已在线，但未能建立 WebRTC 直连。请检查浏览器的本地网络权限、代理或网络限制，然后重试。'
      : offlinePeerMessage
    if (receivedOffer && hasLocalMdnsCandidate && hasRemoteMdnsCandidate && peer.iceGatheringState === 'complete') {
      try {
        const stats = await peer.getStats()
        const hasRemoteAddress = Array.from(stats.values()).some(report => report.type === 'remote-candidate')
        if (!hasRemoteAddress && navigator.userAgent.includes('Macintosh')) {
          message = 'WebRTC 本地网络连接失败。请检查 macOS「系统设置 → 隐私与安全 → 本地网络」中的 Chrome 权限，然后重新连接。'
        }
      } catch {
        // Keep the general connection error if the browser cannot provide stats.
      }
    }
    if (!closed && !connected) fail(message)
  }, connectionTimeout)
  setupIce(peer, socket, candidate => {
    if (candidate?.address?.endsWith('.local')) hasLocalMdnsCandidate = true
  })
  socket.onmessage = async ({ data }) => {
    try {
      const signal = JSON.parse(String(data))
      if (signal.type === 'offer') {
        receivedOffer = true
        callbacks.onStatus?.('已联系到发送方，正在建立 WebRTC 点对点连接…')
        await peer.setRemoteDescription(signal.payload)
        const answer = await peer.createAnswer()
        await peer.setLocalDescription(answer)
        sendSignal(socket, 'answer', answer)
      }
      if (signal.type === 'candidate') {
        if (signal.payload?.candidate?.includes('.local')) hasRemoteMdnsCandidate = true
        await peer.addIceCandidate(signal.payload)
      }
    } catch {
      fail('无法建立点对点连接')
    }
  }
  socket.onerror = () => fail('无法连接信令服务')
  socket.onclose = () => {
    if (!connected) fail('信令连接已断开，请重新尝试连接')
  }
  socket.onopen = () => sendSignal(socket, 'ready')
  peer.ondatachannel = ({ channel }) => {
    connected = true
    if (timeout) window.clearTimeout(timeout)
    channel.onmessage = async ({ data }) => {
      if (typeof data === 'string') {
        const control = JSON.parse(data)
        if (control.kind === 'metadata') {
          if (!Number.isInteger(control.size) || control.size < 0 || control.size > maxAttachmentSize) {
            fail('附件大小无效')
            return
          }
          metadata = control
        }
        if (control.kind === 'complete' && metadata && received === metadata.size) {
          channel.send(JSON.stringify({ kind: 'received' }))
          callbacks.onComplete?.(new File([new Blob(chunks, { type: metadata.type })], metadata.name, { type: metadata.type }))
          window.setTimeout(close, 100)
        }
        return
      }
      const chunk = data instanceof Blob ? data : new Blob([data])
      chunks.push(chunk)
      received += chunk.size
      if (metadata) callbacks.onProgress?.(received, metadata.size)
    }
  }
  return close
}
