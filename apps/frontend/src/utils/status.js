export const BOOKING_STATUS_RU = {
  PENDING: 'В ожидании',
  APPROVED: 'Подтверждено',
  REJECTED: 'Отклонено',
  CANCELED: 'Отменено',
}

export function statusRu(s) {
  const key = String(s || '').trim().toUpperCase()
  return BOOKING_STATUS_RU[key] || s
}

// export function statusRu(status) {
//   switch (status) {
//     case 'PENDING': return 'В ожидании'
//     case 'APPROVED': return 'Подтверждено'
//     case 'REJECTED': return 'Отклонено'
//     case 'CANCELED': return 'Отменено'
//     default: return status
//   }
// }

export function statusClass(status) {
  switch (status) {
    case 'PENDING': return 'status-pending'
    case 'APPROVED': return 'status-approved'
    case 'REJECTED': return 'status-rejected'
    case 'CANCELED': return 'status-canceled'
    default: return ''
  }
}