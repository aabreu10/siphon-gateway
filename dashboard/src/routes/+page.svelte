<script lang="ts">
	import Balatro from '$lib/components/svelte-bits/Balatro.svelte';
	import ShinyText from '$lib/components/landing/ShinyText.svelte';
	import BorderGlowButton from '$lib/components/landing/BorderGlowButton.svelte';
	import BentoGrid from '$lib/components/landing/BentoGrid.svelte';
	import WebhookSimulator from '$lib/components/landing/WebhookSimulator.svelte';

	let copied = $state(false);

	let copied = $state(false);

	function copyCommand() {
		navigator.clipboard.writeText('curl -X POST https://siphon.io/api/v1/ingest/YOUR_ENDPOINT_ID');
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 2000);
	}
</script>

<div class="landing-page">
	<!-- Hero Section -->
	<section class="hero-section">
		<Balatro color1="#6c5ce7" color2="#5f3dc4" color3="#0a0b0f" mouseInteraction={true} />

		<div class="hero-content">
			<div class="pill-badge animate-fade-in-up">
				<span class="pulse-dot"></span>
				<span>Enterprise-Grade Ingestion Engine</span>
			</div>

			<h1 class="hero-title animate-fade-in-up" style="animation-delay: 100ms">
				<ShinyText text="Fault-Tolerant Webhook Infrastructure" speed={5} />
			</h1>

			<p class="hero-subtitle animate-fade-in-up" style="animation-delay: 200ms">
				A resilient shock absorber between third-party webhook providers (Stripe, GitHub, Shopify) and your downstream internal servers. Prevent data loss with asynchronous queuing, exponential backoff retries, and granular cryptographic security.
			</p>

			<div class="hero-actions animate-fade-in-up" style="animation-delay: 300ms">
				<BorderGlowButton href="/login">
					Log In to Dashboard
					<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
						<line x1="5" y1="12" x2="19" y2="12"></line>
						<polyline points="12 5 19 12 12 19"></polyline>
					</svg>
				</BorderGlowButton>

				<button class="cli-command-pill" onclick={copyCommand} title="Copy Example API Call">
					<span class="cli-prompt">$</span>
					<code class="cli-code">curl -X POST /api/v1/ingest/...</code>
					<span class="cli-copy">
						{#if copied}
							<span class="copied-text">✓ Copied</span>
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

	<!-- Interactive Sandbox Section -->
	<section class="section-container animate-fade-in-up" style="animation-delay: 400ms">
		<div class="section-heading">
			<h2 class="section-title">Test The Architecture Live</h2>
			<p class="section-subtitle">
				Click below to simulate real webhook payloads traveling through our highly available ingestion gateways and automated retry workers.
			</p>
		</div>

		<WebhookSimulator />
	</section>

	<!-- Bento Grid Feature Highlights -->
	<section class="section-container">
		<div class="section-heading">
			<h2 class="section-title">Engineered For Resilience</h2>
			<p class="section-subtitle">
				Every component is optimized for sub-millisecond ingestion, guaranteed message durability, and developer observability.
			</p>
		</div>

		<BentoGrid />
	</section>

	<!-- Footer -->
	<footer class="landing-footer">
		<div class="footer-content">
			<span class="footer-brand">Siphon Gateway SaaS</span>
			<span class="footer-meta">Enterprise Webhook Delivery Network &copy; 2026</span>
		</div>
	</footer>
</div>

<style>
	.landing-page {
		display: flex;
		flex-direction: column;
		gap: 96px;
		padding-bottom: 48px;
	}

	/* Hero Section */
	.hero-section {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 70vh;
		padding: 48px 20px;
		border-radius: var(--radius-xl);
		background: radial-gradient(circle at 50% 20%, rgba(108, 92, 231, 0.12), transparent 70%);
		border: 1px solid rgba(255, 255, 255, 0.05);
		overflow: hidden;
	}

	.hero-content {
		position: relative;
		z-index: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		max-width: 860px;
		gap: 24px;
	}

	.pill-badge {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		padding: 6px 16px;
		border-radius: 24px;
		background: rgba(108, 92, 231, 0.12);
		border: 1px solid rgba(108, 92, 231, 0.3);
		color: var(--color-accent);
		font-size: 0.8rem;
		font-weight: 600;
		letter-spacing: 0.02em;
	}

	.pulse-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--color-accent);
		box-shadow: 0 0 8px var(--color-accent);
		animation: pulse-glow 2s ease-in-out infinite;
	}

	.hero-title {
		font-size: clamp(2.2rem, 5vw, 3.8rem);
		font-weight: 800;
		line-height: 1.15;
		letter-spacing: -0.03em;
		max-width: 780px;
	}

	.hero-subtitle {
		font-size: clamp(1rem, 1.8vw, 1.25rem);
		color: var(--color-text-secondary);
		line-height: 1.7;
		max-width: 700px;
	}

	.hero-actions {
		display: flex;
		align-items: center;
		gap: 16px;
		flex-wrap: wrap;
		justify-content: center;
		margin-top: 8px;
	}

	/* CLI Pill */
	.cli-command-pill {
		display: inline-flex;
		align-items: center;
		gap: 10px;
		padding: 12px 20px;
		background: rgba(0, 0, 0, 0.5);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text-primary);
		font-family: var(--font-mono);
		font-size: 0.85rem;
		cursor: pointer;
		transition: all var(--transition-fast);
	}

	.cli-command-pill:hover {
		border-color: var(--color-border-active);
		background: rgba(0, 0, 0, 0.7);
		transform: translateY(-2px);
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
		gap: 36px;
	}

	.section-heading {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: 12px;
	}

	.section-title {
		font-size: 2rem;
		font-weight: 700;
		color: var(--color-text-primary);
		letter-spacing: -0.02em;
	}

	.section-subtitle {
		font-size: 1rem;
		color: var(--color-text-secondary);
		max-width: 620px;
		line-height: 1.6;
	}

	/* Footer */
	.landing-footer {
		margin-top: 24px;
		padding-top: 32px;
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
</style>
