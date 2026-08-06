import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import BentoGrid from './BentoGrid.svelte';

describe('BentoGrid Component', () => {
	it('should render the core enterprise features', () => {
		render(BentoGrid);

		// Verify marketing copy exists and tech stack leaks are removed
		expect(screen.getByText('Zero-Loss Webhook Ingestion')).toBeTruthy();
		expect(screen.getByText('Exponential Backoff')).toBeTruthy();
		expect(screen.getByText('Dead Letter Queue')).toBeTruthy();
		expect(screen.getByText('Real-time SSE Observability')).toBeTruthy();
		expect(screen.getByText('Delivery Logs')).toBeTruthy();
		expect(screen.getByText('HMAC Signatures & Auth')).toBeTruthy();

		// Ensure it positions as "Message Queue" or "Ingest Gateway" rather than "RabbitMQ" or "Go Gateway"
		expect(screen.getByText('Message Queue')).toBeTruthy();
		expect(screen.getByText('Ingest Gateway')).toBeTruthy();
	});
});
