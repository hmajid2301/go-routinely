/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/templates/**/*.{html,templ,go}",
    "./internal/transport/http/**/*.go",
  ],
  theme: {
    extend: {
      colors: {
        // Catppuccin Latte (light theme)
        latte: {
          base: '#eff1f5',
          mantle: '#e6e9ef',
          surface: '#ccd0da',
          text: '#4c4f69',
          overlay: '#9ca0b0',
          pink: '#ea76cb',
          mauve: '#8839ef',
          peach: '#fe640b',
          red: '#e78284',
          yellow: '#e5c890',
          green: '#a6d189',
          teal: '#81c8be',
        },
        // Catppuccin Mocha (dark theme)
        mocha: {
          base: '#1e1e2e',
          mantle: '#181825',
          surface: '#313244',
          text: '#cdd6f4',
          overlay: '#6c7086',
          pink: '#f5c2e7',
          mauve: '#cba6f7',
          peach: '#fab387',
          red: '#e78284',
          yellow: '#e5c890',
          green: '#a6d189',
          teal: '#81c8be',
        },
      },
      borderRadius: {
        'xl': '12px',
      },
    },
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: [
      {
        light: {
          "primary": "#ea76cb",
          "secondary": "#8839ef",
          "accent": "#fe640b",
          "neutral": "#9ca0b0",
          "base-100": "#eff1f5",
          "base-200": "#e6e9ef",
          "base-300": "#ccd0da",
          "info": "#81c8be",
          "success": "#a6d189",
          "warning": "#e5c890",
          "error": "#e78284",
        },
        dark: {
          "primary": "#f5c2e7",
          "secondary": "#cba6f7",
          "accent": "#fab387",
          "neutral": "#6c7086",
          "base-100": "#1e1e2e",
          "base-200": "#181825",
          "base-300": "#313244",
          "info": "#81c8be",
          "success": "#a6d189",
          "warning": "#e5c890",
          "error": "#e78284",
        },
      },
    ],
  },
}
