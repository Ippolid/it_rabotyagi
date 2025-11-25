
import React from 'react';
import { questions, Difficulty } from '../../lib/data';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Input } from '../ui/input';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/card';
import { Search, MessageSquare, Eye, ThumbsUp, Filter } from 'lucide-react';
import { motion } from 'motion/react';

export function Questions() {
  const [search, setSearch] = React.useState('');

  return (
    <div className="space-y-8">
      <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">База вопросов</h1>
          <p className="text-gray-500 mt-2">Ответы сообщества на сложные инженерные задачи.</p>
        </div>
        <Button className="bg-gray-900 text-white">Задать вопрос</Button>
      </div>

      {/* Search & Filter */}
      <div className="flex gap-4 bg-white p-4 rounded-xl border border-gray-200 shadow-sm">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 h-5 w-5" />
          <Input 
            placeholder="Поиск по вопросам, тегам или ошибкам..." 
            className="pl-10 border-gray-200 bg-gray-50/50" 
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <Button variant="outline" className="gap-2">
          <Filter size={18} />
          Фильтры
        </Button>
      </div>

      {/* Questions List */}
      <div className="space-y-4">
        {questions.map((q) => (
          <motion.div 
             key={q.id}
             initial={{ opacity: 0, y: 10 }}
             animate={{ opacity: 1, y: 0 }}
          >
            <Card className="hover:shadow-md transition-shadow cursor-pointer border-gray-200">
              <CardHeader className="pb-2">
                <div className="flex justify-between items-start">
                   <div className="space-y-1">
                      <div className="flex gap-2 mb-2">
                        {q.tags.map(tag => (
                          <Badge key={tag} variant="secondary" className="bg-gray-100 text-gray-600 hover:bg-gray-200 border-none font-normal">
                            {tag}
                          </Badge>
                        ))}
                        <Badge variant="outline" className={`
                          ${q.difficulty === 'Beginner' ? 'text-green-600 border-green-200 bg-green-50' :
                            q.difficulty === 'Intermediate' ? 'text-blue-600 border-blue-200 bg-blue-50' :
                            'text-orange-600 border-orange-200 bg-orange-50'}
                        `}>{q.difficulty === 'Beginner' ? 'Начальный' : q.difficulty === 'Intermediate' ? 'Средний' : 'Продвинутый'}</Badge>
                      </div>
                      <CardTitle className="text-lg text-blue-600 hover:underline">{q.title}</CardTitle>
                   </div>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600 line-clamp-2 mb-4">{q.preview}</p>
                <div className="flex items-center gap-6 text-sm text-gray-400">
                   <div className="flex items-center gap-1.5 hover:text-gray-600">
                      <MessageSquare size={16} />
                      <span>{q.replies} ответов</span>
                   </div>
                   <div className="flex items-center gap-1.5">
                      <Eye size={16} />
                      <span>{q.views} просмотров</span>
                   </div>
                   <div className="flex items-center gap-1.5 ml-auto">
                      <span className="text-gray-400 text-xs">Активность 2 ч назад</span>
                   </div>
                </div>
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
