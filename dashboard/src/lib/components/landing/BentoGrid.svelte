<script lang="ts">
	import SpotlightCard from '$lib/components/svelte-bits/SpotlightCard.svelte';
</script>

<div class="bento-grid">
	<!-- Card 1: Zero-Loss Ingestion (Large span) -->
	<SpotlightCard class="bento-card span-2 glass-panel hover-glow">
		<div class="card-header">
			<span class="card-badge badge-purple">Instant 200 OK</span>
			<div class="card-icon">
				<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path>
				</svg>
			</div>
		</div>
		<h3 class="card-title">Zero-Loss Webhook Ingestion</h3>
		<p class="card-desc">
			Never drop a third-party webhook again. Siphon Gateway absorbs traffic spikes, immediately queues payloads into durable RabbitMQ exchanges, and returns an instant HTTP 200 OK to providers like Stripe or GitHub.
		</p>
		<div class="card-visual visual-ingest">
			<div class="flow-pill">POST /webhook</div>
			<span class="flow-arrow">→</span>
			<div class="flow-pill pulse-green">200 OK</div>
			<span class="flow-arrow">→</span>
			<div class="flow-pill rabbit-pill">RabbitMQ</div>
		</div>
	</SpotlightCard>

	<!-- Card 2: Exponential Backoff Retries -->
	<SpotlightCard class="bento-card glass-panel hover-glow">
		<div class="card-header">
			<span class="card-badge badge-amber">Resilience</span>
			<div class="card-icon">
				<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<polyline points="23 4 23 10 17 10"></polyline>
					<path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
				</svg>
			</div>
		</div>
		<h3 class="card-title">Exponential Backoff</h3>
		<p class="card-desc">
			Downstream server offline? Automated retries back off exponentially (2s → 4s → 8s → 16s → 32s) without overwhelming your API.
		</p>
		<div class="timeline-bars">
			<div class="bar" style="--delay: 2s; --width: 20%;">2s</div>
			<div class="bar" style="--delay: 4s; --width: 40%;">4s</div>
			<div class="bar" style="--delay: 8s; --width: 65%;">8s</div>
			<div class="bar" style="--delay: 16s; --width: 90%;">16s</div>
		</div>
	</SpotlightCard>

	<!-- Card 3: Dead Letter Queue (DLQ) -->
	<SpotlightCard class="bento-card glass-panel hover-glow">
		<div class="card-header">
			<span class="card-badge badge-red">Isolation</span>
			<div class="card-icon">
				<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
					<line x1="12" y1="9" x2="12" y2="13"></line>
					<line x1="12" y1="17" x2="12.01" y2="17"></line>
				</svg>
			</div>
		</div>
		<h3 class="card-title">Dead Letter Queue</h3>
		<p class="card-desc">
			Webhooks that exceed max retry attempts are automatically isolated in the DLQ. Inspect payloads and re-inject them with 1-click replay.
		</p>
		<div class="dlq-badge-visual">
			<span class="dlq-tag">FAILED_DLQ</span>
			<span class="replay-action">↻ Replay Ready</span>
		</div>
	</SpotlightCard>

	<!-- Card 4: Real-time Observability (Large span) -->
	<SpotlightCard class="bento-card span-2 glass-panel hover-glow">
		<div class="card-header">
			<span class="card-badge badge-emerald">Live Stream</span>
			<div class="card-icon">
				<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
				</svg>
			</div>
		</div>
		<h3 class="card-title">Real-time SSE Observability</h3>
		<p class="card-desc">
			Watch webhooks arrive, retry, and succeed in real time. Our Go backend broadcasts live event updates over HTTP Server-Sent Events (SSE) directly to your SvelteKit dashboard with zero polling overhead.
		</p>
		<div class="sse-visual">
			<div class="sse-node go-node">Go Gateway</div>
			<div class="sse-stream">
				<span class="stream-dot"></span>
				<span class="stream-dot"></span>
				<span class="stream-dot"></span>
				<span class="stream-label">SSE /api/v1/events</span>
			</div>
			<div class="sse-node svelte-node">SvelteKit UI</div>
		</div>
	</SpotlightCard>

	<!-- Card 5: Extreme Visibility & Delivery Logs -->
	<SpotlightCard class="bento-card glass-panel hover-glow">
		<div class="card-header">
			<span class="card-badge badge-blue">Transparency</span>
			<div class="card-icon">
				<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
					<polyline points="14 2 14 8 20 8"></polyline>
					<line x1="16" y1="13" x2="8" y2="13"></line>
					<line x1="16" y1="17" x2="8" y2="17"></line>
					<polyline points="10 9 9 9 8 9"></polyline>
				</svg>
			</div>
		</div>
		<h3 class="card-title">Delivery Logs</h3>
		<p class="card-desc">
			Stop guessing why a webhook failed. We store the exact HTTP response body, status code, and timestamp for every single delivery attempt so you can debug downstream errors instantly.
		</p>
		<div class="log-visual">
			<div class="log-line">HTTP 500</div>
			<div class="log-line">"error": "database lock"</div>
		</div>
	</SpotlightCard>

	<!-- Card 6: Enterprise Security -->
	<SpotlightCard class="bento-card glass-panel hover-glow">
		<div class="card-header">
			<span class="card-badge badge-gray">Security</span>
			<div class="card-icon">
				<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
					<path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
				</svg>
			</div>
		</div>
		<h3 class="card-title">HMAC Signatures & Auth</h3>
		<p class="card-desc">
			Downstream servers can cryptographically verify payloads using SHA-256 HMAC signatures to prevent spoofing. The ingest API is protected by rate limiters and API keys.
		</p>
		<div class="hmac-visual">
			<span class="header-name">X-Siphon-Signature:</span>
			<span class="header-value">e3b0c44298fc1c14...</span>
		</div>
	</SpotlightCard>
</div>

<style>
	.bento-grid {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 20px;
		width: 100%;
	}

	@media (max-width: 960px) {
		.bento-grid {
			grid-template-columns: repeat(2, 1fr);
		}
	}

	@media (max-width: 640px) {
		.bento-grid {
			grid-template-columns: 1fr;
		}
	}

	:global(.bento-card) {
		border-radius: var(--radius-xl);
		transition: transform 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease;
		background: rgba(18, 19, 26, 0.6);
	}

	:global(.bento-card .spotlight-content) {
		padding: 30px;
		display: flex;
		flex-direction: column;
		gap: 14px;
		width: 100%;
		height: 100%;
		box-sizing: border-box;
	}

	:global(.bento-card.span-2) {
		grid-column: span 2;
	}

	@media (max-width: 960px) {
		:global(.bento-card.span-2) {
			grid-column: span 1;
		}
	}

	:global(.bento-card:hover) {
		transform: translateY(-4px);
		border-color: rgba(108, 92, 231, 0.35);
		box-shadow:
			0 12px 40px rgba(0, 0, 0, 0.5),
			0 0 50px rgba(108, 92, 231, 0.1);
	}

	.card-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.card-badge {
		font-size: 0.72rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		padding: 4px 10px;
		border-radius: 20px;
	}

	.badge-purple {
		background: var(--color-accent-glow);
		color: var(--color-accent);
		border: 1px solid rgba(108, 92, 231, 0.2);
	}

	.badge-amber {
		background: var(--color-warning-glow);
		color: var(--color-warning);
		border: 1px solid rgba(255, 171, 64, 0.2);
	}

	.badge-red {
		background: var(--color-danger-glow);
		color: var(--color-danger);
		border: 1px solid rgba(255, 82, 82, 0.2);
	}

	.badge-emerald {
		background: var(--color-success-glow);
		color: var(--color-success);
		border: 1px solid rgba(0, 230, 118, 0.2);
	}

	.badge-blue {
		background: rgba(9, 132, 227, 0.1);
		color: #0984e3;
		border: 1px solid rgba(9, 132, 227, 0.2);
	}

	.badge-gray {
		background: rgba(178, 190, 195, 0.1);
		color: #b2bec3;
		border: 1px solid rgba(178, 190, 195, 0.2);
	}

	.card-icon {
		color: var(--color-text-secondary);
	}

	.card-title {
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--color-text-primary);
		letter-spacing: -0.02em;
	}

	.card-desc {
		font-size: 0.9rem;
		color: var(--color-text-secondary);
		line-height: 1.6;
	}

	/* Card Visuals */
	.card-visual {
		display: flex;
		align-items: center;
		gap: 10px;
		margin-top: auto;
		padding-top: 14px;
		flex-wrap: wrap;
	}

	.flow-pill {
		font-family: var(--font-mono);
		font-size: 0.78rem;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		background: rgba(255, 255, 255, 0.04);
		border: 1px solid var(--color-border);
		color: var(--color-text-primary);
	}

	.flow-arrow {
		color: var(--color-text-muted);
		font-weight: bold;
	}

	.pulse-green {
		color: var(--color-success);
		border-color: rgba(0, 230, 118, 0.3);
		background: rgba(0, 230, 118, 0.08);
	}

	.rabbit-pill {
		color: #ff9f43;
		border-color: rgba(255, 159, 67, 0.3);
		background: rgba(255, 159, 67, 0.08);
	}

	/* Timeline Bars */
	.timeline-bars {
		margin-top: auto;
		padding-top: 14px;
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.bar {
		font-family: var(--font-mono);
		font-size: 0.72rem;
		color: var(--color-text-secondary);
		background: linear-gradient(90deg, rgba(255, 171, 64, 0.2), rgba(255, 171, 64, 0.05));
		border-left: 2px solid var(--color-warning);
		width: var(--width);
		padding: 3px 8px;
		border-radius: 0 4px 4px 0;
	}

	/* DLQ Visual */
	.dlq-badge-visual {
		margin-top: auto;
		padding-top: 14px;
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.dlq-tag {
		font-family: var(--font-mono);
		font-size: 0.8rem;
		font-weight: 700;
		color: var(--color-danger);
		background: var(--color-danger-glow);
		padding: 5px 12px;
		border-radius: 4px;
		border: 1px solid rgba(255, 82, 82, 0.2);
	}

	.replay-action {
		font-size: 0.8rem;
		font-weight: 600;
		color: var(--color-accent);
	}

	/* SSE Visual */
	.sse-visual {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-top: auto;
		padding-top: 14px;
		gap: 16px;
	}

	.sse-node {
		font-family: var(--font-mono);
		font-size: 0.82rem;
		font-weight: 600;
		padding: 8px 16px;
		border-radius: var(--radius-sm);
		border: 1px solid var(--color-border);
	}

	.go-node {
		background: rgba(0, 173, 216, 0.1);
		color: #00add8;
		border-color: rgba(0, 173, 216, 0.3);
	}

	.svelte-node {
		background: rgba(255, 62, 0, 0.1);
		color: #ff3e00;
		border-color: rgba(255, 62, 0, 0.3);
	}

	.sse-stream {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
		position: relative;
		height: 2px;
		background: linear-gradient(90deg, rgba(0, 173, 216, 0.3), rgba(0, 230, 118, 0.5), rgba(255, 62, 0, 0.3));
	}

	.stream-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-success);
		box-shadow: 0 0 6px var(--color-success-glow);
		animation: pulse-glow 1.5s ease-in-out infinite;
	}

	.stream-label {
		position: absolute;
		top: -18px;
		font-family: var(--font-mono);
		font-size: 0.7rem;
		color: var(--color-text-muted);
		white-space: nowrap;
	}

	/* Log Visual */
	.log-visual {
		margin-top: auto;
		padding-top: 14px;
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.log-line {
		font-family: var(--font-mono);
		font-size: 0.72rem;
		color: var(--color-danger);
		background: var(--color-danger-glow);
		padding: 4px 8px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(255, 82, 82, 0.2);
	}

	/* HMAC Visual */
	.hmac-visual {
		margin-top: auto;
		padding-top: 14px;
		display: flex;
		align-items: center;
		gap: 8px;
		font-family: var(--font-mono);
		font-size: 0.75rem;
		background: rgba(255, 255, 255, 0.05);
		padding: 8px 12px;
		border-radius: var(--radius-sm);
		border: 1px solid var(--color-border);
	}

	.header-name {
		color: var(--color-accent);
		font-weight: 600;
	}

	.header-value {
		color: var(--color-text-secondary);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>
