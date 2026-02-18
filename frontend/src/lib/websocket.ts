import { writable, type Writable } from 'svelte/store';

interface WebSocketMessage {
	type: string;
	data: any;
	user_id?: number;
}

type MessageHandler = (data: any) => void;

class WebSocketService {
	private ws: WebSocket | null = null;
	private reconnectTimeout: NodeJS.Timeout | null = null;
	private messageHandlers: Map<string, Set<MessageHandler>> = new Map();
	private reconnectAttempts = 0;
	private maxReconnectAttempts = 5;
	private reconnectDelay = 2000;
	
	public connected: Writable<boolean> = writable(false);

	connect(userId?: number) {
		if (this.ws?.readyState === WebSocket.OPEN) {
			return;
		}

		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const host = window.location.hostname;
		const port = '50002'; // Backend port
		const userParam = userId ? `?user_id=${userId}` : '';
		const wsUrl = `${protocol}//${host}:${port}/ws${userParam}`;

		try {
			this.ws = new WebSocket(wsUrl);

			this.ws.onopen = () => {
				console.log('WebSocket connected');
				this.connected.set(true);
				this.reconnectAttempts = 0;
			};

			this.ws.onmessage = (event) => {
				try {
					const message: WebSocketMessage = JSON.parse(event.data);
					this.handleMessage(message);
				} catch (error) {
					console.error('Error parsing WebSocket message:', error);
				}
			};

			this.ws.onclose = () => {
				console.log('WebSocket disconnected');
				this.connected.set(false);
				this.scheduleReconnect(userId);
			};

			this.ws.onerror = (error) => {
				console.error('WebSocket error:', error);
			};
		} catch (error) {
			console.error('Error creating WebSocket:', error);
			this.scheduleReconnect(userId);
		}
	}

	private scheduleReconnect(userId?: number) {
		if (this.reconnectAttempts >= this.maxReconnectAttempts) {
			console.error('Max reconnect attempts reached');
			return;
		}

		if (this.reconnectTimeout) {
			clearTimeout(this.reconnectTimeout);
		}

		this.reconnectAttempts++;
		const delay = this.reconnectDelay * this.reconnectAttempts;
		
		console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);
		
		this.reconnectTimeout = setTimeout(() => {
			this.connect(userId);
		}, delay);
	}

	private handleMessage(message: WebSocketMessage) {
		const handlers = this.messageHandlers.get(message.type);
		if (handlers) {
			handlers.forEach(handler => handler(message.data));
		}
	}

	on(messageType: string, handler: MessageHandler) {
		if (!this.messageHandlers.has(messageType)) {
			this.messageHandlers.set(messageType, new Set());
		}
		this.messageHandlers.get(messageType)!.add(handler);
	}

	off(messageType: string, handler: MessageHandler) {
		const handlers = this.messageHandlers.get(messageType);
		if (handlers) {
			handlers.delete(handler);
		}
	}

	disconnect() {
		if (this.reconnectTimeout) {
			clearTimeout(this.reconnectTimeout);
			this.reconnectTimeout = null;
		}
		
		if (this.ws) {
			this.ws.close();
			this.ws = null;
		}
		
		this.connected.set(false);
		this.reconnectAttempts = 0;
	}
}

export const wsService = new WebSocketService();
export type { WebSocketMessage };
