<script lang="ts">
	import { authAPI, setAuthToken } from '$lib/api';
	
	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);
	
	async function handleLogin() {
		loading = true;
		error = '';
		
		try {
			const response = await authAPI.login({ email, password });
			setAuthToken(response.token);
			window.location.href = '/';
		} catch (err: any) {
			error = err.message || 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-950">
	<div class="max-w-md w-full bg-gray-900 rounded-lg shadow-md border border-gray-800 p-8">
		<div class="flex flex-col items-center mb-6">
			<h1 class="text-2xl font-bold text-white">Tonish</h1>
			<p class="text-sm text-gray-400 mt-1">Task management made simple</p>
		</div>
		<h2 class="text-xl font-semibold text-center mb-6 text-gray-300">Sign In</h2>
		
		{#if error}
			<div class="bg-red-950 border border-red-800 text-red-300 px-4 py-3 rounded mb-4">
				{error}
			</div>
		{/if}
		
		<form onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
			<div class="mb-4">
				<label for="email" class="block text-sm font-medium text-gray-300 mb-2">
					Email
				</label>
				<input
					type="email"
					id="email"
					bind:value={email}
					required
					class="w-full px-3 py-2 bg-gray-800 border border-gray-700 text-white rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			
			<div class="mb-6">
				<label for="password" class="block text-sm font-medium text-gray-300 mb-2">
					Password
				</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					required
					class="w-full px-3 py-2 bg-gray-800 border border-gray-700 text-white rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			
			<button
				type="submit"
				disabled={loading}
				class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-blue-300"
			>
				{loading ? 'Logging in...' : 'Login'}
			</button>
		</form>
		
		<!-- Registration disabled -->
		<!-- <p class="text-center text-sm text-gray-600 mt-4">
			Don't have an account? <a href="/register" class="text-blue-600 hover:underline">Register</a>
		</p> -->
	</div>
</div>
