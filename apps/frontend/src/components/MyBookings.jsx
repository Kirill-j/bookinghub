export default function MyBookings({ me, myBookings, onCancelBooking }) {
  if (!me) return null

  return (
    <div style={{ padding: 12, border: '1px solid #ddd', marginBottom: 16 }}>
      <h2 style={{ marginTop: 0 }}>Мои бронирования</h2>

      {(myBookings?.length || 0) === 0 ? (
        <div style={{ opacity: 0.75 }}>Пока нет бронирований</div>
      ) : (
        <ul>
          {myBookings.map((b) => (
            <li key={b.id} style={{ marginBottom: 8 }}>
              <b>#{b.id}</b> ресурс #{b.resourceId} — {String(b.startAt)} → {String(b.endAt)} — <b>{b.status}</b>{' '}
              {(b.status === 'PENDING' || b.status === 'APPROVED') && (
                <button style={{ marginLeft: 8 }} onClick={() => onCancelBooking(b.id)}>
                  Отменить
                </button>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
