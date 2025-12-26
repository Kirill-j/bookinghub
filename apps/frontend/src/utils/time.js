export function pad2(n) {
  return String(n).padStart(2, '0')
}

export function toTimeHHMM(value) {
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return String(value)
  return `${pad2(d.getHours())}:${pad2(d.getMinutes())}`
}

export function hhmmToMinutes(hhmm) {
  const [h, m] = String(hhmm).split(':').map(Number)
  return h * 60 + m
}

export function minutesToHHMM(min) {
  const h = Math.floor(min / 60)
  const m = min % 60
  return `${pad2(h)}:${pad2(m)}`
}

export function overlaps(aStart, aEnd, bStart, bEnd) {
  return aStart < bEnd && aEnd > bStart
}
