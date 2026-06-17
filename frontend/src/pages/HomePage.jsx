import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getEvents } from '../services/eventService'
import { getEventImage } from '../utils/eventImages'

function HomePage() {
  const [events, setEvents] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function loadFeaturedEvents() {
      try {
        const result = await getEvents()
        setEvents((result.data || []).slice(0, 3))
      } catch {
        setError('No se pudieron cargar los eventos destacados')
      } finally {
        setLoading(false)
      }
    }

    loadFeaturedEvents()
  }, [])

  return (
    <main className="page home-page">
      <section className="home-hero">
        <div>
          <p>Tickets simples, eventos inolvidables</p>
          <h1>Vivi tus eventos favoritos</h1>
          <span>
            Encontra recitales, festivales y experiencias deportivas en un solo
            lugar.
          </span>
          <Link className="primary-button" to="/eventos">
            Ver eventos
          </Link>
        </div>
      </section>

      <section className="featured-section">
        <div className="section-heading">
          <p>Agenda TickGo</p>
          <h2>Eventos destacados</h2>
        </div>

        {loading && <p className="message">Cargando eventos destacados...</p>}
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
                <h2>{event.title}</h2>
                <p className="event-location">{event.location}</p>
                <Link className="primary-button" to={`/eventos/${event.id}`}>
                  Ver detalle
                </Link>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </main>
  )
}

export default HomePage
