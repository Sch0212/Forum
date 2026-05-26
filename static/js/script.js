document.addEventListener("DOMContentLoaded", () => {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach((entry) => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                observer.unobserve(entry.target);
            }
        });
    }, { threshold: 0.1 });

    const fadeElements = document.querySelectorAll('.fade-in');
    
    fadeElements.forEach((el, index) => {
        setTimeout(() => {
            observer.observe(el);
        }, index * 120);
    });
});