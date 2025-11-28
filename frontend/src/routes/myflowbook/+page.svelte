<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { notebookAPI } from '$lib/api';
	import { Plus, Search, Pin, Menu, X, Trash2, Home } from 'lucide-svelte';

	type Notebook = {
		id: number;
		name: string;
		tags?: string;
		is_pinned?: boolean;
		pages?: { id: number }[];
	};

	let notebooks = $state<Notebook[]>([]);
	let loading = $state(true);
	let showAddNotebook = $state(false);
	let searchQuery = $state('');
	let sidebarCollapsed = $state(false);
	let showNotification = $state(false);
	let notificationMessage = $state('');
	let newNotebook = $state<{ name: string; tags: string }>({ name: '', tags: '' });
	
	const filteredNotebooks = $derived(
		!searchQuery
			? notebooks
			: notebooks.filter((n) => {
				const q = searchQuery.toLowerCase();
				return n.name.toLowerCase().includes(q) || 
					   (n.tags ?? '').toLowerCase().includes(q);
			})
	);

	const pinnedNotebooks = $derived(filteredNotebooks.filter((n) => n.is_pinned));
	const regularNotebooks = $derived(filteredNotebooks.filter((n) => !n.is_pinned));
	
	onMount(loadNotebooks);
	
	async function loadNotebooks() {
		// Check if user is authenticated
		const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
		if (!token) {
			goto('/login');
			return;
		}
		try {
			notebooks = await notebookAPI.getAll();
		} catch (error: any) {
			console.error('Failed to load notebooks:', error);
			// If unauthorized, redirect to login
			if (error.message?.includes('Authorization') || error.message?.includes('401')) {
				localStorage.removeItem('authToken');
				goto('/login');
			}
		} finally {
			loading = false;
		}
	}
	
	async function handleAddNotebook() {
		try {
			await notebookAPI.create(newNotebook);
			await loadNotebooks();
			showAddNotebook = false;
			newNotebook = { name: '', tags: '' };
		} catch (error) {
			console.error('Failed to create notebook:', error);
		}
	}
	
	async function handleDeleteNotebook(id: number) {
		if (confirm('Delete this notebook?')) {
			try {
				await notebookAPI.delete(id);
				await loadNotebooks();
				showNotification = true;
				notificationMessage = 'Deleted';
				setTimeout(() => showNotification = false, 2000);
			} catch (error) {
				console.error('Failed to delete:', error);
			}
		}
	}
	
	async function handleTogglePin(notebook: Notebook) {
		try {
			await notebookAPI.update(notebook.id, {
				...notebook,
				is_pinned: !notebook.is_pinned
			});
			await loadNotebooks();
		} catch (error) {
			console.error('Failed to toggle pin:', error);
		}
	}
</script>

<div class="flex flex-col lg:flex-row min-h-screen bg-gray-50">
	<!-- Sidebar -->
	<div class={`${sidebarCollapsed ? 'lg:w-16' : 'lg:w-72'} w-full bg-white border-r border-gray-200 transition-all duration-200 lg:flex-shrink-0 overflow-hidden`}>
		<div class="p-4 h-full flex flex-col">
			<!-- Header -->
			<div class="flex justify-between items-center mb-4">
				{#if !sidebarCollapsed}
					<h2 class="text-lg font-semibold text-gray-900">Notebooks</h2>
				{/if}
				<button
					onclick={() => sidebarCollapsed = !sidebarCollapsed}
					class="text-gray-600 hover:bg-gray-100 p-2 rounded transition"
					title="Toggle"
				>
					{#if sidebarCollapsed}
						<Menu size={18} />
					{:else}
						<X size={18} />
					{/if}
				</button>
			</div>
			
			{#if !sidebarCollapsed}
				<!-- New Button -->
				<button
					onclick={() => showAddNotebook = !showAddNotebook}
					class="w-full flex items-center justify-center gap-2 p-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition text-sm font-medium mb-4"
				>
					<Plus size={16} /> New
				</button>
				
				<!-- List -->
				<div class="flex-1 overflow-y-auto space-y-2">
					{#if pinnedNotebooks.length > 0}
						<div class="mb-3">
							<h3 class="text-xs font-semibold text-gray-600 mb-2 px-2 uppercase">Pinned</h3>
							{#each pinnedNotebooks as nb}
								<a
									href="/myflowbook/{nb.id}"
									class="block p-2 bg-yellow-100 border border-yellow-300 rounded-lg hover:bg-yellow-50 transition text-sm"
								>
									<div class="font-medium text-gray-900 truncate">{nb.name}</div>
									<div class="text-xs text-gray-600">{nb.pages?.length || 0}p</div>
								</a>
							{/each}
						</div>
					{/if}
					
					{#if regularNotebooks.length > 0}
						<div>
							<h3 class="text-xs font-semibold text-gray-600 mb-2 px-2 uppercase">All</h3>
							{#each regularNotebooks as nb}
								<a
									href="/myflowbook/{nb.id}"
									class="block p-2 bg-gray-100 rounded-lg hover:bg-gray-200 transition text-sm"
								>
									<div class="font-medium text-gray-900 truncate">{nb.name}</div>
									<div class="text-xs text-gray-600">{nb.pages?.length || 0}p</div>
								</a>
							{/each}
						</div>
					{/if}
				</div>
				
				<!-- Footer -->
				<a href="/" class="flex items-center gap-2 p-2 hover:bg-gray-100 rounded text-sm text-gray-700 border-t mt-3 pt-3">
					<Home size={16} /> Home
				</a>
			{:else}
				<div class="flex flex-col items-center gap-3 mt-4">
					<button onclick={() => showAddNotebook = !showAddNotebook} class="p-2 hover:bg-gray-100 rounded" title="New">
						<Plus size={18} />
					</button>
					<a href="/" class="p-2 hover:bg-gray-100 rounded" title="Home">
						<Home size={18} />
					</a>
				</div>
			{/if}
		</div>
	</div>
	
	<!-- Main -->
	<div class="flex-1 bg-white overflow-auto">
		<div class="max-w-4xl mx-auto p-6 space-y-4">
		<!-- Notification -->
		{#if showNotification}
			<div class="fixed top-4 right-4 bg-green-600 text-white px-4 py-2 rounded-lg shadow-md z-50 text-sm">
				{notificationMessage}
			</div>
		{/if}
		
		<!-- Header -->
		<div class="flex justify-between items-center">
			<h1 class="text-3xl font-bold text-gray-900">MyFlowBook</h1>
			<button
				onclick={() => showAddNotebook = !showAddNotebook}
				class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition flex items-center gap-2 text-sm font-medium"
			>
				<Plus size={16} /> New
			</button>
		</div>
		
		<!-- Search -->
		<div class="bg-gray-100 rounded-lg p-2 flex items-center gap-2">
			<Search size={16} class="text-gray-500 ml-1" />
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search..."
				class="flex-1 px-2 py-1 border-none bg-gray-100 focus:outline-none text-sm"
			/>
			{#if searchQuery}
				<button
					onclick={() => searchQuery = ''}
					class="mr-1 p-1 hover:bg-gray-200 rounded text-gray-500"
				>
					<X size={16} />
				</button>
			{/if}
		</div>
		
		<!-- Create Form -->
		{#if showAddNotebook}
			<div class="bg-gray-100 rounded-lg border border-gray-300 p-4">
				<h2 class="font-bold text-gray-900 mb-3">New Notebook</h2>
				<form onsubmit={(e) => { e.preventDefault(); handleAddNotebook(); }} class="space-y-2">
					<input
						type="text"
						bind:value={newNotebook.name}
						required
						placeholder="Name"
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
					/>
					<input
						type="text"
						bind:value={newNotebook.tags}
						placeholder="Tags"
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
					/>
					<div class="flex gap-2 pt-1">
						<button
							type="submit"
							class="flex-1 bg-blue-600 text-white px-3 py-2 rounded-lg hover:bg-blue-700 transition text-sm font-medium"
						>
							Create
						</button>
						<button
							type="button"
							onclick={() => { showAddNotebook = false; newNotebook = { name: '', tags: '' }; }}
							class="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition text-sm"
						>
							Cancel
						</button>
					</div>
				</form>
			</div>
		{/if}		
		<!-- Content -->
		{#if loading}
			<div class="py-8 text-center text-gray-500 text-sm">Loading...</div>
		{:else if notebooks.length === 0}
			<div class="bg-gray-100 rounded-lg p-8 text-center">
				<p class="text-gray-600 text-sm mb-3">No notebooks</p>
				<button
					onclick={() => showAddNotebook = true}
					class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition text-sm font-medium inline-flex items-center gap-2"
				>
					<Plus size={14} /> Create
				</button>
			</div>
		{:else if filteredNotebooks.length === 0}
			<div class="bg-gray-100 rounded-lg p-8 text-center text-gray-600 text-sm">
				No match
			</div>
		{:else}
			<div class="space-y-4">
				{#if pinnedNotebooks.length > 0}
					<div>
						<h3 class="text-sm font-semibold text-gray-700 mb-2">Pinned ({pinnedNotebooks.length})</h3>
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
							{#each pinnedNotebooks as nb}
								<div class="bg-yellow-100 border-2 border-yellow-400 rounded-lg p-4">
									<div class="flex justify-between items-start gap-2 mb-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 min-w-0">
											<h3 class="font-semibold text-gray-900 truncate hover:text-blue-600 text-sm">{nb.name}</h3>
										</a>
										<button
											onclick={() => handleTogglePin(nb)}
											class="text-yellow-600 hover:text-yellow-700 p-1 flex-shrink-0"
										>
											<Pin size={16} class="fill-current" />
										</button>
									</div>
									<p class="text-xs text-gray-600 mb-3">{nb.pages?.length || 0} pages</p>
									<div class="flex gap-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 bg-blue-600 text-white text-center px-3 py-1.5 rounded text-xs font-medium hover:bg-blue-700 transition">
											Open
										</a>
										<button
											onclick={() => handleDeleteNotebook(nb.id)}
											class="bg-red-200 text-red-700 hover:bg-red-300 px-2 py-1.5 rounded text-xs"
										>
											<Trash2 size={14} />
										</button>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
				
				{#if regularNotebooks.length > 0}
					<div>
						<h3 class="text-sm font-semibold text-gray-700 mb-2">All ({regularNotebooks.length})</h3>
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
							{#each regularNotebooks as nb}
								<div class="bg-white border border-gray-300 rounded-lg p-4">
									<div class="flex justify-between items-start gap-2 mb-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 min-w-0">
											<h3 class="font-semibold text-gray-900 truncate hover:text-blue-600 text-sm">{nb.name}</h3>
										</a>
										<button
											onclick={() => handleTogglePin(nb)}
											class="text-gray-400 hover:text-yellow-500 p-1 flex-shrink-0"
										>
											<Pin size={16} />
										</button>
									</div>
									<p class="text-xs text-gray-600 mb-3">{nb.pages?.length || 0} pages</p>
									<div class="flex gap-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 bg-blue-600 text-white text-center px-3 py-1.5 rounded text-xs font-medium hover:bg-blue-700 transition">
											Open
										</a>
										<button
											onclick={() => handleDeleteNotebook(nb.id)}
											class="bg-red-100 text-red-600 hover:bg-red-200 px-2 py-1.5 rounded text-xs"
										>
											<Trash2 size={14} />
										</button>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		{/if}
		</div>
	</div>

	<button
		type="button"
		onclick={() => showAddNotebook = !showAddNotebook}
		class="lg:hidden fixed bottom-24 right-5 inline-flex items-center justify-center rounded-full bg-gradient-to-r from-purple-600 to-purple-700 text-white text-2xl h-14 w-14 shadow-2xl focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-purple-600 hover:shadow-lg transition"
		aria-label="Add notebook"
	>
		+
	</button>
</div>

<style>
	@keyframes slide-in {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes scale-in {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	.animate-slide-in {
		animation: slide-in 0.3s ease-out;
	}

	.animate-scale-in {
		animation: scale-in 0.2s ease-out;
	}
</style>
