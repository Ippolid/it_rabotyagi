import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getMe } from '../lib/api'

export default function AuthBar() {
  const nav = useNavigate()
  const [user, setUser] = useState<any>(null)

  useEffect(() => {
    getMe()
      .then(setUser)
      .catch(() => setUser(null))
  }, [])

  function handleLogout() {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('expires_in')
    setUser(null)
    nav('/')
  }

  if (user) {
    return (
      <div className="flex items-center gap-3 text-sm">
        <span className="text-slate-600 dark:text-slate-300">{user.fullName || user.email}</span>
        <span className="rounded-md bg-slate-100 px-2 py-1 text-slate-700 dark:bg-slate-800 dark:text-slate-200">{user.role}</span>
        <button onClick={handleLogout} className="btn-secondary text-sm">Выйти</button>
      </div>
    )
  }

  return (
    <Link className="btn-primary" to="/auth">
      Войти / Регистрация
    </Link>
  )
}


