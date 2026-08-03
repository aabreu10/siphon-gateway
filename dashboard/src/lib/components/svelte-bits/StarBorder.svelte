<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children?: Snippet;
		as?: string;
		class?: string;
		color?: string;
		speed?: string;
		thickness?: number;
		[key: string]: unknown;
	};

	let {
		children,
		as = 'button',
		class: className = '',
		color = '#6c5ce7',
		speed = '6s',
		thickness = 1,
		...rest
	}: Props = $props();

	const gradientBg = $derived(`radial-gradient(circle, ${color}, transparent 10%)`);
</script>

<svelte:element
	this={as}
	class="star-border-wrapper {className}"
	style:padding="{thickness}px 0"
	{...rest}
>
	<div
		class="star-sweep-bottom"
		style:background={gradientBg}
		style:animation-duration={speed}
	></div>
	<div
		class="star-sweep-top"
		style:background={gradientBg}
		style:animation-duration={speed}
	></div>
	<div class="star-border-inner">
		{@render children?.()}
	</div>
</svelte:element>

<style>
	.star-border-wrapper {
		position: relative;
		display: inline-block;
		overflow: hidden;
		border-radius: 20px;
		cursor: pointer;
		background: transparent;
		border: none;
	}

	.star-sweep-bottom {
		position: absolute;
		width: 300%;
		height: 50%;
		opacity: 0.7;
		bottom: -11px;
		right: -250%;
		border-radius: 9999px;
		z-index: 0;
		animation-name: star-movement-bottom;
		animation-timing-function: linear;
		animation-iteration-count: infinite;
		animation-direction: alternate;
	}

	.star-sweep-top {
		position: absolute;
		width: 300%;
		height: 50%;
		opacity: 0.7;
		top: -10px;
		left: -250%;
		border-radius: 9999px;
		z-index: 0;
		animation-name: star-movement-top;
		animation-timing-function: linear;
		animation-iteration-count: infinite;
		animation-direction: alternate;
	}

	.star-border-inner {
		position: relative;
		z-index: 1;
		background: linear-gradient(180deg, #181926 0%, #12131a 100%);
		border: 1px solid rgba(108, 92, 231, 0.35);
		color: #ffffff;
		text-align: center;
		padding: 14px 28px;
		border-radius: 20px;
		font-weight: 600;
		transition: all 0.2s ease;
	}

	.star-border-wrapper:hover .star-border-inner {
		border-color: rgba(108, 92, 231, 0.7);
		transform: translateY(-1px);
		box-shadow: 0 4px 15px rgba(108, 92, 231, 0.25);
	}

	@keyframes star-movement-bottom {
		0% {
			transform: translate(0%, 0%);
			opacity: 1;
		}
		100% {
			transform: translate(-100%, 0%);
			opacity: 0;
		}
	}

	@keyframes star-movement-top {
		0% {
			transform: translate(0%, 0%);
			opacity: 1;
		}
		100% {
			transform: translate(100%, 0%);
			opacity: 0;
		}
	}
</style>
