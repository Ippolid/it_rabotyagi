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

export function listQuestions(filters?: import('./data').QuestionFilters) {
  const params = new URLSearchParams();

  if (filters?.search) params.append('search', filters.search);
  if (filters?.technology) params.append('technology', filters.technology);
  if (filters?.difficulty) params.append('difficulty', filters.difficulty);
  if (filters?.company) params.append('company', filters.company);
  if (filters?.limit) params.append('limit', String(filters.limit));
  if (filters?.offset) params.append('offset', String(filters.offset));

  const query = params.toString();
  const path = query ? `/questions?${query}` : '/questions';

  return request<import('./data').QuestionListResponse>(path);
}

export function getQuestionById(id: number | string) {
  return request<import('./data').QuestionDetail>(`/questions/${id}`);
}

export function submitQuestionAnswer(id: number | string, userAnswer: string) {
  return request<{
    isCorrect: boolean;
    correctAnswer: string;
    explanation?: string;
    alreadySolved?: boolean;
  }>(`/questions/${id}/submit`, {
    method: 'POST',
    body: JSON.stringify({ userAnswer }),
    auth: true,
  });
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

// ===== ML Service API =====
// В продакшене используем nginx proxy, в разработке - прямой доступ
const ML_API_BASE = import.meta.env.DEV 
  ? 'http://localhost:8001/api/v1' 
  : '/ml-api';

export type MLQuestion = {
  id: number;
  question: string;
  ideal_answer: string;
};

export type TranscriptionResponse = {
  text: string;
  language?: string;
  duration?: number;
};

export type EvaluationResponse = {
  score: number;
  feedback: string;
  strengths: string[];
  improvements: string[];
  overall_comment: string;
};

export type CompleteInterviewResponse = {
  question_id?: number;
  question: string;
  transcription: TranscriptionResponse;
  evaluation: EvaluationResponse;
};

async function mlRequest<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${ML_API_BASE}${path}`, options);

  if (!res.ok) {
    let message = res.statusText;
    try {
      const data = await res.json() as { detail?: string };
      if (data?.detail) message = data.detail;
    } catch {
      // ignore parse errors
    }
    throw new Error(message);
  }

  return await res.json() as T;
}

export function getRandomMLQuestion() {
  return mlRequest<MLQuestion>('/questions/random');
}

export function getMLQuestionById(id: number) {
  return mlRequest<MLQuestion>(`/questions/${id}`);
}

export function getAllMLQuestions() {
  return mlRequest<MLQuestion[]>('/questions');
}

export async function completeInterview(
  audioBlob: Blob,
  questionId?: number,
  question?: string,
  idealAnswer?: string,
  language: string = 'ru'
): Promise<CompleteInterviewResponse> {
  const formData = new FormData();
  formData.append('audio', audioBlob, 'answer.webm');
  
  if (questionId !== undefined) {
    formData.append('question_id', String(questionId));
  } else if (question && idealAnswer) {
    formData.append('question', question);
    formData.append('ideal_answer', idealAnswer);
  }
  
  formData.append('language', language);

  return mlRequest<CompleteInterviewResponse>('/interview/complete', {
    method: 'POST',
    body: formData,
  });
}
