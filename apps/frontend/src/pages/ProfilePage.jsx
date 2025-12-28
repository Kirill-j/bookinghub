import { useEffect, useMemo, useState } from 'react'
import { apiJson, saveToken } from '../api/client'

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

export default function ProfilePage({ token, me, categories, resources, onMeUpdated }) {
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  const [myResources, setMyResources] = useState([])
  const [myBookings, setMyBookings] = useState([])

  // формы
  const [profileForm, setProfileForm] = useState({
    name: me?.name || '',
    email: me?.email || '',
  })

  const [passForm, setPassForm] = useState({
    currentPassword: '',
    newPassword: '',
    newPassword2: '',
  })

  useEffect(() => {
    // если me обновился — обновим форму
    setProfileForm({ name: me?.name || '', email: me?.email || '' })
  }, [me])

  const categoryNameById = useMemo(() => {
    const m = new Map()
    for (const c of categories) m.set(String(c.id), c.name)
    return m
  }, [categories])

  const resourceTitleById = useMemo(() => {
    const m = new Map()
    for (const r of (resources || [])) m.set(String(r.id), r.title)
    for (const r of (myResources || [])) m.set(String(r.id), r.title)
    return m
  }, [resources, myResources])

  const loadSummary = async () => {
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
    loadSummary().catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const saveProfile = async (e) => {
    e.preventDefault()
    setError('')
    setOk('')

    const name = profileForm.name.trim()
    const email = profileForm.email.trim().toLowerCase()
    if (!name) return setError('Введите имя')
    if (!email || !email.includes('@')) return setError('Введите корректный email')

    try {
      const updated = await apiJson(
        '/api/auth/me',
        {
          method: 'PATCH',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name, email }),
        },
        token
      )

      onMeUpdated?.(updated) // обновим me в App.jsx
      setOk('Профиль обновлён')
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const changePassword = async (e) => {
    e.preventDefault()
    setError('')
    setOk('')

    if (!passForm.currentPassword) return setError('Введите текущий пароль')
    if (!passForm.newPassword || passForm.newPassword.length < 6) return setError('Новый пароль минимум 6 символов')
    if (passForm.newPassword !== passForm.newPassword2) return setError('Новые пароли не совпадают')

    try {
      await apiJson(
        '/api/auth/password',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            currentPassword: passForm.currentPassword,
            newPassword: passForm.newPassword,
          }),
        },
        token
      )

      setPassForm({ currentPassword: '', newPassword: '', newPassword2: '' })
      setOk('Пароль изменён')
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const [deleteText, setDeleteText] = useState('')

  const deleteAccount = async () => {
    setError('')
    setOk('')

    if (deleteText.trim().toUpperCase() !== 'УДАЛИТЬ') {
      return setError('Для подтверждения введи: УДАЛИТЬ')
    }

    try {
      await apiJson('/api/auth/me', { method: 'DELETE' }, token)
      // после удаления просто разлогиниваемся
      saveToken('')
      window.location.reload()
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  return (
    <div>
      <div className="profile-head">
        <div>
          <h2 style={{ margin: 0 }}>Профиль</h2>
          <div className="muted">{me?.name} · {me?.email} · {me?.role}</div>
        </div>
        <button className="btn-ui" onClick={loadSummary}>Обновить</button>
      </div>

      {error ? <div className="alert-ui">{error}</div> : null}
      {ok ? <div className="ok-ui">{ok}</div> : null}

      <div className="grid grid-2" style={{ marginTop: 12 }}>
        <div className="card">
          <h3 style={{ margin: '0 0 10px' }}>Личные данные</h3>

          <form onSubmit={saveProfile} className="form-col">
            <label className="field-ui">
              <span className="label-ui">Имя</span>
              <input
                className="input-ui"
                value={profileForm.name}
                onChange={(e) => setProfileForm({ ...profileForm, name: e.target.value })}
              />
            </label>

            <label className="field-ui">
              <span className="label-ui">Email</span>
              <input
                className="input-ui"
                value={profileForm.email}
                onChange={(e) => setProfileForm({ ...profileForm, email: e.target.value })}
              />
            </label>

            <button className="btn-ui" type="submit">Сохранить</button>
          </form>
        </div>

        <div className="card">
          <h3 style={{ margin: '0 0 10px' }}>Безопасность</h3>

          <form onSubmit={changePassword} className="form-col">
            <label className="field-ui">
              <span className="label-ui">Текущий пароль</span>
              <input
                className="input-ui"
                type="password"
                value={passForm.currentPassword}
                onChange={(e) => setPassForm({ ...passForm, currentPassword: e.target.value })}
              />
            </label>

            <label className="field-ui">
              <span className="label-ui">Новый пароль</span>
              <input
                className="input-ui"
                type="password"
                value={passForm.newPassword}
                onChange={(e) => setPassForm({ ...passForm, newPassword: e.target.value })}
              />
            </label>

            <label className="field-ui">
              <span className="label-ui">Повтор нового пароля</span>
              <input
                className="input-ui"
                type="password"
                value={passForm.newPassword2}
                onChange={(e) => setPassForm({ ...passForm, newPassword2: e.target.value })}
              />
            </label>

            <button className="btn-ui" type="submit">Сменить пароль</button>
          </form>
        </div>
      </div>

      <div className="card" style={{ marginTop: 14 }}>
        <h3 style={{ margin: '0 0 10px' }}>Опасная зона</h3>
        <div className="muted" style={{ marginBottom: 10 }}>
          Удаление аккаунта удалит ваши объявления и связанные бронирования. Отменить нельзя.
        </div>

        <label className="field-ui">
          <span className="label-ui">Введи “УДАЛИТЬ” для подтверждения</span>
          <input className="input-ui" value={deleteText} onChange={(e) => setDeleteText(e.target.value)} />
        </label>

        <button type="button" className="btn-ui" onClick={deleteAccount}>
          Удалить аккаунт
        </button>
      </div>

      {/* Сводка (по желанию можно оставить ниже) */}
      <div className="grid grid-2" style={{ marginTop: 14 }}>
        <div className="card">
          <h3 style={{ margin: '0 0 10px' }}>Мои объявления (сводка)</h3>

          {myResources.length === 0 ? (
            <div className="muted">Пока нет объявлений. Зайди в раздел “Разместить объявление”.</div>
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
          <h3 style={{ margin: '0 0 10px' }}>Мои бронирования (сводка)</h3>

          {myBookings.length === 0 ? (
            <div className="muted">Пока нет бронирований.</div>
          ) : (
            <div className="list-col">
              {myBookings.map((b) => (
                <div key={b.id} className="list-item" style={{ alignItems: 'flex-start' }}>
                  <div>
                    <div style={{ fontWeight: 900 }}>
                      {resourceTitleById.get(String(b.resourceId)) || `Ресурс #${b.resourceId}`}
                    </div>
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
