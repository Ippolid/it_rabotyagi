import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Progress } from '../ui/progress';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Textarea } from '../ui/textarea';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs';
import { Separator } from '../ui/separator';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Activity, BookOpen, Clock, Star, Trophy, User, Settings } from 'lucide-react';
import {
  getUserStatistics,
  getCourseStatistics,
  getQuestionStatistics,
  getProfile,
  updateProfile,
  updateAvatar,
  changePassword,
  saveTokens,
  type UserStatistics,
  type CourseStatisticsItem,
  type UserQuestionStatistics,
  type UserProfile,
  type UserProfileUpdate,
} from '../../lib/api';

type DashboardProps = { onNavigate: (view: any) => void; userName?: string | null };

export function Dashboard({ onNavigate, userName }: DashboardProps) {
  const [activeTab, setActiveTab] = useState('overview');
  const [stats, setStats] = useState<UserStatistics | null>(null);
  const [courseStats, setCourseStats] = useState<CourseStatisticsItem[]>([]);
  const [questionStats, setQuestionStats] = useState<UserQuestionStatistics | null>(null);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadAllData();
  }, []);

  const loadAllData = async () => {
    let cancelled = false;
    setLoading(true);
    try {
      const [statsRes, courseStatsRes, questionStatsRes, profileRes] = await Promise.all([
        getUserStatistics().catch(() => null),
        getCourseStatistics().catch(() => ({ items: [] })),
        getQuestionStatistics().catch(() => null),
        getProfile().catch(() => null),
      ]);
      if (cancelled) return;

      setStats(statsRes || {
        coursesEnrolled: 0,
        coursesCompleted: 0,
        overallProgressPct: 0,
        questionsStatistics: {
          totalSolved: 0,
          correctAnswersPct: 0,
          totalAttempts: 0,
          avgTimeSpent: 0,
        },
      });
      setCourseStats(courseStatsRes?.items || []);
      setQuestionStats(questionStatsRes || {
        overall: {
          totalSolved: 0,
          correctAnswersPct: 0,
          totalAttempts: 0,
          avgTimeSpent: 0,
        },
        byCourse: [],
      });
      setProfile(profileRes);
    } catch (err) {
      // Keep fallback data
    } finally {
      if (!cancelled) setLoading(false);
    }
    return () => { cancelled = true; };
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin h-8 w-8 border-4 border-blue-600 border-t-transparent rounded-full" />
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      {/* Profile Header with Avatar */}
      <div className="flex items-center gap-4 p-6 bg-white rounded-xl shadow-sm border border-gray-100">
        <Avatar className="h-16 w-16">
          <AvatarImage src={profile?.avatarUrl || (profile?.email ? `https://api.dicebear.com/7.x/avataaars/svg?seed=${profile.email}` : undefined)} />
          <AvatarFallback className="text-lg font-bold bg-blue-100 text-blue-700">
            {profile?.fullName ? profile.fullName.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) :
             userName ? userName.slice(0, 2).toUpperCase() : 'У'}
          </AvatarFallback>
        </Avatar>
        <div className="flex-1">
          <h1 className="text-xl font-bold text-gray-900">{profile?.fullName || userName || 'Пользователь'}</h1>
          <p className="text-sm text-gray-500">{profile?.email || 'Email не указан'}</p>
          {profile?.description && (
            <p className="text-xs text-gray-400 mt-1">{profile.description}</p>
          )}
          <div className="flex gap-2 mt-2">
            <Badge variant="secondary" className="text-xs">{profile?.role === 'admin' ? 'Администратор' : 'Пользователь'}</Badge>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="overview">Обзор</TabsTrigger>
          <TabsTrigger value="settings">
            <Settings className="h-4 w-4 mr-2" />
            Настройки профиля
          </TabsTrigger>
        </TabsList>

        <TabsContent value="overview">
          <OverviewTab
            stats={stats}
            courseStats={courseStats}
            questionStats={questionStats}
            onNavigate={onNavigate}
            userName={userName}
            profile={profile}
          />
        </TabsContent>

        <TabsContent value="settings">
          <SettingsTab profile={profile} onProfileUpdated={loadAllData} />
        </TabsContent>
      </Tabs>
    </div>
  );
}

// Overview Tab - похож на старый Dashboard
function OverviewTab({
  stats,
  courseStats,
  questionStats,
  onNavigate,
  userName,
  profile
}: {
  stats: UserStatistics | null;
  courseStats: CourseStatisticsItem[];
  questionStats: UserQuestionStatistics | null;
  onNavigate: (view: any) => void;
  userName?: string | null;
  profile: UserProfile | null;
}) {
  if (!stats) return null;

  // Use only real API data - no fallback data
  const learningTimeHours = courseStats.reduce((acc, c) => acc + c.timeSpent, 0);

  return (
    <div className="space-y-8 pt-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Дашборд</h1>
          <p className="text-gray-500 mt-1">
            С возвращением, {userName || profile?.fullName || 'Пользователь'}!
          </p>
        </div>
        <Button onClick={() => onNavigate('courses')}>Найти курс</Button>
      </div>

      {/* Stats Overview - данные из API */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card className="bg-gray-900 text-white border-none">
          <CardContent className="p-6 flex flex-col justify-between h-full">
            <Activity className="h-8 w-8 opacity-80" />
            <div>
              <div className="text-3xl font-bold">{learningTimeHours.toFixed(0)}ч</div>
              <div className="text-gray-400 text-sm">Время обучения</div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6 flex flex-col justify-between h-full">
            <BookOpen className="h-8 w-8 text-blue-600" />
            <div>
              <div className="text-3xl font-bold">{stats.coursesEnrolled}</div>
              <div className="text-gray-500 text-sm">Записано курсов</div>
              <div className="text-xs text-gray-400 mt-1">
                Завершено: {stats.coursesCompleted}
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6 flex flex-col justify-between h-full">
            <Trophy className="h-8 w-8 text-amber-500" />
            <div>
              <div className="text-3xl font-bold">{stats.overallProgressPct.toFixed(0)}%</div>
              <div className="text-gray-500 text-sm">Общий прогресс</div>
              <div className="text-xs text-gray-400 mt-1">
                По всем курсам
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6 flex flex-col justify-between h-full">
            <Star className="h-8 w-8 text-purple-500" />
            <div>
              <div className="text-3xl font-bold">{stats.questionsStatistics.totalSolved}</div>
              <div className="text-gray-500 text-sm">Вопросов решено</div>
              <div className="text-xs text-gray-400 mt-1">
                {stats.questionsStatistics.correctAnswersPct.toFixed(0)}% правильно
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Active Learning */}
        <div className="lg:col-span-2 space-y-6">
          <h2 className="text-xl font-bold text-gray-900">Продолжить обучение</h2>
          <div className="space-y-4">
            {courseStats.map((course) => (
              <Card
                key={course.courseId}
                className="overflow-hidden hover:bg-gray-50 transition-colors cursor-pointer border-gray-200"
                onClick={() => onNavigate('courses')}
              >
                <div className="flex flex-col sm:flex-row">
                  <div className="w-full sm:w-48 h-32 sm:h-auto relative bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                    <div className="text-white text-4xl font-bold opacity-80">
                      {course.courseTitle.charAt(0)}
                    </div>
                  </div>
                  <div className="flex-1 p-6 flex flex-col justify-between">
                    <div className="mb-4">
                      <div className="flex justify-between items-start mb-2">
                        <h3 className="font-bold text-gray-900">{course.courseTitle}</h3>
                      </div>
                      <div className="flex items-center gap-4 text-sm text-gray-500">
                        <div className="flex items-center gap-1">
                          <Clock size={14} />
                          <span>{course.timeSpent?.toFixed(1) || '0'}ч потрачено</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <Star size={14} className="text-amber-500" />
                          <span>Начат: {new Date(course.startedAt).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short', year: 'numeric' })}</span>
                        </div>
                      </div>
                    </div>
                    <div className="space-y-2">
                      <div className="flex justify-between text-xs font-medium text-gray-500">
                        <span>Модулей: {course.completedModules || 0}/{course.totalModules || 0}</span>
                        <span>{course.progressPct?.toFixed(0) || 0}%</span>
                      </div>
                      <Progress value={course.progressPct || 0} className="h-2" />
                    </div>
                  </div>
                </div>
              </Card>
            ))}

            {courseStats.length === 0 && (
              <div className="text-center py-12 bg-white border border-dashed border-gray-300 rounded-xl">
                <BookOpen className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500 mb-4">Вы еще не начали ни один курс.</p>
                <Button variant="outline" onClick={() => onNavigate('courses')}>
                  Открыть каталог
                </Button>
              </div>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card className="border-none bg-blue-50">
            <CardHeader>
              <CardTitle className="text-blue-900">Ежедневный челлендж</CardTitle>
              <CardDescription className="text-blue-700">
                Решите случайную задачу по алгоритмам, чтобы держать навыки в тонусе.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="bg-white p-4 rounded-lg shadow-sm mb-4">
                <h4 className="font-bold text-gray-900 mb-2">Инвертировать бинарное дерево</h4>
                <div className="flex gap-2 mb-2">
                  <Badge variant="secondary" className="text-xs">Легкая</Badge>
                  <Badge variant="secondary" className="text-xs">Деревья</Badge>
                </div>
              </div>
              <Button className="w-full bg-blue-600 hover:bg-blue-700">
                Начать челлендж
              </Button>
            </CardContent>
          </Card>

          <div className="space-y-4">
            <h3 className="font-bold text-gray-900">Статистика по курсам</h3>
            {questionStats && questionStats.byCourse.length > 0 ? (
              questionStats.byCourse.slice(0, 3).map((course) => (
                <Card
                  key={course.courseId}
                  className="cursor-pointer hover:shadow-md transition-shadow"
                  onClick={() => onNavigate('courses')}
                >
                  <CardContent className="p-4">
                    <h4 className="font-medium text-sm text-gray-900 mb-2 line-clamp-1">
                      {course.courseTitle}
                    </h4>
                    <div className="grid grid-cols-2 gap-2 text-xs">
                      <div>
                        <span className="text-gray-500">Решено:</span>
                        <span className="font-bold text-gray-900 ml-1">{course.statistics.totalSolved}</span>
                      </div>
                      <div>
                        <span className="text-gray-500">Правильно:</span>
                        <span className="font-bold text-green-600 ml-1">{course.statistics.correctAnswersPct.toFixed(0)}%</span>
                      </div>
                      <div>
                        <span className="text-gray-500">Попыток:</span>
                        <span className="font-bold text-gray-900 ml-1">{course.statistics.totalAttempts}</span>
                      </div>
                      <div>
                        <span className="text-gray-500">~Время:</span>
                        <span className="font-bold text-blue-600 ml-1">{course.statistics.avgTimeSpent.toFixed(0)}с</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))
            ) : (
              <Card className="cursor-pointer hover:shadow-md transition-shadow">
                <CardContent className="p-4 text-center">
                  <p className="text-sm text-gray-500">Нет статистики по вопросам</p>
                  <Button
                    variant="outline"
                    size="sm"
                    className="mt-2"
                    onClick={() => onNavigate('questions')}
                  >
                    Начать решать
                  </Button>
                </CardContent>
              </Card>
            )}

            {questionStats && questionStats.overall.totalSolved > 0 && (
              <Card className="bg-green-50 border-green-200">
                <CardContent className="p-4">
                  <h4 className="font-medium text-sm text-green-900 mb-2">
                    Общая статистика
                  </h4>
                  <div className="text-2xl font-bold text-green-700">
                    {questionStats.overall.correctAnswersPct.toFixed(1)}%
                  </div>
                  <p className="text-xs text-green-600">
                    правильных ответов из {questionStats.overall.totalAttempts} попыток
                  </p>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

// Settings Tab
function SettingsTab({ profile, onProfileUpdated }: { profile: UserProfile | null; onProfileUpdated: () => void }) {
  const [profileForm, setProfileForm] = useState({
    username: '',
    name: '',
    email: '',
    description: '',
  });
  const [avatarUrl, setAvatarUrl] = useState('');
  const [passwordForm, setPasswordForm] = useState({
    oldPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    if (profile) {
      setProfileForm({
        username: profile.fullName || '',
        name: profile.fullName || '',
        email: profile.email || '',
        description: profile.description || '',
      });
    }
  }, [profile]);

  const handleProfileUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setSaving(true);

    try {
      const updateData: UserProfileUpdate = {};
      if (profileForm.username) updateData.username = profileForm.username;
      if (profileForm.name) updateData.name = profileForm.name;
      if (profileForm.email) updateData.email = profileForm.email;
      if (profileForm.description) updateData.description = profileForm.description;

      const result = await updateProfile(updateData);

      if (result.tokens) {
        saveTokens(result.tokens);
      }

      setSuccess('Профиль успешно обновлен!');
      onProfileUpdated();
    } catch (err) {
      if (err instanceof Error) {
        if (err.message.includes('USERNAME_CONFLICT')) {
          setError('Username уже существует');
        } else if (err.message.includes('EMAIL_CONFLICT')) {
          setError('Email уже существует');
        } else {
          setError(err.message);
        }
      }
    } finally {
      setSaving(false);
    }
  };

  const handleAvatarUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (!/^https?:\/\/.+/.test(avatarUrl)) {
      setError('Введите корректный URL (http:// или https://)');
      return;
    }

    setSaving(true);
    try {
      await updateAvatar(avatarUrl);
      setSuccess('Аватар успешно обновлен!');
      setAvatarUrl('');
      onProfileUpdated();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления аватара');
    } finally {
      setSaving(false);
    }
  };

  const handlePasswordChange = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (passwordForm.newPassword.length < 8) {
      setError('Новый пароль должен содержать минимум 8 символов');
      return;
    }

    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }

    setSaving(true);
    try {
      await changePassword(passwordForm.oldPassword, passwordForm.newPassword);
      setSuccess('Пароль успешно изменен!');
      setPasswordForm({ oldPassword: '', newPassword: '', confirmPassword: '' });
    } catch (err) {
      if (err instanceof Error && err.message.includes('INVALID_OLD_PASSWORD')) {
        setError('Неверный старый пароль');
      } else {
        setError(err instanceof Error ? err.message : 'Ошибка изменения пароля');
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-8 pt-6 max-w-3xl mx-auto">
      {error && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-600 text-sm">{error}</p>
        </div>
      )}
      {success && (
        <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
          <p className="text-green-600 text-sm">{success}</p>
        </div>
      )}

      {/* Profile Section */}
      <Card className="border-gray-200">
        <CardHeader>
          <CardTitle>Основные данные профиля</CardTitle>
          <CardDescription>Обновите информацию о себе</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleProfileUpdate} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                value={profileForm.username}
                onChange={(e) => setProfileForm({ ...profileForm, username: e.target.value })}
                maxLength={50}
                placeholder="username"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Полное имя</Label>
              <Input
                id="name"
                value={profileForm.name}
                onChange={(e) => setProfileForm({ ...profileForm, name: e.target.value })}
                maxLength={50}
                placeholder="Иван Иванов"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={profileForm.email}
                onChange={(e) => setProfileForm({ ...profileForm, email: e.target.value })}
                maxLength={250}
                placeholder="email@example.com"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Описание</Label>
              <Textarea
                id="description"
                value={profileForm.description}
                onChange={(e) => setProfileForm({ ...profileForm, description: e.target.value })}
                maxLength={150}
                rows={3}
                placeholder="Расскажите о себе..."
              />
            </div>
            <Button type="submit" disabled={saving}>
              {saving ? 'Сохранение...' : 'Сохранить изменения'}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Separator />

      {/* Avatar Section */}
      <Card className="border-gray-200">
        <CardHeader>
          <CardTitle>Аватар</CardTitle>
          <CardDescription>Обновите URL вашего аватара</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleAvatarUpdate} className="space-y-4">
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <Avatar className="h-16 w-16">
                  <AvatarImage src={profile?.avatarUrl || avatarUrl || `https://api.dicebear.com/7.x/avataaars/svg?seed=${profile?.email}`} />
                  <AvatarFallback className="bg-blue-100 text-blue-700">
                    <User className="h-8 w-8" />
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1">
                  <p className="text-sm text-gray-500 mb-1">Текущий аватар</p>
                  {profile?.avatarUrl && (
                    <p className="text-xs text-gray-400 truncate max-w-xs">{profile.avatarUrl}</p>
                  )}
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="avatarUrl">Новый URL аватара</Label>
                <Input
                  id="avatarUrl"
                  type="url"
                  value={avatarUrl}
                  onChange={(e) => setAvatarUrl(e.target.value)}
                  placeholder="https://example.com/avatar.png"
                />
                <p className="text-xs text-gray-500">Введите URL изображения (http:// или https://)</p>
              </div>
            </div>
            <Button type="submit" disabled={saving || !avatarUrl}>
              {saving ? 'Обновление...' : 'Обновить аватар'}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Separator />

      {/* Password Section */}
      <Card className="border-gray-200">
        <CardHeader>
          <CardTitle>Изменить пароль</CardTitle>
          <CardDescription>Обновите свой пароль</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handlePasswordChange} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="oldPassword">Текущий пароль</Label>
              <Input
                id="oldPassword"
                type="password"
                value={passwordForm.oldPassword}
                onChange={(e) => setPasswordForm({ ...passwordForm, oldPassword: e.target.value })}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="newPassword">Новый пароль</Label>
              <Input
                id="newPassword"
                type="password"
                value={passwordForm.newPassword}
                onChange={(e) => setPasswordForm({ ...passwordForm, newPassword: e.target.value })}
                required
                minLength={8}
              />
              <p className="text-xs text-gray-500">Минимум 8 символов</p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Подтвердите новый пароль</Label>
              <Input
                id="confirmPassword"
                type="password"
                value={passwordForm.confirmPassword}
                onChange={(e) => setPasswordForm({ ...passwordForm, confirmPassword: e.target.value })}
                required
              />
            </div>
            <Button type="submit" disabled={saving}>
              {saving ? 'Изменение...' : 'Изменить пароль'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}