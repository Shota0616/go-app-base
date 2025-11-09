import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import enTranslation from '../../locales/en.json';
import jaTranslation from '../../locales/ja.json';

const resources = {
    en: {
        translation: enTranslation,
    },
    ja: {
        translation: jaTranslation,
    },
};

const storedLang = localStorage.getItem('language');

i18n
    .use(initReactI18next)
    .init({
        resources,
        lng: storedLang || import.meta.env.VITE_APP_LANG || 'en', // Use stored language, then .env, then default to 'en'
        fallbackLng: 'en',
        interpolation: {
            escapeValue: false,
        },
    });

i18n.on('languageChanged', (lng) => {
    localStorage.setItem('language', lng);
});

export default i18n;