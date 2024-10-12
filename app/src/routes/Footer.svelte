<script lang="ts">
	import { Footer, FooterCopyright, FooterLinkGroup, FooterLink } from 'flowbite-svelte';
</script>

<Footer>
	<div class="sm:flex sm:items-center sm:justify-between">
		<FooterCopyright href="/" by="Onysh" year={2024} />
		<FooterLinkGroup
			ulClass="flex flex-wrap items-center mt-3 text-sm text-gray-500 dark:text-gray-400 sm:mt-0"
		>
			<FooterLink href="/">Privacy Policy</FooterLink>
		</FooterLinkGroup>
	</div>
</Footer>
