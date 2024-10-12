<script lang="ts">
	import { page } from '$app/stores';
	import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, Button } from 'flowbite-svelte';
	
	import { _ } from 'svelte-i18n';

	$: activeUrl = $page.url.pathname;
</script>

<Navbar>
	<NavBrand href="/">
		<img src="/images/autocc.svg" class="me-1 mt-0.5 h-16 sm:h-12" alt="AutoCC Logo" />
		<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">{$_("app.name")}</span>
	</NavBrand>
	<div class="flex md:order-2">
		<Button size="sm">{$_("nav.sign_in")}</Button>
		<NavHamburger />
	</div>
	<NavUl {activeUrl} class="order-1">
		<NavLi href="/">{$_("nav.home")}</NavLi>
		<NavLi href="/credentials">{$_("nav.credentials")}</NavLi>
		<NavLi href="/settings">{$_("nav.settings")}</NavLi>
	</NavUl>
</Navbar>
