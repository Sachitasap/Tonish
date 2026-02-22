<script lang="ts">
	import { onMount } from 'svelte';
	import { taskAPI } from '$lib/api';
	import { Calendar, CalendarDays, RotateCcw, Trash2, Archive } from 'lucide-svelte';

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
	<!-- Filters -->
		<div class="bg-gray-900 rounded-lg border border-gray-800 p-3 sm:p-4 mb-3">
			<div class="flex flex-col gap-3 md:flex-row md:items-center">
				<div class="flex flex-wrap gap-2">
					<button
						onclick={() => filterType = 'all'}
						class="px-3 py-2 min-h-[44px] rounded-md text-sm font-medium transition {filterType === 'all' ? 'bg-purple-600 text-white' : 'bg-gray-800 text-gray-300 hover:bg-gray-700'}"
					>
						All ({archivedTasks.length})
					</button>
					<button
						onclick={() => filterType = 'kanban'}
						class="px-3 py-2 min-h-[44px] rounded-md text-sm font-medium transition {filterType === 'kanban' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-300 hover:bg-gray-700'}"
					>
						Kanban ({kanbanTasks.length})
					</button>
					<button
						onclick={() => filterType = 'matrix'}
						class="px-3 py-2 min-h-[44px] rounded-md text-sm font-medium transition {filterType === 'matrix' ? 'bg-green-600 text-white' : 'bg-gray-800 text-gray-300 hover:bg-gray-700'}"
					>
						Matrix ({matrixTasks.length})
					</button>
					<button
						onclick={() => filterType = 'calendar'}
						class="px-3 py-2 min-h-[44px] rounded-md text-sm font-medium transition flex items-center gap-1.5 {filterType === 'calendar' ? 'bg-cyan-600 text-white' : 'bg-gray-800 text-gray-300 hover:bg-gray-700'}"
					>
						<CalendarDays size={14} />
						Calendar ({archivedTasks.filter(t => t.task_type === 'calendar').length})
					</button>
				</div>
				
				<div class="flex gap-2 w-full md:w-auto md:ml-auto">
					<select
						bind:value={sortBy}
						class="w-full md:w-auto px-3 py-2 min-h-[44px] bg-gray-800 border border-gray-700 text-white text-sm rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
					>
						<option value="date">Sort by Date</option>
						<option value="type">Sort by Type</option>
					</select>
				</div>
			</div>
		</div>

		<!-- Bulk Action Bar -->
		{#if selectedTaskIds.size > 0}
			<div class="bg-blue-950 border border-blue-800 rounded-lg p-3 mb-3 flex items-center justify-between">
				<div class="flex items-center gap-3">
					<div class="text-white font-semibold text-sm">
						{selectedTaskIds.size} item{selectedTaskIds.size > 1 ? 's' : ''} selected
					</div>
					<button
						onclick={clearSelection}
						class="text-xs text-blue-300 hover:text-blue-200 underline touch-manipulation"
					>
						Clear
					</button>
				</div>
				<button
					onclick={handleBulkDelete}
					class="min-h-[44px] px-4 py-2 bg-red-600 hover:bg-red-700 active:bg-red-800 text-white rounded-md font-medium text-sm transition flex items-center gap-2 touch-manipulation"
				>
					<Trash2 size={16} />
					<span>Delete</span>
				</button>
			</div>
		{/if}

		<!-- Tasks Table -->
		{#if loading}
			<div class="text-center py-12">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
				<p class="mt-4 text-gray-400 text-sm">Loading archived tasks...</p>
			</div>
		{:else if filteredTasks().length === 0}
		<div class="bg-gray-900 border border-gray-800 rounded-lg p-10 text-center">
			<Archive class="mx-auto text-gray-600 mb-3" size={36} />
			<p class="text-gray-400 text-sm">No archived tasks</p>
		</div>
		{:else}
			<div class="space-y-3">
				<div class="hidden md:block bg-gray-900 border border-gray-800 rounded-lg overflow-hidden">
					<div class="overflow-x-auto">
						<table class="w-full">
							<thead class="bg-gray-800 border-b border-gray-700">
								<tr>
									<th class="px-4 py-3 text-left w-12">
										<input
											type="checkbox"
											checked={isSelectAllChecked}
											onchange={toggleSelectAll}
											class="w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-2 focus:ring-blue-500 cursor-pointer"
											aria-label="Select all tasks"
										/>
									</th>
									<th class="px-4 py-3 text-left text-[10px] font-semibold text-gray-300 uppercase tracking-wider">Date</th>
									<th class="px-4 py-3 text-left text-[10px] font-semibold text-gray-300 uppercase tracking-wider">Kanban Tasks</th>
									<th class="px-4 py-3 text-left text-[10px] font-semibold text-gray-300 uppercase tracking-wider">Matrix</th>							<th class="px-4 py-3 text-left text-[10px] font-semibold text-gray-300 uppercase tracking-wider">Calendar</th>									<th class="px-4 py-3 text-left text-[10px] font-semibold text-gray-300 uppercase tracking-wider">Status</th>
									<th class="px-4 py-3 text-right text-[10px] font-semibold text-gray-300 uppercase tracking-wider">Actions</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-gray-800">
								{#each filteredTasks() as task}
									<tr class="hover:bg-gray-800 transition {selectedTaskIds.has(task.id) ? 'bg-blue-950/30' : ''}">
										<td class="px-4 py-3">
											<input
												type="checkbox"
												checked={selectedTaskIds.has(task.id)}
												onchange={() => toggleTaskSelection(task.id)}
												class="w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-2 focus:ring-blue-500 cursor-pointer"
												aria-label="Select task"
											/>
										</td>
										<td class="px-4 py-3 whitespace-nowrap">
											<div class="flex items-center gap-2">
												<Calendar size={14} class="text-gray-500" />
												<span class="text-xs text-gray-200">{formatDate(task.deleted_at || task.completed_at || task.updated_at)}</span>
											</div>
										</td>
										<td class="px-4 py-3">
											{#if task.task_type === 'kanban'}
												<div class="max-w-xs">
													<div class="font-medium text-gray-200 text-sm">{task.title}</div>
													{#if task.description}
														<div class="text-xs text-gray-400 truncate">{task.description}</div>
													{/if}
													<div class="flex gap-1.5 mt-1">
															<span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-blue-950 text-blue-300">
																{task.status}
															</span>
															<span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-gray-800 text-gray-300">
															{task.priority}
														</span>
													</div>
												</div>
											{:else}
													<span class="text-gray-500 text-sm">—</span>
											{/if}
									</td>
									<td class="px-6 py-4">
										{#if task.task_type === 'matrix'}
											<div class="max-w-xs">
													<div class="font-medium text-gray-200">{task.title}</div>
													{#if task.description}
														<div class="text-sm text-gray-400 truncate">{task.description}</div>
													{/if}
													<div class="flex gap-2 mt-1">
														<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-950 text-green-300">
														{getQuadrantLabel(task.quadrant)}
													</span>
												</div>
											</div>
										{:else}
											<span class="text-gray-400 text-sm">—</span>
										{/if}
									</td>								<td class="px-4 py-3">
									{#if task.task_type === 'calendar'}
										<div class="max-w-xs">
											<div class="font-medium text-gray-200 text-sm">{task.title}</div>
											{#if task.description}
												<div class="text-xs text-gray-400 truncate">{task.description}</div>
											{/if}
											<div class="flex gap-1.5 mt-1 flex-wrap">
												<span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-cyan-950 text-cyan-300">
													{task.status}
												</span>
												<span class="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-gray-800 text-gray-300">
													{task.priority}
												</span>
												{#if task.due_date}
													<span class="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-medium bg-blue-950 text-blue-300">
														<CalendarDays size={9} />
														{new Date(task.due_date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
													</span>
												{/if}
											</div>
										</div>
									{:else}
										<span class="text-gray-500 text-sm">—</span>
									{/if}
								</td>									<td class="px-6 py-4 whitespace-nowrap">
										{#if task.deleted_at}
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-950 text-red-300">
												Deleted
											</span>
										{:else if task.is_archived}
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-950 text-yellow-300">
												Archived
											</span>
										{:else if task.completed_at}
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-950 text-indigo-300">
												Completed
											</span>
										{/if}
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
										<div class="flex gap-2 justify-end">
											<button
												onclick={() => handleRestore(task.id)}
												class="inline-flex items-center gap-1 min-h-[44px] px-4 py-2.5 bg-green-600 text-white rounded-lg hover:bg-green-700 active:bg-green-800 transition touch-manipulation"
												title="Restore task"
											>
												<RotateCcw size={16} />
												<span>Restore</span>
											</button>
											<button
												onclick={() => handlePermanentDelete(task.id)}
												class="inline-flex items-center gap-1 min-h-[44px] px-4 py-2.5 bg-red-600 text-white rounded-lg hover:bg-red-700 active:bg-red-800 transition touch-manipulation"
												title="Permanently delete"
											>
												<Trash2 size={16} />
												<span>Delete</span>
											</button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>

				<div class="md:hidden space-y-4">
					{#each filteredTasks() as task}
						<div class="bg-gray-900 border rounded-2xl shadow-sm p-4 space-y-3 {selectedTaskIds.has(task.id) ? 'border-blue-500 bg-blue-950/20' : 'border-gray-800'}">
							<div class="flex items-start gap-3">
								<input
									type="checkbox"
									checked={selectedTaskIds.has(task.id)}
									onchange={() => toggleTaskSelection(task.id)}
									class="mt-1 w-6 h-6 min-w-[24px] rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-2 focus:ring-blue-500 cursor-pointer touch-manipulation"
									aria-label="Select task"
								/>
								<div class="flex-1">
									<div class="flex items-start justify-between gap-3">
										<div>
											<p class="text-xs uppercase tracking-wide text-gray-400">
							{task.task_type === 'matrix' ? 'Matrix task' : task.task_type === 'calendar' ? 'Calendar task' : 'Kanban task'}
						</p>
											<h3 class="text-lg font-semibold text-gray-200">{task.title}</h3>
										</div>
										<div class="text-right text-xs text-gray-400 flex flex-col items-end">
											<span class="inline-flex items-center gap-1">
												<Calendar size={14} class="text-gray-500" />
												{formatDate(task.deleted_at || task.completed_at || task.updated_at)}
											</span>
										</div>
									</div>
								</div>
							</div>
							{#if task.description}
								<p class="text-sm text-gray-400">{task.description}</p>
							{/if}
							{#if task.task_type === 'matrix'}
								<span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-semibold bg-green-950 text-green-300">
									{getQuadrantLabel(task.quadrant)}
								</span>						{:else if task.task_type === 'calendar' && task.due_date}
							<span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-cyan-950 text-cyan-300">
								<CalendarDays size={12} />
								Due: {new Date(task.due_date).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
							</span>							{/if}
							<div class="flex flex-wrap gap-2">
								{#if task.deleted_at}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-950 text-red-300">Deleted</span>
								{:else if task.is_archived}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-950 text-yellow-300">Archived</span>
								{:else if task.completed_at}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-950 text-indigo-300">Completed</span>
								{/if}
							{#if task.task_type === 'kanban' || task.task_type === 'calendar'}
								<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {task.task_type === 'calendar' ? 'bg-cyan-950 text-cyan-300' : 'bg-blue-950 text-blue-300'}">{task.status}</span>
								{/if}
							</div>
							<div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
								<button
									onclick={() => handleRestore(task.id)}
									class="inline-flex items-center justify-center gap-2 min-h-[44px] px-4 py-2.5 bg-green-600 text-white rounded-lg font-semibold hover:bg-green-700 active:bg-green-800 transition touch-manipulation"
								>
									<RotateCcw size={18} /> Restore
								</button>
								<button
									onclick={() => handlePermanentDelete(task.id)}
									class="inline-flex items-center justify-center gap-2 min-h-[44px] px-4 py-2.5 bg-red-600 text-white rounded-lg font-semibold hover:bg-red-700 active:bg-red-800 transition touch-manipulation"
								>
									<Trash2 size={18} /> Delete
								</button>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
</div>
