import { useState } from 'react'
import { registerUser } from '../api/auth'

function RegisterPage() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    role: 'user',
  })

  const [message, setMessage] = useState('')
  const [isLoading, setIsLoading] = useState(false)

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
      const result = await registerUser(formData)
      console.log('Register success:', result)
      setMessage('Регистрация прошла успешно')
    } catch (error) {
      console.error('Register error:', error.message)
      setMessage(`Ошибка: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div>
      <h1>Регистрация</h1>

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

        <div style={{ marginBottom: '12px' }}>
          <label htmlFor="role">Роль</label>
          <br />
          <select
            id="role"
            name="role"
            value={formData.role}
            onChange={handleChange}
          >
            <option value="user">user</option>
            <option value="admin">admin</option>
          </select>
        </div>

        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Регистрируем...' : 'Зарегистрироваться'}
        </button>
      </form>

      {message && <p style={{ marginTop: '16px' }}>{message}</p>}
    </div>
  )
}

export default RegisterPage