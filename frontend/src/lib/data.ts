
export type Difficulty = 'Beginner' | 'Intermediate' | 'Advanced';

export interface Course {
  id: string;
  title: string;
  description: string;
  modulesCount: number;
  difficulty: Difficulty;
  image: string;
  tags: string[];
  progress?: number;
}

export interface Mentor {
  id: string;
  name: string;
  role: string;
  company: string;
  avatar: string;
  specialization: string[];
  bio: string;
  available: boolean;
}

export interface Question {
  id: string;
  title: string;
  preview: string;
  tags: string[];
  difficulty: Difficulty;
  replies: number;
  views: number;
}

export interface Module {
  id: string;
  title: string;
  duration: string;
  completed: boolean;
}

export const courses: Course[] = [
  {
    id: '1',
    title: 'Fullstack React',
    description: 'Создавайте продакшен-приложения на React, Node.js и PostgreSQL.',
    modulesCount: 12,
    difficulty: 'Intermediate',
    image: 'photo-1633356122544-f134324a6cee',
    tags: ['React', 'Node.js', 'PostgreSQL'],
    progress: 45,
  },
  {
    id: '2',
    title: 'Python для Data Science',
    description: 'Освойте анализ данных, визуализацию и ML на Python.',
    modulesCount: 15,
    difficulty: 'Beginner',
    image: 'photo-1526379095098-d400fd0bf935',
    tags: ['Python', 'Data Science', 'Pandas'],
    progress: 0,
  },
  {
    id: '3',
    title: 'Продвинутый System Design',
    description: 'Проектируйте масштабируемые распределённые системы под высокую нагрузку.',
    modulesCount: 8,
    difficulty: 'Advanced',
    image: 'photo-1519389950473-47ba0277781c',
    tags: ['System Design', 'Architecture', 'Scalability'],
    progress: 0,
  },
  {
    id: '4',
    title: 'Основы UI/UX',
    description: 'Создавайте удобные интерфейсы с Figma и базовыми принципами дизайна.',
    modulesCount: 10,
    difficulty: 'Beginner',
    image: 'photo-1561070791-2526d30994b5',
    tags: ['Design', 'Figma', 'UI/UX'],
    progress: 0,
  }
];

export const mentors: Mentor[] = [
  {
    id: '1',
    name: 'Sarah Chen',
    role: 'Staff Engineer',
    company: 'TechFlow',
    avatar: 'photo-1472099645785-5658abf4ff4e',
    specialization: ['React', 'TypeScript', 'Performance'],
    bio: 'Делаю доступные и быстрые веб-приложения.',
    available: true,
  },
  {
    id: '2',
    name: 'Alexei Petrov',
    role: 'Senior Backend Engineer',
    company: 'DataSystems',
    avatar: 'photo-1494790108377-be9c29b29330',
    specialization: ['Go', 'Kubernetes', 'Distributed Systems'],
    bio: 'Эксперт в облачной архитектуре и масштабируемом бэкенде.',
    available: false,
  },
  {
    id: '3',
    name: 'Emma Williams',
    role: 'Product Design Lead',
    company: 'CreativeStudio',
    avatar: 'photo-1500648767791-00dcc994a43e',
    specialization: ['UI/UX', 'Design Systems', 'Prototyping'],
    bio: 'Помогаю инженерам понимать дизайн и лучше работать в команде.',
    available: true,
  }
];

export const questions: Question[] = [
  {
    id: '1',
    title: 'Как управлять состоянием в сложных формах React?',
    preview: 'Трудности с multi-step формой и валидацией...',
    tags: ['React', 'Forms', 'State Management'],
    difficulty: 'Intermediate',
    replies: 12,
    views: 340,
  },
  {
    id: '2',
    title: 'Лучшие практики обработки ошибок в REST API',
    preview: 'Возвращать 200 с ошибкой в теле или правильные коды статуса?',
    tags: ['API', 'Backend', 'Best Practices'],
    difficulty: 'Beginner',
    replies: 25,
    views: 890,
  },
  {
    id: '3',
    title: 'Оптимизация запросов PostgreSQL для больших таблиц',
    preview: 'Запрос выполняется 5 секунд на таблице с 1М строк...',
    tags: ['Database', 'SQL', 'Performance'],
    difficulty: 'Advanced',
    replies: 5,
    views: 120,
  }
];
