import { useEffect, useState } from 'react'
import {
  cancelTicket,
  getMyTickets,
  transferTicket,
} from '../services/ticketService'

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

function MyTicketsPage() {
  const [tickets, setTickets] = useState([])
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(true)
  const [actionLoadingId, setActionLoadingId] = useState(null)
  const [transferEmails, setTransferEmails] = useState({})

  async function loadTickets() {
    try {
      const result = await getMyTickets()
      setTickets(result.data || [])
    } catch {
      setError('No se pudieron cargar tus entradas')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTickets()
  }, [])

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
              </div>
            )}
          </li>
        ))}
      </ul>
    </main>
  )
}

export default MyTicketsPage
