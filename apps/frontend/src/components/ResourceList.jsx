export default function ResourceList({ resources, categoryNameById }) {
  return (
    <>
      <h2>Ресурсы</h2>
      <ul>
        {resources.map((r) => (
          <li key={r.id}>
            <b>{r.title}</b>{' '}
            <span style={{ opacity: 0.7 }}>
              ({categoryNameById.get(String(r.categoryId)) || `категория #${r.categoryId}`})
            </span>
            {r.location ? ` — ${r.location}` : ''}
          </li>
        ))}
      </ul>
    </>
  )
}
