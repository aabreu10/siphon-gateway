<script lang="ts">
	import { onMount } from 'svelte';

	let mouseX = $state(50);
	let mouseY = $state(30);

	function handleMouseMove(e: MouseEvent) {
		const target = e.currentTarget as HTMLElement;
		const rect = target.getBoundingClientRect();
		mouseX = ((e.clientX - rect.left) / rect.width) * 100;
		mouseY = ((e.clientY - rect.top) / rect.height) * 100;
	}
</script>

<div
	class="grid-background"
	role="presentation"
	onmousemove={handleMouseMove}
	style="--mouse-x: {mouseX}%; --mouse-y: {mouseY}%;"
>
	<!-- Retro Perspective Grid -->
	<div class="grid-lines"></div>

	<!-- Interactive Radial Glow tracking mouse -->
	<div class="radial-glow"></div>

	<!-- Ambient Purple/Indigo Orbs -->
	<div class="ambient-orb orb-1"></div>
	<div class="ambient-orb orb-2"></div>
</div>

<style>
	.grid-background {
		position: absolute;
		inset: 0;
		overflow: hidden;
		pointer-events: auto;
		z-index: 0;
	}

	.grid-lines {
		position: absolute;
		inset: 0;
		background-size: 50px 50px;
		background-image:
			linear-gradient(to right, rgba(108, 92, 231, 0.08) 1px, transparent 1px),
			linear-gradient(to bottom, rgba(108, 92, 231, 0.08) 1px, transparent 1px);
		mask-image: radial-gradient(circle 800px at 50% 30%, black 20%, transparent 100%);
		-webkit-mask-image: radial-gradient(circle 800px at 50% 30%, black 20%, transparent 100%);
		animation: pulse-grid 8s ease-in-out infinite alternate;
	}

	.radial-glow {
		position: absolute;
		inset: 0;
		background: radial-gradient(
			600px circle at var(--mouse-x) var(--mouse-y),
			rgba(108, 92, 231, 0.15),
			rgba(0, 230, 118, 0.04) 40%,
			transparent 80%
		);
		transition: background 0.15s ease-out;
	}

	.ambient-orb {
		position: absolute;
		border-radius: 50%;
		filter: blur(90px);
		opacity: 0.35;
		pointer-events: none;
	}

	.orb-1 {
		width: 450px;
		height: 450px;
		background: #6c5ce7;
		top: -10%;
		left: 15%;
		animation: float-slow 12s ease-in-out infinite alternate;
	}

	.orb-2 {
		width: 500px;
		height: 500px;
		background: #00e676;
		bottom: 10%;
		right: 10%;
		opacity: 0.15;
		animation: float-slow 16s ease-in-out infinite alternate-reverse;
	}

	@keyframes float-slow {
		0% {
			transform: translate(0, 0) scale(1);
		}
		50% {
			transform: translate(-30px, 40px) scale(1.1);
		}
		100% {
			transform: translate(20px, -20px) scale(0.95);
		}
	}

	@keyframes pulse-grid {
		0% {
			opacity: 0.6;
		}
		100% {
			opacity: 1;
		}
	}
</style>
