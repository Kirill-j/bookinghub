import { useEffect, useMemo, useState } from 'react'
import { apiJson } from '../api/client'

function rub(n) {
  const v = Number(n) || 0
  return new Intl.NumberFormat('ru-RU').format(v)
}

function fmt(dt) {
  if (!dt) return ''
  try {
    const d = new Date(dt)
    return d.toLocaleString('ru-RU')
  } catch {
    return String(dt)
  }
}

export default function ProfilePage({ token, me, categories }) {
  const [error, setError] = useState('')
  const [myResources, setMyResources] = useState([])
  const [myBookings, setMyBookings] = useState([])

  const categoryNameById = useMemo(() => {
    const m = new Map()
    for (const c of categories) m.set(String(c.id), c.name)
    return m
  }, [categories])

  const load = async () => {
    setError('')
    try {
      const [res, bookings] = await Promise.all([
        apiJson('/api/resources/my', {}, token),
        apiJson('/api/bookings/my', {}, token),
      ])
      setMyResources(Array.isArray(res) ? res : [])
      setMyBookings(Array.isArray(bookings) ? bookings : [])
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  useEffect(() => {
    load().catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div>
      <div className="profile-head">
        <div>
          <h2 style={{ margin: 0 }}>Профиль</h2>
          <div className="muted">
            {me?.name} · {me?.email} · {me?.role}
          </div>
        </div>
        <button className="btn-ui" onClick={load}>Обновить</button>
      </div>

      {error ? <div className="alert-ui">{error}</div> : null}

      <div className="grid grid-2" style={{ marginTop: 12 }}>
        <div className="card">
          <h3 style={{ margin: '0 0 10px' }}>Мои объявления</h3>

          {myResources.length === 0 ? (
            <div className="muted">Пока нет объявлений. Нажми “Разместить” в шапке.</div>
          ) : (
            <div className="list-col">
              {myResources.map((r) => (
                <div key={r.id} className="list-item">
                  <div>
                    <div style={{ fontWeight: 900 }}>{r.title}</div>
                    <div className="muted" style={{ fontSize: 12 }}>
                      {categoryNameById.get(String(r.categoryId)) || 'Категория'} · {r.location || 'Локация не указана'}
                    </div>
                  </div>
                  <div style={{ fontWeight: 900 }}>{rub(r.pricePerHour)} ₽/час</div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="card">
          <h3 style={{ margin: '0 0 10px' }}>Мои бронирования</h3>

          {myBookings.length === 0 ? (
            <div className="muted">Пока нет бронирований.</div>
          ) : (
            <div className="list-col">
              {myBookings.map((b) => (
                <div key={b.id} className="list-item" style={{ alignItems: 'flex-start' }}>
                  <div>
                    <div style={{ fontWeight: 900 }}>Ресурс #{b.resourceId}</div>
                    <div className="muted" style={{ fontSize: 12 }}>
                      {fmt(b.startAt)} — {fmt(b.endAt)}
                    </div>
                    <div className="muted" style={{ fontSize: 12 }}>
                      Статус: <b>{b.status}</b>
                      {b.managerComment ? ` · Комментарий: ${b.managerComment}` : ''}
                    </div>
                  </div>
                  <div className="pill pill-muted">{b.status}</div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
