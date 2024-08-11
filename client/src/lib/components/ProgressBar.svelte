<script lang="ts">
	export let total = 0;
	export let current = 0;
	export let label = '';
	export let labelRight = '';
	export let done = false;
	let id = `${label}-progress-bar`;
	$: title = labelRight ? `${label} (${labelRight})` : label;
</script>

<div class="progress-bar" class:done>
	<label for={id} {title}>
		<span>{label}</span>
		<span>{labelRight}</span>
	</label>
	<progress {id} value={current} max={total} style={`--progress: ${(current / total) * 100}%;`} />
</div>

<style lang="scss">
	.progress-bar {
		--padding: 0.5em;
		--fill-color: var(--purple);

		position: relative;
		box-sizing: border-box;
		overflow: hidden;
		border-radius: var(--border-radius);
		background: var(--gray);

		label {
			position: absolute;
			top: 50%;
			left: calc(var(--padding) + 1ch);
			transform: translateY(-50%);
			z-index: calc(var(--default-level) + 1);

			width: calc(100% - 2 * var(--padding) - 2ch);
			display: flex;
			justify-content: space-between;
			overflow: hidden;
			color: var(--white);

			span {
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}
		}

		progress {
			--progress: 0%;
			width: 100%;
			height: calc(1em + 2 * var(--padding));

			border: none;
			appearance: none;

			background: var(--gray);
			color: var(--white);

			&::-webkit-progress-bar {
				background-color: var(--gray);
			}

			&::-webkit-progress-value {
				background: var(--fill-color);
				transition: width 1s;
			}

			&::-moz-progress-bar {
				background: var(--fill-color);
			}
		}

		&.done {
			--fill-color: var(--green);
		}
	}
</style>
