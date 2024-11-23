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

i18n
    .use(initReactI18next)
    .init({
        resources,
        lng: import.meta.env.VITE_APP_LANG || 'en', // .envファイルのAPP_LANGを使用
        fallbackLng: 'en',
        interpolation: {
            escapeValue: false,
        },
    });

export default i18n;