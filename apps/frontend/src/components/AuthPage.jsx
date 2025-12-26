export default function AuthPage({
  title = 'BookingHub',
  subtitle = 'Бронирование переговорных, студий и оборудования',
  serverStatus,
  error,
  loginForm,
  setLoginForm,
  onLogin,
}) {
  return (
    <div className="auth-bg">
      <div className="auth-shell">
        <div className="auth-card">
          <div className="auth-header">
            <div className="auth-logo">BH</div>
            <div>
              <h1 className="auth-title">{title}</h1>
              <p className="auth-subtitle">{subtitle}</p>
            </div>
          </div>

          <div className="auth-meta">
            <span className={`pill ${serverStatus === 'ок' ? 'pill-ok' : 'pill-bad'}`}>
              Сервер: {serverStatus}
            </span>
            <span className="pill pill-muted">Демо-проект</span>
          </div>

          {error ? <div className="auth-alert">{error}</div> : null}

          <form onSubmit={onLogin} className="auth-form">
            <label className="field">
              <span className="label">Email</span>
              <input
                className="input"
                value={loginForm.email}
                onChange={(e) => setLoginForm({ ...loginForm, email: e.target.value })}
                placeholder="user@bookinghub.local"
                autoComplete="username"
              />
            </label>

            <label className="field">
              <span className="label">Пароль</span>
              <input
                className="input"
                type="password"
                value={loginForm.password}
                onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })}
                placeholder="••••••"
                autoComplete="current-password"
              />
            </label>

            <button className="btn btn-primary" type="submit">
              Войти
            </button>

            <div className="auth-hint">
              <div className="hint-title">Тестовые аккаунты</div>
              <div className="hint-list">
                <div><b>admin@bookinghub.local</b> / 123456</div>
                <div><b>manager@bookinghub.local</b> / 123456</div>
                <div><b>user@bookinghub.local</b> / 123456</div>
              </div>
            </div>
          </form>
        </div>

        <div className="auth-footer">
          <span>© {new Date().getFullYear()} BookingHub</span>
          <span className="dot">•</span>
          <span>Учебный проект</span>
        </div>
      </div>
    </div>
  )
}
