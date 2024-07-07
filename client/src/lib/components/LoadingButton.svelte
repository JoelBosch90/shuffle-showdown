<script lang="ts">
	type ButtonTypes = 'button' | 'submit' | 'reset';

	export let isDisabled = false;
	export let isLoading = true;
	export let onClick = () => {};
	export let type: ButtonTypes = 'button';
</script>

<div class={`button-wrapper ${isLoading ? 'loading' : ''}`}>
	<button {type} class="filled" disabled={isLoading || isDisabled} on:click={onClick}>
		<slot />
	</button>
</div>

<style lang="scss">
	.button-wrapper {
		border-radius: var(--border-radius);
		box-shadow: 0 0 1em var(--white) inset;
		padding: 3px;

		&.loading {
			background: var(--rainbow);
			background-size: var(--animation-background-flow-size);
			animation: animation-background-flow 10s linear infinite reverse;

			button {
				border: 1px solid var(--white);
			}
		}
	}
</style>
