// static/js/admin.js
let lastScrollY = 0;

document.addEventListener('DOMContentLoaded', function() {
    // Profile dropdown toggle
    const userMenuButton = document.querySelector('[data-dropdown-toggle="dropdown-user"]');
    const userMenu = document.getElementById('dropdown-user');


    if (userMenuButton && userMenu) {
        // Position the dropdown properly
        userMenu.style.position = 'absolute';
        userMenu.style.right = '0';
        userMenu.style.top = '100%';

        userMenuButton.addEventListener('click', function(e) {
            e.stopPropagation();
            userMenu.classList.toggle('hidden');
        });

        // Close dropdown when clicking outside
        document.addEventListener('click', function(event) {
            if (!userMenuButton.contains(event.target) && !userMenu.contains(event.target)) {
                userMenu.classList.add('hidden');
            }
        });
    }

    // Mobile sidebar toggle
    const sidebarButton = document.querySelector('[data-drawer-toggle="logo-sidebar"]');
    const sidebar = document.getElementById('logo-sidebar');

    if (sidebarButton && sidebar) {
        sidebarButton.addEventListener('click', function() {
            sidebar.classList.toggle('-translate-x-full');
        });
    }

    // Get all sidebar menu items
    const menuItems = [
        document.getElementById('pregled'),
        document.getElementById('kategorije'),
        document.getElementById('artikli'),
        document.getElementById('korisnici'),
        document.getElementById('reklame')
    ];

    // Add click event to each menu item that exists
    menuItems.forEach(item => {
        if (item) {
            item.addEventListener('click', function() {
                // Only close the sidebar if we're on mobile (check window width or check if sidebar is visible)
                if (window.innerWidth < 640) { // sm breakpoint in Tailwind
                    sidebar.classList.add('-translate-x-full');
                }
            });
        }
    });

    const userMenuItems = [
        document.getElementById('user-menu-item-overview'),
        document.getElementById('user-menu-item-settings'),
        document.getElementById('user-menu-item-logout')
    ]

    userMenuItems.forEach(item => {
        item.addEventListener('click', function() {
            userMenu.classList.add('hidden');
        });
    });
});

document.body.addEventListener('htmx:beforeSwap', function() {
    lastScrollY = window.scrollY;
});

document.body.addEventListener('htmx:afterSwap', function() {
    window.scrollTo({ top: lastScrollY, behavior: 'instant' });
});
