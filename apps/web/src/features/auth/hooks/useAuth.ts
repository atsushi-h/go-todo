'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import axiosInstance from '@/api/axios-instance'

interface User {
  id: number
  email: string
  name: string
  avatar_url: string
  provider: string
}

const getMeQueryKey = ['me'] as const

export function useAuth() {
  const queryClient = useQueryClient()

  const {
    data: user,
    isLoading,
    error,
  } = useQuery({
    queryKey: getMeQueryKey,
    queryFn: async () => {
      const response = await axiosInstance.get<User>('/me')
      return response.data
    },
    retry: false,
  })

  const logoutMutation = useMutation({
    mutationFn: async () => {
      await axiosInstance.post('/logout')
    },
    onSuccess: () => {
      queryClient.setQueryData(getMeQueryKey, null)
      queryClient.invalidateQueries({ queryKey: getMeQueryKey })
    },
  })

  const deleteAccountMutation = useMutation({
    mutationFn: async () => {
      await axiosInstance.delete('/users/me')
    },
    onSuccess: () => {
      queryClient.setQueryData(getMeQueryKey, null)
      queryClient.invalidateQueries({ queryKey: getMeQueryKey })
    },
  })

  const isAuthenticated = !!user && !error

  return {
    user,
    isLoading,
    isAuthenticated,
    logout: logoutMutation.mutate,
    isLoggingOut: logoutMutation.isPending,
    deleteAccount: deleteAccountMutation.mutate,
    isDeletingAccount: deleteAccountMutation.isPending,
  }
}

export function getLoginUrl() {
  return `${process.env.NEXT_PUBLIC_API_URL}/auth/google`
}
