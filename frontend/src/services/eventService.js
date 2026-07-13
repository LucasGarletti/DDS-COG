import axios from 'axios'
import { API_URL } from './httpClient'

export async function getEvents(filters = {}) {
  const params = new URLSearchParams()

  Object.entries(filters).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      params.set(key, String(value))
    }
  })

  const query = params.toString()
  const response = await axios.get(`${API_URL}/eventos${query ? `?${query}` : ''}`)
  return response.data
}

export async function getEventById(id) {
  const response = await axios.get(`${API_URL}/eventos/${id}`)
  return response.data
}
