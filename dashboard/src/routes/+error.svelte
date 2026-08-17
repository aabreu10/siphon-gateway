<script lang="ts">
	import { page } from '$app/stores';
</script>

<svelte:head>
	<title>Page Not Found — Siphon Gateway</title>
</svelte:head>

<div class="error-container">
	<div class="error-card">
		<span class="error-code">
			{$page.status}
		</span>
		<h1 class="error-title">Page not found</h1>
		<p class="error-desc">
			{#if $page.error?.message}
				{$page.error.message}
			{:else}
				The page you're looking for doesn't exist or has been moved.
			{/if}
		</p>
		<a href="/" class="btn-home">Back to home</a>
	</div>
</div>

<style>
	.error-container {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 70vh;
	}

	.error-card {
		text-align: center;
		max-width: 400px;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
	}

	.error-code {
		font-family: var(--font-mono);
		font-size: 4rem;
		font-weight: 800;
		color: var(--color-text-muted);
		line-height: 1;
		letter-spacing: -0.04em;
	}

	.error-title {
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--color-text-primary);
	}

	.error-desc {
		font-size: 0.95rem;
		color: var(--color-text-secondary);
		line-height: 1.6;
	}

	.btn-home {
		display: inline-flex;
		align-items: center;
		padding: 10px 20px;
		border-radius: var(--radius-md);
		background: var(--color-accent);
		color: #ffffff;
		font-weight: 600;
		font-size: 0.9rem;
		text-decoration: none;
		transition: background var(--transition-fast);
		margin-top: 8px;
	}

	.btn-home:hover {
		background: #0d9488;
	}
</style>
