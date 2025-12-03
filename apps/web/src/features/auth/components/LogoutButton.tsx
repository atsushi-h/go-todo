'use client'

import { Button } from '@/components/ui/button'
import { useAuth } from '../hooks/useAuth'

interface LogoutButtonProps {
  className?: string
}

export function LogoutButton({ className }: LogoutButtonProps) {
  const { logout, isLoggingOut } = useAuth()

  return (
    <Button
      variant="outline"
      onClick={() => logout()}
      disabled={isLoggingOut}
      className={className}
    >
      {isLoggingOut ? 'Logging out...' : 'Logout'}
    </Button>
  )
}
