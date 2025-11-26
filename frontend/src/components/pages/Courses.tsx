
import React, { useEffect, useMemo, useState } from 'react';
import { courses as fallbackCourses, Difficulty } from '../../lib/data';
import { Badge } from '../ui/badge';
import { Button } from '../ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../ui/card';
import { Input } from '../ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { BookOpen, Search } from 'lucide-react';
import { motion } from 'motion/react';
import { listCourses } from '../../lib/api';

export function Courses({ onSelectCourse }: { onSelectCourse: (id: string) => void }) {
  const [searchTerm, setSearchTerm] = useState('');
  const [difficultyFilter, setDifficultyFilter] = useState<Difficulty | 'All'>('All');
  const [items, setItems] = useState(fallbackCourses);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const difficultyLabel: Record<Difficulty, string> = {
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
        const res = await listCourses();
        if (cancelled) return;
        const mapped = res.items.map((c, idx) => {
          const fallback = fallbackCourses[idx % fallbackCourses.length];
          return {
            id: String(c.id),
            title: c.title,
            description: c.description,
            modulesCount: c.modules?.length ?? c.totalModules ?? fallback.modulesCount ?? 0,
            difficulty: (['Beginner', 'Intermediate', 'Advanced'] as Difficulty[])[idx % 3],
            image: fallback.image,
            tags: c.tags || fallback.tags || [],
            progress: 0,
          };
        });
        setItems(mapped as any);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Не удалось загрузить курсы');
        setItems(fallbackCourses);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    return () => {
      cancelled = true;
    };
  }, []);

  const filteredCourses = useMemo(() => {
    return items.filter(course => {
      const matchesSearch = course.title.toLowerCase().includes(searchTerm.toLowerCase()) || 
                            course.description.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesDiff = difficultyFilter === 'All' || course.difficulty === difficultyFilter;
      return matchesSearch && matchesDiff;
    });
  }, [items, searchTerm, difficultyFilter]);

  return (
    <div className="space-y-8">
      <div className="flex flex-col md:flex-row justify-between items-end md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-gray-900">Каталог курсов</h1>
          <p className="text-gray-500 mt-2">Программа, собранная под реальные задачи современного инженера.</p>
        </div>
        
        <div className="flex gap-3 w-full md:w-auto">
          <div className="relative flex-1 md:w-64">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 h-4 w-4" />
            <Input 
              placeholder="Поиск курсов..." 
              className="pl-9 bg-white"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
          <Select value={difficultyFilter} onValueChange={(v) => setDifficultyFilter(v as Difficulty | 'All')}>
            <SelectTrigger className="w-[140px] bg-white">
              <SelectValue placeholder="Сложность" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="All">Все уровни</SelectItem>
              <SelectItem value="Beginner">Начальный</SelectItem>
              <SelectItem value="Intermediate">Средний</SelectItem>
              <SelectItem value="Advanced">Продвинутый</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {loading && <p className="text-sm text-gray-500">Загружаем курсы...</p>}
      {error && <p className="text-sm text-red-600">{error}</p>}

      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredCourses.map((course) => (
          <motion.div
            key={course.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
          >
            <Card className="h-full flex flex-col overflow-hidden hover:shadow-lg transition-all duration-300 cursor-pointer group border-gray-200" onClick={() => onSelectCourse(course.id)}>
              <div className="relative h-48 w-full overflow-hidden bg-gray-100">
                 {/* Using unsplash images via ImageWithFallback - assuming URLs from data.ts are partial or keywords if unsplash_tool was used, 
                     but here I used hardcoded unsplash IDs in data.ts. I need to construct full URL or just use unsplash search.
                     Actually data.ts has partials like 'photo-16333...'. I'll construct the full URL.
                 */}
                 <img 
                    src={`https://images.unsplash.com/${course.image}?auto=format&fit=crop&w=800&q=80`}
                    alt={course.title}
                    className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
                 />
                 <div className="absolute top-4 right-4">
                   <Badge variant="secondary" className="bg-white/90 backdrop-blur text-gray-900 font-medium">
                     {difficultyLabel[course.difficulty]}
                   </Badge>
                 </div>
              </div>
              <CardHeader>
                <CardTitle className="text-xl group-hover:text-blue-600 transition-colors">{course.title}</CardTitle>
                <CardDescription className="line-clamp-2">{course.description}</CardDescription>
              </CardHeader>
              <CardContent className="flex-1">
                <div className="flex flex-wrap gap-2">
                  {course.tags.map(tag => (
                    <Badge key={tag} variant="outline" className="text-gray-500 font-normal border-gray-200">
                      {tag}
                    </Badge>
                  ))}
                </div>
              </CardContent>
              <CardFooter className="border-t border-gray-100 pt-4 text-sm text-gray-500 flex justify-between">
                <div className="flex items-center gap-1">
                  <BookOpen size={16} />
                  <span>{course.modulesCount} модулей</span>
                </div>
                {course.progress !== undefined && course.progress > 0 && (
                   <div className="flex items-center gap-2 text-blue-600 font-medium">
                      <span className="text-xs">В процессе</span>
                      <div className="w-16 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                        <div className="h-full bg-blue-600 rounded-full" style={{ width: `${course.progress}%` }} />
                      </div>
                   </div>
                )}
              </CardFooter>
            </Card>
          </motion.div>
        ))}
      </div>
      
      {filteredCourses.length === 0 && (
        <div className="text-center py-24 text-gray-500">
          <p className="text-lg">По вашему запросу курсы не найдены.</p>
          <Button variant="link" onClick={() => { setSearchTerm(''); setDifficultyFilter('All'); }}>
            Сбросить фильтры
          </Button>
        </div>
      )}
    </div>
  );
}
