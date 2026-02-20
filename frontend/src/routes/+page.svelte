<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { taskAPI } from '$lib/api';
	import { wsService } from '$lib/websocket';
	import { BarChart3, Zap, CheckCircle, Kanban, LayoutGrid, Calendar } from 'lucide-svelte';

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
		quadrant?: string | null;
	};

	let tasks = $state<Task[]>([]);
	let loading = $state(true);

	async function loadTasks() {
		try {
			const tasksData = await taskAPI.getAll();
			tasks = tasksData;
		} catch (error: any) {
			console.error('Failed to load data:', error);
			if (error.message?.includes('Authorization') || error.message?.includes('401')) {
				localStorage.removeItem('authToken');
				goto('/login');
			}
		} finally {
			loading = false;
		}
	}

	onMount(async () => {
		// Check if user is authenticated
		const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
		if (!token) {
			goto('/login');
			return;
		}
		await loadTasks();

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
	
	// Helper function to check if a date is today
	function isToday(date: string | null | undefined): boolean {
		if (!date) return false;
		const today = new Date();
		const taskDate = new Date(date);
		return (
			taskDate.getDate() === today.getDate() &&
			taskDate.getMonth() === today.getMonth() &&
			taskDate.getFullYear() === today.getFullYear()
		);
	}

	// Helper function to check if a date is within next N days
	function isWithinNextDays(date: string | null | undefined, days: number): boolean {
		if (!date) return false;
		const today = new Date();
		today.setHours(0, 0, 0, 0);
		const taskDate = new Date(date);
		taskDate.setHours(0, 0, 0, 0);
		const diffTime = taskDate.getTime() - today.getTime();
		const diffDays = diffTime / (1000 * 60 * 60 * 24);
		return diffDays > 0 && diffDays <= days;
	}

	// Helper function to check if a date is overdue
	function isOverdue(date: string | null | undefined, status: TaskStatus): boolean {
		if (!date || status === 'done') return false;
		const today = new Date();
		today.setHours(0, 0, 0, 0);
		const taskDate = new Date(date);
		taskDate.setHours(0, 0, 0, 0);
		return taskDate < today;
	}

	// Derived state for task statistics
	const todayTasks = $derived(
		tasks.filter((task) => isToday(task.due_date) && task.status !== 'done')
	);
	const todayCompletedTasks = $derived(
		tasks.filter((task) => isToday(task.due_date) && task.status === 'done')
	);
	const upcomingTasks = $derived(
		tasks.filter((task) => isWithinNextDays(task.due_date, 7) && task.status !== 'done')
	);
	const overdueTasks = $derived(
		tasks.filter((task) => isOverdue(task.due_date, task.status))
	);
	const inProgressTasks = $derived(tasks.filter((task) => task.status === 'in-progress').length);
	const completedTasks = $derived(tasks.filter((task) => task.status === 'done').length);

	// Kanban board task breakdown
	const kanbanTasks = $derived(tasks.filter((task) => task.task_type === 'kanban' || !task.quadrant));
	const kanbanTodo = $derived(kanbanTasks.filter((task) => task.status === 'todo').length);
	const kanbanInProgress = $derived(kanbanTasks.filter((task) => task.status === 'in-progress').length);
	const kanbanDone = $derived(kanbanTasks.filter((task) => task.status === 'done').length);

	// Eisenhower Matrix task breakdown
	const matrixTasks = $derived(tasks.filter((task) => task.task_type === 'matrix' && task.quadrant));
	const urgentImportant = $derived(matrixTasks.filter((task) => task.quadrant === 'urgent-important' && task.status !== 'done').length);
	const notUrgentImportant = $derived(matrixTasks.filter((task) => task.quadrant === 'not-urgent-important' && task.status !== 'done').length);
	const urgentNotImportant = $derived(matrixTasks.filter((task) => task.quadrant === 'urgent-not-important' && task.status !== 'done').length);
	const notUrgentNotImportant = $derived(matrixTasks.filter((task) => task.quadrant === 'not-urgent-not-important' && task.status !== 'done').length);
	const matrixCompleted = $derived(matrixTasks.filter((task) => task.status === 'done').length);
</script>

<div class="max-w-[1600px] mx-auto">
	<!-- Compact 3-Section Layout -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2 sm:gap-3">
			
			<!-- Section 1: Task Summary (compact list) -->
			<div class="lg:col-span-1">
				<div class="bg-gray-900 rounded-lg border border-gray-800 p-2">
					<h2 class="text-xs font-bold text-white mb-1.5 flex items-center gap-1.5 px-1">
						<BarChart3 size={12} class="text-blue-600" />
						Task Summary
					</h2>
					<div class="space-y-0.5">
						<!-- Today's Tasks -->
						<a href="/myflow" class="flex items-center gap-2 px-2 py-1.5 rounded-md hover:bg-gray-800 transition-colors">
							<div class="w-5 h-5 rounded bg-amber-950 flex items-center justify-center flex-shrink-0">
								<CheckCircle size={10} class="text-amber-400" />
							</div>
							<span class="text-[10px] text-gray-400 flex-1 leading-none">Today</span>
							{#if todayCompletedTasks.length > 0}
								<span class="text-[9px] text-gray-500">{todayCompletedTasks.length} done</span>
							{/if}
							<span class="text-sm font-bold text-white w-6 text-right leading-none">{todayTasks.length}</span>
						</a>

						<!-- Active Tasks -->
						<a href="/myflow" class="flex items-center gap-2 px-2 py-1.5 rounded-md hover:bg-gray-800 transition-colors">
							<div class="w-5 h-5 rounded bg-emerald-950 flex items-center justify-center flex-shrink-0">
								<Zap size={10} class="text-emerald-400" />
							</div>
							<span class="text-[10px] text-gray-400 flex-1 leading-none">Active</span>
							<span class="text-[9px] text-gray-500">{completedTasks} done</span>
							<span class="text-sm font-bold text-white w-6 text-right leading-none">{inProgressTasks}</span>
						</a>

						<!-- Upcoming Tasks -->
						<a href="/calendar" class="flex items-center gap-2 px-2 py-1.5 rounded-md hover:bg-gray-800 transition-colors">
							<div class="w-5 h-5 rounded bg-blue-950 flex items-center justify-center flex-shrink-0">
								<Calendar size={10} class="text-blue-400" />
							</div>
							<span class="text-[10px] text-gray-400 flex-1 leading-none">Upcoming</span>
							<span class="text-[9px] text-gray-500">7 days</span>
							<span class="text-sm font-bold text-white w-6 text-right leading-none">{upcomingTasks.length}</span>
						</a>

						<!-- Overdue Tasks -->
						<a href="/myflow" class="flex items-center gap-2 px-2 py-1.5 rounded-md hover:bg-gray-800 transition-colors">
							<div class="w-5 h-5 rounded bg-red-950 flex items-center justify-center flex-shrink-0">
								<Zap size={10} class="text-red-400" />
							</div>
							<span class="text-[10px] text-gray-400 flex-1 leading-none">Overdue</span>
							{#if overdueTasks.length > 0}
								<span class="text-[9px] text-red-400 font-medium">urgent</span>
							{/if}
							<span class="text-sm font-bold {overdueTasks.length > 0 ? 'text-red-400' : 'text-white'} w-6 text-right leading-none">{overdueTasks.length}</span>
						</a>
					</div>
				</div>
			</div>

			<!-- Section 2: Kanban Board Summary -->
			<div class="lg:col-span-1">
				<a href="/myflow" class="block bg-gray-900 rounded-lg border border-gray-800 p-2.5 hover:border-gray-700 hover:shadow-md transition-all h-full">
					<div class="flex items-center justify-between mb-2">
						<h2 class="text-sm font-bold text-white flex items-center gap-1.5">
							<Kanban size={13} class="text-blue-400" />
							Kanban Board
						</h2>
						<span class="text-[10px] text-gray-400">{kanbanTasks.length} tasks</span>
					</div>

					<!-- Compact Kanban Columns -->
					<div class="grid grid-cols-3 gap-1.5 mb-2">
						<div class="bg-gray-800 rounded-md p-2 border border-gray-700">
							<div class="flex items-center gap-1 mb-1">
								<div class="w-1.5 h-1.5 rounded-full bg-gray-400"></div>
								<h4 class="text-[9px] font-bold text-gray-300 uppercase tracking-wide">To Do</h4>
							</div>
							<p class="text-xl font-bold text-white leading-none">{kanbanTodo}</p>
						</div>

						<div class="bg-blue-950 rounded-md p-2 border border-blue-900">
							<div class="flex items-center gap-1 mb-1">
								<div class="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse"></div>
								<h4 class="text-[9px] font-bold text-blue-300 uppercase tracking-wide">Active</h4>
							</div>
							<p class="text-xl font-bold text-white leading-none">{kanbanInProgress}</p>
						</div>

						<div class="bg-emerald-950 rounded-md p-2 border border-emerald-900">
							<div class="flex items-center gap-1 mb-1">
								<div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
								<h4 class="text-[9px] font-bold text-emerald-300 uppercase tracking-wide">Done</h4>
							</div>
							<p class="text-xl font-bold text-white leading-none">{kanbanDone}</p>
						</div>
					</div>

					<!-- Compact Progress Bar -->
					{#if kanbanTasks.length > 0}
						<div class="pt-2 border-t border-gray-800">
							<div class="flex justify-between text-[10px] text-gray-400 mb-1.5">
								<span>Progress</span>
								<span class="font-semibold text-white">{Math.round((kanbanDone / kanbanTasks.length) * 100)}%</span>
							</div>
							<div class="h-1 bg-gray-800 rounded-full overflow-hidden">
								<div 
									class="h-full bg-gradient-to-r from-blue-500 to-emerald-500 transition-all duration-700 rounded-full"
									style="width: {(kanbanDone / kanbanTasks.length) * 100}%"
								></div>
							</div>
						</div>
					{/if}
				</a>
			</div>

			<!-- Section 3: Eisenhower Matrix -->
			<div class="lg:col-span-1 sm:col-span-2 lg:col-auto">
				<a href="/myflow" class="block bg-gray-900 rounded-lg border border-gray-800 p-2.5 hover:border-gray-700 hover:shadow-md transition-all h-full">
					<div class="flex items-center justify-between mb-2">
						<h2 class="text-sm font-bold text-white flex items-center gap-1.5">
							<LayoutGrid size={13} class="text-purple-400" />
							Eisenhower Matrix
						</h2>
						<span class="text-[10px] text-gray-400">{matrixTasks.length} tasks</span>
					</div>

					<!-- Compact Matrix Grid -->
					<div class="grid grid-cols-4 gap-1.5 mb-2">
						<!-- Do First -->
						<div class="bg-red-950 rounded-md p-1.5 border border-red-900 flex flex-col items-center">
							<div class="w-1.5 h-1.5 rounded-full bg-red-500 mb-1"></div>
							<p class="text-lg font-bold text-white leading-none">{urgentImportant}</p>
							<p class="text-[8px] text-red-400 mt-0.5 text-center leading-tight">Do First</p>
						</div>

						<!-- Schedule -->
						<div class="bg-amber-950 rounded-md p-1.5 border border-amber-900 flex flex-col items-center">
							<div class="w-1.5 h-1.5 rounded-full bg-amber-500 mb-1"></div>
							<p class="text-lg font-bold text-white leading-none">{notUrgentImportant}</p>
							<p class="text-[8px] text-amber-400 mt-0.5 text-center leading-tight">Schedule</p>
						</div>

						<!-- Delegate -->
						<div class="bg-blue-950 rounded-md p-1.5 border border-blue-900 flex flex-col items-center">
							<div class="w-1.5 h-1.5 rounded-full bg-blue-500 mb-1"></div>
							<p class="text-lg font-bold text-white leading-none">{urgentNotImportant}</p>
							<p class="text-[8px] text-blue-400 mt-0.5 text-center leading-tight">Delegate</p>
						</div>

						<!-- Eliminate -->
						<div class="bg-gray-800 rounded-md p-1.5 border border-gray-700 flex flex-col items-center">
							<div class="w-1.5 h-1.5 rounded-full bg-gray-400 mb-1"></div>
							<p class="text-lg font-bold text-white leading-none">{notUrgentNotImportant}</p>
							<p class="text-[8px] text-gray-400 mt-0.5 text-center leading-tight">Eliminate</p>
						</div>
					</div>

					<!-- Completed -->
					{#if matrixCompleted > 0}
						<div class="pt-2 border-t border-gray-800">
							<div class="flex items-center justify-between bg-emerald-950 rounded-md px-2.5 py-1.5 border border-emerald-900">
								<div class="flex items-center gap-1.5">
									<CheckCircle size={11} class="text-emerald-400" />
									<span class="text-[9px] font-bold text-emerald-300 uppercase tracking-wide">Completed</span>
								</div>
								<span class="text-sm font-bold text-emerald-400">{matrixCompleted}</span>
							</div>
						</div>
					{/if}
				</a>
			</div>
	</div>
</div>