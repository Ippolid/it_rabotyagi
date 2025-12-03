
import React, { useEffect, useState, useRef } from 'react';
import { Button } from '../ui/button';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/card';
import { Badge } from '../ui/badge';
import { Alert, AlertDescription } from '../ui/alert';
import { Progress } from '../ui/progress';
import { Mic, Square, Play, RotateCw, Award, TrendingUp, CheckCircle2, XCircle, Loader2, Brain, ArrowRight, Sparkles } from 'lucide-react';
import { motion, AnimatePresence } from 'motion/react';
import { getRandomMLQuestion, completeInterview, MLQuestion, CompleteInterviewResponse } from '../../lib/api';

type RecordingState = 'idle' | 'recording' | 'processing' | 'completed' | 'error';
type QuestionBase = 'ml' | null;
type ViewState = 'selection' | 'configuration' | 'interview';

// Question base configurations
const QUESTION_BASES = [
  {
    id: 'ml' as const,
    title: 'Machine Learning & Deep Learning',
    description: 'Вопросы по ML и DL',
    icon: Brain,
    questionsCount: 39,
    difficulty: 'Средний - Сложный',
    color: 'from-purple-500 to-pink-500',
    bgColor: 'bg-purple-50',
    borderColor: 'border-purple-200',
    iconColor: 'text-purple-600',
  },
  // Можно легко добавить другие базы в будущем
  // {
  //   id: 'python',
  //   title: 'Python',
  //   description: 'Вопросы по Python, синтаксису, библиотекам',
  //   icon: Code,
  //   questionsCount: 0,
  //   difficulty: 'Все уровни',
  //   color: 'from-blue-500 to-cyan-500',
  //   bgColor: 'bg-blue-50',
  //   borderColor: 'border-blue-200',
  //   iconColor: 'text-blue-600',
  //   disabled: true,
  // },
];

export function InterviewTrainer() {
  const [viewState, setViewState] = useState<ViewState>('selection');
  const [selectedBase, setSelectedBase] = useState<QuestionBase>(null);
  const [questionsCount, setQuestionsCount] = useState<number>(5);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState<number>(0);
  const [question, setQuestion] = useState<MLQuestion | null>(null);
  const [loadingQuestion, setLoadingQuestion] = useState(false);
  const [recordingState, setRecordingState] = useState<RecordingState>('idle');
  const [recordingTime, setRecordingTime] = useState(0);
  const [result, setResult] = useState<CompleteInterviewResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const audioChunksRef = useRef<Blob[]>([]);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  const handleSelectBase = (baseId: QuestionBase) => {
    setSelectedBase(baseId);
    setViewState('configuration');
  };

  const handleStartInterview = () => {
    setViewState('interview');
    setCurrentQuestionIndex(0);
    loadNewQuestion();
  };

  const handleBackToSelection = () => {
    setViewState('selection');
    setSelectedBase(null);
    setQuestion(null);
    setResult(null);
    setError(null);
    setRecordingState('idle');
    setRecordingTime(0);
    setCurrentQuestionIndex(0);
  };

  const handleCancelRecording = () => {
    if (mediaRecorderRef.current && mediaRecorderRef.current.state === 'recording') {
      mediaRecorderRef.current.stop();
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    }
    audioChunksRef.current = [];
    setRecordingState('idle');
    setRecordingTime(0);
  };

  const loadNewQuestion = async () => {
    setLoadingQuestion(true);
    setError(null);
    setResult(null);
    setRecordingState('idle');
    setRecordingTime(0);

    try {
      const q = await getRandomMLQuestion();
      setQuestion(q);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не удалось загрузить вопрос';
      setError(message);
    } finally {
      setLoadingQuestion(false);
    }
  };

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const mediaRecorder = new MediaRecorder(stream, {
        mimeType: 'audio/webm',
      });

      audioChunksRef.current = [];

      mediaRecorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          audioChunksRef.current.push(event.data);
        }
      };

      mediaRecorder.onstop = async () => {
        stream.getTracks().forEach(track => track.stop());
        await processRecording();
      };

      mediaRecorder.start();
      mediaRecorderRef.current = mediaRecorder;
      setRecordingState('recording');
      setRecordingTime(0);

      // Start timer
      timerRef.current = setInterval(() => {
        setRecordingTime(prev => prev + 1);
      }, 1000);

    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не удалось получить доступ к микрофону';
      setError(message);
      setRecordingState('error');
    }
  };

  const stopRecording = () => {
    if (mediaRecorderRef.current && mediaRecorderRef.current.state === 'recording') {
      mediaRecorderRef.current.stop();
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    }
  };

  const processRecording = async () => {
    if (!question || audioChunksRef.current.length === 0) {
      setError('Нет записанного аудио');
      setRecordingState('error');
      return;
    }

    setRecordingState('processing');
    setError(null);

    try {
      const audioBlob = new Blob(audioChunksRef.current, { type: 'audio/webm' });
      const response = await completeInterview(
        audioBlob,
        question.id,
        undefined,
        undefined,
        'ru'
      );
      
      setResult(response);
      setRecordingState('completed');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Ошибка обработки записи';
      setError(message);
      setRecordingState('error');
    }
  };

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const getScoreColor = (score: number): string => {
    if (score >= 80) return 'text-green-600';
    if (score >= 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getScoreBgColor = (score: number): string => {
    if (score >= 80) return 'bg-green-100';
    if (score >= 60) return 'bg-yellow-100';
    return 'bg-red-100';
  };

  // 1. Course Selection Screen
  if (viewState === 'selection') {
    return (
      <div className="max-w-5xl mx-auto space-y-8">
        {/* Header */}
        <div className="text-center space-y-4">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
          >
            <h1 className="text-4xl md:text-5xl font-bold text-gray-900">Тренажер собеседований</h1>
            <p className="text-lg text-gray-500 mt-4 max-w-2xl mx-auto">
              Практикуйтесь устно отвечать на вопросы и получайте детальный фидбек от AI
            </p>
          </motion.div>
        </div>

        {/* Course Selection */}
        <div className="space-y-4">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.2 }}
          >
            <h2 className="text-xl font-semibold text-gray-900 mb-6">Выберите курс</h2>
          </motion.div>

          <div className="grid md:grid-cols-1 gap-6">
            {QUESTION_BASES.map((base, index) => (
              <motion.div
                key={base.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 * (index + 1) }}
              >
                <Card
                  className={`${base.bgColor} ${base.borderColor} border-2 hover:shadow-xl transition-all duration-300 cursor-pointer group relative overflow-hidden`}
                  onClick={() => handleSelectBase(base.id)}
                >
                  {/* Gradient Background */}
                  <div className={`absolute inset-0 bg-gradient-to-br ${base.color} opacity-0 group-hover:opacity-10 transition-opacity duration-300`} />
                  
                  <CardContent className="p-8 relative">
                    <div className="flex items-start gap-6">
                      {/* Icon */}
                      <div className={`flex-shrink-0 w-16 h-16 rounded-2xl bg-white shadow-sm flex items-center justify-center group-hover:scale-110 transition-transform duration-300`}>
                        <base.icon className={`w-8 h-8 ${base.iconColor}`} />
                      </div>

                      {/* Content */}
                      <div className="flex-1 space-y-3">
                        <div className="flex items-center gap-3">
                          <h3 className="text-2xl font-bold text-gray-900 group-hover:text-purple-600 transition-colors">
                            {base.title}
                          </h3>
                          <Badge className="bg-purple-600 text-white border-0">
                            {base.questionsCount} вопросов
                          </Badge>
                        </div>
                        
                        <p className="text-gray-600 leading-relaxed">
                          {base.description}
                        </p>

                        <div className="flex items-center gap-4 text-sm">
                          <div className="flex items-center gap-2 text-gray-500">
                            <Sparkles className="w-4 h-4" />
                            <span>{base.difficulty}</span>
                          </div>
                        </div>
                      </div>

                      {/* Arrow */}
                      <div className="flex-shrink-0 self-center">
                        <ArrowRight className="w-6 h-6 text-gray-400 group-hover:text-purple-600 group-hover:translate-x-1 transition-all duration-300" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </div>

        {/* Info Section */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4 }}
          className="text-center"
        >
          <Card className="bg-white border-gray-200 shadow-sm">
            <CardContent className="p-8">
              <div className="space-y-4">
                <div className="flex items-center justify-center gap-2 mb-4">
                  <Sparkles className="w-5 h-5 text-purple-600" />
                  <h3 className="text-lg font-semibold text-gray-900">Как это работает</h3>
                </div>
                <div className="grid md:grid-cols-3 gap-4">
                  <div className="p-4 rounded-xl border-2 border-gray-200 hover:border-purple-300 transition-colors">
                    <div className="w-8 h-8 rounded-lg bg-purple-600 text-white flex items-center justify-center font-bold mb-3">1</div>
                    <p className="text-sm text-gray-700 font-medium">Выберите курс для тренировки</p>
                  </div>
                  <div className="p-4 rounded-xl border-2 border-gray-200 hover:border-purple-300 transition-colors">
                    <div className="w-8 h-8 rounded-lg bg-purple-600 text-white flex items-center justify-center font-bold mb-3">2</div>
                    <p className="text-sm text-gray-700 font-medium">Настройте количество вопросов на собеседование</p>
                  </div>
                  <div className="p-4 rounded-xl border-2 border-gray-200 hover:border-purple-300 transition-colors">
                    <div className="w-8 h-8 rounded-lg bg-purple-600 text-white flex items-center justify-center font-bold mb-3">3</div>
                    <p className="text-sm text-gray-700 font-medium">Получите вопрос и нажмите "Начать запись"</p>
                  </div>
                  <div className="p-4 rounded-xl border-2 border-gray-200 hover:border-purple-300 transition-colors">
                    <div className="w-8 h-8 rounded-lg bg-purple-600 text-white flex items-center justify-center font-bold mb-3">4</div>
                    <p className="text-sm text-gray-700 font-medium">Ответьте устно на вопрос</p>
                  </div>
                  <div className="p-4 rounded-xl border-2 border-gray-200 hover:border-purple-300 transition-colors">
                    <div className="w-8 h-8 rounded-lg bg-purple-600 text-white flex items-center justify-center font-bold mb-3">5</div>
                    <p className="text-sm text-gray-700 font-medium">Нажмите "Остановить и отправить" для получения AI-фидбека</p>
                  </div>
                  <div className="p-4 rounded-xl border-2 border-gray-200 hover:border-purple-300 transition-colors">
                    <div className="w-8 h-8 rounded-lg bg-purple-600 text-white flex items-center justify-center font-bold mb-3">6</div>
                    <p className="text-sm text-gray-700 font-medium">Получите оценку 0-100, разбор сильных сторон и рекомендации по улучшению</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    );
  }

  // 2. Configuration Screen
  if (viewState === 'configuration') {
    const selectedCourse = QUESTION_BASES.find(b => b.id === selectedBase);
    
    return (
      <div className="max-w-3xl mx-auto space-y-8">
        <div className="flex items-center gap-4">
          <Button
            onClick={handleBackToSelection}
            variant="ghost"
            size="sm"
          >
            ← Назад
          </Button>
        </div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center space-y-4"
        >
          <h1 className="text-3xl font-bold text-gray-900">{selectedCourse?.title}</h1>
          <p className="text-gray-500">{selectedCourse?.description}</p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
        >
          <Card>
            <CardHeader>
              <CardTitle>Настройки собеседования</CardTitle>
            </CardHeader>
            <CardContent className="space-y-8">
              <div className="space-y-6">
                <div className="text-center space-y-6">
                  <span className="text-base font-semibold text-gray-800 block">
                    Количество вопросов
                  </span>
                  <div className="flex flex-wrap justify-center gap-4">
                    {[5, 10, 20, 30].map((num) => (
                      <button
                        key={num}
                        onClick={() => setQuestionsCount(num)}
                        style={questionsCount === num ? { backgroundColor: '#9333ea', color: 'white' } : {}}
                        className={`min-w-[100px] px-8 py-6 rounded-2xl font-bold text-2xl transition-all duration-200 ${
                          questionsCount === num
                            ? 'shadow-2xl scale-105 ring-4 ring-purple-200'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200 hover:scale-105 shadow-md'
                        }`}
                      >
                        {num}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <Button
                onClick={handleStartInterview}
                size="lg"
                style={{ backgroundColor: '#9333ea', color: 'white' }}
                className="w-full hover:bg-purple-700"
              >
                <Play className="mr-2 h-5 w-5" />
                Начать собеседование
              </Button>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    );
  }

  // 3. Interview Screen
  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header with Progress */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <Button
            onClick={handleBackToSelection}
            variant="ghost"
            size="sm"
          >
            ← Завершить собеседование
          </Button>
          <Badge className="bg-purple-100 text-purple-700 border-purple-200 text-sm px-3 py-1">
            Вопрос {currentQuestionIndex + 1} из {questionsCount}
          </Badge>
        </div>
        
        <div className="space-y-2">
          <div className="flex justify-between text-sm text-gray-500">
            <span>Прогресс собеседования</span>
            <span>{Math.round(((currentQuestionIndex + 1) / questionsCount) * 100)}%</span>
          </div>
          <Progress value={((currentQuestionIndex + 1) / questionsCount) * 100} className="h-2" />
        </div>
      </div>

      {/* Error Alert */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
          >
            <Alert variant="destructive">
              <XCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Loading State */}
      {loadingQuestion && (
        <Card>
          <CardContent className="py-12 text-center">
            <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4 text-gray-400" />
            <p className="text-gray-500">Загрузка вопроса...</p>
          </CardContent>
        </Card>
      )}

      {/* Question Card */}
      {question && !loadingQuestion && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <Card className="border-2 border-blue-200 bg-blue-50/30">
            <CardHeader>
              <div className="flex items-center justify-between mb-2">
                <Badge className="bg-blue-100 text-blue-700 border-blue-200">
                  Вопрос #{question.id}
                </Badge>
              </div>
              <CardTitle className="text-xl text-gray-900 leading-relaxed">
                {question.question}
              </CardTitle>
            </CardHeader>
          </Card>
        </motion.div>
      )}

      {/* Recording Controls */}
      {question && !loadingQuestion && recordingState !== 'completed' && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
        >
          <Card className="border-2">
            <CardContent className="py-12">
              <div className="flex flex-col items-center gap-8">
                
                {/* Idle State */}
                {recordingState === 'idle' && (
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="text-center space-y-6"
                  >
                    <div className="relative">
                      <motion.div
                        animate={{ scale: [1, 1.1, 1] }}
                        transition={{ duration: 2, repeat: Infinity }}
                        className="w-24 h-24 rounded-full bg-red-100 flex items-center justify-center mx-auto"
                      >
                        <Mic className="h-12 w-12 text-red-600" />
                      </motion.div>
                    </div>
                    <div>
                      <p className="text-xl font-medium text-gray-900 mb-2">
                        Готовы ответить на вопрос?
                      </p>
                      <p className="text-gray-500">
                        Нажмите кнопку ниже и начните говорить
                      </p>
                    </div>
                    <Button
                      onClick={startRecording}
                      size="lg"
                      style={{ backgroundColor: '#dc2626', color: 'white', padding: '1.5rem 2.5rem', fontSize: '1.125rem' }}
                      className="rounded-xl hover:bg-red-700"
                    >
                      <Mic className="mr-2 h-6 w-6" />
                      Начать запись
                    </Button>
                  </motion.div>
                )}

                {/* Recording State */}
                {recordingState === 'recording' && (
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="text-center space-y-6 w-full"
                  >
                    <div className="relative">
                      <motion.div
                        animate={{ 
                          scale: [1, 1.2, 1],
                          boxShadow: [
                            '0 0 0 0 rgba(220, 38, 38, 0.7)',
                            '0 0 0 20px rgba(220, 38, 38, 0)',
                            '0 0 0 0 rgba(220, 38, 38, 0)'
                          ]
                        }}
                        transition={{ duration: 1.5, repeat: Infinity }}
                        className="w-32 h-32 rounded-full bg-red-600 flex items-center justify-center mx-auto"
                      >
                        <Mic className="h-16 w-16 text-white" />
                      </motion.div>
                    </div>
                    
                    <div>
                      <div className="text-6xl font-bold text-red-600 mb-3 font-mono">
                        {formatTime(recordingTime)}
                      </div>
                      <p className="text-lg text-gray-600 mb-6">Идет запись...</p>
                    </div>

                    <Button
                      onClick={stopRecording}
                      size="lg"
                      style={{ backgroundColor: '#16a34a', color: 'white', padding: '1.5rem 2.5rem', fontSize: '1.125rem' }}
                      className="rounded-xl shadow-lg hover:bg-green-700"
                    >
                      <CheckCircle2 className="mr-2 h-6 w-6" />
                      Остановить и отправить на проверку
                    </Button>
                  </motion.div>
                )}

                {/* Processing State */}
                {recordingState === 'processing' && (
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="text-center space-y-6"
                  >
                    <Loader2 className="h-16 w-16 animate-spin mx-auto text-purple-600" />
                    <div>
                      <p className="text-2xl font-bold text-gray-900 mb-2">
                        Обрабатываем ваш ответ
                      </p>
                      <p className="text-gray-500 mb-4">
                        Это может занять 10-30 секунд
                      </p>
                      <div className="text-sm text-gray-400 space-y-1">
                        <p>• Транскрибируем аудио с помощью Whisper</p>
                        <p>• Анализируем ответ с помощью Llama 3.3 70B</p>
                        <p>• Формируем детальный фидбек</p>
                      </div>
                    </div>
                  </motion.div>
                )}

                {/* Error State */}
                {recordingState === 'error' && (
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    className="text-center space-y-4"
                  >
                    <XCircle className="h-16 w-16 mx-auto text-red-500" />
                    <p className="text-lg text-gray-700">
                      Попробуйте записать ответ еще раз
                    </p>
                    <Button
                      onClick={() => setRecordingState('idle')}
                      variant="outline"
                    >
                      Попробовать снова
                    </Button>
                  </motion.div>
                )}

              </div>
            </CardContent>
          </Card>
        </motion.div>
      )}

      {/* Results */}
      {result && recordingState === 'completed' && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="space-y-6"
        >
          {/* Score Card */}
          <Card className={`border-2 ${getScoreBgColor(result.evaluation.score)}`}>
            <CardContent className="py-8">
              <div className="text-center">
                <Award className={`h-16 w-16 mx-auto mb-4 ${getScoreColor(result.evaluation.score)}`} />
                <div className={`text-6xl font-bold mb-2 ${getScoreColor(result.evaluation.score)}`}>
                  {result.evaluation.score}
                </div>
                <p className="text-gray-600 mb-4">баллов из 100</p>
                <Progress value={result.evaluation.score} className="h-3" />
              </div>
            </CardContent>
          </Card>

          {/* Transcription */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Play className="h-5 w-5 text-blue-600" />
                Ваш ответ (транскрипция)
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-gray-700 leading-relaxed bg-gray-50 p-4 rounded-lg">
                {result.transcription.text}
              </p>
            </CardContent>
          </Card>

          {/* Strengths */}
          {result.evaluation.strengths.length > 0 && (
            <Card className="border-green-200 bg-green-50/30">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-green-700">
                  <CheckCircle2 className="h-5 w-5" />
                  Что сказали правильно
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2">
                  {result.evaluation.strengths.map((strength, idx) => (
                    <li key={idx} className="flex items-start gap-2">
                      <span className="text-green-600 mt-1">✓</span>
                      <span className="text-gray-700">{strength}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          )}

          {/* Improvements */}
          {result.evaluation.improvements.length > 0 && (
            <Card className="border-yellow-200 bg-yellow-50/30">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-yellow-700">
                  <TrendingUp className="h-5 w-5" />
                  Что можно улучшить
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2">
                  {result.evaluation.improvements.map((improvement, idx) => (
                    <li key={idx} className="flex items-start gap-2">
                      <span className="text-yellow-600 mt-1">→</span>
                      <span className="text-gray-700">{improvement}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          )}

          {/* Overall Comment */}
          {result.evaluation.overall_comment && (
            <Card>
              <CardHeader>
                <CardTitle>Общая рекомендация</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-700 leading-relaxed">
                  {result.evaluation.overall_comment}
                </p>
              </CardContent>
            </Card>
          )}

          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            {currentQuestionIndex + 1 < questionsCount ? (
              <Button 
                onClick={() => {
                  setCurrentQuestionIndex(prev => prev + 1);
                  loadNewQuestion();
                }} 
                size="lg"
                style={{ backgroundColor: '#9333ea', color: 'white', padding: '1.5rem 2rem' }}
                className="hover:bg-purple-700 shadow-lg"
              >
                <ArrowRight className="mr-2 h-5 w-5" />
                Следующий вопрос
              </Button>
            ) : (
              <Button 
                onClick={handleBackToSelection}
                size="lg"
                style={{ backgroundColor: '#16a34a', color: 'white', padding: '1.5rem 2rem' }}
                className="hover:bg-green-700 shadow-lg"
              >
                <CheckCircle2 className="mr-2 h-5 w-5" />
                Завершить собеседование
              </Button>
            )}
            <Button
              onClick={() => {
                setRecordingState('idle');
                setResult(null);
                setRecordingTime(0);
              }}
              variant="outline"
              size="lg"
              style={{ padding: '1.5rem 2rem' }}
              className="border-2 border-gray-300 hover:bg-gray-100"
            >
              <RotateCw className="mr-2 h-5 w-5" />
              Попробовать снова
            </Button>
          </div>
        </motion.div>
      )}
    </div>
  );
}
