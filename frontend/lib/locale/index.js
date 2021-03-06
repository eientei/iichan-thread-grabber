//@flow
const NextI18Next = require('next-i18next/dist/commonjs');

module.exports = new NextI18Next({
    defaultLanguage: 'en',
    otherLanguages: ['ru'],
    fallbackLng: ['en', 'ru'],
    detection: {
        lookupCookie: 'locale',
    },
});
