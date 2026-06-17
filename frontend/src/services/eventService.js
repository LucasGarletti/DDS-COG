import axios from 'axios'

const API_URL = 'http://localhost:8080'

export async function getEvents() {
  const response = await axios.get(`${API_URL}/eventos`)
  return response.data
}

export async function getEventById(id) {
  const response = await axios.get(`${API_URL}/eventos/${id}`)
  return response.data
}
