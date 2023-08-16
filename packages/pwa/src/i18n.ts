import i18n from "i18next"
import { initReactI18next } from 'react-i18next'
import detector from 'i18next-browser-languagedetector'

import en from './locale/en'

i18n
    .use(detector)
    .use(initReactI18next)
    .init({
        resources: {
            'en': en,
            'en-GB': en,
            'en-IE': en,
            'en-US': en,
        },
        debug: true,
        fallbackLng: 'en',
    })