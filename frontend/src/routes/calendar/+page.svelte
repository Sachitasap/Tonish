<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { taskAPI } from '$lib/api';
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
</script>

<div class="min-h-screen bg-gray-950 p-4 md:p-8">
	<div class="max-w-7xl mx-auto space-y-6">
		<!-- Header -->
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<div class="rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 p-3">
					<CalendarIcon size={32} class="text-white" />
				</div>
				<div>
					<h1 class="text-3xl md:text-4xl font-bold text-white">Task Calendar</h1>
					<p class="text-gray-400">View and manage your tasks by date</p>
				</div>
			</div>
			<div class="flex items-center gap-3">
				<button
					onclick={() => openAddTaskModal()}
					class="min-h-[44px] px-5 py-3 bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-lg shadow-md hover:shadow-lg active:scale-95 transition flex items-center gap-2 touch-manipulation"
				>
					<Plus size={22} />
					<span class="hidden md:inline font-medium">Add Task</span>
				</button>
				<a
					href="/"
					class="min-h-[44px] px-5 py-3 bg-gray-800 text-gray-200 rounded-lg shadow-sm hover:shadow-md active:scale-95 transition border border-gray-700 flex items-center touch-manipulation"
				>
					← Back
				</a>
			</div>
		</div>

		{#if loading}
			<div class="text-center py-12">
				<p class="text-gray-400">Loading calendar...</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
				<!-- Calendar View -->
				<div class="lg:col-span-2 bg-gray-900 border border-gray-800 rounded-3xl shadow-xl ring-1 ring-gray-800 p-6 md:p-8">
					<!-- Calendar Header -->
					<div class="flex items-center justify-between mb-6">
						<h2 class="text-2xl font-bold text-white">
							{monthNames[currentMonth]} {currentYear}
						</h2>
						<div class="flex items-center gap-2">
							<button
								onclick={goToToday}
								class="min-h-[44px] px-4 py-2.5 text-sm bg-blue-950 text-blue-300 rounded-lg hover:bg-blue-900 active:scale-95 font-medium transition touch-manipulation"
							>
								Today
							</button>
							<button
								onclick={previousMonth}
								class="min-w-[44px] min-h-[44px] p-2.5 hover:bg-gray-800 active:bg-gray-700 rounded-lg transition text-gray-300 touch-manipulation flex items-center justify-center"
								aria-label="Previous month"
							>
								<ChevronLeft size={26} />
							</button>
							<button
								onclick={nextMonth}
								class="min-w-[44px] min-h-[44px] p-2.5 hover:bg-gray-800 active:bg-gray-700 rounded-lg transition text-gray-300 touch-manipulation flex items-center justify-center"
								aria-label="Next month"
							>
								<ChevronRight size={26} />
							</button>
						</div>
					</div>

					<!-- Calendar Grid -->
					<div class="grid grid-cols-7 gap-2 md:gap-3">
						<!-- Day Names -->
						{#each dayNames as dayName}
							<div class="text-center text-sm font-semibold text-gray-400 py-2">
								{dayName}
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
								class="aspect-square min-h-[50px] md:min-h-[60px] rounded-xl p-3 md:p-4 transition-all relative group touch-manipulation
									{isToday ? 'bg-blue-500 text-white font-bold ring-2 ring-blue-600' : ''}
									{isSelected && !isToday ? 'bg-blue-950 ring-2 ring-blue-400' : ''}
									{!isToday && !isSelected ? 'hover:bg-gray-800 active:bg-gray-700 text-gray-300' : ''}
										{dayTasks.length > 0 ? 'font-semibold' : ''}"
								>
									<div class="flex flex-col h-full">
										<span class="text-base md:text-lg font-medium {isToday ? 'text-white' : 'text-gray-300'}">
											{day}
										</span>
										{#if dayTasks.length > 0}
											<div class="flex-1 flex items-end justify-center mt-1">
												<div class="flex gap-1">
													{#each dayTasks.slice(0, 3) as task}
														<div 
															class="w-2 h-2 rounded-full
																{task.status === 'done' ? 'bg-green-500' : ''}
																{task.status === 'in-progress' ? 'bg-blue-500' : ''}
																{task.status === 'todo' ? 'bg-gray-400' : ''}
																{isToday ? 'bg-white' : ''}"
														></div>
													{/each}
													{#if dayTasks.length > 3}
														<span class="text-xs font-bold {isToday ? 'text-white' : 'text-gray-400'}">+</span>
													{/if}
												</div>
											</div>
										{/if}
									</div>
								</button>
							{/if}
						{/each}
					</div>

					<!-- Legend -->
					<div class="mt-6 pt-6 border-t border-gray-800">
						<h3 class="text-sm font-semibold text-gray-300 mb-3">Status Legend</h3>
						<div class="flex flex-wrap gap-4 text-sm">
							<div class="flex items-center gap-2">
								<div class="w-3.5 h-3.5 rounded-full bg-gray-400"></div>
								<span class="text-gray-400">To Do</span>
							</div>
							<div class="flex items-center gap-2">
								<div class="w-3.5 h-3.5 rounded-full bg-blue-500"></div>
								<span class="text-gray-400">In Progress</span>
							</div>
							<div class="flex items-center gap-2">
								<div class="w-3.5 h-3.5 rounded-full bg-green-500"></div>
								<span class="text-gray-400">Completed</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Task Details Sidebar -->
				<div class="bg-gray-900 border border-gray-800 rounded-3xl shadow-xl ring-1 ring-gray-800 p-6">
					<!-- Monthly Summary at Top -->
					{#if tasks.length > 0}
					{@const monthTasks = tasks.filter(task => {
						if (!task.due_date) return false;
						const taskDate = new Date(task.due_date);
						return taskDate.getMonth() === currentMonth && taskDate.getFullYear() === currentYear;
					})}
					<div class="mb-6 pb-6 border-b border-gray-800 bg-gradient-to-br from-gray-800 to-gray-900 rounded-xl p-4">
						<h4 class="text-lg font-bold text-white mb-3 flex items-center gap-2">
							<CalendarIcon size={20} class="text-blue-400" />
							{monthNames[currentMonth]} Summary
						</h4>
						<div class="grid grid-cols-2 gap-3">
							<div class="bg-gray-900 rounded-lg p-3 border border-gray-700">
								<div class="text-2xl font-bold text-white">{monthTasks.length}</div>
								<div class="text-xs text-gray-400 mt-1">Total Tasks</div>
							</div>
							<div class="bg-gray-900 rounded-lg p-3 border border-gray-700">
								<div class="text-2xl font-bold text-green-400">{monthTasks.filter(t => t.status === 'done').length}</div>
								<div class="text-xs text-gray-400 mt-1">Completed</div>
							</div>
							<div class="bg-gray-900 rounded-lg p-3 border border-gray-700">
								<div class="text-2xl font-bold text-blue-400">{monthTasks.filter(t => t.status === 'in-progress').length}</div>
								<div class="text-xs text-gray-400 mt-1">In Progress</div>
							</div>
							<div class="bg-gray-900 rounded-lg p-3 border border-gray-700">
								<div class="text-2xl font-bold text-gray-400">{monthTasks.filter(t => t.status === 'todo').length}</div>
								<div class="text-xs text-gray-400 mt-1">To Do</div>
							</div>
						</div>
					</div>
					{/if}

					<h3 class="text-xl font-bold text-white mb-4">
						{#if selectedDate}
							Tasks for {monthNames[currentMonth]} {selectedDate}
						{:else}
							Select a date
						{/if}
					</h3>

					{#if selectedDate}
						{#if selectedDateTasks.length > 0}
							<div class="space-y-4">
								{#each selectedDateTasks as task}
									{@const StatusIcon = getStatusIcon(task.status)}
									{#if editingTask && editingTask.id === task.id}
										<!-- Edit Mode -->
										<div class="border-2 border-blue-500 rounded-lg p-4 shadow-md bg-gray-800">
											<div class="space-y-4">
												<input
													type="text"
													bind:value={editingTask.title}
													class="w-full min-h-[44px] px-4 py-3 bg-gray-900 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base"
													placeholder="Task title"
												/>
												<textarea
													bind:value={editingTask.description}
													class="w-full px-4 py-3 bg-gray-900 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base"
													placeholder="Description (optional)"
													rows="2"
												></textarea>
												<div class="grid grid-cols-2 gap-3">
													<select
														bind:value={editingTask.priority}
														class="min-h-[44px] px-4 py-3 bg-gray-900 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 text-base"
													>
														<option value="low">Low</option>
														<option value="medium">Medium</option>
														<option value="high">High</option>
													</select>
													<select
														bind:value={editingTask.status}
														class="min-h-[44px] px-4 py-3 bg-gray-900 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 text-base"
													>
														<option value="todo">To Do</option>
														<option value="in-progress">In Progress</option>
														<option value="done">Done</option>
													</select>
												</div>
												<input
													type="date"
													bind:value={editingTask.due_date}
													class="w-full min-h-[44px] px-4 py-3 bg-gray-900 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 text-base"
												/>
												<div class="flex gap-3">
													<button
														onclick={handleSaveEditedTask}
														class="flex-1 min-h-[48px] px-5 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 active:scale-95 transition font-medium text-base touch-manipulation"
													>
														Save
													</button>
													<button
														onclick={cancelEditTask}
														class="flex-1 min-h-[48px] px-5 py-3 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 active:bg-gray-500 transition font-medium text-base touch-manipulation"
													>
														Cancel
													</button>
												</div>
											</div>
										</div>
									{:else}
										<!-- View Mode -->
								<div class="border-l-4 {getPriorityColor(task.priority)} bg-gray-800 rounded-lg p-4 shadow-sm">
									<div class="flex items-start gap-3">
										<StatusIcon size={20} class="{getStatusColor(task.status)} mt-0.5" />
										<div class="flex-1">
											<h4 class="font-semibold text-gray-200">{task.title}</h4>
											{#if task.description}
												<p class="text-sm text-gray-400 mt-1">{task.description}</p>
													{/if}
													<div class="flex items-center gap-2 mt-2">
														<span class="text-xs px-2 py-1 rounded-full
													{task.status === 'done' ? 'bg-green-950 text-green-300' : ''}
													{task.status === 'in-progress' ? 'bg-blue-950 text-blue-300' : ''}
													{task.status === 'todo' ? 'bg-gray-800 text-gray-300' : ''}">
															{task.status.replace('-', ' ')}
														</span>
														{#if task.priority === 'high'}
															<span class="text-xs px-2 py-1 rounded-full bg-red-100 text-red-700">
																High Priority
															</span>
														{/if}
													</div>
												</div>
												<div class="flex gap-2">
													<button
														onclick={() => startEditTask(task)}
														class="min-w-[44px] min-h-[44px] p-2.5 hover:bg-blue-950 active:bg-blue-900 rounded-lg transition text-blue-400 touch-manipulation flex items-center justify-center"
														aria-label="Edit task"
													>
														<Edit2 size={20} />
													</button>
													<button
														onclick={() => handleDeleteTask(task.id)}
														class="min-w-[44px] min-h-[44px] p-2.5 hover:bg-red-950 active:bg-red-900 rounded-lg transition text-red-400 touch-manipulation flex items-center justify-center"
														aria-label="Delete task"
													>
														<Trash2 size={20} />
													</button>
												</div>
											</div>
										</div>
									{/if}
								{/each}
							</div>
						{:else}
							<div class="text-center py-8">
							<CalendarIcon size={48} class="text-gray-600 mx-auto mb-3" />
							<p class="text-gray-400">No tasks scheduled for this day</p>
								<button
									onclick={() => openAddTaskModal(selectedDate)}
									class="inline-flex items-center justify-center gap-2 mt-4 min-h-[44px] px-5 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 active:scale-95 transition text-sm font-medium touch-manipulation"
								>
									<Plus size={20} />
									Add Task
								</button>
							</div>
						{/if}
					{:else}
						<div class="text-center py-8">
							<CalendarIcon size={48} class="text-gray-600 mx-auto mb-3" />
							<p class="text-gray-400">Click on a date to view tasks</p>
						</div>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Add Task Modal -->
		{#if showAddTask}
			<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
				<div class="bg-gray-900 border border-gray-800 rounded-2xl shadow-2xl max-w-md w-full p-6">
					<div class="flex items-center justify-between mb-4">
						<h3 class="text-xl font-bold text-white">Add New Task</h3>
						<button
							onclick={closeAddTaskModal}
							class="min-w-[44px] min-h-[44px] p-2.5 hover:bg-gray-800 active:bg-gray-700 rounded-lg transition text-gray-400 touch-manipulation flex items-center justify-center"
							aria-label="Close"
						>
							<X size={24} />
						</button>
					</div>

					<div class="space-y-4">
						<div>
							<label for="task-title" class="block text-sm font-medium text-gray-300 mb-1">
								Task Title *
							</label>
							<input
								id="task-title"
								type="text"
								bind:value={newTask.title}
							class="w-full min-h-[44px] px-4 py-3 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base"
								placeholder="Enter task title"
							/>
						</div>

						<div>
							<label for="task-description" class="block text-sm font-medium text-gray-300 mb-1">
								Description
							</label>
							<textarea
								id="task-description"
								bind:value={newTask.description}
							class="w-full px-4 py-3 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-base"
								placeholder="Add details (optional)"
								rows="3"
							></textarea>
						</div>

						<div class="grid grid-cols-2 gap-3">
							<div>
								<label for="task-priority" class="block text-sm font-medium text-gray-300 mb-1">
									Priority
								</label>
								<select
									id="task-priority"
									bind:value={newTask.priority}
									class="w-full min-h-[44px] px-4 py-3 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 text-base"
								>
									<option value="low">Low</option>
									<option value="medium">Medium</option>
									<option value="high">High</option>
								</select>
							</div>

							<div>
								<label for="task-status" class="block text-sm font-medium text-gray-300 mb-1">
									Status
								</label>
								<select
									id="task-status"
									bind:value={newTask.status}
									class="w-full min-h-[44px] px-4 py-3 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 text-base"
								>
									<option value="todo">To Do</option>
									<option value="in-progress">In Progress</option>
									<option value="done">Done</option>
								</select>
							</div>
						</div>

						<div>
							<label for="task-due-date" class="block text-sm font-medium text-gray-300 mb-1">
								Due Date
							</label>
							<input
								id="task-due-date"
								type="date"
								bind:value={newTask.due_date}
								class="w-full min-h-[44px] px-4 py-3 bg-gray-800 border border-gray-700 text-white rounded-lg focus:ring-2 focus:ring-blue-500 text-base"
							/>
						</div>

						<div class="flex gap-3 pt-4">
							<button
								onclick={handleAddTask}
								class="flex-1 min-h-[48px] px-5 py-3 bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-lg hover:shadow-lg active:scale-95 transition font-medium text-base touch-manipulation"
							>
								Create Task
							</button>
							<button
								onclick={closeAddTaskModal}
								class="flex-1 min-h-[48px] px-5 py-3 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 active:bg-gray-500 transition font-medium text-base touch-manipulation"
							>
								Cancel
							</button>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
