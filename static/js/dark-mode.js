// /static/js/dark-mode.js

document.addEventListener('DOMContentLoaded', function() {
    // Check for saved theme preference or use system preference
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

    // Initialize theme
    applyTheme(getThemePreference());

    // Create and append toggle button to navbar
    const createToggleButton = () => {
        const dropdownUser = document.querySelector("#dropdown-user");
        if (!dropdownUser) return;

        const toggleButton = document.createElement('button');
        toggleButton.className = 'p-2 text-gray-500 rounded-lg hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600 mr-3';
        toggleButton.setAttribute('id', 'theme-toggle');
        toggleButton.setAttribute('type', 'button');
        toggleButton.innerHTML = `
      <svg id="theme-toggle-dark-icon" class="w-5 h-5 hidden" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
        <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z"></path>
      </svg>
      <svg id="theme-toggle-light-icon" class="w-5 h-5 hidden" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
        <path d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" fill-rule="evenodd" clip-rule="evenodd"></path>
      </svg>
    `;

        dropdownUser.appendChild(toggleButton);

        // Update icon visibility
        const updateIcons = () => {
            const isDark = document.documentElement.classList.contains('dark');
            document.getElementById('theme-toggle-dark-icon').classList.toggle('hidden', isDark);
            document.getElementById('theme-toggle-light-icon').classList.toggle('hidden', !isDark);
        };

        // Add click handler to toggle
        toggleButton.addEventListener('click', () => {
            const isDark = document.documentElement.classList.contains('dark');
            applyTheme(isDark ? 'light' : 'dark');
            updateIcons();
        });

        // Initialize icons
        updateIcons();
    };

    createToggleButton();
});
