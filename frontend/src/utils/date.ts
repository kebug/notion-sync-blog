import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

export const formatDate = (date: string | Date, format = 'YYYY-MM-DD HH:mm:ss') => {
  return dayjs(date).format(format)
}

export const formatRelativeTime = (date: string | Date) => {
  return dayjs(date).fromNow()
}

export const formatDuration = (seconds: number) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60

  if (hours > 0) {
    return `${hours}小时${minutes}分钟`
  } else if (minutes > 0) {
    return `${minutes}分钟${secs}秒`
  } else {
    return `${secs}秒`
  }
}

export const isToday = (date: string | Date) => {
  return dayjs(date).isSame(dayjs(), 'day')
}

export const isYesterday = (date: string | Date) => {
  return dayjs(date).isSame(dayjs().subtract(1, 'day'), 'day')
}

export const getDateRange = (days: number) => {
  const end = dayjs()
  const start = dayjs().subtract(days, 'day')
  return { start, end }
}
