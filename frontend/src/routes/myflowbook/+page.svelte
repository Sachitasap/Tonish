<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { notebookAPI } from '$lib/api';
	import { wsService } from '$lib/websocket';
	import { Plus, Search, Pin, Trash2, BookOpen, X, ChevronRight } from 'lucide-svelte';

	type Notebook = {
		id: number;
		name: string;
		tags?: string;
		is_pinned?: boolean;
		pages?: { id: number }[];
	};

	const NOTEBOOK_COLORS: { key: string; bg: string; ring: string; label: string }[] = [
		{ key: 'blue',    bg: 'bg-blue-600',    ring: 'ring-blue-500',    label: 'Blue'    },
		{ key: 'violet',  bg: 'bg-violet-600',  ring: 'ring-violet-500',  label: 'Violet'  },
		{ key: 'emerald', bg: 'bg-emerald-600', ring: 'ring-emerald-500', label: 'Green'   },
		{ key: 'rose',    bg: 'bg-rose-600',    ring: 'ring-rose-500',    label: 'Rose'    },
		{ key: 'amber',   bg: 'bg-amber-500',   ring: 'ring-amber-400',   label: 'Amber'   },
		{ key: 'cyan',    bg: 'bg-cyan-600',    ring: 'ring-cyan-500',    label: 'Cyan'    },
		{ key: 'orange',  bg: 'bg-orange-600',  ring: 'ring-orange-500',  label: 'Orange'  },
		{ key: 'pink',    bg: 'bg-pink-600',    ring: 'ring-pink-500',    label: 'Pink'    },
	];
	const COLOR_BG_FALLBACK = [
		'bg-blue-600','bg-violet-600','bg-emerald-600','bg-rose-600',
		'bg-amber-500','bg-cyan-600','bg-orange-600','bg-pink-600',
	];

	function getNotebookColor(nb: Notebook): string {
		const colorTag = (nb.tags ?? '').split(',').map(t => t.trim()).find(t => t.startsWith('_clr:'));
		if (colorTag) {
			const key = colorTag.slice(5);
			return NOTEBOOK_COLORS.find(c => c.key === key)?.bg ?? COLOR_BG_FALLBACK[nb.id % COLOR_BG_FALLBACK.length];
		}
		return COLOR_BG_FALLBACK[nb.id % COLOR_BG_FALLBACK.length];
	}

	function visibleTags(nb: Notebook): string[] {
		return (nb.tags ?? '').split(',').map(t => t.trim()).filter(t => t && !t.startsWith('_clr:'));
	}
	let notebooks:    Notebook[] = $state([]);
	let loading     = $state(true);
	let showModal   = $state(false);
	let searchQuery = $state('');
	let toast       = $state('');
	let newNotebook = $state({ name: '', tags: '' });
	let selectedColor = $state('blue');

	const filtered = $derived(
		!searchQuery ? notebooks : notebooks.filter(n => {
			const q = searchQuery.toLowerCase();
			const tags = visibleTags(n).join(' ');
			return n.name.toLowerCase().includes(q) || tags.toLowerCase().includes(q);
		})
	);
	const pinned  = $derived(filtered.filter(n => n.is_pinned));
	const regular = $derived(filtered.filter(n => !n.is_pinned));

	// Group regular notebooks by their chosen color, preserving NOTEBOOK_COLORS order
	const colorGroups = $derived((() => {
		const groups = new Map<string, Notebook[]>();
		for (const c of NOTEBOOK_COLORS) groups.set(c.key, []);
		groups.set('_none', []); // fallback for notebooks without a color tag
		for (const nb of regular) {
			const tag = (nb.tags ?? '').split(',').map(t => t.trim()).find(t => t.startsWith('_clr:'));
			const key = tag ? tag.slice(5) : '_none';
			(groups.get(key) ?? groups.get('_none')!).push(nb);
		}
		// Return only non-empty groups in color order, then _none last
		const ordered: { key: string; label: string; bg: string; nbs: Notebook[] }[] = [];
		for (const c of NOTEBOOK_COLORS) {
			const nbs = groups.get(c.key) ?? [];
			if (nbs.length) ordered.push({ key: c.key, label: c.label, bg: c.bg, nbs });
		}
		const none = groups.get('_none') ?? [];
		if (none.length) ordered.push({ key: '_none', label: 'Other', bg: 'bg-gray-600', nbs: none });
		return ordered;
	})());

	onMount(() => {
		loadNotebooks();
		const refresh = () => loadNotebooks();
		wsService.on('notebook_update', refresh);
		wsService.on('notebook_create', refresh);
		wsService.on('notebook_delete', refresh);
		return () => {
			wsService.off('notebook_update', refresh);
			wsService.off('notebook_create', refresh);
			wsService.off('notebook_delete', refresh);
		};
	});

	async function loadNotebooks() {
		const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
		if (!token) { goto('/login'); return; }
		try {
			notebooks = await notebookAPI.getAll();
		} catch (error: any) {
			if (error.message?.includes('Authorization') || error.message?.includes('401')) {
				localStorage.removeItem('authToken');
				goto('/login');
			}
		} finally {
			loading = false;
		}
	}

	async function handleCreate() {
		if (!newNotebook.name.trim()) return;
		try {
			const colorTag = `_clr:${selectedColor}`;
			const baseTags = newNotebook.tags.trim();
			const fullTags = baseTags ? `${colorTag},${baseTags}` : colorTag;
			await notebookAPI.create({ ...newNotebook, tags: fullTags });
			await loadNotebooks();
			showModal = false;
			newNotebook = { name: '', tags: '' };
			selectedColor = 'blue';
		} catch (e) { console.error(e); }
	}

	async function handleDelete(id: number) {
		if (!confirm('Delete this notebook and all its pages?')) return;
		try {
			await notebookAPI.delete(id);
			await loadNotebooks();
			showToast('Notebook deleted');
		} catch (e) { console.error(e); }
	}

	async function handleTogglePin(nb: Notebook) {
		try {
			await notebookAPI.update(nb.id, { ...nb, is_pinned: !nb.is_pinned });
			await loadNotebooks();
		} catch (e) { console.error(e); }
	}

	function showToast(msg: string) {
		toast = msg;
		setTimeout(() => toast = '', 2500);
	}
</script>

<svelte:head>
	<title>MyFlowBook – Tonish</title>
</svelte:head>

<!-- Toast -->
{#if toast}
	<div class="fixed top-16 right-4 z-50 bg-gray-800 border border-gray-700 text-gray-200 text-xs px-3 py-2 rounded-lg shadow-lg animate-fade-in">
		{toast}
	</div>
{/if}

<div class="max-w-[1400px] mx-auto space-y-4">

	<!-- ── Page Header ── -->
	<div class="flex flex-col sm:flex-row sm:items-center gap-3">
		<div class="flex-1 min-w-0">
			<h1 class="text-lg font-bold text-white leading-tight">MyFlowBook</h1>
			<p class="text-xs text-gray-500 mt-0.5">
				{#if loading}&nbsp;{:else}{notebooks.length} notebook{notebooks.length !== 1 ? 's' : ''}{/if}
			</p>
		</div>

		<div class="flex items-center gap-2 w-full sm:w-auto">
			<div class="flex-1 sm:w-56 flex items-center gap-2 h-9 bg-gray-900 border border-gray-800 rounded-lg px-3 focus-within:border-gray-600 transition">
				<Search size={13} class="text-gray-500 flex-shrink-0" />
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search…"
					class="flex-1 bg-transparent text-sm text-white placeholder:text-gray-600 outline-none"
				/>
				{#if searchQuery}
					<button onclick={() => searchQuery = ''} class="text-gray-500 hover:text-gray-300 transition">
						<X size={13} />
					</button>
				{/if}
			</div>
			<button
				onclick={() => showModal = true}
				class="h-9 px-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium flex items-center gap-1.5 flex-shrink-0 transition"
			><Plus size={14} /> New</button>
		</div>
	</div>

	<!-- ── Content ── -->
	{#if loading}
		<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3">
			{#each [1,2,3,4,5,6] as _}
				<div class="h-28 bg-gray-900 border border-gray-800 rounded-xl animate-pulse"></div>
			{/each}
		</div>

	{:else if notebooks.length === 0}
		<div class="flex flex-col items-center justify-center py-24 text-center">
			<div class="w-16 h-16 bg-gray-900 border border-gray-800 rounded-2xl flex items-center justify-center mb-4">
				<BookOpen size={28} class="text-gray-600" />
			</div>
			<h2 class="text-base font-semibold text-gray-300 mb-1">No notebooks yet</h2>
			<p class="text-xs text-gray-600 mb-5 max-w-xs">Create your first notebook to start organising your notes and ideas.</p>
			<button
				onclick={() => showModal = true}
				class="h-9 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium flex items-center gap-2 transition"
			><Plus size={14} /> Create notebook</button>
		</div>

	{:else if filtered.length === 0}
		<div class="py-16 text-center">
			<Search size={28} class="text-gray-700 mx-auto mb-3" />
			<p class="text-sm text-gray-500">No notebooks match "<span class="text-gray-300">{searchQuery}</span>"</p>
			<button onclick={() => searchQuery = ''} class="mt-3 text-xs text-blue-400 hover:text-blue-300 transition">Clear search</button>
		</div>

	{:else}
		<div class="space-y-5">

			<!-- Pinned -->
			{#if pinned.length > 0}
				<section>
					<div class="flex items-center gap-2 mb-2.5">
						<Pin size={12} class="text-yellow-500 fill-current" />
						<span class="text-[11px] font-semibold text-gray-400 uppercase tracking-wider">Pinned</span>
					</div>
					<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3">
						{#each pinned as nb}
							{@render notebookCard(nb, true)}
						{/each}
					</div>
				</section>
			{/if}

			<!-- Color groups -->
			{#each colorGroups as group}
				<section>
					<div class="flex items-center gap-2 mb-2.5">
						<div class="w-2.5 h-2.5 rounded-full {group.bg} flex-shrink-0"></div>
						<span class="text-[11px] font-semibold text-gray-400 uppercase tracking-wider">{group.label}</span>
						<span class="text-[10px] text-gray-600">({group.nbs.length})</span>
					</div>
					<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3">
						{#each group.nbs as nb}
							{@render notebookCard(nb, false)}
						{/each}
					</div>
				</section>
			{/each}

		</div>
	{/if}
</div>

<!-- ── Notebook Card Snippet ── -->
{#snippet notebookCard(nb: Notebook, isPinned: boolean)}
	<div class="group relative bg-gray-900 border {isPinned ? 'border-yellow-800/60' : 'border-gray-800'} rounded-xl p-4 hover:border-gray-700 transition-all duration-150 flex flex-col gap-3">

		<!-- Top row -->
		<div class="flex items-start gap-3">
			<div class="flex-shrink-0 w-9 h-9 rounded-lg {getNotebookColor(nb)} flex items-center justify-center text-white font-bold text-sm shadow-sm select-none">
				{nb.name.charAt(0).toUpperCase()}
			</div>
			<div class="flex-1 min-w-0 pt-0.5">
				<a href="/myflowbook/{nb.id}" class="block">
					<h3 class="text-sm font-semibold text-white truncate hover:text-blue-400 transition leading-snug">{nb.name}</h3>
				</a>
				<p class="text-[11px] text-gray-500 mt-0.5">{nb.pages?.length ?? 0} page{(nb.pages?.length ?? 0) !== 1 ? 's' : ''}</p>
			</div>
			<button
				onclick={() => handleTogglePin(nb)}
				class="flex-shrink-0 p-1.5 rounded-md transition
					{isPinned
						? 'text-yellow-400 hover:bg-yellow-950'
						: 'text-gray-600 hover:text-yellow-400 hover:bg-gray-800 opacity-0 group-hover:opacity-100'}"
				title="{isPinned ? 'Unpin' : 'Pin'}"
			><Pin size={13} class="{isPinned ? 'fill-current' : ''}" /></button>
		</div>

		<!-- Tags (inline to avoid {@const} placement error) -->
		{#if visibleTags(nb).length > 0}
			<div class="flex flex-wrap gap-1">
				{#each visibleTags(nb).slice(0, 4) as tag}
					<span class="text-[10px] bg-gray-800 border border-gray-700 text-gray-400 px-1.5 py-0.5 rounded-full">{tag}</span>
				{/each}
			</div>
		{/if}

		<!-- Actions -->
		<div class="flex items-center gap-1.5 mt-auto pt-1 border-t border-gray-800">
			<a
				href="/myflowbook/{nb.id}"
				class="flex-1 flex items-center justify-center gap-1.5 h-8 bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-white rounded-lg text-xs font-medium transition"
			>Open <ChevronRight size={12} /></a>
			<button
				onclick={() => handleDelete(nb.id)}
				class="h-8 w-8 flex items-center justify-center bg-gray-800 hover:bg-red-950 text-gray-500 hover:text-red-400 rounded-lg transition opacity-0 group-hover:opacity-100"
				title="Delete"
			><Trash2 size={13} /></button>
		</div>
	</div>
{/snippet}

<!-- ── Create Modal ── -->
{#if showModal}
	<div
		class="fixed inset-0 z-50 bg-black/60 flex items-end sm:items-center justify-center sm:p-4"
		onclick={(e) => { if (e.target === e.currentTarget) { showModal = false; newNotebook = { name: '', tags: '' }; selectedColor = 'blue'; }}}
	>
		<div class="bg-gray-900 border border-gray-800 rounded-t-2xl sm:rounded-2xl shadow-2xl w-full sm:max-w-sm p-5 space-y-4">
			<div class="flex items-center justify-between">
				<h2 class="text-sm font-bold text-white">New Notebook</h2>
				<button
					onclick={() => { showModal = false; newNotebook = { name: '', tags: '' }; selectedColor = 'blue'; }}
					class="w-7 h-7 flex items-center justify-center text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition"
				><X size={16} /></button>
			</div>
			<form onsubmit={(e) => { e.preventDefault(); handleCreate(); }} class="space-y-3">
				<!-- Color preview strip -->
				<div class="flex items-center gap-3 p-3 bg-gray-800 rounded-lg border border-gray-700">
					<div class="w-10 h-10 rounded-lg {NOTEBOOK_COLORS.find(c=>c.key===selectedColor)?.bg ?? 'bg-blue-600'} flex items-center justify-center text-white font-bold text-base select-none flex-shrink-0">
						{newNotebook.name.charAt(0).toUpperCase() || '?'}
					</div>
					<div class="flex-1 min-w-0">
						<p class="text-xs text-gray-400">Preview · <span class="text-gray-300">{NOTEBOOK_COLORS.find(c=>c.key===selectedColor)?.label ?? 'Blue'}</span></p>
						<p class="text-[10px] text-gray-600 mt-0.5">This color will appear on your notebook card</p>
					</div>
				</div>

				<!-- Color swatches -->
				<div>
					<label class="block text-xs font-medium text-gray-400 mb-2">Color</label>
					<div class="grid grid-cols-8 gap-1.5">
						{#each NOTEBOOK_COLORS as c}
							<button
								type="button"
								onclick={() => selectedColor = c.key}
								title={c.label}
								class="w-full aspect-square rounded-lg {c.bg} transition-all
									{selectedColor === c.key ? `ring-2 ring-offset-2 ring-offset-gray-900 ${c.ring} scale-110` : 'opacity-70 hover:opacity-100 hover:scale-105'}"
							></button>
						{/each}
					</div>
				</div>

				<div>
					<label class="block text-xs font-medium text-gray-400 mb-1.5">Name <span class="text-red-500">*</span></label>
					<input
						type="text"
						bind:value={newNotebook.name}
						required
						placeholder="e.g. Project Ideas"
						class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-600 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
					/>
				</div>
				<div>
					<label class="block text-xs font-medium text-gray-400 mb-1.5">Tags <span class="text-gray-600 font-normal">(comma separated)</span></label>
					<input
						type="text"
						bind:value={newNotebook.tags}
						placeholder="e.g. work, personal"
						class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white placeholder:text-gray-600 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
					/>
				</div>
				<div class="flex gap-2 pt-1">
					<button type="submit" class="flex-1 h-10 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition">Create</button>
					<button type="button" onclick={() => { showModal = false; newNotebook = { name: '', tags: '' }; selectedColor = 'blue'; }} class="flex-1 h-10 bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg text-sm transition">Cancel</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Mobile FAB -->
<button
	onclick={() => showModal = true}
	class="md:hidden fixed bottom-24 right-4 z-40 w-12 h-12 bg-blue-600 hover:bg-blue-700 text-white rounded-full shadow-xl flex items-center justify-center transition"
	aria-label="New notebook"
><Plus size={22} /></button>

<style>
	@keyframes fade-in {
		from { opacity: 0; transform: translateY(-6px); }
		to   { opacity: 1; transform: translateY(0); }
	}
	.animate-fade-in { animation: fade-in 0.2s ease-out; }
</style>
