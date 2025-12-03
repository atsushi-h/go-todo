'use client'

import { LoginButton, UserProfile, useAuth } from '@/features/auth'

export default function Home() {
  const { isAuthenticated, isLoading } = useAuth()

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-center gap-8 bg-white px-16 py-32 dark:bg-black sm:items-start">
        <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50">
          Go Todo
        </h1>

        {isLoading ? (
          <p className="text-zinc-600 dark:text-zinc-400">Loading...</p>
        ) : isAuthenticated ? (
          <UserProfile />
        ) : (
          <div className="flex flex-col gap-4">
            <p className="text-zinc-600 dark:text-zinc-400">Sign in to manage your todos.</p>
            <LoginButton />
          </div>
        )}
      </main>
    </div>
  )
}
