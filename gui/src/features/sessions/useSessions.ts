import { useState, useEffect, useCallback } from 'react'
import { api } from '@/api'
import type { Session } from '@/domain'

interface UseSessionsResult {
  sessions: Session[]
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export function useSessions(projectId: string): UseSessionsResult {
  const [sessions, setSessions] = useState<Session[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchSessions = useCallback(async () => {
    if (!projectId) return

    try {
      setLoading(true)
      setError(null)
      const data = await api.projects.sessions(projectId)
      setSessions(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch sessions')
    } finally {
      setLoading(false)
    }
  }, [projectId])

  useEffect(() => {
    fetchSessions()
  }, [fetchSessions])

  return { sessions, loading, error, refetch: fetchSessions }
}
