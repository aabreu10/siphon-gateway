<script lang="ts">
	import type { Webhook } from '$lib/types';
	import { replayWebhook } from '$lib/api';
	import StatusBadge from './StatusBadge.svelte';

	interface Props {
		webhooks: Webhook[];
		oninspect: (webhook: Webhook) => void;
	}

	let { webhooks, oninspect }: Props = $props();
	let replayingIds = $state<Set<string>>(new Set());

	async function handleReplay(id: string) {
		replayingIds.add(id);
		replayingIds = replayingIds; // trigger reactivity
		try {
			await replayWebhook(id);
		} catch (err) {
			console.error('Replay failed:', err);
		} finally {
			replayingIds.delete(id);
			replayingIds = replayingIds;
		}
	}

	function truncateId(id: string): string {
		return id.slice(0, 8) + '…';
	}

	function timeAgo(dateStr: string): string {
		const diff = Date.now() - new Date(dateStr).getTime();
		const secs = Math.floor(diff / 1000);
		if (secs < 60) return `${secs}s ago`;
		const mins = Math.floor(secs / 60);
		if (mins < 60) return `${mins}m ago`;
		const hrs = Math.floor(mins / 60);
		if (hrs < 24) return `${hrs}h ago`;
		return new Date(dateStr).toLocaleDateString();
	}
</script>

<div class="table-container glass-panel">
	<div class="table-scroll">
		<table id="webhook-table">
			<thead>
				<tr>
					<th>Status</th>
					<th>Source</th>
					<th>Webhook ID</th>
					<th>Retries</th>
					<th>Target</th>
					<th>Received</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each webhooks as webhook, i (webhook.id)}
					<tr
						class="webhook-row"
						style="animation-delay: {Math.min(i * 30, 300)}ms"
						onclick={() => oninspect(webhook)}
					>
						<td>
							<StatusBadge status={webhook.status} retryCount={webhook.retry_count} />
						</td>
						<td>
							<span class="source-chip">{webhook.source}</span>
						</td>
						<td>
							<code class="webhook-id" title={webhook.id}>{truncateId(webhook.id)}</code>
						</td>
						<td>
							<span class="retry-count" class:retrying={webhook.status === 'PENDING' && webhook.retry_count > 0}>
								{webhook.retry_count}
							</span>
						</td>
						<td>
							<span class="target-url" title={webhook.target_url}>
								{webhook.target_url.replace(/^https?:\/\//, '').slice(0, 30)}
							</span>
						</td>
						<td>
							<span class="timestamp" title={new Date(webhook.created_at).toLocaleString()}>
								{timeAgo(webhook.created_at)}
							</span>
						</td>
						<td>
							<div class="action-group">
								<button
									class="btn-inspect"
									onclick={(e) => { e.stopPropagation(); oninspect(webhook); }}
									title="Inspect payload"
								>
									<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
										<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
										<circle cx="12" cy="12" r="3"></circle>
									</svg>
								</button>
								{#if webhook.status === 'FAILED_DLQ'}
									<button
										class="btn-replay"
										onclick={(e) => { e.stopPropagation(); handleReplay(webhook.id); }}
										disabled={replayingIds.has(webhook.id)}
										title="Replay webhook"
									>
										{#if replayingIds.has(webhook.id)}
											<svg class="spin" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
												<path d="M21 12a9 9 0 1 1-6.219-8.56"></path>
											</svg>
										{:else}
											<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
												<polyline points="1 4 1 10 7 10"></polyline>
												<path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
											</svg>
										{/if}
										Replay
									</button>
								{/if}
							</div>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="7" class="empty-state">
							<div class="empty-content">
								<svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" opacity="0.3">
									<path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
									<polyline points="13 2 13 9 20 9"></polyline>
								</svg>
								<span>No webhooks received yet</span>
								<span class="empty-hint">Send a POST request to /api/v1/webhook to get started</span>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>

<style>
	.table-container {
		overflow: hidden;
	}

	.table-scroll {
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.85rem;
	}

	thead th {
		position: sticky;
		top: 0;
		background: var(--color-bg-elevated);
		padding: 12px 16px;
		text-align: left;
		font-size: 0.7rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--color-text-muted);
		border-bottom: 1px solid var(--color-border);
		white-space: nowrap;
	}

	tbody td {
		padding: 12px 16px;
		border-bottom: 1px solid var(--color-border);
		vertical-align: middle;
	}

	.webhook-row {
		cursor: pointer;
		transition: background var(--transition-fast);
		animation: fadeInUp 0.3s ease both;
	}

	.webhook-row:hover {
		background: rgba(255, 255, 255, 0.02);
	}

	.source-chip {
		font-family: var(--font-mono);
		font-size: 0.78rem;
		font-weight: 500;
		color: var(--color-accent);
		background: var(--color-accent-glow);
		padding: 3px 8px;
		border-radius: 4px;
	}

	.webhook-id {
		font-family: var(--font-mono);
		font-size: 0.78rem;
		color: var(--color-text-secondary);
		background: rgba(255, 255, 255, 0.03);
		padding: 2px 6px;
		border-radius: 3px;
	}

	.retry-count {
		font-family: var(--font-mono);
		font-weight: 600;
		color: var(--color-text-secondary);
	}

	.retry-count.retrying {
		color: var(--color-warning);
	}

	.target-url {
		font-family: var(--font-mono);
		font-size: 0.75rem;
		color: var(--color-text-muted);
		max-width: 200px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		display: block;
	}

	.timestamp {
		font-size: 0.78rem;
		color: var(--color-text-muted);
		white-space: nowrap;
	}

	.action-group {
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.btn-inspect,
	.btn-replay {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		padding: 5px 10px;
		border-radius: var(--radius-sm);
		font-size: 0.75rem;
		font-family: var(--font-sans);
		font-weight: 500;
		cursor: pointer;
		transition: all var(--transition-fast);
		border: 1px solid var(--color-border);
		background: transparent;
		color: var(--color-text-secondary);
	}

	.btn-inspect:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.05);
		border-color: var(--color-border-active);
	}

	.btn-replay {
		color: var(--color-accent);
		border-color: rgba(108, 92, 231, 0.2);
	}

	.btn-replay:hover:not(:disabled) {
		background: var(--color-accent-glow);
		border-color: rgba(108, 92, 231, 0.4);
	}

	.btn-replay:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.spin {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from { transform: rotate(0deg); }
		to { transform: rotate(360deg); }
	}

	.empty-state {
		text-align: center;
		padding: 60px 20px !important;
	}

	.empty-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 10px;
		color: var(--color-text-muted);
	}

	.empty-hint {
		font-family: var(--font-mono);
		font-size: 0.75rem;
		opacity: 0.5;
	}
</style>
