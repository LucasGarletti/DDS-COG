import { getAuth } from './httpClient'

export async function getAdminSummary() {
  return getAuth('/admin/reportes/resumen')
}

export async function getAdminEventReports() {
  return getAuth('/admin/reportes/eventos')
}
