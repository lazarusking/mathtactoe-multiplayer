/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#2bee79",
        secondary: "#a3e635",
        accent: "#f43f5e",
        "background-light": "#f6f8f7",
        "background-dark": "#0B1121",
        "surface-light": "#ffffff",
        "surface-dark": "#1E293B",
        "cell-light": "#f1f5f9",
        "cell-dark": "#334155"
      },
      fontFamily: {
        display: ["Space Grotesk", "sans-serif"],
        body: ["Plus Jakarta Sans", "sans-serif"]
      }
    }
  },
  plugins: [],
}