// Проверка авторизации при загрузке страницы
document.addEventListener('DOMContentLoaded', async () => {
    // Проверяем URL параметры для сообщений
    const urlParams = new URLSearchParams(window.location.search);
    const authStatus = urlParams.get('auth');
    
    if (authStatus === 'success') {
        showMessage('Вы успешно авторизованы! 🎉', 'success');
        // Очищаем URL от параметров
        window.history.replaceState({}, document.title, window.location.pathname);
    } else if (authStatus === 'error') {
        showMessage('Ошибка авторизации. Попробуйте снова.', 'error');
        window.history.replaceState({}, document.title, window.location.pathname);
    }

    // Проверяем авторизован ли пользователь
    await checkAuth();
});

// Проверка авторизации
async function checkAuth() {
    try {
        console.log('Checking authentication...');
        const response = await fetch('/api/v1/me', {
            method: 'GET',
            credentials: 'include', // Важно для отправки cookies
            headers: {
                'Accept': 'application/json',
            }
        });

        console.log('Auth check response status:', response.status);

        if (response.ok) {
            const data = await response.json();
            console.log('Auth check data:', data);
            
            if (data.success && data.user) {
                showUserInfo(data.user);
            } else {
                console.log('No user in response');
                showAuthButtons();
            }
        } else {
            console.log('Not authenticated, status:', response.status);
            const errorText = await response.text();
            console.log('Error response:', errorText);
            showAuthButtons();
        }
    } catch (error) {
        console.error('Error checking auth:', error);
        showAuthButtons();
    }
}

// Показать информацию о пользователе
function showUserInfo(user) {
    // Скрыть кнопки авторизации
    document.getElementById('auth-buttons').style.display = 'none';
    
    // Показать информацию о пользователе
    const userInfoDiv = document.getElementById('user-info');
    userInfoDiv.style.display = 'block';

    // Установить данные пользователя
    document.getElementById('user-name').textContent = user.name || 'Пользователь';
    
    const emailElement = document.getElementById('user-email');
    if (user.email) {
        emailElement.textContent = user.email;
        emailElement.style.display = 'block';
    } else {
        emailElement.style.display = 'none';
    }

    // Установить аватар
    const avatarImg = document.getElementById('user-avatar-img');
    if (user.avatar_url && user.avatar_url !== '') {
        avatarImg.src = user.avatar_url;
        avatarImg.style.display = 'block';
        avatarImg.onerror = function() {
            // Если аватар не загрузился, показываем placeholder
            this.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect fill="%233b82f6" width="100" height="100"/%3E%3Ctext fill="white" font-size="48" font-weight="bold" x="50%25" y="50%25" text-anchor="middle" dy=".3em"%3E' + (user.name ? user.name[0].toUpperCase() : '?') + '%3C/text%3E%3C/svg%3E';
        };
    } else {
        // Показываем placeholder с первой буквой имени
        avatarImg.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect fill="%233b82f6" width="100" height="100"/%3E%3Ctext fill="white" font-size="48" font-weight="bold" x="50%25" y="50%25" text-anchor="middle" dy=".3em"%3E' + (user.name ? user.name[0].toUpperCase() : '?') + '%3C/text%3E%3C/svg%3E';
        avatarImg.style.display = 'block';
    }

    // Показать бейджи провайдеров
    if (user.github_id) {
        document.getElementById('github-badge').style.display = 'inline-block';
    }
    if (user.google_id) {
        document.getElementById('google-badge').style.display = 'inline-block';
    }
    if (user.telegram_id) {
        document.getElementById('telegram-badge').style.display = 'inline-block';
    }

    console.log('User info displayed:', user);
}

// Показать кнопки авторизации
function showAuthButtons() {
    document.getElementById('auth-buttons').style.display = 'block';
    document.getElementById('user-info').style.display = 'none';
}

// Показать сообщение
function showMessage(message, type = 'info') {
    const messageDiv = document.getElementById('auth-message');
    messageDiv.textContent = message;
    messageDiv.className = 'auth-message ' + type;
    messageDiv.style.display = 'block';

    // Автоматически скрыть через 5 секунд
    setTimeout(() => {
        messageDiv.style.display = 'none';
    }, 5000);
}

// Выход из системы
async function logout() {
    try {
        const response = await fetch('/api/v1/logout', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Accept': 'application/json',
            }
        });

        if (response.ok) {
            showMessage('Вы вышли из системы', 'success');
            // Обновить страницу через секунду
            setTimeout(() => {
                window.location.reload();
            }, 1000);
        } else {
            showMessage('Ошибка при выходе', 'error');
        }
    } catch (error) {
        console.error('Logout error:', error);
        showMessage('Ошибка при выходе', 'error');
    }
}

// Utility функция для получения cookie
function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
    return null;
}
