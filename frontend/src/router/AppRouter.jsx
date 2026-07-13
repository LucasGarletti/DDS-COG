import { useEffect, useState } from 'react'
import {
  BrowserRouter,
  Link,
  Navigate,
  Route,
  Routes,
  useNavigate,
} from 'react-router-dom'
import AdminEventFormPage from '../pages/AdminEventFormPage'
import AdminEventReportPage from '../pages/AdminEventReportPage'
import AdminEventsPage from '../pages/AdminEventsPage'
import AdminFestivalSchedulePage from '../pages/AdminFestivalSchedulePage'
import AdminHomePage from '../pages/AdminHomePage'
import AdminReportsPage from '../pages/AdminReportsPage'
import EventDetailPage from '../pages/EventDetailPage'
import EventsPage from '../pages/EventsPage'
import HomePage from '../pages/HomePage'
import LoginPage from '../pages/LoginPage'
import MyFestivalItineraryPage from '../pages/MyFestivalItineraryPage'
import MyTicketsPage from '../pages/MyTicketsPage'
import { getMe } from '../services/authService'

function ProtectedRoute({ isLoggedIn, children }) {
  if (!isLoggedIn) {
    return <Navigate to="/login" replace />
  }

  return children
}

function AdminRoute({ authChecked, isLoggedIn, user, children }) {
  if (!isLoggedIn) {
    return <Navigate to="/login" replace />
  }

  if (!authChecked) {
    return <p className="message">Verificando permisos...</p>
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (user && user.role !== 'admin') {
    return <Navigate to="/eventos" replace />
  }

  return children
}

function AppContent() {
  const navigate = useNavigate()
  const [authToken, setAuthToken] = useState(localStorage.getItem('token'))
  const [user, setUser] = useState(null)
  const [authChecked, setAuthChecked] = useState(!localStorage.getItem('token'))
  const isLoggedIn = Boolean(authToken)

  useEffect(() => {
    async function loadCurrentUser() {
      if (!authToken) {
        setUser(null)
        setAuthChecked(true)
        return
      }

      setAuthChecked(false)
      try {
        const result = await getMe()
        setUser(result.data)
      } catch {
        setUser(null)
      } finally {
        setAuthChecked(true)
      }
    }

    loadCurrentUser()
  }, [authToken])

  function handleLogin(token) {
    setAuthToken(token || localStorage.getItem('token'))
  }

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('tickets')
    localStorage.removeItem('myTickets')
    localStorage.removeItem('misEntradas')
    setUser(null)
    setAuthChecked(true)
    setAuthToken(null)
    navigate('/login')
  }

  return (
    <div className="app-shell">
      <nav className="app-nav">
        <Link className="brand" to="/">
          TickGo
        </Link>
        <div className="nav-links">
          <Link to="/eventos">Eventos</Link>
          {isLoggedIn ? (
            <>
              <Link to="/mis-entradas">Mis entradas</Link>
              {user?.role === 'admin' && <Link to="/admin">Admin</Link>}
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
        <Route path="/" element={<HomePage />} />
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
        <Route
          path="/mis-entradas/:eventoId/grilla"
          element={
            <ProtectedRoute isLoggedIn={isLoggedIn}>
              <MyFestivalItineraryPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="/admin"
          element={
            <AdminRoute authChecked={authChecked} isLoggedIn={isLoggedIn} user={user}>
              <AdminHomePage />
            </AdminRoute>
          }
        />
        <Route
          path="/admin/eventos"
          element={
            <AdminRoute authChecked={authChecked} isLoggedIn={isLoggedIn} user={user}>
              <AdminEventsPage />
            </AdminRoute>
          }
        />
        <Route
          path="/admin/eventos/nuevo"
          element={
            <AdminRoute authChecked={authChecked} isLoggedIn={isLoggedIn} user={user}>
              <AdminEventFormPage mode="create" />
            </AdminRoute>
          }
        />
        <Route
          path="/admin/eventos/:id/editar"
          element={
            <AdminRoute authChecked={authChecked} isLoggedIn={isLoggedIn} user={user}>
              <AdminEventFormPage mode="edit" />
            </AdminRoute>
          }
        />
        <Route
          path="/admin/eventos/:id/grilla"
          element={
            <AdminRoute authChecked={authChecked} isLoggedIn={isLoggedIn} user={user}>
              <AdminFestivalSchedulePage />
            </AdminRoute>
          }
        />
        <Route
          path="/admin/eventos/:id/reporte"
          element={
            <AdminRoute authChecked={authChecked} isLoggedIn={isLoggedIn} user={user}>
              <AdminEventReportPage />
            </AdminRoute>
          }
        />
        <Route
          path="/admin/reportes"
          element={
            <AdminRoute authChecked={authChecked} isLoggedIn={isLoggedIn} user={user}>
              <AdminReportsPage />
            </AdminRoute>
          }
        />
      </Routes>

      <footer className="app-footer">
        <strong>TickGo</strong>
        <span>Sistema de gestion y compra de entradas</span>
      </footer>
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
