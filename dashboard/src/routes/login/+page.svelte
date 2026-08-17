<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

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

	async function handleLogin() {
		if (!email.trim() || !password.trim()) {
			error = 'Email and password are required';
			return;
		}

		isLoading = true;
		error = '';

		try {
			const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';
			const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});

			if (res.ok) {
				const data = await res.json();
				localStorage.setItem('siphon_jwt_token', data.token);
				goto('/dashboard');
			} else {
				error = 'Invalid credentials';
			}
		} catch (err) {
			error = 'Failed to connect to the server';
		} finally {
			isLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Log In — Siphon Gateway</title>
</svelte:head>

<div class="auth-container">
	<div class="auth-card">
		<div class="auth-header">
			<h2>Siphon Gateway</h2>
			<p>Log in to manage your webhook infrastructure.</p>
		</div>

		<form onsubmit={(e) => { e.preventDefault(); handleLogin(); }} class="auth-form">
			<div class="input-group">
				<input
					type="email"
					bind:value={email}
					placeholder="Email Address"
					class="auth-input"
				/>
				<input
					type="password"
					bind:value={password}
					placeholder="Password"
					class="auth-input"
				/>
			</div>

			{#if error}
				<div class="error-msg">{error}</div>
			{/if}

			<button type="submit" class="btn-submit" disabled={isLoading}>
				{isLoading ? 'Authenticating...' : 'Sign In'}
			</button>

			<div class="auth-link">
				Don't have an account? <a href="/signup">Sign up</a>
			</div>
		</form>
	</div>
</div>

<style>
	.auth-container {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 80vh;
		width: 100%;
	}

	.auth-card {
		width: 100%;
		max-width: 400px;
		background: var(--color-bg-card);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		padding: 40px;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.auth-header h2 {
		font-size: 1.6rem;
		margin: 0 0 8px 0;
		font-weight: 700;
		text-align: center;
		color: var(--color-text-primary);
	}

	.auth-header p {
		color: var(--color-text-secondary);
		font-size: 0.9rem;
		text-align: center;
		margin: 0;
	}

	.auth-form {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.input-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.auth-input {
		background: var(--color-bg-elevated);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		padding: 12px 16px;
		color: var(--color-text-primary);
		font-family: var(--font-sans);
		font-size: 0.9rem;
		transition: border-color var(--transition-fast);
		width: 100%;
		box-sizing: border-box;
	}

	.auth-input:focus {
		outline: none;
		border-color: var(--color-accent);
	}

	.error-msg {
		color: var(--color-danger);
		font-size: 0.85rem;
		background: var(--color-danger-subtle);
		padding: 8px 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(239, 68, 68, 0.2);
		text-align: center;
	}

	.btn-submit {
		width: 100%;
		padding: 12px;
		border-radius: var(--radius-md);
		border: none;
		background: var(--color-accent);
		color: #ffffff;
		font-family: var(--font-sans);
		font-weight: 600;
		font-size: 0.95rem;
		cursor: pointer;
		transition: background var(--transition-fast);
		margin-top: 8px;
	}

	.btn-submit:hover:not(:disabled) {
		background: #0d9488;
	}

	.btn-submit:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.auth-link {
		text-align: center;
		font-size: 0.9rem;
		color: var(--color-text-secondary);
	}

	.auth-link a {
		color: var(--color-accent);
		text-decoration: none;
		font-weight: 500;
	}

	.auth-link a:hover {
		text-decoration: underline;
	}
</style>
