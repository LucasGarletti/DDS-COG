import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import EventCard from '../components/EventCard'
import { getEvents } from '../services/eventService'

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
            <EventCard event={event} key={event.id} />
          ))}
        </ul>

        <Link className="secondary-button" to="/eventos">
          Ver todos
        </Link>
      </section>
    </main>
  )
}

export default HomePage
