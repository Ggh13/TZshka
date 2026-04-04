import { useState } from 'react'

function RegisterPage() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    role: 'user',
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
    console.log('Register form data:', formData)
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

        <button type="submit">Зарегистрироваться</button>
      </form>
    </div>
  )
}

export default RegisterPage