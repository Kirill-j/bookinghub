import { useEffect, useMemo, useState } from 'react'

function getToken() {
  return localStorage.getItem('accessToken') || ''
}

function setToken(token) {
  if (token) localStorage.setItem('accessToken', token)
  else localStorage.removeItem('accessToken')
}

async function apiText(path, token) {
  const r = await fetch(path, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  })
  if (!r.ok) throw new Error(await r.text())
  return r.text()
}

async function apiJson(path, opts = {}, token) {
  const headers = { ...(opts.headers || {}) }
  if (token) headers.Authorization = `Bearer ${token}`
  const r = await fetch(path, { ...opts, headers })
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}

export default function App() {
  const [serverStatus, setServerStatus] = useState('загрузка...')
  const [error, setError] = useState('')

  const [token, setTokenState] = useState(getToken())
  const [me, setMe] = useState(null)

  const [categories, setCategories] = useState([])
  const [resources, setResources] = useState([])

  const [myBookings, setMyBookings] = useState([])

  const [pendingBookings, setPendingBookings] = useState([])
  const [managerComment, setManagerComment] = useState('')

  const [loginForm, setLoginForm] = useState({
    email: 'manager@bookinghub.local',
    password: '123456',
  })

  const [resourceForm, setResourceForm] = useState({
    categoryId: '',
    title: '',
    description: '',
    location: '',
  })

  const [bookingForm, setBookingForm] = useState({
    resourceId: '',
    date: '',
    start: '10:00',
    end: '11:00',
  })

  const categoryNameById = useMemo(() => {
    const m = new Map()
    for (const c of categories) m.set(String(c.id), c.name)
    return m
  }, [categories])

  const loadPublic = async () => {
    const [cats, res] = await Promise.all([
      apiJson('/api/categories', {}, token),
      apiJson('/api/resources', {}, token),
    ])
    setCategories(cats)
    setResources(res)

    if (!resourceForm.categoryId && cats.length) {
      setResourceForm((f) => ({ ...f, categoryId: String(cats[0].id) }))
    }
    if (!bookingForm.resourceId && res.length) {
      setBookingForm((f) => ({ ...f, resourceId: String(res[0].id) }))
    }
  }

  const loadMe = async (t) => {
    if (!t) {
      setMe(null)
      return
    }
    try {
      const u = await apiJson('/api/auth/me', {}, t)
      setMe(u)
    } catch {
      setMe(null)
      setToken('')
      setTokenState('')
    }
  }

  const loadMyBookings = async (t) => {
    if (!t) {
      setMyBookings([])
      return
    }
    const items = await apiJson('/api/bookings/my', {}, t)
    setMyBookings(Array.isArray(items) ? items : [])
  }

  const loadPending = async (t) => {
    if (!t) {
      setPendingBookings([])
      return
    }
    const items = await apiJson('/api/bookings/pending', {}, t)
    setPendingBookings(Array.isArray(items) ? items : [])
  }

  const updateBookingStatus = async (id, status) => {
    setError('')
    try {
      await apiJson(
        `/api/bookings/${id}/status`,
        {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            status,
            managerComment: managerComment.trim() ? managerComment.trim() : null,
          }),
        },
        token
      )

      setManagerComment('')
      await loadPending(token)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  useEffect(() => {
    apiText('/api/health')
      .then(() => setServerStatus('ок'))
      .catch(() => setServerStatus('ошибка'))

    loadPublic().catch((e) => setError(String(e.message || e)))
    loadMe(token).catch(() => {})
    loadMyBookings(token).catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const onLogin = async (e) => {
    e.preventDefault()
    setError('')

    try {
      const data = await apiJson(
        '/api/auth/login',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            email: loginForm.email.trim().toLowerCase(),
            password: loginForm.password,
          }),
        },
        ''
      )

      const t = data.accessToken
      setToken(t)
      setTokenState(t)
      await loadMe(t)
      await loadMyBookings(t)
      await loadPending(t)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const onLogout = async () => {
    setToken('')
    setTokenState('')
    setMe(null)
    setMyBookings([])
    setPendingBookings([])
    setManagerComment('')
  }

  const canCreateResources = me && (me.role === 'MANAGER' || me.role === 'ADMIN')

  const onCreateResource = async (e) => {
    e.preventDefault()
    setError('')

    const payload = {
      categoryId: Number(resourceForm.categoryId),
      title: resourceForm.title.trim(),
      description: resourceForm.description.trim() ? resourceForm.description.trim() : null,
      location: resourceForm.location.trim() ? resourceForm.location.trim() : null,
    }

    if (!payload.categoryId) return setError('Выберите категорию')
    if (!payload.title) return setError('Название обязательно')

    try {
      await apiJson(
        '/api/resources',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        },
        token
      )
      setResourceForm((f) => ({ ...f, title: '', description: '', location: '' }))
      const fresh = await apiJson('/api/resources', {}, token)
      setResources(fresh)
      if (!bookingForm.resourceId && fresh.length) {
        setBookingForm((f) => ({ ...f, resourceId: String(fresh[0].id) }))
      }
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const onCreateBooking = async (e) => {
    e.preventDefault()
    setError('')

    if (!me) return setError('Сначала войдите в систему')
    if (!bookingForm.resourceId) return setError('Выберите ресурс')
    if (!bookingForm.date) return setError('Выберите дату')

    const startAt = `${bookingForm.date}T${bookingForm.start}:00`
    const endAt = `${bookingForm.date}T${bookingForm.end}:00`

    try {
      await apiJson(
        '/api/bookings',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            resourceId: Number(bookingForm.resourceId),
            startAt,
            endAt,
          }),
        },
        token
      )

      await loadMyBookings(token)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const onCancelBooking = async (id) => {
    setError('')
    try {
      await apiJson(`/api/bookings/${id}/cancel`, { method: 'POST' }, token)
      await loadMyBookings(token)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  return (
    <div style={{ padding: 20, maxWidth: 900 }}>
      <h1>BookingHub</h1>
      <p>Статус сервера: {serverStatus}</p>

      {error && (
        <div style={{ padding: 10, background: '#ffe3e3', marginBottom: 12 }}>
          {error}
        </div>
      )}

      <div style={{ padding: 12, border: '1px solid #ddd', marginBottom: 16 }}>
        <h2 style={{ marginTop: 0 }}>Вход</h2>

        {me ? (
          <>
            <p>
              Вы вошли как: <b>{me.name}</b> ({me.email}), роль: <b>{me.role}</b>
            </p>
            <button onClick={onLogout}>Выйти</button>
          </>
        ) : (
          <form onSubmit={onLogin} style={{ display: 'grid', gap: 8 }}>
            <label>
              Email:
              <input
                value={loginForm.email}
                onChange={(e) => setLoginForm({ ...loginForm, email: e.target.value })}
              />
            </label>
            <label>
              Пароль:
              <input
                type="password"
                value={loginForm.password}
                onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })}
              />
            </label>
            <button type="submit">Войти</button>
            <div style={{ fontSize: 12, opacity: 0.75 }}>
              Тестовые аккаунты: admin@bookinghub.local / manager@bookinghub.local / user@bookinghub.local,
              пароль 123456
            </div>
          </form>
        )}
      </div>

      {canCreateResources ? (
        <>
          <h2>Создать ресурс</h2>
          <form onSubmit={onCreateResource} style={{ display: 'grid', gap: 8, marginBottom: 16 }}>
            <label>
              Категория:
              <select
                value={resourceForm.categoryId}
                onChange={(e) => setResourceForm({ ...resourceForm, categoryId: e.target.value })}
                style={{ marginLeft: 8 }}
              >
                {categories.map((c) => (
                  <option key={c.id} value={String(c.id)}>
                    {c.name}
                  </option>
                ))}
              </select>
            </label>

            <label>
              Название:
              <input
                value={resourceForm.title}
                onChange={(e) => setResourceForm({ ...resourceForm, title: e.target.value })}
              />
            </label>

            <label>
              Описание:
              <input
                value={resourceForm.description}
                onChange={(e) => setResourceForm({ ...resourceForm, description: e.target.value })}
              />
            </label>

            <label>
              Место/локация:
              <input
                value={resourceForm.location}
                onChange={(e) => setResourceForm({ ...resourceForm, location: e.target.value })}
              />
            </label>

            <button type="submit">Создать</button>
          </form>
        </>
      ) : (
        <div style={{ padding: 12, border: '1px dashed #bbb', marginBottom: 16 }}>
          <b>Создание ресурсов доступно только менеджеру или администратору.</b>
          <div style={{ fontSize: 12, opacity: 0.75 }}>
            Войдите как manager@bookinghub.local или admin@bookinghub.local.
          </div>
        </div>
      )}

      {me && (
        <div style={{ padding: 12, border: '1px solid #ddd', marginBottom: 16 }}>
          <h2 style={{ marginTop: 0 }}>Создать бронирование</h2>

          <form onSubmit={onCreateBooking} style={{ display: 'grid', gap: 8 }}>
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

            <button type="submit">Забронировать</button>
            <div style={{ fontSize: 12, opacity: 0.75 }}>
              Правило: минимум 30 минут. Отмена — не позднее чем за 2 часа до начала.
            </div>
          </form>
        </div>
      )}

      {me && (
        <div style={{ padding: 12, border: '1px solid #ddd', marginBottom: 16 }}>
          <h2 style={{ marginTop: 0 }}>Мои бронирования</h2>

          {(myBookings?.length || 0) === 0 ? (
            <div style={{ opacity: 0.75 }}>Пока нет бронирований</div>
          ) : (
            <ul>
              {myBookings.map((b) => (
                <li key={b.id} style={{ marginBottom: 8 }}>
                  <b>#{b.id}</b> ресурс #{b.resourceId} —{' '}
                  {new Date(b.startAt).toLocaleString()} → {new Date(b.endAt).toLocaleString()} —{' '}
                  <b>{b.status}</b>{' '}
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
      )}

      {me && (me.role === 'MANAGER' || me.role === 'ADMIN') && (
        <div style={{ padding: 12, border: '1px solid #ddd', marginBottom: 16 }}>
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
                    <button onClick={() => updateBookingStatus(b.id, 'APPROVED')}>
                      Подтвердить
                    </button>
                    <button
                      style={{ marginLeft: 8 }}
                      onClick={() => updateBookingStatus(b.id, 'REJECTED')}
                    >
                      Отклонить
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}

      <h2>Ресурсы</h2>
      <ul>
        {resources.map((r) => (
          <li key={r.id}>
            <b>{r.title}</b>{' '}
            <span style={{ opacity: 0.7 }}>
              ({categoryNameById.get(String(r.categoryId)) || `категория #${r.categoryId}`})
            </span>
            {r.location ? ` — ${r.location}` : ''}
          </li>
        ))}
      </ul>
    </div>
  )
}
