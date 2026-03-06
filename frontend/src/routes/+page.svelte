<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { taskAPI } from '$lib/api';
	import { wsService } from '$lib/websocket';
	import { BarChart3, Zap, CheckCircle, Kanban, LayoutGrid, Calendar, AlarmClock, AlertTriangle, Clock, ChevronRight, Bell, CreditCard, Star, Circle } from 'lucide-svelte';

	type TaskStatus = 'todo' | 'in-progress' | 'done';
	type TaskPriority = 'low' | 'medium' | 'high';
	type CalendarSub = 'regular' | 'payment' | 'reminder' | 'event';

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
		task_type?: string | null;
		calendar_subtype?: CalendarSub | null;
		is_payment?: boolean;
		amount?: number;
		currency?: string;
		is_paid?: boolean;
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

	// Only kanban and matrix tasks feed into the Task Summary
	const flowTasks = $derived(
		tasks.filter((task) => task.task_type !== 'calendar')
	);

	// Derived state for task statistics
	const todayTasks = $derived(
		flowTasks.filter((task) => isToday(task.due_date) && task.status !== 'done')
	);
	const todayCompletedTasks = $derived(
		flowTasks.filter((task) => isToday(task.due_date) && task.status === 'done')
	);
	const upcomingTasks = $derived(
		flowTasks.filter((task) => isWithinNextDays(task.due_date, 7) && task.status !== 'done')
	);
	const overdueTasks = $derived(
		flowTasks.filter((task) => isOverdue(task.due_date, task.status))
	);
	const inProgressTasks = $derived(flowTasks.filter((task) => task.status === 'in-progress').length);
	const completedTasks = $derived(flowTasks.filter((task) => task.status === 'done').length);
	const totalTodayTasks = $derived(todayTasks.length + todayCompletedTasks.length);
	const todayAllDone = $derived(totalTodayTasks > 0 && todayTasks.length === 0);
	const todoTasks = $derived(flowTasks.filter((task) => task.status === 'todo').length);

	// Kanban board task breakdown
	const kanbanTasks = $derived(tasks.filter((task) => (task.task_type === 'kanban' || (!task.task_type && !task.quadrant))));
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

	// Calendar task breakdown (tasks created from the Calendar view)
	const calendarTasks = $derived(tasks.filter((task) => task.task_type === 'calendar'));
	const calendarOverdue = $derived(calendarTasks.filter((task) => isOverdue(task.due_date, task.status)));

	// Next 5 upcoming calendar tasks sorted by date
	const upcomingCalendarTasks = $derived(
		[...calendarTasks]
			.filter((task) => task.due_date && task.status !== 'done')
			.sort((a, b) => new Date(a.due_date!).getTime() - new Date(b.due_date!).getTime())
			.slice(0, 5)
	);

	function calendarSubIcon(task: Task): string {
		if (task.calendar_subtype === 'payment' || task.is_payment) return 'payment';
		if (task.calendar_subtype === 'reminder') return 'reminder';
		if (task.calendar_subtype === 'event') return 'event';
		return 'task';
	}

	function formatUpcomingDate(date: string): string {
		const d = new Date(date);
		const today = new Date();
		today.setHours(0, 0, 0, 0);
		d.setHours(0, 0, 0, 0);
		const diff = Math.round((d.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
		if (diff === 0) return 'Today';
		if (diff === 1) return 'Tomorrow';
		if (diff <= 7) return `${diff}d`;
		return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
	}
</script>

<div class="max-w-[1600px] mx-auto">
	<!-- Compact 4-Section Layout -->
		<div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-2 sm:gap-3">
			
			<!-- Section 1: Task Summary (compact list) -->
			<div>
				<div class="bg-gray-900 rounded-lg border border-gray-800 p-2 h-full">
					<h2 class="text-xs font-bold text-white mb-1.5 flex items-center gap-1.5 px-1">
						<BarChart3 size={12} class="text-blue-600" />
						Task Summary
					</h2>
					<div class="space-y-1">
						<!-- ── TODAY'S TASKS – hero row ── -->
						<a href="/myflow"
							class="group relative flex items-center gap-2 px-2 py-2 rounded-md border transition-all cursor-pointer
								{todayAllDone
									? 'bg-emerald-950/60 border-emerald-800 hover:border-emerald-600'
									: todayTasks.length > 0
										? 'bg-amber-950/70 border-amber-700 hover:border-amber-500'
										: 'bg-gray-800/60 border-gray-700 hover:border-gray-600'}">

							<!-- pulsing indicator when tasks remain -->
							<div class="relative flex-shrink-0">
								<div class="w-6 h-6 rounded-md flex items-center justify-center
									{todayAllDone ? 'bg-emerald-800' : todayTasks.length > 0 ? 'bg-amber-800' : 'bg-gray-700'}">
									{#if todayAllDone}
										<CheckCircle size={12} class="text-emerald-400" />
									{:else}
										<AlarmClock size={12} class="{todayTasks.length > 0 ? 'text-amber-400' : 'text-gray-400'}" />
									{/if}
								</div>
								{#if todayTasks.length > 0 && !todayAllDone}
									<span class="absolute -top-0.5 -right-0.5 w-2 h-2 bg-amber-400 rounded-full animate-ping opacity-75"></span>
									<span class="absolute -top-0.5 -right-0.5 w-2 h-2 bg-amber-400 rounded-full"></span>
								{/if}
							</div>

							<div class="flex-1 min-w-0">
								<p class="text-[11px] font-bold leading-none
									{todayAllDone ? 'text-emerald-300' : todayTasks.length > 0 ? 'text-amber-300' : 'text-gray-300'}">
									Today
								</p>
								{#if totalTodayTasks > 0}
									<!-- progress bar -->
									<div class="mt-1 h-1 rounded-full bg-gray-700 overflow-hidden">
										<div class="h-full rounded-full transition-all duration-500
											{todayAllDone ? 'bg-emerald-500' : 'bg-amber-500'}"
											style="width: {totalTodayTasks > 0 ? (todayCompletedTasks.length / totalTodayTasks) * 100 : 0}%">
										</div>
									</div>
									<p class="text-[9px] mt-0.5
										{todayAllDone ? 'text-emerald-400' : 'text-amber-400/80'}">
										{todayCompletedTasks.length} of {totalTodayTasks} done
									</p>
								{:else}
									<p class="text-[9px] text-gray-500 mt-0.5">No tasks due today</p>
								{/if}
							</div>

							<div class="text-right flex-shrink-0">
								{#if todayAllDone}
									<span class="text-sm font-bold text-emerald-400">✓</span>
								{:else}
									<span class="text-base font-bold
										{todayTasks.length > 0 ? 'text-amber-300' : 'text-white'}">
										{todayTasks.length}
									</span>
									{#if todayTasks.length > 0}
										<p class="text-[8px] text-amber-500 leading-none">left</p>
									{/if}
								{/if}
							</div>

							<ChevronRight size={11} class="flex-shrink-0 transition-colors
								{todayAllDone ? 'text-emerald-700 group-hover:text-emerald-400'
								: todayTasks.length > 0 ? 'text-amber-700 group-hover:text-amber-400'
								: 'text-gray-600 group-hover:text-gray-400'}" />
						</a>

						<!-- ── IN PROGRESS ── -->
						<a href="/myflow" class="group flex items-center gap-2 px-2 py-1.5 rounded-md hover:bg-gray-800 border border-transparent hover:border-gray-700 transition-all cursor-pointer">
							<div class="w-5 h-5 rounded bg-blue-950 flex items-center justify-center flex-shrink-0">
								<Zap size={10} class="text-blue-400" />
							</div>
							<span class="text-[10px] text-gray-400 flex-1 leading-none">In Progress</span>
							<span class="text-[9px] text-gray-600">{todoTasks} to-do</span>
							<span class="text-sm font-bold {inProgressTasks > 0 ? 'text-blue-300' : 'text-white'} w-5 text-right leading-none">{inProgressTasks}</span>
							<ChevronRight size={11} class="text-gray-600 group-hover:text-gray-400 flex-shrink-0 transition-colors" />
						</a>

						<!-- ── UPCOMING (7 days) ── -->
						<a href="/calendar" class="group flex items-center gap-2 px-2 py-1.5 rounded-md hover:bg-gray-800 border border-transparent hover:border-gray-700 transition-all cursor-pointer">
							<div class="w-5 h-5 rounded bg-cyan-950 flex items-center justify-center flex-shrink-0">
								<Calendar size={10} class="text-cyan-400" />
							</div>
							<span class="text-[10px] text-gray-400 flex-1 leading-none">Upcoming</span>
							<span class="text-[9px] text-gray-600">next 7 days</span>
							<span class="text-sm font-bold text-white w-5 text-right leading-none">{upcomingTasks.length}</span>
							<ChevronRight size={11} class="text-gray-600 group-hover:text-gray-400 flex-shrink-0 transition-colors" />
						</a>

						<!-- ── OVERDUE ── -->
						<a href="/myflow" class="group flex items-center gap-2 px-2 py-1.5 rounded-md transition-all cursor-pointer
							{overdueTasks.length > 0
								? 'hover:bg-red-950/40 border border-red-900/50 hover:border-red-800'
								: 'hover:bg-gray-800 border border-transparent hover:border-gray-700'}">
							<div class="w-5 h-5 rounded flex items-center justify-center flex-shrink-0
								{overdueTasks.length > 0 ? 'bg-red-950' : 'bg-gray-800'}">
								<AlertTriangle size={10} class="{overdueTasks.length > 0 ? 'text-red-400' : 'text-gray-500'}" />
							</div>
							<span class="text-[10px] flex-1 leading-none {overdueTasks.length > 0 ? 'text-red-300 font-semibold' : 'text-gray-400'}">Overdue</span>
							{#if overdueTasks.length > 0}
								<span class="text-[9px] text-red-500 font-bold uppercase tracking-wide">urgent</span>
							{/if}
							<span class="text-sm font-bold w-5 text-right leading-none {overdueTasks.length > 0 ? 'text-red-400' : 'text-white'}">{overdueTasks.length}</span>
							<ChevronRight size={11} class="flex-shrink-0 transition-colors {overdueTasks.length > 0 ? 'text-red-700 group-hover:text-red-400' : 'text-gray-600 group-hover:text-gray-400'}" />
						</a>

						<!-- ── COMPLETED TOTAL ── -->
						<div class="flex items-center gap-2 px-2 py-1 mt-0.5">
							<div class="flex-1 h-px bg-gray-800"></div>
							<span class="text-[9px] text-gray-600">{completedTasks} completed total</span>
							<div class="flex-1 h-px bg-gray-800"></div>
						</div>
					</div>
				</div>
			</div>

			<!-- Section 2: Calendar – Upcoming -->
			<div>
				<a href="/calendar" class="block bg-gray-900 rounded-lg border border-gray-800 p-2.5 hover:border-gray-700 hover:shadow-md transition-all h-full">
					<div class="flex items-center justify-between mb-2">
						<h2 class="text-sm font-bold text-white flex items-center gap-1.5">
							<Calendar size={13} class="text-cyan-400" />
							Upcoming
						</h2>
						{#if calendarOverdue.length > 0}
							<span class="flex items-center gap-1 text-[9px] font-bold text-red-400 bg-red-950 border border-red-900 rounded px-1.5 py-0.5">
								<AlertTriangle size={9} />{calendarOverdue.length} overdue
							</span>
						{:else}
							<span class="text-[10px] text-gray-500">{calendarTasks.length} tasks</span>
						{/if}
					</div>

					{#if upcomingCalendarTasks.length > 0}
						<div class="space-y-1">
							{#each upcomingCalendarTasks as task}
								{@const kind = calendarSubIcon(task)}
								{@const label = formatUpcomingDate(task.due_date!)}
								{@const isImminent = label === 'Today' || label === 'Tomorrow'}
								<div class="flex items-center gap-2 px-2 py-1.5 rounded-md bg-gray-800/60 border border-gray-700/60">
									<!-- type icon -->
									<div class="flex-shrink-0">
										{#if kind === 'payment'}
											<CreditCard size={11} class="text-amber-400" />
										{:else if kind === 'reminder'}
											<Bell size={11} class="text-purple-400" />
										{:else if kind === 'event'}
											<Star size={11} class="text-blue-400" />
										{:else}
											<Circle size={11} class="text-gray-400" />
										{/if}
									</div>
									<!-- title -->
									<span class="flex-1 text-[11px] text-gray-200 truncate leading-none">{task.title}</span>
									<!-- date badge -->
									<span class="flex-shrink-0 text-[9px] font-semibold px-1.5 py-0.5 rounded leading-none
										{isImminent ? 'bg-amber-900/60 text-amber-300' : 'bg-gray-700 text-gray-400'}">
										{label}
									</span>
								</div>
							{/each}
						</div>
					{:else}
						<div class="flex flex-col items-center justify-center py-6 gap-2">
							<Calendar size={22} class="text-gray-700" />
							<p class="text-[10px] text-gray-600">No upcoming tasks</p>
						</div>
					{/if}
				</a>
			</div>

			<!-- Section 3: Kanban Board Summary -->
			<div>
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

			<!-- Section 4: Eisenhower Matrix -->
			<div>
				<a href="/myflow?view=matrix" class="block bg-gray-900 rounded-lg border border-gray-800 p-2.5 hover:border-gray-700 hover:shadow-md transition-all h-full">
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