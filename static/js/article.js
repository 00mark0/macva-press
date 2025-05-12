const getThemePreference = () => {
    if (typeof localStorage !== 'undefined' && localStorage.getItem('theme')) {
        return localStorage.getItem('theme');
    }
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
};

// Apply the current theme
const applyTheme = (theme) => {
    if (theme === 'dark') {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
    localStorage.setItem('theme', theme);
};

applyTheme(getThemePreference());

const createThemeToggle = document.getElementById('create-theme-toggle');

createThemeToggle.addEventListener("click", () => {
    const isDark = document.documentElement.classList.contains('dark');
    applyTheme(isDark ? 'light' : 'dark');
})

// Initialize TinyMCE
document.addEventListener('DOMContentLoaded', function() {
    initTinyMCE();
});

document.body.addEventListener('htmx:beforeSwap', function(evt) {
    // Check if there's an HX-Retarget header
    const retargetHeader = evt.detail.xhr.getResponseHeader("HX-Retarget");

    if (retargetHeader) {
        // Change the target of the swap
        evt.detail.target = document.querySelector(retargetHeader);
    }
});

window.addEventListener("beforeunload", function() {
    fetch("/api/cookie", { method: "DELETE", keepalive: true });
});

function initTinyMCE() {
    tinymce.init({
        selector: '#editor',
        plugins: 'anchor autolink charmap codesample emoticons image link lists media searchreplace table visualblocks wordcount',
        toolbar: 'undo redo | blocks | bold italic underline strikethrough | link image media table | align lineheight | numlist bullist indent outdent | emoticons charmap | removeformat',
        height: 800,
        skin: 'oxide-dark',
        content_css: 'dark',
        setup: function(editor) {
            editor.on('change', function(e) {
                editor.save();
            });
        }
    });
}
