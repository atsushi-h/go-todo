'use client'

import { CheckCheck, Trash2, X } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface BatchOperationToolbarProps {
  selectedCount: number
  onBatchComplete: () => void
  onBatchDelete: () => void
  onClearSelection: () => void
  isLoading?: boolean
}

export function BatchOperationToolbar({
  selectedCount,
  onBatchComplete,
  onBatchDelete,
  onClearSelection,
  isLoading = false,
}: BatchOperationToolbarProps) {
  if (selectedCount === 0) {
    return null
  }

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-background/95 backdrop-blur-sm border-t shadow-lg transition-transform duration-200 ease-in-out">
      <div className="container mx-auto max-w-2xl px-4 py-4">
        <div className="flex items-center justify-between gap-4">
          <p className="text-sm font-medium">
            {selectedCount} {selectedCount === 1 ? 'item' : 'items'} selected
          </p>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" onClick={onBatchComplete} disabled={isLoading}>
              <CheckCheck className="size-4 mr-2" />
              Complete Selected
            </Button>
            <Button variant="outline" size="sm" onClick={onBatchDelete} disabled={isLoading}>
              <Trash2 className="size-4 mr-2" />
              Delete Selected
            </Button>
            <Button variant="ghost" size="sm" onClick={onClearSelection} disabled={isLoading}>
              <X className="size-4 mr-2" />
              Clear
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
