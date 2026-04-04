import { useNavigate } from 'react-router-dom'

function HomePage() {
  const navigate = useNavigate()
  const token = localStorage.getItem('token')

  function handleLogout() {
    localStorage.removeItem('token')
    navigate('/login')
  }

  return (
    <div>
      <h1>Главная страница</h1>

      {token ? (
        <>
          <p>Пользователь авторизован</p>
          <button onClick={handleLogout}>Выйти</button>
        </>
      ) : (
        <>
          <p>Вы не авторизованы</p>
          <button onClick={() => navigate('/login')}>Перейти ко входу</button>
        </>
      )}
    </div>
  )
}

export default HomePage