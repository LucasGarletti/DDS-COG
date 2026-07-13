import { deleteAuth, getAuth, patchAuth, postAuth } from './httpClient'

export async function getMyTickets() {
  return getAuth('/mis-entradas')
}

export async function purchaseTicket(eventId) {
  return postAuth(`/entradas/comprar/${eventId}`, {})
}

export async function cancelTicket(ticketId) {
  return patchAuth(`/entradas/${ticketId}/cancelar`, {})
}

export async function transferTicket(ticketId, recipientEmail) {
  return patchAuth(`/entradas/${ticketId}/transferir`, {
    email_destinatario: recipientEmail,
  })
}

export async function deleteTicketRequest(path) {
  return deleteAuth(path)
}
