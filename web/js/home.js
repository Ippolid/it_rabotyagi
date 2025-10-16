// Check if user is authenticated
document.addEventListener('DOMContentLoaded', async () => {
    await checkAuthStatus();
});

async function checkAuthStatus() {
    try {
        const response = await fetch('/api/v1/me', {
            method: 'GET',
            credentials: 'include',
            headers: {
                'Accept': 'application/json',
            }
        });

        const authBtn = document.getElementById('auth-btn');
        if (response.ok) {
            const data = await response.json();
            if (data.success && data.user) {
                // User is authenticated
                authBtn.textContent = data.user.name || 'Профиль';
                authBtn.href = 'index.html'; // or dashboard.html
                authBtn.classList.add('authenticated');
            }
        }
    } catch (error) {
        console.error('Error checking auth:', error);
    }
}

// Theme toggle
const themeToggle = document.querySelector('.theme-toggle');
if (themeToggle) {
    themeToggle.addEventListener('click', () => {
        document.body.style.background = document.body.style.background === 'rgb(30, 41, 59)' ? '#f8f9fa' : '#1e293b';
        themeToggle.textContent = themeToggle.textContent === '🌙' ? '☀️' : '🌙';
    });
}

// Smooth scroll
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({ behavior: 'smooth' });
        }
    });
});
