export const prerender = true

import { register, init, getLocaleFromNavigator } from 'svelte-i18n';

register('en', () => import('$lib/locales/en.json'));

init({
	fallbackLocale: 'en',
	initialLocale: getLocaleFromNavigator(),
})