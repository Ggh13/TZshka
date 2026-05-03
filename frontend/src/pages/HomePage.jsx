import { useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { sendFileToLLM, sendTextToLLM } from '../api/llm'
import './HomePage.css'

const modeOptions = [
  {
    label: 'Instant',
    description: 'Quick responses',
    value: 'Instant',
  },
  {
    label: 'Thinking',
    description: 'Solves complex tasks',
    value: 'Thinking',
  },
]

const standardOptions = [
  {
    label: 'No standard',
    description: 'Checking without standard',
    value: '',
  },
  {
    label: 'GOST-19',
    description: 'Checking according to GOST-19',
    value: 'GOST-19',
  },
  {
    label: 'GOST-2',
    description: 'Checking according to GOST-34',
    value: 'GOST-2',
  },
]

function createChat(id, title) {
  return {
    id,
    title,
    text: '',
    selectedFile: null,
    isSending: false,
    responseData: null,
    responseError: '',
    selectedMode: 'Instant',
    selectedStandard: 'GOST-19',
  }
}

function HomePage() {
  const navigate = useNavigate()
  const fileInputRef = useRef(null)
  const modeDropdownRef = useRef(null)
  const standardDropdownRef = useRef(null)
  const nextChatNumberRef = useRef(1)

  const [modeOpen, setModeOpen] = useState(false)
  const [standardOpen, setStandardOpen] = useState(false)
  const [chats, setChats] = useState([createChat(1, 'Chat')])
  const [activeChatId, setActiveChatId] = useState(1)

  const displayName = useMemo(() => {
    return (
      localStorage.getItem('userLogin') ||
      localStorage.getItem('userEmail') ||
      'Nickname'
    )
  }, [])

  const activeChat =
    chats.find((chat) => chat.id === activeChatId) || chats[0] || createChat(1, 'Chat')

  const hasAnswerBlock =
    activeChat.isSending ||
    Boolean(activeChat.responseData) ||
    Boolean(activeChat.responseError)

  useEffect(() => {
    function handleOutsideClick(event) {
      if (
        modeDropdownRef.current &&
        !modeDropdownRef.current.contains(event.target)
      ) {
        setModeOpen(false)
      }

      if (
        standardDropdownRef.current &&
        !standardDropdownRef.current.contains(event.target)
      ) {
        setStandardOpen(false)
      }
    }

    document.addEventListener('mousedown', handleOutsideClick)
    return () => {
      document.removeEventListener('mousedown', handleOutsideClick)
    }
  }, [])

  function updateActiveChat(patch) {
    setChats((prev) =>
      prev.map((chat) =>
        chat.id === activeChatId
          ? {
              ...chat,
              ...patch,
            }
          : chat
      )
    )
  }

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userEmail')
    localStorage.removeItem('userLogin')
    navigate('/login')
  }

  function handleAttachClick() {
    fileInputRef.current?.click()
  }

  function handleFileChange(event) {
    const file = event.target.files?.[0] || null
    updateActiveChat({
      selectedFile: file,
      responseError: '',
    })
  }

  function handleNewChat() {
    const id = Date.now()
    const title = `Chat ${nextChatNumberRef.current}`

    nextChatNumberRef.current += 1

    const newChat = createChat(id, title)

    setChats((prev) => [...prev, newChat])
    setActiveChatId(id)
    setModeOpen(false)
    setStandardOpen(false)

    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  function handleDeleteChat(chatId) {
    setChats((prev) => {
      const filtered = prev.filter((chat) => chat.id !== chatId)

      if (filtered.length === 0) {
        const fallbackChat = createChat(Date.now(), 'Chat')
        nextChatNumberRef.current = 1
        setActiveChatId(fallbackChat.id)
        return [fallbackChat]
      }

      if (chatId === activeChatId) {
        setActiveChatId(filtered[0].id)
      }

      return filtered
    })

    setModeOpen(false)
    setStandardOpen(false)

    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  async function handleSend() {
    if (!activeChat.text.trim() && !activeChat.selectedFile) {
      return
    }

    updateActiveChat({
      isSending: true,
      responseError: '',
      responseData: null,
    })

    try {
      let result

      if (activeChat.selectedFile) {
        result = await sendFileToLLM({
          mode: activeChat.selectedMode,
          standard: activeChat.selectedStandard || undefined,
          file: activeChat.selectedFile,
        })
      } else {
        result = await sendTextToLLM({
          mode: activeChat.selectedMode,
          standard: activeChat.selectedStandard || undefined,
          content: activeChat.text.trim(),
        })
      }

      updateActiveChat({
        isSending: false,
        responseData: result,
        responseError: '',
        text: '',
        selectedFile: null,
      })

      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    } catch (error) {
      updateActiveChat({
        isSending: false,
        responseError: error.message,
        responseData: null,
      })
    }
  }

  function renderDropdownMenu(type) {
    const options = type === 'mode' ? modeOptions : standardOptions

    return (
      <div className="workspace-dropdown-menu">
        <div className="workspace-dropdown-menu-title">Ambiguity Detector</div>

        {options.map((option) => (
          <button
            key={option.label}
            type="button"
            className="workspace-dropdown-option"
            onClick={() => {
              if (type === 'mode') {
                updateActiveChat({ selectedMode: option.value })
                setModeOpen(false)
              } else {
                updateActiveChat({ selectedStandard: option.value })
                setStandardOpen(false)
              }
            }}
          >
            <span className="workspace-dropdown-option-label">
              {option.label}
            </span>
            <span className="workspace-dropdown-option-description">
              {option.description}
            </span>
          </button>
        ))}
      </div>
    )
  }

  function renderIssues(issues) {
    if (!issues || issues.length === 0) {
      return <p className="workspace-answer-empty">No issues found.</p>
    }

    return (
      <div className="workspace-answer-issues">
        {issues.map((issue, index) => (
          <div key={`${issue.rule_id}-${index}`} className="workspace-answer-issue">
            <div className="workspace-answer-issue-title">{issue.rule_id}</div>
            <div className="workspace-answer-issue-text">{issue.problem}</div>
            <div className="workspace-answer-issue-text workspace-answer-issue-text--muted">
              {issue.explanation}
            </div>
          </div>
        ))}
      </div>
    )
  }

  function renderResponse() {
    if (activeChat.isSending) {
      return (
        <div className="workspace-answer-state">
          Processing your request...
        </div>
      )
    }

    if (activeChat.responseError) {
      return (
        <div className="workspace-answer-state workspace-answer-state--error">
          {activeChat.responseError}
        </div>
      )
    }

    if (!activeChat.responseData?.data) {
      return (
        <div className="workspace-answer-state">
          No response yet.
        </div>
      )
    }

    const rulesChecker = activeChat.responseData.data.rules_checker
    const standardChecker = activeChat.responseData.data.standard_checker

    return (
      <div className="workspace-answer-sections">
        {rulesChecker && (
          <section className="workspace-answer-section">
            <h3 className="workspace-answer-section-title">Rules check</h3>
            <p className="workspace-answer-status">
              Status: {rulesChecker.status}
            </p>
            <p className="workspace-answer-feedback">{rulesChecker.feedback}</p>
            {renderIssues(rulesChecker.issues)}
          </section>
        )}

        {standardChecker && (
          <section className="workspace-answer-section">
            <h3 className="workspace-answer-section-title">Standard check</h3>
            <p className="workspace-answer-status">
              Status: {standardChecker.status}
            </p>
            <p className="workspace-answer-feedback">{standardChecker.feedback}</p>
            {renderIssues(standardChecker.issues)}
          </section>
        )}
      </div>
    )
  }

  const selectedModeOption =
    modeOptions.find((option) => option.value === activeChat.selectedMode) || modeOptions[0]

  const selectedStandardOption =
    standardOptions.find((option) => option.value === activeChat.selectedStandard) || standardOptions[0]

  return (
    <div className="workspace-page">
      <aside className="workspace-sidebar">
        <button className="workspace-new-chat" type="button" onClick={handleNewChat}>
          <span className="workspace-new-chat__icon">✎</span>
          <span>New chat</span>
        </button>

        <nav className="workspace-sidebar-nav">
          <div className="workspace-sidebar-title">Chat</div>

          {chats.map((chat) => (
            <div
              key={chat.id}
              className={`workspace-chat-row ${chat.id === activeChatId ? 'workspace-chat-row--active' : ''}`}
            >
              <button
                className="workspace-chat-link"
                type="button"
                onClick={() => {
                  setActiveChatId(chat.id)
                  setModeOpen(false)
                  setStandardOpen(false)
                }}
              >
                {chat.title}
              </button>

              <button
                className="workspace-chat-delete"
                type="button"
                onClick={() => handleDeleteChat(chat.id)}
                aria-label={`Delete ${chat.title}`}
              >
                ×
              </button>
            </div>
          ))}
        </nav>

        <button className="workspace-logout" type="button" onClick={handleLogout}>
          Log out
        </button>
      </aside>

      <main className="workspace-main">
        <header className="workspace-header">
          <div className="workspace-profile">
            <span className="workspace-profile__name">{displayName}</span>
            <div className="workspace-profile__avatar" aria-hidden="true" />
          </div>
        </header>

        <section className={`workspace-content ${hasAnswerBlock ? 'workspace-content--with-answer' : ''}`}>
          {!hasAnswerBlock && <h1 className="workspace-title">Let’s get started!</h1>}

          <div className="workspace-input-card">
            <textarea
              className="workspace-textarea"
              placeholder="Enter the text"
              value={activeChat.text}
              onChange={(event) => updateActiveChat({ text: event.target.value })}
            />

            {activeChat.selectedFile && (
              <div className="workspace-file-chip">
                {activeChat.selectedFile.name}
                <button
                  type="button"
                  className="workspace-file-chip-remove"
                  onClick={() => {
                    updateActiveChat({ selectedFile: null })
                    if (fileInputRef.current) {
                      fileInputRef.current.value = ''
                    }
                  }}
                >
                  ×
                </button>
              </div>
            )}

            <div className="workspace-input-footer">
              <div className="workspace-left-controls">
                <input
                  ref={fileInputRef}
                  className="workspace-hidden-file-input"
                  type="file"
                  onChange={handleFileChange}
                />

                <button
                  className="workspace-attach-button"
                  type="button"
                  onClick={handleAttachClick}
                  aria-label="Attach file"
                >
                  📎
                </button>
              </div>

              <div className="workspace-right-controls">
                <div className="workspace-dropdown" ref={modeDropdownRef}>
                  <button
                    type="button"
                    className="workspace-dropdown-trigger"
                    onClick={() => {
                      setModeOpen((prev) => !prev)
                      setStandardOpen(false)
                    }}
                  >
                    <span>{selectedModeOption.label}</span>
                    <span className="workspace-dropdown-arrow">⌄</span>
                  </button>

                  {modeOpen && renderDropdownMenu('mode')}
                </div>

                <div className="workspace-dropdown" ref={standardDropdownRef}>
                  <button
                    type="button"
                    className="workspace-dropdown-trigger"
                    onClick={() => {
                      setStandardOpen((prev) => !prev)
                      setModeOpen(false)
                    }}
                  >
                    <span>{selectedStandardOption.label}</span>
                    <span className="workspace-dropdown-arrow">⌄</span>
                  </button>

                  {standardOpen && renderDropdownMenu('standard')}
                </div>

                <button
                  className="workspace-send-button"
                  type="button"
                  onClick={handleSend}
                  aria-label="Send"
                  disabled={activeChat.isSending}
                >
                  ➤
                </button>
              </div>
            </div>
          </div>

          {hasAnswerBlock && (
            <div className="workspace-answer-card">
              <div className="workspace-answer-card-title">Answer</div>
              <div className="workspace-answer-card-content">
                {renderResponse()}
              </div>
            </div>
          )}
        </section>
      </main>
    </div>
  )
}

export default HomePage