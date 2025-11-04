import { useEffect, useState } from 'react'
import { Routes, Route, Link } from 'react-router-dom'
import AuthBar from './components/AuthBar'
import { listMentors, MentorCard } from './lib/api'
import AuthPage from './pages/Auth'

function Header() {
  const [dark, setDark] = useState<boolean>(() => (
    window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
  ))

  useEffect(() => {
    document.documentElement.classList.toggle('dark', dark)
  }, [dark])

  return (
    <header className="sticky top-0 z-40 border-b border-slate-200/60 bg-white/80 backdrop-blur dark:border-slate-800/60 dark:bg-slate-950/60">
      <div className="container flex h-16 items-center justify-between">
        <div className="flex items-center gap-3">
          <img src="/logo.jpg" alt="IT‑Rabotyagi" className="h-9 w-9 rounded-full ring-2 ring-white/70 object-cover" />
          <span className="text-lg font-bold">IT‑Rabotyagi</span>
        </div>
        <nav className="hidden items-center gap-6 md:flex">
          <a className="text-sm text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white" href="#questions">База вопросов</a>
          <a className="text-sm text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white" href="#courses">Курсы</a>
          <a className="text-sm text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white" href="#mentors">Менторы</a>
        </nav>
        <div className="flex items-center gap-3">
          <AuthBar />
          <button aria-label="Toggle theme" onClick={() => setDark(d => !d)} className="ml-1 rounded-md p-2 hover:bg-slate-100 dark:hover:bg-slate-800">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" className="text-slate-700 dark:text-slate-200"><path d="M12 3a9 9 0 1 0 9 9c0-.34-.02-.68-.06-1.01A7 7 0 0 1 12 3Z" stroke="currentColor" strokeWidth="2"/></svg>
          </button>
        </div>
      </div>
    </header>
  )
}

function Hero() {
  return (
    <section className="relative overflow-hidden">
      <div className="container py-16 md:py-24">
        <div className="mx-auto max-w-4xl">
          <h1 className="headline leading-tight">
            Прокачай свои
            <span className="block text-gradient">навыки в программировании</span>
          </h1>
          <p className="subhead mt-5 max-w-2xl">
            Изучай программирование с лучшими менторами. Практические курсы, реальные проекты и поддержка на каждом шаге обучения.
          </p>
          <div className="mt-8 flex flex-col items-start gap-3 sm:flex-row">
            <Link to="/auth" className="btn-primary">Начать обучение</Link>
            <a href="#courses" className="btn-secondary">Посмотреть курсы</a>
          </div>
        </div>

        {/* floating decorative elements (desktop) */}
        <div className="pointer-events-none relative mt-10 hidden h-[360px] md:block">
          {/* left purple card */}
          <div className="hero-card hero-card--purple glow absolute left-6 top-6 rotate-[-8deg] px-6 py-4">
            <div className="flex items-center gap-3">
              <span className="text-2xl">💻</span>
              <div className="text-sm leading-tight">
                <div className="font-semibold">Python</div>
              </div>
            </div>
          </div>

          {/* right blue card */}
          <div className="hero-card hero-card--blue glow absolute right-10 top-2 rotate-[8deg] px-6 py-4">
            <div className="flex items-center gap-3">
              <span className="text-2xl">⚡</span>
              <div className="text-sm leading-tight">
                <div className="font-semibold">JavaScript</div>
              </div>
            </div>
          </div>

          {/* bottom pill */}
          <div className="hero-pill absolute left-1/2 top-1/2 h-16 w-44 -translate-x-1/2 translate-y-24 rounded-2xl opacity-90 blur-[0.5px]" />

          {/* center circle with rocket */}
          <div className="hero-circle glow absolute left-1/2 top-1/2 flex h-56 w-56 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full text-white">
            <div className="text-center">
              <div className="text-5xl">🚀</div>
              <div className="mt-2 text-lg font-semibold">Старт обучения</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

function Features() {
  const items = [
    {
      title: 'База вопросов',
      desc: 'Постоянно обновляемая база вопросов от реальных собеседований в крупнейших компаниях. Готовьтесь эффективно.',
      icon: '📚',
      iconBg: 'bg-purple-500',
    },
    {
      title: 'Тренажер собеседований',
      desc: 'Отвечайте на вопросы и получайте обратную связь. Тренируйтесь столько, сколько нужно для уверенности.',
      icon: '💬',
      iconBg: 'bg-blue-500',
    },
    {
      title: 'Личные менторы',
      desc: 'Опытные разработчики помогут вам разобраться в сложных темах и подготовиться к карьерному росту.',
      icon: '👥',
      iconBg: 'bg-teal-500',
    },
  ]
  return (
    <section id="features" className="relative overflow-hidden">
      <div className="container py-16 md:py-24">
        <h2 className="text-3xl font-bold">Всё для успешной карьеры в IT</h2>
        <p className="subhead mt-2">Инструменты и ресурсы для подготовки к собеседованиям и развития навыков</p>
        <div className="mt-10 grid gap-6 md:grid-cols-3">
          {items.map((it) => (
            <div key={it.title} className="glass-card p-6 transition-all duration-300 hover:scale-105 hover:shadow-xl">
              <div className={`mb-4 flex h-16 w-16 items-center justify-center rounded-xl ${it.iconBg} text-3xl shadow-md`}>
                {it.icon}
              </div>
              <div className="text-xl font-semibold text-slate-900">{it.title}</div>
              <div className="mt-2 text-slate-600 dark:text-slate-300">{it.desc}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

function Mentors() {
  const [items, setItems] = useState<MentorCard[]>([])
  const [state, setState] = useState<'loading' | 'ready' | 'error'>('loading')

  useEffect(() => {
    listMentors()
      .then(({ items }) => {
        setItems(items)
        setState('ready')
      })
      .catch(() => setState('error'))
  }, [])

  return (
    <section id="mentors" className="relative overflow-hidden">
      <div className="container py-16 md:py-24">
        <div className="flex items-end justify-between">
          <div>
            <h2 className="text-3xl font-bold">Топ‑менторы</h2>
            <p className="subhead mt-2">Учитесь у тех, кто делает</p>
          </div>
          <a className="btn-secondary hidden md:inline-flex" href="#">Все менторы</a>
        </div>
        {state === 'loading' && <div className="mt-8 text-slate-500">Загрузка...</div>}
        {state === 'error' && <div className="mt-8 text-red-600">Не удалось загрузить менторов</div>}
        {state === 'ready' && (
          <div className="mt-8 grid gap-6 md:grid-cols-3">
            {items.map((m) => (
              <div key={m.id} className="glass-card p-6 transition-all duration-300 hover:scale-105 hover:shadow-xl">
                <div className="h-16 w-16 rounded-full bg-gradient-to-br from-brand-400 to-brand-600" />
                <div className="mt-4 text-xl font-semibold">{m.fullName}</div>
                <div className="text-slate-500">{m.title}</div>
                <div className="mt-3 flex flex-wrap gap-2">
                  {m.skills.map(s => (
                    <span key={s} className="rounded-md bg-slate-100 px-2.5 py-1 text-xs text-slate-700 dark:bg-slate-800 dark:text-slate-200">{s}</span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}

function Courses() {
  const [courses, setCourses] = useState<{ id: number; title: string; description: string }[]>([])
  const [state, setState] = useState<'loading' | 'ready' | 'error'>('loading')

  const imageByIndex = (i: number) => {
    const imgs = [
      'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?q=80&w=1200&auto=format&fit=crop',
      'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?q=80&w=1200&auto=format&fit=crop',
      'https://images.unsplash.com/photo-1518779578993-ec3579fee39f?q=80&w=1200&auto=format&fit=crop',
      'https://images.unsplash.com/photo-1519389950473-47ba0277781c?q=80&w=1200&auto=format&fit=crop',
    ]
    return imgs[i % imgs.length]
  }

  useEffect(() => {
    import('./lib/api').then(({ listCourses }) =>
      listCourses()
        .then(({ items }) => { setCourses(items as any); setState('ready') })
        .catch(() => setState('error'))
    )
  }, [])
  const courseIcons = ['🐍', '⚡', '🗄️']
  const courseGradients = [
    'from-blue-500 to-blue-600',
    'from-purple-500 to-pink-500',
    'from-teal-500 to-cyan-500',
  ]
  const courseDifficulties = ['Начальный', 'Средний', 'Средний']
  const courseStats = [
    { weeks: '8 недель', students: '2.5k', rating: '4.9' },
    { weeks: '16 недель', students: '3.2k', rating: '4.9' },
    { weeks: '6 недель', students: '1.8k', rating: '4.9' },
  ]

  return (
    <section id="courses" className="relative overflow-hidden">
      <div className="container py-16 md:py-24">
        <h2 className="text-3xl font-bold">Лайв курсы</h2>
        <p className="subhead mt-2">Выберите направление и начните свой путь в программировании уже сегодня</p>
        {state === 'loading' && <div className="mt-8 text-slate-500">Загрузка...</div>}
        {state === 'error' && <div className="mt-8 text-red-600">Не удалось загрузить курсы</div>}
        {state === 'ready' && (
        <div className="mt-8 grid gap-6 md:grid-cols-3">
          {courses.map((c, idx) => (
            <article key={c.title} className="group glass-card p-0 transition-all duration-300 hover:scale-105 hover:shadow-xl">
              <div className={`relative aspect-[16/9] w-full overflow-hidden rounded-t-2xl bg-gradient-to-br ${courseGradients[idx % courseGradients.length]}`}>
                <div className="flex h-full items-center justify-center text-6xl">
                  {courseIcons[idx % courseIcons.length]}
                </div>
                <div className="absolute left-3 top-3">
                  <span className="rounded-full bg-white/90 px-3 py-1 text-xs font-medium text-slate-700">
                    {courseDifficulties[idx % courseDifficulties.length]}
                  </span>
                </div>
              </div>
              <div className="p-6">
                <h3 className="text-xl font-semibold text-slate-900">{c.title}</h3>
                <p className="mt-2 text-sm text-slate-600">{c.description}</p>
                <div className="mt-4 flex items-center gap-4 text-xs text-slate-500">
                  <span className="flex items-center gap-1">
                    <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {courseStats[idx % courseStats.length].weeks}
                  </span>
                  <span className="flex items-center gap-1">
                    <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                    </svg>
                    {courseStats[idx % courseStats.length].students}
                  </span>
                  <span className="flex items-center gap-1">
                    <svg className="h-4 w-4 text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                    {courseStats[idx % courseStats.length].rating}
                  </span>
                </div>
                <div className="mt-4 flex items-center justify-between">
                  <span className="text-sm font-medium text-slate-700">Бесплатно</span>
                  <button className="rounded-lg bg-slate-800 px-4 py-2 text-sm font-semibold text-white transition hover:bg-slate-700">
                    Начать
                  </button>
                </div>
              </div>
            </article>
          ))}
        </div>
        )}
      </div>
    </section>
  )
}

function CTA() {
  return (
    <section className="container pb-20 pt-8">
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-600 to-brand-800 p-8 text-white shadow-xl">
        <div className="absolute right-0 top-0 h-40 w-40 -translate-y-1/3 translate-x-1/3 rounded-full bg-white/10 blur-2xl" />
        <h2 className="text-2xl font-bold md:text-3xl">Готовы начать путь в IT?</h2>
        <p className="mt-2 text-white/95">Присоединяйтесь к сообществу и учитесь у лучших.</p>
        <div className="mt-6 flex flex-col gap-3 sm:flex-row">
          <Link to="/auth" className="inline-flex items-center justify-center gap-2 rounded-lg border-2 border-white/50 bg-white/10 px-6 py-3 font-semibold text-white backdrop-blur transition hover:bg-white/20 focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-brand-600">
            Зарегистрироваться
          </Link>
          <Link to="/auth" className="inline-flex items-center justify-center gap-2 rounded-lg bg-white px-6 py-3 font-semibold text-brand-700 shadow-md transition hover:bg-slate-50 focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-brand-600">
            Попробовать бесплатно
          </Link>
        </div>
      </div>
    </section>
  )
}

function Footer() {
  return (
    <footer className="border-t border-slate-200 py-8 text-sm text-slate-600 dark:border-slate-800 dark:text-slate-300">
      <div className="container flex flex-col items-center justify-between gap-4 md:flex-row">
        <div className="flex items-center gap-3">
          <div className="h-8 w-8 rounded-md bg-brand-600" />
          <span className="font-semibold">IT‑RABOTYAGI</span>
        </div>
        <div className="text-slate-500">© {new Date().getFullYear()} IT‑RABOTYAGI. Все права защищены.</div>
      </div>
    </footer>
  )
}

function Home() {
  return (
    <>
      <Hero />
      <Features />
      <Mentors />
      <Courses />
      <CTA />
    </>
  )
}

export default function App() {
  return (
    <div>
      <Header />
      <main className="gradient-hero">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/auth" element={<AuthPage />} />
        </Routes>
      </main>
      <Footer />
    </div>
  )
}


