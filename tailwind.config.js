/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "templates/*.html",   // or *.html
        "static/**/*.js",        // JS using Tailwind classes
    ],
    theme: {
        extend: {},
    },
    plugins: [],
}
