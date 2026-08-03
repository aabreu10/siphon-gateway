<script lang="ts">
	import { onMount } from 'svelte';
	import type { Webhook, SSEEvent } from '$lib/types';
	import { fetchWebhooks, getSSEUrl } from '$lib/api';
	import WebhookTable from '$lib/components/WebhookTable.svelte';
	import PayloadModal from '$lib/components/PayloadModal.svelte';
	import WebhookTesterModal from '$lib/components/WebhookTesterModal.svelte';
	import Balatro from '$lib/components/svelte-bits/Balatro.svelte';
	import CountUp from '$lib/components/svelte-bits/CountUp.svelte';
	import SpotlightCard from '$lib/components/svelte-bits/SpotlightCard.svelte';

	let webhooks = $state<Webhook[]>([]);
	let inspecting = $state<Webhook | null>(null);
	let showTesterModal = $state(false);
	let sseConnected = $state(false);
	let totalCount = $state(0);

	// Computed stats
	let stats = $derived({
		total: totalCount,
		success: webhooks.filter((w) => w.status === 'SUCCESS').length,
		pending: webhooks.filter((w) => w.status === 'PENDING').length,
		failed: webhooks.filter((w) => w.status === 'FAILED_DLQ').length
	});

	onMount(() => {
		// Initial load (fire-and-forget — no cleanup needed)
		loadInitialData();

		// SSE connection
		const evtSource = new EventSource(getSSEUrl());

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
		} catch (err) {
			console.error('Failed to load webhooks:', err);
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

<div class="dashboard">
	<div class="dashboard-bg-wrapper">
		<Balatro color1="#6c5ce7" color2="#1e1b4b" color3="#0a0b0f" mouseInteraction={true} />
	</div>

	<!-- Stats Cards -->
	<section class="stats-grid" aria-label="Webhook Statistics">
		<SpotlightCard class="stat-card glass-panel" contentClass="stat-card-inner">
			<div class="stat-icon stat-icon-total">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
					<path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>
				</svg>
			</div>
			<div class="stat-content">
				<CountUp to={stats.total} class="stat-value" />
				<span class="stat-label">Total Webhooks</span>
			</div>
		</SpotlightCard>

		<SpotlightCard class="stat-card glass-panel" contentClass="stat-card-inner">
			<div class="stat-icon stat-icon-success">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
					<polyline points="20 6 9 17 4 12"></polyline>
				</svg>
			</div>
			<div class="stat-content">
				<CountUp to={stats.success} class="stat-value stat-value-success" />
				<span class="stat-label">Delivered</span>
			</div>
		</SpotlightCard>

		<SpotlightCard class="stat-card glass-panel" contentClass="stat-card-inner">
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
		</SpotlightCard>

		<SpotlightCard class="stat-card glass-panel" contentClass="stat-card-inner">
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
		</SpotlightCard>
	</section>

	<!-- Webhook Table -->
	<section class="table-section animate-fade-in-up" style="animation-delay: 150ms" aria-label="Webhook Log">
		<div class="section-header">
			<div class="header-left">
				<h2 class="section-title">Live Event Log</h2>
				<button class="btn-test-webhook-sm" onclick={() => (showTesterModal = true)}>
					⚡ Test Webhook
				</button>
			</div>
			<div class="sse-indicator" class:connected={sseConnected}>
				<span class="sse-dot"></span>
				<span class="sse-label">{sseConnected ? 'Live' : 'Reconnecting…'}</span>
			</div>
		</div>
		<WebhookTable {webhooks} oninspect={(w) => (inspecting = w)} />
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
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.dashboard-bg-wrapper {
		position: fixed;
		inset: 0;
		opacity: 0.25;
		pointer-events: none;
		z-index: -1;
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

	:global(.stat-card) {
		animation: fadeInUp 0.4s ease both;
	}

	:global(.stat-card-inner) {
		display: flex !important;
		flex-direction: row !important;
		align-items: center !important;
		justify-content: flex-start !important;
		gap: 20px !important;
		padding: 24px 28px !important;
		width: 100% !important;
		box-sizing: border-box !important;
	}

	:global(.stat-card:nth-child(1)) { animation-delay: 0ms; }
	:global(.stat-card:nth-child(2)) { animation-delay: 50ms; }
	:global(.stat-card:nth-child(3)) { animation-delay: 100ms; }
	:global(.stat-card:nth-child(4)) { animation-delay: 150ms; }

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
		background: var(--color-accent-glow);
		color: var(--color-accent);
		border: 1px solid rgba(108, 92, 231, 0.25);
	}

	.stat-icon-success {
		background: var(--color-success-glow);
		color: var(--color-success);
		border: 1px solid rgba(0, 230, 118, 0.25);
	}

	.stat-icon-pending {
		background: var(--color-warning-glow);
		color: var(--color-warning);
		border: 1px solid rgba(255, 171, 64, 0.25);
	}

	.stat-icon-failed {
		background: var(--color-danger-glow);
		color: var(--color-danger);
		border: 1px solid rgba(255, 82, 82, 0.25);
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
		border-radius: 20px;
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-warning);
		background: var(--color-warning-glow);
		border: 1px solid rgba(255, 171, 64, 0.3);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-test-webhook-sm:hover {
		transform: translateY(-1px);
		box-shadow: 0 3px 10px rgba(255, 171, 64, 0.25);
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
		border-radius: 20px;
		border: 1px solid var(--color-border);
		background: rgba(255, 255, 255, 0.02);
	}

	.sse-indicator.connected {
		border-color: rgba(0, 230, 118, 0.15);
		background: rgba(0, 230, 118, 0.05);
	}

	.sse-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-text-muted);
	}

	.sse-indicator.connected .sse-dot {
		background: var(--color-success);
		box-shadow: 0 0 6px var(--color-success-glow);
		animation: pulse-glow 2s ease-in-out infinite;
	}

	.sse-label {
		font-size: 0.72rem;
		font-weight: 600;
		color: var(--color-text-muted);
	}

	.sse-indicator.connected .sse-label {
		color: var(--color-success);
	}
</style>
