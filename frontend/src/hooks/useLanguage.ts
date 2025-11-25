import { useTranslation } from 'react-i18next'
import { supportedLanguages, SupportedLanguage } from '../i18n/config'

/**
 * Хук для работы с языками интерфейса
 * Предоставляет функции для переключения языка и получения текущего языка
 */
export function useLanguage() {
  const { i18n } = useTranslation()

  const currentLanguage = i18n.language as SupportedLanguage

  const changeLanguage = (lang: SupportedLanguage) => {
    i18n.changeLanguage(lang)
  }

  const toggleLanguage = () => {
    const newLang = currentLanguage === 'ru' ? 'en' : 'ru'
    changeLanguage(newLang)
  }

  const isRussian = currentLanguage === 'ru'
  const isEnglish = currentLanguage === 'en'

  return {
    currentLanguage,
    changeLanguage,
    toggleLanguage,
    isRussian,
    isEnglish,
    supportedLanguages,
  }
}

