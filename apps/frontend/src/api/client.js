// Проверяем, что Vite видит переменную. 
// Если нет — выводим в консоль для отладки (потом удалим)
console.log("API URL:", import.meta.env.VITE_API_URL);

const BASE_URL = import.meta.env.VITE_API_URL || '';

export function getToken() {
  return localStorage.getItem('accessToken') || ''
}

export function saveToken(token) {
  if (token) localStorage.setItem('accessToken', token)
  else localStorage.removeItem('accessToken')
}

export async function apiText(path, token = '') {
  // Добавляем BASE_URL перед путем
  const r = await fetch(`${BASE_URL}${path}`, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  })
  if (!r.ok) throw new Error(await r.text())
  return r.text()
}

export async function apiJson(path, opts = {}, token = '') {
  const headers = { ...(opts.headers || {}) }
  if (token) headers.Authorization = `Bearer ${token}`

  // Добавляем BASE_URL перед путем
  const r = await fetch(`${BASE_URL}${path}`, { ...opts, headers })
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}