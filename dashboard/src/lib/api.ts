import type { WebhookListResponse } from './types';

const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

export async function fetchWebhooks(limit = 100, offset = 0): Promise<WebhookListResponse> {
	const res = await fetch(`${API_BASE}/api/v1/webhooks?limit=${limit}&offset=${offset}`);
	if (!res.ok) throw new Error(`Failed to fetch webhooks: ${res.status}`);
	return res.json();
}

export async function replayWebhook(id: string): Promise<{ id: string; status: string }> {
	const res = await fetch(`${API_BASE}/api/v1/webhook/${id}/replay`, { method: 'POST' });
	if (!res.ok) throw new Error(`Failed to replay webhook: ${res.status}`);
	return res.json();
}

export async function sendWebhook(
	source: string,
	payload: unknown,
	targetUrl?: string
): Promise<{ id: string; status: string }> {
	const url = new URL(`${API_BASE}/api/v1/webhook`, window.location.origin);
	if (targetUrl) {
		url.searchParams.set('target_url', targetUrl);
	}

	const res = await fetch(url.toString(), {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'X-Webhook-Source': source
		},
		body: JSON.stringify(payload)
	});

	if (!res.ok) {
		const errText = await res.text().catch(() => 'Unknown error');
		throw new Error(`Failed to send webhook (${res.status}): ${errText}`);
	}

	return res.json();
}

export function getSSEUrl(): string {
	return `${API_BASE}/api/v1/events`;
}
