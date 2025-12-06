
import React, { useState, useEffect } from 'react';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Card, CardContent } from '../ui/card';
import {
  ArrowLeft, CheckCircle, XCircle,
  Building2, Laptop
} from 'lucide-react';
import { motion } from 'motion/react';
import { getQuestionById, listQuestions, submitQuestionAnswer } from '../../lib/api';
import { QuestionDetail as QuestionDetailType, QuestionListItem } from '../../lib/data';

interface QuestionDetailProps {
  questionId: string;
  onBack: () => void;
  onSelectQuestion?: (id: string) => void;
  courseId?: string;
  moduleId?: string;
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

export function QuestionDetail({ questionId, onBack, onSelectQuestion, courseId, moduleId }: QuestionDetailProps) {
  const [question, setQuestion] = useState<QuestionDetailType | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedAnswer, setSelectedAnswer] = useState<string | null>(null);
  const [userTextAnswer, setUserTextAnswer] = useState<string>('');
  const [showAnswer, setShowAnswer] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitResult, setSubmitResult] = useState<{
    isCorrect: boolean;
    correctAnswer: string;
    explanation?: string;
    alreadySolved?: boolean;
  } | null>(null);
  const [relatedQuestions, setRelatedQuestions] = useState<QuestionListItem[]>([]);

  useEffect(() => {
    let cancelled = false;

    const load = async () => {
      setLoading(true);
      setError(null);
      setSelectedAnswer(null);
      setUserTextAnswer('');
      setShowAnswer(false);
      setIsSubmitted(false);
      setIsSubmitting(false);
      setSubmitResult(null);

      try {
        const data = await getQuestionById(questionId);
        if (cancelled) return;
        setQuestion(data);
      } catch (err) {
        if (cancelled) return;
        const message = err instanceof Error ? err.message : 'Не удалось загрузить вопрос';
        setError(message);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    load();
    return () => {
      cancelled = true;
    };
  }, [questionId]);

  // Load related questions
  useEffect(() => {
    if (!question) return;

    let cancelled = false;

    const loadRelated = async () => {
      try {
        const res = await listQuestions({
          technology: question.technology,
          limit: 10
        });
        if (cancelled) return;

        const filtered = res.items
          .filter(q => q.id !== question.id)
          .slice(0, 5);
        setRelatedQuestions(filtered);
      } catch {
        // Ignore errors for related questions
      }
    };

    loadRelated();
    return () => {
      cancelled = true;
    };
  }, [question]);

  const handleSubmit = async () => {
    // Для открытых вопросов не требуется выбранный ответ
    if (!question?.isOpenEnded && !selectedAnswer) return;
    if (!questionId) return;

    setIsSubmitting(true);
    try {
      // Для открытых вопросов отправляем пустую строку
      const answer = question?.isOpenEnded ? '' : selectedAnswer || '';
      const result = await submitQuestionAnswer(
        questionId,
        answer,
        courseId ? Number(courseId) : undefined,
        moduleId ? Number(moduleId) : undefined
      );
      setSubmitResult(result);
      setIsSubmitted(true);

      if (result.alreadySolved) {
        console.log('This question was already solved before');
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не удалось отправить ответ';
      alert(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-gray-500">Загрузка вопроса...</p>
      </div>
    );
  }

  if (error) {
    const is404 = error.toLowerCase().includes('404') || error.toLowerCase().includes('not found');

    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-4">
          <h2 className="text-2xl font-bold text-gray-900">
            {is404 ? 'Вопрос не найден' : 'Ошибка загрузки'}
          </h2>
          <p className="text-gray-500 max-w-md mx-auto">
            {is404
              ? `Вопрос с ID ${questionId} не существует или был удален.`
              : error}
          </p>
          <Button onClick={onBack} variant="outline">
            Вернуться к списку вопросов
          </Button>
        </div>
      </div>
    );
  }

  if (!question) return null;

  const isCorrect = isSubmitted && submitResult?.isCorrect;
  const isWrong = isSubmitted && !submitResult?.isCorrect;

  return (
    <div className="min-h-screen bg-[#F5F5F7] pb-24">
      {/* Sticky Top Nav */}
      <div className="sticky top-0 z-40 bg-white/80 backdrop-blur-md border-b border-gray-200/50 px-4 md:px-8 py-3">
        <div className="flex items-center gap-3 max-w-7xl mx-auto">
          <Button
            variant="ghost"
            size="sm"
            onClick={onBack}
            className="text-gray-500 hover:text-gray-900 -ml-2"
          >
            <ArrowLeft className="w-5 h-5 mr-1" />
            {courseId && moduleId ? 'К модулю' : 'Назад'}
          </Button>
          <div className="h-4 w-[1px] bg-gray-300 mx-1" />
          {courseId && moduleId && (
            <>
              <span className="text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded-md font-medium">
                Модуль курса
              </span>
              <div className="h-4 w-[1px] bg-gray-300 mx-1" />
            </>
          )}
          <span className="text-sm font-medium text-gray-500 truncate flex-1">
            {question.title}
          </span>
        </div>
      </div>

      <main className="max-w-7xl mx-auto px-4 md:px-6 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* Main Content Column */}
          <div className="lg:col-span-8 space-y-6">
            {/* Header Section */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-white rounded-xl p-6 md:p-8 shadow-sm border border-gray-100"
            >
              <div className="flex flex-wrap items-center gap-3 mb-4">
                <Badge className={difficultyConfig[question.difficulty].color}>
                  {difficultyConfig[question.difficulty].label}
                </Badge>

                {question.tags?.map((tag) => (
                  <div key={tag} className="flex items-center gap-1.5 text-xs text-gray-500 bg-gray-50 px-2 py-1 rounded-md">
                    <Building2 size={12} />
                    {tag}
                  </div>
                ))}
              </div>

              <h1 className="text-3xl font-bold text-gray-900 tracking-tight leading-tight">
                {question.title}
              </h1>
            </motion.div>

            {/* Question Content */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.1 }}
              className="bg-white rounded-xl p-6 md:p-8 shadow-sm border border-gray-100"
            >
              <div className="prose prose-gray max-w-none">
                <p className="text-[17px] leading-[1.6] text-gray-700 whitespace-pre-wrap">
                  {question.content}
                </p>
              </div>

              {/* Technology Chip */}
              <div className="mt-6 pt-6 border-t border-gray-50">
                <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-50 rounded-full border border-blue-200/50 text-sm font-medium text-blue-700 w-fit">
                  <Laptop size={14} className="text-blue-600" />
                  {question.technology}
                </div>
              </div>
            </motion.div>

            {/* Answer Section */}
            {question.isOpenEnded ? (
              /* Open-ended question */
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.2 }}
                className="bg-white rounded-xl p-6 md:p-8 shadow-sm border border-gray-100"
              >
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Ваш ответ:</h3>

                <textarea
                  className="w-full min-h-[150px] p-4 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-y"
                  placeholder="Напишите ваш ответ здесь..."
                  value={userTextAnswer}
                  onChange={(e) => setUserTextAnswer(e.target.value)}
                  disabled={showAnswer}
                />

                {!showAnswer && (
                  <div className="mt-4 flex justify-end">
                    <Button
                      onClick={() => setShowAnswer(true)}
                      className="bg-blue-600 hover:bg-blue-700 text-white px-8"
                    >
                      Показать примерный ответ
                    </Button>
                  </div>
                )}

                {showAnswer && question.explanation && (
                  <motion.div
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="mt-6 p-6 bg-blue-50 border border-blue-200 rounded-lg"
                  >
                    <h4 className="font-semibold text-blue-900 mb-3">Примерный ответ:</h4>
                    <p className="text-blue-800 leading-relaxed whitespace-pre-wrap">
                      {question.explanation}
                    </p>
                  </motion.div>
                )}
              </motion.div>
            ) : (
              /* Multiple choice question */
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.2 }}
                className="bg-white rounded-xl p-6 md:p-8 shadow-sm border border-gray-100"
              >
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Варианты ответа:</h3>

                <div className="space-y-3">
                  {question.options?.map((option, idx) => {
                  const optionLetter = String.fromCharCode(65 + idx);
                  const isSelected = selectedAnswer === option;
                  const showCorrect = isSubmitted && submitResult && option === submitResult.correctAnswer;
                  const showWrong = isSubmitted && isSelected && submitResult && option !== submitResult.correctAnswer;

                  return (
                    <label
                      key={idx}
                      className={`
                        group relative flex items-start gap-4 p-4 rounded-xl border-2 cursor-pointer transition-all duration-200
                        ${!isSubmitted ? 'hover:border-blue-300 hover:bg-blue-50/30' : ''}
                        ${isSelected && !isSubmitted ? 'border-blue-500 bg-blue-50' : ''}
                        ${!isSelected && !isSubmitted ? 'border-gray-200 bg-white' : ''}
                        ${showCorrect ? 'border-green-500 bg-green-50' : ''}
                        ${showWrong ? 'border-red-500 bg-red-50' : ''}
                        ${isSubmitted && !showCorrect && !showWrong ? 'border-gray-200 bg-gray-50 opacity-60' : ''}
                      `}
                    >
                      <input
                        type="radio"
                        name="answer"
                        value={option}
                        checked={isSelected}
                        onChange={(e) => !isSubmitted && setSelectedAnswer(e.target.value)}
                        disabled={isSubmitted}
                        className="sr-only"
                      />

                      <div className={`
                        flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold transition-colors
                        ${!isSubmitted && !isSelected ? 'bg-gray-100 text-gray-600' : ''}
                        ${!isSubmitted && isSelected ? 'bg-blue-500 text-white' : ''}
                        ${showCorrect ? 'bg-green-500 text-white' : ''}
                        ${showWrong ? 'bg-red-500 text-white' : ''}
                        ${isSubmitted && !showCorrect && !showWrong ? 'bg-gray-200 text-gray-500' : ''}
                      `}>
                        {optionLetter}
                      </div>

                      <div className="flex-1">
                        <p className="text-[15px] text-gray-800 font-medium leading-relaxed">
                          {option}
                        </p>
                      </div>

                      {showCorrect && (
                        <CheckCircle className="flex-shrink-0 w-6 h-6 text-green-600" />
                      )}
                      {showWrong && (
                        <XCircle className="flex-shrink-0 w-6 h-6 text-red-600" />
                      )}
                    </label>
                  );
                })}
              </div>

              {!isSubmitted && (
                <div className="mt-6 flex justify-end">
                  <Button
                    onClick={handleSubmit}
                    disabled={!selectedAnswer || isSubmitting}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-8"
                  >
                    {isSubmitting ? 'Отправка...' : 'Ответить'}
                  </Button>
                </div>
              )}

                {isSubmitted && submitResult?.alreadySolved && (
                  <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
                    <p className="text-sm text-blue-700">
                      Вы уже решали этот вопрос ранее. Статистика не обновится.
                    </p>
                  </div>
                )}
              </motion.div>
            )}

            {/* Result Section (only for multiple choice) */}
            {isSubmitted && !question.isOpenEnded && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className={`rounded-xl p-6 md:p-8 border-2 ${
                  isCorrect ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'
                }`}
              >
                <div className="flex items-start gap-4">
                  <div className={`w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 ${
                    isCorrect ? 'bg-green-100' : 'bg-red-100'
                  }`}>
                    {isCorrect ? (
                      <CheckCircle className="w-6 h-6 text-green-600" />
                    ) : (
                      <XCircle className="w-6 h-6 text-red-600" />
                    )}
                  </div>
                  <div className="flex-1">
                    <h3 className={`text-lg font-bold mb-2 ${
                      isCorrect ? 'text-green-900' : 'text-red-900'
                    }`}>
                      {isCorrect ? 'Правильно!' : 'Неправильно'}
                    </h3>
                    {!isCorrect && submitResult && (
                      <p className="text-red-800 mb-4">
                        Правильный ответ: <strong>{submitResult.correctAnswer}</strong>
                      </p>
                    )}
                    {submitResult?.explanation && (
                      <div className={`mt-4 pt-4 border-t ${
                        isCorrect ? 'border-green-200' : 'border-red-200'
                      }`}>
                        <h4 className={`font-semibold mb-2 ${
                          isCorrect ? 'text-green-900' : 'text-red-900'
                        }`}>
                          Объяснение:
                        </h4>
                        <p className={`leading-relaxed ${
                          isCorrect ? 'text-green-800' : 'text-red-800'
                        }`}>
                          {submitResult.explanation}
                        </p>
                      </div>
                    )}
                  </div>
                </div>
              </motion.div>
            )}
          </div>

          {/* Right Sidebar - Related Questions */}
          <div className="hidden lg:block lg:col-span-4">
            <div className="sticky top-24 space-y-4">
              <h3 className="text-lg font-semibold text-gray-900">Похожие вопросы</h3>

              {relatedQuestions.length === 0 ? (
                <p className="text-sm text-gray-500">Нет похожих вопросов</p>
              ) : (
                <div className="space-y-3">
                  {relatedQuestions.map((q) => (
                    <Card
                      key={q.id}
                      className="border-none shadow-sm hover:shadow-md transition-shadow cursor-pointer bg-white group"
                      onClick={() => onSelectQuestion ? onSelectQuestion(String(q.id)) : null}
                    >
                      <CardContent className="p-4 space-y-3">
                        <div className="flex items-start justify-between gap-2">
                          <Badge
                            variant="outline"
                            className={`text-[10px] px-1.5 py-0 ${difficultyConfig[q.difficulty].color}`}
                          >
                            {difficultyConfig[q.difficulty].label}
                          </Badge>
                        </div>
                        <h4 className="font-medium text-gray-900 line-clamp-2 group-hover:text-blue-600 transition-colors text-sm">
                          {q.title}
                        </h4>
                        <div className="flex gap-2 flex-wrap">
                          <span className="text-xs text-gray-400 bg-gray-50 px-1.5 py-0.5 rounded">
                            {q.technology}
                          </span>
                          {q.companyTags?.slice(0, 1).map(tag => (
                            <span key={tag} className="text-xs text-gray-400 bg-gray-50 px-1.5 py-0.5 rounded">
                              {tag}
                            </span>
                          ))}
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Mobile Related Questions */}
        <div className="lg:hidden mt-8 space-y-4">
          <h3 className="text-lg font-semibold text-gray-900">Похожие вопросы</h3>
          <div className="overflow-x-auto pb-4 -mx-4 px-4">
            <div className="flex gap-4 w-max">
              {relatedQuestions.map((q) => (
                <div
                  key={q.id}
                  className="w-[280px] bg-white rounded-xl p-4 shadow-sm border border-gray-100 cursor-pointer hover:shadow-md transition-shadow"
                  onClick={() => onSelectQuestion ? onSelectQuestion(String(q.id)) : null}
                >
                  <Badge
                    variant="outline"
                    className={`mb-2 text-[10px] ${difficultyConfig[q.difficulty].color}`}
                  >
                    {difficultyConfig[q.difficulty].label}
                  </Badge>
                  <h4 className="font-medium text-gray-900 mb-3 line-clamp-2 text-sm">
                    {q.title}
                  </h4>
                  <div className="flex gap-2 flex-wrap">
                    <span className="text-xs text-gray-400 bg-gray-50 px-1.5 py-0.5 rounded">
                      {q.technology}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </main>

      {/* Fixed Bottom Action Bar */}
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 px-4 md:px-8 py-4 shadow-[0_-4px_20px_rgba(0,0,0,0.05)] z-50">
        <div className="max-w-7xl mx-auto flex items-center justify-end gap-4">

          {/* Показываем только одну кнопку в зависимости от состояния */}
          {isSubmitted ? (
            // После отправки ответа показываем "Следующий вопрос" или "К модулю"
            <Button
              onClick={onBack}
              className="flex-1 md:flex-none md:min-w-[200px] bg-gray-900 hover:bg-gray-800 text-white font-semibold h-12 rounded-xl shadow-lg transition-transform active:scale-[0.98]"
            >
              {courseId && moduleId ? 'К модулю' : 'Следующий вопрос'}
            </Button>
          ) : question.isOpenEnded ? (
            // Для открытых вопросов
            showAnswer ? (
              // Если ответ показан - кнопка "Изучил вопрос"
              <Button
                onClick={handleSubmit}
                disabled={isSubmitting}
                className="flex-1 md:flex-none md:min-w-[200px] bg-green-600 hover:bg-green-700 text-white font-semibold h-12 rounded-xl shadow-lg shadow-green-500/20 transition-transform active:scale-[0.98] disabled:opacity-50"
              >
                {isSubmitting ? 'Сохранение...' : 'Изучил вопрос'}
              </Button>
            ) : (
              // Если ответ не показан - кнопка "Показать ответ"
              <Button
                onClick={() => setShowAnswer(true)}
                className="flex-1 md:flex-none md:min-w-[200px] bg-[#007AFF] hover:bg-[#0063cc] text-white font-semibold h-12 rounded-xl shadow-lg shadow-blue-500/20 transition-transform active:scale-[0.98]"
              >
                Показать ответ
              </Button>
            )
          ) : (
            // Для вопросов с вариантами ответов - кнопка "Решить"
            <Button
              onClick={handleSubmit}
              disabled={!selectedAnswer || isSubmitting}
              className="flex-1 md:flex-none md:min-w-[200px] bg-[#007AFF] hover:bg-[#0063cc] text-white font-semibold h-12 rounded-xl shadow-lg shadow-blue-500/20 transition-transform active:scale-[0.98] disabled:bg-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed disabled:shadow-none"
            >
              {isSubmitting ? 'Отправка...' : 'Решить'}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
