'use client'

import { useQueryClient } from '@tanstack/react-query'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { getGetTodosQueryKey, type ModelTodo, useDeleteTodosId } from '../hooks'

interface TodoDeleteDialogProps {
  todo: ModelTodo | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

export function TodoDeleteDialog({ todo, open, onOpenChange, onSuccess }: TodoDeleteDialogProps) {
  const queryClient = useQueryClient()
  const deleteMutation = useDeleteTodosId()

  const handleDelete = () => {
    if (todo?.id === undefined) return

    deleteMutation.mutate(
      { id: todo.id },
      {
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: getGetTodosQueryKey() })
          onSuccess()
        },
      },
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Todo</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete &quot;{todo?.title}&quot;? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={deleteMutation.isPending}
          >
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={deleteMutation.isPending}>
            {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
