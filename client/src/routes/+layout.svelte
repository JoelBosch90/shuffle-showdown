<script>
	import '@fortawesome/fontawesome-free/css/all.min.css'
	import '../styles/global.scss';
</script>

<div class="app">
	<slot />
</div>

<style>
	.app {
		box-sizing: border-box;
		display: flex;
		flex-direction: row;
		justify-content: center;

		height: 100dvh;
		width: 100dvw;

		background-color: var(--white);
	}
</style>
