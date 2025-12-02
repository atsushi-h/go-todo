import { QueryClient } from '@tanstack/react-query'

// Type guard for axios error
function isAxiosError(error: unknown): error is { response?: { status?: number } } {
  return error !== null && typeof error === 'object' && 'response' in error
}

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5分間データを新鮮とみなす
      gcTime: 10 * 60 * 1000, // 10分間キャッシュを保持（以前のcacheTime）
      retry: (failureCount, error: unknown) => {
        // 401エラー（認証エラー）の場合はリトライしない
        if (isAxiosError(error) && error.response?.status === 401) {
          return false
        }
        // その他のエラーは最大3回までリトライ
        return failureCount < 3
      },
      retryDelay: (attemptIndex) => Math.min(1000 * (attemptIndex + 1), 5000),
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: (failureCount, error: unknown) => {
        // 401エラー（認証エラー）の場合はリトライしない
        if (isAxiosError(error) && error.response?.status === 401) {
          return false
        }
        // その他のエラーは最大1回までリトライ
        return failureCount < 1
      },
    },
  },
})
