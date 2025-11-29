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
  avatarUrl?: string;
  description?: string;
};

// Statistics types
export type QuestionStatisticsSummary = {
  totalSolved: number;
  correctAnswersPct: number;
  totalAttempts: number;
  avgTimeSpent: number;
};

export type UserStatistics = {
  coursesEnrolled: number;
  coursesCompleted: number;
  overallProgressPct: number;
  questionsStatistics: QuestionStatisticsSummary;
};

export type CourseStatisticsItem = {
  courseId: number;
  courseTitle: string;
  totalModules: number;
  completedModules: number;
  progressPct: number;
  startedAt: string;
  timeSpent: number;
};

export type QuestionStatisticsByCourse = {
  courseId: number;
  courseTitle: string;
  statistics: QuestionStatisticsSummary;
};

export type UserQuestionStatistics = {
  overall: QuestionStatisticsSummary;
  byCourse: QuestionStatisticsByCourse[];
};

// Profile update types
export type UserProfileUpdate = {
  username?: string;
  name?: string;
  email?: string;
  description?: string;
};

export type UserProfileUpdateResponse = {
  profile: UserProfile;
  tokens?: AuthTokens;
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

// Statistics API
export function getUserStatistics() {
  return request<UserStatistics>('/users/me/statistics', { auth: true });
}

export function getCourseStatistics() {
  return request<{ items: CourseStatisticsItem[] }>('/users/me/statistics/courses', { auth: true });
}

export function getQuestionStatistics(courseId?: number) {
  const query = courseId ? `?courseId=${courseId}` : '';
  return request<UserQuestionStatistics>(`/users/me/statistics/questions${query}`, { auth: true });
}

// Profile management API
export function updateProfile(data: UserProfileUpdate) {
  return request<UserProfileUpdateResponse>('/users/me/profile', {
    method: 'PATCH',
    body: JSON.stringify(data),
    auth: true,
  });
}

export function updateAvatar(avatarUrl: string) {
  return request<{ avatarUrl: string }>('/users/me/avatar', {
    method: 'PATCH',
    body: JSON.stringify({ avatarUrl }),
    auth: true,
  });
}

export function changePassword(oldPassword: string, newPassword: string) {
  return request<{ message: string }>('/users/me/password', {
    method: 'POST',
    body: JSON.stringify({ oldPassword, newPassword }),
    auth: true,
  });
}
