import cosquiRockImage from '../assets/cosquirock.png'
import pastillasImage from '../assets/pastillasdelabuelo.jfif'
import pumasImage from '../assets/pumas.jfif'

const fallbackEventImage =
  'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 450"%3E%3Crect width="800" height="450" fill="%23006fff"/%3E%3Ctext x="50%25" y="50%25" dominant-baseline="middle" text-anchor="middle" font-family="Arial" font-size="56" font-weight="700" fill="white"%3ETickGo%3C/text%3E%3C/svg%3E'

function normalizeTitle(title) {
  return (title || '')
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
}

export function getEventImage(event) {
  const title = normalizeTitle(event.title)

  if (title.includes('cosquin rock')) {
    return cosquiRockImage
  }

  if (title.includes('pumas') && title.includes('kempes')) {
    return pumasImage
  }

  if (title.includes('pastillas del abuelo')) {
    return pastillasImage
  }

  return event.image_url || fallbackEventImage
}
