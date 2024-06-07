<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { fade } from 'svelte/transition';
  import { ToastType, type Toast } from '$lib/types/Toast';

  const dispatch = createEventDispatcher();

  export let toast: Toast = {
    id: 0,
    type: ToastType.Error,
    message: 'An error occurred',
  };

  let type = {
    [ToastType.Error]: 'error',
    [ToastType.Info]: 'info',
    [ToastType.Success]: 'success',
    [ToastType.Warning]: 'warning',
  }[toast.type];
</script>

<article class={type} role="alert" transition:fade>
  <p>
    <slot />
  </p>

  <button on:click={() => dispatch('dismiss')} >
    <i class="fas fa-times"></i>
  </button>
</article>

<style lang="scss">
  article {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 1rem;
    border-radius: var(--border-radius);
    box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.25);
    gap: 1rem;

    p {
      margin: 0;
    }

    button {
      background: none;
      border: none;
      color: inherit;
      cursor: pointer;
      font-size: 1.5rem;
    }

    &.error {
      background-color: var(--red);
      color: var(--white);
    }

    &.info {
      background-color: var(--blue);
      color: var(--white);
    }

    &.success {
      background-color: var(--green);
      color: var(--white);
    }

    &.warning {
      background-color: var(--orange);
      color: var(--white);
    }
  }

</style>