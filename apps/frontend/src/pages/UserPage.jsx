import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { apiJson } from '../api/client'

export default function UserPage({ token }) {
  const { id } = useParams()
  const nav = useNavigate()
  const [u, setU] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    let alive = true

    ;(async () => {
      try {
        const data = await apiJson(`/api/users/${id}`, {}, token)
        if (!alive) return
        setU(data)
        setError('') // ✅ очищаем ошибку тут, а не синхронно в начале эффекта
      } catch {
        if (!alive) return
        setU(null)
        setError('Пользователь не найден')
      }
    })()

    return () => {
      alive = false
    }
  }, [id, token])

  return (
    <div>
      <div className="profile-head">
        <div>
          <h2 style={{ margin: 0 }}>Профиль пользователя</h2>
          <div className="muted">Информация об аккаунте</div>
        </div>
        <button className="btn-ui" onClick={() => nav(-1)}>Назад</button>
      </div>

      {error ? <div className="alert-ui">{error}</div> : null}

      {u ? (
        <div className="card" style={{ marginTop: 12 }}>
          <div style={{ fontWeight: 900, fontSize: 16 }}>{u.name}</div>
          <div className="muted" style={{ marginTop: 6 }}>
            Роль: <b>{u.role}</b>
          </div>
          {u.email ? (
            <div className="muted" style={{ marginTop: 6 }}>
              Email: <b>{u.email}</b>
            </div>
          ) : null}
        </div>
      ) : null}
    </div>
  )
}
