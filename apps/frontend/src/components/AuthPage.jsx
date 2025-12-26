export default function AuthPage({
  title = 'BookingHub',
  subtitle = 'Бронирование переговорных, студий и оборудования',
  serverStatus,
  error,

  mode, // 'login' | 'register'
  setMode,

  loginForm,
  setLoginForm,
  onLogin,

  registerForm,
  setRegisterForm,
  onRegister,
}) {
  const isLogin = mode === 'login'

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

          {/* Переключатель */}
          <div className="auth-tabs">
            <button
              type="button"
              className={`tab ${isLogin ? 'tab-active' : ''}`}
              onClick={() => setMode('login')}
            >
              Вход
            </button>
            <button
              type="button"
              className={`tab ${!isLogin ? 'tab-active' : ''}`}
              onClick={() => setMode('register')}
            >
              Регистрация
            </button>
          </div>

          {error ? <div className="auth-alert">{error}</div> : null}

          {isLogin ? (
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
          ) : (
            <form onSubmit={onRegister} className="auth-form">
              <label className="field">
                <span className="label">Имя</span>
                <input
                  className="input"
                  value={registerForm.name}
                  onChange={(e) => setRegisterForm({ ...registerForm, name: e.target.value })}
                  placeholder="Например: Кирилл"
                  autoComplete="name"
                />
              </label>

              <label className="field">
                <span className="label">Email</span>
                <input
                  className="input"
                  value={registerForm.email}
                  onChange={(e) => setRegisterForm({ ...registerForm, email: e.target.value })}
                  placeholder="you@example.com"
                  autoComplete="email"
                />
              </label>

              <label className="field">
                <span className="label">Пароль</span>
                <input
                  className="input"
                  type="password"
                  value={registerForm.password}
                  onChange={(e) => setRegisterForm({ ...registerForm, password: e.target.value })}
                  placeholder="Минимум 6 символов"
                  autoComplete="new-password"
                />
              </label>

              <label className="field">
                <span className="label">Повтор пароля</span>
                <input
                  className="input"
                  type="password"
                  value={registerForm.password2}
                  onChange={(e) => setRegisterForm({ ...registerForm, password2: e.target.value })}
                  placeholder="Повтори пароль"
                  autoComplete="new-password"
                />
              </label>

              <button className="btn btn-primary" type="submit">
                Зарегистрироваться
              </button>

              <div className="auth-hint">
                <div className="hint-title">Подсказка</div>
                <div className="hint-list">
                  <div>По умолчанию новая учетная запись будет с ролью <b>USER</b>.</div>
                </div>
              </div>
            </form>
          )}
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
