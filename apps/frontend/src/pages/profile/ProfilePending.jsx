import { useEffect, useMemo, useState } from 'react'
import { apiJson } from '../../api/client'

function fmt(dt) {
  if (!dt) return ''
  try {
    return new Date(dt).toLocaleString('ru-RU')
  } catch {
    return String(dt)
  }
}

export default function ProfilePending({ token, resources }) {
  const [error, setError] = useState('')
  const [items, setItems] = useState([])
  const [comment, setComment] = useState('')

  const resourceTitleById = useMemo(() => {
    const m = new Map()
    for (const r of (resources || [])) m.set(String(r.id), r.title)
    return m
  }, [resources])

  const reload = async () => {
    setError('')
    const data = await apiJson('/api/bookings/pending', {}, token)
    setItems(Array.isArray(data) ? data : [])
  }

  useEffect(() => {
    let alive = true

    ;(async () => {
      try {
        const data = await apiJson('/api/bookings/pending', {}, token)
        if (!alive) return
        setItems(Array.isArray(data) ? data : [])
      } catch {
        if (!alive) return
        setError('Не удалось загрузить заявки на подтверждение')
        setItems([])
      }
    })()

    return () => {
      alive = false
    }
  }, [token])

  const setStatus = async (id, status) => {
    setError('')
    try {
      await apiJson(
        `/api/bookings/${id}/status`,
        {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            status,
            managerComment: comment.trim() ? comment.trim() : null,
          }),
        },
        token
      )
      setComment('')
      await reload()
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  return (
    <div className="card">
      <h3 style={{ margin: '0 0 10px' }}>Подтверждение брони</h3>
      {error ? <div className="alert-ui">{error}</div> : null}

      {items.length === 0 ? (
        <div className="muted">Нет заявок на подтверждение.</div>
      ) : (
        <div className="list-col">
          <label className="field-ui">
            <span className="label-ui">Комментарий (опционально)</span>
            <input
              className="input-ui"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Например: подтверждаю, ключи у охраны"
            />
          </label>

          {items.map((b) => (
            <div key={b.id} className="list-item" style={{ alignItems: 'flex-start' }}>
              <div>
                <div style={{ fontWeight: 900 }}>
                  {resourceTitleById.get(String(b.resourceId)) || `Ресурс #${b.resourceId}`}
                </div>
                <div className="muted" style={{ fontSize: 12 }}>
                  {fmt(b.startAt)} — {fmt(b.endAt)}
                </div>

                <div style={{ display: 'flex', gap: 10, marginTop: 10, flexWrap: 'wrap' }}>
                  <button className="btn-ui" type="button" onClick={() => setStatus(b.id, 'APPROVED')}>
                    Подтвердить
                  </button>
                  <button className="btn-ui" type="button" onClick={() => setStatus(b.id, 'REJECTED')}>
                    Отклонить
                  </button>
                </div>
              </div>

              <div className="role-tag">{b.status}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
