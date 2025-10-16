// Add animation to dots
const dots = document.querySelectorAll('.dot');
let currentDot = 0;

setInterval(() => {
    dots[currentDot].classList.remove('active');
    currentDot = (currentDot + 1) % dots.length;
    dots[currentDot].classList.add('active');
}, 2000);

// Add click handlers for auth buttons
document.querySelectorAll('.auth-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        if (!btn.href || btn.href === window.location.href + '#') {
            e.preventDefault();
            alert('Функция авторизации будет доступна после настройки OAuth провайдеров');
        }
    });
});

