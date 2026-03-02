<script lang="ts">
	// Registration is currently disabled
</script>

<svelte:head><title>Register – Tonish</title></svelte:head>

<div class="min-h-screen bg-gray-950 flex items-center justify-center p-4">
	<div class="w-full max-w-sm">

		<div class="flex flex-col items-center mb-6">
			<div class="w-12 h-12 rounded-2xl bg-gray-800 border border-gray-700 flex items-center justify-center mb-3">
				<span class="text-gray-400 font-black text-xl tracking-tight">T</span>
			</div>
			<h1 class="text-lg font-bold text-white tracking-tight">Tonish</h1>
			<p class="text-[11px] text-gray-500 mt-0.5">Task &amp; flow management</p>
		</div>

		<div class="bg-gray-900 border border-gray-800 rounded-2xl shadow-2xl overflow-hidden">
			<div class="h-0.5 bg-gradient-to-r from-yellow-600 via-yellow-400 to-yellow-600"></div>
			<div class="p-6 space-y-4">
				<h2 class="text-sm font-semibold text-white">Registration Disabled</h2>
				<div class="flex items-start gap-3 bg-yellow-950/50 border border-yellow-800/50 rounded-lg p-3">
					<span class="text-yellow-400 text-sm mt-0.5 flex-shrink-0">⚠</span>
					<div>
						<p class="text-xs font-medium text-yellow-300">Access restricted</p>
						<p class="text-[11px] text-yellow-500 mt-1">User registration is currently disabled. Please contact your administrator for access.</p>
					</div>
				</div>
				<a href="/login" class="flex items-center justify-center w-full h-10 bg-gray-800 hover:bg-gray-700 border border-gray-700 text-gray-300 text-sm font-medium rounded-lg transition">
					Go to Sign In
				</a>
			</div>
		</div>

		<p class="text-center text-[10px] text-gray-600 mt-5">Tonish · Private workspace</p>
	</div>
</div>
