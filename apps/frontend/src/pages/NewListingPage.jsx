import { useEffect, useMemo, useState } from 'react'
import { apiJson } from '../api/client'

export default function NewListingPage({ token, categories, onCreated }) {
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  const firstCategoryId = useMemo(() => {
    if (!Array.isArray(categories) || categories.length === 0) return ''
    return String(categories[0].id)
  }, [categories])

  const [form, setForm] = useState({
    categoryId: '',
    title: '',
    description: '',
    location: '',
    pricePerHour: 0,
  })

  useEffect(() => {
    if (!form.categoryId && firstCategoryId) {
      setForm((f) => ({ ...f, categoryId: firstCategoryId }))
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [firstCategoryId])

  const submit = async (e) => {
    e.preventDefault()
    setError('')
    setOk('')

    const payload = {
      categoryId: Number(form.categoryId),
      title: form.title.trim(),
      description: form.description.trim() ? form.description.trim() : null,
      location: form.location.trim() ? form.location.trim() : null,
      pricePerHour: Number(form.pricePerHour) || 0,
    }

    if (!payload.categoryId) return setError('Выберите категорию')
    if (!payload.title) return setError('Введите название')
    if (payload.pricePerHour < 0) return setError('Цена не может быть отрицательной')

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

      setOk('Объявление опубликовано!')
      setForm((f) => ({ ...f, title: '', description: '', location: '', pricePerHour: 0 }))
      await onCreated?.()
    } catch (e) {
      setError(String(e.message || e))
    }
  }

  return (
    <div>
      <h2 style={{ margin: '8px 0 12px' }}>Разместить объявление</h2>

      <div className="card">
        {error ? <div className="alert-ui">{error}</div> : null}
        {ok ? <div className="ok-ui">{ok}</div> : null}

        <form onSubmit={submit} className="form-col">
          <label className="field-ui">
            <span className="label-ui">Категория</span>
            <select
              className="select-ui"
              value={form.categoryId}
              onChange={(e) => setForm({ ...form, categoryId: e.target.value })}
            >
              {categories.map((c) => (
                <option key={c.id} value={String(c.id)}>{c.name}</option>
              ))}
            </select>
          </label>

          <label className="field-ui">
            <span className="label-ui">Название</span>
            <input
              className="input-ui"
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              placeholder="Например: Переговорная на 8 человек"
            />
          </label>

          <label className="field-ui">
            <span className="label-ui">Локация</span>
            <input
              className="input-ui"
              value={form.location}
              onChange={(e) => setForm({ ...form, location: e.target.value })}
              placeholder="Например: Иркутск, 1 этаж"
            />
          </label>

          <label className="field-ui">
            <span className="label-ui">Цена (₽/час)</span>
            <input
              className="input-ui"
              type="number"
              min="0"
              value={form.pricePerHour}
              onChange={(e) => setForm({ ...form, pricePerHour: e.target.value })}
            />
          </label>

          <label className="field-ui">
            <span className="label-ui">Описание</span>
            <textarea
              className="textarea-ui"
              rows={4}
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              placeholder="Коротко опиши, что входит, правила, время работы…"
            />
          </label>

          <button className="btn-ui" type="submit">
            Опубликовать
          </button>
        </form>
      </div>
    </div>
  )
}
