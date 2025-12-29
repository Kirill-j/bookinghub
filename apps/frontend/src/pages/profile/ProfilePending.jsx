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

export default function ProfilePending({ token, resources }) {
  const [error, setError] = useState('')
  const [items, setItems] = useState([])

  // комментарий отдельно для каждой брони
  const [commentById, setCommentById] = useState({})

  // кэш пользователей (кто бронирует)
  const [userById, setUserById] = useState({})

  const resourceTitleById = useMemo(() => {
    const m = new Map()
    for (const r of resources || []) m.set(String(r.id), r.title)
    return m
  }, [resources])

  const reload = async () => {
    setError('')
    try {
      const data = await apiJson('/api/bookings/pending', {}, token)
      setItems(Array.isArray(data) ? data : [])
    } catch (e) {
      setError(String(e.message || e))
      setItems([])
    }
  }

  // грузим pending при заходе
  useEffect(() => {
    reload().catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token])

  // подгрузка данных по пользователям, которые встречаются в items
  useEffect(() => {
    let alive = true

    ;(async () => {
      try {
        const uniqueIds = Array.from(new Set((items || []).map((x) => x.userId))).filter(Boolean)
        const missing = uniqueIds.filter((id) => !userById[id])

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
        // ничего, это просто доп. информация
      }
    })()

    return () => {
      alive = false
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [items, token]) // userById намеренно не добавляем, чтобы не зациклить

  const setStatus = async (bookingId, status) => {
    setError('')
    try {
      const comment = (commentById[bookingId] || '').trim()

      await apiJson(
        `/api/bookings/${bookingId}/status`,
        {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            status,
            managerComment: comment ? comment : null,
          }),
        },
        token
      )

      // очистим комментарий именно у этой карточки
      setCommentById((m) => {
        const x = { ...m }
        delete x[bookingId]
        return x
      })

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
          {items.map((b) => {
            const booker = userById[b.userId]

            return (
              <div key={b.id} className="list-item" style={{ alignItems: 'flex-start' }}>
                <div style={{ width: '100%' }}>
                  <div style={{ fontWeight: 900 }}>
                    {resourceTitleById.get(String(b.resourceId)) || `Ресурс #${b.resourceId}`}
                  </div>

                  <div className="muted" style={{ fontSize: 12 }}>
                    Бронирует:{' '}
                    {booker ? (
                      <Link className="link-btn" to={`/users/${booker.id}`}>
                        {booker.name}
                      </Link>
                    ) : (
                      `Пользователь #${b.userId}`
                    )}
                  </div>

                  <div className="muted" style={{ fontSize: 12 }}>
                    {fmt(b.startAt)} — {fmt(b.endAt)}
                  </div>

                  <div style={{ marginTop: 8 }}>
                    <input
                      className="input-ui"
                      placeholder="Комментарий (опционально)"
                      value={commentById[b.id] || ''}
                      onChange={(e) => setCommentById((m) => ({ ...m, [b.id]: e.target.value }))}
                    />
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

                <div className={`status-badge ${statusClass(b.status)}`}>
                  {statusRu(b.status)}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
