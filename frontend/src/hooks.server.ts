import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	// Proxy API requests to the backend
	if (event.url.pathname.startsWith('/api')) {
		// Use Docker service name for internal communication
		const backendUrl = process.env.BACKEND_URL_INTERNAL || 'http://backend:8080';
		const apiPath = event.url.pathname; // Keep /api prefix
		const targetUrl = `${backendUrl}${apiPath}${event.url.search}`;

		try {
			const headers = new Headers(event.request.headers);
			// Remove host header to avoid conflicts
			headers.delete('host');
			headers.delete('connection');

			const response = await fetch(targetUrl, {
				method: event.request.method,
				headers: headers,
				body: event.request.method !== 'GET' && event.request.method !== 'HEAD' 
					? await event.request.text() 
					: undefined,
			});

			// Create a new response with the backend's response
			const responseHeaders = new Headers(response.headers);
			
			return new Response(response.body, {
				status: response.status,
				statusText: response.statusText,
				headers: responseHeaders,
			});
		} catch (error) {
			console.error('Proxy error:', error);
			return new Response(JSON.stringify({ error: 'Failed to connect to backend' }), {
				status: 503,
				headers: {
					'Content-Type': 'application/json',
				},
			});
		}
	}

	// For all other requests, proceed normally
	return resolve(event);
};
