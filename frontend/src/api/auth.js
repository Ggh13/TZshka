const BASE_URL = 'http://localhost:8080'

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

export async function registerUser({ login, email, password }) {
  const response = await fetch(`${BASE_URL}/api/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      login,
      email,
      password,
    }),
  })

  const result = await parseJson(response)

  if (!response.ok) {
    throw new Error(result.error?.message || 'Registration failed')
  }

  return result
}

export async function loginUser({ login, password }) {
  const response = await fetch(`${BASE_URL}/api/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      login,
      password,
    }),
  })

  const result = await parseJson(response)

  if (!response.ok) {
    throw new Error(result.error?.message || 'Login failed')
  }

  return result
}