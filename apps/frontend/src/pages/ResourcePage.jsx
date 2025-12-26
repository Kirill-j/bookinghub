import { useEffect, useMemo, useState } from 'react'
import { apiJson } from '../api/client'
import OccupancyList from '../components/OccupancyList'

function rub(n) {
  const v = Number(n) || 0
  return new Intl.NumberFormat('ru-RU').format(v)
}

export default function ResourcePage({ id, token, me, resources, onBack, onRefreshAfterBooking }) {
  const resource = useMemo(
    () => (Array.isArray(resources) ? resources.find((x) => String(x.id) === String(id)) : null),
    [resources, id]
  )

  const [date, setDate] = useState('')
  const [resourceBookings, setResourceBookings] = useState([])
  const [start, setStart] = useState('10:00')
  const [end, setEnd] = useState('11:00')
  const [error, setError] = useState('')

  const loadBookings = async () => {
    setError('')
    if (!date) {
      setResourceBookings([])
      return
    }
    const items = await apiJson(`/api/resources/${id}/bookings?from=${date}&to=${date}`, {}, token)
    setResourceBookings(Array.isArray(items) ? items : [])
  }

  useEffect(() => {
    loadBookings().catch((e) => setError(String(e.message || e)))
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, date])

  const onBook = async () => {
    setError('')
    if (!me) return setError('Нужно войти, чтобы бронировать')
    if (!date) return setError('Выберите дату')

    try {
      await apiJson(
        '/api/bookings',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            resourceId: Number(id),
            startAt: `${date}T${start}:00`,
            endAt: `${date}T${end}:00`,
          }),
        },
        token
      )

      await loadBookings()
      await onRefreshAfterBooking?.()
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  if (!resource) {
    return (
      <div className="card">
        <button className="link-btn" onClick={onBack}>← Назад</button>
        <div style={{ marginTop: 10 }}>Ресурс не найден.</div>
      </div>
    )
  }

  return (
    <div>
      <button className="link-btn" onClick={onBack}>← Назад к каталогу</button>

      <div className="resource-hero">
        <div>
          <h2 style={{ margin: 0 }}>{resource.title}</h2>
          <div className="muted">{resource.location || 'Локация не указана'}</div>
        </div>
        <div className="resource-price">{rub(resource.pricePerHour)} ₽/час</div>
      </div>

      {resource.description ? (
        <div className="card" style={{ marginTop: 12 }}>
          <h3 style={{ margin: '0 0 8px' }}>Описание</h3>
          <div style={{ opacity: 0.92, lineHeight: 1.45 }}>{resource.description}</div>
        </div>
      ) : null}

      <div className="grid grid-2" style={{ marginTop: 12 }}>
        <div className="card">
          <h3 style={{ margin: '0 0 10px' }}>Занятость</h3>
          <input className="input-ui" type="date" value={date} onChange={(e) => setDate(e.target.value)} />
          <div style={{ marginTop: 10 }}>
            <OccupancyList resourceBookings={resourceBookings} />
          </div>
        </div>

        <div className="card">
          <h3 style={{ margin: '0 0 10px' }}>Бронирование</h3>
          {!me && <div className="muted" style={{ marginBottom: 8 }}>Войдите, чтобы бронировать</div>}
          {error && <div className="alert-ui">{error}</div>}

          <div className="form-row">
            <label className="label-ui">Начало</label>
            <input className="input-ui" type="time" value={start} onChange={(e) => setStart(e.target.value)} />
          </div>

          <div className="form-row">
            <label className="label-ui">Конец</label>
            <input className="input-ui" type="time" value={end} onChange={(e) => setEnd(e.target.value)} />
          </div>

          <button className="btn-ui" onClick={onBook} disabled={!me}>
            Забронировать
          </button>

          <div className="muted" style={{ marginTop: 8 }}>
            Бронь уйдёт менеджеру на подтверждение.
          </div>
        </div>
      </div>
    </div>
  )
}
