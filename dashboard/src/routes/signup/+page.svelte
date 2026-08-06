<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import ShinyText from '$lib/components/landing/ShinyText.svelte';
	import BorderGlowButton from '$lib/components/landing/BorderGlowButton.svelte';
	import SpotlightCard from '$lib/components/svelte-bits/SpotlightCard.svelte';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let isLoading = $state(false);

	onMount(() => {
		const existingKey = localStorage.getItem('siphon_jwt_token');
		if (existingKey) {
			goto('/dashboard');
		}
	});

	async function handleSignup() {
		if (!email.trim() || !password.trim()) {
			error = 'Email and password are required';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';
			const res = await fetch(`${API_BASE}/api/v1/auth/signup`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});

			if (res.ok) {
				const data = await res.json();
				localStorage.setItem('siphon_jwt_token', data.token);
				goto('/dashboard');
			} else {
				const errData = await res.json().catch(() => ({}));
				error = errData.error || 'Failed to create account';
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
			<p>Create an account to manage your Webhook Infrastructure.</p>
		</div>

		<form onsubmit={(e) => { e.preventDefault(); handleSignup(); }} class="login-form">
			<div class="input-group">
				<input
					type="email"
					bind:value={email}
					placeholder="Email Address"
					class="api-input"
				/>
				<input
					type="password"
					bind:value={password}
					placeholder="Password"
					class="api-input"
				/>
			</div>

			{#if error}
				<div class="error-msg">{error}</div>
			{/if}

			<div class="submit-btn" role="button" tabindex="0" onclick={handleSignup} onkeydown={(e) => { if (e.key === 'Enter') handleSignup(); }}>
				<BorderGlowButton>
					{isLoading ? 'Creating Account...' : 'Sign Up'}
				</BorderGlowButton>
			</div>

			<div class="signup-link">
				Already have an account? <a href="/login">Log in</a>
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

	.signup-link {
		text-align: center;
		font-size: 0.9rem;
		color: var(--color-text-secondary);
		margin-top: 8px;
	}

	.signup-link a {
		color: var(--color-accent);
		text-decoration: none;
		font-weight: 500;
	}

	.signup-link a:hover {
		text-decoration: underline;
	}
</style>
