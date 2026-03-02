<script lang="ts">
	import { authAPI, setAuthToken } from '$lib/api';
	import { Lock, Mail } from 'lucide-svelte';

	let email    = $state('');
	let password = $state('');
	let error    = $state('');
	let loading  = $state(false);

	async function handleLogin() {
		loading = true;
		error   = '';
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

<svelte:head><title>Sign In – Tonish</title></svelte:head>

<div class="min-h-screen bg-gray-950 flex items-center justify-center p-4">
	<div class="w-full max-w-sm">

		<!-- Logo / Brand -->
		<div class="flex flex-col items-center mb-6">
			<div class="w-12 h-12 rounded-2xl bg-blue-600 flex items-center justify-center mb-3 shadow-lg shadow-blue-900/40">
				<span class="text-white font-black text-xl tracking-tight">T</span>
			</div>
			<h1 class="text-lg font-bold text-white tracking-tight">Tonish</h1>
			<p class="text-[11px] text-gray-500 mt-0.5">Task &amp; flow management</p>
		</div>

		<!-- Card -->
		<div class="bg-gray-900 border border-gray-800 rounded-2xl shadow-2xl overflow-hidden">
			<div class="h-0.5 bg-gradient-to-r from-blue-600 via-blue-400 to-blue-600"></div>
			<div class="p-6 space-y-4">
				<h2 class="text-sm font-semibold text-white">Sign in to your workspace</h2>

				{#if error}
					<div class="flex items-start gap-2 bg-red-950/60 border border-red-800/60 text-red-300 px-3 py-2.5 rounded-lg text-xs">
						<span class="mt-0.5 flex-shrink-0">⚠</span>
						<span>{error}</span>
					</div>
				{/if}

				<form onsubmit={(e) => { e.preventDefault(); handleLogin(); }} class="space-y-3">
					<div class="space-y-1.5">
						<label for="email" class="block text-[11px] font-medium text-gray-400 uppercase tracking-wide">Email</label>
						<div class="relative">
							<Mail size={13} class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 pointer-events-none" />
							<input
								type="email" id="email" bind:value={email} required placeholder="you@example.com"
								class="w-full h-10 pl-8 pr-3 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-600 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
							/>
						</div>
					</div>

					<div class="space-y-1.5">
						<label for="password" class="block text-[11px] font-medium text-gray-400 uppercase tracking-wide">Password</label>
						<div class="relative">
							<Lock size={13} class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 pointer-events-none" />
							<input
								type="password" id="password" bind:value={password} required placeholder="••••••••"
								class="w-full h-10 pl-8 pr-3 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-600 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
							/>
						</div>
					</div>

					<button
						type="submit" disabled={loading}
						class="w-full h-10 mt-1 bg-blue-600 hover:bg-blue-500 disabled:bg-blue-800 disabled:cursor-not-allowed text-white text-sm font-semibold rounded-lg transition flex items-center justify-center gap-2 shadow-md shadow-blue-900/30"
					>
						{#if loading}
							<span class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
							Signing in…
						{:else}
							Sign In
						{/if}
					</button>
				</form>
			</div>
		</div>

		<p class="text-center text-[10px] text-gray-600 mt-5">Tonish · Private workspace</p>
	</div>
</div>
