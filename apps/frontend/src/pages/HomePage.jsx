import { useMemo, useState } from 'react'
import Select from '../components/ui/Select'

function rub(n) {
  const v = Number(n) || 0
  return new Intl.NumberFormat('ru-RU').format(v)
}

export default function HomePage({ categories, resources, onOpenResource }) {
  const [q, setQ] = useState('')
  const [catId, setCatId] = useState('all')
  const [sort, setSort] = useState('popular')
  const [minPrice, setMinPrice] = useState('')
  const [maxPrice, setMaxPrice] = useState('')

  const filtered = useMemo(() => {
    let list = Array.isArray(resources) ? resources : []

    const query = q.trim().toLowerCase()
    if (query) {
      list = list.filter((r) =>
        String(r.title || '').toLowerCase().includes(query) ||
        String(r.location || '').toLowerCase().includes(query)
      )
    }

    if (catId !== 'all') {
      list = list.filter((r) => String(r.categoryId) === String(catId))
    }

    const min = minPrice === '' ? null : Number(minPrice)
    const max = maxPrice === '' ? null : Number(maxPrice)
    if (min !== null && !Number.isNaN(min)) {
      list = list.filter((r) => (Number(r.pricePerHour) || 0) >= min)
    }
    if (max !== null && !Number.isNaN(max)) {
      list = list.filter((r) => (Number(r.pricePerHour) || 0) <= max)
    }

    if (sort === 'price_asc') {
      list = [...list].sort((a, b) => (Number(a.pricePerHour) || 0) - (Number(b.pricePerHour) || 0))
    } else if (sort === 'price_desc') {
      list = [...list].sort((a, b) => (Number(b.pricePerHour) || 0) - (Number(a.pricePerHour) || 0))
    } else {
      list = [...list]
    }

    return list
  }, [resources, q, catId, sort, minPrice, maxPrice])

  const categoryOptions = useMemo(() => {
  const base = [{ value: 'all', label: 'Все категории' }]
  const rest = (Array.isArray(categories) ? categories : []).map((c) => ({
    value: String(c.id),
    label: c.name,
  }))
  return base.concat(rest)
}, [categories])

const sortOptions = useMemo(() => ([
  { value: 'popular', label: 'Сначала актуальные' },
  { value: 'price_asc', label: 'Цена: по возрастанию' },
  { value: 'price_desc', label: 'Цена: по убыванию' },
]), [])


  return (
    <div>
      <div className="catalog-head">
        <h2 style={{ margin: 0 }}>Каталог</h2>
        <div className="muted">Выбери ресурс, посмотри занятость и забронируй</div>
      </div>

      <div className="filters">
        <input
          className="input-ui"
          placeholder="Поиск по названию или локации…"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />

        <Select
          value={catId}
          onChange={(v) => setCatId(String(v))}
          options={categoryOptions}
        />

        <Select
          value={sort}
          onChange={(v) => setSort(String(v))}
          options={sortOptions}
        />

        <input
          className="input-ui"
          type="number"
          placeholder="Цена от"
          value={minPrice}
          onChange={(e) => setMinPrice(e.target.value)}
        />
        <input
          className="input-ui"
          type="number"
          placeholder="Цена до"
          value={maxPrice}
          onChange={(e) => setMaxPrice(e.target.value)}
        />
      </div>

      <div className="cards-grid">
        {filtered.map((r) => (
          <button
            key={r.id}
            className="resource-card"
            type="button"
            onClick={() => onOpenResource(r.id)}
          >
            <div className="card-title">{r.title}</div>
            <div className="card-sub">{r.location || 'Локация не указана'}</div>
            <div className="card-price">{rub(r.pricePerHour)} ₽/час</div>
          </button>
        ))}

        {filtered.length === 0 && (
          <div className="muted" style={{ padding: 12 }}>
            Ничего не найдено.
          </div>
        )}
      </div>
    </div>
  )
}
