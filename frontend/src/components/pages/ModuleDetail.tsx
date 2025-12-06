import React, { useState, useEffect } from 'react';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Progress } from '../ui/progress';
import {
  ArrowLeft, CheckCircle, PlayCircle, Lock, AlertCircle
} from 'lucide-react';
import { motion } from 'motion/react';
import { getModuleByIdFull, completeModule } from '../../lib/api';
import type { ModuleDetailFull, ModuleQuestionItem } from '../../lib/api';

interface ModuleDetailProps {
  courseId: string;
  moduleId: string;
  onBack: () => void;
  onNavigateToQuestion: (questionId: string, courseId: string, moduleId: string) => void;
}

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

export function ModuleDetail({ courseId, moduleId, onBack, onNavigateToQuestion }: ModuleDetailProps) {
  const [module, setModule] = useState<ModuleDetailFull | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [completing, setCompleting] = useState(false);

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      setLoading(true);
      setError(null);

      try {
        const data = await getModuleByIdFull(courseId, moduleId);
        if (cancelled) return;
        setModule(data);
      } catch (err) {
        if (cancelled) return;
        const message = err instanceof Error ? err.message : 'Не удалось загрузить модуль';
        setError(message);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    load();
    return () => {
      cancelled = true;
    };
  }, [courseId, moduleId]);

  const handleCompleteModule = async () => {
    if (!module) return;

    const allAnswered = module.questions?.every(q => q.isAnswered) || false;
    if (!allAnswered) {
      alert('Пожалуйста, ответьте на все вопросы перед завершением модуля');
      return;
    }

    setCompleting(true);
    try {
      await completeModule(courseId, moduleId);
      alert('Модуль успешно завершен!');
      onBack();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не удалось завершить модуль';
      alert(message);
    } finally {
      setCompleting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-500">Загрузка модуля...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[400px]">
        <AlertCircle className="w-16 h-16 text-red-500 mb-4" />
        <p className="text-xl font-semibold text-gray-800 mb-2">Ошибка загрузки</p>
        <p className="text-gray-500 mb-4">{error}</p>
        <Button onClick={onBack} variant="outline">
          <ArrowLeft className="w-4 h-4 mr-2" />
          Назад к курсу
        </Button>
      </div>
    );
  }

  if (!module) {
    return null;
  }

  const allAnswered = module.questions?.every(q => q.isAnswered) || false;
  const progressPercentage = module.progress
    ? (module.progress.answeredQuestions / module.progress.totalQuestions) * 100
    : 0;

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button onClick={onBack} variant="ghost" size="sm">
          <ArrowLeft className="w-4 h-4 mr-2" />
          Назад к курсу
        </Button>
      </div>

      {/* Module Info */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="space-y-4"
      >
        <div className="flex items-center gap-3">
          <div className="bg-blue-100 p-3 rounded-lg">
            <PlayCircle className="w-6 h-6 text-blue-600" />
          </div>
          <div className="flex-1">
            <h1 className="text-3xl font-bold text-gray-900">{module.title}</h1>
            <p className="text-gray-600 mt-1">{module.description}</p>
          </div>
          {module.isCompleted && (
            <Badge className="bg-green-100 text-green-700 border-green-200">
              <CheckCircle className="w-4 h-4 mr-1" />
              Завершено
            </Badge>
          )}
          {module.isLocked && (
            <Badge variant="outline" className="bg-gray-100 text-gray-600">
              <Lock className="w-4 h-4 mr-1" />
              Заблокировано
            </Badge>
          )}
        </div>
      </motion.div>

      {/* Progress Card */}
      {module.progress && (
        <Card>
          <CardHeader>
            <CardTitle>Ваш прогресс</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex justify-between text-sm mb-2">
              <span className="text-gray-600">
                {module.progress.answeredQuestions} / {module.progress.totalQuestions} отвечено
              </span>
              <span className="text-gray-600">
                {module.progress.correctnessPct?.toFixed(0) || 0}% правильных
              </span>
            </div>
            <Progress value={progressPercentage} className="h-2" />
            <div className="grid grid-cols-3 gap-4 pt-2">
              <div className="text-center">
                <p className="text-2xl font-bold text-blue-600">{module.progress.totalQuestions}</p>
                <p className="text-xs text-gray-500">Всего вопросов</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-green-600">{module.progress.correctAnswers}</p>
                <p className="text-xs text-gray-500">Правильных</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-yellow-600">
                  {module.progress.answeredQuestions - module.progress.correctAnswers}
                </p>
                <p className="text-xs text-gray-500">Неправильных</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Questions List */}
      <div className="space-y-4">
        <h2 className="text-2xl font-bold text-gray-900">Вопросы</h2>
        {!module.questions || module.questions.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-gray-500">В этом модуле пока нет вопросов</p>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-3">
            {module.questions.map((q, idx) => (
              <motion.div
                key={q.questionId}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: idx * 0.05 }}
              >
                <Card
                  className="cursor-pointer hover:shadow-md transition-shadow hover:border-blue-300"
                  onClick={() => onNavigateToQuestion(String(q.questionId), courseId, moduleId)}
                >
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-sm font-semibold text-gray-500">#{idx + 1}</span>
                          <CardTitle className="text-lg">{q.title}</CardTitle>
                        </div>
                      </div>
                      <div className="flex gap-2 items-center">
                        <Badge className={difficultyConfig[q.difficulty].color}>
                          {difficultyConfig[q.difficulty].label}
                        </Badge>
                        {q.isAnswered && (
                          <Badge
                            variant={q.isCorrect ? 'default' : 'destructive'}
                            className={q.isCorrect ? 'bg-green-500' : ''}
                          >
                            {q.isCorrect ? (
                              <>
                                <CheckCircle className="w-3 h-3 mr-1" />
                                Правильно
                              </>
                            ) : (
                              <>
                                <AlertCircle className="w-3 h-3 mr-1" />
                                Неправильно
                              </>
                            )}
                          </Badge>
                        )}
                      </div>
                    </div>
                  </CardHeader>
                </Card>
              </motion.div>
            ))}
          </div>
        )}
      </div>

      {/* Complete Module Button */}
      {!module.isCompleted && (
        <Card className="border-2 border-blue-200 bg-blue-50">
          <CardContent className="py-6">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="font-semibold text-lg mb-1">Завершить модуль</h3>
                <p className="text-sm text-gray-600">
                  {allAnswered
                    ? 'Вы ответили на все вопросы. Можете завершить модуль.'
                    : 'Ответьте на все вопросы, чтобы завершить модуль'}
                </p>
              </div>
              <Button
                onClick={handleCompleteModule}
                disabled={!allAnswered || completing}
                size="lg"
                className="bg-blue-600 hover:bg-blue-700"
              >
                {completing ? 'Завершение...' : 'Завершить модуль'}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
