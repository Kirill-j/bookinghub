export default function ResourceCreateForm({
  categories,
  resourceForm,
  setResourceForm,
  onCreateResource,
  canCreateResources,
}) {
  if (!canCreateResources) {
    return (
      <div style={{ padding: 12, border: '1px dashed #bbb', marginBottom: 16 }}>
        <b>Создание ресурсов доступно только менеджеру или администратору.</b>
        <div style={{ fontSize: 12, opacity: 0.75 }}>
          Войдите как manager@bookinghub.local или admin@bookinghub.local.
        </div>
      </div>
    )
  }

  return (
    <>
      <h2>Создать ресурс</h2>
      <form onSubmit={onCreateResource} style={{ display: 'grid', gap: 8, marginBottom: 16 }}>
        <label>
          Категория:
          <select
            value={resourceForm.categoryId}
            onChange={(e) => setResourceForm({ ...resourceForm, categoryId: e.target.value })}
            style={{ marginLeft: 8 }}
          >
            {categories.map((c) => (
              <option key={c.id} value={String(c.id)}>
                {c.name}
              </option>
            ))}
          </select>
        </label>

        <label>
          Название:
          <input
            value={resourceForm.title}
            onChange={(e) => setResourceForm({ ...resourceForm, title: e.target.value })}
          />
        </label>

        <label>
          Описание:
          <input
            value={resourceForm.description}
            onChange={(e) => setResourceForm({ ...resourceForm, description: e.target.value })}
          />
        </label>

        <label>
          Место/локация:
          <input
            value={resourceForm.location}
            onChange={(e) => setResourceForm({ ...resourceForm, location: e.target.value })}
          />
        </label>

        <button type="submit">Создать</button>
      </form>
    </>
  )
}
