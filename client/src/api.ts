import { API_BASE_URL } from './config'

export interface ApiError {
  code: string
  message: string
}

export interface ApiResponse<T> {
  success: boolean
  data: T | null
  error: ApiError | null
}

interface AuthRequest {
  username: string
  password: string
}

async function request<T>(path: string, body: AuthRequest): Promise<ApiResponse<T>> {
  const url = API_BASE_URL ? `${API_BASE_URL}${path}` : path

  let response: Response
  try {
    response = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(body),
    })
  } catch {
    throw new Error('Could not reach the backend. Make sure it is running and the API base URL is correct.')
  }

  let payload: ApiResponse<T>
  try {
    payload = (await response.json()) as ApiResponse<T>
  } catch {
    throw new Error(`Backend returned an invalid response (HTTP ${response.status}).`)
  }

  if (!payload.success) {
    const message = payload.error?.message
    throw new Error(message || `Request failed (HTTP ${response.status}).`)
  }

  return payload
}

export function register(username: string, password: string): Promise<ApiResponse<{ message: string }>> {
  return request('/register', { username, password })
}

export function login(username: string, password: string): Promise<ApiResponse<{ message: string }>> {
  return request('/login', { username, password })
}
