<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { taskAPI } from '$lib/api';
	import { wsService } from '$lib/websocket';
	import { Calendar as CalendarIcon, ChevronLeft, ChevronRight, Circle, CheckCircle, Plus, Edit2, Trash2, X } from 'lucide-svelte';

	type TaskStatus = 'todo' | 'in-progress' | 'done';
	type TaskPriority = 'low' | 'medium' | 'high';

	type Task = {
		id: number;
		title: string;
		description?: string | null;
		priority: TaskPriority;
		status: TaskStatus;
		due_date?: string | null;
		tags?: string | null;
		is_quick_task?: boolean;
		is_archived?: boolean;
		quadrant?: string | null;
		task_type?: string;
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
		task_type: string;
	};

	let tasks = $state<Task[]>([]);
	let loading = $state(true);
	let currentMonth = $state(new Date().getMonth());
	let currentYear = $state(new Date().getFullYear());
	let selectedDate = $state<number | null>(null);
	let showAddTask = $state(false);
	let editingTask = $state<Task | null>(null);
	let newTask = $state<NewTask>({
		title: '',
		description: '',
		priority: 'medium',
		status: 'todo',
		tags: '',
		due_date: '',
		is_quick_task: false,
		quadrant: '',
		task_type: 'kanban'
	});

	const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
		'July', 'August', 'September', 'October', 'November', 'December'];
	const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

	onMount(async () => {
		const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
		if (!token) {
			goto('/login');
			return;
		}
		try {
			await loadTasks();
		} catch (error: any) {
			console.error('Failed to load tasks:', error);
			if (error.message?.includes('Authorization') || error.message?.includes('401')) {
				localStorage.removeItem('authToken');
				goto('/login');
			}
		} finally {
			loading = false;
		}

		// Re-fetch whenever another device creates / updates / deletes a task
		const refresh = () => loadTasks();
		wsService.on('task_create', refresh);
		wsService.on('task_update', refresh);
		wsService.on('task_delete', refresh);

		return () => {
			wsService.off('task_create', refresh);
			wsService.off('task_update', refresh);
			wsService.off('task_delete', refresh);
		};
	});

	function previousMonth() {
		if (currentMonth === 0) {
			currentMonth = 11;
			currentYear--;
		} else {
			currentMonth--;
		}
		selectedDate = null;
	}

	function nextMonth() {
		if (currentMonth === 11) {
			currentMonth = 0;
			currentYear++;
		} else {
			currentMonth++;
		}
		selectedDate = null;
	}

	function goToToday() {
		const today = new Date();
		currentMonth = today.getMonth();
		currentYear = today.getFullYear();
		selectedDate = today.getDate();
	}

	function getDaysInMonth(month: number, year: number): number {
		return new Date(year, month + 1, 0).getDate();
	}

	function getFirstDayOfMonth(month: number, year: number): number {
		return new Date(year, month, 1).getDay();
	}

	function getTasksForDate(day: number): Task[] {
		const dateStr = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
		return tasks.filter(task => task.due_date?.startsWith(dateStr));
	}

	function isDateToday(day: number): boolean {
		const today = new Date();
		return day === today.getDate() && 
			   currentMonth === today.getMonth() && 
			   currentYear === today.getFullYear();
	}

	function selectDate(day: number | null) {
		selectedDate = day;
	}

	function openAddTaskModal(day?: number) {
		const targetDay = day || selectedDate;
		if (targetDay) {
			const dateStr = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(targetDay).padStart(2, '0')}`;
			newTask.due_date = dateStr;
		}
		showAddTask = true;
	}

	function closeAddTaskModal() {
		showAddTask = false;
		resetNewTask();
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
			quadrant: '',
			task_type: 'kanban'
		};
	}

	async function handleAddTask() {
		if (!newTask.title.trim()) return;
		
		try {
			// Convert date string to ISO 8601 format if present
			let formattedDueDate = undefined;
			if (newTask.due_date && newTask.due_date.trim() !== '') {
				// Convert YYYY-MM-DD to ISO timestamp at midnight UTC
				formattedDueDate = new Date(newTask.due_date + 'T00:00:00Z').toISOString();
			}
			
			const taskData = {
				title: newTask.title,
				description: newTask.description,
				priority: newTask.priority,
				status: newTask.status,
				tags: newTask.tags,
				due_date: formattedDueDate,
				is_quick_task: newTask.is_quick_task,
				quadrant: newTask.quadrant,
				task_type: newTask.task_type
			};
			
			console.log('Creating task with data:', taskData);
			await taskAPI.create(taskData);
			await loadTasks();
			closeAddTaskModal();
		} catch (error) {
			console.error('Failed to create task:', error);
			alert('Failed to create task. Please try again.');
		}
	}

	function startEditTask(task: Task) {
		editingTask = { ...task };
	}

	function cancelEditTask() {
		editingTask = null;
	}

	async function handleSaveEditedTask() {
		if (!editingTask || !editingTask.title.trim()) return;
		
		try {
			// Convert date string to ISO 8601 format if present
			let formattedDueDate = undefined;
			if (editingTask.due_date && editingTask.due_date.trim() !== '') {
				// Extract just the date part if it already has time component
				const datePart = editingTask.due_date.split('T')[0];
				formattedDueDate = new Date(datePart + 'T00:00:00Z').toISOString();
			}
			
			const taskData = {
				title: editingTask.title,
				description: editingTask.description || '',
				priority: editingTask.priority,
				status: editingTask.status,
				tags: editingTask.tags || '',
				due_date: formattedDueDate,
				is_quick_task: editingTask.is_quick_task || false,
				quadrant: editingTask.quadrant || '',
				task_type: editingTask.task_type || 'kanban'
			};
			
			console.log('Updating task with data:', taskData);
			await taskAPI.update(editingTask.id, taskData);
			await loadTasks();
			cancelEditTask();
		} catch (error) {
			console.error('Failed to update task:', error);
			alert('Failed to update task. Please try again.');
		}
	}

	async function handleDeleteTask(id: number) {
		if (!confirm('Are you sure you want to delete this task?')) return;
		
		try {
			await taskAPI.delete(id);
			await loadTasks();
		} catch (error) {
			console.error('Failed to delete task:', error);
			alert('Failed to delete task. Please try again.');
		}
	}

	async function loadTasks() {
		try {
			const tasksData = await taskAPI.getAll();
			tasks = tasksData.filter((task: Task) => !task.is_archived && task.due_date);
		} catch (error) {
			console.error('Failed to load tasks:', error);
		}
	}

	function getCalendarDays(): (number | null)[] {
		const daysInMonth = getDaysInMonth(currentMonth, currentYear);
		const firstDay = getFirstDayOfMonth(currentMonth, currentYear);
		const days: (number | null)[] = [];
		
		for (let i = 0; i < firstDay; i++) {
			days.push(null);
		}
		
		for (let day = 1; day <= daysInMonth; day++) {
			days.push(day);
		}
		
		return days;
	}

	function getStatusIcon(status: TaskStatus) {
		if (status === 'done') return CheckCircle;
		if (status === 'in-progress') return Clock;
		return Circle;
	}

	function getStatusColor(status: TaskStatus) {
		if (status === 'done') return 'text-green-600';
		if (status === 'in-progress') return 'text-blue-600';
		return 'text-gray-400';
	}

	function getPriorityColor(priority: TaskPriority) {
		if (priority === 'high') return 'border-l-red-500 bg-red-950';
		if (priority === 'medium') return 'border-l-yellow-500 bg-yellow-950';
		return 'border-l-gray-300 bg-gray-800';
	}

	const calendarDays = $derived(getCalendarDays());
	const selectedDateTasks = $derived(
		selectedDate ? getTasksForDate(selectedDate) : []
	);
	// Short single-letter day names for mobile
	const dayNamesShort = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];
</script>

<svelte:head>
	<title>Calendar - Tonish</title>
</svelte:head>

<div class="max-w-[1600px] mx-auto">
	{#if loading}
		<div class="text-center py-12">
			<p class="text-gray-400 text-sm">Loading calendar...</p>
		</div>
	{:else}
		<!-- Main layout: 1-col mobile → 3-col tablet → 3-col desktop -->
		<div class="grid grid-cols-1 sm:grid-cols-3 gap-2 sm:gap-3">

			<!-- ── Calendar Panel ── -->
			<div class="sm:col-span-2 bg-gray-900 border border-gray-800 rounded-lg p-2.5 sm:p-3">

				<!-- Calendar Header -->
				<div class="flex items-center justify-between mb-2.5">
					<h2 class="text-sm font-bold text-white leading-none">
						{monthNames[currentMonth]} <span class="text-gray-400 font-normal">{currentYear}</span>
					</h2>
					<div class="flex items-center gap-1">
						<button
							onclick={() => openAddTaskModal()}
							class="h-8 px-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition flex items-center gap-1 text-xs font-medium touch-manipulation"
						><Plus size={12} /> <span class="hidden sm:inline">Add</span></button>
						<button
							onclick={goToToday}
							class="h-8 px-2 text-xs bg-blue-950 text-blue-300 rounded-md hover:bg-blue-900 font-medium transition touch-manipulation"
						>Today</button>
						<button
							onclick={previousMonth}
							class="w-8 h-8 flex items-center justify-center hover:bg-gray-800 rounded-md transition text-gray-300 touch-manipulation"
							aria-label="Previous month"
						><ChevronLeft size={15} /></button>
						<button
							onclick={nextMonth}
							class="w-8 h-8 flex items-center justify-center hover:bg-gray-800 rounded-md transition text-gray-300 touch-manipulation"
							aria-label="Next month"
						><ChevronRight size={15} /></button>
					</div>
				</div>

				<!-- Calendar Grid -->
				<div class="grid grid-cols-7 gap-px sm:gap-1">
					<!-- Day Names -->
					{#each dayNames as dayName, i}
						<div class="text-center py-1.5">
							<!-- Single letter on mobile, 3-letter on sm+ -->
							<span class="text-[10px] font-semibold text-gray-500 sm:hidden">{dayNamesShort[i]}</span>
							<span class="hidden sm:inline text-xs font-semibold text-gray-400">{dayName}</span>
						</div>
					{/each}

					<!-- Calendar Days -->
					{#each calendarDays as day}
						{#if day === null}
							<div class="aspect-square"></div>
						{:else}
							{@const dayTasks = getTasksForDate(day)}
							{@const isToday = isDateToday(day)}
							{@const isSelected = selectedDate === day}
							<button
								onclick={() => selectDate(day)}
								class="aspect-square rounded sm:rounded-md transition-all touch-manipulation flex flex-col items-center justify-start pt-1 sm:pt-1.5 px-0.5
									{isToday ? 'bg-blue-600 ring-1 ring-blue-400' : ''}
									{isSelected && !isToday ? 'bg-blue-950 ring-1 ring-blue-500' : ''}
									{!isToday && !isSelected ? 'hover:bg-gray-800 active:bg-gray-700' : ''}"
							>
								<span class="text-[11px] sm:text-sm font-medium leading-none
									{isToday ? 'text-white' : isSelected ? 'text-blue-300' : 'text-gray-300'}">
									{day}
								</span>
								{#if dayTasks.length > 0}
									<div class="flex gap-px mt-0.5 sm:mt-1">
										{#each dayTasks.slice(0, 3) as task}
											<div class="w-1 h-1 sm:w-1.5 sm:h-1.5 rounded-full
												{isToday ? 'bg-white' :
												task.status === 'done' ? 'bg-green-500' :
												task.status === 'in-progress' ? 'bg-blue-400' : 'bg-gray-500'}">
											</div>
										{/each}
										{#if dayTasks.length > 3}
											<span class="text-[8px] font-bold {isToday ? 'text-white' : 'text-gray-500'}">+</span>
										{/if}
									</div>
								{/if}
							</button>
						{/if}
					{/each}
				</div>

				<!-- Legend + status bar (compact, single row) -->
				<div class="mt-2 pt-2 border-t border-gray-800 flex items-center gap-3 flex-wrap">
					<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-gray-500"></div><span class="text-[10px] text-gray-500">To Do</span></div>
					<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-blue-400"></div><span class="text-[10px] text-gray-500">Active</span></div>
					<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-green-500"></div><span class="text-[10px] text-gray-500">Done</span></div>
				</div>
			</div>

			<!-- ── Sidebar ── -->
			<div class="sm:col-span-1 bg-gray-900 border border-gray-800 rounded-lg p-2.5 sm:p-3 flex flex-col gap-2">

				<!-- Monthly Summary (compact 4-col horizontal row) -->
				{#if tasks.length > 0}
					{@const monthTasks = tasks.filter(task => {
						if (!task.due_date) return false;
						const taskDate = new Date(task.due_date);
						return taskDate.getMonth() === currentMonth && taskDate.getFullYear() === currentYear;
					})}
					<div class="bg-gray-800 rounded-md p-2 border border-gray-700">
						<p class="text-[10px] font-semibold text-gray-400 mb-1.5 flex items-center gap-1">
							<CalendarIcon size={10} class="text-blue-400" />
							{monthNames[currentMonth]} Summary
						</p>
						<div class="grid grid-cols-4 gap-1">
							<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
								<span class="text-base font-bold text-white leading-none">{monthTasks.length}</span>
								<span class="text-[9px] text-gray-500 mt-0.5">Total</span>
							</div>
							<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
								<span class="text-base font-bold text-gray-400 leading-none">{monthTasks.filter(t => t.status === 'todo').length}</span>
								<span class="text-[9px] text-gray-500 mt-0.5">Todo</span>
							</div>
							<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
								<span class="text-base font-bold text-blue-400 leading-none">{monthTasks.filter(t => t.status === 'in-progress').length}</span>
								<span class="text-[9px] text-gray-500 mt-0.5">Active</span>
							</div>
							<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
								<span class="text-base font-bold text-green-400 leading-none">{monthTasks.filter(t => t.status === 'done').length}</span>
								<span class="text-[9px] text-gray-500 mt-0.5">Done</span>
							</div>
						</div>
					</div>
				{/if}

				<!-- Selected Date Tasks -->
				<div class="flex-1 min-h-0">
					<div class="flex items-center justify-between mb-2">
						<h3 class="text-xs font-bold text-white leading-none">
							{#if selectedDate}
								{monthNames[currentMonth]} {selectedDate}
							{:else}
								Select a date
							{/if}
						</h3>
						{#if selectedDate}
							<button
								onclick={() => openAddTaskModal(selectedDate)}
								class="h-7 px-2 bg-blue-600 hover:bg-blue-700 text-white rounded text-[10px] font-medium flex items-center gap-1 touch-manipulation"
							><Plus size={10} /> Add</button>
						{/if}
					</div>

					{#if selectedDate}
						{#if selectedDateTasks.length > 0}
							<div class="space-y-1.5 overflow-y-auto max-h-[50vh] sm:max-h-[calc(100vh-320px)] pr-0.5">
								{#each selectedDateTasks as task}
									{@const StatusIcon = getStatusIcon(task.status)}
									{#if editingTask && editingTask.id === task.id}
										<!-- Edit Mode -->
										<div class="border border-blue-500 rounded-lg p-2.5 bg-gray-800">
											<div class="space-y-2">
												<input
													type="text"
													bind:value={editingTask.title}
													class="w-full min-h-[36px] px-3 py-2 bg-gray-900 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-blue-500 focus:border-transparent"
													placeholder="Task title"
												/>
												<textarea
													bind:value={editingTask.description}
													class="w-full px-3 py-2 bg-gray-900 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-blue-500 focus:border-transparent"
													placeholder="Description"
													rows="2"
												></textarea>
												<div class="grid grid-cols-2 gap-1.5">
													<select bind:value={editingTask.priority}
														class="min-h-[36px] px-2 py-1.5 bg-gray-900 border border-gray-700 text-white rounded-md text-xs focus:ring-1 focus:ring-blue-500">
														<option value="low">Low</option>
														<option value="medium">Medium</option>
														<option value="high">High</option>
													</select>
													<select bind:value={editingTask.status}
														class="min-h-[36px] px-2 py-1.5 bg-gray-900 border border-gray-700 text-white rounded-md text-xs focus:ring-1 focus:ring-blue-500">
														<option value="todo">To Do</option>
														<option value="in-progress">In Progress</option>
														<option value="done">Done</option>
													</select>
												</div>
												<input type="date" bind:value={editingTask.due_date}
													class="w-full min-h-[36px] px-3 py-2 bg-gray-900 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-blue-500" />
												<div class="flex gap-1.5">
													<button onclick={handleSaveEditedTask}
														class="flex-1 h-9 bg-blue-600 text-white rounded-md text-xs font-medium hover:bg-blue-700 touch-manipulation">Save</button>
													<button onclick={cancelEditTask}
														class="flex-1 h-9 bg-gray-700 text-gray-200 rounded-md text-xs font-medium hover:bg-gray-600 touch-manipulation">Cancel</button>
												</div>
											</div>
										</div>
									{:else}
										<!-- View Mode -->
										<div class="border-l-2 {getPriorityColor(task.priority)} rounded-md p-2 flex items-start gap-2">
											<StatusIcon size={14} class="{getStatusColor(task.status)} mt-0.5 flex-shrink-0" />
											<div class="flex-1 min-w-0">
												<p class="text-xs font-semibold text-gray-200 truncate">{task.title}</p>
												{#if task.description}
													<p class="text-[10px] text-gray-400 mt-0.5 line-clamp-1">{task.description}</p>
												{/if}
												<div class="flex items-center gap-1.5 mt-1 flex-wrap">
													<span class="text-[9px] px-1.5 py-0.5 rounded-full
														{task.status === 'done' ? 'bg-green-950 text-green-300' :
														task.status === 'in-progress' ? 'bg-blue-950 text-blue-300' : 'bg-gray-700 text-gray-300'}">
														{task.status.replace('-', ' ')}
													</span>
													{#if task.priority === 'high'}
														<span class="text-[9px] px-1.5 py-0.5 rounded-full bg-red-950 text-red-300">High</span>
													{/if}
												</div>
											</div>
											<div class="flex gap-1 flex-shrink-0">
												<button onclick={() => startEditTask(task)}
													class="w-7 h-7 flex items-center justify-center hover:bg-blue-950 rounded text-blue-400 touch-manipulation"
													aria-label="Edit"><Edit2 size={13} /></button>
												<button onclick={() => handleDeleteTask(task.id)}
													class="w-7 h-7 flex items-center justify-center hover:bg-red-950 rounded text-red-400 touch-manipulation"
													aria-label="Delete"><Trash2 size={13} /></button>
											</div>
										</div>
									{/if}
								{/each}
							</div>
						{:else}
							<div class="text-center py-6">
								<CalendarIcon size={32} class="text-gray-700 mx-auto mb-2" />
								<p class="text-xs text-gray-500">No tasks for this day</p>
								<button
									onclick={() => openAddTaskModal(selectedDate)}
									class="inline-flex items-center gap-1 mt-3 h-8 px-3 bg-blue-600 text-white rounded-md text-xs font-medium hover:bg-blue-700 touch-manipulation"
								><Plus size={12} /> Add Task</button>
							</div>
						{/if}
					{:else}
						<div class="text-center py-6">
							<CalendarIcon size={32} class="text-gray-700 mx-auto mb-2" />
							<p class="text-xs text-gray-500">Tap a date to view tasks</p>
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}

	<!-- Add Task Modal (bottom-sheet style on mobile, centered on sm+) -->
	{#if showAddTask}
		<div class="fixed inset-0 bg-black/60 z-50 flex items-end sm:items-center justify-center sm:p-4" onclick={(e) => { if (e.target === e.currentTarget) closeAddTaskModal(); }}>
			<div class="bg-gray-900 border border-gray-800 rounded-t-2xl sm:rounded-2xl shadow-2xl w-full sm:max-w-md p-4 sm:p-5 max-h-[90vh] overflow-y-auto">
				<div class="flex items-center justify-between mb-3">
					<h3 class="text-sm font-bold text-white">Add New Task</h3>
					<button onclick={closeAddTaskModal}
						class="w-8 h-8 flex items-center justify-center hover:bg-gray-800 rounded-lg text-gray-400 touch-manipulation"
						aria-label="Close"><X size={18} /></button>
				</div>

				<div class="space-y-3">
					<div>
						<label for="task-title" class="block text-xs font-medium text-gray-400 mb-1">Title *</label>
						<input id="task-title" type="text" bind:value={newTask.title}
							class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="Task title" />
					</div>

					<div>
						<label for="task-description" class="block text-xs font-medium text-gray-400 mb-1">Description</label>
						<textarea id="task-description" bind:value={newTask.description}
							class="w-full px-3 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
							placeholder="Optional" rows="2"></textarea>
					</div>

					<div class="grid grid-cols-2 gap-2">
						<div>
							<label for="task-priority" class="block text-xs font-medium text-gray-400 mb-1">Priority</label>
							<select id="task-priority" bind:value={newTask.priority}
								class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500">
								<option value="low">Low</option>
								<option value="medium">Medium</option>
								<option value="high">High</option>
							</select>
						</div>
						<div>
							<label for="task-status" class="block text-xs font-medium text-gray-400 mb-1">Status</label>
							<select id="task-status" bind:value={newTask.status}
								class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500">
								<option value="todo">To Do</option>
								<option value="in-progress">In Progress</option>
								<option value="done">Done</option>
							</select>
						</div>
					</div>

					<div>
						<label for="task-due-date" class="block text-xs font-medium text-gray-400 mb-1">Due Date</label>
						<input id="task-due-date" type="date" bind:value={newTask.due_date}
							class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500" />
					</div>

					<div class="flex gap-2 pt-1">
						<button onclick={handleAddTask}
							class="flex-1 h-10 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 active:scale-95 transition touch-manipulation">
							Create Task
						</button>
						<button onclick={closeAddTaskModal}
							class="flex-1 h-10 bg-gray-700 text-gray-200 rounded-lg text-sm font-medium hover:bg-gray-600 touch-manipulation">
							Cancel
						</button>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
