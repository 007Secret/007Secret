export const maxAttachmentSize = 10 * 1024 * 1024

export type TransferCredentials = { key: string; password: string }
export type TransferCallbacks = {
  onProgress?: (transferred: number, total: number) => void
  onComplete?: (file: File) => void
  onError?: (message: string) => void
}

const chunkSize = 64 * 1024

function signalingURL(credentials: TransferCredentials, role: 'sender' | 'receiver') {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const query = new URLSearchParams({ ...credentials, role })
  return `${protocol}//${window.location.host}/api/transfer/signal?${query}`
}

function sendSignal(socket: WebSocket, type: 'offer' | 'answer' | 'candidate', payload: unknown) {
  socket.send(JSON.stringify({ type, payload }))
}

function setupIce(peer: RTCPeerConnection, socket: WebSocket) {
  peer.onicecandidate = (event) => {
    if (event.candidate) sendSignal(socket, 'candidate', event.candidate)
  }
}

export function startSendingAttachment(file: File, credentials: TransferCredentials, callbacks: TransferCallbacks = {}) {
	if (typeof WebSocket === 'undefined' || typeof RTCPeerConnection === 'undefined') {
		callbacks.onError?.('当前浏览器不支持 WebRTC 点对点传输')
		return () => undefined
	}
  const socket = new WebSocket(signalingURL(credentials, 'sender'))
  const peer = new RTCPeerConnection({ iceServers: [{ urls: 'stun:stun.l.google.com:19302' }] })
  const channel = peer.createDataChannel('file-transfer', { ordered: true })
  let closed = false

  const close = () => {
    if (closed) return
    closed = true
    channel.close()
    peer.close()
    socket.close()
  }

  setupIce(peer, socket)
  socket.onmessage = async ({ data }) => {
    try {
      const signal = JSON.parse(String(data))
      if (signal.type === 'answer') await peer.setRemoteDescription(signal.payload)
      if (signal.type === 'candidate') await peer.addIceCandidate(signal.payload)
    } catch {
      callbacks.onError?.('无法处理点对点连接信令')
      close()
    }
  }
  socket.onerror = () => callbacks.onError?.('无法连接信令服务')
  socket.onopen = async () => {
    try {
      const offer = await peer.createOffer()
      await peer.setLocalDescription(offer)
      sendSignal(socket, 'offer', offer)
    } catch {
      callbacks.onError?.('无法创建点对点连接')
      close()
    }
  }
  channel.onopen = async () => {
    try {
      channel.send(JSON.stringify({ kind: 'metadata', name: file.name, type: file.type, size: file.size }))
      let transferred = 0
      for (let offset = 0; offset < file.size; offset += chunkSize) {
        while (channel.bufferedAmount > chunkSize * 8) await new Promise(resolve => setTimeout(resolve, 20))
        channel.send(await file.slice(offset, offset + chunkSize).arrayBuffer())
        transferred += Math.min(chunkSize, file.size - offset)
        callbacks.onProgress?.(transferred, file.size)
      }
      channel.send(JSON.stringify({ kind: 'complete' }))
    } catch {
      callbacks.onError?.('附件发送失败')
      close()
    }
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

  const close = () => { peer.close(); socket.close() }
  setupIce(peer, socket)
  socket.onmessage = async ({ data }) => {
    try {
      const signal = JSON.parse(String(data))
      if (signal.type === 'offer') {
        await peer.setRemoteDescription(signal.payload)
        const answer = await peer.createAnswer()
        await peer.setLocalDescription(answer)
        sendSignal(socket, 'answer', answer)
      }
      if (signal.type === 'candidate') await peer.addIceCandidate(signal.payload)
    } catch {
      callbacks.onError?.('无法建立点对点连接')
      close()
    }
  }
  socket.onerror = () => callbacks.onError?.('无法连接信令服务')
  peer.ondatachannel = ({ channel }) => {
    channel.onmessage = async ({ data }) => {
      if (typeof data === 'string') {
        const control = JSON.parse(data)
        if (control.kind === 'metadata') metadata = control
        if (control.kind === 'complete' && metadata && received === metadata.size) {
          callbacks.onComplete?.(new File([new Blob(chunks, { type: metadata.type })], metadata.name, { type: metadata.type }))
          close()
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
