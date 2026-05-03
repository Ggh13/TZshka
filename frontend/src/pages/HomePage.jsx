import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './HomePage.css'

function HomePage() {
  const navigate = useNavigate()
  const storedEmail = localStorage.getItem('userEmail') || ''

  const nickname = useMemo(() => {
    if (!storedEmail) {
      return 'Nickname'
    }

    return storedEmail.split('@')[0]
  }, [storedEmail])

  const [text, setText] = useState('')
  const [mode, setMode] = useState('Mod')
  const [answer, setAnswer] = useState('Answer')

  function handleSend() {
    if (!text.trim()) {
      setAnswer(`Mode: ${mode}. Enter the text to get a response.`)
      return
    }

    setAnswer(`Mode: ${mode}\n\n${text}`)
  }

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userEmail')
    navigate('/')
  }

  return (
    <div className="home-page">
      <aside className="home-sidebar">
        <button className="home-new-chat" type="button">
          <span className="home-new-chat-icon">✎</span>
          <span>New chat</span>
        </button>

        <div className="home-chat-list">
          <button className="home-chat-item" type="button">
            Chat
          </button>
          <button className="home-chat-item" type="button">
            Chat 1
          </button>
          <button className="home-chat-item" type="button">
            Chat 2
          </button>
        </div>

        <button className="home-logout" type="button" onClick={handleLogout}>
          Log out
        </button>
      </aside>

      <div className="home-content">
        <header className="home-topbar">
          <div className="home-user-block">
            <span className="home-user-name">{nickname}</span>
            <div className="home-user-avatar" aria-hidden="true" />
          </div>
        </header>

        <main className="home-main">
          <section className="home-input-card">
            <textarea
              className="home-textarea"
              placeholder="Enter the text"
              value={text}
              onChange={(event) => setText(event.target.value)}
            />

            <div className="home-input-footer">
              <div className="home-mode-buttons">
                <button
                  type="button"
                  className={`home-mode-button ${mode === 'Mod' ? 'home-mode-button--active' : ''}`}
                  onClick={() => setMode('Mod')}
                >
                  Mod
                </button>

                <button
                  type="button"
                  className={`home-mode-button ${mode === 'Standart' ? 'home-mode-button--active' : ''}`}
                  onClick={() => setMode('Standart')}
                >
                  Standart
                </button>
              </div>

              <button className="home-send-button" type="button" onClick={handleSend}>
                →
              </button>
            </div>
          </section>

          <section className="home-answer-card">
            <div className="home-answer-body">
              <div className="home-answer-content">{answer}</div>
              <div className="home-answer-scroll">
                <div className="home-answer-scroll-thumb" />
              </div>
            </div>
          </section>
        </main>
      </div>
    </div>
  )
}

export default HomePage