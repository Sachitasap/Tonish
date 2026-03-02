<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { notebookAPI, pageAPI } from '$lib/api';	import { Plus, Search, Pin, Download, Trash2, Edit2, Eye, Menu, X, Copy, ArrowLeft, Tag, Book, Lightbulb, ChevronDown, FileDown } from 'lucide-svelte';

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
	let showExportMenu = $state(false);

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
	let lastSavedAt = $state<string>('');

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
		lastSavedAt = '';
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
			lastSavedAt = new Date().toLocaleString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: 'numeric', minute: '2-digit' });
			await loadNotebook();
			showNotification = true;
			notificationMessage = '✓ Saved';
			notificationType = 'success';
			setTimeout(() => showNotification = false, 2000);
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
				.replace(/`(.*?)`/g, '<code class="bg-gray-700 text-gray-200 px-1 rounded">$1</code>')
			.replace(/^# (.*$)/gim, '<h1 class="text-2xl font-bold mt-4 mb-2">$1</h1>')
			.replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mt-3 mb-2">$1</h2>')
			.replace(/^### (.*$)/gim, '<h3 class="text-lg font-bold mt-2 mb-1">$1</h3>')
			.replace(/\n/g, '<br>');
	}
	
	// ── Export helpers ──────────────────────────────────────────
	function safeFileName(name: string) {
		return name.replace(/[^a-z0-9_\-]/gi, '_').replace(/_+/g, '_').replace(/^_|_$/g, '');
	}

	function downloadBlob(blob: Blob, filename: string) {
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = filename;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}

	function fireToast(msg: string, type: 'success' | 'error' = 'success') {
		showNotification = true;
		notificationMessage = msg;
		notificationType = type;
		showExportMenu = false;
		setTimeout(() => showNotification = false, 2500);
	}

	// ── Notebook-wide exports ─────────────────────────────────────
	function exportToText() {
		if (!notebook) return;
		const divider = '─'.repeat(48);
		let content = `Tonish – MyFlowBook Export\n${divider}\nNotebook : ${notebook.name}\nPages    : ${pages.length}\nExported : ${new Date().toLocaleString()}\n${divider}\n\n`;
		pages.forEach((p, i) => {
			content += `[${i + 1}] ${p.title}\n${'─'.repeat(p.title.length + 4)}\n\n`;
			content += `${p.content}\n\n`;
			if (p.tags) content += `Tags    : ${p.tags}\n`;
			content += `Created : ${formatDateTime(p.created_at)}\nUpdated : ${formatDateTime(p.updated_at)}\n\n${divider}\n\n`;
		});
		downloadBlob(new Blob([content], { type: 'text/plain;charset=utf-8' }), `${safeFileName(notebook.name)}.txt`);
		fireToast('✓ Exported as .TXT');
	}

	function exportToCSV() {
		if (!notebook) return;
		const esc = (s: string) => `"${(s ?? '').replace(/"/g, '""')}"`;
		let csv = `${esc('Notebook')},${esc('Title')},${esc('Content')},${esc('Tags')},${esc('Created')},${esc('Updated')}\n`;
		pages.forEach(p => {
			csv += `${esc(notebook!.name)},${esc(p.title)},${esc(p.content)},${esc(p.tags ?? '')},${esc(toIsoString(p.created_at))},${esc(toIsoString(p.updated_at))}\n`;
		});
		downloadBlob(new Blob([csv], { type: 'text/csv;charset=utf-8' }), `${safeFileName(notebook.name)}.csv`);
		fireToast('✓ Exported as .CSV');
	}

	function exportToMarkdown() {
		if (!notebook) return;
		let md = `# ${notebook.name}\n\n`;
		if (notebook.tags) md += `**Tags:** ${notebook.tags}\n\n`;
		md += `*${pages.length} page${pages.length !== 1 ? 's' : ''} · Exported ${new Date().toLocaleDateString()}*\n\n---\n\n`;
		pages.forEach(p => {
			md += `## ${p.title}\n\n${p.content}\n\n`;
			if (p.tags) md += `> **Tags:** ${p.tags}\n\n`;
			md += `<sub>Created: ${formatDateTime(p.created_at)} · Updated: ${formatDateTime(p.updated_at)}</sub>\n\n---\n\n`;
		});
		md += `\n*Exported from [Tonish](/) · ${new Date().toLocaleString()}*\n`;
		downloadBlob(new Blob([md], { type: 'text/markdown;charset=utf-8' }), `${safeFileName(notebook.name)}.md`);
		fireToast('✓ Exported as .MD');
	}

	function exportToJSON() {
		if (!notebook) return;
		const data = {
			exportedAt: new Date().toISOString(),
			exportedFrom: 'Tonish – MyFlowBook',
			notebook: { id: notebook.id, name: notebook.name, tags: notebook.tags, is_pinned: notebook.is_pinned },
			pages: pages.map(p => ({
				id: p.id, title: p.title, content: p.content,
				tags: p.tags, is_pinned: p.is_pinned,
				created_at: p.created_at, updated_at: p.updated_at,
			})),
		};
		downloadBlob(new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' }), `${safeFileName(notebook.name)}.json`);
		fireToast('✓ Exported as .JSON');
	}

	function exportToPDF() {
		if (!notebook) return;
		const win = window.open('', '_blank');
		if (!win) {
			fireToast('Popup blocked – allow popups to export PDF', 'error');
			return;
		}
		const exportDate = new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
		const renderContent = (text: string) => text
			.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
			.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
			.replace(/\*(.*?)\*/g, '<em>$1</em>')
			.replace(/`(.*?)`/g, '<code>$1</code>');
		let body = '';
		pages.forEach((p, i) => {
			body += `
				<div class="page-block">
					<div class="page-num">Page ${i + 1} of ${pages.length}</div>
					<h2 class="page-title">${p.title}</h2>
					<div class="page-meta">Created: ${formatDateOnly(p.created_at)} &nbsp;·&nbsp; Updated: ${formatDateOnly(p.updated_at)}</div>
					${p.tags ? `<div class="page-tags">${p.tags.split(',').filter(t => t.trim()).map(t => `<span>${t.trim()}</span>`).join('')}</div>` : ''}
					<div class="page-content">${renderContent(p.content)}</div>
				</div>
			`;
		});
		win.document.write(`<!DOCTYPE html>
<html lang="en"><head>
<meta charset="utf-8">
<title>${notebook.name} – Tonish Export</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; color: #1f2937; padding: 24mm 20mm 20mm; font-size: 14px; line-height: 1.6; }
  .export-header { display: flex; justify-content: space-between; align-items: flex-end; border-bottom: 3px solid #7c3aed; padding-bottom: 12px; margin-bottom: 28px; }
  .export-header .brand { font-size: 13px; font-weight: 700; color: #7c3aed; letter-spacing: .5px; }
  .export-header .date { font-size: 11px; color: #9ca3af; }
  .notebook-title { font-size: 30px; font-weight: 800; color: #111827; margin-bottom: 6px; }
  .notebook-subtitle { color: #6b7280; font-size: 13px; margin-bottom: 24px; }
  .divider { border: none; border-top: 1px solid #e5e7eb; margin: 0 0 24px; }
  .page-block { page-break-inside: avoid; margin-bottom: 32px; }
  .page-num { font-size: 10px; color: #d1d5db; text-transform: uppercase; letter-spacing: .8px; margin-bottom: 4px; }
  .page-title { font-size: 20px; font-weight: 700; color: #4f46e5; margin-bottom: 4px; }
  .page-meta { font-size: 11px; color: #9ca3af; margin-bottom: 8px; }
  .page-tags { margin-bottom: 10px; }
  .page-tags span { display: inline-block; background: #ede9fe; color: #5b21b6; border-radius: 999px; font-size: 11px; padding: 2px 10px; margin: 0 4px 4px 0; font-weight: 500; }
  .page-content { white-space: pre-wrap; color: #374151; line-height: 1.75; border-left: 3px solid #e5e7eb; padding-left: 14px; margin-top: 10px; }
  .page-content strong { font-weight: 700; }
  .page-content em { font-style: italic; }
  .page-content code { background: #f3f4f6; padding: 1px 6px; border-radius: 4px; font-family: 'Courier New', monospace; font-size: 12px; }
  .page-div { border: none; border-top: 1px dashed #e5e7eb; margin: 28px 0; }
  .export-footer { margin-top: 40px; border-top: 1px solid #e5e7eb; padding-top: 10px; text-align: center; color: #9ca3af; font-size: 11px; }
  @media print { @page { margin: 16mm 14mm; } body { padding: 0; } }
</style>
</head><body>
<div class="export-header"><span class="brand">Tonish · MyFlowBook</span><span class="date">Exported on ${exportDate}</span></div>
<h1 class="notebook-title">${notebook.name}</h1>
<p class="notebook-subtitle">${pages.length} page${pages.length !== 1 ? 's' : ''}${notebook.tags ? ' · Tags: ' + notebook.tags : ''}</p>
<hr class="divider">
${body}
<div class="export-footer">Exported from Tonish – ${new Date().toLocaleString()}</div>
</body></html>`);
		win.document.close();
		setTimeout(() => { win.focus(); win.print(); }, 300);
		fireToast('✓ Print dialog opened');
	}

	// ── Per-page export ───────────────────────────────────────────
	function exportCurrentPage(format: 'txt' | 'md') {
		if (!selectedPage || !notebook) return;
		let content = '';
		let mime = 'text/plain;charset=utf-8';
		if (format === 'txt') {
			content = `${selectedPage.title}\n${'='.repeat(selectedPage.title.length)}\n\n`;
			content += `${selectedPage.content}\n\n`;
			if (selectedPage.tags) content += `Tags    : ${selectedPage.tags}\n`;
			content += `Created : ${formatDateTime(selectedPage.created_at)}\nUpdated : ${formatDateTime(selectedPage.updated_at)}\n`;
			content += `\n─`.repeat(40) + `\nExported from Tonish – ${notebook.name}\n`;
		} else {
			mime = 'text/markdown;charset=utf-8';
			content = `# ${selectedPage.title}\n\n${selectedPage.content}\n\n`;
			if (selectedPage.tags) content += `> **Tags:** ${selectedPage.tags}\n\n`;
			content += `---\n<sub>Created: ${formatDateTime(selectedPage.created_at)} · Updated: ${formatDateTime(selectedPage.updated_at)}</sub>\n`;
			content += `\n*Exported from [Tonish](/) – ${notebook.name}*\n`;
		}
		downloadBlob(new Blob([content], { type: mime }), `${safeFileName(selectedPage.title)}.${format}`);
		fireToast(`✓ Page exported as .${format.toUpperCase()}`);
	}
</script>

<div class="flex flex-col lg:flex-row min-h-screen">
	<!-- Sidebar -->
	<div class={`${sidebarCollapsed ? 'lg:w-16' : 'lg:w-80'} w-full bg-gray-900 text-white transition-all duration-300 lg:flex-shrink-0 overflow-hidden border-r border-gray-800`}>
		<div class="p-3 h-full flex flex-col">
			<!-- Sidebar Header -->
			<div class="flex justify-between items-center mb-3 pb-2.5 border-b border-gray-800">
				{#if !sidebarCollapsed}
					<h2 class="text-xs font-bold text-white uppercase tracking-wide">Pages</h2>
				{/if}
				<button
					onclick={() => sidebarCollapsed = !sidebarCollapsed}
					class="text-gray-500 hover:bg-gray-700 p-1.5 rounded transition"
					title={sidebarCollapsed ? 'Expand' : 'Collapse'}
				>
					{#if sidebarCollapsed}
						<Menu size={18} />
					{:else}
						<X size={18} />
					{/if}
				</button>
			</div>
			
			{#if !sidebarCollapsed}
				<!-- Back to All Notebooks -->
				<a href="/myflowbook" class="flex items-center gap-1.5 px-2 h-7 hover:bg-gray-800 rounded transition mb-2.5 text-xs font-medium text-blue-400">
					<ArrowLeft size={13} />
					Back
				</a>
				
				<!-- Current Notebook Info -->
				{#if notebook}
					<div class="mb-3 p-2.5 bg-gray-800 rounded-lg border border-gray-700">
						<div class="flex items-start justify-between gap-2 mb-1.5">
							<h3 class="font-semibold text-xs text-white truncate">{notebook.name}</h3>
							<button
								onclick={handleTogglePin}
								class="hover:scale-110 transition p-1 rounded hover:bg-gray-700"
								title={notebook.is_pinned ? 'Unpin' : 'Pin'}
							>
								<Pin size={14} class={notebook.is_pinned ? 'fill-yellow-500 text-yellow-500' : 'text-gray-500'} />
							</button>
						</div>
						<p class="text-xs text-gray-400">{pages.length} {pages.length === 1 ? 'page' : 'pages'}</p>
					</div>
				{/if}
				
				<!-- Quick Actions -->
				<div class="mb-3 pb-2.5 border-b border-gray-800">
					<button
						onclick={() => showAddPage = !showAddPage}
						class="w-full flex items-center gap-1.5 h-8 px-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition text-xs font-medium"
					>
						<Plus size={13} />
						New Page
					</button>
				</div>
				
				<!-- Search -->
				<div class="mb-3 relative">
					<Search size={12} class="absolute left-2 top-2 text-gray-500" />
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search pages…"
						class="w-full h-7 pl-6 pr-2 bg-gray-800 border border-gray-700 rounded-md text-xs text-white placeholder:text-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					/>
				</div>
				
				<!-- Pages List -->
				<div class="flex-1 overflow-y-auto space-y-2">
					{#if searchQuery}
						{@const filteredPages = pages.filter(p => 
							p.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
							p.content.toLowerCase().includes(searchQuery.toLowerCase())
						)}
						{#each filteredPages as p}
							<button
								onclick={() => selectPage(p)}
								class={`w-full text-left p-2 rounded text-xs transition ${selectedPage?.id === p.id ? 'bg-blue-600 text-white font-semibold' : 'bg-gray-800 hover:bg-gray-700 border border-gray-700 text-gray-200'}`}
							>
								<div class="truncate">{p.title}</div>
							</button>
						{/each}
					{:else}
						{#if pages.filter(p => p.is_pinned).length > 0}
							<div class="mb-3">
								<h4 class="text-xs font-bold text-gray-400 mb-2 px-2">Pinned</h4>
								{#each pages.filter(p => p.is_pinned) as p}
									<button
										onclick={() => selectPage(p)}
										class={`w-full text-left p-2 rounded text-xs transition ${selectedPage?.id === p.id ? 'bg-yellow-900 text-yellow-200 font-semibold border border-yellow-700' : 'bg-gray-800 hover:bg-gray-700 border border-gray-700 text-gray-200'}`}
									>
										<div class="truncate">{p.title}</div>
									</button>
								{/each}
							</div>
						{/if}
						
						{#each pages.filter(p => !p.is_pinned) as p}
							<button
								onclick={() => selectPage(p)}
								class={`w-full text-left p-2 rounded text-xs transition ${selectedPage?.id === p.id ? 'bg-blue-600 text-white font-semibold' : 'bg-gray-800 hover:bg-gray-700 border border-gray-700 text-gray-200'}`}
							>
								<div class="truncate">{p.title}</div>
							</button>
						{/each}
					{/if}
				</div>
			{:else}
				<!-- Collapsed view - just icons -->
				<div class="flex flex-col items-center space-y-4">
					<a href="/myflowbook" class="p-3 hover:bg-gray-700 rounded transition text-gray-400" title="All Notebooks">
						<Book size={24} />
					</a>
					<button onclick={() => showAddPage = !showAddPage} class="p-3 hover:bg-gray-700 rounded transition text-gray-400" title="New Page">
						<Plus size={24} />
					</button>
					<div class="text-center text-xs text-gray-500 mt-4">
						{pages.length}
					</div>
				</div>
			{/if}
		</div>
	</div>
	
	<div class="flex-1 bg-gray-950 overflow-auto">
		<div class="p-3 sm:p-4 space-y-3">
			<!-- Notification -->
			{#if showNotification}
				<div class={`fixed top-14 right-3 px-3 py-2 rounded-lg shadow-xl z-50 text-xs font-medium border ${notificationType === 'success' ? 'bg-green-950 border-green-700 text-green-300' : 'bg-red-950 border-red-700 text-red-300'}`}>
					{notificationMessage}
				</div>
			{/if}

			<!-- Header -->
			<div class="flex items-center justify-between gap-3">
				<div class="flex items-center gap-2 min-w-0">
					{#if notebook}
						<h1 class="text-sm font-bold text-white leading-tight truncate">{notebook.name}</h1>
						<span class="text-[10px] text-gray-600 flex-shrink-0">{pages.length} page{pages.length !== 1 ? 's' : ''}</span>
					{/if}
				</div>
				<div class="flex items-center gap-1.5 flex-shrink-0">
					<!-- Export Dropdown -->
					<div class="relative">
						{#if showExportMenu}
							<div class="fixed inset-0 z-10" onclick={() => showExportMenu = false}></div>
						{/if}
						<button
							onclick={() => showExportMenu = !showExportMenu}
							class="h-8 px-2.5 bg-gray-800 border border-gray-700 text-gray-300 rounded-lg hover:bg-gray-700 transition inline-flex items-center gap-1 text-xs font-medium"
						>
							<Download size={12} />
							<span class="hidden sm:inline">Export</span>
							<ChevronDown size={11} class="text-gray-500 {showExportMenu ? 'rotate-180' : ''} transition-transform" />
						</button>
						{#if showExportMenu}
							<div class="absolute right-0 mt-1.5 w-48 bg-gray-800 border border-gray-700 rounded-xl shadow-2xl z-20 overflow-hidden py-1">
								<p class="px-3 py-1.5 text-[9px] font-semibold text-gray-600 uppercase tracking-wider">Whole Notebook</p>
								<button onclick={exportToPDF} class="w-full text-left px-3 py-2 hover:bg-gray-700 text-gray-300 flex items-center gap-2 transition text-xs">
									<span class="w-4 h-4 rounded bg-red-900 text-red-400 flex items-center justify-center flex-shrink-0 text-[9px] font-bold">P</span>PDF
								</button>
								<button onclick={exportToMarkdown} class="w-full text-left px-3 py-2 hover:bg-gray-700 text-gray-300 flex items-center gap-2 transition text-xs">
									<span class="w-4 h-4 rounded bg-blue-900 text-blue-400 flex items-center justify-center flex-shrink-0 text-[9px] font-bold">#</span>Markdown
								</button>
								<button onclick={exportToText} class="w-full text-left px-3 py-2 hover:bg-gray-700 text-gray-300 flex items-center gap-2 transition text-xs">
									<span class="w-4 h-4 rounded bg-gray-700 text-gray-300 flex items-center justify-center flex-shrink-0 text-[9px] font-bold">T</span>Text
								</button>
								<button onclick={exportToCSV} class="w-full text-left px-3 py-2 hover:bg-gray-700 text-gray-300 flex items-center gap-2 transition text-xs">
									<span class="w-4 h-4 rounded bg-emerald-900 text-emerald-400 flex items-center justify-center flex-shrink-0 text-[9px] font-bold">,</span>CSV
								</button>
								<button onclick={exportToJSON} class="w-full text-left px-3 py-2 hover:bg-gray-700 text-gray-300 flex items-center gap-2 transition text-xs">
									<span class="w-4 h-4 rounded bg-yellow-900 text-yellow-400 flex items-center justify-center flex-shrink-0 text-[9px] font-bold">&#123;&#125;</span>JSON
								</button>
							</div>
						{/if}
					</div>
					<button onclick={() => showAddPage = !showAddPage}
						class="h-8 px-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition inline-flex items-center gap-1.5 text-xs font-medium">
						<Plus size={13} /><span class="hidden sm:inline">New Page</span><span class="sm:hidden">Add</span>
					</button>
				</div>
			</div>
	
			<!-- Search Bar -->
			<div class="bg-gray-900 rounded-lg border border-gray-800 flex items-center gap-2 transition">
				<Search size={14} class="text-gray-500 ml-3 flex-shrink-0" />
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search pages by title or content…"
					id={pageSearchInputId}
					class="flex-1 h-9 px-2 bg-transparent border-none focus:outline-none text-sm text-white placeholder-gray-500"
					onkeydown={(e) => e.key === 'Enter' && handleSearch()}
				/>
				{#if searchQuery}
					<button
						onclick={() => { searchQuery = ''; searchResults = []; isSearching = false; }}
						class="mr-3 p-2 hover:bg-gray-800 rounded-lg transition text-gray-400"
						title="Clear search"
					>
						<X size={18} />
					</button>
				{/if}
			</div>
	
			<!-- search bar spacer only -->
	
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
				<div class="bg-gray-900 border border-gray-800 rounded-lg p-12 text-center">
					<div class="w-10 h-10 rounded-xl bg-purple-950 flex items-center justify-center mx-auto mb-3"><Plus size={18} class="text-purple-400" /></div>
					<p class="text-xs font-medium text-gray-400 mb-1">No pages yet</p>
					<p class="text-[11px] text-gray-600 mb-4">Create your first page to start taking notes</p>
					<button onclick={() => showAddPage = true} class="h-8 px-4 bg-purple-600 hover:bg-purple-700 text-white rounded-lg text-xs font-medium transition inline-flex items-center gap-1.5">
						<Plus size={13} /> New Page
					</button>
				</div>
				{:else if !showPageContent || !selectedPage}
				<div class="bg-gray-900 border border-gray-800 rounded-lg p-12 text-center">
					<div class="w-10 h-10 rounded-xl bg-gray-800 flex items-center justify-center mx-auto mb-3"><Book size={18} class="text-gray-500" /></div>
					<p class="text-xs font-medium text-gray-400 mb-1">No page selected</p>
					<p class="text-[11px] text-gray-600">Pick a page from the sidebar to start reading or editing</p>
				</div>
				{:else}
					<div class="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
						<!-- Page Header -->
						<div class="border-b border-gray-800 px-4 py-3">
							<div class="flex items-start justify-between gap-3">
								<!-- Title + meta -->
								<div class="min-w-0 flex-1">
									<h2 class="text-base font-bold text-white leading-snug break-words">{selectedPage.title}</h2>
									<div class="flex flex-wrap items-center gap-2 mt-1">
										<span class="text-[10px] text-gray-600">Created {formatDateOnly(selectedPage.created_at)}</span>
										<span class="text-gray-700 text-[10px]">·</span>
										{#if lastSavedAt}
											<span class="inline-flex items-center gap-1 text-[10px] bg-green-950 text-green-400 border border-green-900 px-1.5 py-0.5 rounded-full">
												<span class="w-1.5 h-1.5 rounded-full bg-green-500 inline-block"></span>
												Saved {lastSavedAt}
											</span>
										{:else}
											<span class="text-[10px] text-gray-600">Last updated {formatDateOnly(selectedPage.updated_at)}</span>
										{/if}
									</div>
									{#if selectedPage.tags}
										<div class="flex flex-wrap gap-1 mt-2">
											{#each selectedPage.tags.split(',').filter(t => t.trim()) as tag}
												<span class="text-[9px] bg-purple-950 text-purple-400 border border-purple-900 px-1.5 py-0.5 rounded-full">{tag.trim()}</span>
											{/each}
										</div>
									{/if}
								</div>
								<!-- Action buttons -->
								<div class="flex items-center gap-1 flex-shrink-0">
									<button onclick={copyPageContent} title="Copy content" class="h-7 w-7 flex items-center justify-center rounded-md bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white transition">
										<Copy size={13} />
									</button>
									<button onclick={() => exportCurrentPage('md')} title="Export as .md" class="h-7 px-2 flex items-center gap-1 rounded-md bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white transition text-[10px] font-mono">
										<FileDown size={11} />.md
									</button>
									<button onclick={() => exportCurrentPage('txt')} title="Export as .txt" class="h-7 px-2 flex items-center gap-1 rounded-md bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white transition text-[10px] font-mono">
										<FileDown size={11} />.txt
									</button>
									<button onclick={duplicatePage} title="Duplicate page" class="h-7 w-7 flex items-center justify-center rounded-md bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white transition">
										<Copy size={13} />
									</button>
									<button
										onclick={() => { isEditMode = !isEditMode; if (!isEditMode) editingContent = selectedPage!.content; }}
										class={`h-7 px-2.5 rounded-md text-xs font-medium transition inline-flex items-center gap-1 ${isEditMode ? 'bg-purple-950 text-purple-300 border border-purple-800' : 'bg-purple-600 text-white hover:bg-purple-700'}`}
									>
										{#if isEditMode}<Eye size={12} /> View{:else}<Edit2 size={12} /> Edit{/if}
									</button>
									<button onclick={() => selectedPage && handleDeletePage(selectedPage.id)} title="Delete page" class="h-7 w-7 flex items-center justify-center rounded-md bg-red-950 hover:bg-red-900 text-red-400 hover:text-red-300 transition">
										<Trash2 size={13} />
									</button>
								</div>
							</div>
						</div>

						<!-- Editor / Viewer -->
						<div class="p-4">
							{#if isEditMode}
								<!-- Compact toolbar -->
								<div class="mb-2 flex items-center gap-1 pb-2 border-b border-gray-800">
									<button onclick={() => formatText('bold')} title="Bold" class="h-7 px-2.5 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded text-xs font-bold text-gray-300 hover:text-white transition"><strong>B</strong></button>
									<button onclick={() => formatText('italic')} title="Italic" class="h-7 px-2.5 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded text-xs italic text-gray-300 hover:text-white transition"><em>I</em></button>
									<button onclick={() => formatText('code')} title="Code" class="h-7 px-2.5 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded text-xs font-mono text-gray-300 hover:text-white transition">&lt;/&gt;</button>
									<button onclick={() => formatText('heading')} title="Heading" class="h-7 px-2.5 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded text-xs font-bold text-gray-300 hover:text-white transition">H1</button>
									<div class="flex-1"></div>
									<span class="text-[10px] text-gray-600">Ctrl+S to save</span>
									<button onclick={handleSaveContent} class="h-7 px-3 bg-green-600 hover:bg-green-700 text-white rounded text-xs font-medium transition inline-flex items-center gap-1">
										<Download size={11} /> Save
									</button>
								</div>
								<textarea
									bind:value={editingContent}
									rows="22"
									class="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 text-gray-200 placeholder:text-gray-600 rounded-lg focus:outline-none focus:ring-1 focus:ring-purple-500 font-mono text-sm leading-relaxed resize-y"
									placeholder="Write your notes here… Markdown is supported: **bold**, *italic*, `code`, # Heading"
								></textarea>
							{:else}
								<!-- Read view -->
								<div class="prose prose-invert prose-sm max-w-none text-gray-300 leading-relaxed">
									{@html renderMarkdown(selectedPage.content)}
								</div>
								{#if !selectedPage.content}
									<p class="text-[11px] text-gray-600 italic text-center py-8">This page is empty — click Edit to add content</p>
								{/if}
							{/if}
						</div>
					</div>
				{/if}
			{/if}
		</div>
	</div>
</div>

<!-- ── New Page Modal ── -->
{#if showAddPage}
	<div
		class="fixed inset-0 z-50 bg-black/70 flex items-end md:items-center justify-center md:p-6"
		onclick={(e) => { if (e.target === e.currentTarget) { showAddPage = false; newPage = { title: '', content: '', tags: '' }; } }}
	>
		<div class="bg-gray-900 border border-gray-800 rounded-t-2xl md:rounded-2xl shadow-2xl w-full md:max-w-2xl flex flex-col max-h-[92vh] md:max-h-[85vh]">

			<!-- Modal header -->
			<div class="flex items-center justify-between px-5 py-4 border-b border-gray-800 flex-shrink-0">
				<div>
					<h2 class="text-sm font-bold text-white">New Page</h2>
					{#if notebook}<p class="text-[11px] text-gray-500 mt-0.5">in {notebook.name}</p>{/if}
				</div>
				<button
					onclick={() => { showAddPage = false; newPage = { title: '', content: '', tags: '' }; }}
					class="w-8 h-8 flex items-center justify-center text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition"
				><X size={16} /></button>
			</div>

			<!-- Scrollable body -->
			<div class="flex-1 overflow-y-auto">
				<form id="new-page-form" onsubmit={(e) => { e.preventDefault(); handleAddPage(); }} class="p-5 space-y-4">

					<!-- Quick templates -->
					<div>
						<p class="text-[10px] font-semibold text-gray-500 uppercase tracking-wider mb-2">Start from a template</p>
						<div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
							{#each [
								{ label: 'Blank page', icon: '📄', title: '', content: '' },
								{ label: 'Meeting notes', icon: '🗓️', title: 'Meeting Notes', content: '## Attendees\n\n## Agenda\n\n## Action Items\n' },
								{ label: 'Daily log', icon: '📝', title: 'Daily Log', content: '## What I did\n\n## Blockers\n\n## Tomorrow\n' },
								{ label: 'Code snippet', icon: '💻', title: 'Code Snippet', content: '## Description\n\n```\n// your code here\n```\n' },
							] as tpl}
								<button
									type="button"
									onclick={() => { newPage.title = tpl.title; newPage.content = tpl.content; }}
									class="flex items-center gap-2 h-9 px-3 bg-gray-800 hover:bg-gray-700 border border-gray-700 hover:border-gray-600 rounded-lg text-xs text-gray-300 hover:text-white transition"
								>
									<span>{tpl.icon}</span><span class="truncate">{tpl.label}</span>
								</button>
							{/each}
						</div>
					</div>

					<div class="border-t border-gray-800"></div>

					<!-- Title -->
					<div>
						<label for={addPageTitleId} class="block text-xs font-medium text-gray-400 mb-1.5">Title <span class="text-red-500">*</span></label>
						<input
							type="text"
							bind:value={newPage.title}
							required
							placeholder="Give this page a title…"
							id={addPageTitleId}
							class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-600 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500 transition"
						/>
					</div>

					<!-- Content -->
					<div>
						<div class="flex items-center justify-between mb-1.5">
							<label for={addPageContentId} class="text-xs font-medium text-gray-400">Content</label>
							<span class="text-[10px] text-gray-600">Supports **bold**, *italic*, `code`, # Heading</span>
						</div>
						<!-- Mini formatting toolbar -->
						<div class="flex gap-1 mb-1.5">
							{#each [
								{ label: 'B', title: 'Bold', wrap: '**', wrap2: '**' },
								{ label: 'I', title: 'Italic', wrap: '*', wrap2: '*' },
								{ label: '`', title: 'Code', wrap: '`', wrap2: '`' },
								{ label: 'H1', title: 'Heading', wrap: '# ', wrap2: '' },
							] as fmt}
								<button
									type="button"
									title={fmt.title}
									onclick={() => { newPage.content += `\n${fmt.wrap}text${fmt.wrap2}`; }}
									class="h-7 px-2.5 bg-gray-800 hover:bg-gray-700 border border-gray-700 rounded text-xs text-gray-400 hover:text-white font-mono transition"
								>{fmt.label}</button>
							{/each}
						</div>
						<textarea
							bind:value={newPage.content}
							rows="8"
							placeholder="Start writing… markdown is supported"
							id={addPageContentId}
							class="w-full px-3 py-2.5 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-600 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-purple-500 transition resize-y min-h-[160px]"
						></textarea>
					</div>

					<!-- Tags -->
					<div>
						<label for={addPageTagsId} class="block text-xs font-medium text-gray-400 mb-1.5">
							<Tag size={11} class="inline mr-1" />Tags <span class="text-gray-600 font-normal">(comma separated, optional)</span>
						</label>
						<input
							type="text"
							bind:value={newPage.tags}
							placeholder="e.g. important, work, research"
							id={addPageTagsId}
							class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-600 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500 transition"
						/>
					</div>

				</form>
			</div>

			<!-- Modal footer -->
			<div class="flex gap-2 px-5 py-4 border-t border-gray-800 flex-shrink-0">
				<button
					type="submit"
					form="new-page-form"
					class="flex-1 h-10 bg-purple-600 hover:bg-purple-700 text-white rounded-lg text-sm font-medium flex items-center justify-center gap-2 transition"
				><Plus size={15} /> Create Page</button>
				<button
					type="button"
					onclick={() => { showAddPage = false; newPage = { title: '', content: '', tags: '' }; }}
					class="flex-1 h-10 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg text-sm transition"
				>Cancel</button>
			</div>
		</div>
	</div>
{/if}

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
		color: #e5e7eb;
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
		background-color: #374151;
		color: #d1d5db;
		padding: 0.375rem 0.375rem;
		border-radius: 0.25rem;
		font-family: monospace;
		font-size: 0.875rem;
	}

	:global(.prose br) {
		margin-top: 0.5rem;
		margin-bottom: 0.5rem;
	}

	/* Custom scrollbar for sidebar */
	.custom-scrollbar::-webkit-scrollbar {
		width: 6px;
	}

	.custom-scrollbar::-webkit-scrollbar-track {
		background: rgba(139, 92, 246, 0.1);
		border-radius: 10px;
	}

	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: rgba(139, 92, 246, 0.5);
		border-radius: 10px;
	}

	.custom-scrollbar::-webkit-scrollbar-thumb:hover {
		background: rgba(139, 92, 246, 0.7);
	}
</style>
