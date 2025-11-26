import React, { useEffect, useState } from 'react';
import { courses as fallbackCourses, questions as fallbackQuestions } from '../../lib/data';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Progress } from '../ui/progress';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Activity, BookOpen, Clock, Star, Trophy } from 'lucide-react';
import { listCourses, listMentors, listQuestions } from '../../lib/api';

type DashboardProps = { onNavigate: (view: any) => void; userName?: string | null };

export function Dashboard({ onNavigate, userName }: DashboardProps) {
  const [activeCourses, setActiveCourses] = useState(fallbackCourses.filter(c => c.progress !== undefined && c.progress > 0));
  const [questionsData, setQuestionsData] = useState(fallbackQuestions.slice(0, 3));
  const [mentorsCount, setMentorsCount] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const difficultyLabel: Record<string, string> = {
    Beginner: 'Начальный',
    Intermediate: 'Средний',
    Advanced: 'Продвинутый',
  };

  useEffect(() => {
    let cancelled = false;
    const load = async () => {
      setLoading(true);
      try {
        const [coursesRes, questionsRes, mentorsRes] = await Promise.all([
          listCourses(),
          listQuestions(),
          listMentors(),
        ]);
        if (cancelled) return;
        const apiCourses = Array.isArray(coursesRes?.items) ? coursesRes.items : [];
        const mappedCourses = apiCourses
          .map((c, idx) => {
            const fb = fallbackCourses[idx % fallbackCourses.length];
            return {
              id: String(c.id ?? fb.id ?? idx),
              title: c.title ?? fb.title,
              description: c.description ?? fb.description,
              modulesCount: c.modules?.length ?? 0,
              difficulty: (['Beginner', 'Intermediate', 'Advanced'] as const)[idx % 3],
              image: fb.image,
              tags: c.tags || fb.tags || [],
              progress: Math.min(100, (idx + 1) * 15),
            };
          })
          .filter(c => c.progress && c.progress > 0);
        if (mappedCourses.length) setActiveCourses(mappedCourses as any);

        const apiQuestions = Array.isArray(questionsRes?.items) ? questionsRes.items : [];
        if (apiQuestions.length) {
          setQuestionsData(
            apiQuestions.slice(0, 3).map((q, idx) => ({
              id: String(q.id ?? idx),
              title: q.title,
              preview: q.technology || '',
              tags: q.tags || [],
              difficulty: (['Beginner', 'Intermediate', 'Advanced'] as const)[idx % 3],
              replies: 0,
              views: 0,
            })) as any,
          );
        }

        if (typeof mentorsRes?.total === 'number') setMentorsCount(mentorsRes.total);
      } catch {
        // keep fallbacks
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      {/* Header */}
      <div className="flex flex-col gap-2">
        <p className="text-sm text-gray-500">С возвращением{userName ? `, ${userName}` : ''}!</p>
        <h1 className="text-3xl font-bold text-gray-900">Ваш прогресс</h1>
        <p className="text-gray-500">Следите за учебой, вызовами и рекомендациями на одной странице.</p>
      </div>

      {/* Top KPI cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="bg-gray-900 text-white border-none">
          <CardContent className="p-6 flex flex-col gap-4">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-white/10"><Activity className="h-5 w-5" /></div>
              <span className="text-sm text-gray-300">Время обучения</span>
            </div>
            <div className="text-3xl font-bold">{activeCourses.length ? `${activeCourses.length * 3} ч` : '—'}</div>
            <div className="h-2 bg-white/10 rounded-full overflow-hidden">
              <div className="h-full bg-white/70" style={{ width: `${Math.min(100, activeCourses.length * 20)}%` }} />
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6 flex flex-col gap-2">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-blue-50"><BookOpen className="h-5 w-5 text-blue-600" /></div>
              <span className="text-sm text-gray-500">Активных курсов</span>
            </div>
            <div className="text-3xl font-bold">{activeCourses.length || 0}</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6 flex flex-col gap-2">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-amber-50"><Trophy className="h-5 w-5 text-amber-500" /></div>
              <span className="text-sm text-gray-500">Доступных менторов</span>
            </div>
            <div className="text-3xl font-bold">{mentorsCount ?? '—'}</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-6 flex flex-col gap-2">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-purple-50"><Star className="h-5 w-5 text-purple-500" /></div>
              <span className="text-sm text-gray-500">Вопросов решено</span>
            </div>
            <div className="text-3xl font-bold">{questionsData.length || 0}</div>
          </CardContent>
        </Card>
      </div>

      {/* Main layout */}
      <div className="grid lg:grid-cols-3 gap-8">
        {/* Courses */}
        <div className="lg:col-span-2 space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-xl font-bold text-gray-900">Продолжить обучение</h2>
              <p className="text-sm text-gray-500">Ваши активные курсы</p>
            </div>
            <Button variant="outline" onClick={() => onNavigate('courses')}>Каталог</Button>
          </div>
          <div className="space-y-4">
            {activeCourses.slice(0, 3).map((course) => (
              <Card key={course.id} className="overflow-hidden border-gray-200">
                <CardContent className="p-4 sm:p-6 flex flex-col gap-4">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="flex items-center gap-2 mb-2">
                        <Badge variant="secondary">{difficultyLabel[course.difficulty] || course.difficulty}</Badge>
                        <span className="text-xs text-gray-400">Модулей: {course.modulesCount}</span>
                      </div>
                      <h3 className="text-lg font-semibold text-gray-900">{course.title}</h3>
                      <p className="text-sm text-gray-500 mt-1 line-clamp-2">{course.description}</p>
                    </div>
                    <div className="text-right">
                      <span className="text-xs text-gray-400">Прогресс</span>
                      <div className="text-xl font-bold text-gray-900">{course.progress}%</div>
                    </div>
                  </div>
                  <div>
                    <div className="flex justify-between text-xs text-gray-500 mb-1">
                      <span>Модуль 4/12</span>
                      <span>{course.progress}%</span>
                    </div>
                    <Progress value={course.progress} className="h-2" />
                  </div>
                  <div className="flex gap-2">
                    <Button size="sm" onClick={() => onNavigate('course-detail')}>Продолжить</Button>
                    <Button size="sm" variant="ghost" onClick={() => onNavigate('course-detail')}>Детали</Button>
                  </div>
                </CardContent>
              </Card>
            ))}
            {activeCourses.length === 0 && (
              <div className="text-center py-8 bg-white border border-dashed border-gray-300 rounded-xl">
                <p className="text-gray-500 mb-4">Нет активных курсов.</p>
                <Button variant="outline" onClick={() => onNavigate('courses')}>Открыть каталог</Button>
              </div>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          <Card className="border-none bg-blue-50">
            <CardHeader>
              <CardTitle className="text-blue-900">Ежедневный челлендж</CardTitle>
              <CardDescription className="text-blue-700">Решите задачу, чтобы держать форму.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="bg-white p-4 rounded-lg shadow-sm mb-4">
                <h4 className="font-bold text-gray-900 mb-2">Инвертировать бинарное дерево</h4>
                <div className="flex gap-2 mb-2">
                  <Badge variant="secondary" className="text-xs">Лёгкая</Badge>
                  <Badge variant="secondary" className="text-xs">Деревья</Badge>
                </div>
              </div>
              <Button className="w-full bg-blue-600 hover:bg-blue-700">Начать</Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-lg">Рекомендации</CardTitle>
              <CardDescription>Свежие вопросы для практики</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              {questionsData.slice(0, 3).map((q) => (
                <div key={q.id} className="p-3 border border-gray-100 rounded-lg hover:shadow-sm transition">
                  <div className="flex items-center justify-between mb-2">
                    <Badge variant="outline" className="text-xs">
                      {q.difficulty === 'Beginner' ? 'Начальный' : q.difficulty === 'Intermediate' ? 'Средний' : 'Продвинутый'}
                    </Badge>
                    <span className="text-xs text-gray-400">ответов: {q.replies}</span>
                  </div>
                  <p className="text-sm font-medium text-gray-900 line-clamp-1">{q.title}</p>
                  <p className="text-xs text-gray-500 line-clamp-1">{q.preview}</p>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
