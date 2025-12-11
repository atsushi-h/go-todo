'use client'

import { type Todo, useListTodos } from '../hooks'
import { TodoItem } from './TodoItem'

interface TodoListProps {
  onEdit: (todo: Todo) => void
  onDelete: (todo: Todo) => void
  selectedIds: Set<number>
  onToggleSelection: (id: number) => void
}

export function TodoList({ onEdit, onDelete, selectedIds, onToggleSelection }: TodoListProps) {
  const { data: todos, isLoading, error } = useListTodos()

  if (isLoading) {
    return <p className="text-muted-foreground">Loading...</p>
  }

  if (error) {
    return <p className="text-destructive">Error: {String(error)}</p>
  }

  if (!todos || todos.length === 0) {
    return <p className="text-muted-foreground">No todos yet.</p>
  }

  return (
    <div className="space-y-3">
      {todos.map((todo) => (
        <TodoItem
          key={todo.id}
          todo={todo}
          onEdit={onEdit}
          onDelete={onDelete}
          isSelected={selectedIds.has(todo.id)}
          onToggleSelection={onToggleSelection}
        />
      ))}
    </div>
  )
}
