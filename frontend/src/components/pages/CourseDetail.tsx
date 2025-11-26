
import React, { useEffect, useState } from 'react';
import { Course, courses, mentors } from '../../lib/data';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Progress } from '../ui/progress';
import { Card, CardContent } from '../ui/card';
import { CheckCircle, PlayCircle, Lock, Clock, Calendar, User } from 'lucide-react';
import { getCourseById, getCourseModules } from '../../lib/api';

export function CourseDetail({ courseId, onBack }: { courseId: string, onBack: () => void }) {
  const fallback = courses.find(c => c.id === courseId) || courses[0];
  const [course, setCourse] = useState<Course | null>(fallback || null);
  const [modules, setModules] = useState<
    { id: string; title: string; duration?: string; status: 'completed' | 'in-progress' | 'locked' }[]
  >([]);
  const [error, setError] = useState<string | null>(null);
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
      setError(null);
      try {
        const data = await getCourseById(courseId);
        if (cancelled) return;
        setCourse({
          id: String(data.id),
          title: data.title,
          description: data.description,
          modulesCount: data.modules?.length ?? fallback?.modulesCount ?? 0,
          difficulty: fallback?.difficulty || 'Intermediate',
          image: fallback?.image,
          tags: fallback?.tags || [],
          progress: 0,
        });
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Не удалось загрузить курс, показываем демо-версию.');
          setCourse(fallback || null);
        }
      }

      try {
        const mods = await getCourseModules(courseId);
        if (cancelled) return;
        const mapped = mods.items.map((m, idx) => ({
          id: String(m.id ?? idx),
          title: m.title,
          duration: m.description || '—',
          status: idx === 0 ? 'completed' : idx === 1 ? 'in-progress' : 'locked',
        }));
        setModules(mapped);
      } catch {
        if (!cancelled) {
          setModules([
            { id: 'm1', title: 'Введение и настройка', duration: '45 мин', status: 'completed' },
            { id: 'm2', title: 'Глубокие базовые концепты', duration: '1 ч 20 мин', status: 'in-progress' },
            { id: 'm3', title: 'Паттерны архитектуры', duration: '55 мин', status: 'locked' },
          ]);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    return () => {
      cancelled = true;
    };
  }, [courseId, fallback]);

  if (!course) return <div>Курс не найден</div>;

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <Button variant="ghost" onClick={onBack} className="pl-0 hover:bg-transparent hover:text-blue-600">
        ← Назад к курсам
      </Button>
      {loading && <p className="text-sm text-gray-500">Загружаем курс...</p>}
      {error && <p className="text-sm text-red-600">{error}</p>}

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-8">
          <div>
            <div className="flex items-center gap-3 mb-4">
              <Badge variant="secondary" className="text-sm">{difficultyLabel[course.difficulty] || course.difficulty}</Badge>
              <span className="text-gray-400">•</span>
              <span className="text-gray-500 text-sm">{course.modulesCount} модулей</span>
            </div>
            <h1 className="text-4xl font-bold text-gray-900 mb-4">{course.title}</h1>
            <p className="text-xl text-gray-500 leading-relaxed">
              {course.description} 
              Этот курс проводит по ключевым аспектам технологии и готовит к продакшену.
            </p>
          </div>

          <div className="space-y-6">
            <h3 className="text-2xl font-bold text-gray-900">Программа</h3>
            <div className="bg-white border border-gray-200 rounded-2xl overflow-hidden shadow-sm">
              {modules.map((module, index) => (
                <div 
                  key={module.id} 
                  className={`p-4 flex items-center justify-between border-b border-gray-100 last:border-0 hover:bg-gray-50 transition-colors ${
                    module.status === 'locked' ? 'opacity-60' : ''
                  }`}
                >
                 <div className="flex items-center gap-4">
                   <div className={`
                      w-8 h-8 rounded-full flex items-center justify-center
                      ${module.status === 'completed' ? 'bg-green-100 text-green-600' : 
                        module.status === 'in-progress' ? 'bg-blue-100 text-blue-600' : 'bg-gray-100 text-gray-400'}
                    `}>
                      {module.status === 'completed' ? <CheckCircle size={18} /> : 
                       module.status === 'in-progress' ? <PlayCircle size={18} /> : <Lock size={18} />}
                    </div>
                    <div>
                      <h4 className="font-medium text-gray-900">{index + 1}. {module.title}</h4>
                      <p className="text-sm text-gray-500">{module.duration}</p>
                    </div>
                  </div>
                  <div>
                    {module.status === 'in-progress' && <Button size="sm" variant="secondary">Продолжить</Button>}
                    {module.status === 'completed' && <Button size="sm" variant="ghost">Пересмотреть</Button>}
                  </div>
                </div>
              ))}
              {modules.length === 0 && (
                <div className="p-4 text-sm text-gray-500">Нет модулей для отображения.</div>
              )}
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card className="border-none shadow-xl bg-gray-900 text-white sticky top-24">
            <CardContent className="p-6 space-y-6">
              <div className="space-y-2">
                <span className="text-gray-400 text-sm font-medium uppercase tracking-wider">Ваш прогресс</span>
                <div className="flex items-end justify-between">
                  <span className="text-3xl font-bold">{course.progress || 0}%</span>
                  <span className="text-gray-400 text-sm mb-1">2/12 завершено</span>
                </div>
                <Progress value={course.progress || 0} className="h-2 bg-gray-700 [&>div]:bg-blue-500" />
              </div>
              
              <Button className="w-full h-12 text-lg font-medium bg-blue-600 hover:bg-blue-500 text-white border-none">
                {course.progress ? 'Продолжить' : 'Записаться'}
              </Button>
              
              <div className="pt-4 border-t border-gray-800 space-y-3">
                <div className="flex items-center gap-3 text-gray-300">
                  <Clock size={18} />
                  <span className="text-sm">Примерно 12 часов</span>
                </div>
                <div className="flex items-center gap-3 text-gray-300">
                  <Calendar size={18} />
                  <span className="text-sm">Обновлено: октябрь 2024</span>
                </div>
                <div className="flex items-center gap-3 text-gray-300">
                   <User size={18} />
                   <span className="text-sm">1 203 студентов</span>
                </div>
              </div>
            </CardContent>
          </Card>

          <div className="bg-white border border-gray-200 rounded-xl p-6 space-y-4 shadow-sm">
            <h4 className="font-bold text-gray-900">Ментор курса</h4>
            <div className="flex items-center gap-4">
               <img 
                  src={`https://images.unsplash.com/${mentors[0].avatar}?auto=format&fit=crop&w=100&q=80`} 
                  alt={mentors[0].name}
                  className="w-12 h-12 rounded-full object-cover"
               />
               <div>
                 <div className="font-medium text-gray-900">{mentors[0].name}</div>
                 <div className="text-sm text-gray-500">{mentors[0].role}</div>
               </div>
            </div>
            <p className="text-sm text-gray-500">
              "Курс выведет вас из зоны комфорта. Удачи!"
            </p>
            <Button variant="outline" className="w-full">Профиль ментора</Button>
          </div>
        </div>
      </div>
    </div>
  );
}
