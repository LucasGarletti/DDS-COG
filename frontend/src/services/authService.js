import axios from 'axios'
import { API_URL, getAuth } from './httpClient'

export async function login(email, password) {
  const response = await axios.post(`${API_URL}/auth/login`, {
    email,
    password,
  })

  return response.data
}

export async function getMe() {
  return getAuth('/auth/me')
}
