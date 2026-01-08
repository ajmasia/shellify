import { useNavigate } from 'react-router-dom'
import { IconButton } from '@/components/IconButton'
import { Tooltip } from '@/components/Tooltip'
import { EditIcon, DeleteIcon, FolderIcon } from '@/components/Icons'
import { useUIStore } from '@/stores'
import type { ProjectWithCount } from '@/domain'
import styles from './ProjectCard.module.css'

interface ProjectCardProps {
  project: ProjectWithCount
  onEdit: (project: ProjectWithCount) => void
  onDelete: (project: ProjectWithCount) => void
}

export function ProjectCard({ project, onEdit, onDelete }: ProjectCardProps) {
  const navigate = useNavigate()
  const uiStore = useUIStore()

  const handleClick = () => {
    uiStore.setCurrentProject(project.id, project.name)
    navigate(`/projects/${project.id}`)
  }

  const handleEdit = (e: React.MouseEvent) => {
    e.stopPropagation()
    onEdit(project)
  }

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation()
    onDelete(project)
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString()
  }

  return (
    <article className={styles.card} onClick={handleClick}>
      <div className={styles.header}>
        <FolderIcon size={24} className={styles.icon} />
        <div className={styles.actions}>
          <Tooltip content="Edit">
            <IconButton aria-label="Edit project" onClick={handleEdit}>
              <EditIcon size={16} />
            </IconButton>
          </Tooltip>
          <Tooltip content="Delete">
            <IconButton aria-label="Delete project" variant="danger" onClick={handleDelete}>
              <DeleteIcon size={16} />
            </IconButton>
          </Tooltip>
        </div>
      </div>
      <h3 className={styles.name}>{project.name}</h3>
      {project.description && <p className={styles.description}>{project.description}</p>}
      <div className={styles.footer}>
        <span className={styles.sessions}>{project.sessionCount} sessions</span>
        <span className={styles.date}>{formatDate(project.updatedAt)}</span>
      </div>
    </article>
  )
}
