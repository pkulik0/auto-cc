<script lang="ts">
	import '../app.css';
	import Footer from './Footer.svelte';
	import Navbar from './Navbar.svelte';
</script>

<div class="flex min-h-screen flex-col">
	<Navbar />

	<div class="container mx-auto my-4 flex-grow">
		<div class="mx-4"><slot /></div>
	</div>

	<div class="container mx-auto mb-4"><Footer /></div>
</div>
