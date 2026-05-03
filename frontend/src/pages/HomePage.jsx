import { useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { sendFileToLLM, sendTextToLLM } from '../api/llm'
import { getSessionHistory } from '../api/history'
import {
  createSession,
  deleteSession,
  getSessions,
  joinSession,
} from '../api/sessions'
import './HomePage.css'

const modeOptions = [
  { label: 'Instant', value: 'Instant' },
  { label: 'Thinking', value: 'Thinking' },
]

const standardOptions = [
  { label: 'No standard', value: '' },
  { label: 'ГОСТ-19', value: 'ГОСТ-19' },
  { label: 'ГОСТ-34', value: 'ГОСТ-34' },
]

function HomePage() {
  const navigate = useNavigate()
  const fileInputRef = useRef(null)
  const modeDropdownRef = useRef(null)
  const standardDropdownRef = useRef(null)
  const nextChatNumberRef = useRef(1)

  const [sessions, setSessions] = useState([])
  const [activeSessionId, setActiveSessionId] = useState(null)
  const [history, setHistory] = useState([])
  const [historyLoading, setHistoryLoading] = useState(false)

  const [text, setText] = useState('')
  const [selectedFile, setSelectedFile] = useState(null)
  const [isSending, setIsSending] = useState(false)
  const [responseData, setResponseData] = useState(null)
  const [responseError, setResponseError] = useState('')

  const [modeOpen, setModeOpen] = useState(false)
  const [standardOpen, setStandardOpen] = useState(false)
  const [selectedMode, setSelectedMode] = useState(modeOptions[0])
  const [selectedStandard, setSelectedStandard] = useState(standardOptions[1])

  const displayName = useMemo(() => {
    return (
      localStorage.getItem('userLogin') ||
      localStorage.getItem('userEmail') ||
      'Nickname'
    )
  }, [])

  const activeSession =
    sessions.find((session) => session.id === activeSessionId) || null

  const hasAnswerBlock =
    historyLoading ||
    isSending ||
    Boolean(responseData) ||
    Boolean(responseError) ||
    history.length > 0

  useEffect(() => {
    loadSessions()
  }, [])

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

  useEffect(() => {
    if (!activeSessionId) {
      return
    }

    loadHistory(activeSessionId)
  }, [activeSessionId])

  async function handleSelectSession(sessionId) {
  try {
    await joinSession(sessionId)
    setActiveSessionId(sessionId)
    setModeOpen(false)
    setStandardOpen(false)
    setResponseError('')
  } catch (error) {
    setResponseError(error.message)
  }
}

  async function loadSessions() {
    try {
      const loadedSessions = await getSessions()

      if (loadedSessions.length === 0) {
        const created = await createSession('Chat')
        setSessions([created])
        setActiveSessionId(created.id)
        nextChatNumberRef.current = 1
        return
      }

      setSessions(loadedSessions)
      setActiveSessionId(loadedSessions[0].id)

      const numbered = loadedSessions
        .map((session) => {
          const match = session.name.match(/^Chat\s+(\d+)$/)
          return match ? Number(match[1]) : 0
        })
        .filter(Boolean)

      nextChatNumberRef.current =
        numbered.length > 0 ? Math.max(...numbered) + 1 : 1
    } catch (error) {
      setResponseError(error.message)
    }
  }

  async function loadHistory(sessionId) {
    setHistoryLoading(true)
    setResponseError('')
    setResponseData(null)

    try {
      const corrections = await getSessionHistory(sessionId)
      setHistory(corrections)
    } catch (error) {
      setHistory([])
      setResponseError(error.message)
    } finally {
      setHistoryLoading(false)
    }
  }

  function handleLogout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userEmail')
    localStorage.removeItem('userLogin')
    navigate('/login')
  }

  async function handleNewChat() {
    try {
      const name =
        nextChatNumberRef.current === 1
          ? 'Chat 1'
          : `Chat ${nextChatNumberRef.current}`

      const created = await createSession(name)

      nextChatNumberRef.current += 1

      setSessions((prev) => [...prev, created])
      setActiveSessionId(created.id)
      setHistory([])
      setText('')
      setSelectedFile(null)
      setResponseData(null)
      setResponseError('')
      setModeOpen(false)
      setStandardOpen(false)

      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    } catch (error) {
      setResponseError(error.message)
    }
  }

  async function handleDeleteSession(sessionId) {
    try {
      await deleteSession(sessionId)

      const nextSessions = sessions.filter((session) => session.id !== sessionId)

      if (nextSessions.length === 0) {
        const created = await createSession('Chat')
        setSessions([created])
        setActiveSessionId(created.id)
        setHistory([])
        return
      }

      setSessions(nextSessions)

      if (sessionId === activeSessionId) {
        setActiveSessionId(nextSessions[0].id)
      }
    } catch (error) {
      setResponseError(error.message)
    }
  }

  function handleAttachClick() {
    fileInputRef.current?.click()
  }

  function handleFileChange(event) {
    const file = event.target.files?.[0] || null
    setSelectedFile(file)
    setResponseError('')
  }

  async function handleSend() {
    if (!activeSessionId || (!text.trim() && !selectedFile)) {
      return
    }

    setIsSending(true)
    setResponseError('')
    setResponseData(null)

    try {
      let result

      if (selectedFile) {
        result = await sendFileToLLM({
          mode: selectedMode.value,
          standard: selectedStandard.value || undefined,
          file: selectedFile,
          sessionId: activeSessionId,
        })
      } else {
        result = await sendTextToLLM({
          mode: selectedMode.value,
          standard: selectedStandard.value || undefined,
          content: text.trim(),
          sessionId: activeSessionId,
        })
      }

      setResponseData(result)
      setText('')
      setSelectedFile(null)

      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }

      await loadHistory(activeSessionId)
    } catch (error) {
      setResponseError(error.message)
    } finally {
      setIsSending(false)
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
                setSelectedMode(option)
                setModeOpen(false)
              } else {
                setSelectedStandard(option)
                setStandardOpen(false)
              }
            }}
          >
            <span className="workspace-dropdown-option-label">
              {option.label}
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
          <div key={`${issue.rule_id || 'issue'}-${index}`} className="workspace-answer-issue">
            {issue.rule_id && (
              <div className="workspace-answer-issue-title">{issue.rule_id}</div>
            )}

            {issue.problem && (
              <div className="workspace-answer-issue-text">{issue.problem}</div>
            )}

            {issue.explanation && (
              <div className="workspace-answer-issue-text workspace-answer-issue-text--muted">
                {issue.explanation}
              </div>
            )}
          </div>
        ))}
      </div>
    )
  }

  function renderLLMBlock(data) {
    const rulesChecker = data?.rules_checker
    const standardChecker = data?.standard_checker

    return (
      <div className="workspace-answer-sections">
        {rulesChecker && (
          <section className="workspace-answer-section">
            <h3 className="workspace-answer-section-title">Rules check</h3>

            {typeof rulesChecker === 'string' ? (
              <p className="workspace-answer-feedback">{rulesChecker}</p>
            ) : (
              <>
                {rulesChecker.status && (
                  <p className="workspace-answer-status">
                    Status: {rulesChecker.status}
                  </p>
                )}
                {rulesChecker.feedback && (
                  <p className="workspace-answer-feedback">{rulesChecker.feedback}</p>
                )}
                {renderIssues(rulesChecker.issues)}
              </>
            )}
          </section>
        )}

        {standardChecker && (
          <section className="workspace-answer-section">
            <h3 className="workspace-answer-section-title">Standard check</h3>

            {typeof standardChecker === 'string' ? (
              <p className="workspace-answer-feedback">{standardChecker}</p>
            ) : (
              <>
                {standardChecker.status && (
                  <p className="workspace-answer-status">
                    Status: {standardChecker.status}
                  </p>
                )}
                {standardChecker.feedback && (
                  <p className="workspace-answer-feedback">{standardChecker.feedback}</p>
                )}
                {renderIssues(standardChecker.issues)}
              </>
            )}
          </section>
        )}
      </div>
    )
  }

  function renderHistory() {
    if (historyLoading) {
      return <div className="workspace-answer-state">Loading history...</div>
    }

    if (responseError) {
      return (
        <div className="workspace-answer-state workspace-answer-state--error">
          {responseError}
        </div>
      )
    }

    if (history.length === 0 && !responseData) {
      return <div className="workspace-answer-state">No response yet.</div>
    }

    return (
      <div className="workspace-history-list">
        {history.map((item) => (
          <div key={item.id} className="workspace-history-item">
            <div className="workspace-history-label">Input</div>
            <div className="workspace-history-input">{item.inputContent}</div>

            <div className="workspace-history-label">Answer</div>
            {renderLLMBlock(item.responseData)}
          </div>
        ))}

        {responseData?.data && (
          <div className="workspace-history-item">
            <div className="workspace-history-label">Latest answer</div>
            {renderLLMBlock(responseData.data)}
          </div>
        )}
      </div>
    )
  }

  return (
    <div className="workspace-page">
      <aside className="workspace-sidebar">
        <button className="workspace-new-chat" type="button" onClick={handleNewChat}>
          <span className="workspace-new-chat__icon">✎</span>
          <span>New chat</span>
        </button>

        <nav className="workspace-sidebar-nav">
          <div className="workspace-sidebar-title">Chat</div>

          {sessions.map((session) => (
            <div
              key={session.id}
              className={`workspace-chat-row ${
                session.id === activeSessionId ? 'workspace-chat-row--active' : ''
              }`}
            >
              <button
                className="workspace-chat-link"
                type="button"
                onClick={() => handleSelectSession(session.id)}
              >
                {session.name}
              </button>

              <button
                className="workspace-chat-delete"
                type="button"
                onClick={() => handleDeleteSession(session.id)}
                aria-label={`Delete ${session.name}`}
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
              value={text}
              onChange={(event) => setText(event.target.value)}
            />

            {selectedFile && (
              <div className="workspace-file-chip">
                {selectedFile.name}
                <button
                  type="button"
                  className="workspace-file-chip-remove"
                  onClick={() => {
                    setSelectedFile(null)
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
                  accept=".txt"
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
                    <span>{selectedMode.label}</span>
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
                    <span>{selectedStandard.label}</span>
                    <span className="workspace-dropdown-arrow">⌄</span>
                  </button>

                  {standardOpen && renderDropdownMenu('standard')}
                </div>

                <button
                  className="workspace-send-button"
                  type="button"
                  onClick={handleSend}
                  aria-label="Send"
                  disabled={isSending || !activeSession}
                >
                  ➤
                </button>
              </div>
            </div>
          </div>

          {hasAnswerBlock && (
            <div className="workspace-answer-card">
              <div className="workspace-answer-card-title">
                {activeSession ? activeSession.name : 'Answer'}
              </div>
              <div className="workspace-answer-card-content">{renderHistory()}</div>
            </div>
          )}
        </section>
      </main>
    </div>
  )
}

export default HomePage