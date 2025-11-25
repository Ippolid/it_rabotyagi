
import React from 'react';
import { Button } from '../ui/button';
import { Menu, X, BookOpen, Users, MessageCircle, LayoutDashboard, Component } from 'lucide-react';

type View = 'landing' | 'courses' | 'course-detail' | 'questions' | 'mentors' | 'dashboard' | 'auth' | 'design-system';

interface LayoutProps {
  children: React.ReactNode;
  currentView: View;
  onNavigate: (view: View) => void;
}

export function Layout({ children, currentView, onNavigate }: LayoutProps) {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = React.useState(false);

  const navItems = [
    { label: 'Курсы', value: 'courses', icon: BookOpen },
    { label: 'Менторы', value: 'mentors', icon: Users },
    { label: 'Вопросы', value: 'questions', icon: MessageCircle },
    { label: 'Дашборд', value: 'dashboard', icon: LayoutDashboard },
    { label: 'Дизайн-система', value: 'design-system', icon: Component },
  ];

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 font-sans">
      {/* Sticky Navigation */}
      <nav className="sticky top-0 z-50 w-full border-b border-gray-200/50 bg-white/80 backdrop-blur-md supports-[backdrop-filter]:bg-white/60">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          {/* Logo */}
          <div 
            className="flex items-center gap-2 cursor-pointer" 
            onClick={() => onNavigate('landing')}
          >
            <div className="h-8 w-8 rounded-lg bg-gray-900 flex items-center justify-center text-white font-bold text-lg">
              R
            </div>
            <span className="text-lg font-semibold tracking-tight text-gray-900">
              IT-RABOTYAGI
            </span>
          </div>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-8">
            {navItems.map((item) => (
              <button
                key={item.value}
                onClick={() => onNavigate(item.value as View)}
                className={`text-sm font-medium transition-colors hover:text-gray-900 ${
                  currentView === item.value ? 'text-gray-900' : 'text-gray-500'
                }`}
              >
                {item.label}
              </button>
            ))}
          </div>

          {/* Desktop Actions */}
          <div className="hidden md:flex items-center gap-4">
            <Button variant="ghost" onClick={() => onNavigate('auth')}>
              Войти
            </Button>
            <Button onClick={() => onNavigate('auth')}>
              Начать
            </Button>
          </div>

          {/* Mobile Menu Button */}
          <button
            className="md:hidden p-2 text-gray-600 hover:bg-gray-100 rounded-md"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>

        {/* Mobile Menu */}
        {isMobileMenuOpen && (
          <div className="md:hidden border-t border-gray-200 bg-white px-4 py-4 shadow-lg">
            <div className="flex flex-col gap-4">
              {navItems.map((item) => (
                <button
                  key={item.value}
                  onClick={() => {
                    onNavigate(item.value as View);
                    setIsMobileMenuOpen(false);
                  }}
                  className={`flex items-center gap-3 text-sm font-medium ${
                    currentView === item.value ? 'text-gray-900' : 'text-gray-500'
                  }`}
                >
                  <item.icon size={18} />
                  {item.label}
                </button>
              ))}
              <div className="border-t border-gray-100 pt-4 flex flex-col gap-3">
                <Button variant="outline" className="w-full" onClick={() => onNavigate('auth')}>
                  Войти
                </Button>
                <Button className="w-full" onClick={() => onNavigate('auth')}>
                  Начать
                </Button>
              </div>
            </div>
          </div>
        )}
      </nav>

      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-auto">
        <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
            <div>
              <h3 className="text-sm font-semibold text-gray-900">Платформа</h3>
              <ul className="mt-4 space-y-2">
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Курсы</button></li>
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Менторы</button></li>
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Тарифы</button></li>
              </ul>
            </div>
            <div>
              <h3 className="text-sm font-semibold text-gray-900">Компания</h3>
              <ul className="mt-4 space-y-2">
                <li><button className="text-sm text-gray-500 hover:text-gray-900">О нас</button></li>
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Вакансии</button></li>
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Блог</button></li>
              </ul>
            </div>
            <div>
              <h3 className="text-sm font-semibold text-gray-900">Ресурсы</h3>
              <ul className="mt-4 space-y-2">
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Сообщество</button></li>
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Центр помощи</button></li>
                <li><button className="text-sm text-gray-500 hover:text-gray-900">Условия</button></li>
              </ul>
            </div>
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="h-6 w-6 rounded bg-gray-900 flex items-center justify-center text-white text-xs font-bold">R</div>
                <span className="font-semibold text-gray-900">IT-RABOTYAGI</span>
              </div>
              <p className="text-sm text-gray-500">
                Помогаем новому поколению айтишников с практикой и менторством.
              </p>
            </div>
          </div>
          <div className="mt-12 border-t border-gray-100 pt-8 text-center text-sm text-gray-400">
            &copy; 2024 IT-RABOTYAGI. Все права защищены.
          </div>
        </div>
      </footer>
    </div>
  );
}
