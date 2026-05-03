import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { loginUser } from '../api/auth'
import './LoginPage.css'

const initialFormData = {
  email: '',
  password: '',
}

const initialErrors = {
  email: '',
  password: '',
}

function LoginPage() {
  const navigate = useNavigate()

  const [formData, setFormData] = useState(initialFormData)
  const [errors, setErrors] = useState(initialErrors)
  const [serverMessage, setServerMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)

  function validateForm(values) {
    const nextErrors = {
      email: '',
      password: '',
    }

    const trimmedEmail = values.email.trim()
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    if (!trimmedEmail) {
      nextErrors.email = 'Email is required'
    } else if (!emailRegex.test(trimmedEmail)) {
      nextErrors.email = 'Enter a valid email'
    }

    if (!values.password) {
      nextErrors.password = 'Password is required'
    }

    return nextErrors
  }

  function handleChange(event) {
    const { name, value } = event.target

    const nextFormData = {
      ...formData,
      [name]: value,
    }

    const nextErrors = {
      ...errors,
      [name]: '',
    }

    if (name === 'email') {
      const trimmedEmail = value.trim()
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

      if (trimmedEmail && !emailRegex.test(trimmedEmail)) {
        nextErrors.email = 'Enter a valid email'
      }
    }

    setFormData(nextFormData)
    setErrors(nextErrors)
    setServerMessage('')
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setServerMessage('')

    const nextErrors = validateForm(formData)
    setErrors(nextErrors)

    const hasErrors = Object.values(nextErrors).some(Boolean)
    if (hasErrors) {
      return
    }

    setIsLoading(true)

    try {
      const result = await loginUser({
        email: formData.email.trim(),
        password: formData.password,
      })

      localStorage.setItem('token', result.token)
      localStorage.setItem('userEmail', formData.email.trim())

      setServerMessage('Login successful. Redirecting...')

      setTimeout(() => {
        navigate('/app')
      }, 1000)
    } catch (error) {
      setServerMessage(`Error: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-card">
        <div className="login-form-side">
          <div className="login-form-content">
            <h1 className="login-title">Sign In</h1>

            <div className="login-social-buttons">
              <button type="button" className="login-social-button">
                <span className="login-social-icon login-social-icon--google">G</span>
                <span>Google</span>
              </button>

              <button type="button" className="login-social-button">
                <span className="login-social-icon login-social-icon--yandex">Я</span>
                <span>Yandex</span>
              </button>
            </div>

            <div className="login-divider">
              <span className="login-divider-line" />
              <span className="login-divider-text">Or</span>
              <span className="login-divider-line" />
            </div>

            <form className="login-form" onSubmit={handleSubmit} noValidate>
              <div className="login-form-group">
                <label className="login-form-label" htmlFor="email">
                  Email
                </label>
                <input
                  id="email"
                  className={`login-form-input ${errors.email ? 'login-form-input--error' : ''}`}
                  type="email"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  placeholder="example@mail.com"
                  autoComplete="email"
                />
                {errors.email && (
                  <p className="login-form-error" role="alert">
                    {errors.email}
                  </p>
                )}
              </div>

              <div className="login-form-group">
                <label className="login-form-label" htmlFor="password">
                  Password
                </label>

                <div className="login-password-wrapper">
                  <input
                    id="password"
                    className={`login-form-input login-form-input--with-button ${errors.password ? 'login-form-input--error' : ''}`}
                    type={showPassword ? 'text' : 'password'}
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="••••••••"
                    autoComplete="current-password"
                  />
                  <button
                    type="button"
                    className="login-password-toggle"
                    onClick={() => setShowPassword((prev) => !prev)}
                  >
                    {showPassword ? 'Hide' : 'Show'}
                  </button>
                </div>

                {errors.password && (
                  <p className="login-form-error" role="alert">
                    {errors.password}
                  </p>
                )}
              </div>

              <button className="login-submit-button" type="submit" disabled={isLoading}>
                {isLoading ? 'Signing In...' : 'Sign In'}
              </button>
            </form>

            {serverMessage && (
              <p
                className={`login-server-message ${
                  serverMessage.startsWith('Error')
                    ? 'login-server-message--error'
                    : 'login-server-message--success'
                }`}
              >
                {serverMessage}
              </p>
            )}

            <p className="login-signup-text">
              Don&apos;t have an account? <Link to="/register">Sign up</Link>
            </p>
          </div>
        </div>

        <div className="login-visual-side" aria-hidden="true">
          <div className="login-visual-card">
            <span className="login-visual-glow login-visual-glow--top" />
            <span className="login-visual-glow login-visual-glow--center" />
            <span className="login-visual-glow login-visual-glow--bottom" />
          </div>
        </div>
      </div>
    </div>
  )
}

export default LoginPage