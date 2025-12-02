'use client'

import { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { TodoDeleteDialog, TodoForm, TodoList } from '@/features/todos'
import type { ModelTodo } from '@/features/todos/hooks'

export default function TodosPage() {
  const [editingTodo, setEditingTodo] = useState<ModelTodo | null>(null)
  const [deletingTodo, setDeletingTodo] = useState<ModelTodo | null>(null)

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

      <TodoList onEdit={setEditingTodo} onDelete={setDeletingTodo} />
    </div>
  )
}
