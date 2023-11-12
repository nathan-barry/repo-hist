/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["views/**/*.html"],
  theme: {
    extend: {
      backgroundColor: {
        'gruvbox-bg': '#282828',
        'gruvbox-red': '#cc241d',
        'gruvbox-green': '#98971a',
        'gruvbox-yellow': '#d79921',
        'gruvbox-blue': '#458588',
        'gruvbox-purple': '#b16286',
        'gruvbox-aqua': '#689d6a',
        'gruvbox-gray': '#a89984',
      },
      textColor: {
        'gruvbox-fg': '#ebdbb2',
        'gruvbox-red': '#cc241d',
        'gruvbox-green': '#98971a',
        'gruvbox-yellow': '#d79921',
        'gruvbox-blue': '#458588',
        'gruvbox-purple': '#b16286',
        'gruvbox-aqua': '#689d6a',
        'gruvbox-gray': '#a89984',
      },
      borderColor: {
        'gruvbox-fg': '#ebdbb2',
        'gruvbox-red': '#cc241d',
        'gruvbox-green': '#98971a',
        'gruvbox-yellow': '#d79921',
        'gruvbox-blue': '#458588',
        'gruvbox-purple': '#b16286',
        'gruvbox-aqua': '#689d6a',
        'gruvbox-gray': '#a89984',
      },
    },
  },
  plugins: [],
}
