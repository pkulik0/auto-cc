<script>
import LanguageDropdown from "./LanguageDropdown.svelte";
</script>

<nav class="navbar navbar-expand-lg bg-primary mb-2">
    <div class="container-fluid">
        <div class="navbar-brand">AutoCC</div>
        <div class="navbar-nav">
            <LanguageDropdown/>
        </div>
    </div>
</nav>

<div class="container-fluid">
    <slot/>
</div>