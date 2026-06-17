import { useState } from 'react'
import {
  BrowserRouter,
  Link,
  Navigate,
  Route,
  Routes,
  useNavigate,
} from 'react-router-dom'
import EventDetailPage from '../pages/EventDetailPage'
import EventsPage from '../pages/EventsPage'
import LoginPage from '../pages/LoginPage'
import MyTicketsPage from '../pages/MyTicketsPage'

function ProtectedRoute({ isLoggedIn, children }) {
  if (!isLoggedIn) {
    return <Navigate to="/login" replace />
  }

  return children
}

function AppContent() {
  const navigate = useNavigate()
  const [authToken, setAuthToken] = useState(localStorage.getItem('token'))
  const isLoggedIn = Boolean(authToken)

  function handleLogin(token) {
    setAuthToken(token || localStorage.getItem('token'))
  }

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('tickets')
    localStorage.removeItem('myTickets')
    localStorage.removeItem('misEntradas')
    setAuthToken(null)
    navigate('/login')
  }

  return (
    <div className="app-shell">
      <nav className="app-nav">
        <Link className="brand" to="/eventos">
          TickGo
        </Link>
        <div className="nav-links">
          <Link to="/eventos">Eventos</Link>
          {isLoggedIn ? (
            <>
              <Link to="/mis-entradas">Mis entradas</Link>
              <button
                className="nav-link-button"
                type="button"
                onClick={handleLogout}
              >
                Cerrar sesion
              </button>
            </>
          ) : (
            <Link to="/login">Login</Link>
          )}
        </div>
      </nav>

      <Routes>
        <Route path="/" element={<Navigate to="/eventos" replace />} />
        <Route path="/login" element={<LoginPage onLogin={handleLogin} />} />
        <Route path="/eventos" element={<EventsPage />} />
        <Route path="/eventos/:id" element={<EventDetailPage />} />
        <Route
          path="/mis-entradas"
          element={
            <ProtectedRoute isLoggedIn={isLoggedIn}>
              <MyTicketsPage key={authToken} authToken={authToken} />
            </ProtectedRoute>
          }
        />
      </Routes>
    </div>
  )
}

function AppRouter() {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  )
}

export default AppRouter
