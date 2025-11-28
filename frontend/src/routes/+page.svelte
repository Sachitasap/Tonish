<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { taskAPI, notebookAPI } from '$lib/api';
	import { BarChart3, Zap, CheckCircle, BookMarked, Kanban } from 'lucide-svelte';

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

	type Notebook = {
		id: number;
		name: string;
		tags?: string | null;
		is_pinned?: boolean;
		pages?: { id: number }[];
	};

	let tasks = $state<Task[]>([]);
	let notebooks = $state<Notebook[]>([]);
	let loading = $state(true);
	
	onMount(async () => {
		// Check if user is authenticated
		const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
		if (!token) {
			goto('/login');
			return;
		}
		try {
			const [tasksData, notebooksData] = await Promise.all([
				taskAPI.getAll(),
				notebookAPI.getAll()
			]);
			tasks = tasksData;
			notebooks = notebooksData;
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
	
	const todayTasks = $derived(tasks.filter((task) => task.status === 'todo').slice(0, 5));
	const inProgressTasks = $derived(tasks.filter((task) => task.status === 'in-progress').length);
	const completedTasks = $derived(tasks.filter((task) => task.status === 'done').length);

	const quickLinks = [
		{
			title: 'Dive into MyFlow',
			description: 'Plan sprints, drag cards, and update priorities in one swipe.',
			href: '/myflow',
			className: 'from-blue-500 to-indigo-500',
			icon: Kanban
		},
		{
			title: 'Open MyFlowBook',
			description: 'Sketch ideas, meeting notes, and decision logs with rich text.',
			href: '/myflowbook',
			className: 'from-purple-500 to-pink-500',
			icon: BookMarked
		}
	];
</script>


	<div class="space-y-10">
		<section class="rounded-3xl bg-gradient-to-r from-blue-600 via-blue-500 to-purple-500 text-white p-6 md:p-10 flex flex-col gap-6 md:flex-row md:items-center md:justify-between shadow-xl">
			<div class="space-y-4">
				<p class="text-sm uppercase tracking-[0.2em] text-white/70">Today&apos;s Command Center</p>
				<h1 class="text-3xl md:text-4xl font-semibold leading-tight">Stay in sync across tasks & knowledge</h1>
				<p class="text-white/80 max-w-2xl">
					Swipe-friendly controls, keyboard-ready shortcuts, and layouts that flex between tablets and widescreen dashboards.
				</p>
				<div class="flex flex-wrap gap-3">
					<a
						href="/myflow"
						class="touch-target inline-flex items-center gap-2 rounded-full bg-white/90 text-blue-700 px-6 py-3.5 text-base font-semibold shadow hover:bg-white touch-feedback focus-ring"
					>
						<Kanban size={18} /> Open MyFlow
					</a>
					<a
						href="/myflowbook"
						class="touch-target inline-flex items-center gap-2 rounded-full border border-white/60 px-6 py-3.5 text-base font-semibold text-white hover:bg-white/10 touch-feedback focus-ring"
					>
						<BookMarked size={18} /> Open MyFlowBook
					</a>
				</div>
			</div>
			<div class="grid grid-cols-2 gap-4 w-full md:max-w-sm">
				<div class="rounded-2xl bg-white/15 backdrop-blur px-4 py-6 text-center">
					<p class="text-sm uppercase tracking-wide text-white/70">Active Tasks</p>
					<p class="text-4xl font-bold">{inProgressTasks}</p>
				</div>
				<div class="rounded-2xl bg-white/15 backdrop-blur px-4 py-6 text-center">
					<p class="text-sm uppercase tracking-wide text-white/70">Completed</p>
					<p class="text-4xl font-bold">{completedTasks}</p>
				</div>
				<div class="col-span-2 rounded-2xl bg-white/15 backdrop-blur px-4 py-6 text-center">
					<p class="text-sm uppercase tracking-wide text-white/70">Total Inventory</p>
					<p class="text-4xl font-bold">{tasks.length}</p>
				</div>
			</div>
		</section>
	
	{#if loading}
		<div class="text-center py-12">
			<p class="text-gray-500">Loading...</p>
		</div>
	{:else}
			<section class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
				<div class="bg-white rounded-2xl shadow-sm ring-1 ring-gray-100 p-5 flex items-center justify-between">
					<div>
						<p class="text-sm text-gray-500">Total Tasks</p>
						<p class="text-3xl font-semibold text-gray-900">{tasks.length}</p>
					</div>
					<BarChart3 size={40} class="text-gray-400" />
				</div>
				<div class="bg-white rounded-2xl shadow-sm ring-1 ring-gray-100 p-5 flex items-center justify-between">
					<div>
						<p class="text-sm text-gray-500">In Progress</p>
						<p class="text-3xl font-semibold text-blue-600">{inProgressTasks}</p>
					</div>
					<Zap size={40} class="text-blue-400" />
				</div>
				<div class="bg-white rounded-2xl shadow-sm ring-1 ring-gray-100 p-5 flex items-center justify-between">
					<div>
						<p class="text-sm text-gray-500">Completed</p>
						<p class="text-3xl font-semibold text-green-600">{completedTasks}</p>
					</div>
					<CheckCircle size={40} class="text-green-400" />
				</div>
			</section>

			<section class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<div class="bg-white rounded-2xl shadow-sm ring-1 ring-gray-100 p-6">
					<div class="flex items-center justify-between mb-4">
						<h2 class="text-xl font-semibold">Today&apos;s Tasks</h2>
						<a href="/myflow" class="text-blue-600 text-sm font-medium hover:underline">Open board →</a>
					</div>
					{#if todayTasks.length > 0}
						<ul class="space-y-3">
							{#each todayTasks as task}
								<li class="touch-target flex items-center gap-3 rounded-2xl border border-gray-100 px-4 py-3.5 touch-feedback">
									<input type="checkbox" class="h-5 w-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 touch-target" aria-label={`Mark ${task.title} complete`} />
									<div class="flex-1">
										<p class="text-base text-gray-800 font-medium">{task.title}</p>
										{#if task.priority === 'high'}
											<span class="inline-flex mt-1 text-xs items-center gap-1 bg-red-100 text-red-700 px-2 py-0.5 rounded-full">
												<Zap size={12} />
												High priority
											</span>
										{/if}
									</div>
								</li>
							{/each}
						</ul>
					{:else}
						<p class="text-gray-500 text-sm">No tasks for today</p>
					{/if}
				</div>

				<div class="bg-white rounded-2xl shadow-sm ring-1 ring-gray-100 p-6">
					<div class="flex items-center justify-between mb-4">
						<h2 class="text-xl font-semibold">Recent Notebooks</h2>
						<a href="/myflowbook" class="text-purple-600 text-sm font-medium hover:underline">Open library →</a>
					</div>
					{#if notebooks.length > 0}
						<ul class="space-y-3">
							{#each notebooks.slice(0, 5) as notebook}
								<li>
									<a href="/myflowbook/{notebook.id}" class="touch-target flex items-center justify-between rounded-2xl border border-purple-100 px-4 py-3.5 text-purple-800 hover:bg-purple-50 touch-feedback focus-ring">
										<span class="font-medium">{notebook.name}</span>
										<span aria-hidden="true">↗</span>
									</a>
								</li>
							{/each}
						</ul>
					{:else}
						<p class="text-gray-500 text-sm">No notebooks yet</p>
					{/if}
				</div>
			</section>

			<section class="grid grid-cols-1 md:grid-cols-2 gap-6">
				{#each quickLinks as card}
					<a
						href={card.href}
						class={`touch-target rounded-3xl bg-gradient-to-br ${card.className} text-white p-8 flex flex-col gap-3 shadow-lg hover:scale-[1.01] transition touch-feedback focus-ring`}
					>
						<card.icon size={40} class="text-white/90" />
						<h3 class="text-2xl font-bold">{card.title}</h3>
						<p class="text-white/80 text-base">{card.description}</p>
					</a>
				{/each}
			</section>
	{/if}
</div>
