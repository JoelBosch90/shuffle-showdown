<script lang="ts">
  import Toast from './Toast.svelte';

  import { dismissToast, toasts } from "$lib/store/toasts";
</script>

{#if $toasts}
  <section>
    {#each $toasts as toast (toast.id)}
      <Toast {toast} on:dismiss={() => dismissToast(toast.id)}>
        {toast.message}
      </Toast>
    {/each}
  </section>
{/if}

<style lang="scss">
  section {
    position: fixed;
    top: 3rem;
    right: 0;
    left: 0;
  }
</style>