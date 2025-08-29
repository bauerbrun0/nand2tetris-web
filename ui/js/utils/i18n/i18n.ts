import i18next from "i18next";
import en from "../../translations/en.json";

i18next.init({
  lng: "en",
  resources: {
    en: { translation: en },
  },
});

export const t = i18next.t;
export default i18next;
