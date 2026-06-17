import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import cosquiRockImage from '../assets/cosquirock.png'
import pastillasImage from '../assets/pastillasdelabuelo.jfif'
import pumasImage from '../assets/pumas.jfif'
import { getEvents } from '../services/eventService'

const fallbackEventImage =
  'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 450"%3E%3Crect width="800" height="450" fill="%23006fff"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle" font-family="Arial" font-size="56" font-weight="700" fill="white"%3ETickGo%3C/text%3E%3C/svg%3E'

function getEventImage(event) {
  const title = event.title || ''

  if (title === 'Festival Cosquín Rock') {
    return cosquiRockImage
  }

  if (title === 'Los Pumas en el Estadio Mario Alberto Kempes') {
    return pumasImage
  }

  if (title === 'Las Pastillas del Abuelo en Plaza de la Música') {
    return pastillasImage
  }

  return event.image_url || fallbackEventImage
}

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
          <h1>Elegí tu próximo evento</h1>
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
