<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		class?: string;
		contentClass?: string;
		spotlightColor?: string;
		children?: Snippet;
	};

	let {
		class: className = '',
		contentClass = '',
		spotlightColor = 'rgba(108, 92, 231, 0.25)',
		children
	}: Props = $props();

	let divRef: HTMLDivElement;
	let isFocused = $state(false);
	let posX = $state(0);
	let posY = $state(0);
	let opacity = $state(0);

	function handleMouseMove(e: MouseEvent) {
		if (!divRef || isFocused) return;
		const rect = divRef.getBoundingClientRect();
		posX = e.clientX - rect.left;
		posY = e.clientY - rect.top;
	}
	function handleFocus() {
		isFocused = true;
		opacity = 0.6;
	}
	function handleBlur() {
		isFocused = false;
		opacity = 0;
	}
	function handleMouseEnter() {
		opacity = 0.6;
	}
	function handleMouseLeave() {
		opacity = 0;
	}
</script>

<div
	bind:this={divRef}
	role="group"
	tabindex="-1"
	onmousemove={handleMouseMove}
	onfocus={handleFocus}
	onblur={handleBlur}
	onmouseenter={handleMouseEnter}
	onmouseleave={handleMouseLeave}
	class="spotlight-card {className}"
>
	<div
		class="spotlight-glow"
		style="opacity:{opacity};background:radial-gradient(400px circle at {posX}px {posY}px, {spotlightColor}, transparent 80%);"
	></div>
	<div class="spotlight-content {contentClass}">
		{@render children?.()}
	</div>
</div>

<style>
	.spotlight-card {
		position: relative;
		border-radius: var(--radius-lg, 16px);
		background: var(--color-surface, rgba(28, 29, 41, 0.7));
		border: 1px solid var(--color-border, rgba(255, 255, 255, 0.08));
		overflow: hidden;
		transition: border-color 0.2s ease, transform 0.2s ease;
	}

	.spotlight-card:hover {
		border-color: rgba(108, 92, 231, 0.35);
	}

	.spotlight-glow {
		pointer-events: none;
		position: absolute;
		inset: 0;
		transition: opacity 0.4s ease-in-out;
		z-index: 0;
	}

	.spotlight-content {
		position: relative;
		z-index: 1;
		width: 100%;
		height: 100%;
	}
</style>
