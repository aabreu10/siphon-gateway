<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		href?: string;
		onclick?: () => void;
		class?: string;
		children: Snippet;
	}

	let { href, onclick, class: className = '', children }: Props = $props();
</script>

{#if href}
	<a {href} class="border-glow-button {className}">
		<span class="glow-border"></span>
		<span class="button-content">
			{@render children()}
		</span>
	</a>
{:else}
	<button type="button" {onclick} class="border-glow-button {className}">
		<span class="glow-border"></span>
		<span class="button-content">
			{@render children()}
		</span>
	</button>
{/if}

<style>
	.border-glow-button {
		position: relative;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 2px;
		border-radius: var(--radius-lg);
		background: transparent;
		border: none;
		text-decoration: none;
		cursor: pointer;
		overflow: hidden;
		transition: transform var(--transition-fast), box-shadow var(--transition-fast);
	}

	.border-glow-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 30px rgba(108, 92, 231, 0.3);
	}

	.border-glow-button:active {
		transform: translateY(0);
	}

	.glow-border {
		position: absolute;
		inset: -150%;
		background: conic-gradient(
			from 0deg,
			transparent 0deg,
			#6c5ce7 90deg,
			#00e676 180deg,
			transparent 270deg,
			transparent 360deg
		);
		animation: rotate-border 4s linear infinite;
	}

	.button-content {
		position: relative;
		z-index: 1;
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 12px 24px;
		background: var(--color-bg-elevated);
		color: var(--color-text-primary);
		font-family: var(--font-sans);
		font-size: 0.95rem;
		font-weight: 600;
		border-radius: calc(var(--radius-lg) - 2px);
		transition: background var(--transition-fast);
	}

	.border-glow-button:hover .button-content {
		background: rgba(18, 19, 26, 0.9);
	}

	@keyframes rotate-border {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
</style>
