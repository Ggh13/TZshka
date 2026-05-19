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

const USE_MOCK_PREVIEW = false

const modeOptions = [
  { label: 'Instant', value: 'Instant' },
  { label: 'Thinking', value: 'Thinking' },
]

const standardOptions = [
  { label: 'No standard', value: '' },
  { label: 'ГОСТ-19', value: 'ГОСТ-19' },
  { label: 'ГОСТ-34', value: 'ГОСТ-34' },
]

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function buildMockResponse(input, mode, standard) {
  const shortInput = input.trim().slice(0, 90) || 'Текст технического задания'

  return {
    rules_checker: {
      status: 'issues_found',
      feedback:
        mode === 'Thinking'
          ? 'Нашел несколько мест, которые могут трактоваться неоднозначно.'
          : 'Есть формулировки, которые стоит уточнить.',
      issues: [
        {
          rule_id: 'R1',
          problem: `Фраза "${shortInput}" слишком общая и не задает четкий измеримый результат.`,
          explanation:
            'Система или исполнитель не смогут однозначно понять, какой именно результат считается успешным. Лучше добавить конкретные критерии, сроки или числовые показатели.',
        },
        {
          rule_id: 'R5',
          problem:
            'В тексте нет явных ограничений и критериев приемки для результата.',
          explanation:
            'Без критериев приемки невозможно проверить, выполнено ли требование корректно. Стоит явно указать ожидаемый результат, формат проверки и допустимые отклонения.',
        },
      ],
    },
    standard_checker: standard
      ? {
          status: 'issues_found',
          feedback: `Проверка по ${standard} выявила неточности оформления и структуры.`,
          issues: [
            {
              rule_id: standard,
              problem:
                'Структура требования выглядит неполной для выбранного стандарта.',
              explanation:
                'Желательно разделить цель, функциональные требования, ограничения и критерии приемки по отдельным блокам, чтобы документ лучше соответствовал стандарту.',
            },
          ],
        }
      : null,
  }
}

function HomePage() {
  const navigate = useNavigate()
  const fileInputRef = useRef(null)
  const modeDropdownRef = useRef(null)
  const standardDropdownRef = useRef(null)
  const nextChatNumberRef = useRef(1)

  const [sessions, setSessions] = useState([])
  const [activeSessionId, setActiveSessionId] = useState(null)

  const [backendHistory, setBackendHistory] = useState([])
  const [mockHistoryMap, setMockHistoryMap] = useState({})

  const [historyLoading, setHistoryLoading] = useState(false)
  const [text, setText] = useState('')
  const [selectedFile, setSelectedFile] = useState(null)
  const [isSending, setIsSending] = useState(false)
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

  const displayedHistory = useMemo(() => {
    const localItems = activeSessionId ? mockHistoryMap[activeSessionId] || [] : []

    return [...backendHistory, ...localItems].sort((a, b) => {
      const first = new Date(a.createdAt).getTime()
      const second = new Date(b.createdAt).getTime()
      return first - second
    })
  }, [backendHistory, mockHistoryMap, activeSessionId])

  const hasAnswerBlock =
    historyLoading ||
    isSending ||
    Boolean(responseError) ||
    displayedHistory.length > 0

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

    try {
      const corrections = await getSessionHistory(sessionId)

      const sortedCorrections = [...corrections].sort((a, b) => {
        const first = new Date(a.createdAt).getTime()
        const second = new Date(b.createdAt).getTime()
        return first - second
      })

      setBackendHistory(sortedCorrections)
    } catch (error) {
      setBackendHistory([])
      if (!USE_MOCK_PREVIEW) {
        setResponseError(error.message)
      }
    } finally {
      setHistoryLoading(false)
    }
  }

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
      setBackendHistory([])
      setText('')
      setSelectedFile(null)
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
        setBackendHistory([])
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

    try {
      if (USE_MOCK_PREVIEW) {
        await sleep(selectedMode.value === 'Thinking' ? 2200 : 800)

        const inputContent = selectedFile
          ? `Файл: ${selectedFile.name}`
          : text.trim()

        const mockItem = {
          id: `mock-${Date.now()}`,
          sessionId: activeSessionId,
          inputContent,
          responseData: buildMockResponse(
            inputContent,
            selectedMode.value,
            selectedStandard.value
          ),
          createdAt: new Date().toISOString(),
        }

        setMockHistoryMap((prev) => ({
          ...prev,
          [activeSessionId]: [...(prev[activeSessionId] || []), mockItem],
        }))
      } else {
        if (selectedFile) {
          await sendFileToLLM({
            mode: selectedMode.value,
            standard: selectedStandard.value || undefined,
            file: selectedFile,
            sessionId: activeSessionId,
          })
        } else {
          await sendTextToLLM({
            mode: selectedMode.value,
            standard: selectedStandard.value || undefined,
            content: text.trim(),
            sessionId: activeSessionId,
          })
        }

        await loadHistory(activeSessionId)
      }

      setText('')
      setSelectedFile(null)

      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
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
      return <p className="workspace-answer-empty">Нет замечаний.</p>
    }

    return (
      <div className="workspace-answer-issues">
        {issues.map((issue, index) => (
          <div key={`${issue.rule_id || 'issue'}-${index}`} className="workspace-answer-issue">
            {issue.rule_id && (
              <div className="workspace-history-tag">{issue.rule_id}</div>
            )}

            <div className="workspace-answer-block">
              <div className="workspace-answer-block-title">Проблема</div>
              <div className="workspace-answer-problem">{issue.problem}</div>
            </div>

            <div className="workspace-answer-block">
              <div className="workspace-answer-block-title workspace-answer-block-title--secondary">
                Объяснение
              </div>
              <div className="workspace-answer-explanation">
                {issue.explanation}
              </div>
            </div>
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

  function renderThinkingState() {
    const title =
      selectedMode.value === 'Thinking' ? 'Думаю' : 'Проверяю'

    const subtext =
      selectedMode.value === 'Thinking'
        ? 'Анализирую текст и формирую ответ'
        : 'Обрабатываю запрос'

    return (
      <div className="workspace-thinking">
        <div className="workspace-thinking-line">
          <span className="workspace-thinking-text">{title}</span>
          <span className="workspace-thinking-dots">
            <span />
            <span />
            <span />
          </span>
        </div>
        <div className="workspace-thinking-subtext">{subtext}</div>
      </div>
    )
  }

  function renderHistory() {
    if (historyLoading) {
      return <div className="workspace-answer-state">Loading history...</div>
    }

    if (isSending) {
      return renderThinkingState()
    }

    if (responseError) {
      return (
        <div className="workspace-answer-state workspace-answer-state--error">
          {responseError}
        </div>
      )
    }

    if (displayedHistory.length === 0) {
      return <div className="workspace-answer-state">No response yet.</div>
    }

    return (
      <div className="workspace-history-list">
        {displayedHistory.map((item) => (
          <div key={item.id} className="workspace-history-item">
            <div className="workspace-history-label">Input</div>
            <div className="workspace-history-input">{item.inputContent}</div>

            <div className="workspace-history-label">Answer</div>
            {renderLLMBlock(item.responseData?.data || item.responseData)}
          </div>
        ))}
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