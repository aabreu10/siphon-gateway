<script lang="ts">
	interface Props {
		text: string;
		speed?: number;
		class?: string;
	}

	let { text, speed = 4, class: className = '' }: Props = $props();
</script>

<span
	class="shiny-text {className}"
	style="--shimmer-duration: {speed}s;"
>
	{text}
</span>

<style>
	.shiny-text {
		display: inline-block;
		background: linear-gradient(
			120deg,
			var(--color-text-primary) 0%,
			var(--color-text-primary) 35%,
			#ffffff 50%,
			#a29bfe 65%,
			var(--color-text-primary) 80%,
			var(--color-text-primary) 100%
		);
		background-size: 250% 100%;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
		animation: shimmer var(--shimmer-duration) linear infinite;
	}

	@keyframes shimmer {
		0% {
			background-position: 100% 0;
		}
		100% {
			background-position: -100% 0;
		}
	}
</style>
