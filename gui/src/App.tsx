import { Routes, Route } from 'react-router-dom'
import { MainLayout } from '@/components/Layout'
import { ProjectList } from '@/features/projects'
import { SessionList } from '@/features/sessions'
import { SessionEditor } from '@/features/editor'
import { useFaviconTheme } from '@/hooks/useFaviconTheme'

export function App() {
  useFaviconTheme()

  return (
    <Routes>
      <Route
        path="/"
        element={
          <MainLayout>
            <ProjectList />
          </MainLayout>
        }
      />
      <Route
        path="/projects/:projectId"
        element={
          <MainLayout>
            <SessionList />
          </MainLayout>
        }
      />
      <Route path="/sessions/:sessionId/edit" element={<SessionEditor />} />
    </Routes>
  )
}
