import { useEffect, useState } from 'react'

export default function App() {
  const [status, setStatus] = useState('loading...')

  useEffect(() => {
    fetch('api/health')
      .then(r => r.text())
      .then(setStatus)
      .catch(() => setStatus('error'))
  }, [])

  return (
    <div style={{ padding: 20 }}>
      <h1>BookingHub</h1>
      <p>Backend status: {status}</p>
    </div>
  )
}