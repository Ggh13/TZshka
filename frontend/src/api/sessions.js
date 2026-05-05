const BASE_URL = 'http://localhost:8080'

function getAuthHeaders(extraHeaders = {}) {
  const token = localStorage.getItem('token')

  return {
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...extraHeaders,
  }
}

async function parseJson(response) {
  const text = await response.text()

  if (!text) {
    return {}
  }

  try {
    return JSON.parse(text)
  } catch {
    throw new Error('Invalid server response')
  }
}

export async function getSessions() {
  const response = await fetch(`${BASE_URL}/api/sessions`, {
    method: 'GET',
    headers: getAuthHeaders(),
  })

  const result = await parseJson(response)

  if (!response.ok) {
    throw new Error(result.error?.message || 'Failed to load sessions')
  }

  return result.sessions || []
}

export async function createSession(name) {
  const response = await fetch(`${BASE_URL}/api/sessions/create`, {
    method: 'POST',
    headers: getAuthHeaders({
      'Content-Type': 'application/json',
    }),
    body: JSON.stringify({ name }),
  })

  const result = await parseJson(response)

  if (!response.ok) {
    throw new Error(result.error?.message || 'Failed to create session')
  }

  return result.session
}

export async function joinSession(sessionId) {
  const response = await fetch(`${BASE_URL}/api/sessions/join`, {
    method: 'POST',
    headers: getAuthHeaders({
      'Content-Type': 'application/json',
    }),
    body: JSON.stringify({ sessionId }),
  })

  const result = await parseJson(response)

  if (!response.ok) {
    throw new Error(result.error?.message || 'Failed to join session')
  }

  return result
}

export async function deleteSession(sessionId) {
  const response = await fetch(`${BASE_URL}/api/sessions`, {
    method: 'DELETE',
    headers: getAuthHeaders({
      'Content-Type': 'application/json',
    }),
    body: JSON.stringify({ sessionId }),
  })

  const result = await parseJson(response)

  if (!response.ok) {
    throw new Error(result.error?.message || 'Failed to delete session')
  }

  return result
}