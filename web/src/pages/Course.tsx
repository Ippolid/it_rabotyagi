import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { listModules, Module, getCourse, Course } from '../lib/api'

export default function CoursePage() {
  const { id } = useParams()
  const courseId = Number(id)
  const [modules, setModules] = useState<Module[]>([])
  const [course, setCourse] = useState<Course | null>(null)
  const [state, setState] = useState<'loading'|'ready'|'error'>('loading')

  useEffect(() => {
    if (!courseId || Number.isNaN(courseId)) return
    Promise.all([
      getCourse(courseId).then(setCourse),
      listModules(courseId).then(({ items }) => setModules(items)),
    ])
      .then(() => setState('ready'))
      .catch(() => setState('error'))
  }, [courseId])

  return (
    <section className="container py-12">
      <div className="mb-6 flex items-center justify-between">
        <div>
          {state === 'ready' && course ? (
            <>
              <h1 className="text-2xl font-bold">{course.title}</h1>
              <p className="text-slate-600 dark:text-slate-300">{course.description}</p>
            </>
          ) : (
            <>
              <div className="skeleton h-6 w-64 rounded" />
              <div className="skeleton mt-2 h-4 w-96 rounded" />
            </>
          )}
        </div>
        <Link to="/" className="btn-secondary">← На главную</Link>
      </div>

      {state === 'loading' && (
        <div className="grid gap-6 md:grid-cols-2">
          {Array.from({length:4}).map((_,i)=> (
            <div key={i} className="glass-card p-6">
              <div className="skeleton h-4 w-48 rounded" />
              <div className="skeleton mt-2 h-3 w-64 rounded" />
            </div>
          ))}
        </div>
      )}

      {state === 'error' && (
        <div className="text-red-600">Не удалось загрузить модули</div>
      )}

      {state === 'ready' && (
        <div className="grid gap-6 md:grid-cols-2">
          {modules.map((m) => (
            <Link key={m.id} to={`/module/${m.id}`} className="glass-card block p-6 transition hover:shadow-xl">
              <div className="text-sm text-slate-500">Модуль {m.moduleOrder}</div>
              <div className="mt-1 text-lg font-semibold text-slate-900">{m.title}</div>
              <div className="mt-2 text-slate-600 dark:text-slate-300">{m.description}</div>
            </Link>
          ))}
        </div>
      )}
    </section>
  )
}
