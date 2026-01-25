import { Button } from '@/components/ui/button'
import { FolderPlus, Plus } from 'lucide-react'

interface EmptyStateProps {
  title: string
  description: string
  buttonText: string
  buttonOnClick: () => void
}

export default function EmptyState({
  title,
  description,
  buttonText,
  buttonOnClick
}: EmptyStateProps) {
  return (
    <div data-testid="empty-state" className="text-center">
      <FolderPlus className="mx-auto h-12 w-12 text-muted-foreground" />
      <h3 className="mt-2 text-sm font-semibold text-foreground">{title}</h3>
      <p className="mt-1 text-sm text-muted-foreground">{description}</p>
      <div className="mt-6">
        <Button onClick={buttonOnClick}>
          <Plus className="mr-1.5 -ml-0.5 h-5 w-5" />
          {buttonText}
        </Button>
      </div>
    </div>
  )
}
