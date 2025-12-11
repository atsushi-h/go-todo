'use client'

import { useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { getListTodosQueryKey, type Todo, useBatchCompleteTodos, useListTodos } from '../hooks'
import { BatchDeleteDialog } from './BatchDeleteDialog'
import { BatchOperationToolbar } from './BatchOperationToolbar'
import { TodoDeleteDialog } from './TodoDeleteDialog'
import { TodoForm } from './TodoForm'
import { TodoList } from './TodoList'

export function TodoPageContent() {
  const queryClient = useQueryClient()
  const { data: todos } = useListTodos()
  const batchCompleteMutation = useBatchCompleteTodos()

  const [editingTodo, setEditingTodo] = useState<Todo | null>(null)
  const [deletingTodo, setDeletingTodo] = useState<Todo | null>(null)
  const [selectedIds, setSelectedIds] = useState<Set<number>>(new Set())
  const [batchDeleteDialogOpen, setBatchDeleteDialogOpen] = useState(false)

  const handleToggleSelection = (id: number) => {
    setSelectedIds((prev) => {
      const next = new Set(prev)
      if (next.has(id)) {
        next.delete(id)
      } else {
        next.add(id)
      }
      return next
    })
  }

  const handleClearSelection = () => {
    setSelectedIds(new Set())
  }

  const handleBatchComplete = () => {
    if (selectedIds.size === 0) return

    batchCompleteMutation.mutate(
      { data: { ids: Array.from(selectedIds) } },
      {
        onSuccess: (response) => {
          queryClient.invalidateQueries({ queryKey: getListTodosQueryKey() })

          if (response.failed && response.failed.length > 0) {
            console.error('Some todos failed to complete:', response.failed)
          } else {
            handleClearSelection()
          }
        },
        onError: (error) => {
          console.error('Batch complete error:', error)
        },
      },
    )
  }

  const handleBatchDelete = () => {
    setBatchDeleteDialogOpen(true)
  }

  return (
    <div className="container mx-auto py-8 max-w-2xl px-4">
      <h1 className="text-2xl font-bold mb-6">Todos</h1>

      <Card className="mb-8">
        <CardHeader>
          <CardTitle>Add New Todo</CardTitle>
        </CardHeader>
        <CardContent>
          <TodoForm onSuccess={() => {}} />
        </CardContent>
      </Card>

      <Dialog open={!!editingTodo} onOpenChange={(open) => !open && setEditingTodo(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Todo</DialogTitle>
          </DialogHeader>
          {editingTodo && (
            <TodoForm
              todo={editingTodo}
              onSuccess={() => setEditingTodo(null)}
              onCancel={() => setEditingTodo(null)}
            />
          )}
        </DialogContent>
      </Dialog>

      <TodoDeleteDialog
        todo={deletingTodo}
        open={!!deletingTodo}
        onOpenChange={(open) => !open && setDeletingTodo(null)}
        onSuccess={() => setDeletingTodo(null)}
      />

      <BatchDeleteDialog
        todos={todos ?? []}
        selectedIds={selectedIds}
        open={batchDeleteDialogOpen}
        onOpenChange={setBatchDeleteDialogOpen}
        onSuccess={handleClearSelection}
      />

      <TodoList
        onEdit={setEditingTodo}
        onDelete={setDeletingTodo}
        selectedIds={selectedIds}
        onToggleSelection={handleToggleSelection}
      />

      <BatchOperationToolbar
        selectedCount={selectedIds.size}
        onBatchComplete={handleBatchComplete}
        onBatchDelete={handleBatchDelete}
        onClearSelection={handleClearSelection}
        isLoading={batchCompleteMutation.isPending}
      />
    </div>
  )
}
