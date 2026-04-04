import { useState } from 'react'

function LoginPage() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  })

  function handleChange(event) {
    const { name, value } = event.target

    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  function handleSubmit(event) {
    event.preventDefault()
    console.log('Login form data:', formData)
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

        <button type="submit">Войти</button>
      </form>
    </div>
  )
}

export default LoginPage