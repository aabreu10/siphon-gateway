export type WebhookStatus = 'PENDING' | 'SUCCESS' | 'FAILED_DLQ';

export type SSEEventType = 'received' | 'retrying' | 'success' | 'failed_dlq' | 'replayed';

export interface Webhook {
	id: string;
	source: string;
	payload: Record<string, unknown>;
	target_url: string;
	status: WebhookStatus;
	retry_count: number;
	created_at: string;
	updated_at: string;
}

export interface SSEEvent {
	type: SSEEventType;
	webhook: Webhook;
}

export interface WebhookListResponse {
	webhooks: Webhook[];
	total: number;
	limit: number;
	offset: number;
}
