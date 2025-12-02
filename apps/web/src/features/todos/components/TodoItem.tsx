'use client'

import { useQueryClient } from '@tanstack/react-query'
import { Pencil, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { getGetTodosQueryKey, type ModelTodo, usePutTodosId } from '../hooks'

interface TodoItemProps {
  todo: ModelTodo
  onEdit: (todo: ModelTodo) => void
  onDelete: (todo: ModelTodo) => void
}

export function TodoItem({ todo, onEdit, onDelete }: TodoItemProps) {
  const queryClient = useQueryClient()
  const updateMutation = usePutTodosId()

  const handleToggleCompleted = () => {
    if (todo.id === undefined) return

    updateMutation.mutate(
      {
        id: todo.id,
        data: {
          completed: !todo.completed,
        },
      },
      {
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: getGetTodosQueryKey() })
        },
      },
    )
  }

  return (
    <Card className="py-4">
      <CardContent className="flex items-start gap-4">
        <Checkbox
          checked={todo.completed ?? false}
          onCheckedChange={handleToggleCompleted}
          disabled={updateMutation.isPending}
          className="mt-1"
        />
        <div className="flex-1 min-w-0">
          <p
            className={`font-medium ${todo.completed ? 'line-through text-muted-foreground' : ''}`}
          >
            {todo.title}
          </p>
          {todo.description && (
            <p className="text-sm text-muted-foreground mt-1">{todo.description}</p>
          )}
        </div>
        <div className="flex gap-2">
          <Button variant="ghost" size="icon-sm" onClick={() => onEdit(todo)}>
            <Pencil className="size-4" />
          </Button>
          <Button variant="ghost" size="icon-sm" onClick={() => onDelete(todo)}>
            <Trash2 className="size-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
