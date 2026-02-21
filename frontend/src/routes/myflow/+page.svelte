<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { taskAPI } from '$lib/api';
	import { wsService } from '$lib/websocket';
	import { Plus, Zap, LayoutGrid, GridIcon, Circle, CheckCircle, Edit2, Trash, ArrowLeft, ArrowRight, AlertCircle, Lightbulb, Calendar, Users, CalendarDays } from 'lucide-svelte';

	type TaskStatus = 'todo' | 'in-progress' | 'done';
	type TaskPriority = 'low' | 'medium' | 'high';
	type TaskType = 'kanban' | 'matrix';
	type MatrixQuadrant = 'urgent-important' | 'not-urgent-important' | 'urgent-not-important' | 'not-urgent-not-important';

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
		task_type?: TaskType;
		is_archived?: boolean;
		completed_at?: string | null;
		deleted_at?: string | null;
	};

	type NewTask = {
		title: string;
		description: string;
		priority: TaskPriority;
		status: TaskStatus;
		tags: string;
		due_date?: string;
		is_quick_task: boolean;
		quadrant: string;
		task_type: TaskType;
	};

	let tasks = $state<Task[]>([]);
	let loading = $state(true);
	let showAddTask = $state(false);
	let showKanbanQuickTask = $state(false);
	let showMatrixQuickTask = $state(false);
	let activeView = $state<'kanban' | 'matrix'>('kanban');
	let quickTaskTitle = $state('');
	let matrixQuickTaskTitle = $state('');
	let matrixQuickQuadrant = $state<MatrixQuadrant>('urgent-important');
	let quickAddQuadrant = $state<MatrixQuadrant | null>(null);
	let quickAddTitle = $state('');
	let draggedTask = $state<Task | null>(null);
	let draggedFromQuadrant = $state<string | null>(null);
	let editingTaskId = $state<number | null>(null);
	let editingTask = $state<Task | null>(null);
	let touchStartX = $state<number>(0);
	let touchStartY = $state<number>(0);
	let isTouchDragging = $state<boolean>(false);
	let touchMoveTask = $state<Task | null>(null);
	let touchTargetQuadrant = $state<string | null>(null);
	let newTask = $state<NewTask>({ 
		title: '', 
		description: '', 
		priority: 'medium', 
		status: 'todo',
		tags: '',
		is_quick_task: false,
		quadrant: '',
		task_type: 'kanban'
	});
	const newTaskTitleId = 'new-task-title';
	const newTaskDescriptionId = 'new-task-description';
	const newTaskPriorityId = 'new-task-priority';
	const newTaskStatusId = 'new-task-status';
	const newTaskTagsId = 'new-task-tags';
	const newTaskTypeId = 'new-task-type';
	const quickTaskInputId = 'quick-task-title';
	const matrixQuadrantOptions = [
		{ value: 'urgent-important', label: 'Do First (Urgent & Important)' },
		{ value: 'not-urgent-important', label: 'Schedule (Not Urgent & Important)' },
		{ value: 'urgent-not-important', label: 'Delegate (Urgent & Not Important)' },
		{ value: 'not-urgent-not-important', label: 'Eliminate (Not Urgent & Not Important)' }
	];
	
	onMount(async () => {
		// Check if user is authenticated
		const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
		if (!token) {
			goto('/login');
			return;
		}
		await loadTasks();

		// Set up WebSocket listeners for real-time updates
		const handleTaskUpdate = (task: Task) => {
			const index = tasks.findIndex(t => t.id === task.id);
			if (index >= 0) {
				tasks[index] = task;
			}
		};

		const handleTaskCreate = (task: Task) => {
			if (!task.is_archived && !task.deleted_at) {
				tasks = [...tasks, task];
			}
		};

		const handleTaskDelete = (data: { id: string }) => {
			tasks = tasks.filter(t => t.id !== parseInt(data.id));
		};

		wsService.on('task_update', handleTaskUpdate);
		wsService.on('task_create', handleTaskCreate);
		wsService.on('task_delete', handleTaskDelete);

		return () => {
			wsService.off('task_update', handleTaskUpdate);
			wsService.off('task_create', handleTaskCreate);
			wsService.off('task_delete', handleTaskDelete);
		};
	});

	function setActiveView(view: 'kanban' | 'matrix') {
		activeView = view;
		if (view === 'kanban') {
			showMatrixQuickTask = false;
		} else {
			showKanbanQuickTask = false;
		}
	}
	
	async function loadTasks() {
		try {
			const results = await taskAPI.getAll();
			tasks = results.filter((task) => !task.is_archived && !task.deleted_at);
		} catch (error: any) {
			console.error('Failed to load tasks:', error);
			// If unauthorized, redirect to login
			if (error.message?.includes('Authorization') || error.message?.includes('401')) {
				localStorage.removeItem('authToken');
				goto('/login');
			}
		} finally {
			loading = false;
		}
	}
	
	async function handleAddTask() {
		try {
			// Auto-calculate quadrant based on priority and due date for matrix tasks
			if (newTask.task_type === 'matrix' && newTask.priority) {
				const isImportant = newTask.priority === 'high';
				
				if (isImportant) newTask.quadrant = 'not-urgent-important';
				else newTask.quadrant = 'not-urgent-not-important';
			}
			
			// Clean up the task data - convert date to ISO timestamp
			let dueDate = undefined;
			if (newTask.due_date && newTask.due_date.trim() !== '') {
				// Convert YYYY-MM-DD to ISO timestamp at midnight UTC
				dueDate = new Date(newTask.due_date + 'T00:00:00Z').toISOString();
			}
			
			const taskData = {
				...newTask,
				due_date: dueDate
			};
			
			await taskAPI.create(taskData);
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
				is_quick_task: true,
				task_type: 'kanban'
			});
			await loadTasks();
			quickTaskTitle = '';
			showKanbanQuickTask = false;
		} catch (error) {
			console.error('Failed to create quick task:', error);
		}
	}

	async function handleMatrixQuickTask() {
		if (!matrixQuickTaskTitle.trim()) return;

		try {
			await taskAPI.create({
				title: matrixQuickTaskTitle,
				priority: 'high',
				status: 'todo',
				is_quick_task: true,
				task_type: 'matrix',
				quadrant: matrixQuickQuadrant
			});
			await loadTasks();
			matrixQuickTaskTitle = '';
			showMatrixQuickTask = false;
		} catch (error) {
			console.error('Failed to create matrix quick task:', error);
		}
	}
	
	function startEditTask(task: Task) {
		editingTaskId = task.id;
		editingTask = { ...task };
	}

	function cancelEditTask() {
		editingTaskId = null;
		editingTask = null;
	}

	async function handleSaveEditedTask() {
		if (!editingTask) return;
		
		try {
			await taskAPI.update(editingTask.id, editingTask);
			await loadTasks();
			cancelEditTask();
		} catch (error) {
			console.error('Failed to update task:', error);
		}
	}
	
	function resetNewTask() {
		newTask = { 
			title: '', 
			description: '', 
			priority: 'medium', 
			status: 'todo',
			tags: '',
			due_date: undefined,
			is_quick_task: false,
			quadrant: '',
			task_type: 'kanban'
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
	
	async function handleMoveToStatus(task: Task, newStatus: TaskStatus) {
		try {
			await taskAPI.update(task.id, { ...task, status: newStatus });
			await loadTasks();
		} catch (error) {
			console.error('Failed to update task:', error);
		}
	}

	async function handleQuickAddToQuadrant() {
		if (!quickAddTitle.trim() || !quickAddQuadrant) return;
		
		try {
			await taskAPI.create({
				title: quickAddTitle,
				priority: 'high',
				status: 'todo',
				is_quick_task: true,
				task_type: 'matrix',
				quadrant: quickAddQuadrant
			});
			await loadTasks();
			quickAddTitle = '';
			quickAddQuadrant = null;
		} catch (error) {
			console.error('Failed to create quick task:', error);
		}
	}
	
	const kanbanTasks = $derived(
		tasks.filter((t) => {
			if (t.task_type === 'matrix') return false;
			if (t.task_type === 'kanban') return true;
			return !t.quadrant || t.quadrant.trim() === '';
		})
	);
	const matrixTasks = $derived(
		tasks.filter((t) => {
			if (t.task_type === 'matrix') return true;
			return !!(t.quadrant && t.quadrant.trim() !== '');
		})
	);
	
	const todoTasks = $derived(kanbanTasks.filter((t) => t.status === 'todo'));
	const inProgressTasks = $derived(kanbanTasks.filter((t) => t.status === 'in-progress'));
	const doneTasks = $derived(kanbanTasks.filter((t) => t.status === 'done'));
	
	// Eisenhower Matrix quadrants
	const urgentImportant = $derived(matrixTasks.filter((t) => t.quadrant === 'urgent-important'));
	const notUrgentImportant = $derived(matrixTasks.filter((t) => t.quadrant === 'not-urgent-important'));
	const urgentNotImportant = $derived(matrixTasks.filter((t) => t.quadrant === 'urgent-not-important'));
	const notUrgentNotImportant = $derived(matrixTasks.filter((t) => t.quadrant === 'not-urgent-not-important'));
	
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

	// Touch event handlers for mobile drag and drop
	function handleTouchStart(e: TouchEvent, task: Task, quadrant: string) {
		const touch = e.touches[0];
		touchStartX = touch.clientX;
		touchStartY = touch.clientY;
		touchMoveTask = task;
		draggedTask = task;
		draggedFromQuadrant = quadrant;
		isTouchDragging = false;
	}

	function handleTouchMove(e: TouchEvent) {
		if (!touchMoveTask) return;

		const touch = e.touches[0];
		const deltaX = Math.abs(touch.clientX - touchStartX);
		const deltaY = Math.abs(touch.clientY - touchStartY);

		// Start dragging if moved more than 10px
		if (deltaX > 10 || deltaY > 10) {
			isTouchDragging = true;
			e.preventDefault();

			// Find element under touch point
			const elementBelow = document.elementFromPoint(touch.clientX, touch.clientY);
			const quadrantElement = elementBelow?.closest('[data-quadrant]');
			
			if (quadrantElement) {
				const quadrant = quadrantElement.getAttribute('data-quadrant');
				touchTargetQuadrant = quadrant;
			} else {
				touchTargetQuadrant = null;
			}
		}
	}

	async function handleTouchEnd(e: TouchEvent) {
		if (!isTouchDragging || !touchMoveTask || !touchTargetQuadrant) {
			// Reset states
			touchMoveTask = null;
			draggedTask = null;
			draggedFromQuadrant = null;
			isTouchDragging = false;
			touchTargetQuadrant = null;
			return;
		}

		e.preventDefault();

		// Don't move if dropped in same quadrant
		if (touchTargetQuadrant === draggedFromQuadrant) {
			touchMoveTask = null;
			draggedTask = null;
			draggedFromQuadrant = null;
			isTouchDragging = false;
			touchTargetQuadrant = null;
			return;
		}

		try {
			await taskAPI.update(touchMoveTask.id, {
				...touchMoveTask,
				quadrant: touchTargetQuadrant
			});
			await loadTasks();
		} catch (error) {
			console.error('Failed to move task:', error);
		} finally {
			touchMoveTask = null;
			draggedTask = null;
			draggedFromQuadrant = null;
			isTouchDragging = false;
			touchTargetQuadrant = null;
		}
	}

	// Render priority badge
	function getPriorityBadgeClass(priority: TaskPriority) {
		if (priority === 'high') {
			return { bg: 'bg-red-950', text: 'text-red-300', color: 'fill-red-600 text-red-600' };
		} else if (priority === 'medium') {
			return { bg: 'bg-yellow-900', text: 'text-yellow-300', color: 'fill-yellow-500 text-yellow-500' };
		}
		return { bg: 'bg-green-950', text: 'text-green-300', color: 'fill-green-600 text-green-600' };
	}

	function getPriorityLabel(priority: TaskPriority) {
		return priority === 'high' ? 'High' : priority === 'medium' ? 'Med' : 'Low';
	}
</script>

<svelte:head>
	<title>MyFlow - Tonish</title>
</svelte:head>

<div class="max-w-[1600px] mx-auto">
	<!-- View Switcher + Actions row + action buttons row -->
	<div class="flex flex-wrap gap-2 bg-gray-900 rounded-lg shadow border border-gray-800 p-2 mb-2" role="tablist" aria-label="View switcher">
			<button
				type="button"
				role="tab"
				aria-selected={activeView === 'kanban'}
				onclick={() => setActiveView('kanban')}
				class={`flex-1 min-w-[130px] px-4 py-2.5 min-h-[44px] rounded-md text-sm font-medium transition inline-flex items-center justify-center gap-2 touch-manipulation ${activeView === 'kanban' ? 'bg-blue-600 text-white' : 'text-gray-300 hover:bg-gray-800 active:bg-gray-700'}`}
			>
				<LayoutGrid size={16} />
				Kanban Board
			</button>
			<button
				type="button"
				role="tab"
				aria-selected={activeView === 'matrix'}
				onclick={() => setActiveView('matrix')}
				class={`flex-1 min-w-[130px] px-4 py-2.5 min-h-[44px] rounded-md text-sm font-medium transition inline-flex items-center justify-center gap-2 touch-manipulation ${activeView === 'matrix' ? 'bg-blue-600 text-white' : 'text-gray-300 hover:bg-gray-800 active:bg-gray-700'}`}
			>
				<GridIcon size={16} />
			Eisenhower
		</button>
		<!-- Action buttons sit in the same bar -->
		{#if activeView === 'kanban'}
			<div class="flex gap-2 ml-auto">
				<button
					onclick={() => { showKanbanQuickTask = !showKanbanQuickTask; if (showKanbanQuickTask) showMatrixQuickTask = false; }}
					class="bg-green-600 text-white px-3 py-2 min-h-[40px] rounded-md text-sm font-medium hover:bg-green-700 inline-flex items-center gap-1.5 touch-manipulation"
				><Zap size={14} /> Quick</button>
				<button
					onclick={() => showAddTask = !showAddTask}
					class="bg-blue-600 text-white px-3 py-2 min-h-[40px] rounded-md text-sm font-medium hover:bg-blue-700 inline-flex items-center gap-1.5 touch-manipulation"
				><Plus size={14} /> Add</button>
			</div>
		{:else if activeView === 'matrix'}
			<button
				onclick={() => { showMatrixQuickTask = !showMatrixQuickTask; if (showMatrixQuickTask) showKanbanQuickTask = false; }}
				class="ml-auto bg-purple-600 text-white px-3 py-2 min-h-[40px] rounded-md text-sm font-medium hover:bg-purple-700 inline-flex items-center gap-1.5 touch-manipulation"
			><Zap size={14} /> Quick</button>
		{/if}
	</div>
		
		<!-- Quick Task Input -->
		{#if showKanbanQuickTask && activeView === 'kanban'}
			<div class="bg-gray-900 rounded-lg border border-gray-800 p-3 sm:p-4">
				<h2 class="text-sm font-semibold mb-3 inline-flex items-center gap-2 text-white"><Zap size={16} /> Quick Task</h2>
				<form onsubmit={(e) => { e.preventDefault(); handleQuickTask(); }} class="flex flex-col gap-2 sm:flex-row">
					<input
						type="text"
						bind:value={quickTaskTitle}
						placeholder="Type task and press Enter..."
						id={quickTaskInputId}
						aria-label="Quick task title"
						class="flex-1 px-4 py-2.5 min-h-[44px] bg-gray-800 border border-gray-700 rounded-lg text-sm text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-green-500 touch-manipulation"
					/>
				<button
					type="submit"
					class="bg-green-600 text-white px-5 py-2.5 min-h-[44px] rounded-lg font-medium text-sm hover:bg-green-700 active:bg-green-800 touch-manipulation"
				>
					Add
				</button>
				<button
					type="button"
					onclick={() => { showKanbanQuickTask = false; quickTaskTitle = ''; }}
					class="bg-gray-700 text-gray-300 px-5 py-2.5 min-h-[44px] rounded-lg text-sm hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
				>
					Cancel
				</button>
			</form>
		</div>
	{/if}

	{#if showMatrixQuickTask && activeView === 'matrix'}
		<div class="bg-gray-900 rounded-lg border border-gray-800 p-3 sm:p-4">
			<h2 class="text-sm font-semibold mb-3 inline-flex items-center gap-2 text-white"><Zap size={16} /> Matrix Quick Task</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleMatrixQuickTask(); }} class="grid gap-3 md:grid-cols-3">
				<input
					type="text"
					bind:value={matrixQuickTaskTitle}
					placeholder="Task title"
					class="md:col-span-2 px-4 py-2.5 min-h-[44px] bg-gray-800 border border-gray-700 rounded-lg text-sm text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 touch-manipulation"
				/>
				<select
					bind:value={matrixQuickQuadrant}
					class="px-4 py-2.5 min-h-[44px] bg-gray-800 border border-gray-700 rounded-lg text-sm text-white focus:outline-none focus:ring-2 focus:ring-purple-500 touch-manipulation"
				>
					{#each matrixQuadrantOptions as option}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
				<div class="md:col-span-3 flex flex-col gap-2 sm:flex-row">
					<button
						type="submit"
						class="flex-1 bg-purple-600 text-white px-5 py-2.5 min-h-[44px] rounded-lg font-medium text-sm hover:bg-purple-700 active:bg-purple-800 touch-manipulation"
					>
						Add
					</button>
					<button
						type="button"
						onclick={() => { showMatrixQuickTask = false; matrixQuickTaskTitle = ''; }}
						class="flex-1 bg-gray-700 text-gray-300 px-5 py-2.5 min-h-[44px] rounded-lg text-sm hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	{/if}
	
	{#if showAddTask && activeView === 'kanban'}
		<div class="bg-gray-800 rounded-lg shadow p-6">
			<h2 class="text-lg font-semibold mb-4 text-white">New Task</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleAddTask(); }} class="space-y-4">
				<div>
					<label for={newTaskTypeId} class="block text-sm font-medium text-gray-300 mb-2">Task Type</label>
					<select
						bind:value={newTask.task_type}
						id={newTaskTypeId}
						class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
					>
						<option value="kanban">Kanban (Personal Projects)</option>
					</select>
				</div>
				
				<div>
					<label for={newTaskTitleId} class="block text-sm font-medium text-gray-300 mb-2">Title *</label>
					<input
						type="text"
						bind:value={newTask.title}
						required
						placeholder="Enter task title"
						id={newTaskTitleId}
						class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
					/>
				</div>
				
				<div>
					<label for={newTaskDescriptionId} class="block text-sm font-medium text-gray-300 mb-2">Description</label>
					<textarea
						bind:value={newTask.description}
						rows="3"
						placeholder="Add details about your task"
						id={newTaskDescriptionId}
						class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
					></textarea>
				</div>
				
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for={newTaskPriorityId} class="block text-sm font-medium text-gray-300 mb-2">Priority</label>
						<select
							bind:value={newTask.priority}
							id={newTaskPriorityId}
							class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							<option value="low">Low</option>
							<option value="medium">Medium</option>
							<option value="high">High</option>
						</select>
					</div>
					
					<div>
						<label for={newTaskStatusId} class="block text-sm font-medium text-gray-300 mb-2">Status</label>
						<select
							bind:value={newTask.status}
							id={newTaskStatusId}
							class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							<option value="todo">To Do</option>
							<option value="in-progress">In Progress</option>
							<option value="done">Done</option>
						</select>
					</div>
				</div>
				
				<div>
					<label for={newTaskTagsId} class="block text-sm font-medium text-gray-300 mb-2">Tags (comma-separated)</label>
					<input
						type="text"
						bind:value={newTask.tags}
						placeholder="work, urgent, personal"
						id={newTaskTagsId}
						class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
					/>
				</div>
				
				<div>
					<label for="task-due-date" class="block text-sm font-medium text-gray-300 mb-2 inline-flex items-center gap-2">
						<CalendarDays size={16} />
						Due Date
					</label>
					<input
						type="date"
						bind:value={newTask.due_date}
						id="task-due-date"
						class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
					/>
				</div>
				
				<div class="flex gap-2 pt-2">
					<button
						type="submit"
						class="bg-blue-600 text-white px-6 py-2 rounded-md hover:bg-blue-700 font-medium"
					>
						Create Task
					</button>
					<button
						type="button"
						onclick={() => { showAddTask = false; resetNewTask(); }}
						class="bg-gray-700 text-gray-300 px-6 py-2 rounded-md hover:bg-gray-600 font-medium"
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	{/if}
	
	{#if loading}
		<div class="text-center py-12">
			<p class="text-gray-400">Loading tasks...</p>
		</div>
	{:else if activeView === 'kanban'}
		<!-- Kanban Board -->
		<div class="overflow-x-auto -mx-4 sm:mx-0 pb-2 touch-pan-x">
			<div class="flex gap-2 px-4 sm:px-0 snap-x snap-mandatory md:grid md:grid-cols-3 md:gap-2">
			<!-- To Do Column -->
			<section class="bg-gray-800 rounded-lg p-2.5 min-w-[250px] snap-start md:min-w-0">
				<h2 class="text-xs font-semibold mb-2 text-gray-300 inline-flex items-center gap-1.5"><Circle size={13} /> To Do <span class="text-gray-500 ml-0.5">({todoTasks.length})</span></h2>
				<div class="space-y-1.5">
					{#each todoTasks as task}
						{@const badgeClass = getPriorityBadgeClass(task.priority)}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded p-2 border-l-2 border-blue-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded mb-1.5 text-xs text-white placeholder:text-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded mb-1.5 text-xs text-white placeholder:text-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
									placeholder="Description"
									rows="2"
								></textarea>
								<div class="flex gap-1.5">
									<button onclick={handleSaveEditedTask} class="flex-1 bg-green-600 text-white px-2 py-1 rounded text-xs hover:bg-green-700 font-medium">Save</button>
									<button onclick={cancelEditTask} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600">Cancel</button>
								</div>
							</div>
						{:else}
							<div class="bg-gray-900 rounded p-2.5 hover:shadow transition">
								<div class="flex items-start justify-between gap-1.5 mb-1">
									<h3 class="text-xs font-medium text-white leading-tight flex-1">
										{#if task.is_quick_task}<Zap size={11} class="text-green-500 mr-0.5 inline" />{/if}{task.title}
									</h3>
									<span class="text-[10px] {badgeClass.bg} {badgeClass.text} px-1.5 py-0.5 rounded font-semibold shrink-0">{getPriorityLabel(task.priority)}</span>
								</div>
								{#if task.description}
									<p class="text-[11px] text-gray-400 mb-1 leading-tight">{task.description}</p>
								{/if}
								{#if task.tags}
									<div class="flex flex-wrap gap-0.5 mb-1">
										{#each task.tags.split(',') as tag}
											<span class="text-[10px] bg-blue-950 text-blue-300 px-1.5 py-0.5 rounded">#{tag.trim()}</span>
										{/each}
									</div>
								{/if}
								<div class="flex items-center justify-between pt-1.5 border-t border-gray-700 mt-1.5">
									<button onclick={() => handleMoveToStatus(task, 'in-progress')} class="text-[11px] text-blue-400 hover:text-blue-300 inline-flex items-center gap-0.5"><ArrowRight size={12} /> Start</button>
									<div class="flex gap-2">
										<button onclick={() => startEditTask(task)} class="text-[11px] text-purple-400 hover:text-purple-300 inline-flex items-center gap-0.5"><Edit2 size={11} /> Edit</button>
										<button onclick={() => handleDeleteTask(task.id)} class="text-[11px] text-red-400 hover:text-red-300">Del</button>
									</div>
								</div>
							</div>
						{/if}
					{/each}
					{#if todoTasks.length === 0}
						<p class="text-xs text-gray-500 text-center py-4">No tasks</p>
					{/if}
				</div>
			</section>
			
			<!-- In Progress Column -->
			<section class="bg-blue-950 rounded-lg p-2.5 border border-blue-900 min-w-[250px] snap-start md:min-w-0">
				<h2 class="text-xs font-semibold mb-2 text-blue-300 inline-flex items-center gap-1.5"><Zap size={13} /> In Progress <span class="text-blue-600 ml-0.5">({inProgressTasks.length})</span></h2>
				<div class="space-y-1.5">
					{#each inProgressTasks as task}
						{@const badgeClass = getPriorityBadgeClass(task.priority)}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded p-2 border-l-2 border-blue-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded mb-1.5 text-xs text-white placeholder:text-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded mb-1.5 text-xs text-white placeholder:text-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
									placeholder="Description"
									rows="2"
								></textarea>
								<div class="flex gap-1.5">
									<button onclick={handleSaveEditedTask} class="flex-1 bg-green-600 text-white px-2 py-1 rounded text-xs hover:bg-green-700 font-medium">Save</button>
									<button onclick={cancelEditTask} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600">Cancel</button>
								</div>
							</div>
						{:else}
							<div class="bg-gray-900 rounded p-2.5 hover:shadow transition border-l-2 border-blue-500">
								<div class="flex items-start justify-between gap-1.5 mb-1">
									<h3 class="text-xs font-medium text-white leading-tight flex-1">
										{#if task.is_quick_task}<Zap size={11} class="text-green-500 mr-0.5 inline" />{/if}{task.title}
									</h3>
									<span class="text-[10px] {badgeClass.bg} {badgeClass.text} px-1.5 py-0.5 rounded font-semibold shrink-0">{getPriorityLabel(task.priority)}</span>
								</div>
								{#if task.description}
									<p class="text-[11px] text-gray-400 mb-1 leading-tight">{task.description}</p>
								{/if}
								{#if task.tags}
									<div class="flex flex-wrap gap-0.5 mb-1">
										{#each task.tags.split(',') as tag}
											<span class="text-[10px] bg-blue-950 text-blue-300 px-1.5 py-0.5 rounded">#{tag.trim()}</span>
										{/each}
									</div>
								{/if}
								<div class="flex items-center justify-between pt-1.5 border-t border-gray-700 mt-1.5">
									<button onclick={() => handleMoveToStatus(task, 'done')} class="text-[11px] text-green-400 hover:text-green-300 inline-flex items-center gap-0.5"><CheckCircle size={12} /> Done</button>
									<div class="flex gap-2">
										<button onclick={() => handleMoveToStatus(task, 'todo')} class="text-[11px] text-gray-400 hover:text-gray-300 inline-flex items-center gap-0.5"><ArrowLeft size={11} /> Back</button>
										<button onclick={() => startEditTask(task)} class="text-[11px] text-purple-400 hover:text-purple-300"><Edit2 size={11} /></button>
										<button onclick={() => handleDeleteTask(task.id)} class="text-[11px] text-red-400 hover:text-red-300">Del</button>
									</div>
								</div>
							</div>
						{/if}
					{/each}
					{#if inProgressTasks.length === 0}
						<p class="text-xs text-gray-500 text-center py-4">No tasks</p>
					{/if}
				</div>
			</section>
			
			<!-- Done Column -->
			<section class="bg-green-950 rounded-lg p-2.5 border border-green-900 min-w-[250px] snap-start md:min-w-0">
				<h2 class="text-xs font-semibold mb-2 text-green-300 inline-flex items-center gap-1.5"><CheckCircle size={13} /> Done <span class="text-green-600 ml-0.5">({doneTasks.length})</span></h2>
				<div class="space-y-1.5">
					{#each doneTasks as task}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded p-2 border-l-2 border-green-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded mb-1.5 text-xs text-white placeholder:text-gray-500 focus:outline-none focus:ring-1 focus:ring-green-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded mb-1.5 text-xs text-white placeholder:text-gray-500 focus:outline-none focus:ring-1 focus:ring-green-500"
									placeholder="Description"
									rows="2"
								></textarea>
								<div class="flex gap-1.5">
									<button onclick={handleSaveEditedTask} class="flex-1 bg-green-600 text-white px-2 py-1 rounded text-xs hover:bg-green-700 font-medium">Save</button>
									<button onclick={cancelEditTask} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600">Cancel</button>
								</div>
							</div>
						{:else}
							<div class="bg-gray-900 rounded p-2.5 opacity-75 hover:opacity-100 transition">
								<div class="flex items-start justify-between gap-1.5 mb-1">
									<h3 class="text-xs font-medium text-white line-through leading-tight flex-1">
										{#if task.is_quick_task}<Zap size={11} class="text-green-500 mr-0.5 inline" />{/if}{task.title}
									</h3>
								</div>
								{#if task.description}
									<p class="text-[11px] text-gray-400 mb-1 line-through leading-tight">{task.description}</p>
								{/if}
								{#if task.tags}
									<div class="flex flex-wrap gap-0.5 mb-1 opacity-60">
										{#each task.tags.split(',') as tag}
											<span class="text-[10px] bg-gray-800 text-gray-400 px-1.5 py-0.5 rounded">#{tag.trim()}</span>
										{/each}
									</div>
								{/if}
								<div class="flex items-center justify-between pt-1.5 border-t border-gray-700 mt-1.5">
									<button onclick={() => handleMoveToStatus(task, 'todo')} class="text-[11px] text-blue-400 hover:text-blue-300 inline-flex items-center gap-0.5"><ArrowLeft size={11} /> Reopen</button>
									<div class="flex gap-2">
										<button onclick={() => startEditTask(task)} class="text-[11px] text-purple-400 hover:text-purple-300"><Edit2 size={11} /></button>
										<button onclick={() => handleDeleteTask(task.id)} class="text-[11px] text-red-400 hover:text-red-300">Del</button>
									</div>
								</div>
							</div>
						{/if}
					{/each}
					{#if doneTasks.length === 0}
						<p class="text-xs text-gray-500 text-center py-4">No tasks</p>
					{/if}
				</div>
			</section>
			</div>
		</div>
	{:else if activeView === 'matrix'}
		<!-- Eisenhower Matrix -->
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
			<!-- Urgent & Important (Do First) -->
			<div
				data-quadrant="urgent-important"
				class="bg-red-950 border border-red-900 rounded-lg p-2.5 transition-all touch-pan-y {draggedFromQuadrant !== 'urgent-important' && draggedTask ? 'ring-2 ring-green-400' : ''} {touchTargetQuadrant === 'urgent-important' ? 'ring-2 ring-green-400 bg-red-900' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'urgent-important')}
			>
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xs font-semibold text-red-200 flex items-center gap-1.5">
						<AlertCircle size={14} class="text-red-400" /> Do First
						<span class="text-[11px] font-normal text-red-400">({urgentImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'urgent-important' ? null : 'urgent-important'}
						class="text-red-300 hover:bg-red-900 p-1 rounded transition touch-manipulation"
						title="Add task"
					><Plus size={14} /></button>
				</div>

				{#if quickAddQuadrant === 'urgent-important'}
					<div class="mb-2 p-2 bg-gray-900 rounded border border-red-800">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Task title..."
							class="w-full px-2 py-1.5 bg-gray-800 border border-red-900 rounded text-xs mb-2 text-white focus:outline-none focus:ring-1 focus:ring-red-500 touch-manipulation"
						/>
						<div class="flex gap-1.5">
							<button onclick={handleQuickAddToQuadrant} class="flex-1 bg-red-600 text-white px-2 py-1 rounded text-xs hover:bg-red-700 font-medium touch-manipulation">Add</button>
							<button onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600 touch-manipulation">Cancel</button>
						</div>
					</div>
				{/if}

				<div class="space-y-1.5">
					{#each urgentImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded p-2 border-l-2 border-red-500">
								<input type="text" bind:value={editingTask.title} class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-red-500" placeholder="Task title" />
								<textarea bind:value={editingTask.description} class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-red-500 resize-none" placeholder="Description" rows="2" />
								<div class="flex gap-1.5">
									<button onclick={handleSaveEditedTask} class="flex-1 bg-red-600 text-white px-2 py-1 rounded text-xs hover:bg-red-700 font-medium">Save</button>
									<button onclick={cancelEditTask} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600">Cancel</button>
								</div>
							</div>
						{:else}
							<div
								class="bg-gray-900 rounded p-2 border-l-2 border-red-500 cursor-move hover:shadow transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'urgent-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'urgent-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex items-center justify-between gap-1">
									<h3 class="text-xs font-medium text-white leading-tight flex-1">{task.title}</h3>
									<div class="flex items-center gap-1.5 shrink-0">
										<button onclick={() => startEditTask(task)} class="text-purple-400 hover:text-purple-300 touch-manipulation"><Edit2 size={12} /></button>
										<button onclick={() => handleDeleteTask(task.id)} class="text-[11px] text-red-400 hover:text-red-300 touch-manipulation">Del</button>
									</div>
								</div>
							</div>
						{/if}
					{/each}
					{#if urgentImportant.length === 0 && quickAddQuadrant !== 'urgent-important'}
						<p class="text-[11px] text-gray-500 text-center py-3">Drop tasks here</p>
					{/if}
				</div>
			</div>

			<!-- Not Urgent & Important (Schedule) -->
			<div
				data-quadrant="not-urgent-important"
				class="bg-blue-950 border border-blue-900 rounded-lg p-2.5 transition-all touch-pan-y {draggedFromQuadrant !== 'not-urgent-important' && draggedTask ? 'ring-2 ring-green-400' : ''} {touchTargetQuadrant === 'not-urgent-important' ? 'ring-2 ring-green-400 bg-blue-900' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'not-urgent-important')}
			>
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xs font-semibold text-blue-200 flex items-center gap-1.5">
						<Calendar size={14} class="text-blue-400" /> Schedule
						<span class="text-[11px] font-normal text-blue-400">({notUrgentImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'not-urgent-important' ? null : 'not-urgent-important'}
						class="text-blue-300 hover:bg-blue-900 p-1 rounded transition touch-manipulation"
						title="Add task"
					><Plus size={14} /></button>
				</div>

				{#if quickAddQuadrant === 'not-urgent-important'}
					<div class="mb-2 p-2 bg-gray-900 rounded border border-blue-800">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Task title..."
							class="w-full px-2 py-1.5 bg-gray-800 border border-blue-900 rounded text-xs mb-2 text-white focus:outline-none focus:ring-1 focus:ring-blue-500 touch-manipulation"
						/>
						<div class="flex gap-1.5">
							<button onclick={handleQuickAddToQuadrant} class="flex-1 bg-blue-600 text-white px-2 py-1 rounded text-xs hover:bg-blue-700 font-medium touch-manipulation">Add</button>
							<button onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600 touch-manipulation">Cancel</button>
						</div>
					</div>
				{/if}

				<div class="space-y-1.5">
					{#each notUrgentImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded p-2 border-l-2 border-blue-500">
								<input type="text" bind:value={editingTask.title} class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-blue-500" placeholder="Task title" />
								<textarea bind:value={editingTask.description} class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-blue-500 resize-none" placeholder="Description" rows="2" />
								<div class="flex gap-1.5">
									<button onclick={handleSaveEditedTask} class="flex-1 bg-blue-600 text-white px-2 py-1 rounded text-xs hover:bg-blue-700 font-medium">Save</button>
									<button onclick={cancelEditTask} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600">Cancel</button>
								</div>
							</div>
						{:else}
							<div
								class="bg-gray-900 rounded p-2 border-l-2 border-blue-500 cursor-move hover:shadow transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'not-urgent-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'not-urgent-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex items-center justify-between gap-1">
									<h3 class="text-xs font-medium text-white leading-tight flex-1">{task.title}</h3>
									<div class="flex items-center gap-1.5 shrink-0">
										<button onclick={() => startEditTask(task)} class="text-purple-400 hover:text-purple-300 touch-manipulation"><Edit2 size={12} /></button>
										<button onclick={() => handleDeleteTask(task.id)} class="text-[11px] text-red-400 hover:text-red-300 touch-manipulation">Del</button>
									</div>
								</div>
							</div>
						{/if}
					{/each}
					{#if notUrgentImportant.length === 0 && quickAddQuadrant !== 'not-urgent-important'}
						<p class="text-[11px] text-gray-500 text-center py-3">Drop tasks here</p>
					{/if}
				</div>
			</div>

			<!-- Urgent & Not Important (Delegate) -->
			<div
				data-quadrant="urgent-not-important"
				class="bg-yellow-950 border border-yellow-800 rounded-lg p-2.5 transition-all touch-pan-y {draggedFromQuadrant !== 'urgent-not-important' && draggedTask ? 'ring-2 ring-green-400' : ''} {touchTargetQuadrant === 'urgent-not-important' ? 'ring-2 ring-green-400 bg-yellow-900' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'urgent-not-important')}
			>
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xs font-semibold text-yellow-200 flex items-center gap-1.5">
						<Users size={14} class="text-yellow-500" /> Delegate
						<span class="text-[11px] font-normal text-yellow-400">({urgentNotImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'urgent-not-important' ? null : 'urgent-not-important'}
						class="text-yellow-300 hover:bg-yellow-900 p-1 rounded transition touch-manipulation"
						title="Add task"
					><Plus size={14} /></button>
				</div>

				{#if quickAddQuadrant === 'urgent-not-important'}
					<div class="mb-2 p-2 bg-gray-900 rounded border border-yellow-800">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Task title..."
							class="w-full px-2 py-1.5 bg-gray-800 border border-yellow-900 rounded text-xs mb-2 text-white focus:outline-none focus:ring-1 focus:ring-yellow-500 touch-manipulation"
						/>
						<div class="flex gap-1.5">
							<button onclick={handleQuickAddToQuadrant} class="flex-1 bg-yellow-600 text-white px-2 py-1 rounded text-xs hover:bg-yellow-700 font-medium touch-manipulation">Add</button>
							<button onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600 touch-manipulation">Cancel</button>
						</div>
					</div>
				{/if}

				<div class="space-y-1.5">
					{#each urgentNotImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded p-2 border-l-2 border-yellow-500">
								<input type="text" bind:value={editingTask.title} class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-yellow-500" placeholder="Task title" />
								<textarea bind:value={editingTask.description} class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-yellow-500 resize-none" placeholder="Description" rows="2" />
								<div class="flex gap-1.5">
									<button onclick={handleSaveEditedTask} class="flex-1 bg-yellow-600 text-white px-2 py-1 rounded text-xs hover:bg-yellow-700 font-medium">Save</button>
									<button onclick={cancelEditTask} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600">Cancel</button>
								</div>
							</div>
						{:else}
							<div
								class="bg-gray-900 rounded p-2 border-l-2 border-yellow-500 cursor-move hover:shadow transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'urgent-not-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'urgent-not-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex items-center justify-between gap-1">
									<h3 class="text-xs font-medium text-white leading-tight flex-1">{task.title}</h3>
									<div class="flex items-center gap-1.5 shrink-0">
										<button onclick={() => startEditTask(task)} class="text-purple-400 hover:text-purple-300 touch-manipulation"><Edit2 size={12} /></button>
										<button onclick={() => handleDeleteTask(task.id)} class="text-[11px] text-red-400 hover:text-red-300 touch-manipulation">Del</button>
									</div>
								</div>
							</div>
						{/if}
					{/each}
					{#if urgentNotImportant.length === 0 && quickAddQuadrant !== 'urgent-not-important'}
						<p class="text-[11px] text-gray-500 text-center py-3">Drop tasks here</p>
					{/if}
				</div>
			</div>

			<!-- Not Urgent & Not Important (Eliminate) -->
			<div
				data-quadrant="not-urgent-not-important"
				class="bg-gray-900 border border-gray-700 rounded-lg p-2.5 transition-all touch-pan-y {draggedFromQuadrant !== 'not-urgent-not-important' && draggedTask ? 'ring-2 ring-green-400' : ''} {touchTargetQuadrant === 'not-urgent-not-important' ? 'ring-2 ring-green-400 bg-gray-800' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'not-urgent-not-important')}
			>
				<div class="flex items-center justify-between mb-2">
					<h2 class="text-xs font-semibold text-gray-300 flex items-center gap-1.5">
						<Trash size={14} class="text-gray-400" /> Eliminate
						<span class="text-[11px] font-normal text-gray-500">({notUrgentNotImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'not-urgent-not-important' ? null : 'not-urgent-not-important'}
						class="text-gray-400 hover:bg-gray-700 p-1 rounded transition touch-manipulation"
						title="Add task"
					><Plus size={14} /></button>
				</div>

				{#if quickAddQuadrant === 'not-urgent-not-important'}
					<div class="mb-2 p-2 bg-gray-800 rounded border border-gray-600">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Task title..."
							class="w-full px-2 py-1.5 bg-gray-700 border border-gray-600 rounded text-xs mb-2 text-white focus:outline-none focus:ring-1 focus:ring-gray-400 touch-manipulation"
						/>
						<div class="flex gap-1.5">
							<button onclick={handleQuickAddToQuadrant} class="flex-1 bg-gray-600 text-white px-2 py-1 rounded text-xs hover:bg-gray-500 font-medium touch-manipulation">Add</button>
							<button onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600 touch-manipulation">Cancel</button>
						</div>
					</div>
				{/if}

				<div class="space-y-1.5">
					{#each notUrgentNotImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-800 rounded p-2 border-l-2 border-gray-500">
								<input type="text" bind:value={editingTask.title} class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-gray-400" placeholder="Task title" />
								<textarea bind:value={editingTask.description} class="w-full px-2 py-1 bg-gray-700 border border-gray-600 rounded text-xs mb-1.5 text-white focus:outline-none focus:ring-1 focus:ring-gray-400 resize-none" placeholder="Description" rows="2" />
								<div class="flex gap-1.5">
									<button onclick={handleSaveEditedTask} class="flex-1 bg-gray-600 text-white px-2 py-1 rounded text-xs hover:bg-gray-500 font-medium">Save</button>
									<button onclick={cancelEditTask} class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-xs hover:bg-gray-600">Cancel</button>
								</div>
							</div>
						{:else}
							<div
								class="bg-gray-800 rounded p-2 border-l-2 border-gray-500 cursor-move hover:shadow transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'not-urgent-not-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'not-urgent-not-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex items-center justify-between gap-1">
									<h3 class="text-xs font-medium text-white leading-tight flex-1">{task.title}</h3>
									<div class="flex items-center gap-1.5 shrink-0">
										<button onclick={() => startEditTask(task)} class="text-purple-400 hover:text-purple-300 touch-manipulation"><Edit2 size={12} /></button>
										<button onclick={() => handleDeleteTask(task.id)} class="text-[11px] text-red-400 hover:text-red-300 touch-manipulation">Del</button>
									</div>
								</div>
							</div>
						{/if}
					{/each}
					{#if notUrgentNotImportant.length === 0 && quickAddQuadrant !== 'not-urgent-not-important'}
						<p class="text-[11px] text-gray-500 text-center py-3">Drop tasks here</p>
					{/if}
				</div>
			</div>
		</div>

		<div class="bg-blue-950 border border-blue-900 rounded-lg p-2.5 mt-2">
			<h3 class="text-xs font-semibold text-blue-200 mb-1.5 flex items-center gap-1.5"><Lightbulb size={13} class="text-blue-400" /> Matrix Guide</h3>
			<div class="grid grid-cols-2 gap-x-4 gap-y-0.5 text-[11px] text-blue-300">
				<div><strong>Do First:</strong> Critical deadlines</div>
				<div><strong>Schedule:</strong> Long-term goals</div>
				<div><strong>Delegate:</strong> Interruptions</div>
				<div><strong>Eliminate:</strong> Time wasters</div>
			</div>
		</div>
	{/if}
</div>