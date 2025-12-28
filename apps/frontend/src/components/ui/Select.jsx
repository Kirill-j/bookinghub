import { useEffect, useMemo, useRef, useState } from 'react'

export default function Select({
  value,
  onChange,
  options, // [{ value: '1', label: 'Все категории' }]
  placeholder = 'Выберите…',
}) {
  const [open, setOpen] = useState(false)
  const rootRef = useRef(null)

  const currentLabel = useMemo(() => {
    const f = options?.find((o) => String(o.value) === String(value))
    return f ? f.label : ''
  }, [options, value])

  useEffect(() => {
    const onDoc = (e) => {
      if (!rootRef.current) return
      if (!rootRef.current.contains(e.target)) setOpen(false)
    }
    document.addEventListener('mousedown', onDoc)
    return () => document.removeEventListener('mousedown', onDoc)
  }, [])

  return (
    <div ref={rootRef} className="selectx">
      <button
        type="button"
        className="selectx-btn"
        onClick={() => setOpen((v) => !v)}
      >
        <span className={`selectx-value ${currentLabel ? '' : 'muted'}`}>
          {currentLabel || placeholder}
        </span>
        <span className={`selectx-caret ${open ? 'selectx-caret-up' : ''}`}>▾</span>
      </button>

      {open ? (
        <div className="selectx-menu">
          {options?.map((o) => {
            const active = String(o.value) === String(value)
            return (
              <button
                key={String(o.value)}
                type="button"
                className={`selectx-item ${active ? 'selectx-item-active' : ''}`}
                onClick={() => {
                  onChange?.(o.value)
                  setOpen(false)
                }}
              >
                {o.label}
              </button>
            )
          })}
        </div>
      ) : null}
    </div>
  )
}
