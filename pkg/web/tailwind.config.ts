import type { Config } from "tailwindcss";
import fluid, { extract, fontSize, screens } from "fluid-tailwind";

const config: Config = {
  content: {
    files: [
      "./app/**/*.{js,ts,jsx,tsx,mdx}",
    ],
    extract,
  },
  darkMode: "class",
  theme: {
    screens,
    fontSize,
    extend: {
      colors: {},
      fontFamily: {
        heading: ["var(--font-heading)"],
        sans: ["var(--font-sans)"],
        mono: ["var(--font-mono)"],
      },
    },
  },
  corePlugins: {
    preflight: false,
  },
  plugins: [fluid],
};

export default config;
