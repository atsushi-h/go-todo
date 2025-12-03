'use client'

import { Button } from '@/components/ui/button'
import { getLoginUrl } from '../hooks/useAuth'

interface LoginButtonProps {
  className?: string
}

export function LoginButton({ className }: LoginButtonProps) {
  const handleLogin = () => {
    window.location.href = getLoginUrl()
  }

  return (
    <Button onClick={handleLogin} className={className}>
      Login with Google
    </Button>
  )
}
