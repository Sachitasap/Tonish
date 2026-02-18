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

<div class="space-y-6">
	<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
		<div class="flex items-center gap-3">
			<h1 class="text-3xl font-bold text-white">MyFlow</h1>
		</div>
		{#if activeView === 'kanban'}
			<div class="flex flex-wrap gap-3">
				<button
					onclick={() => {
						showKanbanQuickTask = !showKanbanQuickTask;
						if (showKanbanQuickTask) showMatrixQuickTask = false;
					}}
					class="bg-green-600 text-white px-6 py-3 min-h-[48px] rounded-full text-base font-medium hover:bg-green-700 active:bg-green-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600 inline-flex items-center gap-2 touch-manipulation"
				>
					<Zap size={20} /> Quick Task
				</button>
				<button
					onclick={() => showAddTask = !showAddTask}
					class="bg-blue-600 text-white px-6 py-3 min-h-[48px] rounded-full text-base font-medium hover:bg-blue-700 active:bg-blue-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 inline-flex items-center gap-2 touch-manipulation"
				>
					<Plus size={20} /> Add Task
				</button>
			</div>
		{:else if activeView === 'matrix'}
			<div class="flex flex-wrap gap-3">
				<button
					onclick={() => {
						showMatrixQuickTask = !showMatrixQuickTask;
						if (showMatrixQuickTask) showKanbanQuickTask = false;
					}}
					class="bg-purple-600 text-white px-6 py-3 min-h-[48px] rounded-full text-base font-medium hover:bg-purple-700 active:bg-purple-800 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-purple-600 inline-flex items-center gap-2 touch-manipulation"
				>
					<Zap size={20} /> Matrix Quick Task
				</button>
			</div>
		{/if}
	</div>
	
	<!-- View Switcher -->
	<div class="flex flex-wrap gap-2 bg-gray-900 rounded-2xl shadow p-2" role="tablist" aria-label="View switcher">
		<button
			type="button"
			role="tab"
			aria-selected={activeView === 'kanban'}
			onclick={() => setActiveView('kanban')}
			class={`flex-1 min-w-[140px] px-5 py-3.5 min-h-[48px] rounded-xl text-base font-medium transition inline-flex items-center justify-center gap-2 touch-manipulation ${activeView === 'kanban' ? 'bg-blue-600 text-white' : 'text-gray-300 hover:bg-gray-800 active:bg-gray-700'}`}
		>
			<LayoutGrid size={18} />
			Kanban Board
		</button>
		<button
			type="button"
			role="tab"
			aria-selected={activeView === 'matrix'}
			onclick={() => setActiveView('matrix')}
			class={`flex-1 min-w-[140px] px-5 py-3.5 min-h-[48px] rounded-xl text-base font-medium transition inline-flex items-center justify-center gap-2 touch-manipulation ${activeView === 'matrix' ? 'bg-blue-600 text-white' : 'text-gray-300 hover:bg-gray-800 active:bg-gray-700'}`}
		>
			<GridIcon size={18} />
			Eisenhower Matrix
		</button>
	</div>
	
	<!-- Quick Task Input -->
	{#if showKanbanQuickTask && activeView === 'kanban'}
		<div class="bg-gray-800 rounded-lg shadow p-5">
			<h2 class="text-lg font-semibold mb-4 inline-flex items-center gap-2 text-white"><Zap size={22} /> Quick Task</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleQuickTask(); }} class="flex flex-col gap-3 sm:flex-row">
				<input
					type="text"
					bind:value={quickTaskTitle}
					placeholder="Type task and press Enter..."
						id={quickTaskInputId}
						aria-label="Quick task title"
					class="flex-1 px-5 py-4 min-h-[48px] bg-gray-800 border border-gray-700 rounded-xl text-base text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-green-500 touch-manipulation"
				/>
				<button
					type="submit"
					class="bg-green-600 text-white px-8 py-4 min-h-[48px] rounded-xl font-semibold hover:bg-green-700 active:bg-green-800 touch-manipulation"
				>
					Add
				</button>
				<button
					type="button"
					onclick={() => { showKanbanQuickTask = false; quickTaskTitle = ''; }}
					class="bg-gray-700 text-gray-300 px-6 py-4 min-h-[48px] rounded-xl hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
				>
					Cancel
				</button>
			</form>
		</div>
	{/if}

	{#if showMatrixQuickTask && activeView === 'matrix'}
		<div class="bg-gray-800 rounded-lg shadow p-5">
			<h2 class="text-lg font-semibold mb-4 inline-flex items-center gap-2 text-white"><Zap size={22} /> Matrix Quick Task</h2>
			<form onsubmit={(e) => { e.preventDefault(); handleMatrixQuickTask(); }} class="grid gap-4 md:grid-cols-3">
				<input
					type="text"
					bind:value={matrixQuickTaskTitle}
					placeholder="Task title"
					class="md:col-span-2 px-5 py-4 min-h-[48px] bg-gray-800 border border-gray-700 rounded-xl text-base text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 touch-manipulation"
				/>
				<select
					bind:value={matrixQuickQuadrant}
					class="px-5 py-4 min-h-[48px] bg-gray-800 border border-gray-700 rounded-xl text-base text-white focus:outline-none focus:ring-2 focus:ring-purple-500 touch-manipulation"
				>
					{#each matrixQuadrantOptions as option}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
				<div class="md:col-span-3 flex flex-col gap-3 sm:flex-row">
					<button
						type="submit"
						class="flex-1 bg-purple-600 text-white px-8 py-4 min-h-[48px] rounded-xl font-semibold hover:bg-purple-700 active:bg-purple-800 touch-manipulation"
					>
						Add
					</button>
					<button
						type="button"
						onclick={() => { showMatrixQuickTask = false; matrixQuickTaskTitle = ''; }}
						class="flex-1 bg-gray-700 text-gray-300 px-6 py-4 min-h-[48px] rounded-xl hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
					>
						Cancel
					</button>
				</div>
			</form>
		</div>
	{/if}
	
	{#if showAddTask && activeView === 'kanban'}
		<div class="bg-gray-800 rounded-lg shadow p-6">
			<h2 class="text-xl font-semibold mb-6 text-white">New Task</h2>
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
		<div class="overflow-x-auto -mx-4 sm:mx-0 pb-4 touch-pan-y">
			<div class="flex gap-4 px-4 sm:px-0 snap-x snap-mandatory md:grid md:grid-cols-3 md:gap-6">
			<!-- To Do Column -->
			<section class="bg-gray-800 rounded-2xl p-4 min-w-[280px] snap-start md:min-w-0">
				<h2 class="text-lg font-semibold mb-4 text-gray-300 inline-flex items-center gap-2"><Circle size={18} /> To Do ({todoTasks.length})</h2>
				<div class="space-y-3">
					{#each todoTasks as task}
						{@const badgeClass = getPriorityBadgeClass(task.priority)}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded-lg shadow p-4 border-l-4 border-blue-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
									placeholder="Description"
									rows="2"
								></textarea>
								<div class="flex gap-2">
									<button
										onclick={handleSaveEditedTask}
										class="flex-1 bg-green-600 text-white px-3 py-2 rounded hover:bg-green-700 text-sm font-medium"
									>
										Save
									</button>
									<button
										onclick={cancelEditTask}
										class="flex-1 bg-gray-700 text-gray-300 px-3 py-2 rounded hover:bg-gray-600 text-sm font-medium"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<div class="bg-gray-900 rounded-lg shadow p-4 hover:shadow-md transition">
								<div class="flex justify-between items-start mb-2">
									<div class="flex-1">
										<h3 class="font-medium text-white mb-1">
											{#if task.is_quick_task}
												<Zap size={14} class="text-green-600 mr-1 inline" />
											{/if}
											{task.title}
										</h3>
									</div>
									<span class="text-xs {badgeClass.bg} {badgeClass.text} px-2 py-1 rounded font-semibold inline-flex items-center gap-1">
										<Circle size={12} class={badgeClass.color} /> {getPriorityLabel(task.priority)}
									</span>
								</div>
								{#if task.description}
									<p class="text-sm text-gray-400 mb-2">{task.description}</p>
								{/if}
								{#if task.tags}
									<div class="flex flex-wrap gap-1 mb-2">
										{#each task.tags.split(',') as tag}
											<span class="text-xs bg-blue-950 text-blue-300 px-2 py-0.5 rounded">#{tag.trim()}</span>
										{/each}
									</div>
								{/if}
								<div class="flex justify-between items-center mt-3 pt-2 border-t border-gray-700 flex-wrap gap-2">
									<button
										onclick={() => handleMoveToStatus(task, 'in-progress')}
										class="text-xs text-blue-600 hover:underline font-medium inline-flex items-center gap-1">
										<ArrowRight size={14} /> Start
									</button>
									<div class="flex gap-2">
										<button
											onclick={() => startEditTask(task)}
											class="text-xs text-purple-600 hover:underline font-medium inline-flex items-center gap-1"
										>
											<Edit2 size={14} /> Edit
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
						{/if}
					{/each}
					{#if todoTasks.length === 0}
						<p class="text-sm text-gray-400 text-center py-8">No tasks to do</p>
					{/if}
				</div>
			</section>
			
			<!-- In Progress Column -->
			<section class="bg-blue-950 rounded-2xl p-4 border-2 border-blue-900 min-w-[280px] snap-start md:min-w-0">
				<h2 class="text-lg font-semibold mb-4 text-blue-200 inline-flex items-center gap-2"><Zap size={18} /> In Progress ({inProgressTasks.length})</h2>
				<div class="space-y-3">
					{#each inProgressTasks as task}
						{@const badgeClass = getPriorityBadgeClass(task.priority)}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded-lg shadow p-4 border-l-4 border-blue-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
									placeholder="Description"
									rows="2"
								></textarea>
								<div class="flex gap-2">
									<button
										onclick={handleSaveEditedTask}
										class="flex-1 bg-green-600 text-white px-3 py-2 rounded hover:bg-green-700 text-sm font-medium"
									>
										Save
									</button>
									<button
										onclick={cancelEditTask}
										class="flex-1 bg-gray-700 text-gray-300 px-3 py-2 rounded hover:bg-gray-600 text-sm font-medium"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<div class="bg-gray-900 rounded-lg shadow p-4 hover:shadow-md transition border-l-4 border-blue-500">
								<div class="flex justify-between items-start mb-2">
									<div class="flex-1">
										<h3 class="font-medium text-white mb-1">
											{#if task.is_quick_task}
												<Zap size={14} class="text-green-600 mr-1 inline" />
											{/if}
											{task.title}
										</h3>
									</div>
									<span class="text-xs {badgeClass.bg} {badgeClass.text} px-2 py-1 rounded font-semibold inline-flex items-center gap-1">
										<Circle size={12} class={badgeClass.color} /> {getPriorityLabel(task.priority)}
									</span>
								</div>
								{#if task.description}
									<p class="text-sm text-gray-400 mb-2">{task.description}</p>
								{/if}
								{#if task.tags}
									<div class="flex flex-wrap gap-1 mb-2">
										{#each task.tags.split(',') as tag}
											<span class="text-xs bg-blue-950 text-blue-300 px-2 py-0.5 rounded">#{tag.trim()}</span>
										{/each}
									</div>
								{/if}
								<div class="flex justify-between items-center mt-3 pt-2 border-t border-gray-700 flex-wrap gap-2">
									<button
										onclick={() => handleMoveToStatus(task, 'done')}
										class="text-xs text-green-600 hover:underline font-medium inline-flex items-center gap-1"
									>
										<CheckCircle size={14} /> Complete
									</button>
									<div class="flex gap-2">
										<button
											onclick={() => handleMoveToStatus(task, 'todo')}
											class="text-xs text-gray-400 hover:underline inline-flex items-center gap-1"
										>
											<ArrowLeft size={14} /> Back
										</button>
										<button
											onclick={() => startEditTask(task)}
											class="text-xs text-purple-600 hover:underline font-medium inline-flex items-center gap-1"
										>
											<Edit2 size={14} /> Edit
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
						{/if}
					{/each}
					{#if inProgressTasks.length === 0}
						<p class="text-sm text-gray-400 text-center py-8">No tasks in progress</p>
					{/if}
				</div>
			</section>
			
			<!-- Done Column -->
			<section class="bg-green-950 rounded-2xl p-4 border-2 border-green-900 min-w-[280px] snap-start md:min-w-0">
				<h2 class="text-lg font-semibold mb-4 text-green-200 inline-flex items-center gap-2"><CheckCircle size={18} /> Done ({doneTasks.length})</h2>
				<div class="space-y-3">
					{#each doneTasks as task}
						{#if editingTaskId === task.id && editingTask}
							<div class="bg-gray-900 rounded-lg shadow p-4 border-l-4 border-green-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-green-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-green-500 text-sm"
									placeholder="Description"
									rows="2"
								></textarea>
								<div class="flex gap-2">
									<button
										onclick={handleSaveEditedTask}
										class="flex-1 bg-green-600 text-white px-3 py-2 rounded hover:bg-green-700 text-sm font-medium"
									>
										Save
									</button>
									<button
										onclick={cancelEditTask}
										class="flex-1 bg-gray-700 text-gray-300 px-3 py-2 rounded hover:bg-gray-600 text-sm font-medium"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<div class="bg-gray-900 rounded-lg shadow p-4 opacity-75 hover:opacity-100 transition">
								<div class="flex justify-between items-start mb-2">
									<div class="flex-1">
										<h3 class="font-medium text-white line-through mb-1">
											{#if task.is_quick_task}
												<Zap size={14} class="text-green-600 mr-1 inline" />
											{/if}
											{task.title}
										</h3></div>
								</div>
								{#if task.description}
									<p class="text-sm text-gray-400 mb-2 line-through">{task.description}</p>
								{/if}
								{#if task.tags}
									<div class="flex flex-wrap gap-1 mb-2 opacity-60">
										{#each task.tags.split(',') as tag}
											<span class="text-xs bg-gray-800 text-gray-400 px-2 py-0.5 rounded">#{tag.trim()}</span>
										{/each}
									</div>
								{/if}
								<div class="flex justify-between items-center mt-3 pt-2 border-t border-gray-700 flex-wrap gap-2">
									<button
										onclick={() => handleMoveToStatus(task, 'todo')}
										class="text-xs text-blue-600 hover:underline inline-flex items-center gap-1"
									>
										<ArrowLeft size={14} /> Reopen
									</button>
									<div class="flex gap-2">
										<button
											onclick={() => startEditTask(task)}
											class="text-xs text-purple-600 hover:underline font-medium inline-flex items-center gap-1"
										>
											<Edit2 size={14} /> Edit
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
						{/if}
					{/each}
					{#if doneTasks.length === 0}
						<p class="text-sm text-gray-400 text-center py-8">No completed tasks</p>
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
				data-quadrant="urgent-important"
				class="bg-red-950 border-2 border-red-900 rounded-lg p-5 min-h-[400px] transition-all touch-pan-y {draggedFromQuadrant !== 'urgent-important' && draggedTask ? 'ring-4 ring-green-400' : ''} {touchTargetQuadrant === 'urgent-important' ? 'ring-4 ring-green-400 bg-red-900' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'urgent-important')}
			>
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-lg font-semibold text-red-200 flex items-center gap-2">
						<AlertCircle size={22} class="text-red-600" /> Do First
						<span class="text-sm font-normal text-red-300 block">({urgentImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'urgent-important' ? null : 'urgent-important'}
						class="text-red-200 hover:bg-red-900 active:bg-red-800 p-3 min-w-[44px] min-h-[44px] rounded-lg transition touch-manipulation flex items-center justify-center"
						title="Quick add task"
					>
						<Plus size={20} />
					</button>
				</div>

				{#if quickAddQuadrant === 'urgent-important'}
					<div class="mb-4 p-4 bg-gray-900 rounded-lg border-2 border-red-300">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Quick task..."
							class="w-full px-4 py-3 min-h-[48px] bg-gray-800 border border-red-900 rounded-lg text-base mb-3 text-white focus:outline-none focus:ring-2 focus:ring-red-500 touch-manipulation"
						/>
						<div class="flex gap-3">
							<button
								onclick={handleQuickAddToQuadrant}
								class="flex-1 bg-red-600 text-white px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-red-700 active:bg-red-800 font-medium touch-manipulation"
							>
								Add
							</button>
							<button
								onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }}
								class="flex-1 bg-gray-700 text-gray-300 px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
							>
								Cancel
							</button>
						</div>
					</div>
				{/if}

				<div class="space-y-2">
					{#each urgentImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<!-- Edit Mode -->
							<div class="bg-gray-900 rounded-lg shadow-sm p-3 border-2 border-red-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 resize-none"
									placeholder="Description"
									rows="2"
								/>
								<div class="flex gap-2">
									<button
										onclick={handleSaveEditedTask}
										class="flex-1 bg-red-600 text-white px-2 py-1 rounded text-sm hover:bg-red-700 font-medium"
									>
										Save
									</button>
									<button
										onclick={cancelEditTask}
										class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-sm hover:bg-gray-600"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<!-- Display Mode -->
							<div 
								class="bg-gray-900 rounded-lg shadow-sm p-4 border-l-4 border-red-500 cursor-move hover:shadow-md transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50 scale-95' : ''} {isTouchDragging && touchMoveTask?.id === task.id ? 'shadow-2xl z-10' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'urgent-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'urgent-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex justify-between items-start gap-3">
									<div class="flex-1">
										<h3 class="font-medium text-white text-base leading-relaxed">{task.title}</h3></div>
									<button
										onclick={() => startEditTask(task)}
										class="text-purple-400 hover:text-purple-300 active:text-purple-500 p-2 min-w-[40px] min-h-[40px] flex items-center justify-center touch-manipulation"
										title="Edit"
									>
										<Edit2 size={18} />
									</button>
								</div>
								<div class="flex gap-3 mt-3">
									<button onclick={() => handleDeleteTask(task.id)} class="text-sm text-red-400 hover:text-red-300 active:underline py-2 px-3 min-h-[40px] touch-manipulation">Delete</button>
								</div>
							</div>
						{/if}
					{/each}
					{#if urgentImportant.length === 0 && quickAddQuadrant !== 'urgent-important'}
						<p class="text-xs text-gray-400 text-center py-4">Drop tasks here</p>
					{/if}
				</div>
			</div>
			
			<!-- Not Urgent & Important (Schedule) -->
			<div 
				data-quadrant="not-urgent-important"
				class="bg-blue-950 border-2 border-blue-900 rounded-lg p-5 min-h-[400px] transition-all touch-pan-y {draggedFromQuadrant !== 'not-urgent-important' && draggedTask ? 'ring-4 ring-green-400' : ''} {touchTargetQuadrant === 'not-urgent-important' ? 'ring-4 ring-green-400 bg-blue-900' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'not-urgent-important')}
			>
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-lg font-semibold text-blue-200 flex items-center gap-2">
						<Calendar size={22} class="text-blue-600" /> Schedule
						<span class="text-sm font-normal text-blue-300 block">({notUrgentImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'not-urgent-important' ? null : 'not-urgent-important'}
						class="text-blue-200 hover:bg-blue-900 active:bg-blue-800 p-3 min-w-[44px] min-h-[44px] rounded-lg transition touch-manipulation flex items-center justify-center"
						title="Quick add task"
					>
						<Plus size={20} />
					</button>
				</div>

				{#if quickAddQuadrant === 'not-urgent-important'}
					<div class="mb-4 p-4 bg-gray-900 rounded-lg border-2 border-blue-300">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Quick task..."
							class="w-full px-4 py-3 min-h-[48px] bg-gray-800 border border-blue-900 rounded-lg text-base mb-3 text-white focus:outline-none focus:ring-2 focus:ring-blue-500 touch-manipulation"
						/>
						<div class="flex gap-3">
							<button
								onclick={handleQuickAddToQuadrant}
								class="flex-1 bg-blue-600 text-white px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-blue-700 active:bg-blue-800 font-medium touch-manipulation"
							>
								Add
							</button>
							<button
								onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }}
								class="flex-1 bg-gray-700 text-gray-300 px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
							>
								Cancel
							</button>
						</div>
					</div>
				{/if}

				<div class="space-y-2">
					{#each notUrgentImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<!-- Edit Mode -->
							<div class="bg-gray-900 rounded-lg shadow-sm p-3 border-2 border-blue-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
									placeholder="Description"
									rows="2"
								/>
								<div class="flex gap-2">
									<button
										onclick={handleSaveEditedTask}
										class="flex-1 bg-blue-600 text-white px-2 py-1 rounded text-sm hover:bg-blue-700 font-medium"
									>
										Save
									</button>
									<button
										onclick={cancelEditTask}
										class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-sm hover:bg-gray-600"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<!-- Display Mode -->
							<div 
								class="bg-gray-900 rounded-lg shadow-sm p-4 border-l-4 border-blue-500 cursor-move hover:shadow-md transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50 scale-95' : ''} {isTouchDragging && touchMoveTask?.id === task.id ? 'shadow-2xl z-10' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'not-urgent-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'not-urgent-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex justify-between items-start gap-3">
									<div class="flex-1">
										<h3 class="font-medium text-white text-base leading-relaxed">{task.title}</h3></div>
									<button
										onclick={() => startEditTask(task)}
										class="text-purple-400 hover:text-purple-300 active:text-purple-500 p-2 min-w-[40px] min-h-[40px] flex items-center justify-center touch-manipulation"
										title="Edit"
									>
										<Edit2 size={18} />
									</button>
								</div>
								<div class="flex gap-3 mt-3">
									<button onclick={() => handleDeleteTask(task.id)} class="text-sm text-red-400 hover:text-red-300 active:underline py-2 px-3 min-h-[40px] touch-manipulation">Delete</button>
								</div>
							</div>
						{/if}
					{/each}
					{#if notUrgentImportant.length === 0 && quickAddQuadrant !== 'not-urgent-important'}
						<p class="text-xs text-gray-400 text-center py-4">Drop tasks here</p>
					{/if}
				</div>
			</div>
			
			<!-- Urgent & Not Important (Delegate) -->
			<div 
				data-quadrant="urgent-not-important"
				class="bg-yellow-950 border-2 border-yellow-800 rounded-lg p-5 min-h-[400px] transition-all touch-pan-y {draggedFromQuadrant !== 'urgent-not-important' && draggedTask ? 'ring-4 ring-green-400' : ''} {touchTargetQuadrant === 'urgent-not-important' ? 'ring-4 ring-green-400 bg-yellow-900' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'urgent-not-important')}
			>
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-lg font-semibold text-yellow-300 flex items-center gap-2">
						<Users size={22} class="text-yellow-600" /> Delegate
						<span class="text-sm font-normal text-yellow-300 block">({urgentNotImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'urgent-not-important' ? null : 'urgent-not-important'}
						class="text-yellow-200 hover:bg-yellow-900 active:bg-yellow-800 p-3 min-w-[44px] min-h-[44px] rounded-lg transition touch-manipulation flex items-center justify-center"
						title="Quick add task"
					>
						<Plus size={20} />
					</button>
				</div>

				{#if quickAddQuadrant === 'urgent-not-important'}
					<div class="mb-4 p-4 bg-gray-900 rounded-lg border-2 border-yellow-300">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Quick task..."
							class="w-full px-4 py-3 min-h-[48px] bg-gray-800 border border-yellow-800 rounded-lg text-base mb-3 text-white focus:outline-none focus:ring-2 focus:ring-yellow-500 touch-manipulation"
						/>
						<div class="flex gap-3">
							<button
								onclick={handleQuickAddToQuadrant}
								class="flex-1 bg-yellow-600 text-white px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-yellow-700 active:bg-yellow-800 font-medium touch-manipulation"
							>
								Add
							</button>
							<button
								onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }}
								class="flex-1 bg-gray-700 text-gray-300 px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
							>
								Cancel
							</button>
						</div>
					</div>
				{/if}

				<div class="space-y-2">
					{#each urgentNotImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<!-- Edit Mode -->
							<div class="bg-gray-900 rounded-lg shadow-sm p-3 border-2 border-yellow-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-yellow-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-yellow-500 resize-none"
									placeholder="Description"
									rows="2"
								/>
								<div class="flex gap-2">
									<button
										onclick={handleSaveEditedTask}
										class="flex-1 bg-yellow-600 text-white px-2 py-1 rounded text-sm hover:bg-yellow-700 font-medium"
									>
										Save
									</button>
									<button
										onclick={cancelEditTask}
										class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-sm hover:bg-gray-600"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<!-- Display Mode -->
							<div 
								class="bg-gray-900 rounded-lg shadow-sm p-4 border-l-4 border-yellow-500 cursor-move hover:shadow-md transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50 scale-95' : ''} {isTouchDragging && touchMoveTask?.id === task.id ? 'shadow-2xl z-10' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'urgent-not-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'urgent-not-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex justify-between items-start gap-3">
									<div class="flex-1">
										<h3 class="font-medium text-white text-base leading-relaxed">{task.title}</h3></div>
									<button
										onclick={() => startEditTask(task)}
										class="text-purple-400 hover:text-purple-300 active:text-purple-500 p-2 min-w-[40px] min-h-[40px] flex items-center justify-center touch-manipulation"
										title="Edit"
									>
										<Edit2 size={18} />
									</button>
								</div>
								<div class="flex gap-3 mt-3">
									<button onclick={() => handleDeleteTask(task.id)} class="text-sm text-red-400 hover:text-red-300 active:underline py-2 px-3 min-h-[40px] touch-manipulation">Delete</button>
								</div>
							</div>
						{/if}
					{/each}
					{#if urgentNotImportant.length === 0 && quickAddQuadrant !== 'urgent-not-important'}
						<p class="text-xs text-gray-400 text-center py-4">Drop tasks here</p>
					{/if}
				</div>
			</div>
			
			<!-- Not Urgent & Not Important (Eliminate) -->
			<div 
				data-quadrant="not-urgent-not-important"
				class="bg-gray-950 border-2 border-gray-700 rounded-lg p-5 min-h-[400px] transition-all touch-pan-y {draggedFromQuadrant !== 'not-urgent-not-important' && draggedTask ? 'ring-4 ring-green-400' : ''} {touchTargetQuadrant === 'not-urgent-not-important' ? 'ring-4 ring-green-400 bg-gray-800' : ''}"
				ondragover={handleDragOver}
				ondrop={(e) => handleDrop(e, 'not-urgent-not-important')}
			>
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-lg font-semibold text-white flex items-center gap-2">
						<Trash size={22} class="text-gray-400" /> Eliminate
						<span class="text-sm font-normal text-gray-400 block">({notUrgentNotImportant.length})</span>
					</h2>
					<button
						onclick={() => quickAddQuadrant = quickAddQuadrant === 'not-urgent-not-important' ? null : 'not-urgent-not-important'}
						class="text-gray-400 hover:bg-gray-800 active:bg-gray-700 p-3 min-w-[44px] min-h-[44px] rounded-lg transition touch-manipulation flex items-center justify-center"
						title="Quick add task"
					>
						<Plus size={20} />
					</button>
				</div>

				{#if quickAddQuadrant === 'not-urgent-not-important'}
					<div class="mb-4 p-4 bg-gray-900 rounded-lg border-2 border-gray-300">
						<input
							type="text"
							bind:value={quickAddTitle}
							placeholder="Quick task..."
							class="w-full px-4 py-3 min-h-[48px] bg-gray-800 border border-gray-700 rounded-lg text-base mb-3 text-white focus:outline-none focus:ring-2 focus:ring-gray-500 touch-manipulation"
						/>
						<div class="flex gap-3">
							<button
								onclick={handleQuickAddToQuadrant}
								class="flex-1 bg-gray-600 text-white px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-gray-700 active:bg-gray-800 font-medium touch-manipulation"
							>
								Add
							</button>
							<button
								onclick={() => { quickAddQuadrant = null; quickAddTitle = ''; }}
								class="flex-1 bg-gray-700 text-gray-300 px-4 py-3 min-h-[48px] rounded-lg text-base hover:bg-gray-600 active:bg-gray-500 touch-manipulation"
							>
								Cancel
							</button>
						</div>
					</div>
				{/if}

				<div class="space-y-2">
					{#each notUrgentNotImportant as task}
						{#if editingTaskId === task.id && editingTask}
							<!-- Edit Mode -->
							<div class="bg-gray-900 rounded-lg shadow-sm p-3 border-2 border-gray-500">
								<input
									type="text"
									bind:value={editingTask.title}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-500"
									placeholder="Task title"
								/>
								<textarea
									bind:value={editingTask.description}
									class="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-sm mb-2 text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-500 resize-none"
									placeholder="Description"
									rows="2"
								/>
								<div class="flex gap-2">
									<button
										onclick={handleSaveEditedTask}
										class="flex-1 bg-gray-600 text-white px-2 py-1 rounded text-sm hover:bg-gray-700 font-medium"
									>
										Save
									</button>
									<button
										onclick={cancelEditTask}
										class="flex-1 bg-gray-700 text-gray-300 px-2 py-1 rounded text-sm hover:bg-gray-600"
									>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<!-- Display Mode -->
							<div 
								class="bg-gray-900 rounded-lg shadow-sm p-4 border-l-4 border-gray-400 cursor-move hover:shadow-md transition touch-none select-none {draggedTask?.id === task.id ? 'opacity-50 scale-95' : ''} {isTouchDragging && touchMoveTask?.id === task.id ? 'shadow-2xl z-10' : ''}"
								draggable="true"
								ondragstart={() => handleDragStart(task, 'not-urgent-not-important')}
								ondragend={handleDragEnd}
								ontouchstart={(e) => handleTouchStart(e, task, 'not-urgent-not-important')}
								ontouchmove={handleTouchMove}
								ontouchend={handleTouchEnd}
							>
								<div class="flex justify-between items-start gap-3">
									<div class="flex-1">
										<h3 class="font-medium text-white text-base leading-relaxed">{task.title}</h3></div>
									<button
										onclick={() => startEditTask(task)}
										class="text-purple-400 hover:text-purple-300 active:text-purple-500 p-2 min-w-[40px] min-h-[40px] flex items-center justify-center touch-manipulation"
										title="Edit"
									>
										<Edit2 size={18} />
									</button>
								</div>
								<div class="flex gap-3 mt-3">
									<button onclick={() => handleDeleteTask(task.id)} class="text-sm text-red-400 hover:text-red-300 active:underline py-2 px-3 min-h-[40px] touch-manipulation">Delete</button>
								</div>
							</div>
						{/if}
					{/each}
					{#if notUrgentNotImportant.length === 0 && quickAddQuadrant !== 'not-urgent-not-important'}
						<p class="text-xs text-gray-400 text-center py-4">Drop tasks here</p>
					{/if}
				</div>
			</div>
		</div>
		
		<div class="bg-blue-950 border border-blue-900 rounded-lg p-4 mt-4">
			<h3 class="font-semibold text-blue-200 mb-2 flex items-center gap-2"><Lightbulb size={20} class="text-blue-600" /> Eisenhower Matrix Guide</h3>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm text-blue-300">
				<div><strong>Do First:</strong> Critical deadlines & crises</div>
				<div><strong>Schedule:</strong> Long-term goals & planning</div>
				<div><strong>Delegate:</strong> Interruptions & some emails</div>
				<div><strong>Eliminate:</strong> Time wasters & distractions</div>
			</div>
			<p class="text-xs text-blue-300 mt-2 flex items-center gap-1">
				<Lightbulb size={14} class="text-blue-600" /> Tasks are automatically categorized based on priority and due date (within 3 days = urgent)
			</p>
		</div>
	{/if}
</div>
