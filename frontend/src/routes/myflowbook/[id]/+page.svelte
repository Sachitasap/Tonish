<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { notebookAPI, pageAPI } from '$lib/api';
	import { Plus, Search, Pin, Download, Trash2, Edit2, Eye, BookMarked, Menu, X, Copy, Clipboard, ArrowLeft, MapPin, Tag, Book, Lightbulb } from 'lucide-svelte';

	type NotebookPage = {
		id: number;
		title: string;
		content: string;
		tags?: string;
		is_pinned?: boolean;
		created_at?: string;
		updated_at?: string;
	};

	type Notebook = {
		id: number;
		name: string;
		tags?: string;
		is_pinned?: boolean;
		pages?: NotebookPage[];
	};

	let notebook = $state<Notebook | null>(null);
	let pages = $state<NotebookPage[]>([]);
	let selectedPage = $state<NotebookPage | null>(null);
	let loading = $state(true);
	let showAddPage = $state(false);
	let showPageContent = $state(false);
	let searchQuery = $state('');
	let searchResults = $state<NotebookPage[]>([]);
	let isSearching = $state(false);
	let sidebarCollapsed = $state(false);
	let showCopyNotification = $state(false);

	let newPage = $state<{ title: string; content: string; tags: string }>({
		title: '',
		content: '',
		tags: ''
	});

	let editingContent = $state('');
	let isEditMode = $state(false);
	let showNotification = $state(false);
	let notificationMessage = $state('');
	let notificationType = $state<'success' | 'error'>('success');

	const notebookId = $derived(parseInt($page.params.id ?? '0', 10));
	const pageSearchInputId = 'page-search-input';
	const addPageTitleId = 'new-page-title';
	const addPageContentId = 'new-page-content';
	const addPageTagsId = 'new-page-tags';

	const pagesByTag = $derived((() => {
		const tagMap = new Map<string, NotebookPage[]>();
		pages.forEach((pageItem) => {
			if (pageItem.tags) {
				pageItem.tags
					.split(',')
					.map((tag) => tag.trim())
					.filter((tag) => tag)
					.forEach((tag) => {
						if (!tagMap.has(tag)) {
							tagMap.set(tag, []);
						}
						tagMap.get(tag)?.push(pageItem);
					});
			}
		});
		return tagMap;
	})());

	const formatDateTime = (value?: string) => (value ? new Date(value).toLocaleString() : 'Unknown');
	const toIsoString = (value?: string) => (value ? new Date(value).toISOString() : '');
	const formatDateOnly = (value?: string) => (value ? new Date(value).toLocaleDateString() : 'Unknown');
	
	onMount(async () => {
		await loadNotebook();
	});
	
	async function loadNotebook() {
		try {
			const data = await notebookAPI.getOne(notebookId);
			notebook = data;
			pages = data.pages || [];
		} catch (error) {
			console.error('Failed to load notebook:', error);
		} finally {
			loading = false;
		}
	}
	
	async function handleTogglePin() {
		if (!notebook) return;
		try {
			const updatedNotebook: Notebook = {
				...notebook,
				is_pinned: !notebook.is_pinned
			};
			await notebookAPI.update(notebookId, updatedNotebook);
			notebook = updatedNotebook;
		} catch (error) {
			console.error('Failed to toggle pin:', error);
		}
	}
	
	async function handleAddPage() {
		try {
			await pageAPI.create({
				notebook_id: notebookId,
				title: newPage.title,
				content: newPage.content,
				tags: newPage.tags
			});
			await loadNotebook();
			showAddPage = false;
			newPage = { title: '', content: '', tags: '' };
		} catch (error) {
			console.error('Failed to create page:', error);
		}
	}
	
	async function handleDeletePage(pageId: number) {
		if (confirm('Are you sure you want to delete this page?')) {
			try {
				await pageAPI.delete(pageId);
				await loadNotebook();
				if (selectedPage?.id === pageId) {
					selectedPage = null;
					showPageContent = false;
				}
			} catch (error) {
				console.error('Failed to delete page:', error);
			}
		}
	}
	
	async function handleTogglePagePin(pageToPin: NotebookPage) {
		try {
			await pageAPI.update(pageToPin.id, {
				...pageToPin,
				is_pinned: !pageToPin.is_pinned
			});
			await loadNotebook();
		} catch (error) {
			console.error('Failed to toggle page pin:', error);
		}
	}
	
	function selectPage(p: NotebookPage) {
		selectedPage = p;
		editingContent = p.content;
		isEditMode = false;
		showPageContent = true;
	}
	
	async function handleSaveContent() {
		if (!selectedPage) return;
		try {
			const updatedPage: NotebookPage = {
				...selectedPage,
				content: editingContent
			};
			await pageAPI.update(selectedPage.id, updatedPage);
			selectedPage = updatedPage;
			isEditMode = false;
			await loadNotebook();
			showNotification = true;
			notificationMessage = 'Content saved successfully';
			notificationType = 'success';
			setTimeout(() => showNotification = false, 3000);
		} catch (error) {
			console.error('Failed to save content:', error);
			showNotification = true;
			notificationMessage = 'Failed to save content';
			notificationType = 'error';
			setTimeout(() => showNotification = false, 3000);
		}
	}
	
	function copyPageContent() {
		if (!selectedPage) return;
		
		const textToCopy = selectedPage.content;
		navigator.clipboard.writeText(textToCopy).then(() => {
			showNotification = true;
			notificationMessage = '✓ Content copied to clipboard';
			notificationType = 'success';
			setTimeout(() => {
				showNotification = false;
			}, 2000);
		}).catch(err => {
			console.error('Failed to copy:', err);
			showNotification = true;
			notificationMessage = 'Failed to copy content';
			notificationType = 'error';
			setTimeout(() => {
				showNotification = false;
			}, 2000);
		});
	}
	
	function duplicatePage() {
		if (!selectedPage) return;
		newPage = {
			title: `${selectedPage.title} (Copy)`,
			content: selectedPage.content,
			tags: selectedPage.tags || ''
		};
		showAddPage = true;
	}

	function focusPageSearchInput() {
		const searchInput = document.getElementById(pageSearchInputId) as HTMLInputElement | null;
		searchInput?.focus();
	}
	
	async function handleSearch() {
		if (!searchQuery.trim()) {
			searchResults = [];
			isSearching = false;
			return;
		}
		
		isSearching = true;
		try {
			const results = await pageAPI.search(searchQuery);
			searchResults = (results as NotebookPage[]) || [];
		} catch (error) {
			console.error('Search failed:', error);
			searchResults = [];
		}
	}
	
	// Text formatting helpers
	function formatText(command: string) {
		const selection = window.getSelection();
		if (!selection || selection.rangeCount === 0) return;
		
		const range = selection.getRangeAt(0);
		const selectedText = range.toString();
		
		let formatted = selectedText;
		switch(command) {
			case 'bold':
				formatted = `**${selectedText}**`;
				break;
			case 'italic':
				formatted = `*${selectedText}*`;
				break;
			case 'code':
				formatted = `\`${selectedText}\``;
				break;
			case 'heading':
				formatted = `# ${selectedText}`;
				break;
		}
		
		// For simplicity, we'll just append to content
		if (formatted !== selectedText) {
			editingContent = editingContent.replace(selectedText, formatted);
		}
	}
	
	function renderMarkdown(content: string) {
		if (!content) return '';
		// Basic markdown rendering
		return content
			.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
			.replace(/\*(.*?)\*/g, '<em>$1</em>')
			.replace(/`(.*?)`/g, '<code class="bg-gray-100 px-1 rounded">$1</code>')
			.replace(/^# (.*$)/gim, '<h1 class="text-2xl font-bold mt-4 mb-2">$1</h1>')
			.replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mt-3 mb-2">$1</h2>')
			.replace(/^### (.*$)/gim, '<h3 class="text-lg font-bold mt-2 mb-1">$1</h3>')
			.replace(/\n/g, '<br>');
	}
	
	// Export functions
	function exportToText() {
		if (!notebook) return;
		
		let content = `${notebook.name}\n${'='.repeat(notebook.name.length)}\n\n`;
		
		pages.forEach(page => {
			content += `\n${page.title}\n${'-'.repeat(page.title.length)}\n\n`;
			content += `${page.content}\n\n`;
			if (page.tags) {
				content += `Tags: ${page.tags}\n\n`;
			}
			content += `Created: ${formatDateTime(page.created_at)}\n`;
			content += `Updated: ${formatDateTime(page.updated_at)}\n\n`;
			content += '---\n';
		});
		
		const blob = new Blob([content], { type: 'text/plain' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${notebook.name.replace(/\s+/g, '_')}.txt`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}
	
	function exportToCSV() {
		if (!notebook) return;
		
		let csv = 'Title,Content,Tags,Created,Updated\n';
		
		pages.forEach(page => {
			const title = `"${page.title.replace(/"/g, '""')}"`;
			const content = `"${page.content.replace(/"/g, '""')}"`;
			const tags = `"${(page.tags || '').replace(/"/g, '""')}"`;
			const created = toIsoString(page.created_at);
			const updated = toIsoString(page.updated_at);
			
			csv += `${title},${content},${tags},${created},${updated}\n`;
		});
		
		const blob = new Blob([csv], { type: 'text/csv' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${notebook.name.replace(/\s+/g, '_')}.csv`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}
	
	function exportToMarkdown() {
		if (!notebook) return;
		
		let content = `# ${notebook.name}\n\n`;
		
		if (notebook.tags) {
			content += `**Tags:** ${notebook.tags}\n\n`;
		}
		
		content += `---\n\n`;
		
		pages.forEach(page => {
			content += `## ${page.title}\n\n`;
			content += `${page.content}\n\n`;
			if (page.tags) {
				content += `*Tags: ${page.tags}*\n\n`;
			}
			content += `<small>Created: ${formatDateTime(page.created_at)} | Updated: ${formatDateTime(page.updated_at)}</small>\n\n`;
			content += '---\n\n';
		});
		
		const blob = new Blob([content], { type: 'text/markdown' });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = `${notebook.name.replace(/\s+/g, '_')}.md`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}
</script>

<div class="flex flex-col lg:flex-row min-h-screen">
	<!-- Sidebar -->
	<div class={`${sidebarCollapsed ? 'lg:w-16' : 'lg:w-80'} w-full bg-gradient-to-b from-purple-900 to-purple-800 text-white transition-all duration-300 lg:flex-shrink-0 overflow-hidden`}>
		<div class="p-4 h-full flex flex-col">
			<!-- Sidebar Header -->
			<div class="flex justify-between items-center mb-6">
				{#if !sidebarCollapsed}
					<h2 class="text-lg font-bold flex items-center">
						<BookMarked size={20} class="mr-2" /> Navigation
					</h2>
				{/if}
				<button
					onclick={() => sidebarCollapsed = !sidebarCollapsed}
					class="text-white hover:bg-purple-700 p-2 rounded transition"
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
				<!-- Back to All Notebooks -->
				<a href="/myflowbook" class="flex items-center space-x-2 p-3 bg-purple-700 rounded-lg hover:bg-purple-600 transition mb-4">
					<ArrowLeft size={20} />
					<span>All Notebooks</span>
				</a>
				
				<!-- Current Notebook Info -->
				{#if notebook}
					<div class="mb-6 p-4 bg-purple-700/50 rounded-lg">
						<div class="flex items-start justify-between mb-2">
							<h3 class="font-bold text-lg">{notebook.name}</h3>
							<button
							onclick={handleTogglePin}
							class="hover:scale-110 transition"
						>
							<Pin size={18} class={notebook.is_pinned ? 'fill-current' : ''} />
							</button>
						</div>
						<p class="text-purple-200 text-sm">{pages.length} pages</p>
						{#if notebook.tags}
							<div class="flex flex-wrap gap-1 mt-2">
								{#each notebook.tags.split(',').filter((tag: string) => tag.trim()) as tag}
									<span class="text-xs bg-purple-600 px-2 py-0.5 rounded">{tag.trim()}</span>
								{/each}
							</div>
						{/if}
					</div>
				{/if}
				
					<!-- Quick Actions -->
				<div class="mb-6">
					<h3 class="text-sm font-semibold text-purple-300 mb-3 uppercase">Quick Actions</h3>
					<div class="space-y-2">
						<button
							onclick={() => showAddPage = !showAddPage}
							class="w-full flex items-center space-x-2 p-2 bg-purple-700 rounded hover:bg-purple-600 transition"
						>
							<Plus size={16} />
							<span>New Page</span>
						</button>
						<button
							onclick={focusPageSearchInput}
							class="w-full flex items-center space-x-2 p-2 bg-purple-700 rounded hover:bg-purple-600 transition"
						>
							<Search size={16} />
							<span>Search</span>
						</button>
					</div>
				</div>				<!-- Pages organized by status -->
				<div class="flex-1 overflow-y-auto">
					<h3 class="text-sm font-semibold text-purple-300 mb-3 uppercase">Pages</h3>
					
					<!-- Pinned Pages -->
					{#if pages.filter(p => p.is_pinned).length > 0}
						<div class="mb-4">
							<h4 class="text-xs text-purple-400 mb-2 flex items-center">
								<Pin size={12} class="mr-1" /> Pinned ({pages.filter(p => p.is_pinned).length})
							</h4>
							<div class="space-y-1">
								{#each pages.filter(p => p.is_pinned) as p}
									<button
										onclick={() => selectPage(p)}
										class={`w-full text-left p-2 rounded transition ${selectedPage?.id === p.id ? 'bg-yellow-500 text-purple-900' : 'bg-purple-700/50 hover:bg-purple-600'}`}
									>
										<div class="font-medium text-sm truncate">{p.title}</div>
										<div class="text-xs text-purple-200 truncate">{p.content?.substring(0, 30)}...</div>
									</button>
								{/each}
							</div>
						</div>
					{/if}
					
					<!-- Recent Pages (unpinned) -->
					{#if pages.filter(p => !p.is_pinned).length > 0}
						<div class="mb-4">
							<h4 class="text-xs text-purple-400 mb-2">Recent Pages</h4>
							<div class="space-y-1">
								{#each pages.filter(p => !p.is_pinned).slice(0, 10) as p}
									<button
										onclick={() => selectPage(p)}
										class={`w-full text-left p-2 rounded transition ${selectedPage?.id === p.id ? 'bg-purple-500' : 'bg-purple-700/50 hover:bg-purple-600'}`}
									>
										<div class="font-medium text-sm truncate">{p.title}</div>
										<div class="text-xs text-purple-200 truncate">{p.content?.substring(0, 30)}...</div>
									</button>
								{/each}
							</div>
						</div>
					{/if}
					
					<!-- Pages by Tag -->
					{#if pagesByTag.size > 0}
						<div class="mb-4">
							<h4 class="text-xs text-purple-400 mb-2">By Tag</h4>
							<div class="space-y-2">
								{#each Array.from(pagesByTag.entries()) as [tag, tagPages]}
									<details class="group">
									<summary class="cursor-pointer p-2 bg-purple-700/50 rounded hover:bg-purple-600 transition text-sm flex justify-between items-center">
										<span class="flex items-center gap-2"><Tag size={16} /> {tag}</span>
										<span class="text-xs text-purple-300">{tagPages.length}</span>
									</summary>
										<div class="ml-2 mt-1 space-y-1">
											{#each tagPages as p}
												<button
													onclick={() => selectPage(p)}
													class={`w-full text-left p-1.5 rounded transition text-sm ${selectedPage?.id === p.id ? 'bg-purple-500' : 'bg-purple-700/30 hover:bg-purple-600'}`}
												>
													{p.title}
												</button>
											{/each}
										</div>
									</details>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			{:else}
				<!-- Collapsed view - just icons -->
				<div class="flex flex-col items-center space-y-4">
					<a href="/myflowbook" class="p-3 hover:bg-purple-700 rounded transition" title="All Notebooks">
						<Book size={24} />
					</a>
					<button onclick={() => showAddPage = !showAddPage} class="p-3 hover:bg-purple-700 rounded transition" title="New Page">
						<Plus size={24} />
					</button>
					<div class="text-center text-xs text-purple-300 mt-4">
						{pages.length}
					</div>
				</div>
			{/if}
		</div>
	</div>
	
	<!-- Main Content -->
	<div class="flex-1 bg-gradient-to-b from-gray-50 to-white overflow-auto">
		<div class="max-w-5xl mx-auto p-6 space-y-6">
			<!-- Notification -->
			{#if showNotification}
				<div class={`fixed top-4 right-4 px-6 py-3 rounded-lg shadow-lg z-50 animate-slide-in font-medium ${notificationType === 'success' ? 'bg-green-500 text-white' : 'bg-red-500 text-white'}`}>
					{notificationMessage}
				</div>
			{/if}
			
			<!-- Header -->
			<div class="flex justify-between items-center">
				<div class="flex items-center space-x-4">
					{#if notebook}
						<h1 class="text-3xl font-bold text-gray-900">{notebook.name}</h1>
					{/if}
				</div>
				<div class="flex flex-wrap gap-2 items-center">
					<!-- Export Dropdown -->
					<div class="relative group">
						<button
							class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition inline-flex items-center gap-2 font-medium"
						>
							<Download size={18} />
							<span class="hidden sm:inline">Export</span>
							<span>▼</span>
						</button>
						<div class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-xl hidden group-hover:block z-10 border border-gray-200 overflow-hidden">
							<button
								onclick={exportToText}
								class="w-full text-left px-4 py-3 hover:bg-gray-50 text-gray-900 flex items-center gap-2 transition border-b border-gray-100"
							>
								<Download size={16} />
								Export as Text
							</button>
							<button
								onclick={exportToCSV}
								class="w-full text-left px-4 py-3 hover:bg-gray-50 text-gray-900 flex items-center gap-2 transition border-b border-gray-100"
							>
								<Download size={16} />
								Export as CSV
							</button>
							<button
								onclick={exportToMarkdown}
								class="w-full text-left px-4 py-3 hover:bg-gray-50 text-gray-900 flex items-center gap-2 transition"
							>
								<Download size={16} />
								Export as Markdown
							</button>
						</div>
					</div>
					<button
						onclick={() => showAddPage = !showAddPage}
						class="bg-gradient-to-r from-purple-600 to-purple-700 text-white px-4 py-2 rounded-lg hover:shadow-lg transition inline-flex items-center gap-2 font-medium"
					>
						<Plus size={18} />
						<span class="hidden sm:inline">New Page</span>
					</button>
				</div>
			</div>
	
			<!-- Search Bar -->
			<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-1 flex items-center gap-2 hover:shadow-md transition">
				<Search size={18} class="text-gray-400 ml-3" />
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search pages by title or content..."
					id={pageSearchInputId}
					class="flex-1 px-3 py-3 bg-white border-none focus:outline-none text-gray-900 placeholder-gray-400"
					onkeydown={(e) => e.key === 'Enter' && handleSearch()}
				/>
				{#if searchQuery}
					<button
						onclick={() => { searchQuery = ''; searchResults = []; isSearching = false; }}
						class="mr-3 p-2 hover:bg-gray-100 rounded-lg transition text-gray-400"
						title="Clear search"
					>
						<X size={18} />
					</button>
				{/if}
			</div>
	
			<!-- Add Page Form -->
			{#if showAddPage}
				<div class="bg-white rounded-lg shadow-lg border border-purple-100 p-6 animate-scale-in">
					<h2 class="text-2xl font-bold text-gray-900 mb-6">Create New Page</h2>
					<form onsubmit={(e) => { e.preventDefault(); handleAddPage(); }} class="space-y-4">
						<div>
							<label for={addPageTitleId} class="block text-sm font-semibold text-gray-700 mb-2">Page Title</label>
							<input
								type="text"
								bind:value={newPage.title}
								required
								placeholder="e.g., Meeting Notes, Research, Ideas"
								id={addPageTitleId}
								class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition"
							/>
						</div>
						
						<div>
							<label for={addPageContentId} class="block text-sm font-semibold text-gray-700 mb-2">Content</label>
							<textarea
								bind:value={newPage.content}
								rows="6"
								placeholder="Start writing... Supports markdown: **bold**, *italic*, `code`, # headings"
								id={addPageContentId}
								class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition font-mono text-sm"
							></textarea>
						</div>
						
						<div>
							<label for={addPageTagsId} class="block text-sm font-semibold text-gray-700 mb-2">Tags (comma-separated, optional)</label>
							<input
								type="text"
								bind:value={newPage.tags}
								placeholder="e.g., important, work, research"
								id={addPageTagsId}
								class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition"
							/>
						</div>
						
						<div class="flex gap-2 pt-2">
							<button
								type="submit"
								class="flex-1 bg-gradient-to-r from-purple-600 to-purple-700 text-white px-4 py-3 rounded-lg hover:shadow-md transition font-medium"
							>
								Create Page
							</button>
							<button
								type="button"
								onclick={() => { showAddPage = false; newPage = { title: '', content: '', tags: '' }; }}
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
						<div class="w-2 h-2 bg-purple-600 rounded-full animate-pulse"></div>
						<p>Loading notebook...</p>
					</div>
				</div>
			{:else}
				<!-- Page Content Viewer -->
				{#if pages.length === 0}
				<div class="bg-gradient-to-br from-purple-50 to-blue-50 rounded-lg shadow-sm border border-purple-100 p-16 text-center">
					<div class="mb-4 flex justify-center"><Pin size={64} class="text-purple-400" /></div>
					<p class="text-gray-600 text-lg mb-2">No pages yet</p>
						<p class="text-gray-500 mb-6">Create your first page to start taking notes</p>
						<button
							onclick={() => showAddPage = true}
							class="bg-gradient-to-r from-purple-600 to-purple-700 text-white px-8 py-3 rounded-lg hover:shadow-lg transition font-medium inline-flex items-center gap-2"
						>
							<Plus size={20} /> Create First Page
						</button>
					</div>
				{:else if !showPageContent || !selectedPage}
				<div class="bg-gradient-to-br from-gray-50 to-white rounded-lg shadow-sm border border-gray-200 p-16 text-center">
					<div class="mb-4 flex justify-center"><ArrowLeft size={64} class="text-gray-400" /></div>
					<p class="text-gray-600 text-lg">Select a page from the sidebar to start editing</p>
					</div>
				{:else}
					<div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
						<!-- Page Header -->
						<div class="border-b border-gray-200 p-6 space-y-4 bg-gray-50">
							<div class="flex justify-between items-start gap-3 flex-wrap">
								<div class="flex-1">
									<h2 class="text-3xl font-bold text-gray-900 break-words">{selectedPage.title}</h2>
									<p class="text-sm text-gray-600 mt-2">
										Created: {formatDateOnly(selectedPage.created_at)} • Updated: {formatDateOnly(selectedPage.updated_at)}
									</p>
								</div>
								<div class="flex flex-wrap gap-2">
									<button
										onclick={copyPageContent}
										class="bg-blue-100 text-blue-700 hover:bg-blue-200 px-4 py-2 rounded-lg transition inline-flex items-center gap-2 font-medium text-sm"
										title="Copy content to clipboard"
									>
										<Copy size={16} />
										<span class="hidden sm:inline">Copy</span>
									</button>
									<button
										onclick={duplicatePage}
										class="bg-indigo-100 text-indigo-700 hover:bg-indigo-200 px-4 py-2 rounded-lg transition inline-flex items-center gap-2 font-medium text-sm"
										title="Duplicate this page"
									>
										<Copy size={16} />
										<span class="hidden sm:inline">Duplicate</span>
									</button>
									<button
										onclick={() => isEditMode = !isEditMode}
										class={`px-4 py-2 rounded-lg transition inline-flex items-center gap-2 font-medium text-sm ${isEditMode ? 'bg-purple-100 text-purple-700 hover:bg-purple-200' : 'bg-purple-600 text-white hover:bg-purple-700'}`}
									>
										{#if isEditMode}
											<Eye size={16} />
											<span class="hidden sm:inline">Preview</span>
										{:else}
											<Edit2 size={16} />
											<span class="hidden sm:inline">Edit</span>
										{/if}
									</button>
									<button
										onclick={() => selectedPage && handleDeletePage(selectedPage.id)}
										class="bg-red-100 text-red-600 hover:bg-red-200 px-4 py-2 rounded-lg transition inline-flex items-center gap-2 font-medium text-sm"
									>
										<Trash2 size={16} />
										<span class="hidden sm:inline">Delete</span>
									</button>
								</div>
							</div>
							
							{#if selectedPage.tags}
								<div class="flex flex-wrap gap-2 pt-2">
									{#each selectedPage.tags.split(',').filter((tag: string) => tag.trim()) as tag}
										<span class="text-xs bg-purple-100 text-purple-800 px-3 py-1 rounded-full font-medium">{tag.trim()}</span>
									{/each}
								</div>
							{/if}
						</div>
						
						<!-- Rich Text Editor / Viewer -->
						<div class="p-6">
							{#if isEditMode}
								<!-- Toolbar -->
								<div class="mb-4 flex flex-wrap gap-2 border-b pb-4">
									<button
										onclick={() => formatText('bold')}
										class="px-3 py-2 bg-gray-100 text-gray-700 hover:bg-gray-200 rounded-lg transition font-semibold text-sm"
										title="Bold (select text)"
									>
										<strong>B</strong>
									</button>
									<button
										onclick={() => formatText('italic')}
										class="px-3 py-2 bg-gray-100 text-gray-700 hover:bg-gray-200 rounded-lg transition font-semibold italic text-sm"
										title="Italic (select text)"
									>
										<em>I</em>
									</button>
									<button
										onclick={() => formatText('code')}
										class="px-3 py-2 bg-gray-100 text-gray-700 hover:bg-gray-200 rounded-lg transition font-mono text-sm"
										title="Code (select text)"
									>
										&lt;/&gt;
									</button>
									<button
										onclick={() => formatText('heading')}
										class="px-3 py-2 bg-gray-100 text-gray-700 hover:bg-gray-200 rounded-lg transition font-bold text-sm"
										title="Heading (select text)"
									>
										H1
									</button>
									<div class="flex-1"></div>
									<button
										onclick={handleSaveContent}
										class="px-4 py-2 bg-green-600 text-white hover:bg-green-700 rounded-lg transition inline-flex items-center gap-2 font-medium text-sm"
									>
										<Download size={16} />
										<span class="hidden sm:inline">Save</span>
									</button>
								</div>
								
								<!-- Editor -->
								<textarea
									bind:value={editingContent}
									rows="20"
									class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 font-mono text-sm bg-white"
									placeholder="Write your content here... (Supports markdown)"
								></textarea>
								
							<div class="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
								<p class="text-sm text-blue-900 flex items-center gap-2">
									<Lightbulb size={16} class="text-blue-600" /> <strong>Tip:</strong> Use <code class="bg-blue-100 px-1 rounded">**text**</code> for bold,
										<code class="bg-blue-100 px-1 rounded">*text*</code> for italic, 
										<code class="bg-blue-100 px-1 rounded">`code`</code> for inline code, 
										and <code class="bg-blue-100 px-1 rounded"># Heading</code> for headings
									</p>
								</div>
							{:else}
								<!-- Preview -->
								<div class="prose prose-sm max-w-none text-gray-800 leading-relaxed">
									{@html renderMarkdown(selectedPage.content)}
								</div>
							{/if}
						</div>
					</div>
				{/if}
			{/if}
		</div>
	</div>
</div>

<style>
	@keyframes fade-in {
		from { opacity: 0; transform: translateY(-10px); }
		to { opacity: 1; transform: translateY(0); }
	}
	
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

	.animate-fade-in {
		animation: fade-in 0.3s ease-out;
	}

	.animate-slide-in {
		animation: slide-in 0.3s ease-out;
	}

	.animate-scale-in {
		animation: scale-in 0.2s ease-out;
	}

	/* Prose styling for markdown preview */
	:global(.prose) {
		color: #374151;
	}

	:global(.prose strong) {
		font-weight: bold;
	}

	:global(.prose em) {
		font-style: italic;
	}

	:global(.prose h1) {
		font-size: 1.875rem;
		font-weight: bold;
		margin-top: 1rem;
		margin-bottom: 1rem;
	}

	:global(.prose h2) {
		font-size: 1.25rem;
		font-weight: bold;
		margin-top: 0.75rem;
		margin-bottom: 0.75rem;
	}

	:global(.prose h3) {
		font-size: 1.125rem;
		font-weight: bold;
		margin-top: 0.5rem;
		margin-bottom: 0.5rem;
	}

	:global(.prose code) {
		background-color: #f3f4f6;
		padding: 0.375rem 0.375rem;
		border-radius: 0.25rem;
		font-family: monospace;
		font-size: 0.875rem;
	}

	:global(.prose br) {
		margin-top: 0.5rem;
		margin-bottom: 0.5rem;
	}
</style>
