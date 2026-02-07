const API_BASE_URL = import.meta.env.PUBLIC_API_URL || '/api';

let authToken: string | null = null;

if (typeof window !== 'undefined') {
	authToken = localStorage.getItem('authToken');
}

export function setAuthToken(token: string) {
	authToken = token;
	if (typeof window !== 'undefined') {
		localStorage.setItem('authToken', token);
	}
}

export function clearAuthToken() {
	authToken = null;
	if (typeof window !== 'undefined') {
		localStorage.removeItem('authToken');
	}
}

async function fetchAPI(endpoint: string, options: RequestInit = {}) {
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(options.headers as Record<string, string>)
	};

	if (authToken) {
		headers['Authorization'] = `Bearer ${authToken}`;
	}

	const response = await fetch(`${API_BASE_URL}${endpoint}`, {
		...options,
		headers
	});

	if (!response.ok) {
		if (response.status === 401) {
			throw new Error('Authorization header required');
		}
		const error = await response.json().catch(() => ({ error: 'Request failed' }));
		throw new Error(error.error || 'Request failed');
	}

	if (response.status === 204) {
		return null;
	}

	return response.json();
}

// Auth API
export const authAPI = {
	register: (data: { email: string; password: string; name: string }) =>
		fetchAPI('/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	login: (data: { email: string; password: string }) =>
		fetchAPI('/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		}),

	getCurrentUser: () => fetchAPI('/user/me')
};

// Task API
export const taskAPI = {
	getAll: () => fetchAPI('/tasks'),
	getArchived: () => fetchAPI('/tasks/archived'),
	getOne: (id: number) => fetchAPI(`/tasks/${id}`),
	create: (data: any) =>
		fetchAPI('/tasks', {
			method: 'POST',
			body: JSON.stringify(data)
		}),
	update: (id: number, data: any) =>
		fetchAPI(`/tasks/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),
	delete: (id: number) =>
		fetchAPI(`/tasks/${id}`, {
			method: 'DELETE'
		}),
	restore: (id: number) =>
		fetchAPI(`/tasks/${id}/restore`, {
			method: 'POST'
		}),
	permanentDelete: (id: number) =>
		fetchAPI(`/tasks/${id}/permanent`, {
			method: 'DELETE'
		}),
	archive: (id: number) =>
		fetchAPI(`/tasks/${id}/archive`, {
			method: 'POST'
		}),
	getByStatus: (status: string) => fetchAPI(`/tasks/status?status=${status}`),
	getByQuadrant: (quadrant: string) => fetchAPI(`/tasks/quadrant/${quadrant}`)
};

// Notebook API
export const notebookAPI = {
	getAll: () => fetchAPI('/notebooks'),
	getOne: (id: number) => fetchAPI(`/notebooks/${id}`),
	create: (data: any) =>
		fetchAPI('/notebooks', {
			method: 'POST',
			body: JSON.stringify(data)
		}),
	update: (id: number, data: any) =>
		fetchAPI(`/notebooks/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),
	delete: (id: number) =>
		fetchAPI(`/notebooks/${id}`, {
			method: 'DELETE'
		})
};

// Page API
export const pageAPI = {
	getOne: (id: number) => fetchAPI(`/pages/${id}`),
	create: (data: any) =>
		fetchAPI('/pages', {
			method: 'POST',
			body: JSON.stringify(data)
		}),
	update: (id: number, data: any) =>
		fetchAPI(`/pages/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),
	delete: (id: number) =>
		fetchAPI(`/pages/${id}`, {
			method: 'DELETE'
		}),
	search: (query: string) => fetchAPI(`/pages/search?q=${encodeURIComponent(query)}`)
};
