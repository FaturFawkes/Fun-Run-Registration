import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#FF6B35',
          light: '#FF8C61',
          dark: '#E55527',
        },
        secondary: {
          DEFAULT: '#004E89',
          light: '#0066B3',
          dark: '#003A67',
        },
        accent: {
          DEFAULT: '#F7B801',
          light: '#FFD34E',
          dark: '#C99200',
        },
      },
    },
  },
  plugins: [],
}

export default config
