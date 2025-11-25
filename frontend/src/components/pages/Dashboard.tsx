
import { courses, questions } from '../../lib/data';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Progress } from '../ui/progress';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Activity, BookOpen, Clock, Star, Trophy } from 'lucide-react';

export function Dashboard({ onNavigate }: { onNavigate: (view: any) => void }) {
  // Mock user progress
  const activeCourses = courses.filter(c => c.progress !== undefined && c.progress > 0);
  const difficultyLabel: Record<string, string> = {
    Beginner: 'Начальный',
    Intermediate: 'Средний',
    Advanced: 'Продвинутый',
  };

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Дашборд</h1>
          <p className="text-gray-500 mt-1">С возвращением, Алекс! Серия — 3 дня.</p>
        </div>
        <Button onClick={() => onNavigate('courses')}>Выбрать новый курс</Button>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card className="bg-gray-900 text-white border-none">
           <CardContent className="p-6 flex flex-col justify-between h-full">
              <Activity className="h-8 w-8 opacity-80" />
            <div>
                <div className="text-3xl font-bold">12 ч</div>
                 <div className="text-gray-400 text-sm">Время обучения</div>
            </div>
           </CardContent>
        </Card>
        <Card>
           <CardContent className="p-6 flex flex-col justify-between h-full">
              <BookOpen className="h-8 w-8 text-blue-600" />
             <div>
                <div className="text-3xl font-bold">4</div>
                 <div className="text-gray-500 text-sm">Активных курсов</div>
             </div>
           </CardContent>
        </Card>
        <Card>
           <CardContent className="p-6 flex flex-col justify-between h-full">
              <Trophy className="h-8 w-8 text-amber-500" />
             <div>
                <div className="text-3xl font-bold">850</div>
                 <div className="text-gray-500 text-sm">Очки XP</div>
             </div>
           </CardContent>
        </Card>
        <Card>
           <CardContent className="p-6 flex flex-col justify-between h-full">
              <Star className="h-8 w-8 text-purple-500" />
             <div>
                <div className="text-3xl font-bold">12</div>
                 <div className="text-gray-500 text-sm">Решённых вопросов</div>
             </div>
           </CardContent>
        </Card>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Active Learning */}
        <div className="lg:col-span-2 space-y-6">
          <h2 className="text-xl font-bold text-gray-900">Продолжить обучение</h2>
          <div className="space-y-4">
             {activeCourses.map(course => (
                <Card key={course.id} className="overflow-hidden hover:bg-gray-50 transition-colors cursor-pointer border-gray-200">
                   <div className="flex flex-col sm:flex-row">
                      <div className="w-full sm:w-48 h-32 sm:h-auto relative">
                         <img 
                            src={`https://images.unsplash.com/${course.image}?auto=format&fit=crop&w=400&q=80`}
                            alt={course.title}
                            className="absolute inset-0 w-full h-full object-cover"
                         />
                      </div>
                      <div className="flex-1 p-6 flex flex-col justify-between">
                         <div className="mb-4">
                            <div className="flex justify-between items-start mb-2">
                               <h3 className="font-bold text-gray-900">{course.title}</h3>
                               <Badge variant="secondary">{difficultyLabel[course.difficulty] || course.difficulty}</Badge>
                            </div>
                            <div className="flex items-center gap-2 text-sm text-gray-500">
                               <Clock size={14} />
                               <span>Осталось 2 ч 15 м</span>
                            </div>
                         </div>
                         <div className="space-y-2">
                            <div className="flex justify-between text-xs font-medium text-gray-500">
                               <span>Модуль 4/12</span>
                               <span>{course.progress}%</span>
                            </div>
                            <Progress value={course.progress} className="h-2" />
                         </div>
                      </div>
                   </div>
                </Card>
             ))}
             
             {activeCourses.length === 0 && (
               <div className="text-center py-12 bg-white border border-dashed border-gray-300 rounded-xl">
                 <p className="text-gray-500 mb-4">Вы ещё не начали курсы.</p>
                 <Button variant="outline" onClick={() => onNavigate('courses')}>Открыть каталог</Button>
               </div>
             )}
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
           <Card className="border-none bg-blue-50">
              <CardHeader>
                 <CardTitle className="text-blue-900">Ежедневный челлендж</CardTitle>
                 <CardDescription className="text-blue-700">Решите случайную задачу по алгоритмам, чтобы держать форму.</CardDescription>
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

           <div className="space-y-4">
              <h3 className="font-bold text-gray-900">Рекомендовано вам</h3>
              {questions.slice(0, 3).map(q => (
                 <Card key={q.id} className="cursor-pointer hover:shadow-md transition-shadow">
                    <CardContent className="p-4">
                       <h4 className="font-medium text-sm text-gray-900 mb-2 line-clamp-1">{q.title}</h4>
                       <div className="flex items-center justify-between text-xs text-gray-500">
                          <span>{q.difficulty === 'Beginner' ? 'Начальный' : q.difficulty === 'Intermediate' ? 'Средний' : 'Продвинутый'}</span>
                          <span>{q.replies} ответов</span>
                       </div>
                    </CardContent>
                 </Card>
              ))}
           </div>
        </div>
      </div>
    </div>
  );
}
