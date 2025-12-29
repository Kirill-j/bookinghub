import { useEffect, useState } from 'react'
import { apiJson } from '../api/client'

export default function AdminCategoriesPage({ token, me }) {
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')
  const [items, setItems] = useState([])
  const [newName, setNewName] = useState('')
  const [editById, setEditById] = useState({}) // { [id]: "new name" }

  const isAdmin = me?.role === 'ADMIN'

  const load = async () => {
    setError('')
    setOk('')
    try {
      const cats = await apiJson('/api/categories', {}, token)
      setItems(Array.isArray(cats) ? cats : [])
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  useEffect(() => {
    load().catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const create = async (e) => {
    e.preventDefault()
    setError('')
    setOk('')

    const name = newName.trim()
    if (!name) return setError('Введите название категории')

    try {
      await apiJson(
        '/api/categories',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name }),
        },
        token
      )
      setNewName('')
      setOk('Категория создана')
      await load()
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const save = async (id) => {
    setError('')
    setOk('')

    const name = String(editById[id] || '').trim()
    if (!name) return setError('Название не может быть пустым')

    try {
      await apiJson(
        `/api/categories/${id}`,
        {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name }),
        },
        token
      )
      setEditById((m) => {
        const x = { ...m }
        delete x[id]
        return x
      })
      setOk('Категория обновлена')
      await load()
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const remove = async (id) => {
    setError('')
    setOk('')

    if (!confirm('Удалить категорию? Если к ней привязаны объявления — сервер не даст удалить.')) return

    try {
      await apiJson(`/api/categories/${id}`, { method: 'DELETE' }, token)
      setOk('Категория удалена')
      await load()
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  if (!isAdmin) {
    return (
      <div className="card">
        <h3 style={{ margin: '0 0 10px' }}>Админ</h3>
        <div className="muted">Доступно только администратору.</div>
      </div>
    )
  }

  return (
    <div className="grid" style={{ gap: 14 }}>
      <div className="card">
        <h3 style={{ margin: '0 0 10px' }}>Админ · Категории</h3>
        <div className="muted">Создание, переименование и удаление категорий.</div>

        {error ? <div className="alert-ui">{error}</div> : null}
        {ok ? <div className="ok-ui">{ok}</div> : null}

        <form onSubmit={create} className="form-col" style={{ marginTop: 12 }}>
          <label className="field-ui">
            <span className="label-ui">Новая категория</span>
            <input className="input-ui" value={newName} onChange={(e) => setNewName(e.target.value)} placeholder="Например: Автомобиль" />
          </label>
          <button className="btn-ui" type="submit">Добавить</button>
        </form>
      </div>

      <div className="card">
        <div className="profile-head" style={{ marginTop: 0 }}>
          <h3 style={{ margin: 0 }}>Список категорий</h3>
          <button className="btn-ui" type="button" onClick={load}>Обновить</button>
        </div>

        {items.length === 0 ? (
          <div className="muted">Категорий пока нет.</div>
        ) : (
          <div className="list-col" style={{ marginTop: 12 }}>
            {items.map((c) => {
              const edit = editById[c.id] ?? c.name
              return (
                <div key={c.id} className="list-item" style={{ alignItems: 'center' }}>
                  <div style={{ width: '100%' }}>
                    <div className="muted" style={{ fontSize: 12, marginBottom: 6 }}>#{c.id}</div>
                    <input
                      className="input-ui"
                      value={edit}
                      onChange={(e) => setEditById((m) => ({ ...m, [c.id]: e.target.value }))}
                    />
                    <div style={{ display: 'flex', gap: 10, marginTop: '12px', flexWrap: 'wrap' }}>
                        <button className="btn-ui" type="button" onClick={() => save(c.id)}>Сохранить</button>
                        <button className="btn-del" type="button" onClick={() => remove(c.id)}>Удалить</button>
                    </div>  
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
