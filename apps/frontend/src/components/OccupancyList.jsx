import { toTimeHHMM } from '../utils/time'

export default function OccupancyList({ resourceBookings }) {
  return (
    <div style={{ marginTop: 10, paddingTop: 10, borderTop: '1px dashed #bbb' }}>
      <b>Занятость на выбранную дату:</b>

      {(resourceBookings?.length || 0) === 0 ? (
        <div style={{ opacity: 0.75 }}>Свободно (броней нет)</div>
      ) : (
        <ul>
          {resourceBookings.map((b) => (
            <li key={b.id}>
              {toTimeHHMM(b.startAt)}–{toTimeHHMM(b.endAt)} — {b.status}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
