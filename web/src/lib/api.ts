export type AuthTokens = {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

const API_BASE = '/api/v1'

function getAccessToken() {
  return localStorage.getItem('access_token') || ''
}

function getRefreshToken() {
  return localStorage.getItem('refresh_token') || ''
}

function setTokens(tokens: AuthTokens) {
  localStorage.setItem('access_token', tokens.accessToken)
  localStorage.setItem('refresh_token', tokens.refreshToken)
  localStorage.setItem('expires_in', String(tokens.expiresIn))
}

function clearTokens() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('expires_in')
}

// Internal helper: fetch with optional auth and auto-refresh on 401
async function apiFetch(input: string, init: RequestInit = {}, opts: { auth?: boolean; retry?: boolean } = {}) {
  const { auth = false, retry = true } = opts
  const headers: Record<string, string> = {
    Accept: 'application/json',
    ...(init.headers as Record<string, string>),
  }
  if (auth) {
    const token = getAccessToken()
    if (token) headers.Authorization = `Bearer ${token}`
  }
  const res = await fetch(`${API_BASE}${input}`, { ...init, headers })
  if (res.status === 401 && retry && getRefreshToken()) {
    // try refresh once then retry original request
    try {
      await refresh()
      return apiFetch(input, init, { auth, retry: false })
    } catch (e) {
      clearTokens()
      throw e
    }
  }
  return res
}

export async function register(email: string, nickname: string, password: string) {
  const res = await apiFetch('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, nickname, password })
  })
  if (!res.ok) throw new Error('Register failed')
  const data = (await res.json()) as AuthTokens
  setTokens(data)
  return data
}

export async function login(email: string, password: string) {
  const res = await apiFetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  })
  if (!res.ok) throw new Error('Login failed')
  const data = (await res.json()) as AuthTokens
  setTokens(data)
  return data
}

export async function refresh() {
  const refreshToken = getRefreshToken()
  if (!refreshToken) throw new Error('No refresh token')
  const res = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
    body: JSON.stringify({ refreshToken })
  })
  if (!res.ok) throw new Error('Refresh failed')
  const data = (await res.json()) as AuthTokens
  setTokens(data)
  return data
}

export function logout() {
  // Server has no dedicated logout route registered yet; perform client-side logout
  clearTokens()
}

export async function getMe() {
  const res = await apiFetch('/users/me', { method: 'GET' }, { auth: true })
  if (!res.ok) throw new Error('Unauthorized')
  return res.json()
}

export type MentorCard = {
  id: number
  fullName: string
  title: string
  skills: string[]
  yearsOfExperience?: number
}

export async function listMentors(): Promise<{ items: MentorCard[]; total?: number }> {
  const res = await apiFetch('/mentors', { method: 'GET' })
  if (!res.ok) throw new Error('Failed to load mentors')
  return res.json()
}

export type Course = {
  id: number
  title: string
  description: string
}

export async function listCourses(): Promise<{ items: Course[] }> {
  const res = await apiFetch('/courses', { method: 'GET' })
  if (!res.ok) throw new Error('Failed to load courses')
  return res.json()
}

export async function getCourse(id: number): Promise<Course> {
  const res = await apiFetch(`/courses/${id}`, { method: 'GET' })
  if (!res.ok) throw new Error('Failed to load course')
  return res.json()
}

export type Module = {
  id: number
  courseId: number
  title: string
  description: string
  moduleOrder: number
}

export async function listModules(courseId: number): Promise<{ items: Module[] }> {
  const res = await apiFetch(`/courses/${courseId}/modules`, { method: 'GET' })
  if (!res.ok) throw new Error('Failed to load modules')
  return res.json()
}

export type Question = {
  id: number
  title: string
  content: string
  difficulty?: string
  options?: string[]
  correctAnswer?: string
  explanation?: string
}

export async function listQuestions(moduleId: number): Promise<{ items: Question[] }> {
  const res = await apiFetch(`/modules/${moduleId}/questions`, { method: 'GET' })
  if (!res.ok) throw new Error('Failed to load questions')
  return res.json()
}

export async function getModule(id: number): Promise<Module> {
  const res = await apiFetch(`/modules/${id}`, { method: 'GET' })
  if (!res.ok) throw new Error('Failed to load module')
  return res.json()
}

export async function getQuestionByModule(moduleId: number, questionId: number): Promise<Question | null> {
  const { items } = await listQuestions(moduleId)
  return items.find(q => q.id === questionId) ?? null
}


