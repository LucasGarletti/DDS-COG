import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { getEventById } from '../services/eventService'
import { getPublicSchedule } from '../services/festivalService'
import { getErrorMessage } from '../services/httpClient'
import { purchaseTicket } from '../services/ticketService'
import { getEventImage } from '../utils/eventImages'
import { formatDateTime, formatPrice, formatTime } from '../utils/formatters'

function EventDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [event, setEvent] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [purchaseMessage, setPurchaseMessage] = useState('')
  const [purchaseError, setPurchaseError] = useState('')
  const [purchasing, setPurchasing] = useState(false)
  const [schedule, setSchedule] = useState([])
  const [scheduleError, setScheduleError] = useState('')
  const [scheduleLoading, setScheduleLoading] = useState(false)

  useEffect(() => {
    async function loadEvent() {
      try {
        const result = await getEventById(id)
        setEvent(result.data)
      } catch {
        setError('No se pudo cargar el evento')
      } finally {
        setLoading(false)
      }
    }

    loadEvent()
  }, [id])

  useEffect(() => {
    if (!event?.is_festival) {
      return
    }

    async function loadSchedule() {
      setScheduleError('')
      setScheduleLoading(true)

      try {
        const result = await getPublicSchedule(id)
        setSchedule(result.data || [])
      } catch (requestError) {
        setScheduleError(getErrorMessage(requestError, 'No se pudo cargar la grilla'))
      } finally {
        setScheduleLoading(false)
      }
    }

    loadSchedule()
  }, [event?.is_festival, id])

  async function handlePurchase() {
    setPurchaseMessage('')
    setPurchaseError('')

    if (!localStorage.getItem('token')) {
      navigate('/login')
      return
    }

    setPurchasing(true)

    try {
      await purchaseTicket(id)
      setPurchaseMessage('Entrada comprada correctamente')

      const result = await getEventById(id)
      setEvent(result.data)
    } catch {
      setPurchaseError('No se pudo comprar la entrada')
    } finally {
      setPurchasing(false)
    }
  }

  return (
    <main className="page">
      {loading && <p className="message">Cargando evento...</p>}
      {error && <p className="message error">{error}</p>}

      {event && (
        <article className="detail-layout">
          <img
            src={getEventImage(event)}
            alt={event.title}
            className="detail-image"
          />

          <section className="detail-panel">
            <p className="event-date">{formatDateTime(event.date)}</p>
            <h1>{event.title}</h1>
            {event.is_festival && <span className="status-badge">Festival</span>}
            <p className="detail-description">{event.description}</p>

            <dl className="detail-list">
              <div>
                <dt>Lugar</dt>
                <dd>{event.location}</dd>
              </div>
              <div>
                <dt>Disponibilidad</dt>
                <dd>{event.available_capacity} entradas</dd>
              </div>
              <div>
                <dt>Precio</dt>
                <dd>{formatPrice(event.price)}</dd>
              </div>
            </dl>

            {purchaseMessage && (
              <p className="message success">{purchaseMessage}</p>
            )}
            {purchaseError && <p className="message error">{purchaseError}</p>}

            <button
              className="primary-button"
              type="button"
              onClick={handlePurchase}
              disabled={purchasing}
            >
              {purchasing ? 'Comprando...' : 'Comprar entrada'}
            </button>
          </section>
        </article>
      )}

      {event?.is_festival && (
        <section className="festival-section">
          <div className="section-heading">
            <p>Festival</p>
            <h2>Grilla oficial</h2>
          </div>

          {scheduleLoading && <p className="message">Cargando grilla...</p>}
          {scheduleError && <p className="message error">{scheduleError}</p>}

          <div className="schedule-list">
            {schedule.map((show) => (
              <article className="schedule-card" key={show.id}>
                {show.image_url && <img src={show.image_url} alt={show.artist} />}
                <div>
                  <p className="event-date">
                    {formatTime(show.start_time)} - {formatTime(show.end_time)}
                  </p>
                  <h3>{show.artist}</h3>
                  <p>{show.stage}</p>
                </div>
              </article>
            ))}
          </div>

          {schedule.length === 0 && !scheduleLoading && (
            <p className="message">Todavía no hay shows cargados.</p>
          )}
        </section>
      )}
    </main>
  )
}

export default EventDetailPage
