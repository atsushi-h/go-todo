'use client'

import { type ModelTodo, useGetTodos } from '../hooks'
import { TodoItem } from './TodoItem'

interface TodoListProps {
  onEdit: (todo: ModelTodo) => void
  onDelete: (todo: ModelTodo) => void
}

export function TodoList({ onEdit, onDelete }: TodoListProps) {
  const { data: todos, isLoading, error } = useGetTodos()

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
        <TodoItem key={todo.id} todo={todo} onEdit={onEdit} onDelete={onDelete} />
      ))}
    </div>
  )
}
