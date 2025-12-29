import { useEffect, useMemo, useState } from 'react'
import { apiJson } from '../../api/client'
import { Link } from 'react-router-dom'
import { statusRu, statusClass } from '../../utils/status'

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

  const [userById, setUserById] = useState({})

  const resourceTitleById = useMemo(() => {
    const m = new Map()
    for (const r of resources || []) m.set(String(r.id), r.title)
    return m
  }, [resources])

  const resourceOwnerById = useMemo(() => {
    const m = new Map()
    for (const r of resources || []) m.set(String(r.id), r.ownerUserId)
    return m
  }, [resources])

  const cancelBooking = async (id) => {
  if (!window.confirm('Отменить эту бронь?')) return

  setError('')
  try {
    await apiJson(
      `/api/bookings/${id}/cancel`,
      { method: 'POST' },
      token
    )

    // обновим список
    setMyBookings((prev) =>
      prev.map((b) =>
        b.id === id ? { ...b, status: 'CANCELED' } : b
      )
    )
  } catch (e) {
    setError(String(e.message || e))
  }
}

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

  // подгрузим владельцев объявлений, которые встречаются в моих бронированиях
  useEffect(() => {
    let alive = true

    ;(async () => {
      try {
        const ownerIds = Array.from(
          new Set(
            (myBookings || [])
              .map((b) => resourceOwnerById.get(String(b.resourceId)))
              .filter(Boolean)
          )
        )

        const missing = ownerIds.filter((id) => !userById[id])
        if (missing.length === 0) return

        const results = await Promise.all(
          missing.map(async (uid) => {
            try {
              const u = await apiJson(`/api/users/${uid}`, {}, token)
              return [uid, u]
            } catch {
              return [uid, null]
            }
          })
        )

        if (!alive) return

        setUserById((prev) => {
          const next = { ...prev }
          for (const [uid, u] of results) next[uid] = u
          return next
        })
      } catch {
        // ок
      }
    })()

    return () => {
      alive = false
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [myBookings, token, resourceOwnerById])

  return (
    <div className="card">
      <h3 style={{ margin: '0 0 10px' }}>Мои бронирования</h3>
      {error ? <div className="alert-ui">{error}</div> : null}

      {myBookings.length === 0 ? (
        <div className="muted">Пока нет бронирований.</div>
      ) : (
        <div className="list-col">
          {myBookings.map((b) => {
            const ownerId = resourceOwnerById.get(String(b.resourceId))
            const owner = ownerId ? userById[ownerId] : null

            return (
              <div key={b.id} className="list-item" style={{ alignItems: 'flex-start' }}>
                <div>
                  <div style={{ fontWeight: 900 }}>
                    {resourceTitleById.get(String(b.resourceId)) || `Ресурс #${b.resourceId}`}
                  </div>

                  <div className="muted" style={{ fontSize: 12 }}>
                    Владелец:{' '}
                    {owner ? (
                      <Link className="link-btn" to={`/users/${owner.id}`}>
                        {owner.name}
                      </Link>
                    ) : ownerId ? (
                      `Пользователь #${ownerId}`
                    ) : (
                      '—'
                    )}
                  </div>

                  <div className="muted" style={{ fontSize: 12 }}>
                    {fmt(b.startAt)} — {fmt(b.endAt)}
                  </div>

                  <div className="muted" style={{ fontSize: 12 }}>
                    Статус: <b>{statusRu(b.status)}</b>
                    {b.managerComment ? ` · Комментарий: ${b.managerComment}` : ''}
                  </div>
                </div>
                <div className="list-col">
                  <div className={`status-badge ${statusClass(b.status)}`}>
                    {statusRu(b.status)}
                  </div>

                  {(b.status === 'PENDING' || b.status === 'APPROVED') && (
                    <button
                      type="button"
                      className="btn-del"
                      style={{ marginTop: 8 }}
                      onClick={() => cancelBooking(b.id)}
                    >
                      Отменить бронь
                    </button>
                  )}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
