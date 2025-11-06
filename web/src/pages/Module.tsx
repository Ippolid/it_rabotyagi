import { useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { listQuestions, Question, getModule, Module as ModuleType } from '../lib/api'

export default function ModulePage() {
  const { id } = useParams()
  const moduleId = Number(id)
  const [moduleInfo, setModuleInfo] = useState<ModuleType | null>(null)
  const [questions, setQuestions] = useState<Question[]>([])
  const [state, setState] = useState<'loading'|'ready'|'error'>('loading')
  const [openId, setOpenId] = useState<number|null>(null)
  const [answers, setAnswers] = useState<Record<number, string>>({})

  useEffect(() => {
    if (!moduleId || Number.isNaN(moduleId)) return
    Promise.all([
      getModule(moduleId).then(setModuleInfo),
      listQuestions(moduleId).then(({ items }) => setQuestions(items)),
    ])
      .then(() => setState('ready'))
      .catch(() => setState('error'))
  }, [moduleId])

  const title = useMemo(() => moduleInfo?.title ? moduleInfo.title : `Модуль #${moduleId}`, [moduleInfo, moduleId])
  const subtitle = useMemo(() => moduleInfo?.description ?? 'Вопросы модуля', [moduleInfo])

  function selectAnswer(qid: number, value: string) {
    setAnswers(a => ({ ...a, [qid]: value }))
  }

  return (
    <section className="container py-12">
      <div className="mb-6 flex items-center justify-between">
        <div>
          {state === 'ready' ? (
            <>
              <h1 className="text-2xl font-bold">{title}</h1>
              <p className="text-slate-600 dark:text-slate-300">{subtitle}</p>
            </>
          ) : (
            <>
              <div className="skeleton h-6 w-56 rounded" />
              <div className="skeleton mt-2 h-4 w-80 rounded" />
            </>
          )}
        </div>
        <div className="flex gap-2">
          <Link to="/" className="btn-secondary">← На главную</Link>
        </div>
      </div>

      {state === 'loading' && (
        <div className="grid gap-6 md:grid-cols-2">
          {Array.from({length:6}).map((_,i)=> (
            <div key={i} className="glass-card p-6">
              <div className="skeleton h-4 w-72 rounded" />
              <div className="skeleton mt-2 h-3 w-48 rounded" />
            </div>
          ))}
        </div>
      )}

      {state === 'error' && (
        <div className="text-red-600">Не удалось загрузить вопросы</div>
      )}

      {state === 'ready' && (
        <div className="grid gap-6 md:grid-cols-2">
          {questions.map((q) => {
            const selected = answers[q.id]
            const isChoice = (q.options && q.options.length > 0)
            const isCorrect = isChoice && q.correctAnswer && selected ? (selected.trim() === q.correctAnswer.trim()) : undefined
            return (
              <div key={q.id} className="glass-card p-6">
                <div className="flex items-center justify-between">
                  <div className="text-lg font-semibold text-slate-900">{q.title}</div>
                  {q.difficulty && (
                    <span className="rounded-full px-3 py-1 text-xs font-medium text-white" style={{background:'rgba(0,0,0,0.35)'}}>
                      {q.difficulty}
                    </span>
                  )}
                </div>

                <div className="mt-3 flex gap-2">
                <button className="btn-primary" onClick={() => setOpenId(openId===q.id?null:q.id)}>
                  {openId===q.id ? 'Скрыть' : 'Открыть'}
                </button>
                <Link to={`/trainer/${moduleId}/${q.id}`} className="btn-secondary">Тренироваться</Link>
                </div>

                {openId===q.id && (
                  <div className="mt-4 space-y-4">
                    <div className="whitespace-pre-wrap text-slate-700 dark:text-slate-200">{q.content}</div>

                    {isChoice && (
                      <div className="space-y-2">
                        {q.options!.map(opt => (
                          <label key={opt} className="flex cursor-pointer items-center gap-3 rounded-md border border-slate-200 p-3 dark:border-slate-700">
                            <input
                              type="radio"
                              name={`q-${q.id}`}
                              checked={selected === opt}
                              onChange={() => selectAnswer(q.id, opt)}
                            />
                            <span className="text-sm">{opt}</span>
                          </label>
                        ))}
                      </div>
                    )}

                    {isChoice && selected && (
                      <div className={`rounded-md px-3 py-2 text-sm ${isCorrect ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-200' : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-200'}`}>
                        {isCorrect ? 'Верно' : 'Неверно'}
                      </div>
                    )}

                    {q.explanation && (
                      <details className="rounded-md border border-slate-200 p-3 text-sm dark:border-slate-700">
                        <summary className="cursor-pointer select-none font-medium">Пояснение</summary>
                        <div className="mt-2 whitespace-pre-wrap text-slate-700 dark:text-slate-200">{q.explanation}</div>
                      </details>
                    )}
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}
    </section>
  )
}
