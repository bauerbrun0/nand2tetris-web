import Alpine from "alpinejs";
import "flowbite";

import type { Theme, ThemeStore } from "../utils/theme.ts";
import {
  getStoredThemeInLocaleStorage,
  getSystemTheme,
  storeThemeInLocaleStorage,
  toggleHtmlDarkClass,
} from "../utils/theme.ts";
import { showToast } from "../utils/toast.ts";

export function init() {
  // @ts-expect-error: showToast is not part of the Window type
  window.showToast = showToast;
  // @ts-expect-error: Alpine is not part of the Window type
  window.Alpine = Alpine;
  Alpine.store("profileDropDown", false);
  Alpine.store("theme", {
    current: "dark",

    switchTheme(theme: Theme) {
      this.current = theme;
      storeThemeInLocaleStorage(theme);
      toggleHtmlDarkClass(theme);
    },

    init() {
      const storedTheme = getStoredThemeInLocaleStorage();
      if (storedTheme !== null) {
        this.current = storedTheme;
        return;
      }

      this.current = getSystemTheme();
    },
  } satisfies ThemeStore);

  Alpine.start();
}
