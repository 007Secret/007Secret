export function formatFileSize(size: number) {
  if (size < 1024) return `${size} B`
  const unit = size < 1024 * 1024 ? 'KB' : 'MB'
  const value = size / (unit === 'KB' ? 1024 : 1024 * 1024)
  return `${new Intl.NumberFormat('zh-CN', { maximumFractionDigits: 2 }).format(value)} ${unit}`
}
