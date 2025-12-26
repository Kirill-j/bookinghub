import { useMemo, useState } from 'react'

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

        <select className="select-ui" value={catId} onChange={(e) => setCatId(e.target.value)}>
          <option value="all">Все категории</option>
          {categories.map((c) => (
            <option key={c.id} value={String(c.id)}>{c.name}</option>
          ))}
        </select>

        <select className="select-ui" value={sort} onChange={(e) => setSort(e.target.value)}>
          <option value="popular">Сначала актуальные</option>
          <option value="price_asc">Цена: по возрастанию</option>
          <option value="price_desc">Цена: по убыванию</option>
        </select>

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
