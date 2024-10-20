<script lang="ts">
	import { page } from '$app/stores';
	import { login, logout, userStore } from '$lib/auth';
	import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Button } from 'flowbite-svelte';

	import { _ } from 'svelte-i18n';

	$: activeUrl = $page.url.pathname;
</script>

<Navbar>
	<NavBrand href="/">
		<img src="/images/autocc.svg" class="me-1 mt-0.5 h-16 w-16 sm:h-12 sm:w-12" alt="AutoCC Logo" />
		<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
			>{$_('app.name')}</span
		>
	</NavBrand>
	<div class="flex md:order-2">
		{#if $userStore}
			<Button on:click={() => logout()} size="sm" outline>{$_('nav.sign_out')}</Button>
		{:else}
			<Button on:click={() => login()} size="sm">{$_('nav.sign_in')}</Button>
		{/if}
		<NavHamburger />
	</div>
	<NavUl {activeUrl} class="order-1">
		<NavLi href="/videos">{$_('nav.videos')}</NavLi>
		<NavLi href="/credentials">{$_('nav.credentials')}</NavLi>
		<NavLi href="/settings">{$_('nav.settings')}</NavLi>
	</NavUl>
</Navbar>
