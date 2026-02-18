<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { taskAPI } from '$lib/api';
	import TonishLogo from '$lib/components/TonishLogo.svelte';
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
	
	onMount(async () => {
		// Check if user is authenticated
		const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
		if (!token) {
			goto('/login');
			return;
		}
		try {
			const tasksData = await taskAPI.getAll();
			tasks = tasksData;
		} catch (error: any) {
			console.error('Failed to load data:', error);
			// If unauthorized, redirect to login
			if (error.message?.includes('Authorization') || error.message?.includes('401')) {
				localStorage.removeItem('authToken');
				goto('/login');
			}
		} finally {
			loading = false;
		}
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

<div class="min-h-screen bg-gray-950">
	<div class="max-w-[1600px] mx-auto px-3 sm:px-4 lg:px-6">
		<!-- Header with Logo -->
		<div class="flex items-center justify-between py-6">
			<div class="flex items-center gap-3">
				<TonishLogo size="md" variant="icon" />
				<div>
					<h1 class="text-2xl font-bold text-white">Dashboard</h1>
					<p class="text-sm text-gray-400">Welcome back to Tonish</p>
				</div>
			</div>
		</div>
		<!-- Compact 3-Section Layout -->
		<div class="grid grid-cols-1 lg:grid-cols-3 gap-3 sm:gap-4">
			
			<!-- Section 1: Task Summary (2x2 Grid) -->
			<div class="lg:col-span-1 space-y-2">
				<div class="bg-gray-900 rounded-lg border border-gray-800 p-2.5">
					<h2 class="text-sm font-bold text-white mb-2 flex items-center gap-1.5">
						<BarChart3 size={14} class="text-blue-600" />
						Task Summary
					</h2>
					<div class="grid grid-cols-2 gap-1">
						<!-- Today's Tasks -->
						<a href="/myflow" class="flex flex-col justify-between p-2 rounded-md hover:bg-gray-800 transition-colors group">
							<div class="flex items-center gap-1.5 mb-1">
								<div class="w-6 h-6 rounded-md bg-amber-950 flex items-center justify-center flex-shrink-0">
									<CheckCircle size={12} class="text-amber-400" />
								</div>
								<p class="text-[10px] text-gray-400 leading-none">Today</p>
							</div>
							<div class="flex items-baseline justify-between">
								<p class="text-xl font-bold text-white leading-tight">{todayTasks.length}</p>
								{#if todayCompletedTasks.length > 0}
									<span class="text-[9px] text-gray-400">{todayCompletedTasks.length} done</span>
								{/if}
							</div>
						</a>

						<!-- Active Tasks -->
						<a href="/myflow" class="flex flex-col justify-between p-2 rounded-md hover:bg-gray-800 transition-colors group">
							<div class="flex items-center gap-1.5 mb-1">
								<div class="w-6 h-6 rounded-md bg-emerald-950 flex items-center justify-center flex-shrink-0">
									<Zap size={12} class="text-emerald-400" />
								</div>
								<p class="text-[10px] text-gray-400 leading-none">Active</p>
							</div>
							<div class="flex items-baseline justify-between">
								<p class="text-xl font-bold text-white leading-tight">{inProgressTasks}</p>
								<span class="text-[9px] text-gray-400">{completedTasks} done</span>
							</div>
						</a>

						<!-- Upcoming Tasks -->
						<a href="/calendar" class="flex flex-col justify-between p-2 rounded-md hover:bg-gray-800 transition-colors group">
							<div class="flex items-center gap-1.5 mb-1">
								<div class="w-6 h-6 rounded-md bg-blue-950 flex items-center justify-center flex-shrink-0">
									<Calendar size={12} class="text-blue-400" />
								</div>
								<p class="text-[10px] text-gray-400 leading-none">Upcoming</p>
							</div>
							<div class="flex items-baseline justify-between">
								<p class="text-xl font-bold text-white leading-tight">{upcomingTasks.length}</p>
								<span class="text-[9px] text-gray-400">7 days</span>
							</div>
						</a>

						<!-- Overdue Tasks -->
						<a href="/myflow" class="flex flex-col justify-between p-2 rounded-md hover:bg-gray-800 transition-colors group">
							<div class="flex items-center gap-1.5 mb-1">
								<div class="w-6 h-6 rounded-md bg-red-950 flex items-center justify-center flex-shrink-0">
									<Zap size={12} class="text-red-400" />
								</div>
								<p class="text-[10px] text-gray-400 leading-none">Overdue</p>
							</div>
							<div class="flex items-baseline justify-between">
								<p class="text-xl font-bold text-white leading-tight">{overdueTasks.length}</p>
								{#if overdueTasks.length > 0}
									<span class="text-[9px] text-red-400 font-medium">Urgent</span>
								{/if}
							</div>
						</a>
					</div>
				</div>
			</div>

			<!-- Section 2: Kanban Board Summary -->
			<div class="lg:col-span-1">
				<a href="/myflow" class="block bg-gray-900 rounded-lg border border-gray-800 p-3 sm:p-4 hover:border-gray-700 hover:shadow-md transition-all h-full">
					<div class="flex items-center justify-between mb-3">
						<h2 class="text-base sm:text-lg font-bold text-white flex items-center gap-2">
							<Kanban size={18} class="text-blue-400" />
							Kanban Board
						</h2>
						<span class="text-[10px] text-gray-400">{kanbanTasks.length} tasks</span>
					</div>

					<!-- Compact Kanban Columns -->
					<div class="grid grid-cols-3 gap-2 mb-3">
						<div class="bg-gray-800 rounded-lg p-2.5 border border-gray-700">
							<div class="flex items-center gap-1.5 mb-2">
								<div class="w-1.5 h-1.5 rounded-full bg-gray-400"></div>
								<h4 class="text-[10px] font-bold text-gray-300 uppercase tracking-wide">To Do</h4>
							</div>
							<p class="text-2xl font-bold text-white">{kanbanTodo}</p>
						</div>

						<div class="bg-blue-950 rounded-lg p-2.5 border border-blue-900">
							<div class="flex items-center gap-1.5 mb-2">
								<div class="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse"></div>
								<h4 class="text-[10px] font-bold text-blue-300 uppercase tracking-wide">Active</h4>
							</div>
							<p class="text-2xl font-bold text-white">{kanbanInProgress}</p>
						</div>

						<div class="bg-emerald-950 rounded-lg p-2.5 border border-emerald-900">
							<div class="flex items-center gap-1.5 mb-2">
								<div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
								<h4 class="text-[10px] font-bold text-emerald-300 uppercase tracking-wide">Done</h4>
							</div>
							<p class="text-2xl font-bold text-white">{kanbanDone}</p>
						</div>
					</div>

					<!-- Compact Progress Bar -->
					{#if kanbanTasks.length > 0}
						<div class="pt-3 border-t border-gray-800">
							<div class="flex justify-between text-xs text-gray-400 mb-2">
								<span>Progress</span>
								<span class="font-semibold text-white">{Math.round((kanbanDone / kanbanTasks.length) * 100)}%</span>
							</div>
							<div class="h-1.5 bg-gray-800 rounded-full overflow-hidden">
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
			<div class="lg:col-span-1">
				<a href="/myflow" class="block bg-gray-900 rounded-lg border border-gray-800 p-3 sm:p-4 hover:border-gray-700 hover:shadow-md transition-all h-full">
					<div class="flex items-center justify-between mb-3">
						<h2 class="text-base sm:text-lg font-bold text-white flex items-center gap-2">
							<LayoutGrid size={18} class="text-purple-400" />
							Eisenhower Matrix
						</h2>
						<span class="text-[10px] text-gray-400">{matrixTasks.length} tasks</span>
					</div>

					<!-- Compact Matrix Grid -->
					<div class="grid grid-cols-2 gap-2 mb-3">
						<!-- Do First -->
						<div class="bg-red-950 rounded-lg p-2.5 border border-red-900">
							<div class="flex items-center gap-1.5 mb-2">
								<div class="w-1.5 h-1.5 rounded-full bg-red-500"></div>
								<h4 class="text-[10px] font-bold text-red-300 uppercase tracking-wide">Do First</h4>
							</div>
							<p class="text-2xl font-bold text-white">{urgentImportant}</p>
							<p class="text-[9px] text-red-400 mt-1">Urgent & Important</p>
						</div>

						<!-- Schedule -->
						<div class="bg-amber-950 rounded-lg p-2.5 border border-amber-900">
							<div class="flex items-center gap-1.5 mb-2">
								<div class="w-1.5 h-1.5 rounded-full bg-amber-500"></div>
								<h4 class="text-[10px] font-bold text-amber-300 uppercase tracking-wide">Schedule</h4>
							</div>
							<p class="text-2xl font-bold text-white">{notUrgentImportant}</p>
							<p class="text-[9px] text-amber-400 mt-1">Important</p>
						</div>

						<!-- Delegate -->
						<div class="bg-blue-950 rounded-lg p-2.5 border border-blue-900">
							<div class="flex items-center gap-1.5 mb-2">
								<div class="w-1.5 h-1.5 rounded-full bg-blue-500"></div>
								<h4 class="text-[10px] font-bold text-blue-300 uppercase tracking-wide">Delegate</h4>
							</div>
							<p class="text-2xl font-bold text-white">{urgentNotImportant}</p>
							<p class="text-[9px] text-blue-400 mt-1">Urgent</p>
						</div>

						<!-- Eliminate -->
						<div class="bg-gray-800 rounded-lg p-2.5 border border-gray-700">
							<div class="flex items-center gap-1.5 mb-2">
								<div class="w-1.5 h-1.5 rounded-full bg-gray-400"></div>
								<h4 class="text-[10px] font-bold text-gray-300 uppercase tracking-wide">Eliminate</h4>
							</div>
							<p class="text-2xl font-bold text-white">{notUrgentNotImportant}</p>
							<p class="text-[9px] text-gray-400 mt-1">Low Priority</p>
						</div>
					</div>

					<!-- Completed -->
					{#if matrixCompleted > 0}
						<div class="pt-3 border-t border-gray-800">
							<div class="flex items-center justify-between bg-emerald-950 rounded-lg px-3 py-2 border border-emerald-900">
								<div class="flex items-center gap-2">
									<CheckCircle size={14} class="text-emerald-400" />
									<span class="text-[10px] font-bold text-emerald-300 uppercase tracking-wide">Completed</span>
								</div>
								<span class="text-lg font-bold text-emerald-400">{matrixCompleted}</span>
							</div>
						</div>
					{/if}
				</a>
			</div>
		</div>
	</div>
</div>