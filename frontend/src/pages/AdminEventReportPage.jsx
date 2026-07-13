import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import AdminNav from '../components/AdminNav'
import { getEventReport } from '../services/adminEventService'
import { formatPrice } from '../utils/formatters'

function AdminEventReportPage() {
  const { id } = useParams()
  const [report, setReport] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function loadReport() {
      try {
        const result = await getEventReport(id)
        setReport(result.data)
      } catch {
        setError('No se pudo cargar el reporte del evento')
      } finally {
        setLoading(false)
      }
    }

    loadReport()
  }, [id])

  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Administración</p>
          <h1>Reporte del evento</h1>
        </div>
      </section>
      <AdminNav />
      {loading && <p className="message">Cargando reporte...</p>}
      {error && <p className="message error">{error}</p>}
      {report && (
        <section className="metric-grid">
          <Metric label="Evento" value={report.title} />
          <Metric label="Estado" value={report.status} />
          <Metric label="Capacidad" value={report.capacity} />
          <Metric label="Disponibles" value={report.available_capacity} />
          <Metric label="Entradas activas" value={report.active_tickets} />
          <Metric label="Ocupación" value={`${Number(report.occupancy_percentage).toFixed(1)}%`} />
          <Metric label="Recaudación" value={formatPrice(report.estimated_revenue)} />
        </section>
      )}
    </main>
  )
}

function Metric({ label, value }) {
  return (
    <article className="metric-card">
      <span>{label}</span>
      <strong>{value}</strong>
    </article>
  )
}

export default AdminEventReportPage
