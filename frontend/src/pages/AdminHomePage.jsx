import { Link } from 'react-router-dom'
import AdminNav from '../components/AdminNav'

function AdminHomePage() {
  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Administración</p>
          <h1>Panel administrativo</h1>
        </div>
      </section>
      <AdminNav />
      <section className="admin-grid">
        <Link className="admin-card" to="/admin/eventos">
          <h2>Eventos</h2>
          <p>Editar, cancelar y revisar eventos publicados.</p>
        </Link>
        <Link className="admin-card" to="/admin/eventos/nuevo">
          <h2>Crear evento</h2>
          <p>Cargar nuevos shows, festivales y experiencias.</p>
        </Link>
        <Link className="admin-card" to="/admin/reportes">
          <h2>Reportes</h2>
          <p>Ver métricas generales y rendimiento por evento.</p>
        </Link>
      </section>
    </main>
  )
}

export default AdminHomePage
