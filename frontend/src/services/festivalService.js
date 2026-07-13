import { deleteAuth, getAuth, patchAuth, postAuth } from './httpClient'
import { API_URL } from './httpClient'
import axios from 'axios'

export async function getPublicSchedule(eventId) {
  const response = await axios.get(`${API_URL}/eventos/${eventId}/grilla`)
  return response.data
}

export async function getAdminSchedule(eventId) {
  return getAuth(`/admin/eventos/${eventId}/grilla`)
}

export async function createSchedule(eventId, scheduleData) {
  return postAuth(`/admin/eventos/${eventId}/grilla`, scheduleData)
}

export async function updateSchedule(scheduleId, scheduleData) {
  return patchAuth(`/admin/grilla/${scheduleId}`, scheduleData)
}

export async function deleteSchedule(scheduleId) {
  return deleteAuth(`/admin/grilla/${scheduleId}`)
}
