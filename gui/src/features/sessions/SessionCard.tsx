import { useNavigate } from 'react-router-dom'
import { IconButton } from '@/components/IconButton'
import { Tooltip } from '@/components/Tooltip'
import { Chip } from '@/components/Chip'
import { DeleteIcon, EditIcon } from '@/components/Icons'
import type { Session } from '@/domain'
import styles from './SessionCard.module.css'

interface SessionCardProps {
  session: Session
  onDelete: (session: Session) => void
}

export function SessionCard({ session, onDelete }: SessionCardProps) {
  const navigate = useNavigate()

  const handleClick = () => {
    navigate(`/sessions/${session.id}/edit`)
  }

  const handleEdit = (e: React.MouseEvent) => {
    e.stopPropagation()
    navigate(`/sessions/${session.id}/edit`)
  }

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation()
    onDelete(session)
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString()
  }

  const windowCount = session.windows?.length ?? 0

  return (
    <article className={styles.card} onClick={handleClick}>
      <div className={styles.header}>
        <Chip variant={session.targetMultiplexer}>{session.targetMultiplexer}</Chip>
        <div className={styles.actions}>
          <Tooltip content="Edit">
            <IconButton aria-label="Edit session" variant="ghost" onClick={handleEdit}>
              <EditIcon size={16} />
            </IconButton>
          </Tooltip>
          <Tooltip content="Delete">
            <IconButton aria-label="Delete session" variant="danger" onClick={handleDelete}>
              <DeleteIcon size={16} />
            </IconButton>
          </Tooltip>
        </div>
      </div>
      <h3 className={styles.name}>{session.name}</h3>
      {session.description && <p className={styles.description}>{session.description}</p>}
      <div className={styles.footer}>
        <span className={styles.windows}>{windowCount} windows</span>
        <span className={styles.date}>{formatDate(session.updatedAt)}</span>
      </div>
    </article>
  )
}
