export function formatDate(date, options = {}) {
  if (!date) {
    return 'Fecha a confirmar'
  }

  return new Intl.DateTimeFormat('es-AR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    ...options,
  }).format(new Date(date))
}

export function formatDateTime(date) {
  return formatDate(date, {
    month: 'long',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatTime(date) {
  if (!date) {
    return '--:--'
  }

  return new Intl.DateTimeFormat('es-AR', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(date))
}

export function formatPrice(price) {
  if (price === undefined || price === null) {
    return 'Precio a confirmar'
  }

  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
  }).format(price)
}

export function toDateTimeLocalValue(date) {
  if (!date) {
    return ''
  }

  const parsedDate = new Date(date)
  const timezoneOffset = parsedDate.getTimezoneOffset() * 60000
  return new Date(parsedDate.getTime() - timezoneOffset).toISOString().slice(0, 16)
}

export function fromDateTimeLocalValue(value) {
  if (!value) {
    return ''
  }

  return new Date(value).toISOString()
}
