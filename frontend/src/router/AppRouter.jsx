import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import HomePage from '../pages/HomePage'
import LoginPage from '../pages/LoginPage'
import RegisterPage from '../pages/RegisterPage'

function AppRouter() {
  return (
    <BrowserRouter>
      <div style={{ padding: '24px', fontFamily: 'Arial, sans-serif' }}>
        <nav style={{ marginBottom: '24px', display: 'flex', gap: '16px' }}>
          <Link to="/">Главная</Link>
          <Link to="/login">Вход</Link>
          <Link to="/register">Регистрация</Link>
        </nav>

        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}

export default AppRouter