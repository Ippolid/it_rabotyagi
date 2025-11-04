import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { login, register, getMe } from '../lib/api'

export default function AuthPage() {
  const nav = useNavigate()
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [email, setEmail] = useState('')
  const [nickname, setNickname] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      if (mode === 'login') {
        await login(email, password)
      } else {
        await register(email, nickname || email.split('@')[0], password)
      }
      await getMe()
      nav('/')
      window.location.reload()
    } catch (err: any) {
      setError(err?.message || (mode === 'login' ? 'Ошибка входа' : 'Ошибка регистрации'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen">
      <div className="container flex min-h-screen items-center justify-center py-16">
        <div className="w-full max-w-md rounded-2xl bg-white/80 p-8 shadow-xl backdrop-blur">
            <div className="mb-6 flex items-center gap-3">
              <img src="/logo.jpg" alt="Logo" className="h-10 w-10 rounded-full" />
              <div>
                <div className="text-xl font-bold text-slate-900">{mode === 'login' ? 'Вход' : 'Регистрация'}</div>
                <div className="text-sm text-slate-600">
                  {mode === 'login' ? 'Войдите в свой аккаунт' : 'Создайте аккаунт и начните обучение'}
                </div>
              </div>
            </div>
            <form className="space-y-4" onSubmit={handleSubmit}>
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Email</label>
                <input
                  required
                  type="email"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  className="w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 placeholder-slate-400 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
                  placeholder="you@example.com"
                />
              </div>
              {mode === 'register' && (
                <div>
                  <label className="mb-1 block text-sm font-medium text-slate-700">Никнейм</label>
                  <input
                    required
                    value={nickname}
                    onChange={e => setNickname(e.target.value)}
                    className="w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 placeholder-slate-400 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
                    placeholder="username"
                  />
                </div>
              )}
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Пароль</label>
                <input
                  required
                  type="password"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  className="w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 placeholder-slate-400 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500"
                  placeholder="••••••••"
                />
              </div>
              {error && <div className="text-sm text-red-600">{error}</div>}
              <button disabled={loading} className="btn-primary w-full justify-center">
                {loading ? (mode === 'login' ? 'Вход...' : 'Создание...') : (mode === 'login' ? 'Войти' : 'Зарегистрироваться')}
              </button>
            </form>
            <div className="mt-6 text-center text-sm text-slate-600">
              {mode === 'login' ? (
                <>
                  Нет аккаунта?{' '}
                  <button
                    onClick={() => {
                      setMode('register')
                      setError('')
                    }}
                    className="font-medium text-brand-600 hover:text-brand-700"
                  >
                    Зарегистрироваться
                  </button>
                </>
              ) : (
                <>
                  Уже есть аккаунт?{' '}
                  <button
                    onClick={() => {
                      setMode('login')
                      setError('')
                    }}
                    className="font-medium text-brand-600 hover:text-brand-700"
                  >
                    Войти
                  </button>
                </>
              )}
            </div>
            <div className="mt-4 text-center">
              <Link to="/" className="text-sm text-slate-500 hover:text-slate-700">
                ← Вернуться на главную
              </Link>
            </div>
          </div>
        </div>
    </div>
  )
}

