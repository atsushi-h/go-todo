'use client'

import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useAuth } from '../hooks/useAuth'

interface DeleteAccountDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function DeleteAccountDialog({ open, onOpenChange }: DeleteAccountDialogProps) {
  const router = useRouter()
  const { deleteAccount, isDeletingAccount } = useAuth()
  const [error, setError] = useState<string | null>(null)

  const handleDelete = () => {
    setError(null)
    deleteAccount(undefined, {
      onSuccess: () => {
        onOpenChange(false)
        // Redirect to home page after successful deletion
        router.push('/')
      },
      onError: (err) => {
        console.error('Failed to delete account:', err)
        setError('Failed to delete account. Please try again.')
      },
    })
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      setError(null)
    }
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Account</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete your account? This action cannot be undone. All your
            todos will be permanently deleted.
          </DialogDescription>
        </DialogHeader>
        {error && (
          <div className="rounded-md bg-red-50 p-3 text-sm text-red-800 dark:bg-red-900/20 dark:text-red-400">
            {error}
          </div>
        )}
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => handleOpenChange(false)}
            disabled={isDeletingAccount}
          >
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={isDeletingAccount}>
            {isDeletingAccount ? 'Deleting...' : 'Delete Account'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
