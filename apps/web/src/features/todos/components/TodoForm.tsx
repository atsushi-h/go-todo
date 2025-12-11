'use client'

import { useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { getListTodosQueryKey, type Todo, useCreateTodo, useUpdateTodo } from '../hooks'

interface TodoFormProps {
  todo?: Todo
  onSuccess: () => void
  onCancel?: () => void
}

export function TodoForm({ todo, onSuccess, onCancel }: TodoFormProps) {
  const [title, setTitle] = useState(todo?.title ?? '')
  const [description, setDescription] = useState(todo?.description ?? '')

  const queryClient = useQueryClient()
  const createMutation = useCreateTodo()
  const updateMutation = useUpdateTodo()

  const isEditing = !!todo
  const isPending = createMutation.isPending || updateMutation.isPending
  const isValid = title.trim().length > 0

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!isValid) return

    if (isEditing && todo) {
      updateMutation.mutate(
        {
          id: todo.id,
          data: {
            title: title.trim(),
            description: description.trim() || undefined,
          },
        },
        {
          onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: getListTodosQueryKey() })
            onSuccess()
          },
        },
      )
    } else {
      createMutation.mutate(
        {
          data: {
            title: title.trim(),
            description: description.trim() || undefined,
          },
        },
        {
          onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: getListTodosQueryKey() })
            setTitle('')
            setDescription('')
            onSuccess()
          },
        },
      )
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="title">Title</Label>
        <Input
          id="title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Enter todo title"
          disabled={isPending}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <Textarea
          id="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Enter todo description (optional)"
          disabled={isPending}
        />
      </div>
      <div className="flex gap-2">
        <Button type="submit" disabled={!isValid || isPending}>
          {isPending ? 'Saving...' : isEditing ? 'Update' : 'Add Todo'}
        </Button>
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancel
          </Button>
        )}
      </div>
    </form>
  )
}
