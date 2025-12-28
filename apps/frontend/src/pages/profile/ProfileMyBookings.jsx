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

export default function ProfileMyBookings({ token, resources }) {
  const [error, setError] = useState('')
  const [myBookings, setMyBookings] = useState([])

  const resourceTitleById = useMemo(() => {
    const m = new Map()
    for (const r of (resources || [])) m.set(String(r.id), r.title)
    return m
  }, [resources])

  useEffect(() => {
    let alive = true

    ;(async () => {
      setError('')
      try {
        const items = await apiJson('/api/bookings/my', {}, token)
        if (!alive) return
        setMyBookings(Array.isArray(items) ? items : [])
      } catch {
        if (!alive) return
        setError('Не удалось загрузить мои бронирования')
        setMyBookings([])
      }
    })()

    return () => {
      alive = false
    }
  }, [token])

  return (
    <div className="card">
      <h3 style={{ margin: '0 0 10px' }}>Мои бронирования</h3>
      {error ? <div className="alert-ui">{error}</div> : null}

      {myBookings.length === 0 ? (
        <div className="muted">Пока нет бронирований.</div>
      ) : (
        <div className="list-col">
          {myBookings.map((b) => (
            <div key={b.id} className="list-item" style={{ alignItems: 'flex-start' }}>
              <div>
                <div style={{ fontWeight: 900 }}>
                  {resourceTitleById.get(String(b.resourceId)) || `Ресурс #${b.resourceId}`}
                </div>
                <div className="muted" style={{ fontSize: 12 }}>
                  {fmt(b.startAt)} — {fmt(b.endAt)}
                </div>
                <div className="muted" style={{ fontSize: 12 }}>
                  Статус: <b>{b.status}</b>
                  {b.managerComment ? ` · Комментарий: ${b.managerComment}` : ''}
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
