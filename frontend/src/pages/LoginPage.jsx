import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { loginUser } from '../api/auth'

function LoginPage() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  })

  const [message, setMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()

  function handleChange(event) {
    const { name, value } = event.target

    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setMessage('')
    setIsLoading(true)

    try {
    const result = await loginUser(formData)
    console.log('Login success:', result)

    localStorage.setItem('token', result.token)

    setMessage('Вход выполнен успешно')
    navigate('/')
    } catch (error) {
      console.error('Login error:', error.message)
      setMessage(`Ошибка: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div>
      <h1>Вход</h1>

      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '12px' }}>
          <label htmlFor="email">Email</label>
          <br />
          <input
            id="email"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="Введите email"
          />
        </div>

        <div style={{ marginBottom: '12px' }}>
          <label htmlFor="password">Пароль</label>
          <br />
          <input
            id="password"
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="Введите пароль"
          />
        </div>

        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Входим...' : 'Войти'}
        </button>
      </form>

      {message && <p style={{ marginTop: '16px' }}>{message}</p>}
    </div>
  )
}

export default LoginPage