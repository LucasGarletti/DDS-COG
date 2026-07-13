import { useCallback, useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getEventById } from '../services/eventService'
import { getPublicSchedule } from '../services/festivalService'
import { getErrorMessage } from '../services/httpClient'
import {
  addPersonalActivity,
  addShowToItinerary,
  deleteItineraryItem,
  getItinerary,
} from '../services/itineraryService'
import {
  formatTime,
  fromDateTimeLocalValue,
} from '../utils/formatters'

const emptyActivity = {
  title: '',
  location: '',
  notes: '',
  start_time: '',
  end_time: '',
}

function MyFestivalItineraryPage() {
  const { eventoId } = useParams()
  const [event, setEvent] = useState(null)
  const [schedule, setSchedule] = useState([])
  const [itinerary, setItinerary] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')
  const [actionLoadingId, setActionLoadingId] = useState(null)
  const [activityForm, setActivityForm] = useState(emptyActivity)

  const loadData = useCallback(async () => {
    setError('')
    setLoading(true)

    try {
      const [eventResult, scheduleResult, itineraryResult] = await Promise.all([
        getEventById(eventoId),
        getPublicSchedule(eventoId),
        getItinerary(eventoId),
      ])

      setEvent(eventResult.data)
      setSchedule(scheduleResult.data || [])
      setItinerary(itineraryResult.data)
    } catch (requestError) {
      setError(getItineraryErrorMessage(requestError))
    } finally {
      setLoading(false)
    }
  }, [eventoId])

  useEffect(() => {
    Promise.resolve().then(loadData)
  }, [loadData])

  function updateActivityField(field, value) {
    setActivityForm((currentForm) => ({ ...currentForm, [field]: value }))
  }

  async function handleAddShow(scheduleId) {
    setError('')
    setMessage('')
    setActionLoadingId(scheduleId)

    try {
      const result = await addShowToItinerary(eventoId, scheduleId)
      setMessage(
        result.data?.has_conflict
          ? 'Show agregado. Atención: se superpone con otra actividad.'
          : 'Show agregado a tu grilla',
      )
      await loadData()
    } catch (requestError) {
      setError(getItineraryErrorMessage(requestError))
    } finally {
      setActionLoadingId(null)
    }
  }

  async function handleAddActivity(submitEvent) {
    submitEvent.preventDefault()
    setError('')
    setMessage('')

    if (new Date(activityForm.end_time) <= new Date(activityForm.start_time)) {
      setError('La hora de fin debe ser posterior al inicio.')
      return
    }

    setActionLoadingId('activity')

    try {
      const result = await addPersonalActivity(eventoId, {
        title: activityForm.title,
        location: activityForm.location,
        notes: activityForm.notes,
        start_time: fromDateTimeLocalValue(activityForm.start_time),
        end_time: fromDateTimeLocalValue(activityForm.end_time),
      })
      setActivityForm(emptyActivity)
      setMessage(
        result.data?.has_conflict
          ? 'Actividad agregada. Atención: se superpone con otra actividad.'
          : 'Actividad agregada correctamente',
      )
      await loadData()
    } catch (requestError) {
      setError(getItineraryErrorMessage(requestError))
    } finally {
      setActionLoadingId(null)
    }
  }

  async function handleDeleteItem(itemId) {
    if (!window.confirm('¿Eliminar este item de tu grilla?')) {
      return
    }

    setError('')
    setMessage('')
    setActionLoadingId(itemId)

    try {
      await deleteItineraryItem(eventoId, itemId)
      setMessage('Item eliminado correctamente')
      await loadData()
    } catch (requestError) {
      setError(getItineraryErrorMessage(requestError))
    } finally {
      setActionLoadingId(null)
    }
  }

  return (
    <main className="page">
      <section className="page-hero compact">
        <div>
          <p>Mis entradas</p>
          <h1>Mi grilla del festival</h1>
        </div>
      </section>

      {loading && <p className="message">Cargando grilla...</p>}
      {message && <p className="message success">{message}</p>}
      {error && <p className="message error">{error}</p>}

      {event && !event.is_festival && (
        <p className="message error">
          Esta función solo está disponible para festivales.
        </p>
      )}

      {event?.is_festival && (
        <>
          <section className="festival-section">
            <div className="section-heading">
              <p>{event.title}</p>
              <h2>Grilla oficial</h2>
            </div>

            <div className="schedule-list">
              {schedule.map((show) => {
                const alreadyAdded = itinerary?.items?.some(
                  (item) => item.festival_schedule_id === show.id,
                )

                return (
                  <article className="schedule-card" key={show.id}>
                    {show.image_url && <img src={show.image_url} alt={show.artist} />}
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
                      onClick={() => handleAddShow(show.id)}
                      disabled={alreadyAdded || actionLoadingId === show.id}
                    >
                      {alreadyAdded ? 'En tu grilla' : 'Agregar a mi grilla'}
                    </button>
                  </article>
                )
              })}
            </div>

            {schedule.length === 0 && !loading && (
              <p className="message">Todavía no hay shows cargados.</p>
            )}
          </section>

          <section className="itinerary-panel">
            <div className="section-heading">
              <p>Hoja de ruta personal</p>
              <h2>Mi itinerario</h2>
            </div>

            {itinerary?.has_conflicts && (
              <p className="message warning">
                Tu itinerario tiene actividades superpuestas.
              </p>
            )}

            <ul className="itinerary-list">
              {(itinerary?.items || []).map((item) => (
                <li className="itinerary-item" key={item.id}>
                  <div>
                    <p className="event-date">
                      {formatTime(item.start_time)} - {formatTime(item.end_time)}
                    </p>
                    <h3>{item.title}</h3>
                    <p>{item.location || item.type}</p>
                    {item.notes && <p>{item.notes}</p>}
                  </div>
                  <button
                    className="secondary-button"
                    type="button"
                    onClick={() => handleDeleteItem(item.id)}
                    disabled={actionLoadingId === item.id}
                  >
                    Eliminar
                  </button>
                </li>
              ))}
            </ul>

            <form className="form itinerary-form" onSubmit={handleAddActivity}>
              <h3>Agregar actividad personal</h3>
              <div className="form-grid">
                <label>
                  Título
                  <input
                    value={activityForm.title}
                    onChange={(event) =>
                      updateActivityField('title', event.target.value)
                    }
                    required
                  />
                </label>
                <label>
                  Ubicación
                  <input
                    value={activityForm.location}
                    onChange={(event) =>
                      updateActivityField('location', event.target.value)
                    }
                  />
                </label>
                <label>
                  Inicio
                  <input
                    type="datetime-local"
                    value={activityForm.start_time}
                    onChange={(event) =>
                      updateActivityField('start_time', event.target.value)
                    }
                    required
                  />
                </label>
                <label>
                  Fin
                  <input
                    type="datetime-local"
                    value={activityForm.end_time}
                    onChange={(event) =>
                      updateActivityField('end_time', event.target.value)
                    }
                    required
                  />
                </label>
              </div>
              <label>
                Notas
                <textarea
                  value={activityForm.notes}
                  onChange={(event) =>
                    updateActivityField('notes', event.target.value)
                  }
                  rows="3"
                />
              </label>
              <button
                className="primary-button"
                type="submit"
                disabled={Boolean(actionLoadingId)}
              >
                Agregar actividad
              </button>
            </form>
          </section>
        </>
      )}

      <Link className="secondary-button" to="/mis-entradas">
        Volver a mis entradas
      </Link>
    </main>
  )
}

function getItineraryErrorMessage(error) {
  if (error?.response?.status === 401) {
    return 'Necesitás iniciar sesión para usar tu grilla.'
  }

  if (error?.response?.status === 403) {
    return 'Necesitás una entrada activa para crear la grilla de este festival.'
  }

  if (error?.response?.status === 400) {
    return 'Esta función solo está disponible para festivales.'
  }

  return getErrorMessage(error, 'No se pudo actualizar tu grilla')
}

export default MyFestivalItineraryPage
