export default function ManagerPending({
  me,
  pendingBookings,
  managerComment,
  setManagerComment,
  loadPending,
  updateBookingStatus,
  token,
}) {
  const isManager = me && (me.role === 'MANAGER' || me.role === 'ADMIN')
  if (!isManager) return null

  return (
    <div>
      <h2 style={{ marginTop: 0 }}>Заявки на бронирование</h2>

      <div style={{ display: 'grid', gap: 8, marginBottom: 10 }}>
        <label>
          Комментарий менеджера (необязательно):
          <input
            value={managerComment}
            onChange={(e) => setManagerComment(e.target.value)}
            placeholder="Например: подтверждено / занято / не подходит время"
            style={{ width: '100%' }}
          />
        </label>

        <button onClick={() => loadPending(token)}>Обновить список</button>
      </div>

      {(pendingBookings?.length || 0) === 0 ? (
        <div style={{ opacity: 0.75 }}>Нет заявок</div>
      ) : (
        <ul>
          {pendingBookings.map((b) => (
            <li key={b.id} style={{ marginBottom: 10 }}>
              <b>#{b.id}</b> — ресурс #{b.resourceId}, пользователь #{b.userId}
              <div style={{ opacity: 0.75 }}>
                {String(b.startAt)} → {String(b.endAt)}
              </div>
              <div style={{ marginTop: 6 }}>
                <button onClick={() => updateBookingStatus(b.id, 'APPROVED')}>Подтвердить</button>
                <button style={{ marginLeft: 8 }} onClick={() => updateBookingStatus(b.id, 'REJECTED')}>
                  Отклонить
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
