export type Theme = "dark" | "light";

export type ThemeStore = {
  current: Theme;
  switchTheme: (theme: Theme) => void;

  init: () => void;
};

export function getSystemTheme() {
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light";
}

export function getStoredThemeInLocaleStorage(): Theme | null {
  const theme = localStorage.getItem("color-theme");

  if (theme === null) {
    return null;
  }

  if (theme !== "dark" && theme !== "light") {
    return null;
  }

  return theme;
}

export function storeThemeInLocaleStorage(theme: Theme) {
  localStorage.setItem("color-theme", theme);
}

export function toggleHtmlDarkClass(theme: Theme) {
  document.documentElement.classList.toggle("dark", theme === "dark");
}
