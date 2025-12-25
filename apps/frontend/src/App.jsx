import { useEffect, useMemo, useState } from 'react'

async function apiText(path) {
  const r = await fetch(path)
  if (!r.ok) throw new Error(await r.text())
  return r.text()
}

async function apiJson(path, opts) {
  const r = await fetch(path, opts)
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}

export default function App() {
  const [status, setStatus] = useState('loading...')
  const [error, setError] = useState('')

  const [categories, setCategories] = useState([])
  const [resources, setResources] = useState([])

  const [form, setForm] = useState({
    categoryId: '',
    title: '',
    description: '',
    location: '',
  })

  const categoryNameById = useMemo(() => {
    const m = new Map()
    for (const c of categories) m.set(String(c.id), c.name)
    return m
  }, [categories])

  const loadAll = async () => {
    const [cats, res] = await Promise.all([
      apiJson('/api/categories'),
      apiJson('/api/resources'),
    ])
    setCategories(cats)
    setResources(res)

    // если categoryId пустой — выберем первую категорию
    if (!form.categoryId && cats.length) {
      setForm((f) => ({ ...f, categoryId: String(cats[0].id) }))
    }
  }

  useEffect(() => {
    apiText('/api/health')
      .then(setStatus)
      .catch(() => setStatus('error'))

    loadAll().catch((e) => setError(String(e.message || e)))
    // eslint-disable-next-line react-hooks/exhaustive-deps
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

    if (!payload.categoryId) {
      setError('Choose a category')
      return
    }
    if (!payload.title) {
      setError('Title is required')
      return
    }

    try {
      await apiJson('/api/resources', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      setForm((f) => ({ ...f, title: '', description: '', location: '' }))
      const fresh = await apiJson('/api/resources')
      setResources(fresh)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  return (
    <div style={{ padding: 20, maxWidth: 900 }}>
      <h1>BookingHub</h1>
      <p>Backend status: {status}</p>

      {error && (
        <div style={{ padding: 10, background: '#ffe3e3', marginBottom: 12 }}>
          {error}
        </div>
      )}

      <h2>Create resource</h2>
      <form onSubmit={onCreate} style={{ display: 'grid', gap: 8, marginBottom: 16 }}>
        <label>
          Category:
          <select
            value={form.categoryId}
            onChange={(e) => setForm({ ...form, categoryId: e.target.value })}
            style={{ marginLeft: 8 }}
          >
            {categories.map((c) => (
              <option key={c.id} value={String(c.id)}>
                {c.name}
              </option>
            ))}
          </select>
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

        <button type="submit">Create</button>
      </form>

      <h2>Resources</h2>
      <ul>
        {resources.map((r) => (
          <li key={r.id}>
            <b>{r.title}</b>{' '}
            <span style={{ opacity: 0.7 }}>
              ({categoryNameById.get(String(r.categoryId)) || `categoryId:${r.categoryId}`})
            </span>
            {r.location ? ` — ${r.location}` : ''}
          </li>
        ))}
      </ul>
    </div>
  )
}
