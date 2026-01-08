export function App() {
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100%',
        gap: 'var(--space-md)',
      }}
    >
      <h1 style={{ color: 'var(--ctp-mauve)' }}>Shellify</h1>
      <p style={{ color: 'var(--text-secondary)' }}>Terminal Session Manager</p>
    </div>
  )
}
