export type AuthTokens = {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

const API_BASE = '/api/v1'

function getAccessToken() {
  return localStorage.getItem('access_token') || ''
}

function setTokens(tokens: AuthTokens) {
  localStorage.setItem('access_token', tokens.accessToken)
  localStorage.setItem('refresh_token', tokens.refreshToken)
  localStorage.setItem('expires_in', String(tokens.expiresIn))
}

export async function register(email: string, nickname: string, password: string) {
  const res = await fetch(`${API_BASE}/auth/register`, {
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
  const res = await fetch(`${API_BASE}/auth/login`, {
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
  const refreshToken = localStorage.getItem('refresh_token') || ''
  const res = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refreshToken })
  })
  if (!res.ok) throw new Error('Refresh failed')
  const data = (await res.json()) as AuthTokens
  setTokens(data)
  return data
}

export async function getMe() {
  const res = await fetch(`${API_BASE}/users/me`, {
    headers: { Authorization: `Bearer ${getAccessToken()}` },
  })
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
  const res = await fetch(`${API_BASE}/mentors`)
  if (!res.ok) throw new Error('Failed to load mentors')
  return res.json()
}

export type Course = {
  id: number
  title: string
  description: string
}

export async function listCourses(): Promise<{ items: Course[] }> {
  const res = await fetch(`${API_BASE}/courses`)
  if (!res.ok) throw new Error('Failed to load courses')
  return res.json()
}


