
import React, { useEffect, useState } from 'react';
import { Layout } from './components/layout/Layout';
import { Landing } from './components/pages/Landing';
import { Courses } from './components/pages/Courses';
import { CourseDetail } from './components/pages/CourseDetail';
import { Questions } from './components/pages/Questions';
import { Mentors } from './components/pages/Mentors';
import { Auth } from './components/pages/Auth';
import { Dashboard } from './components/pages/Dashboard';
import { DesignSystem } from './components/pages/DesignSystem';
import { clearTokens, getProfile, getStoredTokens } from './lib/api';

type View = 'landing' | 'courses' | 'course-detail' | 'questions' | 'mentors' | 'dashboard' | 'auth' | 'design-system';

export default function App() {
  const [currentView, setCurrentView] = useState<View>('landing');
  const [selectedCourseId, setSelectedCourseId] = useState<string | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userName, setUserName] = useState<string | null>(null);

  useEffect(() => {
    const tokens = getStoredTokens();
    if (tokens?.accessToken) {
      setIsAuthenticated(true);
      loadProfile();
    }
    syncFromLocation();
    const onPop = () => syncFromLocation();
    window.addEventListener('popstate', onPop);
    return () => window.removeEventListener('popstate', onPop);
  }, []);

  const syncFromLocation = () => {
    const path = window.location.pathname;
    if (path.startsWith('/courses/')) {
      const id = path.split('/')[2];
      setSelectedCourseId(id || null);
      setCurrentView('course-detail');
      return;
    }
    const map: Record<string, View> = {
      '/': 'landing',
      '/courses': 'courses',
      '/questions': 'questions',
      '/mentors': 'mentors',
      '/dashboard': 'dashboard',
      '/auth': 'auth',
      '/design-system': 'design-system',
    };
    setCurrentView(map[path] ?? 'landing');
    setSelectedCourseId(null);
  };

  const pushPath = (view: View, courseId?: string | null) => {
    const path =
      view === 'landing'
        ? '/'
        : view === 'courses'
          ? '/courses'
          : view === 'course-detail' && courseId
            ? `/courses/${courseId}`
            : view === 'questions'
              ? '/questions'
              : view === 'mentors'
                ? '/mentors'
                : view === 'dashboard'
                  ? '/dashboard'
                  : view === 'auth'
                    ? '/auth'
                    : view === 'design-system'
                      ? '/design-system'
                      : '/';
    window.history.pushState({ view, courseId }, '', path);
  };

  const loadProfile = async () => {
    try {
      const profile = await getProfile();
      setUserName(profile.fullName || profile.email);
    } catch {
      // ignore errors, stay in guest mode
    }
  };

  const handleNavigate = (view: View) => {
    if (view === 'dashboard' && !isAuthenticated) {
      setCurrentView('auth');
      pushPath('auth');
      return;
    }
    setCurrentView(view);
    setSelectedCourseId(null);
    pushPath(view);
    window.scrollTo(0, 0);
  };

  const handleCourseSelect = (id: string) => {
    setSelectedCourseId(id);
    setCurrentView('course-detail');
    pushPath('course-detail', id);
    window.scrollTo(0, 0);
  };

  const handleLogin = (name?: string) => {
    setIsAuthenticated(true);
    if (name) setUserName(name);
    loadProfile();
    setCurrentView('dashboard');
    pushPath('dashboard');
  };

  const handleLogout = () => {
    clearTokens();
    setIsAuthenticated(false);
    setUserName(null);
    setCurrentView('landing');
    pushPath('landing');
  };

  return (
    <Layout currentView={currentView} onNavigate={handleNavigate} userName={userName} onLogout={handleLogout}>
      {currentView === 'landing' && <Landing onNavigate={handleNavigate} />}
      
      {currentView === 'courses' && (
        <Courses onSelectCourse={handleCourseSelect} />
      )}
      
      {currentView === 'course-detail' && selectedCourseId && (
        <CourseDetail 
          courseId={selectedCourseId} 
          onBack={() => setCurrentView('courses')} 
        />
      )}
      
      {currentView === 'questions' && <Questions />}
      
      {currentView === 'mentors' && <Mentors />}
      
      {currentView === 'auth' && <Auth onLogin={handleLogin} />}
      
      {currentView === 'dashboard' && <Dashboard onNavigate={handleNavigate} userName={userName} />}
      
      {currentView === 'design-system' && <DesignSystem />}
    </Layout>
  );
}
