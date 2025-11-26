const API_BASE = '/api/v1';
const TOKEN_KEY = 'it_rabotyagi_tokens';

export type AuthTokens = {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
};

export type UserProfile = {
  id: number;
  email: string;
  fullName: string;
  role?: string;
};

type RequestOptions = RequestInit & {
  auth?: boolean;
};

type ApiError = {
  message?: string;
  code?: string;
};

export function getStoredTokens(): AuthTokens | null {
  try {
    const raw = localStorage.getItem(TOKEN_KEY);
    return raw ? (JSON.parse(raw) as AuthTokens) : null;
  } catch {
    return null;
  }
}

export function saveTokens(tokens: AuthTokens) {
  localStorage.setItem(TOKEN_KEY, JSON.stringify(tokens));
}

export function clearTokens() {
  localStorage.removeItem(TOKEN_KEY);
}

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const tokens = getStoredTokens();
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };

  if (options.auth && tokens?.accessToken) {
    headers.Authorization = `Bearer ${tokens.accessToken}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    let message = res.statusText;
    try {
      const data = (await res.json()) as ApiError;
      if (data?.message) message = data.message;
    } catch {
      // ignore parse errors
    }
    throw new Error(message);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return (await res.json()) as T;
}

export function login(email: string, password: string) {
  return request<AuthTokens>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export function register(email: string, password: string, nickname: string) {
  return request<AuthTokens>('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password, nickname }),
  });
}

export function getProfile() {
  return request<UserProfile>('/users/me', { auth: true });
}

export function listCourses() {
  return request<{ items: any[]; total?: number }>('/courses');
}

export function getCourseById(id: string | number) {
  return request(`/courses/${id}`);
}

export function getCourseModules(id: string | number) {
  return request<{ items: any[] }>(`/courses/${id}/modules`);
}

export function listMentors() {
  return request<{ items: any[]; total?: number }>('/mentors');
}

export function listQuestions() {
  return request<{ items: any[]; total?: number }>('/questions');
}
