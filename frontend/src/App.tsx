
import React, { useEffect, useState } from 'react';
import { Layout } from './components/layout/Layout';
import { Landing } from './components/pages/Landing';
import { Courses } from './components/pages/Courses';
import { CourseDetail } from './components/pages/CourseDetail';
import { ModuleDetail } from './components/pages/ModuleDetail';
import { Questions } from './components/pages/Questions';
import { QuestionDetail } from './components/pages/QuestionDetail';
import { Mentors } from './components/pages/Mentors';
import { Auth } from './components/pages/Auth';
import { Dashboard } from './components/pages/Dashboard';
import { DesignSystem } from './components/pages/DesignSystem';
import { InterviewTrainer } from './components/pages/InterviewTrainer';
import { clearTokens, getProfile, getStoredTokens } from './lib/api';

type View = 'landing' | 'courses' | 'course-detail' | 'module-detail' | 'questions' | 'question-detail' | 'mentors' | 'dashboard' | 'auth' | 'design-system' | 'interview-trainer';

export default function App() {
  const [currentView, setCurrentView] = useState<View>('landing');
  const [selectedCourseId, setSelectedCourseId] = useState<string | null>(null);
  const [selectedModuleId, setSelectedModuleId] = useState<string | null>(null);
  const [selectedQuestionId, setSelectedQuestionId] = useState<string | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userName, setUserName] = useState<string | null>(null);
  const [userAvatar, setUserAvatar] = useState<string | null>(null);

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

    // Check /courses/:courseId/modules/:moduleId/questions/:questionId
    const moduleQuestionMatch = path.match(/^\/courses\/(\d+)\/modules\/(\d+)\/questions\/(\d+)$/);
    if (moduleQuestionMatch) {
      setSelectedCourseId(moduleQuestionMatch[1]);
      setSelectedModuleId(moduleQuestionMatch[2]);
      setSelectedQuestionId(moduleQuestionMatch[3]);
      setCurrentView('question-detail');
      return;
    }

    // Check /courses/:courseId/modules/:moduleId
    const moduleMatch = path.match(/^\/courses\/(\d+)\/modules\/(\d+)$/);
    if (moduleMatch) {
      setSelectedCourseId(moduleMatch[1]);
      setSelectedModuleId(moduleMatch[2]);
      setCurrentView('module-detail');
      return;
    }

    // Check /questions/:id BEFORE /questions
    if (path.startsWith('/questions/')) {
      const id = path.split('/')[2];
      setSelectedQuestionId(id || null);
      setSelectedCourseId(null);
      setSelectedModuleId(null);
      setCurrentView('question-detail');
      return;
    }

    // Check /courses/:id
    if (path.startsWith('/courses/')) {
      const id = path.split('/')[2];
      setSelectedCourseId(id || null);
      setSelectedModuleId(null);
      setCurrentView('course-detail');
      return;
    }

    const map: Record<string, View> = {
      '/': 'landing',
      '/courses': 'courses',
      '/questions': 'questions',
      '/mentors': 'mentors',
      '/profile': 'dashboard',
      '/dashboard': 'dashboard', // Keep for backwards compatibility
      '/auth': 'auth',
      '/design-system': 'design-system',
      '/interview-trainer': 'interview-trainer',
    };
    setCurrentView(map[path] ?? 'landing');
    setSelectedCourseId(null);
    setSelectedModuleId(null);
    setSelectedQuestionId(null);
  };

  const pushPath = (view: View, courseId?: string | null, moduleId?: string | null, questionId?: string | null) => {
    const path =
      view === 'landing'
        ? '/'
        : view === 'courses'
          ? '/courses'
          : view === 'course-detail' && courseId
            ? `/courses/${courseId}`
            : view === 'module-detail' && courseId && moduleId
              ? `/courses/${courseId}/modules/${moduleId}`
              : view === 'questions'
                ? '/questions'
                : view === 'question-detail' && questionId && courseId && moduleId
                  ? `/courses/${courseId}/modules/${moduleId}/questions/${questionId}`
                  : view === 'question-detail' && questionId
                    ? `/questions/${questionId}`
                    : view === 'mentors'
                      ? '/mentors'
                      : view === 'dashboard'
                        ? '/profile'
                        : view === 'auth'
                          ? '/auth'
                          : view === 'design-system'
                            ? '/design-system'
                            : view === 'interview-trainer'
                              ? '/interview-trainer'
                              : '/';
    window.history.pushState({ view, courseId, moduleId, questionId }, '', path);
  };

  const loadProfile = async () => {
    try {
      const profile = await getProfile();
      // Приоритет: fullName > email, убираем никнейм из отображения
      setUserName(profile.fullName || profile.email);
      setUserAvatar(profile.avatarUrl || null);
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

  const handleQuestionSelect = (id: string) => {
    setSelectedQuestionId(id);
    setSelectedCourseId(null);
    setSelectedModuleId(null);
    setCurrentView('question-detail');
    pushPath('question-detail', null, null, id);
    window.scrollTo(0, 0);
  };

  const handleModuleSelect = (courseId: string, moduleId: string) => {
    setSelectedCourseId(courseId);
    setSelectedModuleId(moduleId);
    setCurrentView('module-detail');
    pushPath('module-detail', courseId, moduleId);
    window.scrollTo(0, 0);
  };

  const handleModuleQuestionSelect = (questionId: string, courseId: string, moduleId: string) => {
    setSelectedQuestionId(questionId);
    setSelectedCourseId(courseId);
    setSelectedModuleId(moduleId);
    setCurrentView('question-detail');
    pushPath('question-detail', courseId, moduleId, questionId);
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
    setUserAvatar(null);
    setCurrentView('landing');
    pushPath('landing');
  };

  return (
    <Layout currentView={currentView} onNavigate={handleNavigate} userName={userName} userAvatar={userAvatar} onLogout={handleLogout}>
      {currentView === 'landing' && <Landing onNavigate={handleNavigate} />}
      
      {currentView === 'courses' && (
        <Courses onSelectCourse={handleCourseSelect} />
      )}
      
      {currentView === 'course-detail' && selectedCourseId && (
        <CourseDetail
          courseId={selectedCourseId}
          onBack={() => {
            setCurrentView('courses');
            pushPath('courses');
          }}
          onNavigateToModule={handleModuleSelect}
        />
      )}

      {currentView === 'module-detail' && selectedCourseId && selectedModuleId && (
        <ModuleDetail
          courseId={selectedCourseId}
          moduleId={selectedModuleId}
          onBack={() => {
            setCurrentView('course-detail');
            pushPath('course-detail', selectedCourseId);
          }}
          onNavigateToQuestion={handleModuleQuestionSelect}
        />
      )}

      {currentView === 'questions' && <Questions onSelectQuestion={handleQuestionSelect} />}

      {currentView === 'question-detail' && selectedQuestionId && (
        <QuestionDetail
          questionId={selectedQuestionId}
          courseId={selectedCourseId || undefined}
          moduleId={selectedModuleId || undefined}
          onBack={() => {
            if (selectedCourseId && selectedModuleId) {
              // Back to module if coming from module context
              setCurrentView('module-detail');
              pushPath('module-detail', selectedCourseId, selectedModuleId);
            } else {
              // Back to questions list if standalone question
              setCurrentView('questions');
              pushPath('questions');
            }
          }}
          onSelectQuestion={handleQuestionSelect}
        />
      )}
      
      {currentView === 'mentors' && <Mentors />}
      
      {currentView === 'auth' && <Auth onLogin={handleLogin} />}
      
      {currentView === 'dashboard' && <Dashboard onNavigate={handleNavigate} userName={userName} />}
      
      {currentView === 'design-system' && <DesignSystem />}
      
      {currentView === 'interview-trainer' && <InterviewTrainer />}
    </Layout>
  );
}
