import { useEffect, useState } from 'react'

export default function App() {
  const [status, setStatus] = useState('loading...')
  const [resources, setResources] = useState([])
  const [form, setForm] = useState({
    categoryId: 1,
    title: '',
    description: '',
    location: '',
  })
  const [error, setError] = useState('')

  const load = async () => {
    const res = await fetch('/api/resources')
    if (!res.ok) throw new Error(await res.text())
    return res.json()
  }

  useEffect(() => {
    fetch('/api/health')
      .then(r => r.text())
      .then(setStatus)
      .catch(() => setStatus('error'))

    load()
      .then(setResources)
      .catch(e => setError(String(e.message || e)))
  }, [])

  const onCreate = async (e) => {
    e.preventDefault()
    setError('')

    const payload = {
      categoryId: Number(form.categoryId),
      title: form.title.trim(),
      description: form.description.trim() ? form.description.trim() : null,
      location: form.location.trim() ? form.location.trim() : null,
    }

    const res = await fetch('/api/resources', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })

    if (!res.ok) {
      setError(await res.text())
      return
    }

    setForm({ categoryId: 1, title: '', description: '', location: '' })
    const fresh = await load()
    setResources(fresh)
  }

  return (
    <div style={{ padding: 20, maxWidth: 900 }}>
      <h1>BookingHub</h1>
      <p>Backend status: {status}</p>

      <h2>Resources</h2>

      {error && (
        <div style={{ padding: 10, background: '#ffe3e3', marginBottom: 12 }}>
          {error}
        </div>
      )}

      <form onSubmit={onCreate} style={{ display: 'grid', gap: 8, marginBottom: 16 }}>
        <label>
          Category ID (пока руками):
          <input
            value={form.categoryId}
            onChange={(e) => setForm({ ...form, categoryId: e.target.value })}
          />
        </label>

        <label>
          Title:
          <input
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
          />
        </label>

        <label>
          Description:
          <input
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
          />
        </label>

        <label>
          Location:
          <input
            value={form.location}
            onChange={(e) => setForm({ ...form, location: e.target.value })}
          />
        </label>

        <button type="submit">Create resource</button>
      </form>

      <ul>
        {resources.map((r) => (
          <li key={r.id}>
            <b>{r.title}</b> (categoryId: {r.categoryId}) {r.location ? `— ${r.location}` : ''}
          </li>
        ))}
      </ul>
    </div>
  )
}