
import React, { useState } from 'react';
import { Layout } from './components/layout/Layout';
import { Landing } from './components/pages/Landing';
import { Courses } from './components/pages/Courses';
import { CourseDetail } from './components/pages/CourseDetail';
import { Questions } from './components/pages/Questions';
import { Mentors } from './components/pages/Mentors';
import { Auth } from './components/pages/Auth';
import { Dashboard } from './components/pages/Dashboard';
import { DesignSystem } from './components/pages/DesignSystem';

type View = 'landing' | 'courses' | 'course-detail' | 'questions' | 'mentors' | 'dashboard' | 'auth' | 'design-system';

export default function App() {
  const [currentView, setCurrentView] = useState<View>('landing');
  const [selectedCourseId, setSelectedCourseId] = useState<string | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const handleNavigate = (view: View) => {
    if (view === 'dashboard' && !isAuthenticated) {
      setCurrentView('auth');
      return;
    }
    setCurrentView(view);
    window.scrollTo(0, 0);
  };

  const handleCourseSelect = (id: string) => {
    setSelectedCourseId(id);
    setCurrentView('course-detail');
    window.scrollTo(0, 0);
  };

  const handleLogin = () => {
    setIsAuthenticated(true);
    setCurrentView('dashboard');
  };

  return (
    <Layout currentView={currentView} onNavigate={handleNavigate}>
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
      
      {currentView === 'dashboard' && <Dashboard onNavigate={handleNavigate} />}
      
      {currentView === 'design-system' && <DesignSystem />}
    </Layout>
  );
}
