import axios from 'axios'

const API_URL = 'http://localhost:8080'

function getAuthHeaders() {
  const token = localStorage.getItem('token')

  return {
    Authorization: `Bearer ${token}`,
    'Cache-Control': 'no-cache',
    Pragma: 'no-cache',
  }
}

export async function getMyTickets() {
  const response = await axios.get(`${API_URL}/mis-entradas`, {
    headers: getAuthHeaders(),
  })

  return response.data
}

export async function purchaseTicket(eventId) {
  const response = await axios.post(
    `${API_URL}/entradas/comprar/${eventId}`,
    {},
    {
      headers: getAuthHeaders(),
    },
  )

  return response.data
}

export async function cancelTicket(ticketId) {
  const response = await axios.patch(
    `${API_URL}/entradas/${ticketId}/cancelar`,
    {},
    {
      headers: getAuthHeaders(),
    },
  )

  return response.data
}

export async function transferTicket(ticketId, recipientEmail) {
  const response = await axios.patch(
    `${API_URL}/entradas/${ticketId}/transferir`,
    {
      email_destinatario: recipientEmail,
    },
    {
      headers: getAuthHeaders(),
    },
  )

  return response.data
}
