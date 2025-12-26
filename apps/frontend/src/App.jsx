import { useEffect, useState } from 'react'
import { BrowserRouter, Routes, Route, useNavigate, useParams } from 'react-router-dom'

import { apiJson, apiText, getToken, saveToken } from './api/client'

import AuthPage from './components/AuthPage'
import AppShell from './components/AppShell'
import HomePage from './pages/HomePage'
import ResourcePage from './pages/ResourcePage'
import NewListingPage from './pages/NewListingPage'
import ProfilePage from './pages/ProfilePage'


function HomeRoute(props) {
  const navigate = useNavigate()
  return <HomePage {...props} onOpenResource={(id) => navigate(`/resources/${id}`)} />
}

function ResourceRoute(props) {
  const { id } = useParams()
  const navigate = useNavigate()
  return <ResourcePage {...props} id={id} onBack={() => navigate('/')} />
}

export default function App() {
  const [serverStatus, setServerStatus] = useState('загрузка...')
  const [error, setError] = useState('')

  const [token, setTokenState] = useState(getToken())
  const [me, setMe] = useState(null)

  const [categories, setCategories] = useState([])
  const [resources, setResources] = useState([])

  const [loginForm, setLoginForm] = useState({
    email: 'manager@bookinghub.local',
    password: '123456',
  })

  const [authMode, setAuthMode] = useState('login') // 'login' | 'register'

  const [registerForm, setRegisterForm] = useState({
    name: '',
    email: '',
    password: '',
    password2: '',
    accountType: 'INDIVIDUAL',
  })

  // ---- loaders ----
  const loadPublic = async () => {
    const [cats, res] = await Promise.all([
      apiJson('/api/categories', {}, token),
      apiJson('/api/resources', {}, token),
    ])

    setCategories(Array.isArray(cats) ? cats : [])
    setResources(Array.isArray(res) ? res : [])
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

  // ---- effects ----
  useEffect(() => {
    apiText('/api/health')
      .then(() => setServerStatus('ок'))
      .catch(() => setServerStatus('ошибка'))

    loadPublic().catch((e) => setError(String(e.message || e)))
    loadMe(token).catch(() => {})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

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

      if (data.user) setMe(data.user)

      await loadMe(t)
      await loadPublic(t)
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const onRegister = async (e) => {
    e.preventDefault()
    setError('')

    const name = registerForm.name.trim()
    const email = registerForm.email.trim().toLowerCase()
    const password = registerForm.password
    const password2 = registerForm.password2

    if (!name) return setError('Введите имя')
    if (!email) return setError('Введите email')
    if (!password || password.length < 6) return setError('Пароль должен быть минимум 6 символов')
    if (password !== password2) return setError('Пароли не совпадают')

    try {
      // 1) регистрация
      await apiJson(
        '/api/auth/register',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ name, email, password, accountType: registerForm.accountType }),
        },
        ''
      )

      // 2) авто-логин после регистрации
      const data = await apiJson(
        '/api/auth/login',
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email, password }),
        },
        ''
      )

      const t = data.accessToken
      saveToken(t)
      setTokenState(t)

      if (data.user) setMe(data.user)

      await loadMe(t)
      await loadPublic()

      // очистим форму
      setRegisterForm({ name: '', email: '', password: '', password2: '' })
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  const onLogout = async () => {
    saveToken('')
    setTokenState('')
    setMe(null)
  }

  if (!me) {
    return (
      <AuthPage
        serverStatus={serverStatus}
        error={error}
        mode={authMode}
        setMode={setAuthMode}
        loginForm={loginForm}
        setLoginForm={setLoginForm}
        onLogin={onLogin}
        registerForm={registerForm}
        setRegisterForm={setRegisterForm}
        onRegister={onRegister}
      />
    )
  }

  return (
    <BrowserRouter>
      <AppShell me={me} onLogout={onLogout}>
        <Routes>
          <Route
            path="/"
            element={
              <HomeRoute
                categories={categories}
                resources={resources}
              />
            }
          />

          <Route
            path="/resources/:id"
            element={
              <ResourceRoute
                token={token}
                me={me}
                resources={resources}
                onRefreshAfterBooking={async () => {
                }}
              />
            }
          />

          <Route
            path="/new"
            element={
              <NewListingPage
                token={token}
                categories={categories}
                onCreated={async () => {
                  await loadPublic()
                }}
              />
            }
          />

          <Route
            path="/profile"
            element={
              <ProfilePage
                token={token}
                me={me}
                categories={categories}
              />
            }
          />
        </Routes>
      </AppShell>
    </BrowserRouter>
  )
}
