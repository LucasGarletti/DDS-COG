import { deleteAuth, getAuth, patchAuth, postAuth } from './httpClient'

export async function createEvent(eventData) {
  return postAuth('/admin/eventos', eventData)
}

export async function updateEvent(eventId, eventData) {
  return patchAuth(`/admin/eventos/${eventId}`, eventData)
}

export async function cancelEvent(eventId) {
  return deleteAuth(`/admin/eventos/${eventId}`)
}

export async function getEventReport(eventId) {
  return getAuth(`/admin/eventos/${eventId}/reporte`)
}
