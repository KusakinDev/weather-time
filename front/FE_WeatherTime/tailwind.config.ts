import type { Config } from "tailwindcss";

export default {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--background)",
        foreground: "var(--foreground)",

        //цвета фонов
        bgPage: '#000000', //'#cad2c5',
        bgComp: '#84a98c',
        bgElem: '#52796f',
        bgElemActive: '#354f52',
        bgHeaderFooter: '#2f3e46',

        //цвета шрифтов для фонов
        txPage: '#4b6166',
        txComp: '#2e4057',
        txElem: '#f4a261',
        txElemActive: '#ffcb77',
        txHeaderFooter: '#edf2f4',
      },
    },
  },
  plugins: [],
} satisfies Config;
