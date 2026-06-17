import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getEvents } from '../services/eventService'
import { getEventImage } from '../utils/eventImages'

function formatDate(date) {
  if (!date) {
    return 'Fecha a confirmar'
  }

  return new Intl.DateTimeFormat('es-AR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(new Date(date))
}

function formatPrice(price) {
  if (price === undefined || price === null) {
    return 'Precio a confirmar'
  }

  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
  }).format(price)
}

function EventsPage() {
  const [events, setEvents] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function loadEvents() {
      try {
        const result = await getEvents()
        setEvents(result.data || [])
      } catch {
        setError('No se pudieron cargar los eventos')
      } finally {
        setLoading(false)
      }
    }

    loadEvents()
  }, [])

  return (
    <main className="page">
      <section className="page-hero">
        <div>
          <p>Entradas para shows y experiencias</p>
          <h1>Elegi tu proximo evento</h1>
        </div>
      </section>

      {loading && <p className="message">Cargando eventos...</p>}
      {error && <p className="message error">{error}</p>}

      <ul className="event-list">
        {events.map((event) => (
          <li className="event-card" key={event.id}>
            <img
              src={getEventImage(event)}
              alt={event.title}
              className="event-image"
            />
            <div className="event-card-body">
              <p className="event-date">{formatDate(event.date)}</p>
              <h2>{event.title}</h2>
              <p className="event-location">{event.location}</p>
              <div className="event-meta">
                <span>{formatPrice(event.price)}</span>
                <span>{event.available_capacity} disponibles</span>
              </div>
              <Link className="primary-button" to={`/eventos/${event.id}`}>
                Ver detalle
              </Link>
            </div>
          </li>
        ))}
      </ul>
    </main>
  )
}

export default EventsPage
