export function getToken() {
  return localStorage.getItem('accessToken') || ''
}

export function saveToken(token) {
  if (token) localStorage.setItem('accessToken', token)
  else localStorage.removeItem('accessToken')
}

export async function apiText(path, token = '') {
  const r = await fetch(path, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  })
  if (!r.ok) throw new Error(await r.text())
  return r.text()
}

export async function apiJson(path, opts = {}, token = '') {
  const headers = { ...(opts.headers || {}) }
  if (token) headers.Authorization = `Bearer ${token}`

  const r = await fetch(path, { ...opts, headers })
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}
