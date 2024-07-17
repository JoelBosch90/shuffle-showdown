<script>
	import '@fortawesome/fontawesome-free/css/all.min.css';
	import Toasts from '$lib/components/Toasts.svelte';
	import AudioProgressBorder from '$lib/components/AudioProgressBorder.svelte';
	import '../styles/global.scss';
</script>

<div class="app">
	<AudioProgressBorder>
		<main>
			<slot />
		</main>
		<Toasts />
	</AudioProgressBorder>
</div>

<style lang="scss">
	.app {
		box-sizing: border-box;
		display: flex;
		flex-direction: row;
		justify-content: center;
		height: 100dvh;
		width: 100dvw;
		background-color: var(--white);

		main {
			flex: 1;
			display: flex;
			flex-direction: column;
			max-width: 64rem;
			box-sizing: border-box;
		}
	}
</style>
