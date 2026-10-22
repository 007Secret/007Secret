import { describe, expect, it } from 'vitest'

import { maxAttachmentSize } from '../attachmentTransfer'

describe('attachment transfer', () => {
  it('limits attachments to 10 MB', () => {
    expect(maxAttachmentSize).toBe(10 * 1024 * 1024)
  })
})
