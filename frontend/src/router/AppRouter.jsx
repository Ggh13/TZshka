import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import HomePage from '../pages/HomePage'
import LoginPage from '../pages/LoginPage'
import RegisterPage from '../pages/RegisterPage'
import ProtectedRoute from './ProtectedRoute'
import PublicRoute from './PublicRoute'

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
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <HomePage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/login"
            element={
              <PublicRoute>
                <LoginPage />
              </PublicRoute>
            }
          />

          <Route
            path="/register"
            element={
              <PublicRoute>
                <RegisterPage />
              </PublicRoute>
            }
          />
        </Routes>
      </div>
    </BrowserRouter>
  )
}

export default AppRouter