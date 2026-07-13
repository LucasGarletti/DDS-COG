import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import AdminNav from '../components/AdminNav'
import EventForm from '../components/EventForm'
import { createEvent, updateEvent } from '../services/adminEventService'
import { getEventById } from '../services/eventService'

function AdminEventFormPage({ mode }) {
  const { id } = useParams()
  const navigate = useNavigate()
  const isEdit = mode === 'edit'
  const [event, setEvent] = useState(null)
  const [loading, setLoading] = useState(false)
  const [pageLoading, setPageLoading] = useState(isEdit)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  useEffect(() => {
    if (!isEdit) {
      return
    }

    async function loadEvent() {
      try {
        const result = await getEventById(id)
        setEvent(result.data)
      } catch {
        setError('No se pudo cargar el evento')
      } finally {
        setPageLoading(false)
      }
    }

    loadEvent()
  }, [id, isEdit])

  async function handleSubmit(eventData) {
    setLoading(true)
    setError('')
    setMessage('')

    try {
      if (isEdit) {
        await updateEvent(id, eventData)
        setMessage('Evento actualizado correctamente')
      } else {
        const result = await createEvent(eventData)
        setMessage('Evento creado correctamente')
        navigate(`/admin/eventos/${result.data.id}/editar`)
      }
    } catch {
      setError('No se pudo guardar el evento')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Administración</p>
          <h1>{isEdit ? 'Editar evento' : 'Crear evento'}</h1>
        </div>
      </section>
      <AdminNav />

      {pageLoading && <p className="message">Cargando evento...</p>}
      {message && <p className="message success">{message}</p>}
      {error && <p className="message error">{error}</p>}
      {(!isEdit || event) && (
        <EventForm
          key={event?.id || 'new'}
          initialEvent={event}
          loading={loading}
          submitLabel={isEdit ? 'Guardar cambios' : 'Crear evento'}
          onSubmit={handleSubmit}
        />
      )}
    </main>
  )
}

export default AdminEventFormPage
