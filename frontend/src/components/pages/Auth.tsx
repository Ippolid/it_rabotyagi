
import React from 'react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../ui/card';
import { Github, Mail } from 'lucide-react';
import { login, register, saveTokens } from '../../lib/api';

export function Auth({ onLogin }: { onLogin: (name?: string) => void }) {
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const handleAuth = async (e: React.FormEvent, mode: 'login' | 'signup') => {
    e.preventDefault();
    setError(null);
    setIsLoading(true);
    const form = new FormData(e.target as HTMLFormElement);
    const email = String(form.get('email') || form.get('email-signup') || '');
    const password = String(form.get('password') || form.get('password-signup') || '');
    const nickname = mode === 'signup' ? String(form.get('nickname') || email.split('@')[0]) : '';

    try {
      if (mode === 'login') {
        const tokens = await login(email, password);
        saveTokens(tokens);
      } else {
        const tokens = await register(email, password, nickname);
        saveTokens(tokens);
      }
      onLogin(mode === 'login' ? email : nickname);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка авторизации');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-[80vh]">
      <Card className="w-full max-w-md shadow-xl border-gray-100">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold tracking-tight text-center">С возвращением</CardTitle>
          <CardDescription className="text-center">
            Введите email, чтобы войти в свой аккаунт
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="login" className="w-full">
            <TabsList className="grid w-full grid-cols-2 mb-8">
              <TabsTrigger value="login">Вход</TabsTrigger>
              <TabsTrigger value="signup">Регистрация</TabsTrigger>
            </TabsList>
            
            <TabsContent value="login">
              <form onSubmit={(e) => handleAuth(e, 'login')} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input id="email" name="email" type="email" placeholder="m@example.com" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">Пароль</Label>
                  <Input id="password" name="password" type="password" required />
                </div>
                <Button className="w-full" type="submit" disabled={isLoading}>
                  {isLoading ? "Входим..." : "Войти"}
                </Button>
              </form>
            </TabsContent>
            
            <TabsContent value="signup">
              <form onSubmit={(e) => handleAuth(e, 'signup')} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="nickname">Никнейм</Label>
                  <Input id="nickname" name="nickname" placeholder="rabotyaga" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="email-signup">Email</Label>
                  <Input id="email-signup" name="email-signup" type="email" placeholder="m@example.com" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password-signup">Пароль</Label>
                  <Input id="password-signup" name="password-signup" type="password" required />
                </div>
                <Button className="w-full" type="submit" disabled={isLoading}>
                  {isLoading ? "Создаем аккаунт..." : "Создать аккаунт"}
                </Button>
              </form>
            </TabsContent>
          </Tabs>

          {error && <p className="text-sm text-red-600 mt-2">{error}</p>}

          <div className="relative my-8">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-white px-2 text-gray-500">Или продолжить через</span>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <Button variant="outline" onClick={() => setError('Соцлогин пока недоступен')}>
              <Github className="mr-2 h-4 w-4" />
              Github
            </Button>
            <Button variant="outline" onClick={() => setError('Соцлогин пока недоступен')}>
              <Mail className="mr-2 h-4 w-4" />
              Google
            </Button>
          </div>
        </CardContent>
        <CardFooter>
          <p className="text-center text-sm text-gray-500 w-full">
            Продолжая, вы соглашаетесь с условиями сервиса и политикой конфиденциальности.
          </p>
        </CardFooter>
      </Card>
    </div>
  );
}
