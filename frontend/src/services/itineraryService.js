import { deleteAuth, getAuth, postAuth } from './httpClient'

export async function getItinerary(eventId) {
  return getAuth(`/mis-itinerarios/${eventId}`)
}

export async function addShowToItinerary(eventId, festivalScheduleId) {
  return postAuth(`/mis-itinerarios/${eventId}/shows`, {
    festival_schedule_id: festivalScheduleId,
  })
}

export async function addPersonalActivity(eventId, activityData) {
  return postAuth(`/mis-itinerarios/${eventId}/actividades`, activityData)
}

export async function deleteItineraryItem(eventId, itemId) {
  return deleteAuth(`/mis-itinerarios/${eventId}/items/${itemId}`)
}
