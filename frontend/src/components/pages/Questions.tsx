
import React, { useEffect, useState } from 'react';
import { QuestionListItem, QuestionDifficulty, QuestionFilters } from '../../lib/data';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Input } from '../ui/input';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/card';
import { Search, Building2, Laptop, CheckCircle2 } from 'lucide-react';
import { motion } from 'motion/react';
import { listQuestions } from '../../lib/api';

interface QuestionsProps {
  onSelectQuestion: (id: string) => void;
}

const LIMIT = 20;

const TECHNOLOGIES = ['All', 'JavaScript', 'Python', 'Go', 'React', 'TypeScript', 'Node.js', 'PostgreSQL'];
const DIFFICULTIES: Array<'All' | QuestionDifficulty> = ['All', 'easy', 'medium', 'hard'];
const COMPANIES = ['All', 'Google', 'Meta', 'Yandex', 'VK', 'Amazon'];

const difficultyConfig = {
  easy: {
    label: 'Легко',
    color: 'bg-green-100 text-green-700 border-green-200'
  },
  medium: {
    label: 'Средне',
    color: 'bg-yellow-100 text-yellow-700 border-yellow-200'
  },
  hard: {
    label: 'Сложно',
    color: 'bg-red-100 text-red-700 border-red-200'
  }
};

export function Questions({ onSelectQuestion }: QuestionsProps) {
  const [items, setItems] = useState<QuestionListItem[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Filter states
  const [search, setSearch] = useState('');
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [selectedTechnology, setSelectedTechnology] = useState<string>('All');
  const [selectedDifficulty, setSelectedDifficulty] = useState<'All' | QuestionDifficulty>('All');
  const [selectedCompany, setSelectedCompany] = useState<string>('All');
  const [offset, setOffset] = useState(0);

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  // Load questions function
  const loadQuestions = async (append = false) => {
    if (append) {
      setLoadingMore(true);
    } else {
      setLoading(true);
    }
    setError(null);

    try {
      const filters: QuestionFilters = {
        limit: LIMIT,
        offset: append ? offset : 0,
      };

      if (debouncedSearch) filters.search = debouncedSearch;
      if (selectedTechnology !== 'All') filters.technology = selectedTechnology;
      if (selectedDifficulty !== 'All') filters.difficulty = selectedDifficulty;
      if (selectedCompany !== 'All') filters.company = selectedCompany;

      const res = await listQuestions(filters);

      if (append) {
        setItems(prev => [...prev, ...res.items]);
        setOffset(prev => prev + LIMIT);
      } else {
        setItems(res.items);
        setOffset(LIMIT);
      }
      setTotal(res.total);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не удалось загрузить вопросы';
      setError(message);
      if (!append) {
        setItems([]);
        setTotal(0);
      }
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  };

  // Initial load
  useEffect(() => {
    loadQuestions(false);
  }, []);

  // Reset and reload when filters change
  useEffect(() => {
    if (offset === 0 && items.length === 0) return; // Skip initial render
    setOffset(0);
    loadQuestions(false);
  }, [debouncedSearch, selectedTechnology, selectedDifficulty, selectedCompany]);

  const handleLoadMore = () => {
    loadQuestions(true);
  };

  const handleClearFilters = () => {
    setSearch('');
    setDebouncedSearch('');
    setSelectedTechnology('All');
    setSelectedDifficulty('All');
    setSelectedCompany('All');
  };

  const hasActiveFilters = debouncedSearch || selectedTechnology !== 'All' || selectedDifficulty !== 'All' || selectedCompany !== 'All';

  return (
    <div className="space-y-8">
      <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">База вопросов</h1>
          <p className="text-gray-500 mt-2">Подготовка к собеседованиям в IT-компаниях.</p>
        </div>
      </div>

      {/* Filters Section */}
      <div className="space-y-4">
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 h-5 w-5" />
          <Input
            placeholder="Поиск по названию и содержанию вопроса..."
            className="pl-10 border-gray-200 bg-white"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        {/* Filter Dropdowns */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          {/* Technology Filter */}
          <div>
            <label className="text-sm font-medium text-gray-700 mb-2 block">Технология</label>
            <select
              value={selectedTechnology}
              onChange={(e) => setSelectedTechnology(e.target.value)}
              className="w-full h-10 px-3 rounded-md border border-gray-200 bg-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {TECHNOLOGIES.map(tech => (
                <option key={tech} value={tech}>{tech}</option>
              ))}
            </select>
          </div>

          {/* Difficulty Filter */}
          <div>
            <label className="text-sm font-medium text-gray-700 mb-2 block">Сложность</label>
            <select
              value={selectedDifficulty}
              onChange={(e) => setSelectedDifficulty(e.target.value as any)}
              className="w-full h-10 px-3 rounded-md border border-gray-200 bg-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {DIFFICULTIES.map(diff => (
                <option key={diff} value={diff}>
                  {diff === 'All' ? 'Все' : difficultyConfig[diff].label}
                </option>
              ))}
            </select>
          </div>

          {/* Company Filter */}
          <div>
            <label className="text-sm font-medium text-gray-700 mb-2 block">Компания</label>
            <select
              value={selectedCompany}
              onChange={(e) => setSelectedCompany(e.target.value)}
              className="w-full h-10 px-3 rounded-md border border-gray-200 bg-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {COMPANIES.map(company => (
                <option key={company} value={company}>{company}</option>
              ))}
            </select>
          </div>

          {/* Clear Filters Button */}
          <div className="flex items-end">
            <Button
              variant="outline"
              onClick={handleClearFilters}
              disabled={!hasActiveFilters}
              className="w-full"
            >
              Сбросить фильтры
            </Button>
          </div>
        </div>
      </div>

      {/* Results Count */}
      {!loading && (
        <div className="text-sm text-gray-500">
          Показано {items.length} из {total} вопросов
        </div>
      )}

      {/* Loading State */}
      {loading && (
        <div className="text-center py-12">
          <p className="text-gray-500">Загрузка вопросов...</p>
        </div>
      )}

      {/* Error State */}
      {error && !loading && (
        <div className="text-center py-12">
          <p className="text-red-600 mb-4">{error}</p>
          <Button onClick={() => loadQuestions(false)} variant="outline">
            Попробовать снова
          </Button>
        </div>
      )}

      {/* Empty State */}
      {!loading && !error && items.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500 mb-4">Вопросы не найдены</p>
          {hasActiveFilters && (
            <Button onClick={handleClearFilters} variant="outline">
              Сбросить фильтры
            </Button>
          )}
        </div>
      )}

      {/* Questions List */}
      {!loading && items.length > 0 && (
        <div className="space-y-4">
          {items.map((q, idx) => (
            <motion.div
              key={q.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: idx * 0.05 }}
            >
              <Card
                className={`hover:shadow-md transition-shadow cursor-pointer border-gray-200 ${
                  q.solved ? 'bg-green-50/50 border-green-200' : ''
                }`}
                onClick={() => onSelectQuestion(String(q.id))}
              >
                <CardHeader className="pb-3">
                  <div className="flex flex-wrap items-center gap-2 mb-3">
                    {q.solved && (
                      <Badge className="bg-green-100 text-green-700 border-green-200">
                        <CheckCircle2 size={12} className="mr-1" />
                        Решено
                      </Badge>
                    )}

                    <Badge className={difficultyConfig[q.difficulty].color}>
                      {difficultyConfig[q.difficulty].label}
                    </Badge>

                    <div className="flex items-center gap-1.5 text-xs text-gray-500 bg-gray-50 px-2 py-1 rounded-md">
                      <Laptop size={12} />
                      {q.technology}
                    </div>

                    {q.companyTags.slice(0, 3).map((company) => (
                      <div key={company} className="flex items-center gap-1.5 text-xs text-gray-500 bg-gray-50 px-2 py-1 rounded-md">
                        <Building2 size={12} />
                        {company}
                      </div>
                    ))}
                  </div>

                  <CardTitle className={`text-lg hover:text-blue-600 transition-colors break-words ${
                    q.solved ? 'text-gray-700' : 'text-gray-900'
                  }`}>
                    {q.title}
                  </CardTitle>
                </CardHeader>
              </Card>
            </motion.div>
          ))}
        </div>
      )}

      {/* Load More Button */}
      {!loading && !error && items.length < total && (
        <div className="flex justify-center pt-4">
          <Button
            onClick={handleLoadMore}
            disabled={loadingMore}
            variant="outline"
            className="min-w-[200px]"
          >
            {loadingMore ? 'Загрузка...' : 'Загрузить еще'}
          </Button>
        </div>
      )}
    </div>
  );
}
