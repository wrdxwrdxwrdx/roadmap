import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'
import enTranslations from './locales/en.json'
import ruTranslations from './locales/ru.json'

// Определяем доступные языки
export const supportedLanguages = ['en', 'ru'] as const
export type SupportedLanguage = typeof supportedLanguages[number]

// Ресурсы переводов
const resources = {
  en: {
    translation: enTranslations,
  },
  ru: {
    translation: ruTranslations,
  },
}

// Конфигурация i18next
i18n
  // Определение языка браузера
  .use(LanguageDetector)
  // Передача экземпляра i18n в react-i18next
  .use(initReactI18next)
  // Инициализация i18next
  .init({
    resources,
    fallbackLng: 'ru', // Язык по умолчанию
    supportedLngs: supportedLanguages,
    
    // Настройки для определения языка
    detection: {
      // Порядок определения языка
      order: ['localStorage', 'navigator', 'htmlTag'],
      // Ключ для хранения языка в localStorage
      lookupLocalStorage: 'i18nextLng',
      // Кэшировать выбранный язык
      caches: ['localStorage'],
    },

    // Настройки интерполяции
    interpolation: {
      escapeValue: false, // React уже экранирует значения
    },

    // Настройки для разработки
    debug: false,

    // Настройки для загрузки переводов
    load: 'languageOnly', // Загружать только язык, без региона (en вместо en-US)
  })

export default i18n

