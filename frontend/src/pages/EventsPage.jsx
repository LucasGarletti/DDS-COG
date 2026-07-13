import { useEffect, useState } from 'react'
import EventCard from '../components/EventCard'
import { getEvents } from '../services/eventService'
import { getErrorMessage } from '../services/httpClient'

const emptyFilters = {
  search: '',
  location: '',
  date_from: '',
  date_to: '',
  min_price: '',
  max_price: '',
  is_festival: false,
  available: false,
  sort: 'date_asc',
}

function EventsPage() {
  const [events, setEvents] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [filters, setFilters] = useState(emptyFilters)
  const [appliedFilters, setAppliedFilters] = useState(emptyFilters)

  useEffect(() => {
    async function loadEvents() {
      setLoading(true)
      setError('')
      try {
        const result = await getEvents(buildEventFilterParams(appliedFilters))
        setEvents(result.data || [])
      } catch (requestError) {
        setError(getErrorMessage(requestError, 'No se pudieron cargar los eventos'))
      } finally {
        setLoading(false)
      }
    }

    loadEvents()
  }, [appliedFilters])

  function updateFilter(field, value) {
    setFilters((currentFilters) => ({ ...currentFilters, [field]: value }))
  }

  function handleApplyFilters(event) {
    event.preventDefault()
    setAppliedFilters(filters)
  }

  function handleClearFilters() {
    setFilters(emptyFilters)
    setAppliedFilters(emptyFilters)
  }

  return (
    <main className="page">
      <section className="page-hero">
        <div>
          <p>Entradas para shows y experiencias</p>
          <h1>Elegi tu proximo evento</h1>
        </div>
      </section>

      <form className="filter-panel" onSubmit={handleApplyFilters}>
        <div className="filter-row filter-row-main">
          <label>
            Buscar
            <input
              value={filters.search}
              onChange={(event) => updateFilter('search', event.target.value)}
              placeholder="Rock, teatro, festival"
            />
          </label>
          <label>
            Ubicación
            <input
              value={filters.location}
              onChange={(event) => updateFilter('location', event.target.value)}
              placeholder="Córdoba"
            />
          </label>
          <label>
            Desde
            <input
              type="date"
              value={filters.date_from}
              onChange={(event) => updateFilter('date_from', event.target.value)}
            />
          </label>
          <label>
            Hasta
            <input
              type="date"
              value={filters.date_to}
              onChange={(event) => updateFilter('date_to', event.target.value)}
            />
          </label>
        </div>
        <div className="filter-row filter-row-secondary">
          <label>
            Precio mínimo
            <input
              type="number"
              min="0"
              value={filters.min_price}
              onChange={(event) => updateFilter('min_price', event.target.value)}
            />
          </label>
          <label>
            Precio máximo
            <input
              type="number"
              min="0"
              value={filters.max_price}
              onChange={(event) => updateFilter('max_price', event.target.value)}
            />
          </label>
          <label>
            Ordenamiento
            <select
              value={filters.sort}
              onChange={(event) => updateFilter('sort', event.target.value)}
            >
              <option value="date_asc">Fecha ascendente</option>
              <option value="date_desc">Fecha descendente</option>
              <option value="price_asc">Precio ascendente</option>
              <option value="price_desc">Precio descendente</option>
            </select>
          </label>
        </div>
        <div className="filter-row filter-row-toggles">
          <label className="checkbox-field">
            <input
              type="checkbox"
              checked={filters.is_festival}
              onChange={(event) =>
                updateFilter('is_festival', event.target.checked)
              }
            />
            Solo festivales
          </label>
          <label className="checkbox-field">
            <input
              type="checkbox"
              checked={filters.available}
              onChange={(event) => updateFilter('available', event.target.checked)}
            />
            Solo con disponibilidad
          </label>
        </div>
        <div className="filter-actions">
          <button className="primary-button" type="submit" disabled={loading}>
            Aplicar filtros
          </button>
          <button
            className="secondary-button"
            type="button"
            onClick={handleClearFilters}
          >
            Limpiar filtros
          </button>
        </div>
      </form>

      {loading && <p className="message">Cargando eventos...</p>}
      {error && <p className="message error">{error}</p>}
      {!loading && !error && events.length === 0 && (
        <p className="message">No encontramos eventos con esos filtros.</p>
      )}

      <ul className="event-list">
        {events.map((event) => (
          <EventCard event={event} key={event.id} />
        ))}
      </ul>
    </main>
  )
}

function buildEventFilterParams(filters) {
  const params = {}

  Object.entries(filters).forEach(([key, value]) => {
    if (key === 'sort' && value === 'date_asc') {
      return
    }

    if (typeof value === 'boolean') {
      if (value) {
        params[key] = value
      }
      return
    }

    if (value !== '') {
      params[key] = value
    }
  })

  return params
}

export default EventsPage
