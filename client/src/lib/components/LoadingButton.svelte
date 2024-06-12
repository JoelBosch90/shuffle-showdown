<script lang="ts">

  type ButtonTypes = 'button' | 'submit' | 'reset';

  export let isDisabled = false;
  export let isLoading = true;
  export let onClick = () => {};
  export let type: ButtonTypes = 'button';
</script>

<div class={`button-wrapper ${isLoading ? 'loading' : ''}`}>
  <button type={type} class="filled" disabled={isLoading || isDisabled} on:click={onClick}>
    <slot />
  </button>
</div>

<style lang="scss">
  .button-wrapper {
    border-radius: var(--border-radius);
    box-shadow: 0 0 1em var(--white) inset;
    padding: 3px;

    &.loading {
      background: linear-gradient(90deg, var(--red), var(--purple), var(--purple), var(--blue), var(--blue), var(--green), var(--green), var(--yellow), var(--yellow), var(--orange), var(--orange), var(--red));
      background-size: 400%;
      animation: flow 10s linear infinite;

      button {
        border: 1px solid var(--white);
      }
    }
  }

  @keyframes flow {
    0% {
      background-position: 0 0;
    }
    100% {
      background-position: 400% 0;
    }
  }
</style>