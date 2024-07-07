<script lang="ts">
</script>

<div class="loading-background loading">
	<div class="loading-slider-top"></div>
	<div class="loading-slider-bottom"></div>
	<div class="content-wrapper">
		<slot />
	</div>
</div>

<style lang="scss">
	.loading-background {
		--loading-background-padding: 0.5rem;
		--loading-background-default: var(--white);
		--loading-percentage: 0;

		display: flex;
		width: 100%;
		padding: var(--loading-background-padding);
		background-color: var(--loading-background-default);

		&.loading {
			background: var(--rainbow);
			background-size: var(--animation-background-flow-size);
			animation: animation-background-flow 10s linear infinite reverse;
		}

		.loading-slider-top {
			content: '';
			position: absolute;
			top: 0;
			left: 0;
			bottom: calc(100% - var(--loading-background-padding));
			right: 0;
			background-color: var(--loading-background-default);
			transform: translateX(calc(var(--loading-percentage) * 1%));
			transition: transform var(--animation-speed-slow);
		}

		.loading-slider-bottom {
			content: '';
			position: absolute;
			top: var(--loading-background-padding);
			left: 0;
			bottom: 0;
			right: 0;
			background-color: var(--loading-background-default);
		}

		.content-wrapper {
			border-radius: var(--loading-background-padding);
			width: 100%;
			display: flex;
			flex-direction: row;
			justify-content: center;
			background-color: var(--loading-background-default);
		}
	}
</style>
