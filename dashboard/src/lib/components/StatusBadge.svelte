<script lang="ts">
	import type { WebhookStatus } from '$lib/types';

	interface Props {
		status: WebhookStatus;
		retryCount?: number;
	}

	let { status, retryCount = 0 }: Props = $props();

	const config = $derived.by(() => {
		switch (status) {
			case 'SUCCESS':
				return { label: 'Delivered', cssClass: 'badge-success', icon: '✓' };
			case 'FAILED_DLQ':
				return { label: 'Failed (DLQ)', cssClass: 'badge-danger', icon: '✕' };
			case 'PENDING':
				if (retryCount > 0) {
					return { label: `Retry #${retryCount}`, cssClass: 'badge-warning', icon: '↻' };
				}
				return { label: 'Pending', cssClass: 'badge-pending', icon: '◉' };
			default:
				return { label: status, cssClass: 'badge-pending', icon: '?' };
		}
	});
</script>

<span class="badge {config.cssClass}" title="{config.label}">
	<span class="badge-icon">{config.icon}</span>
	<span class="badge-label">{config.label}</span>
</span>

<style>
	.badge {
		display: inline-flex;
		align-items: center;
		gap: 5px;
		padding: 4px 10px;
		border-radius: 20px;
		font-size: 0.75rem;
		font-weight: 600;
		letter-spacing: 0.02em;
		white-space: nowrap;
		line-height: 1;
		transition: all var(--transition-fast);
	}

	.badge-icon {
		font-size: 0.7rem;
	}

	.badge-success {
		color: var(--color-success);
		background: var(--color-success-subtle);
		border: 1px solid rgba(34, 197, 94, 0.15);
	}

	.badge-warning {
		color: var(--color-warning);
		background: var(--color-warning-subtle);
		border: 1px solid rgba(245, 158, 11, 0.15);
	}

	.badge-danger {
		color: var(--color-danger);
		background: var(--color-danger-subtle);
		border: 1px solid rgba(239, 68, 68, 0.15);
	}

	.badge-pending {
		color: var(--color-accent);
		background: var(--color-accent-subtle);
		border: 1px solid rgba(20, 184, 166, 0.15);
	}
</style>
