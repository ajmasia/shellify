import { useState } from 'react'
import { toast } from 'sonner'
import { FAB } from '@/components/FAB'
import { useModal } from '@/components/Modal'
import { api } from '@/api'
import type { ProjectWithCount } from '@/domain'
import { useProjects } from './useProjects'
import { ProjectCard } from './ProjectCard'
import styles from './ProjectList.module.css'

export function ProjectList() {
  const { projects, loading, error, refetch } = useProjects()
  const { openCreateProjectModal, openConfirmModal, closeModal } = useModal()
  const [actionLoading, setActionLoading] = useState(false)

  const handleCreate = () => {
    openCreateProjectModal({
      onSubmit: async (data) => {
        setActionLoading(true)
        try {
          await api.projects.create(data)
          closeModal()
          toast.success('Project created')
          await refetch()
        } catch (err) {
          toast.error(err instanceof Error ? err.message : 'Failed to create project')
        } finally {
          setActionLoading(false)
        }
      },
      loading: actionLoading,
    })
  }

  const handleEdit = (project: ProjectWithCount) => {
    toast.info(`Edit project: ${project.name} (not implemented yet)`)
  }

  const handleDelete = (project: ProjectWithCount) => {
    openConfirmModal({
      title: 'Delete Project',
      message: (
        <>
          Are you sure you want to delete <strong>{project.name}</strong>?
          <br />
          This will also delete all sessions in this project.
        </>
      ),
      confirmText: 'Delete',
      variant: 'danger',
      onConfirm: async () => {
        setActionLoading(true)
        try {
          await api.projects.delete(project.id, true)
          closeModal()
          toast.success('Project deleted')
          await refetch()
        } catch (err) {
          toast.error(err instanceof Error ? err.message : 'Failed to delete project')
        } finally {
          setActionLoading(false)
        }
      },
      loading: actionLoading,
    })
  }

  if (loading) {
    return <div className={styles.loading}>Loading projects...</div>
  }

  if (error) {
    return <div className={styles.error}>{error}</div>
  }

  return (
    <div className={styles.container}>
      {projects.length === 0 ? (
        <div className={styles.empty}>
          <p>No projects yet</p>
          <p className={styles.hint}>Click the + button to create your first project</p>
        </div>
      ) : (
        <div className={styles.grid}>
          {projects.map((project) => (
            <ProjectCard
              key={project.id}
              project={project}
              onEdit={handleEdit}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}
      <FAB onClick={handleCreate} tooltip="Create project" />
    </div>
  )
}
