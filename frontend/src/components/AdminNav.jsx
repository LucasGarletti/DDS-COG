import { NavLink } from 'react-router-dom'

function AdminNav() {
  return (
    <nav className="admin-nav" aria-label="Navegacion administrativa">
      <NavLink to="/admin">Inicio</NavLink>
      <NavLink to="/admin/eventos">Eventos</NavLink>
      <NavLink to="/admin/eventos/nuevo">Crear evento</NavLink>
      <NavLink to="/admin/reportes">Reportes</NavLink>
    </nav>
  )
}

export default AdminNav
