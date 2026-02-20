<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { notebookAPI } from '$lib/api';
	import { wsService } from '$lib/websocket';
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
	
	onMount(() => {
		loadNotebooks();

		// Set up WebSocket listeners for real-time updates
		const handleNotebookUpdate = () => {
			loadNotebooks();
		};

		const handleNotebookCreate = () => {
			loadNotebooks();
		};

		const handleNotebookDelete = () => {
			loadNotebooks();
		};

		wsService.on('notebook_update', handleNotebookUpdate);
		wsService.on('notebook_create', handleNotebookCreate);
		wsService.on('notebook_delete', handleNotebookDelete);

		return () => {
			wsService.off('notebook_update', handleNotebookUpdate);
			wsService.off('notebook_create', handleNotebookCreate);
			wsService.off('notebook_delete', handleNotebookDelete);
		};
	});
	
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

<svelte:head>
	<title>MyFlowBook - Tonish</title>
</svelte:head>

<div class="flex flex-col lg:flex-row">
	<!-- Sidebar -->
	<div class={`${sidebarCollapsed ? 'lg:w-16' : 'lg:w-72'} w-full bg-gray-900 border-r border-gray-800 transition-all duration-200 lg:flex-shrink-0 overflow-hidden`}>
		<div class="p-3 sm:p-4 h-full flex flex-col">
			<!-- Header -->
			<div class="flex justify-between items-center mb-3">
				{#if !sidebarCollapsed}
					<div class="flex items-center gap-2">
						<h2 class="text-base font-semibold text-white">Notebooks</h2>
					</div>
				{/if}
				<button
					onclick={() => sidebarCollapsed = !sidebarCollapsed}
					class="text-gray-400 hover:bg-gray-800 p-2 rounded transition"
					title="Toggle"
				>
					{#if sidebarCollapsed}
						<Menu size={16} />
					{:else}
						<X size={16} />
					{/if}
				</button>
			</div>
			
			{#if !sidebarCollapsed}
				<!-- New Button -->
				<button
					onclick={() => showAddNotebook = !showAddNotebook}
					class="w-full flex items-center justify-center gap-2 p-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition text-sm font-medium mb-3"
				>
					<Plus size={14} /> New
				</button>
				
				<!-- List -->
				<div class="flex-1 overflow-y-auto space-y-2">
					{#if pinnedNotebooks.length > 0}
						<div class="mb-3">
							<h3 class="text-[10px] font-semibold text-gray-400 mb-2 px-2 uppercase tracking-wide">Pinned</h3>
							{#each pinnedNotebooks as nb}
								<a
									href="/myflowbook/{nb.id}"
									class="block p-2 bg-yellow-900 border border-yellow-800 rounded-md hover:bg-yellow-800 transition text-sm"
								>
									<div class="font-medium text-white truncate text-sm">{nb.name}</div>
									<div class="text-xs text-gray-400">{nb.pages?.length || 0} pages</div>
								</a>
							{/each}
						</div>
					{/if}
					
					{#if regularNotebooks.length > 0}
						<div>
							<h3 class="text-[10px] font-semibold text-gray-400 mb-2 px-2 uppercase tracking-wide">All</h3>
							{#each regularNotebooks as nb}
								<a
									href="/myflowbook/{nb.id}"
									class="block p-2 bg-gray-800 rounded-md hover:bg-gray-700 transition text-sm border border-gray-700"
								>
									<div class="font-medium text-white truncate text-sm">{nb.name}</div>
									<div class="text-xs text-gray-400">{nb.pages?.length || 0} pages</div>
								</a>
							{/each}
						</div>
					{/if}
				</div>
				
				<!-- Footer -->
				<a href="/" class="flex items-center gap-2 p-2 hover:bg-gray-800 rounded text-sm text-gray-300 border-t border-gray-800 mt-3 pt-3">
					<Home size={14} /> Home
				</a>
			{:else}
				<div class="flex flex-col items-center gap-3 mt-4">
					<button onclick={() => showAddNotebook = !showAddNotebook} class="p-2 hover:bg-gray-800 rounded" title="New">
						<Plus size={16} />
					</button>
					<a href="/" class="p-2 hover:bg-gray-800 rounded" title="Home">
						<Home size={16} />
					</a>
				</div>
			{/if}
		</div>
	</div>
	
	<!-- Main -->
	<div class="flex-1 bg-gray-950 overflow-auto">
		<div class="max-w-[1200px] mx-auto px-3 sm:px-4 lg:px-6 py-3 sm:py-4 space-y-3">
		<!-- Notification -->
		{#if showNotification}
			<div class="fixed top-4 right-4 bg-green-600 text-white px-4 py-2 rounded-lg shadow-md z-50 text-sm">
				{notificationMessage}
			</div>
		{/if}
		
		<!-- Header: just search + new button -->
		<div class="flex gap-2 items-center">
			<div class="flex-1 bg-gray-900 rounded-lg border border-gray-800 p-2 flex items-center gap-2">
				<Search size={14} class="text-gray-400 ml-1" />
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search notebooks..."
					class="flex-1 px-2 py-1.5 border-none bg-gray-900 text-white placeholder:text-gray-500 focus:outline-none text-sm"
				/>
				{#if searchQuery}
					<button onclick={() => searchQuery = ''} class="mr-1 p-1 hover:bg-gray-800 rounded text-gray-400"><X size={14} /></button>
				{/if}
			</div>
			<button
				onclick={() => showAddNotebook = !showAddNotebook}
				class="bg-blue-600 text-white px-3 py-2 min-h-[40px] rounded-lg hover:bg-blue-700 transition flex items-center gap-1.5 text-sm font-medium flex-shrink-0"
			><Plus size={14} /> New</button>
		</div>
		
		<!-- Create Form -->
		{#if showAddNotebook}
			<div class="bg-gray-900 rounded-lg border border-gray-800 p-3 sm:p-4">
				<h2 class="font-semibold text-white mb-3 text-sm">New Notebook</h2>
				<form onsubmit={(e) => { e.preventDefault(); handleAddNotebook(); }} class="space-y-2">
					<input
						type="text"
						bind:value={newNotebook.name}
						required
						placeholder="Name"
						class="w-full px-3 py-2 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
					/>
					<input
						type="text"
						bind:value={newNotebook.tags}
						placeholder="Tags"
						class="w-full px-3 py-2 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
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
							class="px-4 py-2 bg-gray-700 text-gray-300 rounded-lg hover:bg-gray-600 transition text-sm"
						>
							Cancel
						</button>
					</div>
				</form>
			</div>
		{/if}		
		<!-- Content -->
		{#if loading}
			<div class="py-8 text-center text-gray-400 text-sm">Loading...</div>
		{:else if notebooks.length === 0}
			<div class="bg-gray-800 rounded-lg p-8 text-center">
				<p class="text-gray-400 text-sm mb-3">No notebooks</p>
				<button
					onclick={() => showAddNotebook = true}
					class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition text-sm font-medium inline-flex items-center gap-2"
				>
					<Plus size={14} /> Create
				</button>
			</div>
		{:else if filteredNotebooks.length === 0}
			<div class="bg-gray-900 rounded-lg border border-gray-800 p-8 text-center text-gray-400 text-sm">
				No match
			</div>
		{:else}
			<div class="space-y-3">
				{#if pinnedNotebooks.length > 0}
					<div>
						<h3 class="text-sm font-semibold text-gray-300 mb-2">Pinned ({pinnedNotebooks.length})</h3>
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
							{#each pinnedNotebooks as nb}
								<div class="bg-yellow-900 border border-yellow-700 rounded-lg p-3">
									<div class="flex justify-between items-start gap-2 mb-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 min-w-0">
											<h3 class="font-semibold text-white truncate hover:text-blue-400 text-sm">{nb.name}</h3>
										</a>
										<button
											onclick={() => handleTogglePin(nb)}
											class="text-yellow-400 hover:text-yellow-300 p-1 flex-shrink-0"
										>
											<Pin size={14} class="fill-current" />
										</button>
									</div>
									<p class="text-xs text-gray-400 mb-2">{nb.pages?.length || 0} pages</p>
									<div class="flex gap-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 bg-blue-600 text-white text-center px-3 py-1.5 rounded-md text-xs font-medium hover:bg-blue-700 transition">
											Open
										</a>
										<button
											onclick={() => handleDeleteNotebook(nb.id)}
											class="bg-red-950 text-red-300 hover:bg-red-900 px-2 py-1.5 rounded-md text-xs"
										>
											<Trash2 size={12} />
										</button>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
				
				{#if regularNotebooks.length > 0}
					<div>
						<h3 class="text-sm font-semibold text-gray-300 mb-2">All ({regularNotebooks.length})</h3>
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
							{#each regularNotebooks as nb}
								<div class="bg-gray-900 border border-gray-800 rounded-lg p-3 hover:border-gray-700 transition-all">
									<div class="flex justify-between items-start gap-2 mb-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 min-w-0">
											<h3 class="font-semibold text-white truncate hover:text-blue-400 text-sm">{nb.name}</h3>
										</a>
										<button
											onclick={() => handleTogglePin(nb)}
											class="text-gray-400 hover:text-yellow-400 p-1 flex-shrink-0"
										>
											<Pin size={14} />
										</button>
									</div>
									<p class="text-xs text-gray-400 mb-2">{nb.pages?.length || 0} pages</p>
									<div class="flex gap-2">
										<a href="/myflowbook/{nb.id}" class="flex-1 bg-blue-600 text-white text-center px-3 py-1.5 rounded-md text-xs font-medium hover:bg-blue-700 transition">
											Open
										</a>
										<button
											onclick={() => handleDeleteNotebook(nb.id)}
											class="bg-red-950 text-red-300 hover:bg-red-900 px-2 py-1.5 rounded-md text-xs"
										>
											<Trash2 size={12} />
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
