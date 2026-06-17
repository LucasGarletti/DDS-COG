import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import fallbackEventImage from '../assets/hero.png'
import { getEventById } from '../services/eventService'
import { purchaseTicket } from '../services/ticketService'

function formatDate(date) {
  if (!date) {
    return 'Fecha a confirmar'
  }

  return new Intl.DateTimeFormat('es-AR', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
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

function EventDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [event, setEvent] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [purchaseMessage, setPurchaseMessage] = useState('')
  const [purchaseError, setPurchaseError] = useState('')
  const [purchasing, setPurchasing] = useState(false)

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
            src={event.image_url || fallbackEventImage}
            alt={event.title}
            className="detail-image"
          />

          <section className="detail-panel">
            <p className="event-date">{formatDate(event.date)}</p>
            <h1>{event.title}</h1>
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
    </main>
  )
}

export default EventDetailPage
