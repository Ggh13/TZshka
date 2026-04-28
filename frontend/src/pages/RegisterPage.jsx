import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { registerUser } from '../api/auth'
import './RegisterPage.css'

const initialFormData = {
  login: '',
  email: '',
  password: '',
  confirmPassword: '',
}

const initialErrors = {
  login: '',
  email: '',
  password: '',
  confirmPassword: '',
}

function RegisterPage() {
  const navigate = useNavigate()
  const [formData, setFormData] = useState(initialFormData)
  const [errors, setErrors] = useState(initialErrors)
  const [serverMessage, setServerMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)

  function validateForm(values) {
    const nextErrors = {
      login: '',
      email: '',
      password: '',
      confirmPassword: '',
    }

    const trimmedLogin = values.login.trim()
    const trimmedEmail = values.email.trim()
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    if (!trimmedLogin) {
      nextErrors.login = 'Login is required'
    }

    if (!trimmedEmail) {
      nextErrors.email = 'Email is required'
    } else if (!emailRegex.test(trimmedEmail)) {
      nextErrors.email = 'Enter a valid email'
    }

    if (!values.password) {
      nextErrors.password = 'Password is required'
    } else if (values.password.length < 8) {
      nextErrors.password = 'Password must be at least 8 characters'
    }

    if (!values.confirmPassword) {
      nextErrors.confirmPassword = 'Confirm your password'
    } else if (values.confirmPassword !== values.password) {
      nextErrors.confirmPassword = 'Passwords do not match'
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

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  const trimmedEmail = nextFormData.email.trim()

  if (name === 'email') {
    if (trimmedEmail && !emailRegex.test(trimmedEmail)) {
      nextErrors.email = 'Enter a valid email'
    }
  }

  if (name === 'password' || name === 'confirmPassword') {
    if (nextFormData.password && nextFormData.password.length < 8) {
      nextErrors.password = 'Password must be at least 8 characters'
    } else {
      nextErrors.password = ''
    }

    if (nextFormData.confirmPassword) {
      if (nextFormData.confirmPassword !== nextFormData.password) {
        nextErrors.confirmPassword = 'Passwords do not match'
      } else {
        nextErrors.confirmPassword = ''
      }
    } else {
      nextErrors.confirmPassword = ''
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
      await registerUser({
        email: formData.email.trim(),
        password: formData.password,
        role: 'user',
      })

      setServerMessage('Registration completed successfully. Redirecting to sign in...')
      setFormData(initialFormData)
      setErrors(initialErrors)

      setTimeout(() => {
        navigate('/login')
      }, 1500)
    } catch (error) {
      setServerMessage(`Error: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="register-page">
      <div className="register-card">
        <div className="register-form-side">
          <div className="register-form-content">
            <h1 className="register-title">Create an Account</h1>
            <p className="register-subtitle">
              Join us and start your journey today.
            </p>

            <div className="social-buttons">
              <button type="button" className="social-button">
                <span className="social-icon social-icon--google">G</span>
                <span>Google</span>
              </button>

              <button type="button" className="social-button">
                <span className="social-icon social-icon--yandex">Я</span>
                <span>Yandex</span>
              </button>
            </div>

            <div className="register-divider">
              <span className="register-divider-line" />
              <span className="register-divider-text">Or</span>
              <span className="register-divider-line" />
            </div>

            <form className="register-form" onSubmit={handleSubmit} noValidate>
              <div className="form-group">
                <label className="form-label" htmlFor="login">
                  Login
                </label>
                <input
                  id="login"
                  className={`form-input ${errors.login ? 'form-input--error' : ''}`}
                  type="text"
                  name="login"
                  value={formData.login}
                  onChange={handleChange}
                  placeholder="IvanIvanov"
                  autoComplete="username"
                />
                {errors.login && (
                  <p className="form-error" role="alert">
                    {errors.login}
                  </p>
                )}
              </div>

              <div className="form-group">
                <label className="form-label" htmlFor="email">
                  Email
                </label>
                <input
                  id="email"
                  className={`form-input ${errors.email ? 'form-input--error' : ''}`}
                  type="email"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  placeholder="example@mail.com"
                  autoComplete="email"
                />
                {errors.email && (
                  <p className="form-error" role="alert">
                    {errors.email}
                  </p>
                )}
              </div>

              <div className="form-group">
                <label className="form-label" htmlFor="password">
                  Password
                </label>

                <div className="password-input-wrapper">
                  <input
                    id="password"
                    className={`form-input form-input--with-button ${errors.password ? 'form-input--error' : ''}`}
                    type={showPassword ? 'text' : 'password'}
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="••••••••"
                    autoComplete="new-password"
                  />
                  <button
                    type="button"
                    className="password-toggle"
                    onClick={() => setShowPassword((prev) => !prev)}
                  >
                    {showPassword ? 'Hide' : 'Show'}
                  </button>
                </div>

                {errors.password ? (
                  <p className="form-error" role="alert">
                    {errors.password}
                  </p>
                ) : (
                  <p className="form-hint">Must be at least 8 characters</p>
                )}
              </div>

              <div className="form-group">
                <label className="form-label" htmlFor="confirmPassword">
                  Confirm password
                </label>

                <div className="password-input-wrapper">
                  <input
                    id="confirmPassword"
                    className={`form-input form-input--with-button ${errors.confirmPassword ? 'form-input--error' : ''}`}
                    type={showConfirmPassword ? 'text' : 'password'}
                    name="confirmPassword"
                    value={formData.confirmPassword}
                    onChange={handleChange}
                    placeholder="••••••••"
                    autoComplete="new-password"
                  />
                  <button
                    type="button"
                    className="password-toggle"
                    onClick={() => setShowConfirmPassword((prev) => !prev)}
                  >
                    {showConfirmPassword ? 'Hide' : 'Show'}
                  </button>
                </div>

                {errors.confirmPassword && (
                  <p className="form-error" role="alert">
                    {errors.confirmPassword}
                  </p>
                )}
              </div>

              <button className="submit-button" type="submit" disabled={isLoading}>
                {isLoading ? 'Signing Up...' : 'Sign Up'}
              </button>
            </form>

            {serverMessage && (
              <p
                className={`server-message ${
                  serverMessage.startsWith('Error') ? 'server-message--error' : 'server-message--success'
                }`}
              >
                {serverMessage}
              </p>
            )}

            <p className="signin-text">
              Have an account? <Link to="/login">Sign in</Link>
            </p>
          </div>
        </div>

        <div className="register-visual-side" aria-hidden="true">
          <div className="visual-card">
            <span className="visual-glow visual-glow--top" />
            <span className="visual-glow visual-glow--center" />
            <span className="visual-glow visual-glow--bottom" />
          </div>
        </div>
      </div>
    </div>
  )
}

export default RegisterPage