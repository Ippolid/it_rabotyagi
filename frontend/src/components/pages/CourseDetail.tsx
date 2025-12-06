
import React, { useEffect, useState } from 'react';
import { Course, courses, mentors } from '../../lib/data';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Progress } from '../ui/progress';
import { Card, CardContent } from '../ui/card';
import { CheckCircle, PlayCircle, Lock, Clock, Calendar, User } from 'lucide-react';
import { getCourseById, getCourseModules, enrollInCourse, getCourseProgress } from '../../lib/api';

export function CourseDetail({
  courseId,
  onBack,
  onNavigateToModule
}: {
  courseId: string;
  onBack: () => void;
  onNavigateToModule?: (courseId: string, moduleId: string) => void;
}) {
  const fallback = courses.find(c => c.id === courseId) || courses[0];
  const [course, setCourse] = useState<Course | null>(fallback || null);
  const [modules, setModules] = useState<
    { id: string; title: string; duration?: string; status: 'completed' | 'in-progress' | 'locked'; order: number }[]
  >([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [enrolled, setEnrolled] = useState(false);
  const [enrolling, setEnrolling] = useState(false);
  const [progress, setProgress] = useState<{ completed: number; total: number; percentage: number }>({
    completed: 0,
    total: 0,
    percentage: 0
  });

  const difficultyLabel: Record<string, string> = {
    Beginner: 'Начальный',
    Intermediate: 'Средний',
    Advanced: 'Продвинутый',
  };
  
  const loadProgress = async () => {
    try {
      const progressData = await getCourseProgress(courseId);
      setEnrolled(true);
      setProgress({
        completed: progressData.completedModules,
        total: progressData.totalModules,
        percentage: progressData.percentage || 0
      });
      return progressData;
    } catch (err) {
      // Not enrolled or error
      setEnrolled(false);
      return null;
    }
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

      // Load modules
      try {
        const mods = await getCourseModules(courseId);
        if (cancelled) return;

        // Load progress to determine module statuses
        const progressData = await loadProgress();

        const completedModulesSet = new Set(
          progressData?.modules?.filter(m => m.completed).map(m => m.moduleId) || []
        );

        const mapped = mods.items.map((m, idx) => {
          const moduleId = m.id ?? idx;
          const order = m.order ?? idx + 1;
          const isCompleted = completedModulesSet.has(moduleId);

          // Sequential lock logic: module is unlocked if it's first OR previous is completed
          const prevModule = idx > 0 ? mods.items[idx - 1] : null;
          const prevCompleted = !prevModule || completedModulesSet.has(prevModule.id);

          // User is enrolled if progressData is not null
          const isEnrolled = progressData !== null;

          let status: 'completed' | 'in-progress' | 'locked';
          if (isCompleted) {
            status = 'completed';
          } else if (prevCompleted && isEnrolled) {
            status = 'in-progress';
          } else {
            status = 'locked';
          }

          return {
            id: String(moduleId),
            title: m.title,
            duration: m.description || '—',
            status,
            order
          };
        });
        setModules(mapped);
      } catch {
        if (!cancelled) {
          setModules([]);
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

  const handleEnroll = async () => {
    setEnrolling(true);
    try {
      await enrollInCourse(courseId);
      setEnrolled(true);
      // Reload to update module statuses
      window.location.reload();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не удалось записаться на курс';
      alert(message);
    } finally {
      setEnrolling(false);
    }
  };

  const handleModuleClick = (module: typeof modules[0]) => {
    if (!enrolled) {
      alert('Сначала запишитесь на курс');
      return;
    }
    if (module.status === 'locked') {
      alert('Завершите предыдущий модуль, чтобы разблокировать этот');
      return;
    }
    if (onNavigateToModule) {
      onNavigateToModule(courseId, module.id);
    }
  };

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
                  onClick={() => handleModuleClick(module)}
                  className={`p-4 flex items-center justify-between border-b border-gray-100 last:border-0 transition-colors ${
                    module.status === 'locked' ? 'opacity-60 cursor-not-allowed' : 'hover:bg-gray-50 cursor-pointer'
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
                    {module.status === 'in-progress' && <Button size="sm" variant="secondary" onClick={(e) => { e.stopPropagation(); handleModuleClick(module); }}>Продолжить</Button>}
                    {module.status === 'completed' && <Button size="sm" variant="ghost" onClick={(e) => { e.stopPropagation(); handleModuleClick(module); }}>Пересмотреть</Button>}
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
                  <span className="text-3xl font-bold">{enrolled ? Math.round(progress.percentage) : 0}%</span>
                  <span className="text-gray-400 text-sm mb-1">
                    {enrolled ? `${progress.completed}/${progress.total} завершено` : '—'}
                  </span>
                </div>
                <Progress value={enrolled ? progress.percentage : 0} className="h-2 bg-gray-700 [&>div]:bg-blue-500" />
              </div>

              <Button
                onClick={enrolled ? undefined : handleEnroll}
                disabled={enrolling}
                className="w-full h-12 text-lg font-medium bg-blue-600 hover:bg-blue-500 text-white border-none"
              >
                {enrolling ? 'Записываемся...' : enrolled ? 'Вы записаны' : 'Записаться'}
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
