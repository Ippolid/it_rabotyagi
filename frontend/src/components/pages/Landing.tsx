
import React, { useEffect, useState } from 'react';
import { Button } from '../ui/button';
import { ArrowRight, Code, Users, Terminal, Star, TrendingUp } from 'lucide-react';
import { motion } from 'motion/react';
import { courses as fallbackCourses, mentors as fallbackMentors } from '../../lib/data';
import { Badge } from '../ui/badge';
import { listCourses, listMentors } from '../../lib/api';

export function Landing({ onNavigate }: { onNavigate: (view: any) => void }) {
  const [topCourses, setTopCourses] = useState(fallbackCourses.slice(0, 4));
  const [topMentors, setTopMentors] = useState(fallbackMentors.slice(0, 3));

  useEffect(() => {
    let cancelled = false;
    const load = async () => {
      try {
        const [coursesRes, mentorsRes] = await Promise.all([listCourses(), listMentors()]);
        if (cancelled) return;
        const mappedCourses = coursesRes.items.slice(0, 4).map((c, idx) => {
          const fallback = fallbackCourses[idx % fallbackCourses.length];
          return {
            ...fallback,
            id: String(c.id),
            title: c.title,
            description: c.description,
          };
        });
        setTopCourses(mappedCourses);

        const mappedMentors = mentorsRes.items.slice(0, 3).map((m, idx) => {
          const fb = fallbackMentors[idx % fallbackMentors.length];
          return {
            id: String(m.id),
            name: m.fullName || fb.name,
            role: m.title || fb.role,
            company: fb.company,
            specialization: m.skills || fb.specialization,
            bio: fb.bio,
            available: true,
            avatar: fb.avatar,
          };
        });
        setTopMentors(mappedMentors);
      } catch {
        // keep fallbacks
      }
    };
    load();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 pb-24 overflow-hidden font-sans selection:bg-gray-200">
      
      {/* Hero Section: Dynamic Bento Grid */}
      <section className="pt-12 md:pt-20 pb-16 px-4 md:px-6 max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-12 grid-rows-auto md:grid-rows-[minmax(300px,_auto)_minmax(200px,_auto)] gap-6">
          
          {/* Main Value Prop Tile (Largest) */}
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
            className="col-span-1 md:col-span-7 row-span-2 bg-white rounded-[32px] p-8 md:p-12 flex flex-col justify-between shadow-sm border border-gray-100 relative overflow-hidden group"
          >
            <div className="z-10 space-y-8">
              <Badge variant="secondary" className="bg-gray-100 text-gray-600 hover:bg-gray-200 rounded-full px-4 py-1 text-sm font-medium">
                Новое: курс по System Design
              </Badge>
              <h1 className="text-5xl md:text-7xl font-bold tracking-tight text-gray-900 leading-[0.95]">
                Учись так же <br/>
                как работаешь.
              </h1>
              <p className="text-xl text-gray-500 max-w-md leading-relaxed">
                Освой продакшен-инжиниринг с менторством от опытных инженеров. Без воды, только код.
              </p>
              
              <div className="pt-4">
                <Button 
                  size="lg" 
                  className="h-14 px-8 text-lg rounded-full bg-gray-900 text-white hover:bg-gray-800 transition-all duration-300 transform hover:scale-[1.02] hover:shadow-[0_0_20px_rgba(0,0,0,0.2)] active:scale-[0.98]"
                  onClick={() => onNavigate('courses')}
                >
                  Начать обучение <ArrowRight className="ml-2 w-5 h-5" />
                </Button>
              </div>
            </div>

            {/* Subtle Abstract Background Gradient */}
            <div className="absolute -bottom-[20%] -right-[10%] w-[80%] h-[80%] bg-gradient-to-tl from-gray-100/50 to-transparent rounded-full blur-3xl pointer-events-none opacity-0 group-hover:opacity-100 transition-opacity duration-1000" />
          </motion.div>

          {/* Secondary Tile: Code Snippet (Glass-morphic floating effect) */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.2, duration: 0.6 }}
            className="col-span-1 md:col-span-5 row-span-1 relative overflow-hidden flex flex-col justify-center rounded-[32px] bg-white border border-gray-100 shadow-[0_20px_60px_rgba(15,23,42,0.15)]"
          >
            <div className="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-40 pointer-events-none" />
            <div className="absolute inset-0 bg-gradient-to-br from-white via-white to-gray-50" />
            <motion.div
              className="absolute -inset-24 bg-gradient-to-r from-transparent via-gray-200/50 to-transparent"
              animate={{ x: ['-30%', '130%'] }}
              transition={{ duration: 8, repeat: Infinity, ease: 'easeInOut' }}
            />
            <div className="relative z-10 font-mono text-base text-gray-900 space-y-4 p-8">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-3 h-3 rounded-full bg-red-400 shadow-md shadow-red-300/60" />
                <div className="w-3 h-3 rounded-full bg-amber-300 shadow-md shadow-amber-200/60" />
                <div className="w-3 h-3 rounded-full bg-emerald-300 shadow-md shadow-emerald-200/60" />
                <span className="ml-auto text-xs uppercase tracking-[0.18em] text-gray-700 font-semibold">ЖИВОЙ СЕАНС</span>
              </div>
              <div className="space-y-3 leading-relaxed">
                <p>const mastery = await learn(concepts);</p>
                <p>if (mastery) {'{'}</p>
                <p className="pl-4">ship("Production");</p>
                <p className="pl-4">// выкатываем уверенно 🚀</p>
                <p>{'}'}</p>
              </div>
              <div className="flex justify-between pt-6 text-sm font-medium text-gray-800">
                <span>задержка: 12мс</span>
                <span className="flex items-center gap-2">
                  аптайм: 99.99%
                  <Terminal className="w-5 h-5 text-gray-400" />
                </span>
              </div>
            </div>
          </motion.div>

          {/* Tertiary Tile: Mentors (Floating avatars) */}
          <motion.div 
             initial={{ opacity: 0, y: 20 }}
             animate={{ opacity: 1, y: 0 }}
             transition={{ delay: 0.3, duration: 0.6 }}
             className="col-span-1 md:col-span-3 bg-white rounded-[32px] p-6 shadow-sm border border-gray-100 flex flex-col justify-center items-center text-center relative overflow-hidden"
             onClick={() => onNavigate('mentors')}
          >
             <div className="flex -space-x-4 mb-4">
                {topMentors.slice(0,3).map((m, i) => (
                   <div key={i} className="w-12 h-12 rounded-full border-2 border-white overflow-hidden shadow-lg z-10">
                      <img src={`https://images.unsplash.com/${m.avatar}?auto=format&fit=crop&w=100&q=80`} alt="" className="w-full h-full object-cover" />
                   </div>
                ))}
             </div>
            <p className="font-semibold text-gray-900">Экспертные менторы</p>
            <p className="text-sm text-gray-500">Из топовых IT-компаний</p>
          </motion.div>

           {/* Quaternary Tile: Progress/Stats */}
           <motion.div 
             initial={{ opacity: 0, y: 20 }}
             animate={{ opacity: 1, y: 0 }}
             transition={{ delay: 0.4, duration: 0.6 }}
             className="col-span-1 md:col-span-2 bg-blue-50 rounded-[32px] p-6 flex flex-col justify-center items-center text-center relative group cursor-pointer hover:bg-blue-100 transition-colors"
             onClick={() => onNavigate('dashboard')}
          >
             <div className="w-12 h-12 bg-white rounded-2xl flex items-center justify-center shadow-sm mb-3 group-hover:scale-110 transition-transform duration-300">
                <TrendingUp className="text-blue-600 w-6 h-6" />
             </div>
             <div className="text-2xl font-bold text-gray-900">120%</div>
             <p className="text-xs font-medium text-blue-600 uppercase tracking-wide">Рост эффективности</p>
          </motion.div>

        </div>
      </section>

      {/* Horizontal Scrollable Course Strip */}
      <section className="py-12 space-y-6">
         <div className="max-w-7xl mx-auto px-6 flex justify-between items-end">
            <div>
               <h2 className="text-2xl md:text-3xl font-bold text-gray-900 tracking-tight">Популярные курсы</h2>
               <p className="text-gray-500 mt-1">Программа на острие технологий.</p>
            </div>
            <Button variant="ghost" onClick={() => onNavigate('courses')} className="group text-gray-900">
               Все курсы <ArrowRight className="ml-2 w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </Button>
         </div>
         
         {/* Scroll Container */}
         <div className="w-full overflow-x-auto pb-12 pt-4 hide-scrollbar">
            <div className="flex gap-6 px-6 md:px-8 w-max">
               {topCourses.map((course, index) => (
                  <motion.div 
                     key={course.id}
                     initial={{ opacity: 0, x: 50 }}
                     whileInView={{ opacity: 1, x: 0 }}
                     viewport={{ once: true }}
                     transition={{ delay: index * 0.1 }}
                     className="w-[300px] md:w-[380px] flex-shrink-0 group"
                  >
                     <div 
                        className="bg-white rounded-3xl p-4 shadow-sm border border-gray-100 hover:shadow-xl hover:shadow-gray-200/50 transition-all duration-500 cursor-pointer h-full flex flex-col"
                        onClick={() => onNavigate('courses')}
                     >
                        <div className="relative aspect-[4/3] rounded-2xl overflow-hidden mb-4">
                           <img 
                              src={`https://images.unsplash.com/${course.image}?auto=format&fit=crop&w=600&q=80`}
                              alt={course.title}
                              className="w-full h-full object-cover transform group-hover:scale-105 transition-transform duration-700 ease-out"
                           />
                           <div className="absolute top-3 right-3">
                              <Badge className="bg-white/90 backdrop-blur text-gray-900 shadow-sm border-none font-medium">
                                 {course.difficulty === 'Beginner' ? 'Начальный' : course.difficulty === 'Intermediate' ? 'Средний' : 'Продвинутый'}
                              </Badge>
                           </div>
                           <div className="absolute inset-0 bg-black/0 group-hover:bg-black/5 transition-colors duration-300" />
                        </div>
                        
                        <div className="px-2 pb-2 flex-1 flex flex-col">
                           <h3 className="text-xl font-bold text-gray-900 mb-2 group-hover:text-blue-600 transition-colors">
                              {course.title}
                           </h3>
                           <p className="text-gray-500 text-sm line-clamp-2 mb-4 flex-1">
                              {course.description}
                           </p>
                           <div className="flex items-center justify-between text-sm text-gray-400 font-medium pt-4 border-t border-gray-50">
                              <div className="flex items-center gap-1">
                                <Code size={16} />
                                 <span>{course.modulesCount} модулей</span>
                              </div>
                              <div className="flex items-center gap-1">
                                 <Users size={16} />
                                 <span>1.2k студентов</span>
                              </div>
                           </div>
                        </div>
                     </div>
                  </motion.div>
               ))}
            </div>
         </div>
      </section>

      {/* Floating Features Section */}
      <section className="max-w-7xl mx-auto px-6 py-12">
         <div className="grid md:grid-cols-3 gap-8">
            {[
               {
                  title: "Интерактивные лаборатории",
                  desc: "Без настройки окружения: код прямо в браузере через облачный IDE.",
                  icon: Terminal
               },
               {
                  title: "Код-ревью",
                  desc: "Получайте разборы от сеньоров по вашим проектам и кейсам.",
                  icon: Star
               },
               {
                  title: "Пробные интервью",
                  desc: "Практикуйте system design и алгоритмы в формате реальных собесов.",
                  icon: Users
               }
            ].map((feature, i) => (
               <motion.div
                  key={i}
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.1 }}
                  className="p-8 rounded-3xl bg-white border border-gray-100 shadow-sm hover:shadow-2xl hover:shadow-gray-200/50 transition-all duration-500 group"
               >
                  <div className="w-12 h-12 rounded-2xl bg-gray-50 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:bg-gray-900 group-hover:text-white transition-all duration-300">
                     <feature.icon size={24} />
                  </div>
                  <h3 className="text-xl font-bold text-gray-900 mb-3">{feature.title}</h3>
                  <p className="text-gray-500 leading-relaxed">{feature.desc}</p>
               </motion.div>
            ))}
         </div>
      </section>

       {/* CTA Section */}
       <section className="max-w-4xl mx-auto px-6 pt-20 text-center">
         <motion.div 
            initial={{ opacity: 0, scale: 0.95 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            className="bg-gradient-to-b from-gray-900 to-gray-800 rounded-[40px] p-12 md:p-20 shadow-2xl shadow-gray-900/20 relative overflow-hidden"
         >
            {/* Decorative Background Elements */}
            <div className="absolute top-0 left-0 w-full h-full overflow-hidden z-0">
               <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-blue-500/20 rounded-full blur-[100px] transform translate-x-1/2 -translate-y-1/2"></div>
               <div className="absolute bottom-0 left-0 w-[500px] h-[500px] bg-purple-500/20 rounded-full blur-[100px] transform -translate-x-1/2 translate-y-1/2"></div>
            </div>

            <div className="relative z-10 space-y-8">
               <h2 className="text-4xl md:text-6xl font-bold text-white tracking-tight">
                  Будущий ты <br /> скажет спасибо.
               </h2>
               <p className="text-lg text-gray-300 max-w-xl mx-auto">
                  Присоединяйся к платформе, которая меняет подход к обучению инженеров. 
                  Начни бесплатный пробный период.
               </p>
               <Button 
                  size="lg" 
                  className="h-16 px-10 text-xl rounded-full bg-white text-gray-900 hover:bg-gray-100 shadow-[0_0_40px_rgba(255,255,255,0.3)] hover:scale-105 transition-all duration-300"
                  onClick={() => onNavigate('auth')}
               >
                  Начать сейчас
               </Button>
            </div>
         </motion.div>
       </section>
    </div>
  );
}
