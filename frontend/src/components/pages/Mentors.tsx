
import React, { useEffect, useState } from 'react';
import { mentors as fallbackMentors } from '../../lib/data';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Card, CardContent, CardFooter } from '../ui/card';
import { Calendar, Briefcase } from 'lucide-react';
import { motion } from 'motion/react';
import { listMentors } from '../../lib/api';

export function Mentors() {
  const [items, setItems] = useState(fallbackMentors);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    const load = async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await listMentors();
        if (cancelled) return;
        const mapped = res.items.map((m, idx) => {
          const fb = fallbackMentors[idx % fallbackMentors.length];
          return {
            id: String(m.id),
            name: m.fullName || fb.name,
            role: m.title || fb.role,
            company: fb.company,
            avatar: fb.avatar,
            specialization: m.skills || fb.specialization,
            bio: fb.bio,
            available: true,
          };
        });
        setItems(mapped as any);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Не удалось загрузить менторов');
        setItems(fallbackMentors);
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="space-y-8">
      <div className="text-center max-w-2xl mx-auto mb-12">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Найди своего ментора</h1>
        <p className="text-gray-500">
          Свяжись с опытными инженерами, получай ревью, карьерные советы и пробные интервью.
        </p>
      </div>

      {loading && <p className="text-sm text-gray-500 text-center">Загружаем менторов...</p>}
      {error && <p className="text-sm text-red-600 text-center">{error}</p>}

      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        {items.map((mentor, i) => (
          <motion.div
             key={mentor.id}
             initial={{ opacity: 0, scale: 0.95 }}
             animate={{ opacity: 1, scale: 1 }}
             transition={{ delay: i * 0.1 }}
          >
            <Card className="h-full flex flex-col overflow-hidden hover:shadow-xl transition-shadow border-none shadow-md">
               <div className="h-32 bg-gradient-to-r from-gray-800 to-gray-900 relative">
                  <div className="absolute -bottom-12 left-6">
                     <img 
                        src={`https://images.unsplash.com/${mentor.avatar}?auto=format&fit=crop&w=300&q=80`} 
                        alt={mentor.name}
                        className="w-24 h-24 rounded-2xl border-4 border-white object-cover shadow-sm bg-white"
                     />
                  </div>
               </div>
               <CardContent className="pt-16 pb-4 flex-1">
                  <div className="flex justify-between items-start mb-2">
                     <div>
                        <h3 className="font-bold text-xl text-gray-900">{mentor.name}</h3>
                     <div className="flex items-center gap-1 text-sm text-gray-500">
                        <Briefcase size={14} />
                        <span>{mentor.role} в {mentor.company}</span>
                     </div>
                    </div>
                    {mentor.available ? (
                        <Badge className="bg-green-100 text-green-700 hover:bg-green-200 border-none">Доступен</Badge>
                    ) : (
                        <Badge variant="secondary">Нет слотов</Badge>
                    )}
                 </div>
                  
                  <p className="text-gray-600 text-sm mt-4 line-clamp-3 leading-relaxed">
                     {mentor.bio}
                  </p>

                  <div className="mt-6 flex flex-wrap gap-2">
                     {mentor.specialization.map(spec => (
                        <Badge key={spec} variant="secondary" className="font-normal bg-gray-50 border-gray-200 text-gray-600">
                           {spec}
                        </Badge>
                     ))}
                  </div>
               </CardContent>
               <CardFooter className="border-t border-gray-100 pt-4 bg-gray-50/50">
                  <Button className="w-full gap-2" disabled={!mentor.available}>
                     <Calendar size={16} />
                     Забронировать сессию
                  </Button>
               </CardFooter>
            </Card>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
