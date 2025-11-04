import { ref } from 'vue'
import { useMessage } from 'naive-ui'

export const useError = () => {
  const message = useMessage()
  const error = ref<Error | null>(null)
  const hasError = ref(false)

  const handleError = (err: any, fallbackMessage = '操作失败') => {
    console.error('Error:', err)
    error.value = err instanceof Error ? err : new Error(err)
    hasError.value = true

    const errorMessage =
      err?.response?.data?.error?.message ||
      err?.message ||
      fallbackMessage

    message.error(errorMessage)
  }

  const clearError = () => {
    error.value = null
    hasError.value = false
  }

  const resetError = () => {
    clearError()
  }

  return {
    error,
    hasError,
    handleError,
    clearError,
    resetError,
  }
}
