/** @type {import('tailwindcss').Config} */
export default {
	content: [
		'./components/**/*.{html,templ,go,txt}',  // All your Templ files + embedded HTML
		'./static/js/**/*.js',                // If any JS files use Tailwind classes
		'./static/react/**/*.{js,jsx,ts,tsx}', // Just in case you drop React in
	],
	darkMode: 'class',
	theme: {
		extend: {},
	},
	plugins: [],
}

