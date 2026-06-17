import { BrowserRouter, Link, Navigate, Route, Routes } from 'react-router-dom'
import EventDetailPage from '../pages/EventDetailPage'
import EventsPage from '../pages/EventsPage'
import LoginPage from '../pages/LoginPage'
import MyTicketsPage from '../pages/MyTicketsPage'

function AppRouter() {
  return (
    <BrowserRouter>
      <div className="app-shell">
        <nav className="app-nav">
          <Link className="brand" to="/eventos">
            TickGo
          </Link>
          <div className="nav-links">
            <Link to="/eventos">Eventos</Link>
            <Link to="/mis-entradas">Mis entradas</Link>
            <Link to="/login">Login</Link>
          </div>
        </nav>

        <Routes>
          <Route path="/" element={<Navigate to="/eventos" replace />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/eventos" element={<EventsPage />} />
          <Route path="/eventos/:id" element={<EventDetailPage />} />
          <Route path="/mis-entradas" element={<MyTicketsPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}

export default AppRouter
