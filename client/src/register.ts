import { register } from './api'
import './styles.css'

const form = document.getElementById('register-form') as HTMLFormElement
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
  submitButton.textContent = 'Registering...'
  message.textContent = ''
  message.className = 'message'

  try {
    await register(username, password)
    showMessage('success', 'Registration successful. You can now log in.')
  } catch (error) {
    showMessage('error', error instanceof Error ? error.message : 'Registration failed.')
  } finally {
    submitButton.disabled = false
    submitButton.textContent = 'Register'
  }
})
