export const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

export const formatNumber = (num: number) => {
  return new Intl.NumberFormat().format(num)
}

export const truncate = (str: string, length: number) => {
  if (str.length <= length) return str
  return str.slice(0, length) + '...'
}

export const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

export const camelToKebab = (str: string) => {
  return str.replace(/[A-Z]/g, (letter) => `-${letter.toLowerCase()}`)
}

export const kebabToCamel = (str: string) => {
  return str.replace(/-([a-z])/g, (g) => g[1].toUpperCase())
}

export const getInitials = (name: string) => {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

export const generateId = () => {
  return Math.random().toString(36).substr(2, 9)
}
