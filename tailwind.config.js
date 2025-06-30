/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./website/component/**/*.{html,js,templ,txt}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
};
