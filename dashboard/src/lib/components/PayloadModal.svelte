<script lang="ts">
	import { onMount } from 'svelte';
	import type { Webhook } from '$lib/types';
	import { fetchDeliveryLogs } from '$lib/api';

	interface Props {
		webhook: Webhook;
		onclose: () => void;
	}

	let { webhook, onclose }: Props = $props();
	
	let activeTab = $state<'payload' | 'logs'>('payload');
	let deliveryLogs = $state<any[]>([]);
	let loadingLogs = $state(true);

	onMount(async () => {
		try {
			deliveryLogs = await fetchDeliveryLogs(webhook.id);
		} catch (err) {
			console.error('Failed to load logs', err);
		} finally {
			loadingLogs = false;
		}
	});

	function formatJson(obj: unknown): string {
		return JSON.stringify(obj, null, 2);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose();
	}

	function handleBackdropClick(e: MouseEvent) {
		if ((e.target as HTMLElement).classList.contains('modal-backdrop')) {
			onclose();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
<div class="modal-backdrop" role="presentation" onclick={handleBackdropClick}>
	<div class="modal-panel animate-fade-in-up" role="dialog" aria-modal="true" aria-label="Webhook Payload Inspector">
		<header class="modal-header">
			<div class="modal-title-group">
				<h2 class="modal-title">Payload Inspector</h2>
				<span class="modal-id">{webhook.id}</span>
			</div>
			<button class="modal-close" onclick={onclose} aria-label="Close modal">
				<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
					<line x1="18" y1="6" x2="6" y2="18"></line>
					<line x1="6" y1="6" x2="18" y2="18"></line>
				</svg>
			</button>
		</header>

		<div class="modal-meta">
			<div class="meta-item">
				<span class="meta-label">Source</span>
				<span class="meta-value source-tag">{webhook.source}</span>
			</div>
			<div class="meta-item">
				<span class="meta-label">Status</span>
				<span class="meta-value status-{webhook.status.toLowerCase()}">{webhook.status}</span>
			</div>
			<div class="meta-item">
				<span class="meta-label">Retries</span>
				<span class="meta-value">{webhook.retry_count}</span>
			</div>
			<div class="meta-item">
				<span class="meta-label">Target</span>
				<span class="meta-value mono">{webhook.target_url}</span>
			</div>
			<div class="meta-item">
				<span class="meta-label">Received</span>
				<span class="meta-value">{new Date(webhook.created_at).toLocaleString()}</span>
			</div>
			<div class="meta-item">
				<span class="meta-label">Last Update</span>
				<span class="meta-value">{new Date(webhook.updated_at).toLocaleString()}</span>
			</div>
		</div>

		<div class="modal-tabs">
			<button class="tab-btn" class:active={activeTab === 'payload'} onclick={() => (activeTab = 'payload')}>Payload</button>
			<button class="tab-btn" class:active={activeTab === 'logs'} onclick={() => (activeTab = 'logs')}>Delivery Logs</button>
		</div>

		<div class="modal-body">
			{#if activeTab === 'payload'}
				<div class="payload-header">
					<span class="payload-label">Raw JSON Payload</span>
					<button class="copy-btn" onclick={() => navigator.clipboard.writeText(formatJson(webhook.payload))}>
						<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
							<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
							<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
						</svg>
						Copy
					</button>
				</div>
				<pre class="payload-json"><code>{formatJson(webhook.payload)}</code></pre>
			{:else}
				<div class="logs-container">
					{#if loadingLogs}
						<div class="logs-empty">Loading logs...</div>
					{:else if deliveryLogs.length === 0}
						<div class="logs-empty">No delivery attempts yet.</div>
					{:else}
						{#each deliveryLogs as log}
							<div class="log-entry">
								<div class="log-header">
									<span class="log-attempt">Attempt #{log.attempt_number}</span>
									<span class="log-status {log.status_code >= 200 && log.status_code < 300 ? 'status-success' : 'status-failed_dlq'}">
										HTTP {log.status_code}
									</span>
									<span class="log-time">{new Date(log.created_at).toLocaleString()}</span>
								</div>
								{#if log.response_body}
									<pre class="log-body">{log.response_body}</pre>
								{/if}
							</div>
						{/each}
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		z-index: 1000;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--color-bg-overlay);
		backdrop-filter: blur(8px);
		-webkit-backdrop-filter: blur(8px);
		padding: 24px;
	}

	.modal-panel {
		background: var(--color-bg-elevated);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		width: 100%;
		max-width: 640px;
		max-height: 85vh;
		display: flex;
		flex-direction: column;
		overflow: hidden;
		box-shadow:
			0 0 0 1px rgba(255, 255, 255, 0.03),
			0 24px 48px rgba(0, 0, 0, 0.5);
	}

	.modal-header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		padding: 24px 24px 16px;
		border-bottom: 1px solid var(--color-border);
	}

	.modal-title {
		font-size: 1.1rem;
		font-weight: 700;
		color: var(--color-text-primary);
	}

	.modal-id {
		font-family: var(--font-mono);
		font-size: 0.75rem;
		color: var(--color-text-muted);
		margin-top: 4px;
		display: block;
	}

	.modal-close {
		background: none;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		color: var(--color-text-secondary);
		cursor: pointer;
		padding: 6px;
		display: flex;
		transition: all var(--transition-fast);
	}

	.modal-close:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.05);
		border-color: var(--color-border-active);
	}

	.modal-meta {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 12px;
		padding: 16px 24px;
		border-bottom: 1px solid var(--color-border);
	}

	.meta-item {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.meta-label {
		font-size: 0.7rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.06em;
		color: var(--color-text-muted);
	}

	.meta-value {
		font-size: 0.85rem;
		color: var(--color-text-primary);
	}

	.meta-value.mono {
		font-family: var(--font-mono);
		font-size: 0.78rem;
		word-break: break-all;
	}

	.source-tag {
		color: var(--color-accent);
		font-weight: 600;
	}

	.status-success { color: var(--color-success); }
	.status-failed_dlq { color: var(--color-danger); }
	.status-pending { color: var(--color-warning); }

	.modal-body {
		padding: 16px 24px 24px;
		overflow-y: auto;
		flex: 1;
	}

	.payload-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 10px;
	}

	.payload-label {
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.06em;
		color: var(--color-text-muted);
	}

	.copy-btn {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		background: none;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
		color: var(--color-text-secondary);
		font-size: 0.75rem;
		font-family: var(--font-sans);
		padding: 4px 10px;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.copy-btn:hover {
		color: var(--color-text-primary);
		background: rgba(255, 255, 255, 0.05);
		border-color: var(--color-border-active);
	}

	.payload-json {
		background: rgba(0, 0, 0, 0.3);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: 16px;
		overflow-x: auto;
		font-family: var(--font-mono);
		font-size: 0.8rem;
		line-height: 1.7;
		color: var(--color-text-primary);
		tab-size: 2;
	}

	.payload-json code {
		white-space: pre;
	}

	.modal-tabs {
		display: flex;
		gap: 16px;
		padding: 0 24px;
		border-bottom: 1px solid var(--color-border);
		background: rgba(0,0,0,0.1);
	}

	.tab-btn {
		background: none;
		border: none;
		color: var(--color-text-muted);
		padding: 12px 0;
		font-size: 0.85rem;
		font-weight: 600;
		cursor: pointer;
		border-bottom: 2px solid transparent;
		transition: all var(--transition-fast);
	}

	.tab-btn:hover {
		color: var(--color-text-primary);
	}

	.tab-btn.active {
		color: var(--color-accent);
		border-bottom-color: var(--color-accent);
	}

	.logs-container {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.logs-empty {
		padding: 32px;
		text-align: center;
		color: var(--color-text-muted);
		font-size: 0.9rem;
	}

	.log-entry {
		background: rgba(0, 0, 0, 0.2);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		overflow: hidden;
	}

	.log-header {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 10px 16px;
		background: rgba(255, 255, 255, 0.02);
		border-bottom: 1px solid var(--color-border);
		font-size: 0.8rem;
	}

	.log-attempt {
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.log-status {
		font-weight: 700;
		font-family: var(--font-mono);
	}

	.log-time {
		margin-left: auto;
		color: var(--color-text-muted);
		font-size: 0.75rem;
	}

	.log-body {
		padding: 12px 16px;
		font-family: var(--font-mono);
		font-size: 0.75rem;
		color: var(--color-text-secondary);
		white-space: pre-wrap;
		word-break: break-all;
		margin: 0;
	}
</style>
