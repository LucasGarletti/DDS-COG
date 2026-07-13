import { Link } from 'react-router-dom'
import { getEventImage } from '../utils/eventImages'
import { formatDate, formatPrice } from '../utils/formatters'

function EventCard({ event }) {
  return (
    <li className="event-card">
      <img src={getEventImage(event)} alt={event.title} className="event-image" />
      <div className="event-card-body">
        <p className="event-date">{formatDate(event.date)}</p>
        <h2>{event.title}</h2>
        {event.description && (
          <p className="event-summary">{event.description.slice(0, 120)}</p>
        )}
        <p className="event-location">{event.location}</p>
        <div className="event-meta">
          <span>{formatPrice(event.price)}</span>
          <span>{event.available_capacity} disponibles</span>
          {event.is_festival && <span>Festival</span>}
        </div>
        <Link className="primary-button" to={`/eventos/${event.id}`}>
          Ver detalle
        </Link>
      </div>
    </li>
  )
}

export default EventCard
