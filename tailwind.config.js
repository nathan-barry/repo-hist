module.exports = {
  content: ["views/**/*.html"],
  theme: {
    extend: {
      backgroundColor: {
        'github-dark': '#0d1117', // GitHub dark background
        'github-dark-secondary': '#161b22', // Secondary background color
        // ... other colors as needed
      },
      textColor: {
        'github-dark': '#c9d1d9', // GitHub dark text color
        // ... other colors as needed
      },
      borderColor: {
        'github-dark': '#30363d', // GitHub dark border color
        // ... other colors as needed
      },
      // Add other GitHub-like dark theme extensions if needed
    },
  },
  plugins: [],
}
