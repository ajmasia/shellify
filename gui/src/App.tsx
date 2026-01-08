import { Routes, Route } from 'react-router-dom'
import { MainLayout } from '@/components/Layout'
import { ProjectList } from '@/features/projects'
import { SessionList } from '@/features/sessions'

export function App() {
  return (
    <MainLayout>
      <Routes>
        <Route path="/" element={<ProjectList />} />
        <Route path="/projects/:projectId" element={<SessionList />} />
      </Routes>
    </MainLayout>
  )
}
