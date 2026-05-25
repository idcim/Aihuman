import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./lib/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        ink: "#151923",
        muted: "#667085",
        line: "#d9dee8",
        canvas: "#f7f8fb",
        panel: "#ffffff",
        teal: {
          50: "#ecfdf9",
          600: "#0f766e",
          700: "#115e59",
        },
        amber: {
          50: "#fffbeb",
          600: "#b45309",
        },
      },
      boxShadow: {
        panel: "0 1px 2px rgba(16, 24, 40, 0.06)",
      },
    },
  },
  plugins: [],
};

export default config;

