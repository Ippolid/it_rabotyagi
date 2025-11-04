import type { Config } from 'tailwindcss'

export default {
  content: [
    './index.html',
    './src/**/*.{ts,tsx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#eef9ff',
          100: '#d8f0ff',
          200: '#b9e5ff',
          300: '#89d6ff',
          400: '#4bc0ff',
          500: '#1aa6ff',
          600: '#0685e0',
          700: '#0569b3',
          800: '#094f85',
          900: '#0c406a'
        }
      }
    }
  },
  plugins: [],
} satisfies Config


