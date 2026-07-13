import { useEffect, useState } from 'react'
import AdminNav from '../components/AdminNav'
import { getAdminEventReports, getAdminSummary } from '../services/reportService'
import { formatPrice } from '../utils/formatters'

function AdminReportsPage() {
  const [summary, setSummary] = useState(null)
  const [reports, setReports] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    async function loadReports() {
      try {
        const [summaryResult, reportsResult] = await Promise.all([
          getAdminSummary(),
          getAdminEventReports(),
        ])
        setSummary(summaryResult.data)
        setReports(reportsResult.data || [])
      } catch {
        setError('No se pudieron cargar los reportes')
      } finally {
        setLoading(false)
      }
    }

    loadReports()
  }, [])

  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Administración</p>
          <h1>Reportes</h1>
        </div>
      </section>
      <AdminNav />

      {loading && <p className="message">Cargando reportes...</p>}
      {error && <p className="message error">{error}</p>}

      {summary && (
        <section className="metric-grid">
          <Metric label="Eventos" value={summary.total_events} />
          <Metric label="Eventos activos" value={summary.active_events} />
          <Metric label="Cancelados" value={summary.cancelled_events} />
          <Metric label="Entradas activas" value={summary.active_tickets} />
          <Metric label="Entradas canceladas" value={summary.cancelled_tickets} />
          <Metric label="Recaudación estimada" value={formatPrice(summary.estimated_revenue)} />
        </section>
      )}

      <div className="table-wrap">
        <table className="admin-table">
          <thead>
            <tr>
              <th>Evento</th>
              <th>Estado</th>
              <th>Emitidas</th>
              <th>Activas</th>
              <th>Canceladas</th>
              <th>Ocupación</th>
              <th>Recaudación</th>
            </tr>
          </thead>
          <tbody>
            {reports.map((report) => (
              <tr key={report.event_id}>
                <td>{report.title}</td>
                <td>{report.status}</td>
                <td>{report.tickets_issued}</td>
                <td>{report.active_tickets}</td>
                <td>{report.cancelled_tickets}</td>
                <td>{Number(report.occupancy_percentage).toFixed(1)}%</td>
                <td>{formatPrice(report.estimated_revenue)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
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

export default AdminReportsPage
