export default function AuthPanel({
  me,
  loginForm,
  setLoginForm,
  onLogin,
  onLogout,
}) {
  return (
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
            Тестовые аккаунты: admin@bookinghub.local / manager@bookinghub.local / user@bookinghub.local, пароль 123456
          </div>
        </form>
      )}
    </div>
  )
}
