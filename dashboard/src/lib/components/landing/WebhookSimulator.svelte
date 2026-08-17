<script lang="ts">
	type SimState = 'IDLE' | 'SUCCESS' | 'RETRYING' | 'DLQ';

	let simState = $state<SimState>('IDLE');
	let activeStep = $state(0);
	let logMessage = $state('System idle. Select a scenario to simulate.');
	let timer: ReturnType<typeof setTimeout> | null = null;

	function runSimulation(mode: SimState) {
		if (timer) clearTimeout(timer);
		simState = mode;
		activeStep = 1;

		if (mode === 'SUCCESS') {
			logMessage = 'POST /api/v1/webhook from Stripe received. Instant 200 OK sent.';
			timer = setTimeout(() => {
				activeStep = 2;
				logMessage = 'Payload persisted in RabbitMQ exchange [webhooks.direct].';
				timer = setTimeout(() => {
					activeStep = 3;
					logMessage = 'Worker delivered payload to downstream target (HTTP 200 OK).';
				}, 900);
			}, 900);
		} else if (mode === 'RETRYING') {
			logMessage = 'POST /api/v1/webhook received. Enqueueing to RabbitMQ...';
			timer = setTimeout(() => {
				activeStep = 2;
				logMessage = 'Downstream responded with 503 Service Unavailable.';
				timer = setTimeout(() => {
					activeStep = 3;
					logMessage = 'Exponential backoff triggered. Scheduled retry in 2s (Retry #1/5).';
				}, 900);
			}, 900);
		} else if (mode === 'DLQ') {
			logMessage = 'POST /api/v1/webhook received. Enqueueing to RabbitMQ...';
			timer = setTimeout(() => {
				activeStep = 2;
				logMessage = 'All 5 retry attempts exhausted. Route to Dead Letter Queue.';
				timer = setTimeout(() => {
					activeStep = 3;
					logMessage = 'Webhook isolated in FAILED_DLQ. Available for 1-click manual replay.';
				}, 900);
			}, 900);
		}
	}
</script>

<div class="simulator-card">
	<div class="sim-header">
		<div class="sim-title-group">
			<span class="sim-badge">Interactive Demo</span>
			<h3 class="sim-title">Architecture Sandbox</h3>
		</div>
		<div class="sim-controls">
			<button
				class="sim-btn btn-success"
				class:active={simState === 'SUCCESS'}
				onclick={() => runSimulation('SUCCESS')}
			>
				Simulate Delivered
			</button>
			<button
				class="sim-btn btn-warning"
				class:active={simState === 'RETRYING'}
				onclick={() => runSimulation('RETRYING')}
			>
				Simulate Retry
			</button>
			<button
				class="sim-btn btn-danger"
				class:active={simState === 'DLQ'}
				onclick={() => runSimulation('DLQ')}
			>
				Simulate DLQ
			</button>
		</div>
	</div>

	<!-- Interactive Flow Diagram -->
	<div class="diagram-area">
		<!-- Node 1: Provider -->
		<div class="node-box" class:active={activeStep >= 1}>
			<div class="node-icon">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
				</svg>
			</div>
			<span class="node-label">Stripe / GitHub</span>
			<span class="node-sub">Provider</span>
		</div>

		<!-- Connector 1 -->
		<div class="connector" class:active={activeStep >= 1}>
			<div class="packet" class:moving={activeStep === 1}></div>
			<span class="connector-label">200 OK</span>
		</div>

		<!-- Node 2: Gateway API -->
		<div class="node-box" class:active={activeStep >= 1}>
			<div class="node-icon">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<polygon points="12 2 2 7 12 12 22 7 12 2"></polygon>
					<polyline points="2 17 12 22 22 17"></polyline>
					<polyline points="2 12 12 17 22 12"></polyline>
				</svg>
			</div>
			<span class="node-label">Siphon API</span>
			<span class="node-sub">Go Gateway</span>
		</div>

		<!-- Connector 2 -->
		<div class="connector" class:active={activeStep >= 2}>
			<div class="packet" class:moving={activeStep === 2}></div>
			<span class="connector-label">Queue</span>
		</div>

		<!-- Node 3: RabbitMQ -->
		<div class="node-box" class:active={activeStep >= 2}>
			<div class="node-icon">
				<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<rect x="2" y="7" width="20" height="14" rx="2" ry="2"></rect>
					<path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"></path>
				</svg>
			</div>
			<span class="node-label">RabbitMQ</span>
			<span class="node-sub">Message Broker</span>
		</div>

		<!-- Connector 3 -->
		<div class="connector" class:active={activeStep >= 3}>
			<div
				class="packet"
				class:moving={activeStep === 3}
				class:packet-warning={simState === 'RETRYING'}
				class:packet-danger={simState === 'DLQ'}
			></div>
			<span class="connector-label">Deliver</span>
		</div>

		<!-- Node 4: Target / DLQ -->
		<div
			class="node-box"
			class:active={activeStep >= 3}
			class:node-success={simState === 'SUCCESS' && activeStep === 3}
			class:node-warning={simState === 'RETRYING' && activeStep === 3}
			class:node-danger={simState === 'DLQ' && activeStep === 3}
		>
			<div class="node-icon">
				{#if simState === 'SUCCESS' && activeStep === 3}
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<polyline points="20 6 9 17 4 12"></polyline>
					</svg>
				{:else if simState === 'RETRYING' && activeStep === 3}
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<circle cx="12" cy="12" r="10"></circle>
						<polyline points="12 6 12 12 16 14"></polyline>
					</svg>
				{:else if simState === 'DLQ' && activeStep === 3}
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
						<line x1="12" y1="9" x2="12" y2="13"></line>
						<line x1="12" y1="17" x2="12.01" y2="17"></line>
					</svg>
				{:else}
					<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
						<line x1="8" y1="21" x2="16" y2="21"></line>
						<line x1="12" y1="17" x2="12" y2="21"></line>
					</svg>
				{/if}
			</div>
			<span class="node-label">
				{#if simState === 'DLQ' && activeStep === 3}
					Dead Letter Queue
				{:else if simState === 'RETRYING' && activeStep === 3}
					Backoff Retry
				{:else}
					Internal Server
				{/if}
			</span>
			<span class="node-sub">Target Endpoint</span>
		</div>
	</div>

	<!-- Status Bar -->
	<div class="sim-status-bar">
		<span class="status-bar-text">
			{logMessage}
		</span>
	</div>
</div>

<style>
	.simulator-card {
		padding: 28px;
		border-radius: var(--radius-xl);
		display: flex;
		flex-direction: column;
		gap: 24px;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
	}

	.sim-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-wrap: wrap;
		gap: 16px;
	}

	.sim-title-group {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.sim-badge {
		font-family: var(--font-mono);
		font-size: 0.72rem;
		font-weight: 700;
		color: var(--color-accent);
		background: var(--color-accent-subtle);
		padding: 4px 10px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(20, 184, 166, 0.15);
	}

	.sim-title {
		font-size: 1.15rem;
		font-weight: 700;
		color: var(--color-text-primary);
	}

	.sim-controls {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
	}

	.sim-btn {
		font-family: var(--font-sans);
		font-size: 0.8rem;
		font-weight: 600;
		padding: 8px 16px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border);
		background: transparent;
		color: var(--color-text-secondary);
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.sim-btn:hover {
		border-color: var(--color-border-active);
		color: var(--color-text-primary);
	}

	.btn-success.active {
		background: var(--color-success-subtle);
		color: var(--color-success);
		border-color: rgba(34, 197, 94, 0.3);
	}

	.btn-warning.active {
		background: var(--color-warning-subtle);
		color: var(--color-warning);
		border-color: rgba(245, 158, 11, 0.3);
	}

	.btn-danger.active {
		background: var(--color-danger-subtle);
		color: var(--color-danger);
		border-color: rgba(239, 68, 68, 0.3);
	}

	/* Diagram Area */
	.diagram-area {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 20px 0;
		overflow-x: auto;
		gap: 8px;
	}

	.node-box {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8px;
		padding: 18px 20px;
		min-width: 130px;
		border-radius: var(--radius-lg);
		background: var(--color-bg-elevated);
		border: 1px solid var(--color-border);
		transition: all 0.3s ease;
		text-align: center;
	}

	.node-box.active {
		border-color: rgba(20, 184, 166, 0.4);
	}

	.node-box.node-success {
		border-color: rgba(34, 197, 94, 0.4);
	}

	.node-box.node-warning {
		border-color: rgba(245, 158, 11, 0.4);
	}

	.node-box.node-danger {
		border-color: rgba(239, 68, 68, 0.4);
	}

	.node-icon {
		color: var(--color-text-secondary);
	}

	.node-box.active .node-icon {
		color: var(--color-accent);
	}

	.node-success .node-icon { color: var(--color-success); }
	.node-warning .node-icon { color: var(--color-warning); }
	.node-danger .node-icon { color: var(--color-danger); }

	.node-label {
		font-size: 0.85rem;
		font-weight: 700;
		color: var(--color-text-primary);
	}

	.node-sub {
		font-family: var(--font-mono);
		font-size: 0.68rem;
		color: var(--color-text-muted);
	}

	/* Connectors */
	.connector {
		flex: 1;
		height: 2px;
		background: var(--color-border);
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 50px;
		transition: background 0.3s ease;
	}

	.connector.active {
		background: rgba(20, 184, 166, 0.3);
	}

	.connector-label {
		position: absolute;
		top: -18px;
		font-family: var(--font-mono);
		font-size: 0.68rem;
		color: var(--color-text-muted);
	}

	.packet {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		background: var(--color-accent);
		position: absolute;
		opacity: 0;
	}

	.packet.moving {
		animation: travel 0.9s ease-in-out infinite;
	}

	.packet-warning {
		background: var(--color-warning);
	}

	.packet-danger {
		background: var(--color-danger);
	}

	@keyframes travel {
		0% {
			left: 0;
			opacity: 0;
		}
		20% {
			opacity: 1;
		}
		80% {
			opacity: 1;
		}
		100% {
			left: 100%;
			opacity: 0;
		}
	}

	/* Status Bar */
	.sim-status-bar {
		display: flex;
		align-items: center;
		padding: 12px 16px;
		background: var(--color-bg-elevated);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
	}

	.status-bar-text {
		font-family: var(--font-mono);
		font-size: 0.8rem;
		color: var(--color-text-secondary);
	}
</style>
