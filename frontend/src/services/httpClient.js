import axios from 'axios'

export const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function getAuthHeaders() {
  const token = localStorage.getItem('token')

  return {
    Authorization: `Bearer ${token}`,
    'Cache-Control': 'no-cache',
    Pragma: 'no-cache',
  }
}

export function getErrorMessage(error, fallback) {
  return error?.response?.data?.error || fallback
}

export async function getAuth(path) {
  const response = await axios.get(`${API_URL}${path}`, {
    headers: getAuthHeaders(),
  })
  return response.data
}

export async function postAuth(path, body) {
  const response = await axios.post(`${API_URL}${path}`, body, {
    headers: getAuthHeaders(),
  })
  return response.data
}

export async function patchAuth(path, body) {
  const response = await axios.patch(`${API_URL}${path}`, body, {
    headers: getAuthHeaders(),
  })
  return response.data
}

export async function deleteAuth(path) {
  const response = await axios.delete(`${API_URL}${path}`, {
    headers: getAuthHeaders(),
  })
  return response.data
}
