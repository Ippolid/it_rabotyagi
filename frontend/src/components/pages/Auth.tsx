
import React from 'react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../ui/card';
import { Github, Mail } from 'lucide-react';

export function Auth({ onLogin }: { onLogin: () => void }) {
  const [isLoading, setIsLoading] = React.useState(false);

  const handleAuth = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      onLogin();
    }, 1000);
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
              <form onSubmit={handleAuth} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="email">Email</Label>
                  <Input id="email" type="email" placeholder="m@example.com" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">Пароль</Label>
                  <Input id="password" type="password" required />
                </div>
                <Button className="w-full" type="submit" disabled={isLoading}>
                  {isLoading ? "Входим..." : "Войти"}
                </Button>
              </form>
            </TabsContent>
            
            <TabsContent value="signup">
              <form onSubmit={handleAuth} className="space-y-4">
                 <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="first-name">Имя</Label>
                      <Input id="first-name" required />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="last-name">Фамилия</Label>
                      <Input id="last-name" required />
                    </div>
                 </div>
                <div className="space-y-2">
                  <Label htmlFor="email-signup">Email</Label>
                  <Input id="email-signup" type="email" placeholder="m@example.com" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password-signup">Пароль</Label>
                  <Input id="password-signup" type="password" required />
                </div>
                <Button className="w-full" type="submit" disabled={isLoading}>
                  {isLoading ? "Создаем аккаунт..." : "Создать аккаунт"}
                </Button>
              </form>
            </TabsContent>
          </Tabs>

          <div className="relative my-8">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-white px-2 text-gray-500">Или продолжить через</span>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <Button variant="outline" onClick={handleAuth}>
              <Github className="mr-2 h-4 w-4" />
              Github
            </Button>
            <Button variant="outline" onClick={handleAuth}>
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
