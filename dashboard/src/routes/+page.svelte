<script lang="ts">
	import WebhookSimulator from "$lib/components/landing/WebhookSimulator.svelte";

	let copied = $state(false);

	function copyCommand() {
		navigator.clipboard.writeText(
			"curl -X POST https://siphon.io/api/v1/ingest/YOUR_ENDPOINT_ID",
		);
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 2000);
	}

	const features = [
		{
			title: "Zero-Loss Webhook Ingestion",
			badge: "Instant 200 OK",
			description: "Never drop a third-party webhook again. Siphon absorbs traffic spikes, immediately queues payloads into durable persistent storage, and returns an instant HTTP 200 OK to providers like Stripe or GitHub.",
			icon: "ingest"
		},
		{
			title: "Exponential Backoff Retries",
			badge: "Resilience",
			description: "Downstream server offline? Automated retries back off exponentially (2s, 4s, 8s, 16s, 32s) without overwhelming your API. Configurable retry limits per endpoint.",
			icon: "retry"
		},
		{
			title: "Dead Letter Queue",
			badge: "Isolation",
			description: "Webhooks that exceed max retry attempts are automatically isolated in the DLQ. Inspect payloads and re-inject them with 1-click replay from the dashboard.",
			icon: "dlq"
		},
		{
			title: "Real-time SSE Observability",
			badge: "Live Stream",
			description: "Watch webhooks arrive, retry, and succeed in real time. Live event updates over Server-Sent Events (SSE) stream directly to your authenticated dashboard with zero polling overhead.",
			icon: "sse"
		},
		{
			title: "Full Delivery Logs",
			badge: "Transparency",
			description: "Stop guessing why a webhook failed. Every delivery attempt is logged with the exact HTTP response body, status code, and timestamp for instant debugging.",
			icon: "logs"
		},
		{
			title: "HMAC Signatures & Auth",
			badge: "Security",
			description: "Downstream servers can cryptographically verify payloads using SHA-256 HMAC signatures to prevent spoofing. The ingest API is protected by rate limiters and API keys.",
			icon: "security"
		}
	];

	const faqs = [
		{
			q: "What is Siphon Gateway?",
			a: "Siphon Gateway is a fault-tolerant webhook ingestion engine that acts as a buffer between third-party providers (Stripe, GitHub, Shopify) and your internal servers. It prevents data loss by queuing payloads asynchronously and delivering them with automatic retries."
		},
		{
			q: "How does the retry mechanism work?",
			a: "When a downstream delivery fails, Siphon re-queues the webhook with exponential backoff delays (2s, 4s, 8s, 16s, 32s). Each attempt is logged with the full HTTP response for debugging. After 5 failed attempts, the webhook is moved to the Dead Letter Queue."
		},
		{
			q: "What happens when all retries are exhausted?",
			a: "The webhook is moved to the Dead Letter Queue (DLQ) where it is preserved indefinitely. From the dashboard, you can inspect the original payload, review every failed delivery attempt, and replay the webhook with a single click once the issue is resolved."
		},
		{
			q: "Is my webhook data secure?",
			a: "Yes. All payloads are delivered with SHA-256 HMAC signatures so your downstream servers can verify authenticity. The ingestion API is protected by API key authentication and rate limiting. Data is stored in PostgreSQL with standard encryption at rest."
		},
		{
			q: "How do I integrate with Stripe, GitHub, or Shopify?",
			a: "Point your provider's webhook URL to your Siphon endpoint (e.g. https://siphon.io/api/v1/ingest/YOUR_ENDPOINT_ID). Siphon accepts the webhook, immediately returns 200 OK to the provider, and asynchronously delivers the payload to your configured target URL."
		}
	];

	let openFaq = $state<number | null>(null);

	function toggleFaq(index: number) {
		openFaq = openFaq === index ? null : index;
	}
</script>

<svelte:head>
	<title>Siphon Gateway — Fault-Tolerant Webhook Infrastructure</title>
	<meta name="description" content="Fault-tolerant webhook ingestion engine with async queuing, exponential backoff retries, and real-time observability." />
</svelte:head>

<div class="landing-page">
	<!-- Hero Section -->
	<section class="hero-section">
		<div class="hero-content">
			<span class="pill-badge animate-fade-in-up">
				<span class="badge-dot"></span>
				Enterprise-Grade Ingestion Engine
			</span>

			<h1 class="hero-title animate-fade-in-up" style="animation-delay: 80ms">
				Fault-Tolerant Webhook Infrastructure
			</h1>

			<p class="hero-subtitle animate-fade-in-up" style="animation-delay: 160ms">
				A resilient buffer between third-party webhook providers
				(Stripe, GitHub, Shopify) and your downstream servers.
				Asynchronous queuing, exponential backoff retries, and granular
				cryptographic verification.
			</p>

			<div class="hero-actions animate-fade-in-up" style="animation-delay: 240ms">
				<a href="/login" class="btn-primary-lg">
					Log In to Dashboard
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
						<line x1="5" y1="12" x2="19" y2="12"></line>
						<polyline points="12 5 19 12 12 19"></polyline>
					</svg>
				</a>

				<button
					class="cli-command-pill"
					onclick={copyCommand}
					title="Copy Example API Call"
				>
					<span class="cli-prompt">$</span>
					<code class="cli-code">curl -X POST /api/v1/ingest/...</code>
					<span class="cli-copy">
						{#if copied}
							<span class="copied-text">Copied</span>
						{:else}
							<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
								<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
							</svg>
						{/if}
					</span>
				</button>
			</div>
		</div>
	</section>

	<!-- Interactive Demo Section -->
	<section class="section-container animate-fade-in-up" style="animation-delay: 320ms">
		<div class="section-heading">
			<h2 class="section-title">Try the Architecture</h2>
			<p class="section-subtitle">
				Simulate real webhook payloads traveling through the ingestion
				gateway and automated retry workers.
			</p>
		</div>

		<WebhookSimulator />
	</section>

	<!-- Features Section -->
	<section class="section-container">
		<div class="section-heading">
			<h2 class="section-title">Built for Resilience</h2>
			<p class="section-subtitle">
				Every component is optimized for fast ingestion, guaranteed
				message durability, and developer observability.
			</p>
		</div>

		<div class="features-grid">
			{#each features as feature}
				<div class="feature-card">
					<div class="feature-header">
						<span class="feature-badge">{feature.badge}</span>
					</div>
					<h3 class="feature-title">{feature.title}</h3>
					<p class="feature-desc">{feature.description}</p>
				</div>
			{/each}
		</div>
	</section>

	<!-- FAQ Section -->
	<section class="section-container">
		<div class="section-heading">
			<h2 class="section-title">Frequently Asked Questions</h2>
		</div>

		<div class="faq-list">
			{#each faqs as faq, i}
				<button class="faq-item" class:open={openFaq === i} onclick={() => toggleFaq(i)}>
					<div class="faq-question">
						<span>{faq.q}</span>
						<svg class="faq-chevron" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
							<polyline points="6 9 12 15 18 9"></polyline>
						</svg>
					</div>
					{#if openFaq === i}
						<div class="faq-answer">
							<p>{faq.a}</p>
						</div>
					{/if}
				</button>
			{/each}
		</div>
	</section>

	<!-- Footer -->
	<footer class="landing-footer">
		<div class="footer-content">
			<span class="footer-brand">Siphon Gateway</span>
			<div class="footer-links">
				<a href="/privacy">Privacy Policy</a>
			</div>
			<span class="footer-meta">Juan Abreu Navarro &copy; 2026</span>
		</div>
	</footer>
</div>

<style>
	.landing-page {
		display: flex;
		flex-direction: column;
		gap: 80px;
		padding-bottom: 48px;
	}

	/* Hero Section */
	.hero-section {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 65vh;
		padding: 64px 20px;
	}

	.hero-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		max-width: 780px;
		gap: 24px;
	}

	.pill-badge {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		padding: 6px 16px;
		border-radius: var(--radius-md);
		background: var(--color-accent-subtle);
		border: 1px solid rgba(20, 184, 166, 0.2);
		color: var(--color-accent);
		font-size: 0.8rem;
		font-weight: 600;
		letter-spacing: 0.02em;
	}

	.badge-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-accent);
	}

	.hero-title {
		font-size: clamp(2rem, 5vw, 3.2rem);
		font-weight: 700;
		line-height: 1.15;
		letter-spacing: -0.03em;
		color: var(--color-text-primary);
	}

	.hero-subtitle {
		font-size: clamp(1rem, 1.8vw, 1.15rem);
		color: var(--color-text-secondary);
		line-height: 1.7;
		max-width: 650px;
	}

	.hero-actions {
		display: flex;
		align-items: center;
		gap: 16px;
		flex-wrap: wrap;
		justify-content: center;
		margin-top: 8px;
	}

	.btn-primary-lg {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		padding: 12px 24px;
		background: var(--color-accent);
		color: #ffffff;
		font-family: var(--font-sans);
		font-size: 0.95rem;
		font-weight: 600;
		border-radius: var(--radius-lg);
		text-decoration: none;
		transition: background var(--transition-fast);
	}

	.btn-primary-lg:hover {
		background: #0d9488;
	}

	/* CLI Pill */
	.cli-command-pill {
		display: inline-flex;
		align-items: center;
		gap: 10px;
		padding: 12px 20px;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text-primary);
		font-family: var(--font-mono);
		font-size: 0.85rem;
		cursor: pointer;
		transition: border-color var(--transition-fast);
	}

	.cli-command-pill:hover {
		border-color: var(--color-border-active);
	}

	.cli-prompt {
		color: var(--color-accent);
		font-weight: 700;
	}

	.cli-code {
		color: var(--color-text-primary);
	}

	.cli-copy {
		display: flex;
		align-items: center;
		color: var(--color-text-muted);
		margin-left: 4px;
	}

	.copied-text {
		color: var(--color-success);
		font-size: 0.75rem;
		font-weight: 600;
	}

	/* Section Containers */
	.section-container {
		display: flex;
		flex-direction: column;
		gap: 32px;
	}

	.section-heading {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: 10px;
	}

	.section-title {
		font-size: 1.8rem;
		font-weight: 700;
		color: var(--color-text-primary);
		letter-spacing: -0.02em;
	}

	.section-subtitle {
		font-size: 1rem;
		color: var(--color-text-secondary);
		max-width: 560px;
		line-height: 1.6;
	}

	/* Features Grid */
	.features-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 16px;
		width: 100%;
	}

	@media (max-width: 700px) {
		.features-grid {
			grid-template-columns: 1fr;
		}
	}

	.feature-card {
		display: flex;
		flex-direction: column;
		gap: 12px;
		padding: 24px;
		border-radius: var(--radius-xl);
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		transition: border-color var(--transition-fast);
	}

	.feature-card:hover {
		border-color: var(--color-border-active);
	}

	.feature-header {
		display: flex;
		align-items: center;
	}

	.feature-badge {
		font-size: 0.7rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.06em;
		padding: 3px 8px;
		border-radius: var(--radius-sm);
		background: var(--color-accent-subtle);
		color: var(--color-accent);
		border: 1px solid rgba(20, 184, 166, 0.15);
	}

	.feature-title {
		font-size: 1.1rem;
		font-weight: 600;
		color: var(--color-text-primary);
		letter-spacing: -0.01em;
	}

	.feature-desc {
		font-size: 0.9rem;
		color: var(--color-text-secondary);
		line-height: 1.6;
	}

	/* FAQ */
	.faq-list {
		display: flex;
		flex-direction: column;
		gap: 2px;
		max-width: 720px;
		margin: 0 auto;
		width: 100%;
	}

	.faq-item {
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		padding: 0;
		cursor: pointer;
		text-align: left;
		color: inherit;
		font-family: inherit;
		transition: border-color var(--transition-fast);
	}

	.faq-item:first-child {
		border-radius: var(--radius-lg) var(--radius-lg) 0 0;
	}

	.faq-item:last-child {
		border-radius: 0 0 var(--radius-lg) var(--radius-lg);
	}

	.faq-item:hover {
		border-color: var(--color-border-active);
	}

	.faq-question {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 16px 20px;
		font-size: 0.95rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	.faq-chevron {
		flex-shrink: 0;
		color: var(--color-text-muted);
		transition: transform var(--transition-fast);
	}

	.faq-item.open .faq-chevron {
		transform: rotate(180deg);
	}

	.faq-answer {
		padding: 0 20px 16px;
	}

	.faq-answer p {
		font-size: 0.9rem;
		color: var(--color-text-secondary);
		line-height: 1.65;
	}

	/* Footer */
	.landing-footer {
		margin-top: 16px;
		padding-top: 28px;
		border-top: 1px solid var(--color-border);
	}

	.footer-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-wrap: wrap;
		gap: 16px;
		font-size: 0.82rem;
		color: var(--color-text-muted);
	}

	.footer-brand {
		font-weight: 700;
		color: var(--color-text-secondary);
	}

	.footer-links a {
		color: var(--color-text-secondary);
		text-decoration: none;
		transition: color var(--transition-fast);
	}

	.footer-links a:hover {
		color: var(--color-text-primary);
	}
</style>
