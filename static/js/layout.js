let lastScrollY = 0;

document.addEventListener('DOMContentLoaded', function() {
    // User Dropdown Logic
    const userMenuButton = document.getElementById('user-menu-button');
    const userDropdown = document.getElementById('dropdown-user');

    if (userMenuButton && userDropdown) {
        // User Dropdown Toggle
        userMenuButton.addEventListener('click', function(event) {
            event.stopPropagation();

            // Toggle visibility
            userDropdown.classList.toggle('hidden');

            // Position the dropdown precisely
            const buttonRect = userMenuButton.getBoundingClientRect();
            userDropdown.style.position = 'fixed';
            userDropdown.style.top = `${buttonRect.bottom + 10}px`;
            userDropdown.style.right = `${window.innerWidth - buttonRect.right}px`;
        });

        // Close user dropdown when clicking outside
        document.addEventListener('click', function() {
            userDropdown.classList.add('hidden');
        });

        // Prevent user dropdown from closing when clicking inside
        userDropdown.addEventListener('click', function(event) {
            event.stopPropagation();
        });
    }

    // Mobile Navigation Menu Logic
    const mobileMenuButton = document.querySelector('[data-collapse-toggle="navbar-user"]');
    const mobileMenu = document.getElementById('navbar-user');

    if (mobileMenuButton && mobileMenu) {
        mobileMenuButton.addEventListener('click', function() {
            // Toggle mobile menu visibility
            mobileMenu.classList.toggle('hidden');

            // Update aria-expanded attribute
            const isExpanded = this.getAttribute('aria-expanded') === 'true';
            this.setAttribute('aria-expanded', (!isExpanded).toString());
        });
    }
});

document.body.addEventListener('htmx:beforeSwap', function(evt) {
    // Check if there's an HX-Retarget header
    const retargetHeader = evt.detail.xhr.getResponseHeader("HX-Retarget");

    if (retargetHeader) {
        // Change the target of the swap
        evt.detail.target = document.querySelector(retargetHeader);
    }

    lastScrollY = window.scrollY;
});

document.body.addEventListener('htmx:afterSwap', function() {
    window.scrollTo({ top: lastScrollY, behavior: 'instant' });
});

function sendAdClick() {
    htmx.ajax('POST', '/api/increment-ads-clicks', {
        swap: 'none'
    });
}
