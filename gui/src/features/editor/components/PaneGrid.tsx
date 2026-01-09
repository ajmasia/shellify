import { observer } from 'mobx-react-lite'
import type { EditorPane, Direction } from '../types'
import { Pane } from './Pane'
import styles from './PaneGrid.module.css'

interface PaneGridProps {
  panes: EditorPane[]
  direction: Direction
}

export const PaneGrid = observer(function PaneGrid({ panes, direction }: PaneGridProps) {
  return (
    <div
      className={styles.grid}
      style={{
        flexDirection: direction === 'horizontal' ? 'row' : 'column',
      }}
    >
      {panes.map((pane, index) => (
        <Pane key={pane.id} pane={pane} direction={direction} isLast={index === panes.length - 1} />
      ))}
    </div>
  )
})
