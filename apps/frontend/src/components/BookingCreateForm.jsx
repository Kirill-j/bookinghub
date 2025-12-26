import { hhmmToMinutes, minutesToHHMM, overlaps } from '../utils/time'
import OccupancyList from './OccupancyList'

export default function BookingCreateForm({
  me,
  resources,
  bookingForm,
  setBookingForm,
  bookingDurationMin,
  setBookingDurationMin,
  resourceBookings,
  onCreateBooking,
  pickFreeSlot, // опционально: можно передать готовую функцию
  setError,
}) {
  if (!me) return null

  const localPick = () => {
    // Если функцию не передали — сделаем встроенную
    if (pickFreeSlot) return pickFreeSlot()

    setError('')
    if (!bookingForm.date) return setError('Сначала выберите дату')
    if (!bookingForm.resourceId) return setError('Сначала выберите ресурс')

    const duration = Number(bookingDurationMin) || 60
    const dayStart = 8 * 60
    const dayEnd = 20 * 60

    const busy = (Array.isArray(resourceBookings) ? resourceBookings : [])
      .map((b) => {
        const s = new Date(b.startAt)
        const e = new Date(b.endAt)
        if (Number.isNaN(s.getTime()) || Number.isNaN(e.getTime())) return null
        return { start: s.getHours() * 60 + s.getMinutes(), end: e.getHours() * 60 + e.getMinutes() }
      })
      .filter(Boolean)

    for (let t = dayStart; t + duration <= dayEnd; t += 15) {
      const candidateStart = t
      const candidateEnd = t + duration

      const hasConflict = busy.some((x) => overlaps(candidateStart, candidateEnd, x.start, x.end))
      if (!hasConflict) {
        setBookingForm((f) => ({
          ...f,
          start: minutesToHHMM(candidateStart),
          end: minutesToHHMM(candidateEnd),
        }))
        return
      }
    }
    setError('На выбранную дату нет свободного окна в рабочее время')
  }

  const onSubmit = (e) => {
    // фронт-валидация времени
    if (hhmmToMinutes(bookingForm.end) <= hhmmToMinutes(bookingForm.start)) {
      e.preventDefault()
      setError('Время окончания должно быть позже времени начала')
      return
    }
    onCreateBooking(e)
  }

  return (
    <div style={{ padding: 12, border: '1px solid #ddd', marginBottom: 16 }}>
      <h2 style={{ marginTop: 0 }}>Создать бронирование</h2>

      <form onSubmit={onSubmit} style={{ display: 'grid', gap: 8 }}>
        <label>
          Ресурс:
          <select
            value={bookingForm.resourceId}
            onChange={(e) => setBookingForm({ ...bookingForm, resourceId: e.target.value })}
            style={{ marginLeft: 8 }}
          >
            {resources.map((r) => (
              <option key={r.id} value={String(r.id)}>
                {r.title}
              </option>
            ))}
          </select>
        </label>

        <label>
          Дата:
          <input
            type="date"
            value={bookingForm.date}
            onChange={(e) => setBookingForm({ ...bookingForm, date: e.target.value })}
          />
        </label>

        <label>
          Начало:
          <input
            type="time"
            value={bookingForm.start}
            onChange={(e) => setBookingForm({ ...bookingForm, start: e.target.value })}
          />
        </label>

        <label>
          Конец:
          <input
            type="time"
            value={bookingForm.end}
            onChange={(e) => setBookingForm({ ...bookingForm, end: e.target.value })}
          />
        </label>

        <label>
          Длительность (мин):
          <input
            type="number"
            min="30"
            step="15"
            value={bookingDurationMin}
            onChange={(e) => setBookingDurationMin(e.target.value)}
            style={{ marginLeft: 8, width: 90 }}
          />
        </label>

        <button type="button" onClick={localPick}>
          Подобрать свободное время
        </button>

        <button type="submit">Забронировать</button>

        <div style={{ fontSize: 12, opacity: 0.75 }}>
          Правило: минимум 30 минут. Отмена — не позднее чем за 2 часа до начала.
        </div>

        <OccupancyList resourceBookings={resourceBookings} />
      </form>
    </div>
  )
}
