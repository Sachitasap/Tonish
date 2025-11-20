<script lang="ts">
	import { onMount } from 'svelte';
	import { notebookAPI } from '$lib/api';
	import { Plus, Search, Pin, Menu, X, Trash2, BookMarked, FileText, Tag, Home, Loader } from 'lucide-svelte';

	type Page = {
		id: number;
		title: string;
		content: string;
		tags?: string;
		created_at?: string;
		updated_at?: string;
	};

	type Notebook = {
		id: number;
		name: string;
		tags?: string;
		is_pinned?: boolean;
		pages?: Page[];
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
			: notebooks.filter((notebook) => {
				const query = searchQuery.toLowerCase();
				const nameMatch = notebook.name.toLowerCase().includes(query);
				const tagsMatch = (notebook.tags ?? '').toLowerCase().includes(query);
				return nameMatch || tagsMatch;
			})
	);

	const pinnedNotebooks = $derived(filteredNotebooks.filter((notebook) => Boolean(notebook.is_pinned)));
	const regularNotebooks = $derived(filteredNotebooks.filter((notebook) => !notebook.is_pinned));

	const notebooksByTag = $derived((() => {
		const tagMap = new Map<string, Notebook[]>();
		notebooks.forEach((notebook) => {
			if (notebook.tags) {
				notebook.tags
					.split(',')
					.map((tag) => tag.trim())
					.filter((tag) => tag)
					.forEach((tag) => {
						if (!tagMap.has(tag)) {
							tagMap.set(tag, []);
						}
						tagMap.get(tag)?.push(notebook);
					});
			}
		});
		return tagMap;
	})());

	const createNotebookNameId = 'create-notebook-name';
	const createNotebookTagsId = 'create-notebook-tags';
	
	onMount(async () => {
		await loadNotebooks();
	});
	
	async function loadNotebooks() {
		try {
			notebooks = await notebookAPI.getAll();
		} catch (error) {
			console.error('Failed to load notebooks:', error);
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
		if (confirm('Are you sure you want to delete this notebook and all its pages?')) {
			try {
				await notebookAPI.delete(id);
				await loadNotebooks();
				showNotification = true;
				notificationMessage = 'Notebook deleted successfully';
				setTimeout(() => showNotification = false, 3000);
			} catch (error) {
				console.error('Failed to delete notebook:', error);
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

<div class="flex flex-col lg:flex-row min-h-screen">
	<!-- Sidebar with Glass-Morphism Effect -->
	<div class={`${sidebarCollapsed ? 'lg:w-20' : 'lg:w-80'} w-full backdrop-blur-lg bg-white/50 border-r border-white/30 text-gray-900 transition-all duration-300 lg:flex-shrink-0 overflow-hidden shadow-lg`}>
		<div class="p-4 h-full flex flex-col">
			<!-- Sidebar Header -->
			<div class="flex justify-between items-center mb-6">
				{#if !sidebarCollapsed}
					<h2 class="text-lg font-bold flex items-center text-gray-900">
						<BookMarked size={20} class="mr-2 text-indigo-600" /> My Notebooks
					</h2>
				{/if}
				<button
					onclick={() => sidebarCollapsed = !sidebarCollapsed}
					class="text-gray-600 hover:bg-white/50 p-2 rounded-lg transition"
					title={sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
				>
					{#if sidebarCollapsed}
						<Menu size={20} />
					{:else}
						<X size={20} />
					{/if}
				</button>
			</div>
			
			{#if !sidebarCollapsed}
				<!-- Quick Stats -->
				<div class="mb-6 p-4 bg-gradient-to-br from-indigo-500/20 to-blue-500/20 rounded-xl border border-indigo-300/30 backdrop-blur-sm">
					<div class="text-3xl font-bold mb-1 text-gray-900">{notebooks.length}</div>
					<div class="text-gray-700 text-sm font-medium">Total Notebooks</div>
					{#if pinnedNotebooks.length > 0}
						<div class="mt-2 text-gray-600 text-sm flex items-center gap-1">
							<Pin size={14} class="text-yellow-500" /> {pinnedNotebooks.length} pinned
						</div>
					{/if}
				</div>
				
				<!-- Quick Actions -->
				<div class="mb-6">
					<button
						onclick={() => showAddNotebook = !showAddNotebook}
						class="w-full flex items-center justify-center space-x-2 p-3 bg-gradient-to-r from-indigo-600 to-blue-600 text-white rounded-lg hover:shadow-md hover:from-indigo-700 hover:to-blue-700 transition font-medium"
					>
						<Plus size={18} />
						<span>New Notebook</span>
					</button>
				</div>
				
				<!-- Notebooks List -->
				<div class="flex-1 overflow-y-auto space-y-4 pr-2">
					<!-- Pinned Notebooks -->
					{#if pinnedNotebooks.length > 0}
						<div>
							<h3 class="text-xs text-indigo-700 mb-2 uppercase font-bold flex items-center gap-1 px-2">
								<Pin size={12} /> Pinned
							</h3>
							<div class="space-y-2">
								{#each pinnedNotebooks as notebook}
									<a
										href="/myflowbook/{notebook.id}"
										class="block p-3 bg-gradient-to-r from-yellow-300/20 to-yellow-400/20 border border-yellow-300/40 rounded-lg hover:from-yellow-300/30 hover:to-yellow-400/30 transition backdrop-blur-sm group"
									>
										<div class="font-semibold truncate text-gray-900 group-hover:text-indigo-700">{notebook.name}</div>
										<div class="text-xs text-gray-600 mt-1 flex items-center gap-1">
											<FileText size={14} /> {notebook.pages?.length || 0} pages
										</div>
									</a>
								{/each}
							</div>
						</div>
					{/if}
					
					<!-- Recent Notebooks -->
					{#if regularNotebooks.length > 0}
						<div>
							<h3 class="text-xs text-indigo-700 mb-2 uppercase font-bold px-2">All Notebooks</h3>
							<div class="space-y-2">
								{#each regularNotebooks as notebook}
									<a
										href="/myflowbook/{notebook.id}"
										class="block p-3 bg-indigo-500/10 border border-indigo-300/20 rounded-lg hover:bg-indigo-500/20 transition backdrop-blur-sm group"
									>
										<div class="font-semibold truncate text-gray-900 group-hover:text-indigo-700">{notebook.name}</div>
										<div class="text-xs text-gray-600 mt-1 flex items-center gap-1">
											<FileText size={14} /> {notebook.pages?.length || 0} pages
										</div>
									</a>
								{/each}
							</div>
						</div>
					{/if}
					
					<!-- Notebooks by Tag -->
					{#if notebooksByTag.size > 0}
						<div>
							<h3 class="text-xs text-indigo-700 mb-2 uppercase font-bold px-2 flex items-center gap-1"><Tag size={14} /> By Tag</h3>
							<div class="space-y-2">
								{#each Array.from(notebooksByTag.entries()) as [tag, tagNotebooks]}
									<details class="group">
										<summary class="cursor-pointer p-2 bg-blue-500/10 border border-blue-300/20 rounded-lg hover:bg-blue-500/20 transition text-sm font-medium flex justify-between items-center backdrop-blur-sm text-gray-900">
											<span class="flex items-center gap-1"><Tag size={14} /> {tag}</span>
											<span class="text-xs text-gray-600 bg-white/50 px-2 py-0.5 rounded-full">{tagNotebooks.length}</span>
										</summary>
										<div class="ml-2 mt-1 space-y-1">
											{#each tagNotebooks as notebook}
												<a
													href="/myflowbook/{notebook.id}"
													class="block p-2 text-sm bg-blue-500/5 border border-blue-300/10 rounded hover:bg-blue-500/15 transition font-medium text-gray-800 backdrop-blur-sm"
												>
													{notebook.name}
												</a>
											{/each}
										</div>
									</details>
								{/each}
							</div>
						</div>
					{/if}
				</div>
				
					<!-- Back to Dashboard -->
					<div class="mt-4 pt-4 border-t border-indigo-300/20">
						<a href="/" class="flex items-center space-x-2 p-2 hover:bg-indigo-500/20 rounded-lg transition text-gray-700 font-medium group">
							<Home size={18} class="group-hover:-translate-x-1 transition" />
							<span>Dashboard</span>
						</a>
					</div>
			{:else}
				<!-- Collapsed view -->
				<div class="flex flex-col items-center space-y-4 mt-8">
					<div class="text-2xl font-bold text-gray-900 bg-indigo-500/20 w-full text-center py-2 rounded-lg">{notebooks.length}</div>
					<button onclick={() => showAddNotebook = !showAddNotebook} class="p-3 hover:bg-indigo-500/20 rounded-lg transition text-gray-700" title="New Notebook">
						<Plus size={20} />
					</button>
					<a href="/" class="p-3 hover:bg-indigo-500/20 rounded-lg transition text-gray-700" title="Dashboard">
						<Home size={20} />
					</a>
				</div>
			{/if}
		</div>
	</div>
	
	<!-- Main Content -->
	<div class="flex-1 bg-gradient-to-b from-gray-50 to-white overflow-auto w-full">
		<div class="max-w-6xl mx-auto p-6 space-y-6">
		<!-- Notification -->
		{#if showNotification}
			<div class="fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50 animate-slide-in flex items-center gap-2">
				<Plus size={18} class="rotate-45" /> {notificationMessage}
			</div>
		{/if}			<!-- Header Section -->
			<div class="flex justify-between items-start gap-4">
				<div>
					<h1 class="text-4xl font-bold text-gray-900 inline-flex items-center gap-3">
						<BookMarked size={32} class="text-purple-600" /> MyFlowBook
					</h1>
					<p class="text-gray-600 mt-2">Organize, write, and manage your notes in one place</p>
				</div>
				<button
					onclick={() => showAddNotebook = !showAddNotebook}
					class="bg-gradient-to-r from-purple-600 to-purple-700 text-white px-6 py-3 rounded-lg hover:shadow-lg transition inline-flex items-center gap-2 font-medium"
				>
					<Plus size={20} /> New Notebook
				</button>
			</div>
			
			<!-- Search Bar with better UX -->
			<div class="bg-white rounded-xl shadow-sm border border-gray-200 p-1 flex items-center gap-2 hover:shadow-md transition">
				<Search size={20} class="text-gray-400 ml-3" />
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search notebooks by name, tags..."
					class="flex-1 px-4 py-3 border-none focus:outline-none bg-white text-gray-900 placeholder-gray-400"
				/>
				{#if searchQuery}
					<button
						onclick={() => searchQuery = ''}
						class="mr-3 p-2 hover:bg-gray-100 rounded-lg transition text-gray-400"
						title="Clear search"
					>
						<X size={18} />
					</button>
				{/if}
			</div>
			
			<!-- Create Notebook Form -->
			{#if showAddNotebook}
				<div class="bg-white rounded-xl shadow-lg border border-purple-100 p-6 animate-scale-in">
					<h2 class="text-2xl font-bold text-gray-900 mb-6">Create New Notebook</h2>
					<form onsubmit={(e) => { e.preventDefault(); handleAddNotebook(); }} class="space-y-4">
						<div>
							<label for={createNotebookNameId} class="block text-sm font-semibold text-gray-700 mb-2">Notebook Name</label>
							<input
								type="text"
								bind:value={newNotebook.name}
								required
								placeholder="e.g., Work Notes, Project Ideas"
								id={createNotebookNameId}
								class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition"
							/>
						</div>
						
						<div>
							<label for={createNotebookTagsId} class="block text-sm font-semibold text-gray-700 mb-2">Tags (comma-separated, optional)</label>
							<input
								type="text"
								bind:value={newNotebook.tags}
								placeholder="e.g., work, urgent, research"
								id={createNotebookTagsId}
								class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition"
							/>
						</div>
						
						<div class="flex gap-2 pt-2">
							<button
								type="submit"
								class="flex-1 bg-gradient-to-r from-purple-600 to-purple-700 text-white px-4 py-3 rounded-lg hover:shadow-md transition font-medium"
							>
								Create Notebook
							</button>
							<button
								type="button"
								onclick={() => { showAddNotebook = false; newNotebook = { name: '', tags: '' }; }}
								class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium"
							>
								Cancel
							</button>
						</div>
					</form>
				</div>
			{/if}
			
			{#if loading}
			<div class="py-16 text-center">
				<div class="inline-flex items-center gap-2 text-gray-500">
					<Loader size={18} class="animate-spin text-indigo-600" />
					<p>Loading your notebooks...</p>
				</div>
			</div>
			{:else if notebooks.length === 0}
			<div class="bg-gradient-to-br from-indigo-50 to-blue-50 rounded-xl shadow-sm border border-indigo-100 p-16 text-center">
				<div class="mb-4">
					<BookMarked size={48} class="mx-auto text-indigo-600" />
				</div>
				<h3 class="text-2xl font-bold text-gray-900 mb-2">No notebooks yet</h3>
					<p class="text-gray-600 mb-6">Create your first notebook to start organizing your notes</p>
					<button
						onclick={() => showAddNotebook = true}
						class="bg-gradient-to-r from-purple-600 to-purple-700 text-white px-8 py-3 rounded-lg hover:shadow-lg transition font-medium inline-flex items-center gap-2"
					>
						<Plus size={20} /> Create First Notebook
					</button>
				</div>
			{:else if filteredNotebooks.length === 0}
				<div class="bg-gray-50 rounded-xl shadow-sm border border-gray-200 p-12 text-center">
					<p class="text-gray-600 text-lg">No notebooks match <span class="font-semibold">"{searchQuery}"</span></p>
					<button
						onclick={() => searchQuery = ''}
						class="mt-4 text-purple-600 hover:text-purple-700 font-medium"
					>
						Clear search
					</button>
				</div>
			{:else}
				<div class="space-y-8">
					<!-- Pinned Notebooks -->
					{#if pinnedNotebooks.length > 0}
						<div>
							<div class="flex items-center gap-2 mb-4">
								<Pin size={20} class="text-yellow-500" />
								<h2 class="text-xl font-bold text-gray-900">Pinned Notebooks ({pinnedNotebooks.length})</h2>
							</div>
							<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								{#each pinnedNotebooks as notebook (notebook.id)}
									<div class="group bg-gradient-to-br from-yellow-50 to-purple-50 rounded-xl shadow-md hover:shadow-lg transition border-2 border-yellow-300 overflow-hidden">
										<div class="p-6 space-y-4">
											<div class="flex justify-between items-start gap-2">
												<div class="flex-1 min-w-0">
													<h3 class="text-lg font-bold text-gray-900 truncate group-hover:text-indigo-600 transition">{notebook.name}</h3>
													<p class="text-sm text-gray-600 mt-1 flex items-center gap-1">
														<FileText size={16} /> {notebook.pages?.length || 0} pages
													</p>
												</div>
												<button
													onclick={(e) => { e.preventDefault(); handleTogglePin(notebook); }}
													class="p-2 hover:bg-yellow-200 rounded-lg transition"
													title="Unpin notebook"
												>
													<Pin size={20} class="text-yellow-600" />
												</button>
											</div>
											
											{#if notebook.tags}
												<div class="flex flex-wrap gap-2">
								{#each notebook.tags.split(',').filter((tag: string) => tag.trim()).slice(0, 3) as tag}
									<span class="text-xs bg-indigo-200 text-indigo-900 px-3 py-1 rounded-full font-medium flex items-center gap-1"><Tag size={12} /> {tag.trim()}</span>
													{/each}
													{#if notebook.tags.split(',').filter((tag: string) => tag.trim()).length > 3}
														<span class="text-xs bg-gray-200 text-gray-700 px-3 py-1 rounded-full">+{notebook.tags.split(',').filter((tag: string) => tag.trim()).length - 3}</span>
													{/if}
												</div>
											{/if}
											
											<div class="flex gap-2 pt-2">
												<a
													href="/myflowbook/{notebook.id}"
													class="flex-1 bg-purple-600 text-white text-center px-4 py-2 rounded-lg hover:bg-purple-700 transition font-medium text-sm"
												>
													Open
												</a>
												<button
													onclick={(e) => { e.preventDefault(); handleDeleteNotebook(notebook.id); }}
													class="bg-red-100 text-red-600 hover:bg-red-200 p-2 rounded-lg transition"
													title="Delete notebook"
												>
													<Trash2 size={18} />
												</button>
											</div>
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}
					
					<!-- Regular Notebooks -->
					{#if regularNotebooks.length > 0}
						<div>
							<h2 class="text-xl font-bold text-gray-900 mb-4">All Notebooks ({regularNotebooks.length})</h2>
							<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
								{#each regularNotebooks as notebook (notebook.id)}
									<div class="group bg-white rounded-xl shadow-sm hover:shadow-lg transition border border-gray-200 overflow-hidden">
										<div class="p-6 space-y-4">
											<div class="flex justify-between items-start gap-2">
												<div class="flex-1 min-w-0">
													<h3 class="text-lg font-bold text-gray-900 truncate group-hover:text-indigo-600 transition">{notebook.name}</h3>
													<p class="text-sm text-gray-600 mt-1 flex items-center gap-1">
														<FileText size={16} /> {notebook.pages?.length || 0} pages
													</p>
												</div>
												<button
													onclick={(e) => { e.preventDefault(); handleTogglePin(notebook); }}
													class="p-2 hover:bg-gray-100 rounded-lg transition text-gray-400 hover:text-yellow-500"
													title="Pin notebook"
												>
													<Pin size={20} />
												</button>
											</div>
											
											{#if notebook.tags}
												<div class="flex flex-wrap gap-2">
								{#each notebook.tags.split(',').filter((tag: string) => tag.trim()).slice(0, 3) as tag}
									<span class="text-xs bg-indigo-100 text-indigo-800 px-3 py-1 rounded-full font-medium flex items-center gap-1"><Tag size={12} /> {tag.trim()}</span>
													{/each}
													{#if notebook.tags.split(',').filter((tag: string) => tag.trim()).length > 3}
														<span class="text-xs bg-gray-100 text-gray-600 px-3 py-1 rounded-full">+{notebook.tags.split(',').filter((tag: string) => tag.trim()).length - 3}</span>
													{/if}
												</div>
											{/if}
											
											<div class="flex gap-2 pt-2">
												<a
													href="/myflowbook/{notebook.id}"
													class="flex-1 bg-purple-600 text-white text-center px-4 py-2 rounded-lg hover:bg-purple-700 transition font-medium text-sm"
												>
													Open
												</a>
												<button
													onclick={(e) => { e.preventDefault(); handleDeleteNotebook(notebook.id); }}
													class="bg-red-100 text-red-600 hover:bg-red-200 p-2 rounded-lg transition"
													title="Delete notebook"
												>
													<Trash2 size={18} />
												</button>
											</div>
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
