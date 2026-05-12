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

export async function getSessionHistory(sessionId) {
  const response = await fetch(`${BASE_URL}/api/history/${sessionId}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  })

  const result = await parseJson(response)

  if (!response.ok) {
    throw new Error(result.error?.message || 'Failed to load session history')
  }

  return result.corrections || []
}