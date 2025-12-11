'use client'

import { useQueryClient } from '@tanstack/react-query'
import { AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { getListTodosQueryKey, type Todo, useBatchDeleteTodos } from '../hooks'

interface BatchDeleteDialogProps {
  todos: Todo[]
  selectedIds: Set<number>
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
}

export function BatchDeleteDialog({
  todos,
  selectedIds,
  open,
  onOpenChange,
  onSuccess,
}: BatchDeleteDialogProps) {
  const queryClient = useQueryClient()
  const batchDeleteMutation = useBatchDeleteTodos()

  const selectedTodos = todos.filter((todo) => selectedIds.has(todo.id))
  const selectedCount = selectedTodos.length

  const handleBatchDelete = () => {
    if (selectedIds.size === 0) return

    batchDeleteMutation.mutate(
      { data: { ids: Array.from(selectedIds) } },
      {
        onSuccess: (response) => {
          queryClient.invalidateQueries({ queryKey: getListTodosQueryKey() })

          if (response.failed && response.failed.length > 0) {
            console.error('Some todos failed to delete:', response.failed)
          }

          onSuccess()
          onOpenChange(false)
        },
        onError: (error) => {
          console.error('Batch delete error:', error)
        },
      },
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            Delete {selectedCount} {selectedCount === 1 ? 'Todo' : 'Todos'}
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete the following {selectedCount === 1 ? 'todo' : 'todos'}?
            This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        <div className="max-h-60 overflow-y-auto space-y-2 py-4">
          {selectedTodos.map((todo) => (
            <div key={todo.id} className="text-sm border rounded p-2">
              <p className="font-medium">{todo.title}</p>
              {todo.description && (
                <p className="text-muted-foreground text-xs mt-1">{todo.description}</p>
              )}
            </div>
          ))}
        </div>

        {batchDeleteMutation.isError && (
          <div className="flex items-start gap-2 p-3 bg-destructive/10 border border-destructive/20 rounded">
            <AlertCircle className="size-4 text-destructive mt-0.5" />
            <div className="flex-1 text-sm text-destructive">
              <p className="font-medium">Failed to delete todos</p>
              <p className="text-xs mt-1">Please try again</p>
            </div>
          </div>
        )}

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={batchDeleteMutation.isPending}
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleBatchDelete}
            disabled={batchDeleteMutation.isPending}
          >
            {batchDeleteMutation.isPending ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
