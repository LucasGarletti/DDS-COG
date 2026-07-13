import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  cancelTicket,
  getMyTickets,
  transferTicket,
} from '../services/ticketService'
import { getItinerary } from '../services/itineraryService'

function formatDate(date) {
  if (!date) {
    return 'Fecha no disponible'
  }

  return new Intl.DateTimeFormat('es-AR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(new Date(date))
}

function getUserIDFromToken(token) {
  if (!token) {
    return null
  }

  try {
    const payload = token.split('.')[1]
    const normalizedPayload = payload.replace(/-/g, '+').replace(/_/g, '/')
    const decodedPayload = JSON.parse(atob(normalizedPayload))

    return decodedPayload.id
  } catch {
    return null
  }
}

function filterTicketsByUser(tickets, userID) {
  if (!userID) {
    return tickets
  }

  return tickets.filter((ticket) => {
    if (ticket.user_id === undefined || ticket.user_id === null) {
      return true
    }

    return Number(ticket.user_id) === Number(userID)
  })
}

function MyTicketsPage({ authToken }) {
  const [tickets, setTickets] = useState([])
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(true)
  const [actionLoadingId, setActionLoadingId] = useState(null)
  const [transferEmails, setTransferEmails] = useState({})
  const [itineraryItemCounts, setItineraryItemCounts] = useState({})

  async function loadTickets() {
    const currentToken = localStorage.getItem('token')
    const currentUserID = getUserIDFromToken(currentToken)

    setTickets([])
    setItineraryItemCounts({})
    setError('')
    setLoading(true)

    try {
      const result = await getMyTickets()
      if (currentToken !== localStorage.getItem('token')) {
        return
      }

      const userTickets = filterTicketsByUser(result.data || [], currentUserID)
      const counts = await loadItineraryItemCounts(userTickets)
      if (currentToken !== localStorage.getItem('token')) {
        return
      }
      setTickets(userTickets)
      setItineraryItemCounts(counts)
    } catch {
      setError('No se pudieron cargar tus entradas')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    async function resetAndLoadTickets() {
      setTickets([])
      setTransferEmails({})
      setItineraryItemCounts({})
      setMessage('')
      setError('')
      await loadTickets()
    }

    resetAndLoadTickets()
  }, [authToken])

  async function handleCancel(ticketId) {
    setError('')
    setMessage('')
    setActionLoadingId(ticketId)

    try {
      await cancelTicket(ticketId)
      setMessage('Entrada cancelada correctamente')
      await loadTickets()
    } catch {
      setError('No se pudo cancelar la entrada')
    } finally {
      setActionLoadingId(null)
    }
  }

  async function handleTransfer(ticketId) {
    const recipientEmail = (transferEmails[ticketId] || '').trim()
    if (!recipientEmail) {
      setError('Ingresá el email del destinatario')
      return
    }

    setError('')
    setMessage('')
    setActionLoadingId(ticketId)

    try {
      await transferTicket(ticketId, recipientEmail)
      setMessage('Entrada transferida correctamente')
      setTransferEmails((currentEmails) => ({
        ...currentEmails,
        [ticketId]: '',
      }))
      await loadTickets()
    } catch {
      setError('No se pudo transferir la entrada')
    } finally {
      setActionLoadingId(null)
    }
  }

  function handleTransferEmailChange(ticketId, value) {
    setTransferEmails((currentEmails) => ({
      ...currentEmails,
      [ticketId]: value,
    }))
  }

  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Tu cuenta</p>
          <h1>Mis entradas</h1>
        </div>
      </section>

      {loading && <p className="message">Cargando entradas...</p>}
      {message && <p className="message success">{message}</p>}
      {error && <p className="message error">{error}</p>}

      <ul className="ticket-list">
        {tickets.map((ticket) => (
          <li className="ticket-card" key={ticket.id}>
            <div className="ticket-info">
              <div>
                <p className="ticket-code">{ticket.code}</p>
                <h2>{ticket.event?.title || `Entrada ${ticket.id}`}</h2>
                <p>{formatDate(ticket.purchase_date)}</p>
              </div>
              <span className={`status-badge ${ticket.status}`}>
                {ticket.status}
              </span>
            </div>

            {ticket.status === 'active' && (
              <div className="ticket-actions">
                <button
                  className="secondary-button"
                  type="button"
                  onClick={() => handleCancel(ticket.id)}
                  disabled={actionLoadingId === ticket.id}
                >
                  Cancelar
                </button>

                <div className="transfer-form">
                  <label htmlFor={`transfer-${ticket.id}`}>
                    Email destinatario
                  </label>
                  <div>
                    <input
                      id={`transfer-${ticket.id}`}
                      type="email"
                      value={transferEmails[ticket.id] || ''}
                      onChange={(event) =>
                        handleTransferEmailChange(ticket.id, event.target.value)
                      }
                      placeholder="otro@mail.com"
                    />
                    <button
                      className="primary-button"
                      type="button"
                      onClick={() => handleTransfer(ticket.id)}
                      disabled={actionLoadingId === ticket.id}
                    >
                      Transferir
                    </button>
                  </div>
                </div>
                {ticket.event?.is_festival && (
                  <Link
                    className="secondary-button"
                    to={`/mis-entradas/${ticket.event_id}/grilla`}
                  >
                    {getItineraryButtonText(ticket, itineraryItemCounts)}
                  </Link>
                )}
              </div>
            )}
          </li>
        ))}
      </ul>
    </main>
  )
}

async function loadItineraryItemCounts(tickets) {
  const festivalTickets = tickets.filter(
    (ticket) => ticket.status === 'active' && ticket.event?.is_festival,
  )

  if (festivalTickets.length === 0) {
    return {}
  }

  const countEntries = await Promise.all(
    festivalTickets.map(async (ticket) => {
      try {
        const result = await getItinerary(ticket.event_id)
        return [ticket.event_id, result.data?.items?.length || 0]
      } catch {
        return [ticket.event_id, 0]
      }
    }),
  )

  return Object.fromEntries(countEntries)
}

function getItineraryButtonText(ticket, counts) {
  return counts[ticket.event_id] > 0
    ? 'Ver / Editar mi grilla'
    : 'Crear grilla del evento'
}

export default MyTicketsPage
