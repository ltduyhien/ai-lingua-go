import { type InputHTMLAttributes, forwardRef } from 'react'
import { cn } from '@/lib/utils'

const Input = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement>>(
  ({ className, type, ...props }, ref) => (
    <input
      type={type}
      className={cn('app-input w-full', className)}
      ref={ref}
      {...props}
    />
  )
)
Input.displayName = 'Input'

export { Input }
