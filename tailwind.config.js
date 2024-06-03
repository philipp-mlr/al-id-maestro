/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./component/**/*.{html,js,templ,txt}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
};
