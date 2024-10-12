export const prerender = true;
export const ssr = false;

import { addMessages, init, getLocaleFromNavigator } from 'svelte-i18n';
import en from '$lib/locales/en.json';

addMessages('en', en);
init({
	fallbackLocale: 'en',
	initialLocale: getLocaleFromNavigator()
});
