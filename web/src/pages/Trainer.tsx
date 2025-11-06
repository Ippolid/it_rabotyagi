import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Course, listCourses, listModules, listQuestions, Question } from '../lib/api'
import AudioRecorder from '../components/AudioRecorder'

export default function TrainerPage() {
  const [courses, setCourses] = useState<Course[]>([])
  const [selectedCourse, setSelectedCourse] = useState<number | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [current, setCurrent] = useState<Question | null>(null)
  const [showAnswer, setShowAnswer] = useState(false)
  const [transcript, setTranscript] = useState('')

  useEffect(() => {
    listCourses()
      .then(({ items }) => setCourses(items))
      .catch(() => setCourses([]))
  }, [])

  async function nextQuestion() {
    if (!selectedCourse) return
    setError('')
    setLoading(true)
    setShowAnswer(false)
    setTranscript('')
    try {
      // 1) берем список модулей курса
      const mods = await listModules(selectedCourse)
      if (!mods.items.length) throw new Error('В курсе нет модулей')
      // 2) случайный модуль
      const m = mods.items[Math.floor(Math.random() * mods.items.length)]
      // 3) вопросы модуля
      const qs = await listQuestions(m.id)
      if (!qs.items.length) throw new Error('В модуле нет вопросов')
      // 4) случайный вопрос
      setCurrent(qs.items[Math.floor(Math.random() * qs.items.length)])
    } catch (e: any) {
      setError(e?.message || 'Не удалось загрузить вопрос')
    } finally {
      setLoading(false)
    }
  }

  const referenceText = useMemo(() => {
    if (!current) return ''
    return current.explanation || current.correctAnswer || current.content || ''
  }, [current])

  return (
    <section className="container py-12">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Тренажер собеседований</h1>
          <p className="text-slate-600 dark:text-slate-300">Выберите курс и отвечайте на вопросы голосом. Транскрипция и семантическая проверка будут добавлены позже.</p>
        </div>
        <Link to="/" className="btn-secondary">← На главную</Link>
      </div>

      <div className="glass-card p-6">
        <div className="grid gap-4 md:grid-cols-[1fr_auto] md:items-end">
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700 dark:text-slate-200">Курс</label>
            <select
              className="w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-slate-900 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-100"
              value={selectedCourse ?? ''}
              onChange={(e) => setSelectedCourse(Number(e.target.value) || null)}
            >
              <option value="">Выберите курс</option>
              {courses.map(c => (
                <option key={c.id} value={c.id}>{c.title}</option>
              ))}
            </select>
          </div>
          <button disabled={!selectedCourse || loading} onClick={nextQuestion} className="btn-primary justify-center">{loading ? 'Загрузка...' : 'Случайный вопрос'}</button>
        </div>
      </div>

      {error && <div className="mt-4 rounded-md bg-red-100 px-3 py-2 text-sm text-red-800 dark:bg-red-900/30 dark:text-red-200">{error}</div>}

      {current && (
        <div className="mt-6 grid gap-6 md:grid-cols-2">
          <div className="glass-card p-6">
            <div className="flex items-center justify-between">
              <div className="text-lg font-semibold text-slate-900 dark:text-slate-100">{current.title}</div>
              {current.difficulty && (
                <span className="rounded-full px-3 py-1 text-xs font-medium text-white" style={{background:'rgba(0,0,0,0.35)'}}>{current.difficulty}</span>
              )}
            </div>
            <div className="mt-3 whitespace-pre-wrap text-slate-700 dark:text-slate-200">{current.content}</div>

            {current.options && current.options.length > 0 && (
              <ul className="mt-4 list-disc space-y-2 pl-6 text-sm text-slate-700 dark:text-slate-300">
                {current.options.map(o => <li key={o}>{o}</li>)}
              </ul>
            )}

            <div className="mt-4 flex flex-wrap gap-3">
              <button className="btn-secondary" onClick={() => setShowAnswer(s => !s)}>{showAnswer ? 'Скрыть ответ' : 'Показать ответ'}</button>
              <button className="btn-secondary" onClick={nextQuestion}>Следующий</button>
            </div>

            {showAnswer && (
              <div className="mt-4 rounded-md border border-slate-200 p-3 text-sm dark:border-slate-700">
                <div className="font-semibold">Эталонный ответ</div>
                <div className="mt-2 whitespace-pre-wrap">{referenceText}</div>
              </div>
            )}
          </div>

          <div className="glass-card p-6">
            <div className="text-lg font-semibold">Ваш ответ (голосом)</div>
            <p className="mt-1 text-sm text-slate-600 dark:text-slate-300">Нажмите запись и отвечайте устно. Транскрипция появится ниже.</p>
            <div className="mt-3">
              <AudioRecorder onTranscript={(t)=> setTranscript(t)} />
            </div>
            {transcript && (
              <div className="mt-4 rounded-md border border-slate-200 p-3 text-sm dark:border-slate-700">
                <div className="font-semibold">Транскрипт</div>
                <div className="mt-2 whitespace-pre-wrap">{transcript}</div>
              </div>
            )}
          </div>
        </div>
      )}
    </section>
  )
}
