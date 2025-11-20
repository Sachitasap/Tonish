<script lang="ts">
	import "../app.css";
	import favicon from '$lib/assets/favicon.svg';
	import { page } from '$app/stores';
	import { LayoutGrid, Kanban, BookOpen, Menu, X } from 'lucide-svelte';
	import TonishLogo from '$lib/components/TonishLogo.svelte';

	let { children } = $props();
	let mobileNavOpen = $state(false);

	const navLinks = [
		{ href: '/', label: 'Dashboard', icon: LayoutGrid },
		{ href: '/myflow', label: 'MyFlow', icon: Kanban },
		{ href: '/myflowbook', label: 'MyFlowBook', icon: BookOpen }
	];

	$effect(() => {
		// close mobile drawer whenever the route changes
		$page.url.pathname;
		mobileNavOpen = false;
	});
</script>

<svelte:head>
	<link rel="icon" href="{favicon}" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
	{#if !$page.url.pathname.startsWith('/login') && !$page.url.pathname.startsWith('/register')}
		<nav class="bg-white shadow-sm border-b sticky top-0 z-30 touch-pan-y">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex items-center justify-between h-16">
					<a href="/" class="flex items-center gap-2 hover:opacity-80 transition">
						<TonishLogo size="md" />
					</a>
					<div class="hidden md:flex items-center space-x-6">
						{#each navLinks as link}
							<a
								href={link.href}
								class={`inline-flex items-center gap-2 px-3 py-2 rounded-full text-sm font-medium transition ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
									? 'bg-blue-50 text-blue-700'
									: 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}`}
							>
								<link.icon size={18} />
								<span>{link.label}</span>
							</a>
						{/each}
					</div>
					<div class="flex items-center gap-4">
						<span class="hidden md:inline text-sm text-gray-500">Tonish v1.0</span>
						<button
							type="button"
							class="md:hidden inline-flex items-center justify-center rounded-full p-2 text-gray-600 hover:bg-gray-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500"
							onclick={() => mobileNavOpen = !mobileNavOpen}
							aria-label="Toggle navigation"
							aria-expanded={mobileNavOpen}
						>
							{#if mobileNavOpen}
								<X size={24} />
							{:else}
								<Menu size={24} />
							{/if}
						</button>
					</div>
				</div>
				<div class={`md:hidden overflow-hidden transition-[max-height] duration-300 ${mobileNavOpen ? 'max-h-60' : 'max-h-0'}`}>
					<div class="pt-2 pb-4 space-y-2">
						{#each navLinks as link}
							<a
								href={link.href}
								class={`flex items-center gap-3 rounded-xl px-4 py-3 text-base font-medium ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
									? 'bg-blue-50 text-blue-700'
									: 'text-gray-600 hover:bg-gray-100'}`}
							>
								<link.icon size={20} />
								<span>{link.label}</span>
							</a>
						{/each}
					</div>
				</div>
			</div>
		</nav>
	{/if}

	<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 pb-28 md:pb-12">
		{@render children()}
	</main>

	{#if !$page.url.pathname.startsWith('/login') && !$page.url.pathname.startsWith('/register')}
		<div class="md:hidden fixed bottom-0 inset-x-0 bg-white border-t shadow-lg z-40">
			<nav class="flex justify-around items-stretch py-2">
				{#each navLinks as link}
					<a
						href={link.href}
						class={`flex flex-col items-center text-xs font-medium gap-1 px-3 ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
							? 'text-blue-600'
							: 'text-gray-500'}`}
					>
						<div class="text-lg" aria-hidden="true">
							<link.icon size={20} />
						</div>
						<span>{link.label}</span>
					</a>
				{/each}
			</nav>
		</div>
	{/if}
</div>
