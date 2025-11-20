<script lang="ts">
	import { authAPI, setAuthToken } from '$lib/api';
	
	let email = $state('');
	let password = $state('');
	let name = $state('');
	let error = $state('');
	let loading = $state(false);
	
	async function handleRegister() {
		loading = true;
		error = '';
		
		try {
			await authAPI.register({ email, password, name });
			const loginResponse = await authAPI.login({ email, password });
			setAuthToken(loginResponse.token);
			window.location.href = '/';
		} catch (err: any) {
			error = err.message || 'Registration failed';
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100">
	<div class="max-w-md w-full bg-white rounded-lg shadow-md p-8">
		<h1 class="text-2xl font-bold text-center mb-6">Register for Tonish</h1>
		
		{#if error}
			<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
				{error}
			</div>
		{/if}
		
		<form onsubmit={(e) => { e.preventDefault(); handleRegister(); }}>
			<div class="mb-4">
				<label for="name" class="block text-sm font-medium text-gray-700 mb-2">
					Name
				</label>
				<input
					type="text"
					id="name"
					bind:value={name}
					required
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			
			<div class="mb-4">
				<label for="email" class="block text-sm font-medium text-gray-700 mb-2">
					Email
				</label>
				<input
					type="email"
					id="email"
					bind:value={email}
					required
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			
			<div class="mb-6">
				<label for="password" class="block text-sm font-medium text-gray-700 mb-2">
					Password
				</label>
				<input
					type="password"
					id="password"
					bind:value={password}
					required
					minlength="6"
					class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>
			</div>
			
			<button
				type="submit"
				disabled={loading}
				class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-blue-300"
			>
				{loading ? 'Registering...' : 'Register'}
			</button>
		</form>
		
		<p class="text-center text-sm text-gray-600 mt-4">
			Already have an account? <a href="/login" class="text-blue-600 hover:underline">Login</a>
		</p>
	</div>
</div>
