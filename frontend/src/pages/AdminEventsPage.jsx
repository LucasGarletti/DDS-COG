import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import AdminNav from '../components/AdminNav'
import { cancelEvent } from '../services/adminEventService'
import { getEvents } from '../services/eventService'
import { getAdminEventReports } from '../services/reportService'
import { formatDate, formatPrice } from '../utils/formatters'

function AdminEventsPage() {
  const [events, setEvents] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [actionId, setActionId] = useState(null)
  const [festivalEventIds, setFestivalEventIds] = useState(new Set())
  const [publicEventsById, setPublicEventsById] = useState({})

  async function loadEvents() {
    setLoading(true)
    setError('')
    try {
      const [reportsResult, publicEventsResult] = await Promise.all([
        getAdminEventReports(),
        getEvents(),
      ])
      setEvents(reportsResult.data || [])
      setFestivalEventIds(
        new Set(
          (publicEventsResult.data || [])
            .filter((event) => event.is_festival)
            .map((event) => Number(event.id)),
        ),
      )
      setPublicEventsById(
        Object.fromEntries(
          (publicEventsResult.data || []).map((event) => [Number(event.id), event]),
        ),
      )
    } catch {
      setError('No se pudieron cargar los eventos administrativos')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    async function loadInitialEvents() {
      setLoading(true)
      setError('')
      try {
        const [reportsResult, publicEventsResult] = await Promise.all([
          getAdminEventReports(),
          getEvents(),
        ])
        setEvents(reportsResult.data || [])
        setFestivalEventIds(
          new Set(
            (publicEventsResult.data || [])
              .filter((event) => event.is_festival)
              .map((event) => Number(event.id)),
          ),
        )
      } catch {
        setError('No se pudieron cargar los eventos administrativos')
      } finally {
        setLoading(false)
      }
    }

    loadInitialEvents()
  }, [])

  async function handleCancel(eventId) {
    if (!window.confirm('¿Cancelar este evento?')) {
      return
    }

    setActionId(eventId)
    setMessage('')
    setError('')

    try {
      await cancelEvent(eventId)
      setMessage('Evento cancelado correctamente')
      await loadEvents()
    } catch {
      setError('No se pudo cancelar el evento')
    } finally {
      setActionId(null)
    }
  }

  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Administración</p>
          <h1>Eventos</h1>
        </div>
      </section>
      <AdminNav />

      {loading && <p className="message">Cargando eventos...</p>}
      {message && <p className="message success">{message}</p>}
      {error && <p className="message error">{error}</p>}

      <div className="table-wrap">
        <table className="admin-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Título</th>
              <th>Fecha</th>
              <th>Ubicación</th>
              <th>Precio</th>
              <th>Estado</th>
              <th>Festival</th>
              <th>Capacidad</th>
              <th>Disponibles</th>
              <th>Recaudación estimada</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            {events.map((event) => {
              const publicEvent = publicEventsById[Number(event.event_id)]
              const isFestival = festivalEventIds.has(Number(event.event_id))

              return (
                <tr key={event.event_id}>
                  <td>{event.event_id}</td>
                  <td>
                    <strong>{event.title}</strong>
                    {event.status === 'active' && (
                      <span className="status-badge">Activo</span>
                    )}
                  </td>
                  <td>{publicEvent ? formatDate(publicEvent.date) : 'No disponible'}</td>
                  <td>{publicEvent?.location || 'No disponible'}</td>
                  <td>
                    {publicEvent ? formatPrice(publicEvent.price) : 'No disponible'}
                  </td>
                  <td>{event.status}</td>
                  <td>{isFestival ? 'Sí' : 'No'}</td>
                  <td>{event.capacity}</td>
                  <td>{event.available_capacity}</td>
                  <td>{formatPrice(event.estimated_revenue)}</td>
                  <td className="table-actions">
                    <Link to={`/admin/eventos/${event.event_id}/editar`}>Editar</Link>
                    {isFestival && (
                      <Link to={`/admin/eventos/${event.event_id}/grilla`}>Grilla</Link>
                    )}
                    <Link to={`/admin/eventos/${event.event_id}/reporte`}>Reporte</Link>
                    {event.status !== 'cancelled' && (
                      <button
                        type="button"
                        onClick={() => handleCancel(event.event_id)}
                        disabled={actionId === event.event_id}
                      >
                        Cancelar
                      </button>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </main>
  )
}

export default AdminEventsPage
