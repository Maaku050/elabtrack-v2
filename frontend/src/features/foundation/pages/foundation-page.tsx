import { Link } from 'react-router-dom'
import { FlaskConical, Moon, Sun } from 'lucide-react'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useUIStore } from '@/stores/ui-store'

export function FoundationPage() {
  const theme = useUIStore((s) => s.theme)
  const toggleTheme = useUIStore((s) => s.toggleTheme)
  return (
    <main className="mx-auto flex min-h-svh max-w-2xl flex-col justify-center gap-6 p-6">
      <header className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-3"><FlaskConical aria-hidden="true" /><h1 className="text-2xl font-semibold">eLabTrack V2</h1></div>
        <Button variant="outline" size="icon" onClick={toggleTheme} aria-label={theme === 'light' ? 'Switch to dark theme' : 'Switch to light theme'}>
          {theme === 'light' ? <Moon aria-hidden="true" /> : <Sun aria-hidden="true" />}
        </Button>
      </header>
      <Card>
        <CardHeader><CardTitle>Foundation in progress</CardTitle><CardDescription>Modernizing laboratory equipment management for FSMO.</CardDescription></CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-muted-foreground">The application is being prepared for its next development phase.</p>
          <Link to="/status" className={buttonVariants({ variant: 'outline' })}>Check service connection</Link>
        </CardContent>
      </Card>
    </main>
  )
}
