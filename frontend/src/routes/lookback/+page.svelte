<script lang="ts">
	import { onMount } from 'svelte';
	import { taskAPI } from '$lib/api';
	import { Calendar, CalendarDays, RotateCcw, Trash2, Archive, Kanban, LayoutGrid } from 'lucide-svelte';

	interface Task {
		id: number;
		title: string;
		description: string;
		task_type: 'kanban' | 'matrix' | 'calendar';
		status: string;
		quadrant: string;
		priority: string;
		due_date: string | null;
		is_archived: boolean;
		completed_at: string | null;
		deleted_at: string | null;
		created_at: string;
		updated_at: string;
	}

	let archivedTasks = $state<Task[]>([]);
	let loading = $state(true);
	let filterType = $state<'all' | 'kanban' | 'matrix' | 'calendar'>('all');
	let sortBy = $state<'date' | 'type'>('date');
	let selectedTaskIds = $state<Set<number>>(new Set());
	let isSelectAllChecked = $state(false);

	onMount(async () => {
		await loadArchivedTasks();
	});

	async function loadArchivedTasks() {
		try {
			loading = true;
			archivedTasks = await taskAPI.getArchived();
		} catch (error) {
			console.error('Failed to load archived tasks:', error);
		} finally {
			loading = false;
		}
	}

	async function handleRestore(taskId: number) {
		if (confirm('Are you sure you want to restore this task?')) {
			try {
				await taskAPI.restore(taskId);
				await loadArchivedTasks();
			} catch (error) {
				console.error('Failed to restore task:', error);
			}
		}
	}

	async function handlePermanentDelete(taskId: number) {
		if (confirm('⚠️ This will permanently delete the task. This action cannot be undone. Are you sure?')) {
			try {
				await taskAPI.permanentDelete(taskId);
				await loadArchivedTasks();
			} catch (error) {
				console.error('Failed to permanently delete task:', error);
			}
		}
	}

	function toggleTaskSelection(taskId: number) {
		const newSet = new Set(selectedTaskIds);
		if (newSet.has(taskId)) {
			newSet.delete(taskId);
		} else {
			newSet.add(taskId);
		}
		selectedTaskIds = newSet;
		updateSelectAllState();
	}

	function toggleSelectAll() {
		isSelectAllChecked = !isSelectAllChecked;
		if (isSelectAllChecked) {
			selectedTaskIds = new Set(filteredTasks().map(t => t.id));
		} else {
			selectedTaskIds = new Set();
		}
	}

	function updateSelectAllState() {
		const filtered = filteredTasks();
		if (filtered.length === 0) {
			isSelectAllChecked = false;
		} else {
			isSelectAllChecked = filtered.every(t => selectedTaskIds.has(t.id));
		}
	}

	async function handleBulkDelete() {
		const count = selectedTaskIds.size;
		if (count === 0) return;
		
		if (confirm(`⚠️ This will permanently delete ${count} task${count > 1 ? 's' : ''}. This action cannot be undone. Are you sure?`)) {
			try {
				const deletePromises = Array.from(selectedTaskIds).map(id => taskAPI.permanentDelete(id));
				await Promise.all(deletePromises);
				selectedTaskIds = new Set();
				isSelectAllChecked = false;
				await loadArchivedTasks();
			} catch (error) {
				console.error('Failed to delete tasks:', error);
			}
		}
	}

	function clearSelection() {
		selectedTaskIds = new Set();
		isSelectAllChecked = false;
	}

	function formatDate(dateString: string | null): string {
		if (!dateString) return 'N/A';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { 
			year: 'numeric', 
			month: 'short', 
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getQuadrantLabel(quadrant: string): string {
		const labels: Record<string, string> = {
			'urgent-important': 'Urgent & Important',
			'not-urgent-important': 'Not Urgent & Important',
			'urgent-not-important': 'Urgent & Not Important',
			'not-urgent-not-important': 'Not Urgent & Not Important'
		};
		return labels[quadrant] || 'N/A';
	}

	const filteredTasks = $derived(() => {
		let filtered = archivedTasks;
		
		if (filterType !== 'all') {
			filtered = filtered.filter(t => t.task_type === filterType);
		}
		
		if (sortBy === 'date') {
			filtered = [...filtered].sort((a, b) => {
				const dateA = new Date(a.deleted_at || a.completed_at || a.updated_at).getTime();
				const dateB = new Date(b.deleted_at || b.completed_at || b.updated_at).getTime();
				return dateB - dateA;
			});
		} else {
			filtered = [...filtered].sort((a, b) => a.task_type.localeCompare(b.task_type));
		}
		
		return filtered;
	});

	const kanbanTasks = $derived(filteredTasks().filter(t => t.task_type === 'kanban'));
	const matrixTasks = $derived(filteredTasks().filter(t => t.task_type === 'matrix'));
	const calendarTasks = $derived(filteredTasks().filter(t => t.task_type === 'calendar'));
</script>

<svelte:head>
	<title>LookBack - Tonish</title>
</svelte:head>

<div class="max-w-[1600px] mx-auto">

	<!-- Page Header -->
	<div class="flex flex-wrap items-center justify-between gap-2 mb-3">
		<div class="flex items-center gap-2">
			<div class="w-7 h-7 rounded-lg bg-purple-950 border border-purple-800 flex items-center justify-center flex-shrink-0">
				<Archive size={13} class="text-purple-400" />
			</div>
			<div>
				<h1 class="text-sm font-bold text-white leading-tight">LookBack</h1>
				<p class="text-[10px] text-gray-500">Archived &amp; completed tasks</p>
			</div>
		</div>
		<!-- Quick stats -->
		<div class="flex items-center gap-1.5">
			<div class="flex items-center gap-1 bg-gray-900 border border-gray-800 rounded-md px-2 py-1">
				<Kanban size={10} class="text-blue-400" />
				<span class="text-[10px] text-gray-400">Kanban</span>
				<span class="text-[10px] font-bold text-white ml-0.5">{kanbanTasks.length}</span>
			</div>
			<div class="flex items-center gap-1 bg-gray-900 border border-gray-800 rounded-md px-2 py-1">
				<LayoutGrid size={10} class="text-green-400" />
				<span class="text-[10px] text-gray-400">Matrix</span>
				<span class="text-[10px] font-bold text-white ml-0.5">{matrixTasks.length}</span>
			</div>
			<div class="flex items-center gap-1 bg-gray-900 border border-gray-800 rounded-md px-2 py-1">
				<CalendarDays size={10} class="text-cyan-400" />
				<span class="text-[10px] text-gray-400">Calendar</span>
				<span class="text-[10px] font-bold text-white ml-0.5">{calendarTasks.length}</span>
			</div>
		</div>
	</div>

	<!-- Filters -->
	<div class="bg-gray-900 rounded-lg border border-gray-800 p-2 sm:p-2.5 mb-2">
		<div class="flex flex-wrap items-center gap-1.5">
			<button
				onclick={() => filterType = 'all'}
				class="px-2.5 py-1.5 rounded-md text-xs font-medium transition {filterType === 'all' ? 'bg-purple-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700 hover:text-gray-200'}"
			>All ({archivedTasks.length})</button>
			<button
				onclick={() => filterType = 'kanban'}
				class="px-2.5 py-1.5 rounded-md text-xs font-medium transition inline-flex items-center gap-1 {filterType === 'kanban' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700 hover:text-gray-200'}"
			><Kanban size={11} />Kanban ({kanbanTasks.length})</button>
			<button
				onclick={() => filterType = 'matrix'}
				class="px-2.5 py-1.5 rounded-md text-xs font-medium transition inline-flex items-center gap-1 {filterType === 'matrix' ? 'bg-green-700 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700 hover:text-gray-200'}"
			><LayoutGrid size={11} />Matrix ({matrixTasks.length})</button>
			<button
				onclick={() => filterType = 'calendar'}
				class="px-2.5 py-1.5 rounded-md text-xs font-medium transition inline-flex items-center gap-1 {filterType === 'calendar' ? 'bg-cyan-700 text-white' : 'bg-gray-800 text-gray-400 hover:bg-gray-700 hover:text-gray-200'}"
			><CalendarDays size={11} />Calendar ({archivedTasks.filter(t => t.task_type === 'calendar').length})</button>
			<div class="ml-auto">
				<select
					bind:value={sortBy}
					class="h-7 px-2 bg-gray-800 border border-gray-700 text-gray-300 text-xs rounded-md focus:outline-none focus:ring-1 focus:ring-purple-500"
				>
					<option value="date">By Date</option>
					<option value="type">By Type</option>
				</select>
			</div>
		</div>
	</div>

	<!-- Bulk Action Bar -->
	{#if selectedTaskIds.size > 0}
		<div class="bg-blue-950 border border-blue-800 rounded-lg p-2 mb-2 flex items-center justify-between">
			<div class="flex items-center gap-2">
				<span class="text-xs font-semibold text-white">{selectedTaskIds.size} selected</span>
				<button onclick={clearSelection} class="text-[10px] text-blue-300 hover:text-blue-200 underline">Clear</button>
			</div>
			<button
				onclick={handleBulkDelete}
				class="h-7 px-3 bg-red-600 hover:bg-red-700 text-white rounded-md font-medium text-xs transition inline-flex items-center gap-1.5"
			><Trash2 size={12} />Delete</button>
		</div>
	{/if}

		<!-- Tasks Table -->
	{#if loading}
		<div class="flex items-center justify-center py-16">
			<div class="flex flex-col items-center gap-3">
				<div class="w-8 h-8 border-2 border-gray-700 border-t-purple-500 rounded-full animate-spin"></div>
				<p class="text-xs text-gray-500">Loading archive…</p>
			</div>
		</div>
	{:else if filteredTasks().length === 0}
		<div class="bg-gray-900 border border-gray-800 rounded-lg p-10 text-center">
			<Archive class="mx-auto text-gray-700 mb-3" size={28} />
			<p class="text-xs text-gray-500">No archived tasks{filterType !== 'all' ? ` in ${filterType}` : ''}</p>
		</div>
		{:else}
			<div class="space-y-3">
				<div class="hidden md:block bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead class="bg-gray-800/80 border-b border-gray-700">
							<tr>
								<th class="px-3 py-2 text-left w-10">
									<input type="checkbox" checked={isSelectAllChecked} onchange={toggleSelectAll}
										class="w-3.5 h-3.5 rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-1 focus:ring-blue-500 cursor-pointer" aria-label="Select all tasks" />
								</th>
								<th class="px-3 py-2 text-left text-[9px] font-semibold text-gray-500 uppercase tracking-wider">Date</th>
								<th class="px-3 py-2 text-left text-[9px] font-semibold text-gray-500 uppercase tracking-wider">Kanban</th>
								<th class="px-3 py-2 text-left text-[9px] font-semibold text-gray-500 uppercase tracking-wider">Matrix</th>
								<th class="px-3 py-2 text-left text-[9px] font-semibold text-gray-500 uppercase tracking-wider">Calendar</th>
								<th class="px-3 py-2 text-left text-[9px] font-semibold text-gray-500 uppercase tracking-wider">Status</th>
								<th class="px-3 py-2 text-right text-[9px] font-semibold text-gray-500 uppercase tracking-wider">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-gray-800/60">
							{#each filteredTasks() as task}
								<tr class="hover:bg-gray-800/40 transition {selectedTaskIds.has(task.id) ? 'bg-blue-950/20' : ''}">
									<td class="px-3 py-2">
										<input type="checkbox" checked={selectedTaskIds.has(task.id)} onchange={() => toggleTaskSelection(task.id)}
											class="w-3.5 h-3.5 rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-1 focus:ring-blue-500 cursor-pointer" aria-label="Select task" />
									</td>
									<td class="px-3 py-2 whitespace-nowrap">
										<div class="flex items-center gap-1.5">
											<Calendar size={12} class="text-gray-600" />
											<span class="text-xs text-gray-300">{formatDate(task.deleted_at || task.completed_at || task.updated_at)}</span>
										</div>
									</td>
									<td class="px-3 py-2">
										{#if task.task_type === 'kanban'}
											<div class="max-w-xs">
												<div class="font-medium text-gray-200 text-xs">{task.title}</div>
												{#if task.description}<div class="text-[11px] text-gray-500 truncate">{task.description}</div>{/if}
												<div class="flex gap-1 mt-1">
													<span class="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-medium bg-blue-950 text-blue-300">{task.status}</span>
													<span class="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-medium bg-gray-800 text-gray-400">{task.priority}</span>
												</div>
											</div>
										{:else}<span class="text-gray-600 text-xs">—</span>{/if}
									</td>
									<td class="px-3 py-2">
										{#if task.task_type === 'matrix'}
											<div class="max-w-xs">
												<div class="font-medium text-gray-200 text-xs">{task.title}</div>
												{#if task.description}<div class="text-[11px] text-gray-500 truncate">{task.description}</div>{/if}
												<span class="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-medium bg-green-950 text-green-300 mt-1">{getQuadrantLabel(task.quadrant)}</span>
											</div>
										{:else}<span class="text-gray-600 text-xs">—</span>{/if}
									</td>
									<td class="px-3 py-2">
										{#if task.task_type === 'calendar'}
											<div class="max-w-xs">
												<div class="font-medium text-gray-200 text-xs">{task.title}</div>
												{#if task.description}<div class="text-[11px] text-gray-500 truncate">{task.description}</div>{/if}
												<div class="flex gap-1 mt-1 flex-wrap">
													<span class="inline-flex items-center px-1.5 py-0.5 rounded text-[9px] font-medium bg-cyan-950 text-cyan-300">{task.status}</span>
													{#if task.due_date}
														<span class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[9px] font-medium bg-blue-950 text-blue-300">
															<CalendarDays size={8} />{new Date(task.due_date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
														</span>
													{/if}
												</div>
											</div>
										{:else}<span class="text-gray-600 text-xs">—</span>{/if}
									</td>
									<td class="px-3 py-2 whitespace-nowrap">
										{#if task.deleted_at}
											<span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-red-950 text-red-300">Deleted</span>
										{:else if task.is_archived}
											<span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-yellow-950 text-yellow-300">Archived</span>
										{:else if task.completed_at}
											<span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-indigo-950 text-indigo-300">Completed</span>
										{/if}
									</td>
									<td class="px-3 py-2 whitespace-nowrap text-right">
										<div class="flex gap-1.5 justify-end">
											<button onclick={() => handleRestore(task.id)}
												class="inline-flex items-center gap-1 h-7 px-2.5 bg-green-600 text-white rounded-md hover:bg-green-700 transition text-xs font-medium"
												title="Restore task"><RotateCcw size={12} />Restore</button>
											<button onclick={() => handlePermanentDelete(task.id)}
												class="inline-flex items-center gap-1 h-7 px-2.5 bg-red-600 text-white rounded-md hover:bg-red-700 transition text-xs font-medium"
												title="Permanently delete"><Trash2 size={12} />Delete</button>
										</div>
									</td>
								</tr>
							{/each}
					</tbody>
				</table>
			</div>
		</div>

				<div class="md:hidden space-y-2">
				{#each filteredTasks() as task}
					<div class="bg-gray-900 border rounded-lg p-3 {selectedTaskIds.has(task.id) ? 'border-blue-700 bg-blue-950/20' : 'border-gray-800'}">
						<div class="flex items-start gap-2.5 mb-2">
							<input type="checkbox" checked={selectedTaskIds.has(task.id)} onchange={() => toggleTaskSelection(task.id)}
								class="mt-0.5 w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-600 cursor-pointer flex-shrink-0" aria-label="Select task" />
							<div class="flex-1 min-w-0">
								<div class="flex items-start justify-between gap-2">
									<div class="min-w-0">
										<p class="text-[9px] uppercase tracking-wide text-gray-600 mb-0.5">
											{task.task_type === 'matrix' ? 'Matrix' : task.task_type === 'calendar' ? 'Calendar' : 'Kanban'}
										</p>
										<h3 class="text-xs font-semibold text-gray-200 leading-snug">{task.title}</h3>
									</div>
									<span class="text-[9px] text-gray-600 flex-shrink-0 flex items-center gap-0.5 mt-0.5">
										<Calendar size={9} />{formatDate(task.deleted_at || task.completed_at || task.updated_at)}
									</span>
								</div>
							</div>
						</div>
						{#if task.description}
							<p class="text-[11px] text-gray-500 mb-2 ml-6 leading-snug">{task.description}</p>
						{/if}
						<div class="flex flex-wrap gap-1 mb-2 ml-6">
							{#if task.deleted_at}
								<span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-red-950 text-red-300">Deleted</span>
							{:else if task.is_archived}
								<span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-yellow-950 text-yellow-300">Archived</span>
							{:else if task.completed_at}
								<span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-indigo-950 text-indigo-300">Completed</span>
							{/if}
							{#if task.task_type === 'matrix'}
								<span class="inline-flex items-center px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-green-950 text-green-300">{getQuadrantLabel(task.quadrant)}</span>
							{:else if task.task_type === 'calendar' && task.due_date}
								<span class="inline-flex items-center gap-1 px-1.5 py-0.5 rounded-full text-[9px] font-medium bg-cyan-950 text-cyan-300">
									<CalendarDays size={8} />Due: {new Date(task.due_date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
								</span>
							{/if}
						</div>
						<div class="grid grid-cols-2 gap-1.5 ml-6">
							<button onclick={() => handleRestore(task.id)}
								class="inline-flex items-center justify-center gap-1.5 h-8 bg-green-600 text-white rounded-md font-medium text-xs hover:bg-green-700 transition">
								<RotateCcw size={13} />Restore
							</button>
							<button onclick={() => handlePermanentDelete(task.id)}
								class="inline-flex items-center justify-center gap-1.5 h-8 bg-red-600 text-white rounded-md font-medium text-xs hover:bg-red-700 transition">
								<Trash2 size={13} />Delete
							</button>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>
