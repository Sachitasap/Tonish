<script lang="ts">
	import { onMount } from 'svelte';
	import { taskAPI } from '$lib/api';
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
				const dateA = new Date(a.deleted_at || a.updated_at).getTime();
				const dateB = new Date(b.deleted_at || b.updated_at).getTime();
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

<div class="min-h-screen bg-gradient-to-br from-purple-50 to-blue-50 p-6">
	<div class="max-w-7xl mx-auto">
		<!-- Header -->
		<div class="mb-8">
			<div class="flex items-center gap-3 mb-2">
				<Archive class="text-purple-600" size={32} />
				<h1 class="text-4xl font-bold text-gray-800">LookBack</h1>
			</div>
			<p class="text-gray-600">Review, restore, or permanently delete archived and completed tasks</p>
		</div>

		<!-- Filters -->
		<div class="bg-white rounded-lg shadow-sm p-4 mb-6">
			<div class="flex flex-wrap gap-4 items-center">
				<div class="flex gap-2">
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
				
				<div class="ml-auto flex gap-2">
					<select
						bind:value={sortBy}
						class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
					>
						<option value="date">Sort by Date</option>
						<option value="type">Sort by Type</option>
					</select>
				</div>
			</div>
		</div>

		<!-- Tasks Table -->
		{#if loading}
			<div class="text-center py-12">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
				<p class="mt-4 text-gray-600">Loading archived tasks...</p>
			</div>
		{:else if filteredTasks().length === 0}
			<div class="bg-white rounded-lg shadow-sm p-12 text-center">
				<Archive class="mx-auto text-gray-400 mb-4" size={48} />
				<h3 class="text-xl font-semibold text-gray-700 mb-2">No archived tasks</h3>
				<p class="text-gray-500">Deleted or completed tasks will appear here</p>
			</div>
		{:else}
			<div class="bg-white rounded-lg shadow-sm overflow-hidden">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead class="bg-gray-50 border-b border-gray-200">
							<tr>
								<th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Date Deleted/Archived</th>
								<th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Kanban Tasks</th>
								<th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Eisenhower Matrix</th>
								<th class="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Status</th>
								<th class="px-6 py-4 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-gray-200">
							{#each filteredTasks() as task}
								<tr class="hover:bg-gray-50 transition">
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="flex items-center gap-2">
											<Calendar size={16} class="text-gray-400" />
											<span class="text-sm text-gray-900">{formatDate(task.deleted_at || task.updated_at)}</span>
										</div>
									</td>
									<td class="px-6 py-4">
										{#if task.task_type === 'kanban'}
											<div class="max-w-xs">
												<div class="font-medium text-gray-900">{task.title}</div>
												{#if task.description}
													<div class="text-sm text-gray-500 truncate">{task.description}</div>
												{/if}
												<div class="flex gap-2 mt-1">
													<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800">
														{task.status}
													</span>
													<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
														{task.priority}
													</span>
												</div>
											</div>
										{:else}
											<span class="text-gray-400 text-sm">—</span>
										{/if}
									</td>
									<td class="px-6 py-4">
										{#if task.task_type === 'matrix'}
											<div class="max-w-xs">
												<div class="font-medium text-gray-900">{task.title}</div>
												{#if task.description}
													<div class="text-sm text-gray-500 truncate">{task.description}</div>
												{/if}
												<div class="flex gap-2 mt-1">
													<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
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
										{/if}
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
										<div class="flex gap-2 justify-end">
											<button
												onclick={() => handleRestore(task.id)}
												class="inline-flex items-center gap-1 px-3 py-1.5 bg-green-600 text-white rounded-lg hover:bg-green-700 transition"
												title="Restore task"
											>
												<RotateCcw size={14} />
												<span>Restore</span>
											</button>
											<button
												onclick={() => handlePermanentDelete(task.id)}
												class="inline-flex items-center gap-1 px-3 py-1.5 bg-red-600 text-white rounded-lg hover:bg-red-700 transition"
												title="Permanently delete"
											>
												<Trash2 size={14} />
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
		{/if}
	</div>
</div>
