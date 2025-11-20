<script lang="ts">
	import { onMount } from 'svelte';
	import { taskAPI } from '$lib/api';
	import { Plus, Zap, LayoutGrid, GridIcon, Trash2, Circle, CheckCircle, Calendar, ArrowRight, ArrowLeft, Users, Trash, Lightbulb, AlertCircle } from 'lucide-svelte';

	type TaskStatus = 'todo' | 'in-progress' | 'done';
	type TaskPriority = 'low' | 'medium' | 'high';

	type Task = {
		id: number;
		title: string;
		description?: string | null;
		priority: TaskPriority;
		status: TaskStatus;
		tags?: string | null;
		due_date?: string | null;
		is_quick_task?: boolean;
		quadrant?: string | null;
	};

	type NewTask = {
		title: string;
		description: string;
		priority: TaskPriority;
		status: TaskStatus;
		tags: string;
		due_date: string;
		is_quick_task: boolean;
		quadrant: string;
	};

	let tasks = $state<Task[]>([]);
	let loading = $state(true);
	let showAddTask = $state(false);
	let showQuickTask = $state(false);
	let activeView = $state<'kanban' | 'matrix'>('kanban');
	let quickTaskTitle = $state('');
	let draggedTask = $state<Task | null>(null);
	let draggedFromQuadrant = $state<string | null>(null);
	let newTask = $state<NewTask>({ 
		title: '', 
		description: '', 
		priority: 'medium', 
		status: 'todo',
		tags: '',
		due_date: '',
		is_quick_task: false,
		quadrant: ''
	});
	const quickTaskInputId = 'quick-task-title';
	const newTaskTitleId = 'new-task-title';
	const newTaskDescriptionId = 'new-task-description';
	const newTaskPriorityId = 'new-task-priority';
	const newTaskStatusId = 'new-task-status';
	const newTaskDueDateId = 'new-task-due-date';
	const newTaskTagsId = 'new-task-tags';
	
	onMount(async () => {
		await loadTasks();
	});
	
	async function loadTasks() {
		try {
			tasks = await taskAPI.getAll();
		} catch (error) {
			console.error('Failed to load tasks:', error);
		} finally {
			loading = false;
		}
	}
	
	async function handleAddTask() {
		try {
			// Auto-calculate quadrant based on priority and due date
			if (newTask.due_date && newTask.priority) {
				const dueDate = new Date(newTask.due_date);
				const today = new Date();
				const isUrgent = (dueDate.getTime() - today.getTime()) / (1000 * 60 * 60 * 24) <= 3;
				const isImportant = newTask.priority === 'high';
				
				if (isUrgent && isImportant) newTask.quadrant = 'urgent-important';
				else if (!isUrgent && isImportant) newTask.quadrant = 'not-urgent-important';
				else if (isUrgent && !isImportant) newTask.quadrant = 'urgent-not-important';
				else newTask.quadrant = 'not-urgent-not-important';
			}
			
			await taskAPI.create(newTask);
			await loadTasks();
			showAddTask = false;
			resetNewTask();
		} catch (error) {
			console.error('Failed to create task:', error);
		}
	}
	
	async function handleQuickTask() {
		if (!quickTaskTitle.trim()) return;
		
		try {
			await taskAPI.create({
				title: quickTaskTitle,
				priority: 'medium',
				status: 'todo',
				is_quick_task: true
			});
			await loadTasks();
			quickTaskTitle = '';
			showQuickTask = false;
		} catch (error) {
			console.error('Failed to create quick task:', error);
		}
	}
	
	function resetNewTask() {
		newTask = { 
			title: '', 
			description: '', 
			priority: 'medium', 
			status: 'todo',
			tags: '',
			due_date: '',
			is_quick_task: false,
			quadrant: ''
		};
	}
	
	async function handleDeleteTask(id: number) {
		if (confirm('Are you sure you want to delete this task?')) {
			try {
				await taskAPI.delete(id);
				await loadTasks();
			} catch (error) {
				console.error('Failed to delete task:', error);
			}
		}
	}
	
	async function handleToggleStatus(task: Task) {
		let newStatus = task.status;
		if (task.status === 'todo') newStatus = 'in-progress';
		else if (task.status === 'in-progress') newStatus = 'done';
		else newStatus = 'todo';
		
		try {
			await taskAPI.update(task.id, { ...task, status: newStatus });
			await loadTasks();
		} catch (error) {
			console.error('Failed to update task:', error);
		}
	}
	
	async function handleMoveToStatus(task: Task, newStatus: TaskStatus) {
		try {
			await taskAPI.update(task.id, { ...task, status: newStatus });
			await loadTasks();
		} catch (error) {
			console.error('Failed to update task:', error);
		}
	}
	
	function formatDate(dateString?: string | null) {
		if (!dateString) return '';
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
	}
	
	function isOverdue(dateString?: string | null) {
		if (!dateString) return false;
		const dueDate = new Date(dateString);
		const today = new Date();
		today.setHours(0, 0, 0, 0);
		return dueDate < today;
	}
	
	const todoTasks = $derived(tasks.filter((t) => t.status === 'todo'));
	const inProgressTasks = $derived(tasks.filter((t) => t.status === 'in-progress'));
	const doneTasks = $derived(tasks.filter((t) => t.status === 'done'));
	
	// Eisenhower Matrix quadrants
	const urgentImportant = $derived(tasks.filter((t) => t.quadrant === 'urgent-important'));
	const notUrgentImportant = $derived(tasks.filter((t) => t.quadrant === 'not-urgent-important'));
	const urgentNotImportant = $derived(tasks.filter((t) => t.quadrant === 'urgent-not-important'));
	const notUrgentNotImportant = $derived(tasks.filter((t) => t.quadrant === 'not-urgent-not-important'));
	
	// Drag and drop handlers for Eisenhower Matrix
	function handleDragStart(task: Task, quadrant: string) {
		draggedTask = task;
		draggedFromQuadrant = quadrant;
	}
	
	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		e.dataTransfer!.dropEffect = 'move';
	}
	
	async function handleDrop(e: DragEvent, targetQuadrant: string) {
		e.preventDefault();
		
		if (!draggedTask || targetQuadrant === draggedFromQuadrant) {
			draggedTask = null;
			draggedFromQuadrant = null;
			return;
		}
		
		try {
			await taskAPI.update(draggedTask.id, {
				...draggedTask,
				quadrant: targetQuadrant
			});
			await loadTasks();
		} catch (error) {
			console.error('Failed to move task:', error);
		} finally {
			draggedTask = null;
			draggedFromQuadrant = null;
		}
	}
	
	function handleDragEnd() {
		draggedTask = null;
		draggedFromQuadrant = null;
	}

	// Render priority badge
	function getPriorityBadgeClass(priority: TaskPriority) {
		if (priority === 'high') {
			return { bg: 'bg-red-100', text: 'text-red-800', color: 'fill-red-600 text-red-600' };
		} else if (priority === 'medium') {
			return { bg: 'bg-yellow-100', text: 'text-yellow-800', color: 'fill-yellow-500 text-yellow-500' };
		}
		return { bg: 'bg-green-100', text: 'text-green-800', color: 'fill-green-600 text-green-600' };
	}

	function getPriorityLabel(priority: TaskPriority) {
		return priority === 'high' ? 'High' : priority === 'medium' ? 'Med' : 'Low';
	}
</script>

<div class="space-y-6">
	<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
		<h1 class="text-3xl font-bold text-gray-900">MyFlow</h1>
		<div class="flex flex-wrap gap-2">
			<button
				onclick={() => showQuickTask = !showQuickTask}
				class="bg-green-600 text-white px-5 py-2.5 rounded-full text-base font-medium hover:bg-green-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600 inline-flex items-center gap-2"
			>
				<Zap size={18} /> Quick Task
			</button>
			<button
				onclick={() => showAddTask = !showAddTask}
				class="bg-blue-600 text-white px-5 py-2.5 rounded-full text-base font-medium hover:bg-blue-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 inline-flex items-center gap-2"
			>
				<Plus size={18} /> Add Task
			</button>
		</div>
	</div>
	
	<!-- View Switcher -->
	<div class="flex flex-wrap gap-2 bg-white rounded-2xl shadow p-1" role="tablist" aria-label="View switcher">
		<button
			type="button"
			role="tab"
			aria-selected={activeView === 'kanban'}
			onclick={() => activeView = 'kanban'}
			class={`flex-1 min-w-[140px] px-4 py-2.5 rounded-xl text-sm font-medium transition inline-flex items-center justify-center gap-2 ${activeView === 'kanban' ? 'bg-blue-600 text-white' : 'text-gray-700 hover:bg-gray-100'}`}
		>
			<LayoutGrid size={16} />
			Kanban Board
		</button>
		<button
			type="button"
			role="tab"
			aria-selected={activeView === 'matrix'}
			onclick={() => activeView = 'matrix'}
			class={`flex-1 min-w-[140px] px-4 py-2.5 rounded-xl text-sm font-medium transition inline-flex items-center justify-center gap-2 ${activeView === 'matrix' ? 'bg-blue-600 text-white' : 'text-gray-700 hover:bg-gray-100'}`}
		>
			<GridIcon size={16} />
			Eisenhower Matrix
		</button>
	</div>
	
	<!-- Quick Task Input -->
	{#if showQuickTask}
		<div class="bg-white rounded-lg shadow p-4">
			<h2 class="text-lg font-semibold mb-3 inline-flex items-center gap-2"><Zap size={20} /> Quick Task</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleQuickTask(); }} class="flex flex-col gap-3 sm:flex-row">
				<input
					type="text"
					bind:value={quickTaskTitle}
					placeholder="Type task and press Enter..."
						id={quickTaskInputId}
						aria-label="Quick task title"
					class="flex-1 px-4 py-3 border border-gray-300 rounded-xl text-base focus:outline-none focus:ring-2 focus:ring-green-500"
				/>
				<button
					type="submit"
					class="bg-green-600 text-white px-6 py-3 rounded-xl font-semibold hover:bg-green-700"
				>
					Add
				</button>
				<button
					type="button"
					onclick={() => { showQuickTask = false; quickTaskTitle = ''; }}
					class="bg-gray-200 text-gray-700 px-4 py-3 rounded-xl hover:bg-gray-300"
				>
					Cancel
				</button>
			</form>
		</div>
	{/if}
	
	{#if showAddTask}
		<div class="bg-white rounded-lg shadow p-6">
			<h2 class="text-xl font-semibold mb-4">New Task</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleAddTask(); }}>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="md:col-span-2">
						<label for={newTaskTitleId} class="block text-sm font-medium text-gray-700 mb-2">Title *</label>
						<input
							type="text"
							bind:value={newTask.title}
							required
							placeholder="Enter task title"
							id={newTaskTitleId}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					
					<div class="md:col-span-2">
						<label for={newTaskDescriptionId} class="block text-sm font-medium text-gray-700 mb-2">Description</label>
						<textarea
							bind:value={newTask.description}
							rows="3"
							placeholder="Add details about your task"
							id={newTaskDescriptionId}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						></textarea>
					</div>
					
					<div>
						<label for={newTaskPriorityId} class="block text-sm font-medium text-gray-700 mb-2">Priority</label>
						<select
							bind:value={newTask.priority}
							id={newTaskPriorityId}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
						<option value="low">Low</option>
						<option value="medium">Medium</option>
						<option value="high">High</option>
						</select>
					</div>
					
					<div>
						<label for={newTaskStatusId} class="block text-sm font-medium text-gray-700 mb-2">Status</label>
						<select
							bind:value={newTask.status}
							id={newTaskStatusId}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							<option value="todo">To Do</option>
							<option value="in-progress">In Progress</option>
							<option value="done">Done</option>
						</select>
					</div>
					
					<div>
						<label for={newTaskDueDateId} class="block text-sm font-medium text-gray-700 mb-2">Due Date</label>
						<input
							type="date"
							bind:value={newTask.due_date}
							id={newTaskDueDateId}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					
					<div>
						<label for={newTaskTagsId} class="block text-sm font-medium text-gray-700 mb-2">Tags (comma-separated)</label>
						<input
							type="text"
							bind:value={newTask.tags}
							placeholder="work, urgent, personal"
							id={newTaskTagsId}
							class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					
					<div class="md:col-span-2 flex gap-2">
						<button
							type="submit"
							class="bg-blue-600 text-white px-6 py-2 rounded-md hover:bg-blue-700"
						>
							Create Task
						</button>
						<button
							type="button"
							onclick={() => { showAddTask = false; resetNewTask(); }}
							class="bg-gray-300 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-400"
						>
							Cancel
						</button>
					</div>
				</div>
			</form>
		</div>
	{/if}
	
	{#if loading}
		<div class="text-center py-12">
			<p class="text-gray-500">Loading tasks...</p>
		</div>
	{:else if activeView === 'kanban'}
		<!-- Kanban Board -->
		<div class="overflow-x-auto -mx-4 sm:mx-0 pb-4 touch-pan-y">
			<div class="flex gap-4 px-4 sm:px-0 snap-x snap-mandatory md:grid md:grid-cols-3 md:gap-6">
			<!-- To Do Column -->
			<section class="bg-gray-100 rounded-2xl p-4 min-w-[280px] snap-start md:min-w-0">
				<h2 class="text-lg font-semibold mb-4 text-gray-700 inline-flex items-center gap-2"><Circle size={18} /> To Do ({todoTasks.length})</h2>
				<div class="space-y-3">
					{#each todoTasks as task}
						{@const badgeClass = getPriorityBadgeClass(task.priority)}
						<div class="bg-white rounded-lg shadow p-4 hover:shadow-md transition">
							<div class="flex justify-between items-start mb-2">
								<div class="flex-1">
									<h3 class="font-medium text-gray-900 mb-1">
										{#if task.is_quick_task}
											<Zap size={14} class="text-green-600 mr-1" />
										{/if}
										{task.title}
									</h3>
									{#if task.due_date}
										<p class="text-xs {isOverdue(task.due_date) ? 'text-red-600 font-semibold' : 'text-gray-500'}">
											<Calendar size={14} class="inline mr-1" /> {formatDate(task.due_date)}
											{#if isOverdue(task.due_date)}
												<span class="ml-1">(Overdue)</span>
											{/if}
										</p>
									{/if}
								</div>
								<span class="text-xs {badgeClass.bg} {badgeClass.text} px-2 py-1 rounded font-semibold inline-flex items-center gap-1">
									<Circle size={12} class={badgeClass.color} /> {getPriorityLabel(task.priority)}
								</span>
							</div>
							{#if task.description}
								<p class="text-sm text-gray-600 mb-2">{task.description}</p>
							{/if}
							{#if task.tags}
								<div class="flex flex-wrap gap-1 mb-2">
									{#each task.tags.split(',') as tag}
										<span class="text-xs bg-blue-50 text-blue-700 px-2 py-0.5 rounded">#{tag.trim()}</span>
									{/each}
								</div>
							{/if}
							<div class="flex justify-between items-center mt-3 pt-2 border-t">
								<button
									onclick={() => handleMoveToStatus(task, 'in-progress')}
									class="text-xs text-blue-600 hover:underline font-medium inline-flex items-center gap-1">
									<ArrowRight size={14} /> Start
								</button>
								<button
									onclick={() => handleDeleteTask(task.id)}
									class="text-xs text-red-600 hover:underline"
								>
									Delete
								</button>
							</div>
						</div>
					{/each}
					{#if todoTasks.length === 0}
						<p class="text-sm text-gray-500 text-center py-8">No tasks to do</p>
					{/if}
				</div>
			</section>
			
			<!-- In Progress Column -->
			<section class="bg-blue-50 rounded-2xl p-4 border-2 border-blue-200 min-w-[280px] snap-start md:min-w-0">
				<h2 class="text-lg font-semibold mb-4 text-blue-900 inline-flex items-center gap-2"><Zap size={18} /> In Progress ({inProgressTasks.length})</h2>
				<div class="space-y-3">
					{#each inProgressTasks as task}
						{@const badgeClass = getPriorityBadgeClass(task.priority)}
						<div class="bg-white rounded-lg shadow p-4 hover:shadow-md transition border-l-4 border-blue-500">
							<div class="flex justify-between items-start mb-2">
								<div class="flex-1">
									<h3 class="font-medium text-gray-900 mb-1">
										{#if task.is_quick_task}
											<Zap size={14} class="text-green-600 mr-1 inline" />
										{/if}
										{task.title}
									</h3>
									{#if task.due_date}
										<p class="text-xs {isOverdue(task.due_date) ? 'text-red-600 font-semibold' : 'text-gray-500'}">
											<Calendar size={14} class="inline mr-1" /> {formatDate(task.due_date)}
											{#if isOverdue(task.due_date)}
												<span class="ml-1">(Overdue)</span>
											{/if}
										</p>
									{/if}
								</div>
								<span class="text-xs {badgeClass.bg} {badgeClass.text} px-2 py-1 rounded font-semibold inline-flex items-center gap-1">
									<Circle size={12} class={badgeClass.color} /> {getPriorityLabel(task.priority)}
								</span>
							</div>
							{#if task.description}
								<p class="text-sm text-gray-600 mb-2">{task.description}</p>
							{/if}
							{#if task.tags}
								<div class="flex flex-wrap gap-1 mb-2">
									{#each task.tags.split(',') as tag}
										<span class="text-xs bg-blue-50 text-blue-700 px-2 py-0.5 rounded">#{tag.trim()}</span>
									{/each}
								</div>
							{/if}
							<div class="flex justify-between items-center mt-3 pt-2 border-t">
								<button
									onclick={() => handleMoveToStatus(task, 'done')}
									class="text-xs text-green-600 hover:underline font-medium inline-flex items-center gap-1"
								>
									<CheckCircle size={14} /> Complete
								</button>
								<div class="flex gap-2">
									<button
										onclick={() => handleMoveToStatus(task, 'todo')}
										class="text-xs text-gray-600 hover:underline inline-flex items-center gap-1"
									>
										<ArrowLeft size={14} /> Back
									</button>
									<button
										onclick={() => handleDeleteTask(task.id)}
										class="text-xs text-red-600 hover:underline"
									>
										Delete
									</button>
								</div>
							</div>
						</div>
					{/each}
					{#if inProgressTasks.length === 0}
						<p class="text-sm text-gray-500 text-center py-8">No tasks in progress</p>
					{/if}
				</div>
			</section>
			
			<!-- Done Column -->
			<section class="bg-green-50 rounded-2xl p-4 border-2 border-green-200 min-w-[280px] snap-start md:min-w-0">
				<h2 class="text-lg font-semibold mb-4 text-green-900 inline-flex items-center gap-2"><CheckCircle size={18} /> Done ({doneTasks.length})</h2>
				<div class="space-y-3">
					{#each doneTasks as task}
						<div class="bg-white rounded-lg shadow p-4 opacity-75 hover:opacity-100 transition">
							<div class="flex justify-between items-start mb-2">
								<div class="flex-1">
									<h3 class="font-medium text-gray-900 line-through mb-1">
										{#if task.is_quick_task}
											<Zap size={14} class="text-green-600 mr-1 inline" />
										{/if}
										{task.title}
									</h3>
								</div>
							</div>
							{#if task.description}
								<p class="text-sm text-gray-600 mb-2 line-through">{task.description}</p>
							{/if}
							{#if task.tags}
								<div class="flex flex-wrap gap-1 mb-2 opacity-60">
									{#each task.tags.split(',') as tag}
										<span class="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded">#{tag.trim()}</span>
									{/each}
								</div>
							{/if}
							<div class="flex justify-between items-center mt-3 pt-2 border-t">
								<button
									onclick={() => handleMoveToStatus(task, 'todo')}
									class="text-xs text-blue-600 hover:underline inline-flex items-center gap-1"
								>
									<ArrowLeft size={14} /> Reopen
								</button>
								<button
									onclick={() => handleDeleteTask(task.id)}
									class="text-xs text-red-600 hover:underline"
								>
									Delete
								</button>
							</div>
						</div>
					{/each}
					{#if doneTasks.length === 0}
						<p class="text-sm text-gray-500 text-center py-8">No completed tasks</p>
					{/if}
				</div>
			</section>
			</div>
		</div>
	{:else if activeView === 'matrix'}
		<!-- Eisenhower Matrix -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<!-- Urgent & Important (Do First) -->
			<div 
				class="bg-red-50 border-2 border-red-300 rounded-lg p-4 min-h-[300px] transition-all {draggedFromQuadrant !== 'urgent-important' && draggedTask ? 'ring-2 ring-green-400' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'urgent-important')}
			>
				<h2 class="text-lg font-semibold mb-3 text-red-900 flex items-center gap-2">
					<AlertCircle size={20} class="text-red-600" /> Do First
					<span class="text-sm font-normal text-red-700 block">Urgent & Important ({urgentImportant.length})</span>
				</h2>
				<div class="space-y-2">
					{#each urgentImportant as task}
						<div 
							class="bg-white rounded-lg shadow-sm p-3 border-l-4 border-red-500 cursor-move hover:shadow-md transition {draggedTask?.id === task.id ? 'opacity-50' : ''}"
							draggable="true"
							ondragstart={() => handleDragStart(task, 'urgent-important')}
							ondragend={handleDragEnd}
						>
						<h3 class="font-medium text-gray-900 text-sm">{task.title}</h3>
						{#if task.due_date}
							<p class="text-xs text-red-600 mt-1 flex items-center gap-1"><Calendar size={14} /> {formatDate(task.due_date)}</p>
						{/if}
						<div class="flex gap-2 mt-2">
								<button onclick={() => handleDeleteTask(task.id)} class="text-xs text-red-600 hover:underline">Delete</button>
							</div>
						</div>
					{/each}
					{#if urgentImportant.length === 0}
						<p class="text-xs text-gray-500 text-center py-4">Drop tasks here</p>
					{/if}
				</div>
			</div>
			
			<!-- Not Urgent & Important (Schedule) -->
			<div 
				class="bg-blue-50 border-2 border-blue-300 rounded-lg p-4 min-h-[300px] transition-all {draggedFromQuadrant !== 'not-urgent-important' && draggedTask ? 'ring-2 ring-green-400' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'not-urgent-important')}
			>
			<h2 class="text-lg font-semibold mb-3 text-blue-900 flex items-center gap-2">
				<Calendar size={20} class="text-blue-600" /> Schedule
				<span class="text-sm font-normal text-blue-700 block">Not Urgent & Important ({notUrgentImportant.length})</span>
			</h2>
				<div class="space-y-2">
					{#each notUrgentImportant as task}
						<div 
							class="bg-white rounded-lg shadow-sm p-3 border-l-4 border-blue-500 cursor-move hover:shadow-md transition {draggedTask?.id === task.id ? 'opacity-50' : ''}"
							draggable="true"
							ondragstart={() => handleDragStart(task, 'not-urgent-important')}
							ondragend={handleDragEnd}
						>
					<h3 class="font-medium text-gray-900 text-sm">{task.title}</h3>
					{#if task.due_date}
						<p class="text-xs text-blue-600 mt-1 flex items-center gap-1"><Calendar size={14} /> {formatDate(task.due_date)}</p>
					{/if}
					<div class="flex gap-2 mt-2">
						<button onclick={() => handleDeleteTask(task.id)} class="text-xs text-red-600 hover:underline">Delete</button>
					</div>
				</div>
			{/each}
			{#if notUrgentImportant.length === 0}
				<p class="text-xs text-gray-500 text-center py-4">Drop tasks here</p>
			{/if}
		</div>
	</div>			<!-- Urgent & Not Important (Delegate) -->
			<div 
				class="bg-yellow-50 border-2 border-yellow-300 rounded-lg p-4 min-h-[300px] transition-all {draggedFromQuadrant !== 'urgent-not-important' && draggedTask ? 'ring-2 ring-green-400' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'urgent-not-important')}
			>
			<h2 class="text-lg font-semibold mb-3 text-yellow-900 flex items-center gap-2">
				<Users size={20} class="text-yellow-600" /> Delegate
				<span class="text-sm font-normal text-yellow-700 block">Urgent & Not Important ({urgentNotImportant.length})</span>
			</h2>
				<div class="space-y-2">
					{#each urgentNotImportant as task}
						<div 
							class="bg-white rounded-lg shadow-sm p-3 border-l-4 border-yellow-500 cursor-move hover:shadow-md transition {draggedTask?.id === task.id ? 'opacity-50' : ''}"
							draggable="true"
							ondragstart={() => handleDragStart(task, 'urgent-not-important')}
							ondragend={handleDragEnd}
						>
							<h3 class="font-medium text-gray-900 text-sm">{task.title}</h3>
							{#if task.due_date}
								<p class="text-xs text-yellow-600 mt-1 flex items-center gap-1"><Calendar size={14} /> {formatDate(task.due_date)}</p>
							{/if}
							<div class="flex gap-2 mt-2">
								<button onclick={() => handleDeleteTask(task.id)} class="text-xs text-red-600 hover:underline">Delete</button>
							</div>
						</div>
					{/each}
					{#if urgentNotImportant.length === 0}
						<p class="text-xs text-gray-500 text-center py-4">Drop tasks here</p>
					{/if}
				</div>
			</div>
			
			<!-- Not Urgent & Not Important (Eliminate) -->
			<div 
				class="bg-gray-50 border-2 border-gray-300 rounded-lg p-4 min-h-[300px] transition-all {draggedFromQuadrant !== 'not-urgent-not-important' && draggedTask ? 'ring-2 ring-green-400' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'not-urgent-not-important')}
			>
			<h2 class="text-lg font-semibold mb-3 text-gray-900 flex items-center gap-2">
				<Trash size={20} class="text-gray-600" /> Eliminate
				<span class="text-sm font-normal text-gray-600 block">Not Urgent & Not Important ({notUrgentNotImportant.length})</span>
			</h2>
				<div class="space-y-2">
					{#each notUrgentNotImportant as task}
						<div 
							class="bg-white rounded-lg shadow-sm p-3 border-l-4 border-gray-400 cursor-move hover:shadow-md transition {draggedTask?.id === task.id ? 'opacity-50' : ''}"
							draggable="true"
							ondragstart={() => handleDragStart(task, 'not-urgent-not-important')}
							ondragend={handleDragEnd}
						>
							<h3 class="font-medium text-gray-900 text-sm">{task.title}</h3>
							{#if task.due_date}
								<p class="text-xs text-gray-600 mt-1 flex items-center gap-1"><Calendar size={14} /> {formatDate(task.due_date)}</p>
							{/if}
							<div class="flex gap-2 mt-2">
								<button onclick={() => handleDeleteTask(task.id)} class="text-xs text-red-600 hover:underline">Delete</button>
							</div>
						</div>
					{/each}
					{#if notUrgentNotImportant.length === 0}
						<p class="text-xs text-gray-500 text-center py-4">Drop tasks here</p>
					{/if}
				</div>
			</div>
		</div>
		
		<div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mt-4">
			<h3 class="font-semibold text-blue-900 mb-2 flex items-center gap-2"><Lightbulb size={20} class="text-blue-600" /> Eisenhower Matrix Guide</h3>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm text-blue-800">
				<div><strong>Do First:</strong> Critical deadlines & crises</div>
				<div><strong>Schedule:</strong> Long-term goals & planning</div>
				<div><strong>Delegate:</strong> Interruptions & some emails</div>
				<div><strong>Eliminate:</strong> Time wasters & distractions</div>
			</div>
			<p class="text-xs text-blue-700 mt-2 flex items-center gap-1">
				<Lightbulb size={14} class="text-blue-600" /> Tasks are automatically categorized based on priority and due date (within 3 days = urgent)
			</p>
		</div>
	{/if}
</div>
