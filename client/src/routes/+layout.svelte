<script>
	import '../styles/global.scss';
</script>

<div class="app">
	<slot />
</div>

<style>
	.app {
		display: flex;
		flex-direction: row;
		min-height: 100vh;
		justify-content: center;
	}
</style>
