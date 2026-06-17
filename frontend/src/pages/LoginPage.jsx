import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../services/authService'

function LoginPage({ onLogin }) {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(event) {
    event.preventDefault()
    setError('')
    setLoading(true)

    try {
      const result = await login(email, password)
      localStorage.setItem('token', result.data.token)
      onLogin?.(result.data.token)
      navigate('/eventos')
    } catch {
      setError('No se pudo iniciar sesion')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page login-page">
      <section className="login-card">
        <div className="section-heading">
          <p>Acceso de usuarios</p>
          <h1>Ingresar a TickGo</h1>
        </div>

        <form className="form" onSubmit={handleSubmit}>
          <label>
            Email
            <input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="tu@email.com"
              required
            />
          </label>

          <label>
            Password
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="Tu password"
              required
            />
          </label>

          {error && <p className="message error">{error}</p>}

          <button className="primary-button" type="submit" disabled={loading}>
            {loading ? 'Ingresando...' : 'Ingresar'}
          </button>
        </form>
      </section>
    </main>
  )
}

export default LoginPage
