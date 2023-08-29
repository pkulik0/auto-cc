<script lang="ts">
    import {toasts} from "$lib/toast";
</script>

<div class="toast-container position-fixed bottom-0 end-0 p-3">
    {#each $toasts as toast}
        <div id="toast{toast.id}" class="toast" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-header">
                <img src="" class="rounded me-2" alt="">
                <strong class="me-auto">{toast.title}</strong>
                <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
            </div>
            <div class="toast-body">
                {toast.body}
            </div>
        </div>
    {/each}
</div>