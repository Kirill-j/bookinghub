import { Link } from 'react-router-dom'

export default function AppShell({ me, onLogout, children }) {
  return (
    <div className="app-bg">
      <div className="app-topbar">
        <div className="app-topbar-inner">
          <div className="brand">
            <div className="brand-badge">BH</div>
            <div>
              <h1 className="brand-title">BookingHub</h1>
              <p className="brand-sub">Каталог и бронирование ресурсов</p>
            </div>
          </div>

          <div className="topbar-right">
            <div className="topbar-nav">
              <Link className="nav-link" to="/">Каталог</Link>
              <Link className="nav-link" to="/new">Разместить</Link>
              <Link className="nav-link" to="/profile">Профиль</Link>
            </div>

            <div className="user-pill">
              <span>{me?.name || 'Пользователь'}</span>
              <span className="muted">{me?.email}</span>
              <span className="role-tag">{me?.role}</span>
            </div>

            <button className="btn-ui" onClick={onLogout}>
              Выйти
            </button>
          </div>
        </div>
      </div>

      <div className="app-container">
        {children}
      </div>
    </div>
  )
}
