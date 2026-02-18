<script lang="ts">
	import { onMount } from 'svelte';
	import { taskAPI } from '$lib/api';
	import TonishLogo from '$lib/components/TonishLogo.svelte';
	import { Calendar, RotateCcw, Trash2, Archive } from 'lucide-svelte';

	interface Task {
		id: number;
		title: string;
		description: string;
		task_type: 'kanban' | 'matrix';
		status: string;
		quadrant: string;
		priority: string;
		is_archived: boolean;
		completed_at: string | null;
		deleted_at: string | null;
		created_at: string;
		updated_at: string;
	}

	let archivedTasks = $state<Task[]>([]);
	let loading = $state(true);
	let filterType = $state<'all' | 'kanban' | 'matrix'>('all');
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
</script>

<svelte:head>
	<title>LookBack - Tonish</title>
</svelte:head>

<div class="min-h-screen bg-gray-950 p-6">
	<div class="max-w-7xl mx-auto">
		<!-- Header -->
		<div class="mb-8">
			<div class="flex items-center gap-3 mb-2">
				<TonishLogo size="md" variant="icon" />
				<Archive class="text-purple-400" size={32} />
				<h1 class="text-4xl font-bold text-white">LookBack</h1>
			</div>
			<p class="text-gray-400">Review, restore, or permanently delete archived and completed tasks</p>
		</div>

		<!-- Filters -->
		<div class="bg-gray-900 rounded-lg shadow-sm border border-gray-800 p-4 mb-6">
			<div class="flex flex-col gap-4 md:flex-row md:items-center">
				<div class="flex flex-wrap gap-2">
					<button
						onclick={() => filterType = 'all'}
						class="px-4 py-2 rounded-lg font-medium transition {filterType === 'all' ? 'bg-purple-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
					>
						All ({archivedTasks.length})
					</button>
					<button
						onclick={() => filterType = 'kanban'}
						class="px-4 py-2 rounded-lg font-medium transition {filterType === 'kanban' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
					>
						Kanban ({kanbanTasks.length})
					</button>
					<button
						onclick={() => filterType = 'matrix'}
						class="px-4 py-2 rounded-lg font-medium transition {filterType === 'matrix' ? 'bg-green-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
					>
						Matrix ({matrixTasks.length})
					</button>
				</div>
				
				<div class="flex gap-2 w-full md:w-auto md:ml-auto">
					<select
						bind:value={sortBy}
						class="w-full md:w-auto px-4 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
					>
						<option value="date">Sort by Date</option>
						<option value="type">Sort by Type</option>
					</select>
				</div>
			</div>
		</div>

		<!-- Bulk Action Bar -->
		{#if selectedTaskIds.size > 0}
			<div class="bg-blue-950 border border-blue-800 rounded-lg shadow-lg p-4 mb-6 flex items-center justify-between">
				<div class="flex items-center gap-3">
					<div class="text-white font-semibold">
						{selectedTaskIds.size} item{selectedTaskIds.size > 1 ? 's' : ''} selected
					</div>
					<button
						onclick={clearSelection}
						class="text-sm text-blue-300 hover:text-blue-200 underline touch-manipulation"
					>
						Clear selection
					</button>
				</div>
				<button
					onclick={handleBulkDelete}
					class="min-h-[44px] px-5 py-2.5 bg-red-600 hover:bg-red-700 active:bg-red-800 text-white rounded-lg font-medium transition flex items-center gap-2 touch-manipulation"
				>
					<Trash2 size={20} />
					<span>Delete Selected</span>
				</button>
			</div>
		{/if}

		<!-- Tasks Table -->
		{#if loading}
			<div class="text-center py-12">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
				<p class="mt-4 text-gray-400">Loading archived tasks...</p>
			</div>
		{:else if filteredTasks().length === 0}
			<div class="bg-gray-900 border border-gray-800 rounded-lg shadow-sm p-12 text-center">
				<Archive class="mx-auto text-gray-600 mb-4" size={48} />
				<h3 class="text-xl font-semibold text-gray-200 mb-2">No archived tasks</h3>
				<p class="text-gray-400">Deleted or completed tasks will appear here</p>
			</div>
		{:else}
			<div class="space-y-6">
				<div class="hidden md:block bg-gray-900 border border-gray-800 rounded-lg shadow-sm overflow-hidden">
					<div class="overflow-x-auto">
						<table class="w-full">
							<thead class="bg-gray-800 border-b border-gray-700">
								<tr>
									<th class="px-6 py-4 text-left w-12">
										<input
											type="checkbox"
											checked={isSelectAllChecked}
											onchange={toggleSelectAll}
											class="w-5 h-5 rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-2 focus:ring-blue-500 cursor-pointer"
											aria-label="Select all tasks"
										/>
									</th>
									<th class="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider">Date Deleted/Archived</th>
									<th class="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider">Kanban Tasks</th>
									<th class="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider">Eisenhower Matrix</th>
									<th class="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider">Status</th>
									<th class="px-6 py-4 text-right text-xs font-semibold text-gray-300 uppercase tracking-wider">Actions</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-gray-800">
								{#each filteredTasks() as task}
									<tr class="hover:bg-gray-800 transition {selectedTaskIds.has(task.id) ? 'bg-blue-950/30' : ''}">
										<td class="px-6 py-4">
											<input
												type="checkbox"
												checked={selectedTaskIds.has(task.id)}
												onchange={() => toggleTaskSelection(task.id)}
												class="w-5 h-5 rounded border-gray-600 bg-gray-700 text-blue-600 focus:ring-2 focus:ring-blue-500 cursor-pointer"
												aria-label="Select task"
											/>
										</td>
										<td class="px-6 py-4 whitespace-nowrap">
											<div class="flex items-center gap-2">
												<Calendar size={16} class="text-gray-500" />
												<span class="text-sm text-gray-200">{formatDate(task.deleted_at || task.completed_at || task.updated_at)}</span>
											</div>
										</td>
										<td class="px-6 py-4">
											{#if task.task_type === 'kanban'}
												<div class="max-w-xs">
													<div class="font-medium text-gray-200">{task.title}</div>
													{#if task.description}
														<div class="text-sm text-gray-400 truncate">{task.description}</div>
													{/if}
													<div class="flex gap-2 mt-1">
															<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-950 text-blue-300">
																{task.status}
															</span>
															<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-800 text-gray-300">
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
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										{#if task.deleted_at}
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
												Deleted
											</span>
										{:else if task.is_archived}
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
												Archived
											</span>
										{:else if task.completed_at}
											<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
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
											<p class="text-xs uppercase tracking-wide text-gray-400">{task.task_type === 'matrix' ? 'Matrix task' : 'Kanban task'}</p>
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
								</span>
							{/if}
							<div class="flex flex-wrap gap-2">
								{#if task.deleted_at}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">Deleted</span>
								{:else if task.is_archived}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">Archived</span>
								{:else if task.completed_at}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">Completed</span>
								{/if}
								{#if task.task_type === 'kanban'}
									<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-950 text-blue-300">{task.status}</span>
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
</div>
