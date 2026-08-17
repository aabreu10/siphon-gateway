<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import type { Webhook, SSEEvent } from '$lib/types';
	import { fetchWebhooks, getSSEUrl } from '$lib/api';
	import WebhookTable from '$lib/components/WebhookTable.svelte';
	import PayloadModal from '$lib/components/PayloadModal.svelte';
	import WebhookTesterModal from '$lib/components/WebhookTesterModal.svelte';
	import CountUp from '$lib/components/svelte-bits/CountUp.svelte';

	let webhooks = $state<Webhook[]>([]);
	let inspecting = $state<Webhook | null>(null);
	let showTesterModal = $state(false);
	let sseConnected = $state(false);
	let totalCount = $state(0);

	// Filtering
	let searchQuery = $state('');
	let statusFilter = $state('ALL');

	// Computed stats
	let stats = $derived({
		total: totalCount,
		success: webhooks.filter((w) => w.status === 'SUCCESS').length,
		pending: webhooks.filter((w) => w.status === 'PENDING').length,
		failed: webhooks.filter((w) => w.status === 'FAILED_DLQ').length
	});

	let filteredWebhooks = $derived(webhooks);

	let searchTimeout: ReturnType<typeof setTimeout>;
	
	function handleSearchChange() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			loadInitialData();
		}, 300);
	}

	$effect(() => {
		// when statusFilter changes, reload
		if (statusFilter) {
			loadInitialData();
		}
	});

	onMount(() => {
		// Initial load - if unauthenticated, fetchWebhooks will redirect
		loadInitialData();

		// SSE connection
		const evtSource = new EventSource(getSSEUrl(), { withCredentials: true });

		evtSource.addEventListener('connected', () => {
			sseConnected = true;
		});

		evtSource.onmessage = (event) => {
			try {
				const sseEvent: SSEEvent = JSON.parse(event.data);
				handleSSEEvent(sseEvent);
			} catch (err) {
				console.error('Failed to parse SSE event:', err);
			}
		};

		evtSource.onerror = () => {
			sseConnected = false;
		};

		return () => {
			evtSource.close();
		};
	});

	async function loadInitialData() {
		try {
			const data = await fetchWebhooks(200);
			webhooks = data.webhooks;
			totalCount = data.total;
		} catch (err: any) {
			console.error('Failed to load webhooks:', err);
			if (err.message === 'Unauthorized') {
				goto('/login');
			}
		}
	}

	function handleSSEEvent(event: SSEEvent) {
		const existingIdx = webhooks.findIndex((w) => w.id === event.webhook.id);

		if (existingIdx >= 0) {
			// Update existing webhook in-place
			webhooks[existingIdx] = {
				...webhooks[existingIdx],
				...event.webhook,
				updated_at: event.webhook.updated_at ?? new Date().toISOString()
			};
			webhooks = webhooks; // trigger reactivity
		} else {
			// Prepend new webhook
			webhooks = [
				{ ...event.webhook, created_at: event.webhook.created_at ?? new Date().toISOString(), updated_at: event.webhook.updated_at ?? new Date().toISOString() },
				...webhooks
			];
			totalCount++;
		}

		// If we're inspecting the same webhook, update modal data
		if (inspecting && inspecting.id === event.webhook.id) {
			inspecting = { ...inspecting, ...event.webhook };
		}
	}
</script>

<svelte:head>
	<title>Dashboard — Siphon Gateway</title>
</svelte:head>

<div class="dashboard">
	<!-- Stats Cards -->
	<section class="stats-grid" aria-label="Webhook Statistics">
		<div class="stat-card animate-fade-in-up">
			<div class="stat-icon stat-icon-total">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
					<path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>
				</svg>
			</div>
			<div class="stat-content">
				<CountUp to={stats.total} class="stat-value" />
				<span class="stat-label">Total Webhooks</span>
			</div>
		</div>

		<div class="stat-card animate-fade-in-up" style="animation-delay: 50ms">
			<div class="stat-icon stat-icon-success">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
					<polyline points="20 6 9 17 4 12"></polyline>
				</svg>
			</div>
			<div class="stat-content">
				<CountUp to={stats.success} class="stat-value stat-value-success" />
				<span class="stat-label">Delivered</span>
			</div>
		</div>

		<div class="stat-card animate-fade-in-up" style="animation-delay: 100ms">
			<div class="stat-icon stat-icon-pending">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
					<circle cx="12" cy="12" r="10"></circle>
					<polyline points="12 6 12 12 16 14"></polyline>
				</svg>
			</div>
			<div class="stat-content">
				<CountUp to={stats.pending} class="stat-value stat-value-pending" />
				<span class="stat-label">In Progress</span>
			</div>
		</div>

		<div class="stat-card animate-fade-in-up" style="animation-delay: 150ms">
			<div class="stat-icon stat-icon-failed">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
					<path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
					<line x1="12" y1="9" x2="12" y2="13"></line>
					<line x1="12" y1="17" x2="12.01" y2="17"></line>
				</svg>
			</div>
			<div class="stat-content">
				<CountUp to={stats.failed} class="stat-value stat-value-failed" />
				<span class="stat-label">Dead Letter Queue</span>
			</div>
		</div>
	</section>

	<!-- Integration Details -->
	<section class="integration-card animate-fade-in-up" style="animation-delay: 125ms" aria-label="Integration Details">
		<div class="integration-content">
			<div class="integration-text">
				<h2 class="section-title">Integration Details</h2>
				<p class="integration-desc">Point your third-party providers (Stripe, GitHub, etc.) to your Siphon Gateway ingestion URL to start buffering webhooks.</p>
			</div>
			<div class="integration-code-box">
				<span class="method-badge">POST</span>
				<code class="url-code">https://siphon.io/api/v1/webhook</code>
			</div>
		</div>
	</section>

	<!-- Webhook Table -->
	<section class="table-section animate-fade-in-up" style="animation-delay: 150ms" aria-label="Webhook Log">
		<div class="section-header">
			<div class="header-left">
				<h2 class="section-title">Live Event Log</h2>
				<button class="btn-test-webhook-sm" onclick={() => (showTesterModal = true)}>
					Simulate Event
				</button>
			</div>
			<div class="header-filters">
				<input type="text" class="search-input" placeholder="Search ID or Source..." bind:value={searchQuery} oninput={handleSearchChange} />
				<select class="status-select" bind:value={statusFilter} onchange={handleSearchChange}>
					<option value="ALL">All Statuses</option>
					<option value="SUCCESS">Success</option>
					<option value="PENDING">Pending</option>
					<option value="FAILED_DLQ">Failed DLQ</option>
				</select>
			</div>
			<div class="sse-indicator" class:connected={sseConnected}>
				<span class="sse-dot"></span>
				<span class="sse-label">{sseConnected ? 'Live' : 'Reconnecting...'}</span>
			</div>
		</div>
		<WebhookTable webhooks={filteredWebhooks} oninspect={(w) => (inspecting = w)} />
	</section>
</div>

<!-- Payload Modal -->
{#if inspecting}
	<PayloadModal webhook={inspecting} onclose={() => (inspecting = null)} />
{/if}

<!-- Tester Modal -->
{#if showTesterModal}
	<WebhookTesterModal onclose={() => (showTesterModal = false)} />
{/if}

<style>
	.dashboard {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	/* ── Stats Grid ──────────────────────────────────────────────────── */
	.stats-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 16px;
	}

	@media (max-width: 900px) {
		.stats-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}

	@media (max-width: 500px) {
		.stats-grid {
			grid-template-columns: 1fr;
		}
	}

	.stat-card {
		display: flex;
		flex-direction: row;
		align-items: center;
		gap: 20px;
		padding: 24px 28px;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		transition: border-color var(--transition-fast);
	}

	.stat-card:hover {
		border-color: var(--color-border-active);
	}

	.stat-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 54px;
		height: 54px;
		border-radius: var(--radius-lg);
		flex-shrink: 0;
	}

	.stat-icon-total {
		background: var(--color-accent-subtle);
		color: var(--color-accent);
		border: 1px solid rgba(20, 184, 166, 0.2);
	}

	.stat-icon-success {
		background: var(--color-success-subtle);
		color: var(--color-success);
		border: 1px solid rgba(34, 197, 94, 0.2);
	}

	.stat-icon-pending {
		background: var(--color-warning-subtle);
		color: var(--color-warning);
		border: 1px solid rgba(245, 158, 11, 0.2);
	}

	.stat-icon-failed {
		background: var(--color-danger-subtle);
		color: var(--color-danger);
		border: 1px solid rgba(239, 68, 68, 0.2);
	}

	.stat-content {
		display: flex;
		flex-direction: column;
		gap: 4px;
		flex: 1;
		min-width: 0;
	}

	:global(.stat-value) {
		font-size: 2.1rem;
		font-weight: 800;
		letter-spacing: -0.04em;
		line-height: 1.1;
		color: var(--color-text-primary);
		font-family: var(--font-mono);
	}

	:global(.stat-value-success) { color: var(--color-success); }
	:global(.stat-value-pending) { color: var(--color-warning); }
	:global(.stat-value-failed) { color: var(--color-danger); }

	.stat-label {
		font-size: 0.8rem;
		font-weight: 600;
		color: var(--color-text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	/* ── Integration Details ─────────────────────────────────────────── */
	.integration-card {
		padding: 20px 28px;
		border-radius: var(--radius-xl);
		border: 1px solid var(--color-border);
		background: var(--color-bg-card);
	}

	.integration-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 20px;
	}

	.integration-text {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.integration-desc {
		font-size: 0.85rem;
		color: var(--color-text-muted);
		margin: 0;
	}

	.integration-code-box {
		display: flex;
		align-items: center;
		gap: 12px;
		background: var(--color-bg-elevated);
		padding: 8px 16px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
	}

	.method-badge {
		font-size: 0.7rem;
		font-weight: 700;
		color: var(--color-success);
		background: var(--color-success-subtle);
		padding: 3px 8px;
		border-radius: var(--radius-sm);
	}

	.url-code {
		font-family: var(--font-mono);
		font-size: 0.85rem;
		color: var(--color-text-primary);
	}

	@media (max-width: 768px) {
		.integration-content {
			flex-direction: column;
			align-items: flex-start;
		}
	}

	/* ── Table Section ───────────────────────────────────────────────── */
	.section-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 14px;
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.btn-test-webhook-sm {
		display: inline-flex;
		align-items: center;
		padding: 4px 12px;
		border-radius: var(--radius-md);
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-warning);
		background: var(--color-warning-subtle);
		border: 1px solid rgba(245, 158, 11, 0.2);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-test-webhook-sm:hover {
		border-color: var(--color-warning);
	}

	.section-title {
		font-size: 1rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.sse-indicator {
		display: flex;
		align-items: center;
		gap: 6px;
		padding: 4px 12px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: transparent;
	}

	.sse-indicator.connected {
		border-color: rgba(34, 197, 94, 0.15);
		background: var(--color-success-subtle);
	}

	.sse-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-text-muted);
	}

	.sse-indicator.connected .sse-dot {
		background: var(--color-success);
	}

	.sse-label {
		font-size: 0.72rem;
		font-weight: 600;
		color: var(--color-text-muted);
	}

	.sse-indicator.connected .sse-label {
		color: var(--color-success);
	}

	.header-filters {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-left: auto;
		margin-right: 16px;
	}

	.search-input, .status-select {
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		color: var(--color-text-primary);
		padding: 6px 12px;
		border-radius: var(--radius-md);
		font-size: 0.8rem;
		outline: none;
		transition: border-color var(--transition-fast);
	}

	.search-input:focus, .status-select:focus {
		border-color: var(--color-accent);
	}

	.search-input {
		width: 200px;
	}
</style>
