import { useEffect, useMemo, useState } from 'react'
import { apiJson } from '../../api/client'

function rub(n) {
  const v = Number(n) || 0
  return new Intl.NumberFormat('ru-RU').format(v)
}

export default function ProfileMyListings({ token, categories }) {
  const [error, setError] = useState('')
  const [myResources, setMyResources] = useState([])

  const categoryNameById = useMemo(() => {
    const m = new Map()
    for (const c of (categories || [])) m.set(String(c.id), c.name)
    return m
  }, [categories])

  useEffect(() => {
    let alive = true

    ;(async () => {
      setError('')
      try {
        const res = await apiJson('/api/resources/my', {}, token)
        if (!alive) return
        setMyResources(Array.isArray(res) ? res : [])
      } catch {
        if (!alive) return
        setError('Не удалось загрузить мои объявления')
        setMyResources([])
      }
    })()

    return () => {
      alive = false
    }
  }, [token])

  return (
    <div className="card">
      <h3 style={{ margin: '0 0 10px' }}>Мои объявления</h3>
      {error ? <div className="alert-ui">{error}</div> : null}

      {myResources.length === 0 ? (
        <div className="muted">Пока нет объявлений.</div>
      ) : (
        <div className="list-col">
          {myResources.map((r) => (
            <div key={r.id} className="list-item">
              <div>
                <div style={{ fontWeight: 900 }}>{r.title}</div>
                <div className="muted" style={{ fontSize: 12 }}>
                  {categoryNameById.get(String(r.categoryId)) || 'Категория'} ·{' '}
                  {r.location || 'Локация не указана'}
                </div>
              </div>
              <div style={{ fontWeight: 900 }}>{rub(r.pricePerHour)} ₽/час</div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
