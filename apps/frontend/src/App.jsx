import { useEffect, useMemo, useState } from 'react'
import { apiJson, apiText, getToken, saveToken } from './api/client'

import AuthPanel from './components/AuthPanel'
import ResourceCreateForm from './components/ResourceCreateForm'
import BookingCreateForm from './components/BookingCreateForm'
import MyBookings from './components/MyBookings'
import ManagerPending from './components/ManagerPending'
import ResourceList from './components/ResourceList'
import AuthPage from './components/AuthPage'

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

  const [resourceBookings, setResourceBookings] = useState([])

  const [bookingDurationMin, setBookingDurationMin] = useState(60)

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

  const canCreateResources = me && (me.role === 'MANAGER' || me.role === 'ADMIN')

  // ---- loaders ----
  const loadPublic = async () => {
    const [cats, res] = await Promise.all([
      apiJson('/api/categories', {}, token),
      apiJson('/api/resources', {}, token),
    ])

    setCategories(Array.isArray(cats) ? cats : [])
    setResources(Array.isArray(res) ? res : [])

    if (!resourceForm.categoryId && Array.isArray(cats) && cats.length) {
      setResourceForm((f) => ({ ...f, categoryId: String(cats[0].id) }))
    }
    if (!bookingForm.resourceId && Array.isArray(res) && res.length) {
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
      saveToken('')
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

  const loadResourceBookings = async (resourceId, date) => {
    if (!resourceId || !date) {
      setResourceBookings([])
      return
    }
    const url = `/api/resources/${resourceId}/bookings?from=${date}&to=${date}`
    const items = await apiJson(url, {}, token)
    setResourceBookings(Array.isArray(items) ? items : [])
  }

  // ---- effects ----
  useEffect(() => {
    apiText('/api/health')
      .then(() => setServerStatus('ок'))
      .catch(() => setServerStatus('ошибка'))

    loadPublic().catch((e) => setError(String(e.message || e)))
    loadMe(token).catch(() => {})
    loadMyBookings(token).catch(() => {})
    loadPending(token).catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    loadResourceBookings(bookingForm.resourceId, bookingForm.date).catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [bookingForm.resourceId, bookingForm.date])

  // ---- actions ----
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
      saveToken(t)
      setTokenState(t)

      await loadMe(t)
      await loadMyBookings(t)
      await loadPending(t)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const onLogout = async () => {
    saveToken('')
    setTokenState('')
    setMe(null)
    setMyBookings([])
    setPendingBookings([])
    setManagerComment('')
    setResourceBookings([])
  }

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
      const arr = Array.isArray(fresh) ? fresh : []
      setResources(arr)

      if (!bookingForm.resourceId && arr.length) {
        setBookingForm((f) => ({ ...f, resourceId: String(arr[0].id) }))
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
      await loadResourceBookings(bookingForm.resourceId, bookingForm.date)
      await loadPending(token)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const onCancelBooking = async (id) => {
    setError('')
    try {
      await apiJson(`/api/bookings/${id}/cancel`, { method: 'POST' }, token)
      await loadMyBookings(token)
      await loadResourceBookings(bookingForm.resourceId, bookingForm.date)
      await loadPending(token)
    } catch (e) {
      setError(String(e.message || e))
    }
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

  if (!me) {
    return (
      <AuthPage
        serverStatus={serverStatus}
        error={error}
        loginForm={loginForm}
        setLoginForm={setLoginForm}
        onLogin={onLogin}
      />
    )
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

      <AuthPanel
        me={me}
        loginForm={loginForm}
        setLoginForm={setLoginForm}
        onLogin={onLogin}
        onLogout={onLogout}
      />

      <ResourceCreateForm
        categories={categories}
        resourceForm={resourceForm}
        setResourceForm={setResourceForm}
        onCreateResource={onCreateResource}
        canCreateResources={canCreateResources}
      />

      <BookingCreateForm
        me={me}
        resources={resources}
        bookingForm={bookingForm}
        setBookingForm={setBookingForm}
        bookingDurationMin={bookingDurationMin}
        setBookingDurationMin={setBookingDurationMin}
        resourceBookings={resourceBookings}
        onCreateBooking={onCreateBooking}
        setError={setError}
      />

      <MyBookings me={me} myBookings={myBookings} onCancelBooking={onCancelBooking} />

      <ManagerPending
        me={me}
        pendingBookings={pendingBookings}
        managerComment={managerComment}
        setManagerComment={setManagerComment}
        loadPending={loadPending}
        updateBookingStatus={updateBookingStatus}
        token={token}
      />

      <ResourceList resources={resources} categoryNameById={categoryNameById} />
    </div>
  )
}
