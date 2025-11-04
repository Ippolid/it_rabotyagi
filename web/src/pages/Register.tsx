import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { register } from '../lib/api'

export default function RegisterPage() {
  const nav = useNavigate()
  const [email, setEmail] = useState('')
  const [nickname, setNickname] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await register(email, nickname || email.split('@')[0], password)
      nav('/')
    } catch (err: any) {
      setError(err?.message || 'Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <div className="relative isolate">
        {/* background image */}
        <div className="pointer-events-none absolute inset-0 -z-10">
          <img src="/auth-bg.jpg" alt="bg" className="h-full w-full object-cover opacity-90" />
        </div>
        <div className="container flex min-h-screen items-center justify-center py-16">
          <div className="w-full max-w-md rounded-2xl bg-white/80 p-8 shadow-xl backdrop-blur">
            <div className="mb-6 flex items-center gap-3">
              <img src="/logo.jpg" className="h-10 w-10 rounded-full" />
              <div>
                <div className="text-xl font-bold">Регистрация</div>
                <div className="text-sm text-slate-600">Создайте аккаунт и начните обучение</div>
              </div>
            </div>
            <form className="space-y-4" onSubmit={onSubmit}>
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Email</label>
                <input required type="email" value={email} onChange={e=>setEmail(e.target.value)} className="w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 placeholder-slate-400 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500" placeholder="you@example.com" />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Никнейм</label>
                <input required value={nickname} onChange={e=>setNickname(e.target.value)} className="w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 placeholder-slate-400 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500" placeholder="username" />
              </div>
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">Пароль</label>
                <input required type="password" value={password} onChange={e=>setPassword(e.target.value)} className="w-full rounded-md border border-slate-300 px-3 py-2 text-slate-900 placeholder-slate-400 focus:border-brand-500 focus:outline-none focus:ring-1 focus:ring-brand-500" placeholder="••••••••" />
              </div>
              {error && <div className="text-sm text-red-600">{error}</div>}
              <button disabled={loading} className="btn-primary w-full justify-center">{loading ? 'Создание...' : 'Зарегистрироваться'}</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}


