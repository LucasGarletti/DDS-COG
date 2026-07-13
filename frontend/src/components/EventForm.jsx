import { useState } from 'react'
import { fromDateTimeLocalValue, toDateTimeLocalValue } from '../utils/formatters'

function EventForm({ initialEvent, loading, submitLabel, onSubmit }) {
  const [form, setForm] = useState(() => ({
    title: initialEvent?.title || '',
    description: initialEvent?.description || '',
    date: toDateTimeLocalValue(initialEvent?.date),
    location: initialEvent?.location || '',
    capacity: initialEvent?.capacity ?? '',
    price: initialEvent?.price ?? '',
    image_url: initialEvent?.image_url || '',
    is_festival: Boolean(initialEvent?.is_festival),
  }))
  const [validationError, setValidationError] = useState('')

  function updateField(field, value) {
    setForm((currentForm) => ({ ...currentForm, [field]: value }))
  }

  function handleSubmit(event) {
    event.preventDefault()
    setValidationError('')

    if (!form.title.trim() || !form.location.trim() || !form.date) {
      setValidationError('Completá título, fecha y ubicación')
      return
    }

    if (Number(form.capacity) <= 0) {
      setValidationError('La capacidad debe ser mayor que cero')
      return
    }

    if (Number(form.price) < 0) {
      setValidationError('El precio no puede ser negativo')
      return
    }

    onSubmit({
      title: form.title.trim(),
      description: form.description,
      date: fromDateTimeLocalValue(form.date),
      location: form.location.trim(),
      capacity: Number(form.capacity),
      price: Number(form.price),
      image_url: form.image_url.trim(),
      is_festival: form.is_festival,
    })
  }

  return (
    <form className="form admin-form" onSubmit={handleSubmit}>
      <label>
        Título
        <input
          value={form.title}
          onChange={(event) => updateField('title', event.target.value)}
          required
        />
      </label>
      <label>
        Descripción
        <textarea
          value={form.description}
          onChange={(event) => updateField('description', event.target.value)}
          rows="4"
        />
      </label>
      <div className="form-grid">
        <label>
          Fecha y hora
          <input
            type="datetime-local"
            value={form.date}
            onChange={(event) => updateField('date', event.target.value)}
            required
          />
        </label>
        <label>
          Ubicación
          <input
            value={form.location}
            onChange={(event) => updateField('location', event.target.value)}
            required
          />
        </label>
        <label>
          Capacidad
          <input
            type="number"
            min="1"
            value={form.capacity}
            onChange={(event) => updateField('capacity', event.target.value)}
            required
          />
        </label>
        <label>
          Precio
          <input
            type="number"
            min="0"
            step="0.01"
            value={form.price}
            onChange={(event) => updateField('price', event.target.value)}
            required
          />
        </label>
      </div>
      <label>
        URL de imagen
        <input
          value={form.image_url}
          onChange={(event) => updateField('image_url', event.target.value)}
          placeholder="https://..."
        />
      </label>
      <label className="checkbox-field">
        <input
          type="checkbox"
          checked={form.is_festival}
          onChange={(event) => updateField('is_festival', event.target.checked)}
        />
        Es festival
      </label>

      {validationError && <p className="message error">{validationError}</p>}

      <button className="primary-button" type="submit" disabled={loading}>
        {loading ? 'Guardando...' : submitLabel}
      </button>
    </form>
  )
}

export default EventForm
