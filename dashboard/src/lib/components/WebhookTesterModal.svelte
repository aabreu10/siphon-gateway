<script lang="ts">
	import { sendWebhook } from '$lib/api';

	interface Props {
		onclose: () => void;
	}

	let { onclose }: Props = $props();

	type Preset = 'stripe' | 'github' | 'shopify' | 'dlq_test' | 'custom';
	let selectedPreset = $state<Preset>('stripe');

	let source = $state('stripe');
	let targetUrl = $state('');
	let payloadText = $state(
		JSON.stringify(
			{
				event: 'payment_intent.succeeded',
				data: {
					id: 'pi_3MtwBwLkdIwHu7ix28a3tq5x',
					object: 'payment_intent',
					amount: 4999,
					currency: 'usd',
					status: 'succeeded',
					customer: 'cus_N2U7X9Q88a203b'
				}
			},
			null,
			2
		)
	);

	let loading = $state(false);
	let statusMsg = $state<{ type: 'success' | 'error'; text: string } | null>(null);

	// Validate JSON
	let isValidJson = $derived.by(() => {
		try {
			JSON.parse(payloadText);
			return true;
		} catch {
			return false;
		}
	});

	function selectPreset(preset: Preset) {
		selectedPreset = preset;
		statusMsg = null;

		if (preset === 'stripe') {
			source = 'stripe';
			targetUrl = '';
			payloadText = JSON.stringify(
				{
					event: 'payment_intent.succeeded',
					data: {
						id: 'pi_3MtwBwLkdIwHu7ix28a3tq5x',
						object: 'payment_intent',
						amount: 4999,
						currency: 'usd',
						status: 'succeeded',
						customer: 'cus_N2U7X9Q88a203b'
					}
				},
				null,
				2
			);
		} else if (preset === 'github') {
			source = 'github';
			targetUrl = '';
			payloadText = JSON.stringify(
				{
					event: 'push',
					repository: {
						name: 'siphon-gateway',
						full_name: 'aabreu10/siphon-gateway',
						private: true
					},
					pusher: {
						name: 'aabreu10',
						email: 'abreu@example.com'
					},
					ref: 'refs/heads/main'
				},
				null,
				2
			);
		} else if (preset === 'shopify') {
			source = 'shopify';
			targetUrl = '';
			payloadText = JSON.stringify(
				{
					event: 'orders/create',
					order: {
						id: 8209829119461,
						order_number: '#1024',
						total_price: '129.50',
						currency: 'USD',
						email: 'customer@shopify.com'
					}
				},
				null,
				2
			);
		} else if (preset === 'dlq_test') {
			source = 'dlq_simulator';
			targetUrl = 'http://localhost:9999/offline_endpoint';
			payloadText = JSON.stringify(
				{
					event: 'order.payment_failed',
					reason: 'downstream_timeout_test',
					timestamp: new Date().toISOString()
				},
				null,
				2
			);
		} else if (preset === 'custom') {
			source = 'custom_service';
			targetUrl = '';
			payloadText = JSON.stringify(
				{
					event: 'custom.event',
					message: 'Hello from Siphon Gateway Frontend Tester!'
				},
				null,
				2
			);
		}
	}

	async function handleSend(count = 1) {
		if (!isValidJson || loading) return;
		loading = true;
		statusMsg = null;

		try {
			const parsed = JSON.parse(payloadText);
			let lastId = '';

			for (let i = 0; i < count; i++) {
				const res = await sendWebhook(source, parsed, targetUrl);
				lastId = res.id;
				if (count > 1) {
					await new Promise((r) => setTimeout(r, 150));
				}
			}

			statusMsg = {
				type: 'success',
				text: count === 1
					? `✓ Ingested! ID: ${lastId.slice(0, 13)}… (Check dashboard log)`
					: `✓ Ingested ${count} burst webhooks! Check event log.`
			};
		} catch (err) {
			statusMsg = {
				type: 'error',
				text: err instanceof Error ? err.message : 'Failed to send webhook'
			};
		} finally {
			loading = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') onclose();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<!-- Modal Backdrop -->
<div
	class="modal-backdrop"
	role="dialog"
	aria-modal="true"
	aria-label="Webhook Tester Modal"
	tabindex="-1"
	onclick={(e) => {
		if (e.target === e.currentTarget) onclose();
	}}
	onkeydown={(e) => {
		if (e.key === 'Escape') onclose();
	}}
>
	<!-- Modal Box -->
	<div class="modal-box glass-panel animate-scale-in">
		<!-- Header -->
		<div class="modal-header">
			<div class="header-title-group">
				<span class="tester-badge">Frontend Tester</span>
				<h3 class="modal-title">Send Test Webhook</h3>
			</div>
			<button class="btn-close" onclick={onclose} aria-label="Close modal">✕</button>
		</div>

		<!-- Presets -->
		<div class="presets-row">
			<span class="preset-label">Templates:</span>
			<button
				class="preset-btn"
				class:active={selectedPreset === 'stripe'}
				onclick={() => selectPreset('stripe')}
			>
				Stripe Payment
			</button>
			<button
				class="preset-btn"
				class:active={selectedPreset === 'github'}
				onclick={() => selectPreset('github')}
			>
				GitHub Push
			</button>
			<button
				class="preset-btn"
				class:active={selectedPreset === 'shopify'}
				onclick={() => selectPreset('shopify')}
			>
				Shopify Order
			</button>
			<button
				class="preset-btn btn-dlq-preset"
				class:active={selectedPreset === 'dlq_test'}
				onclick={() => selectPreset('dlq_test')}
				title="Sends to offline port 9999 to test retries & DLQ"
			>
				⚠️ Test DLQ Failure
			</button>
			<button
				class="preset-btn"
				class:active={selectedPreset === 'custom'}
				onclick={() => selectPreset('custom')}
			>
				Custom JSON
			</button>
		</div>

		<!-- Form Fields -->
		<div class="fields-grid">
			<div class="field-group">
				<label for="webhook-source" class="field-label">Header: X-Webhook-Source</label>
				<input
					id="webhook-source"
					type="text"
					class="field-input"
					bind:value={source}
					placeholder="stripe, github, etc."
				/>
			</div>

			<div class="field-group">
				<label for="webhook-target" class="field-label">
					Target URL <span class="label-hint">(Optional — override downstream server)</span>
				</label>
				<input
					id="webhook-target"
					type="text"
					class="field-input"
					bind:value={targetUrl}
					placeholder="Leave blank for default server, or http://localhost:9999 to trigger DLQ"
				/>
			</div>
		</div>

		<!-- JSON Editor -->
		<div class="json-group">
			<div class="json-header">
				<label for="webhook-payload" class="field-label">JSON Payload Body</label>
				{#if !isValidJson}
					<span class="json-invalid">⚠️ Invalid JSON syntax</span>
				{:else}
					<span class="json-valid">✓ Valid JSON</span>
				{/if}
			</div>
			<textarea
				id="webhook-payload"
				class="json-editor font-mono"
				bind:value={payloadText}
				rows="11"
				spellcheck="false"
			></textarea>
		</div>

		<!-- Status Result Message -->
		{#if statusMsg}
			<div class="status-box {statusMsg.type === 'success' ? 'status-success' : 'status-error'}">
				<span>{statusMsg.text}</span>
			</div>
		{/if}

		<!-- Footer Actions -->
		<div class="modal-footer">
			<span class="footer-note">
				Sent webhooks insert live via Server-Sent Events (SSE).
			</span>
			<div class="action-buttons">
				<button class="btn-secondary" onclick={onclose}>Close</button>
				<button
					class="btn-burst"
					onclick={() => handleSend(5)}
					disabled={!isValidJson || loading}
					title="Send 5 rapid webhooks to test queue buffering"
				>
					⚡ Send 5 Burst
				</button>
				<button
					class="btn-primary"
					onclick={() => handleSend(1)}
					disabled={!isValidJson || loading}
				>
					{loading ? 'Sending…' : '⚡ Send Webhook'}
				</button>
			</div>
		</div>
	</div>
</div>

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		z-index: 1000;
		background: rgba(10, 11, 15, 0.85);
		backdrop-filter: blur(8px);
		-webkit-backdrop-filter: blur(8px);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 20px;
	}

	.modal-box {
		width: 100%;
		max-width: 680px;
		max-height: 90vh;
		display: flex;
		flex-direction: column;
		gap: 20px;
		padding: 28px;
		border-radius: var(--radius-xl);
		background: rgba(18, 19, 26, 0.95);
		border: 1px solid var(--color-border);
		box-shadow: 0 25px 60px rgba(0, 0, 0, 0.7);
		overflow-y: auto;
	}

	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		border-bottom: 1px solid var(--color-border);
		padding-bottom: 16px;
	}

	.header-title-group {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.tester-badge {
		font-family: var(--font-mono);
		font-size: 0.72rem;
		font-weight: 700;
		color: var(--color-accent);
		background: var(--color-accent-glow);
		padding: 4px 10px;
		border-radius: 20px;
		border: 1px solid rgba(108, 92, 231, 0.2);
	}

	.modal-title {
		font-size: 1.2rem;
		font-weight: 700;
		color: var(--color-text-primary);
	}

	.btn-close {
		width: 32px;
		height: 32px;
		border-radius: var(--radius-sm);
		display: flex;
		align-items: center;
		justify-content: center;
		background: transparent;
		border: 1px solid var(--color-border);
		color: var(--color-text-muted);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-close:hover {
		color: var(--color-text-primary);
		border-color: var(--color-border-active);
		background: rgba(255, 255, 255, 0.05);
	}

	/* Presets Row */
	.presets-row {
		display: flex;
		align-items: center;
		gap: 8px;
		flex-wrap: wrap;
	}

	.preset-label {
		font-size: 0.8rem;
		font-weight: 600;
		color: var(--color-text-muted);
		margin-right: 4px;
	}

	.preset-btn {
		font-size: 0.78rem;
		font-weight: 600;
		padding: 6px 14px;
		border-radius: 20px;
		border: 1px solid var(--color-border);
		background: rgba(255, 255, 255, 0.03);
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.preset-btn:hover {
		border-color: var(--color-border-active);
		color: var(--color-text-primary);
	}

	.preset-btn.active {
		background: var(--color-accent-glow);
		color: var(--color-accent);
		border-color: rgba(108, 92, 231, 0.4);
	}

	.btn-dlq-preset.active {
		background: var(--color-danger-glow);
		color: var(--color-danger);
		border-color: rgba(255, 82, 82, 0.4);
	}

	/* Form Fields */
	.fields-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 16px;
	}

	@media (max-width: 600px) {
		.fields-grid {
			grid-template-columns: 1fr;
		}
	}

	.field-group {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.field-label {
		font-size: 0.8rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.label-hint {
		font-size: 0.7rem;
		font-weight: 400;
		color: var(--color-text-muted);
	}

	.field-input {
		width: 100%;
		padding: 10px 14px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: rgba(0, 0, 0, 0.4);
		color: var(--color-text-primary);
		font-family: var(--font-mono);
		font-size: 0.85rem;
		transition: border-color var(--transition-fast);
	}

	.field-input:focus {
		outline: none;
		border-color: var(--color-accent);
	}

	/* JSON Editor */
	.json-group {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.json-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.json-valid {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-success);
	}

	.json-invalid {
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-danger);
	}

	.json-editor {
		width: 100%;
		padding: 14px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: rgba(0, 0, 0, 0.6);
		color: var(--color-text-primary);
		font-size: 0.82rem;
		line-height: 1.5;
		resize: vertical;
		transition: border-color var(--transition-fast);
	}

	.json-editor:focus {
		outline: none;
		border-color: var(--color-accent);
	}

	/* Status Box */
	.status-box {
		padding: 12px 16px;
		border-radius: var(--radius-md);
		font-size: 0.85rem;
		font-weight: 600;
		font-family: var(--font-mono);
	}

	.status-success {
		background: rgba(0, 230, 118, 0.1);
		border: 1px solid rgba(0, 230, 118, 0.3);
		color: var(--color-success);
	}

	.status-error {
		background: rgba(255, 82, 82, 0.1);
		border: 1px solid rgba(255, 82, 82, 0.3);
		color: var(--color-danger);
	}

	/* Footer */
	.modal-footer {
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-wrap: wrap;
		gap: 12px;
		padding-top: 16px;
		border-top: 1px solid var(--color-border);
	}

	.footer-note {
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}

	.action-buttons {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.btn-secondary {
		padding: 8px 18px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: transparent;
		color: var(--color-text-secondary);
		font-weight: 600;
		font-size: 0.85rem;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-secondary:hover {
		border-color: var(--color-border-active);
		color: var(--color-text-primary);
	}

	.btn-burst {
		padding: 8px 18px;
		border-radius: var(--radius-md);
		border: 1px solid rgba(255, 171, 64, 0.4);
		background: rgba(255, 171, 64, 0.1);
		color: var(--color-warning);
		font-weight: 600;
		font-size: 0.85rem;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.btn-burst:hover:not(:disabled) {
		background: rgba(255, 171, 64, 0.2);
		transform: translateY(-1px);
	}

	.btn-primary {
		padding: 8px 22px;
		border-radius: var(--radius-md);
		border: none;
		background: linear-gradient(135deg, #6c5ce7, #5f3dc4);
		color: #ffffff;
		font-weight: 600;
		font-size: 0.85rem;
		cursor: pointer;
		transition: all var(--transition-fast);
		box-shadow: 0 4px 15px rgba(108, 92, 231, 0.3);
	}

	.btn-primary:hover:not(:disabled) {
		transform: translateY(-1px);
		box-shadow: 0 6px 20px rgba(108, 92, 231, 0.5);
	}

	.btn-primary:disabled,
	.btn-burst:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none;
	}
</style>
