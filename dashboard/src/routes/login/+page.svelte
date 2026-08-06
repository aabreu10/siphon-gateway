<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import ShinyText from '$lib/components/svelte-bits/ShinyText.svelte';
	import BorderGlowButton from '$lib/components/svelte-bits/BorderGlowButton.svelte';
	import SpotlightCard from '$lib/components/svelte-bits/SpotlightCard.svelte';

	let apiKey = '';
	let error = '';
	let isLoading = false;

	onMount(() => {
		const existingKey = localStorage.getItem('siphon_admin_key');
		if (existingKey) {
			goto('/dashboard');
		}
	});

	async function handleLogin() {
		if (!apiKey.trim()) {
			error = 'API Key is required';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';
			// Verify key by fetching endpoints (requires admin key)
			const res = await fetch(`${API_BASE}/api/v1/endpoints`, {
				headers: { Authorization: `Bearer ${apiKey}` }
			});

			if (res.ok) {
				localStorage.setItem('siphon_admin_key', apiKey);
				goto('/dashboard');
			} else {
				error = 'Invalid API Key';
			}
		} catch (err) {
			error = 'Failed to connect to the server';
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="login-container">
	<SpotlightCard class="login-card glass-panel hover-glow">
		<div class="login-header">
			<h2><ShinyText text="Siphon Gateway" speed={3} /></h2>
			<p>Enter your Admin API Key to access the dashboard.</p>
		</div>

		<form on:submit|preventDefault={handleLogin} class="login-form">
			<div class="input-group">
				<input
					type="password"
					bind:value={apiKey}
					placeholder="sk_siphon_..."
					class="api-input"
				/>
			</div>

			{#if error}
				<div class="error-msg">{error}</div>
			{/if}

			<div class="submit-btn" on:click={handleLogin}>
				<BorderGlowButton>
					{isLoading ? 'Verifying...' : 'Sign In'}
				</BorderGlowButton>
			</div>
		</form>
	</SpotlightCard>
</div>

<style>
	.login-container {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 80vh;
		width: 100%;
	}

	:global(.login-card) {
		width: 100%;
		max-width: 400px;
		background: rgba(18, 19, 26, 0.6);
		border-radius: var(--radius-xl);
	}

	:global(.login-card .spotlight-content) {
		padding: 40px;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.login-header h2 {
		font-size: 1.8rem;
		margin: 0 0 8px 0;
		font-weight: 700;
		text-align: center;
	}

	.login-header p {
		color: var(--color-text-secondary);
		font-size: 0.9rem;
		text-align: center;
		margin: 0;
	}

	.login-form {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.input-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.api-input {
		background: rgba(0, 0, 0, 0.2);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: 12px 16px;
		color: var(--color-text-primary);
		font-family: var(--font-mono);
		font-size: 0.9rem;
		transition: all 0.2s ease;
		width: 100%;
		box-sizing: border-box;
	}

	.api-input:focus {
		outline: none;
		border-color: var(--color-accent);
		box-shadow: 0 0 0 2px rgba(108, 92, 231, 0.2);
	}

	.error-msg {
		color: var(--color-danger);
		font-size: 0.85rem;
		background: var(--color-danger-glow);
		padding: 8px 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(255, 82, 82, 0.2);
		text-align: center;
	}

	.submit-btn {
		width: 100%;
		display: flex;
		justify-content: center;
		margin-top: 10px;
	}
</style>
