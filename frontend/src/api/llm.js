const BASE_URL = 'http://localhost:8080'

function getAuthHeaders(extraHeaders = {}) {
  const token = localStorage.getItem('token')

  return {
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...extraHeaders,
  }
}

export async function sendTextToLLM({ mode, standard, content }) {
  const payload = {
    mode,
    content,
  }

  if (standard) {
    payload.standard = standard
  }

  const response = await fetch(`${BASE_URL}/api/llm/text`, {
    method: 'POST',
    headers: getAuthHeaders({
      'Content-Type': 'application/json',
    }),
    body: JSON.stringify(payload),
  })

  const result = await response.json()

  if (!response.ok || result.success === false) {
    throw new Error(result.error?.message || 'Text processing failed')
  }

  return result
}

export async function sendFileToLLM({ mode, standard, file }) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('mode', mode)

  if (standard) {
    formData.append('standard', standard)
  }

  const response = await fetch(`${BASE_URL}/api/llm/file`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: formData,
  })

  const result = await response.json()

  if (!response.ok || result.success === false) {
    throw new Error(result.error?.message || 'File processing failed')
  }

  return result
}