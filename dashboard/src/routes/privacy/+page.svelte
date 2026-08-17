<svelte:head>
	<title>Privacy Policy — Siphon Gateway</title>
	<meta name="description" content="Privacy policy for Siphon Gateway webhook ingestion service." />
</svelte:head>

<div class="privacy-page">
	<h1>Privacy Policy</h1>
	<p class="last-updated">Last updated: August 2026</p>

	<section>
		<h2>1. Information We Collect</h2>
		<p>
			When you create an account, we collect your email address and a hashed password. When you use the webhook ingestion service, we temporarily store webhook payload data that third-party providers send to your configured endpoints.
		</p>
	</section>

	<section>
		<h2>2. How We Use Your Data</h2>
		<p>
			Webhook payloads are stored solely for the purpose of reliable delivery to your configured downstream servers. Payloads in the Dead Letter Queue are retained until you manually replay or delete them. We do not analyze, sell, or share your webhook data with third parties.
		</p>
	</section>

	<section>
		<h2>3. Data Retention</h2>
		<p>
			Successfully delivered webhooks and their metadata are retained for 30 days for observability and debugging purposes. Dead Letter Queue entries are retained indefinitely until you take action. You may request deletion of your account and all associated data at any time.
		</p>
	</section>

	<section>
		<h2>4. Security</h2>
		<p>
			All webhook deliveries support SHA-256 HMAC signature verification. API access is protected by JWT authentication and rate limiting. Data is stored in PostgreSQL with encryption at rest. All connections use TLS.
		</p>
	</section>

	<section>
		<h2>5. Contact</h2>
		<p>
			For questions about this privacy policy or your data, contact the project maintainer through the repository at github.com/aabreu10/siphon-gateway.
		</p>
	</section>
</div>

<style>
	.privacy-page {
		max-width: 680px;
		margin: 0 auto;
		padding: 40px 0;
		display: flex;
		flex-direction: column;
		gap: 28px;
	}

	h1 {
		font-size: 2rem;
		font-weight: 700;
		color: var(--color-text-primary);
		letter-spacing: -0.02em;
	}

	.last-updated {
		font-size: 0.85rem;
		color: var(--color-text-muted);
		margin-top: -16px;
	}

	section {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	h2 {
		font-size: 1.15rem;
		font-weight: 600;
		color: var(--color-text-primary);
	}

	p {
		font-size: 0.95rem;
		color: var(--color-text-secondary);
		line-height: 1.7;
	}
</style>
