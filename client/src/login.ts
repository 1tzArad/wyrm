import { login } from './api'
import './styles.css'

const form = document.getElementById('login-form') as HTMLFormElement
const usernameInput = document.getElementById('username') as HTMLInputElement
const passwordInput = document.getElementById('password') as HTMLInputElement
const submitButton = document.getElementById('submit') as HTMLButtonElement
const message = document.getElementById('message') as HTMLParagraphElement

function showMessage(kind: 'success' | 'error', text: string) {
  message.textContent = text
  message.className = `message ${kind}`
}

form.addEventListener('submit', async (event) => {
  event.preventDefault()

  const username = usernameInput.value.trim()
  const password = passwordInput.value

  if (!username || !password) {
    showMessage('error', 'Username and password are required.')
    return
  }

  submitButton.disabled = true
  submitButton.textContent = 'Logging in...'
  message.textContent = ''
  message.className = 'message'

  try {
    await login(username, password)
    showMessage('success', 'Login successful. The JWT was stored in the "token" HttpOnly cookie.')
  } catch (error) {
    showMessage('error', error instanceof Error ? error.message : 'Login failed.')
  } finally {
    submitButton.disabled = false
    submitButton.textContent = 'Login'
  }
})
