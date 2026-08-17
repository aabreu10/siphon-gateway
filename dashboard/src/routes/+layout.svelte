<script lang="ts">
	import { page } from '$app/stores';
	import WebhookTesterModal from '$lib/components/WebhookTesterModal.svelte';
	import '../app.css';

	let { children } = $props();
	let showTesterModal = $state(false);
</script>

<svelte:head>
	<title>Siphon Gateway — Fault-Tolerant Webhook Ingestion Engine</title>
</svelte:head>

<div class="app-shell">
	<header class="app-header">
		<div class="header-content">
			<a href="/" class="header-brand">
				<div class="logo-mark">
					<svg width="28" height="28" viewBox="0 0 32 32" fill="none">
						<defs>
							<linearGradient id="logo-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
								<stop offset="0%" stop-color="#14b8a6" />
								<stop offset="100%" stop-color="#0d9488" />
							</linearGradient>
						</defs>
						<path d="M16 2L4 8v16l12 6 12-6V8L16 2z" stroke="url(#logo-gradient)" stroke-width="2" fill="none" opacity="0.6"/>
						<path d="M16 6L8 10v12l8 4 8-4V10L16 6z" fill="url(#logo-gradient)" opacity="0.15"/>
						<circle cx="16" cy="16" r="3" fill="url(#logo-gradient)"/>
						<path d="M16 8v5M16 19v5M10 13l3.5 2M18.5 17l3.5 2M10 19l3.5-2M18.5 15l3.5-2" stroke="url(#logo-gradient)" stroke-width="1.5" opacity="0.5"/>
					</svg>
				</div>
				<div class="brand-text">
					<h1 class="brand-name">Siphon Gateway</h1>
					<span class="brand-subtitle">Webhook Engine</span>
				</div>
			</a>

			<nav class="header-nav" aria-label="Main Navigation">
				<a
					href="/"
					class="nav-link"
					class:active={$page.url.pathname === '/'}
				>
					Overview
				</a>
				<a
					href="/dashboard"
					class="nav-link"
					class:active={$page.url.pathname.startsWith('/dashboard')}
				>
					Dashboard
				</a>
			</nav>

			<div class="header-actions">
				<button class="btn-test-webhook" onclick={() => (showTesterModal = true)}>
					Test Webhook
				</button>
				<div class="header-status">
					<span class="status-dot"></span>
					<span class="status-text">SSE Live</span>
				</div>
				{#if $page.url.pathname.startsWith('/dashboard')}
					<button class="btn-launch" style="background: var(--color-surface); border: 1px solid var(--color-border);" onclick={async () => {
						const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';
						await fetch(`${API_BASE}/api/v1/auth/logout`, { method: 'POST', credentials: 'include' });
						window.location.href = '/login';
					}}>
						Log out
					</button>
				{:else}
					<a href="/dashboard" class="btn-launch">
						Launch App
					</a>
				{/if}
			</div>
		</div>
	</header>

	<main class="app-main">
		{@render children()}
	</main>
</div>

{#if showTesterModal}
	<WebhookTesterModal onclose={() => (showTesterModal = false)} />
{/if}

<style>
	.app-shell {
		min-height: 100dvh;
		display: flex;
		flex-direction: column;
	}

	.app-header {
		position: sticky;
		top: 0;
		z-index: 100;
		background: var(--color-bg);
		border-bottom: 1px solid var(--color-border);
	}

	.header-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		max-width: 1400px;
		margin: 0 auto;
		padding: 14px 28px;
		gap: 20px;
	}

	.header-brand {
		display: flex;
		align-items: center;
		gap: 14px;
		text-decoration: none;
	}

	.logo-mark {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		border-radius: var(--radius-md);
		background: var(--color-accent-subtle);
		border: 1px solid rgba(20, 184, 166, 0.2);
	}

	.brand-name {
		font-size: 1.1rem;
		font-weight: 700;
		letter-spacing: -0.02em;
		color: var(--color-text-primary);
	}

	.brand-subtitle {
		font-size: 0.7rem;
		font-weight: 500;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		display: block;
	}

	/* Navigation */
	.header-nav {
		display: flex;
		align-items: center;
		gap: 4px;
		background: var(--color-bg-card);
		padding: 4px;
		border-radius: var(--radius-lg);
		border: 1px solid var(--color-border);
	}

	.nav-link {
		padding: 6px 18px;
		border-radius: var(--radius-md);
		font-size: 0.85rem;
		font-weight: 500;
		color: var(--color-text-secondary);
		text-decoration: none;
		transition: all var(--transition-fast);
	}

	.nav-link:hover {
		color: var(--color-text-primary);
	}

	.nav-link.active {
		background: var(--color-bg-elevated);
		color: var(--color-text-primary);
		border: 1px solid var(--color-border-active);
	}

	.header-actions {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.header-status {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 6px 14px;
		border-radius: var(--radius-md);
		border: 1px solid rgba(34, 197, 94, 0.15);
		background: var(--color-success-subtle);
	}

	.status-dot {
		width: 7px;
		height: 7px;
		border-radius: 50%;
		background: var(--color-success);
	}

	.status-text {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-success);
	}

	.btn-test-webhook {
		display: inline-flex;
		align-items: center;
		padding: 6px 14px;
		border-radius: var(--radius-md);
		font-size: 0.82rem;
		font-weight: 600;
		color: var(--color-warning);
		background: var(--color-warning-subtle);
		border: 1px solid rgba(245, 158, 11, 0.2);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-test-webhook:hover {
		border-color: var(--color-warning);
	}

	.btn-launch {
		display: inline-flex;
		align-items: center;
		padding: 6px 16px;
		border-radius: var(--radius-md);
		font-size: 0.82rem;
		font-weight: 600;
		color: #ffffff;
		background: var(--color-accent);
		text-decoration: none;
		transition: all var(--transition-fast);
	}

	.btn-launch:hover {
		background: #0d9488;
	}

	.app-main {
		flex: 1;
		max-width: 1400px;
		margin: 0 auto;
		padding: 28px;
		width: 100%;
	}

	@media (max-width: 768px) {
		.header-content {
			flex-direction: column;
			gap: 14px;
		}

		.header-actions {
			width: 100%;
			justify-content: space-between;
		}
	}
</style>
