<script lang="ts">
	import { onMount } from 'svelte';
	import '../app.css';
	import Footer from './Footer.svelte';
	import Navbar from './Navbar.svelte';
</script>

<div class="flex flex-col min-h-screen">
	<Navbar />

	<div class="flex-grow container mx-auto my-4">
		<div class="mx-8"><slot /></div>
	</div>

	<div class="container mx-auto mb-4"><Footer /></div>
</div>
