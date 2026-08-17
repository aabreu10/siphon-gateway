import type { WebhookListResponse } from './types';

const API_BASE = (import.meta.env.VITE_API_URL ?? 'http://localhost:8080').replace(/\/+$/, '');

export function getAuthHeaders(): Record<string, string> {
	// We no longer read from localStorage for the API requests;
	// authentication is handled via HttpOnly cookies (credentials: 'include').
	// We return empty headers or fallback logic if needed.
	return {};
}

export async function fetchWebhooks(limit = 100, offset = 0, status = 'ALL', search = ''): Promise<WebhookListResponse> {
	const params = new URLSearchParams({
		limit: limit.toString(),
		offset: offset.toString(),
		status: status,
		search: search
	});
	const res = await fetch(`${API_BASE}/api/v1/webhooks?${params.toString()}`, {
		headers: getAuthHeaders(),
		credentials: 'include'
	});
	if (res.status === 401) throw new Error('Unauthorized');
	if (!res.ok) throw new Error(`Failed to fetch webhooks: ${res.status}`);
	return res.json();
}

export async function replayWebhook(id: string): Promise<{ id: string; status: string }> {
	const res = await fetch(`${API_BASE}/api/v1/webhook/${id}/replay`, { 
		method: 'POST',
		headers: getAuthHeaders(),
		credentials: 'include'
	});
	if (!res.ok) throw new Error(`Failed to replay webhook: ${res.status}`);
	return res.json();
}

export async function fetchDeliveryLogs(id: string): Promise<any[]> {
	const res = await fetch(`${API_BASE}/api/v1/webhook/${id}/logs`, {
		headers: getAuthHeaders(),
		credentials: 'include'
	});
	if (!res.ok) throw new Error(`Failed to fetch logs: ${res.status}`);
	return res.json();
}

// Endpoints CRUD
export async function fetchEndpoints(): Promise<any[]> {
	const res = await fetch(`${API_BASE}/api/v1/endpoints`, {
		headers: getAuthHeaders(),
		credentials: 'include'
	});
	if (!res.ok) throw new Error(`Failed to fetch endpoints: ${res.status}`);
	return res.json();
}

export async function createEndpoint(name: string, target_url: string, secret_key: string): Promise<any> {
	const res = await fetch(`${API_BASE}/api/v1/endpoints`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			...getAuthHeaders()
		},
		credentials: 'include',
		body: JSON.stringify({ name, target_url, secret_key })
	});
	if (!res.ok) throw new Error(`Failed to create endpoint: ${res.status}`);
	return res.json();
}

export async function sendWebhook(
	source: string,
	payload: unknown,
	endpointId: string
): Promise<{ id: string; status: string }> {
	let url = `${API_BASE}/api/v1/ingest`;
	if (endpointId) {
		const isUUID = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/.test(endpointId);
		if (isUUID) {
			url += `/${endpointId}`;
		} else {
			url += `?target_url=${encodeURIComponent(endpointId)}`;
		}
	}

	const res = await fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'X-Webhook-Source': source,
			...getAuthHeaders()
		},
		credentials: 'include',
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
