<script lang="ts">
    import IdentitiesList from "./IdentitiesList.svelte";
    import AddIdentityForm from "./AddIdentityForm.svelte";
</script>

<div class="row mt-4">
    <div class="col-md-8 order-md-1 order-2">
        <IdentitiesList/>
    </div>
    <div class="col-md-4 order-md-2 order-1 mb-4">
        <AddIdentityForm/>
    </div>
</div>