'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { useAuth } from '../hooks/useAuth'
import { LogoutButton } from './LogoutButton'
import { DeleteAccountDialog } from './DeleteAccountDialog'

interface UserProfileProps {
  className?: string
}

export function UserProfile({ className }: UserProfileProps) {
  const { user } = useAuth()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  if (!user) {
    return null
  }

  return (
    <div className={className}>
      <div className="flex items-center gap-4">
        {user.avatar_url && (
          <Image
            src={user.avatar_url}
            alt={user.name}
            width={48}
            height={48}
            className="size-12 rounded-full"
          />
        )}
        <div>
          <p className="font-medium text-zinc-900 dark:text-zinc-50">{user.name}</p>
          <p className="text-sm text-zinc-600 dark:text-zinc-400">{user.email}</p>
        </div>
      </div>
      <div className="mt-6 flex flex-col gap-3 sm:flex-row">
        <Link
          href="/todos"
          className="flex h-10 items-center justify-center rounded-md bg-primary px-6 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
        >
          Go to Todos
        </Link>
        <LogoutButton />
        <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
          Delete Account
        </Button>
      </div>

      <DeleteAccountDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen} />
    </div>
  )
}
