
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '../ui/card';
import { Input } from '../ui/input';
import { Progress } from '../ui/progress';
import { Separator } from '../ui/separator';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs';

export function DesignSystem() {
  return (
    <div className="space-y-16">
      <div className="space-y-4">
        <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">Дизайн-система</h1>
        <p className="text-lg text-gray-500 max-w-2xl">
          Минималистичная система в духе Apple: ясность, много воздуха, мягкие взаимодействия.
        </p>
      </div>

      <Separator />

      {/* Typography */}
      <section className="space-y-8">
        <h2 className="text-2xl font-semibold text-gray-900">Типографика</h2>
        <div className="space-y-4 border rounded-xl p-8 bg-white shadow-sm">
          <div className="space-y-1">
            <p className="text-sm text-gray-400">Дисплей</p>
            <h1 className="text-5xl font-bold tracking-tight text-gray-900">Заголовок 1</h1>
          </div>
          <div className="space-y-1">
            <p className="text-sm text-gray-400">Title 1</p>
            <h2 className="text-4xl font-bold tracking-tight text-gray-900">Заголовок 2</h2>
          </div>
          <div className="space-y-1">
            <p className="text-sm text-gray-400">Title 2</p>
            <h3 className="text-2xl font-semibold tracking-tight text-gray-900">Заголовок 3</h3>
          </div>
          <div className="space-y-1">
            <p className="text-sm text-gray-400">Основной текст</p>
            <p className="text-base text-gray-600 leading-relaxed">
              Базовый текст: читаемость, контраст, комфортный интерлиньяж.
              Чем проще, тем изящнее.
            </p>
          </div>
          <div className="space-y-1">
            <p className="text-sm text-gray-400">Малый / подпись</p>
            <p className="text-sm text-gray-500">
              Маленький текст для подписей и второстепенной информации.
            </p>
          </div>
        </div>
      </section>

      {/* Colors */}
      <section className="space-y-8">
        <h2 className="text-2xl font-semibold text-gray-900">Цветовая палитра</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="space-y-2">
            <div className="h-24 w-full rounded-xl bg-gray-900 shadow-sm"></div>
            <div className="flex justify-between px-1">
              <span className="font-medium text-sm">Основной</span>
              <span className="text-gray-500 text-sm">Gray-900</span>
            </div>
          </div>
          <div className="space-y-2">
            <div className="h-24 w-full rounded-xl bg-gray-100 shadow-sm"></div>
            <div className="flex justify-between px-1">
              <span className="font-medium text-sm">Поверхность</span>
              <span className="text-gray-500 text-sm">Gray-100</span>
            </div>
          </div>
          <div className="space-y-2">
            <div className="h-24 w-full rounded-xl bg-blue-600 shadow-sm"></div>
            <div className="flex justify-between px-1">
              <span className="font-medium text-sm">Акцент</span>
              <span className="text-gray-500 text-sm">Blue-600</span>
            </div>
          </div>
          <div className="space-y-2">
            <div className="h-24 w-full rounded-xl bg-white border border-gray-200 shadow-sm"></div>
            <div className="flex justify-between px-1">
              <span className="font-medium text-sm">Карточка</span>
              <span className="text-gray-500 text-sm">Белый</span>
            </div>
          </div>
        </div>
      </section>

      {/* Buttons */}
      <section className="space-y-8">
        <h2 className="text-2xl font-semibold text-gray-900">Кнопки</h2>
        <div className="flex flex-wrap gap-4 p-8 border rounded-xl bg-white shadow-sm items-center">
          <Button>Основная</Button>
          <Button variant="secondary">Вторичная</Button>
          <Button variant="outline">Обводка</Button>
          <Button variant="ghost">Прозрачная</Button>
          <Button variant="destructive">Опасное</Button>
          <Button disabled>Отключено</Button>
        </div>
      </section>

      {/* Inputs & Controls */}
      <section className="space-y-8">
        <h2 className="text-2xl font-semibold text-gray-900">Инпуты и контролы</h2>
        <div className="grid md:grid-cols-2 gap-8 p-8 border rounded-xl bg-white shadow-sm">
          <div className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium leading-none">Email</label>
              <Input placeholder="name@example.com" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium leading-none">Пароль</label>
              <Input type="password" placeholder="••••••••" />
            </div>
          </div>
          <div className="space-y-4">
             <div className="space-y-2">
              <label className="text-sm font-medium leading-none">Вкладки</label>
              <Tabs defaultValue="account" className="w-full">
                <TabsList className="grid w-full grid-cols-2">
                  <TabsTrigger value="account">Аккаунт</TabsTrigger>
                  <TabsTrigger value="password">Пароль</TabsTrigger>
                </TabsList>
              </Tabs>
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium leading-none">Бейджи / чипы</label>
              <div className="flex gap-2">
                <Badge>Стандарт</Badge>
                <Badge variant="secondary">Вторичный</Badge>
                <Badge variant="outline">Контур</Badge>
                <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-200 border-none">Своя метка</Badge>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Cards & Progress */}
      <section className="space-y-8">
        <h2 className="text-2xl font-semibold text-gray-900">Карточки и компоненты</h2>
        <div className="grid md:grid-cols-2 gap-6">
          <Card className="rounded-xl shadow-sm hover:shadow-md transition-shadow duration-200 border-gray-100">
            <CardHeader>
              <CardTitle>Прогресс по курсам</CardTitle>
              <CardDescription>Продолжайте в том же духе!</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="font-medium">React: основы</span>
                  <span className="text-gray-500">75%</span>
                </div>
                <Progress value={75} className="h-2" />
              </div>
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="font-medium">Системный дизайн</span>
                  <span className="text-gray-500">30%</span>
                </div>
                <Progress value={30} className="h-2" />
              </div>
            </CardContent>
            <CardFooter>
              <Button className="w-full" variant="outline">Открыть дашборд</Button>
            </CardFooter>
          </Card>
          
          <Card className="rounded-xl shadow-sm hover:shadow-md transition-shadow duration-200 border-gray-100">
            <CardHeader>
              <div className="flex items-center gap-4">
                <div className="h-12 w-12 rounded-full bg-gray-100 flex items-center justify-center">
                    <span className="font-bold text-gray-500">AC</span>
                </div>
                <div>
                    <CardTitle className="text-lg">Alex Chen</CardTitle>
                    <CardDescription>Старший фронтенд-инженер</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-600">
                Специалист по React, TypeScript и доступным веб-приложениям.
                Доступен для менторских сессий.
              </p>
              <div className="flex gap-2 mt-4 flex-wrap">
                <Badge variant="secondary" className="rounded-md">React</Badge>
                <Badge variant="secondary" className="rounded-md">TypeScript</Badge>
              </div>
            </CardContent>
            <CardFooter>
              <Button className="w-full">Забронировать сессию</Button>
            </CardFooter>
          </Card>
        </div>
      </section>
    </div>
  );
}
