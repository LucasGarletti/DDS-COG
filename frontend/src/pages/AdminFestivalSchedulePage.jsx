import { useCallback, useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import AdminNav from '../components/AdminNav'
import { getEventById } from '../services/eventService'
import {
  createSchedule,
  deleteSchedule,
  getAdminSchedule,
} from '../services/festivalService'
import { getErrorMessage } from '../services/httpClient'
import {
  formatDate,
  formatTime,
  fromDateTimeLocalValue,
  toDateTimeLocalValue,
} from '../utils/formatters'

const emptyShow = {
  artist: '',
  stage: '',
  start_time: '',
  end_time: '',
  image_url: '',
}

function AdminFestivalSchedulePage() {
  const { id } = useParams()
  const [event, setEvent] = useState(null)
  const [shows, setShows] = useState([])
  const [form, setForm] = useState(emptyShow)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  const loadSchedule = useCallback(async function loadSchedule() {
    setLoading(true)
    setError('')

    try {
      const [eventResult, scheduleResult] = await Promise.all([
        getEventById(id),
        getAdminSchedule(id),
      ])
      setEvent(eventResult.data)
      setShows(scheduleResult.data || [])
    } catch (requestError) {
      setError(getErrorMessage(requestError, 'No se pudo cargar la grilla'))
    } finally {
      setLoading(false)
    }
  }, [id])

  useEffect(() => {
    let ignore = false

    async function loadInitialSchedule() {
      try {
        const [eventResult, scheduleResult] = await Promise.all([
          getEventById(id),
          getAdminSchedule(id),
        ])

        if (ignore) {
          return
        }

        setEvent(eventResult.data)
        setShows(scheduleResult.data || [])
      } catch (requestError) {
        if (!ignore) {
          setError(getErrorMessage(requestError, 'No se pudo cargar la grilla'))
        }
      } finally {
        if (!ignore) {
          setLoading(false)
        }
      }
    }

    loadInitialSchedule()

    return () => {
      ignore = true
    }
  }, [id])

  function updateField(field, value) {
    setForm((currentForm) => ({ ...currentForm, [field]: value }))
  }

  async function handleSubmit(submitEvent) {
    submitEvent.preventDefault()
    setError('')
    setMessage('')

    const festivalDate = getEventDateValue(event?.date)

    if (!festivalDate) {
      setError('No se pudo identificar la fecha del festival')
      return
    }

    if (
      !form.artist.trim() ||
      !form.stage.trim() ||
      !form.start_time ||
      !form.end_time
    ) {
      setError('Completa artista, escenario, inicio y fin')
      return
    }

    const startDateTime = combineDateAndTime(festivalDate, form.start_time)
    const endDateTime = combineDateAndTime(festivalDate, form.end_time)

    if (new Date(endDateTime) <= new Date(startDateTime)) {
      setError('La hora de fin debe ser posterior al inicio')
      return
    }

    setSaving(true)

    try {
      await createSchedule(id, {
        artist: form.artist.trim(),
        stage: form.stage.trim(),
        start_time: fromDateTimeLocalValue(startDateTime),
        end_time: fromDateTimeLocalValue(endDateTime),
        image_url: form.image_url.trim(),
      })
      setMessage('Show creado correctamente')
      setForm(emptyShow)
      await loadSchedule()
    } catch (requestError) {
      setError(getErrorMessage(requestError, 'No se pudo crear el show'))
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(showId) {
    if (!window.confirm('Eliminar este show de la grilla?')) {
      return
    }

    setError('')
    setMessage('')

    try {
      await deleteSchedule(showId)
      setMessage('Show eliminado correctamente')
      await loadSchedule()
    } catch (requestError) {
      setError(getErrorMessage(requestError, 'No se pudo eliminar el show'))
    }
  }

  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Administracion</p>
          <h1>Grilla del festival</h1>
        </div>
      </section>

      <AdminNav />

      {message && <p className="message success">{message}</p>}
      {error && <p className="message error">{error}</p>}
      {event && (
        <p className="message">
          Fecha del festival: {formatDate(event.date, { month: 'long' })}
        </p>
      )}

      <form className="form admin-form" onSubmit={handleSubmit}>
        <div className="form-grid">
          <label>
            Artista
            <input
              value={form.artist}
              onChange={(inputEvent) =>
                updateField('artist', inputEvent.target.value)
              }
            />
          </label>
          <label>
            Escenario
            <input
              value={form.stage}
              onChange={(inputEvent) =>
                updateField('stage', inputEvent.target.value)
              }
            />
          </label>
          <label>
            Hora de inicio
            <input
              type="time"
              value={form.start_time}
              onChange={(inputEvent) =>
                updateField('start_time', inputEvent.target.value)
              }
            />
          </label>
          <label>
            Hora de finalizacion
            <input
              type="time"
              value={form.end_time}
              onChange={(inputEvent) =>
                updateField('end_time', inputEvent.target.value)
              }
            />
          </label>
        </div>

        <label>
          URL de imagen
          <input
            value={form.image_url}
            onChange={(inputEvent) =>
              updateField('image_url', inputEvent.target.value)
            }
          />
        </label>

        <button className="primary-button" type="submit" disabled={saving}>
          {saving ? 'Guardando...' : 'Crear show'}
        </button>
      </form>

      {loading && <p className="message">Cargando shows...</p>}

      <div className="schedule-list">
        {shows.map((show) => (
          <article className="schedule-card" key={show.id}>
            <div>
              <p className="event-date">
                {formatTime(show.start_time)} - {formatTime(show.end_time)}
              </p>
              <h3>{show.artist}</h3>
              <p>{show.stage}</p>
            </div>
            <button
              className="secondary-button"
              type="button"
              onClick={() => handleDelete(show.id)}
            >
              Eliminar
            </button>
          </article>
        ))}
      </div>
    </main>
  )
}

function getEventDateValue(date) {
  return toDateTimeLocalValue(date).slice(0, 10)
}

function combineDateAndTime(date, time) {
  return `${date}T${time}`
}

export default AdminFestivalSchedulePage
