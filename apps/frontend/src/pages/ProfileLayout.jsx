import { NavLink, Outlet } from 'react-router-dom'

export default function ProfileLayout() {
  return (
    <div>
      <div className="profile-grid">
        <div className="profile-menu card">
          <div className="profile-menu-title">Разделы</div>

          <NavLink className="profile-link" to="." end>Профиль</NavLink>
          <NavLink className="profile-link" to="listings">Мои объявления</NavLink>
          <NavLink className="profile-link" to="pending">Подтверждение брони</NavLink>
          <NavLink className="profile-link" to="bookings">Мои бронирования</NavLink>
          <NavLink className="profile-link" to="new">Разместить объявление</NavLink>
        </div>

        <div className="profile-content">
          <Outlet />
        </div>
      </div>
    </div>
  )
}
