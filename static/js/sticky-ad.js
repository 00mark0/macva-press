document.addEventListener('DOMContentLoaded', function() {
    const stickyAdContainer = document.getElementById('sticky-ad-container');
    const lgBreakpoint = 1024; // Match the lg breakpoint from Tailwind (adjust if needed)
    let initialOffsetTop = null;

    function setupStickyBehavior() {
        // Only set up sticky behavior if screen is wide enough
        if (window.innerWidth >= lgBreakpoint && stickyAdContainer) {
            // Cache the initial offset if not already done
            if (initialOffsetTop === null) {
                initialOffsetTop = stickyAdContainer.getBoundingClientRect().top + window.scrollY;
            }
            const scrollLimit = initialOffsetTop + 1000; // Adjust this value as needed
            const scrollPosition = window.scrollY;

            // Apply sticky only within the defined scroll range
            if (scrollPosition >= initialOffsetTop && scrollPosition <= scrollLimit) {
                stickyAdContainer.classList.add('is-sticky');
            } else {
                stickyAdContainer.classList.remove('is-sticky');
            }
        } else if (stickyAdContainer) {
            // Remove sticky class on smaller screens
            stickyAdContainer.classList.remove('is-sticky');
        }
    }

    // Reset cached offset on window resize
    window.addEventListener('resize', function() {
        initialOffsetTop = null;
        setupStickyBehavior();
    });

    window.addEventListener('scroll', setupStickyBehavior);

    // Run once on page load
    setupStickyBehavior();
});

