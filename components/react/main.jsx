import React from 'react';
import { createRoot } from 'react-dom/client';
import TrendingContent from './trendingContent';
import DailyAnalytics from './daily-analytics';

// Registry mapping component names to React components
const components = {
	TrendingContent,
	DailyAnalytics
};

// Function to mount React components in a given root element
export default function mountReactComponents(root = document) {
	root.querySelectorAll('[data-react-component]').forEach((el) => {
		const componentName = el.getAttribute('data-react-component');
		const Component = components[componentName];
		if (Component) {
			createRoot(el).render(<Component />);
		} else {
			console.error(`No component registered for ${componentName}`);
		}
	});
}

// Run on initial load
document.addEventListener("DOMContentLoaded", () => {
	mountReactComponents();

	// Listen for HTMX swaps and mount React inside new content
	document.body.addEventListener("htmx:afterSwap", (event) => {
		mountReactComponents(event.target);
	});
});

