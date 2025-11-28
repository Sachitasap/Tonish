<script lang="ts">
	import "../app.css";
	import favicon from '$lib/assets/favicon.svg';
	import { page } from '$app/stores';
	import { LayoutGrid, Kanban, BookOpen, Archive, Menu, X } from 'lucide-svelte';
	import TonishLogo from '$lib/components/TonishLogo.svelte';
	import { onMount } from 'svelte';

	let { children } = $props();
	let mobileNavOpen = $state(false);
	let isKioskMode = $state(false);
	let idleTimer: number | null = null;

	const navLinks = [
		{ href: '/', label: 'Dashboard', icon: LayoutGrid },
		{ href: '/myflow', label: 'MyFlow', icon: Kanban },
		{ href: '/myflowbook', label: 'MyFlowBook', icon: BookOpen },
		{ href: '/lookback', label: 'LookBack', icon: Archive }
	];

	$effect(() => {
		// close mobile drawer whenever the route changes
		$page.url.pathname;
		mobileNavOpen = false;
	});

	// Kiosk Mode functionality
	onMount(() => {
		// Check for kiosk mode query parameter or localStorage
		const urlParams = new URLSearchParams(window.location.search);
		const kioskParam = urlParams.get('kiosk');
		const kioskStored = localStorage.getItem('tonish-kiosk-mode');
		
		if (kioskParam === 'true' || kioskStored === 'true') {
			enableKioskMode();
		}

		// Listen for kiosk mode toggle (can be triggered by keyboard shortcut)
		const handleKeyDown = (e: KeyboardEvent) => {
			// Ctrl+Shift+K to toggle kiosk mode
			if (e.ctrlKey && e.shiftKey && e.key === 'K') {
				e.preventDefault();
				toggleKioskMode();
			}
		};
		window.addEventListener('keydown', handleKeyDown);
		
		// Prevent default back/forward navigation in kiosk mode
		window.addEventListener('popstate', handlePopState);

		// Reset idle timer on user activity
		const resetIdleTimer = () => {
			if (isKioskMode && idleTimer) {
				clearTimeout(idleTimer);
				idleTimer = window.setTimeout(() => {
					// Navigate to home after 5 minutes of inactivity
					window.location.href = '/';
				}, 5 * 60 * 1000);
			}
		};

		window.addEventListener('mousedown', resetIdleTimer);
		window.addEventListener('touchstart', resetIdleTimer);
		window.addEventListener('keypress', resetIdleTimer);

		return () => {
			window.removeEventListener('keydown', handleKeyDown);
			window.removeEventListener('popstate', handlePopState);
			window.removeEventListener('mousedown', resetIdleTimer);
			window.removeEventListener('touchstart', resetIdleTimer);
			window.removeEventListener('keypress', resetIdleTimer);
			if (idleTimer) clearTimeout(idleTimer);
		};
	});

	function enableKioskMode() {
		isKioskMode = true;
		localStorage.setItem('tonish-kiosk-mode', 'true');
		
		// Request fullscreen
		if (document.documentElement.requestFullscreen) {
			document.documentElement.requestFullscreen().catch((err) => {
				console.log('Fullscreen request failed:', err);
			});
		}

		// Lock orientation to landscape if available (for tablets)
		if (screen.orientation && (screen.orientation as any).lock) {
			(screen.orientation as any).lock('landscape').catch((err: any) => {
				console.log('Orientation lock failed:', err);
			});
		}
	}

	function disableKioskMode() {
		isKioskMode = false;
		localStorage.removeItem('tonish-kiosk-mode');
		
		// Exit fullscreen
		if (document.fullscreenElement && document.exitFullscreen) {
			document.exitFullscreen();
		}

		// Unlock orientation
		if (screen.orientation && (screen.orientation as any).unlock) {
			(screen.orientation as any).unlock();
		}

		if (idleTimer) {
			clearTimeout(idleTimer);
			idleTimer = null;
		}
	}

	function toggleKioskMode() {
		if (isKioskMode) {
			disableKioskMode();
		} else {
			enableKioskMode();
		}
	}

	function handlePopState(event: PopStateEvent) {
		if (isKioskMode) {
			event.preventDefault();
			// Stay on current page
			window.history.pushState(null, '', window.location.pathname);
		}
	}
</script>

<svelte:head>
	<link rel="icon" href="{favicon}" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
	{#if !$page.url.pathname.startsWith('/login') && !$page.url.pathname.startsWith('/register')}
		<nav class="bg-white shadow-sm border-b sticky top-0 z-30 touch-pan-y">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex items-center justify-between h-16">
					<a href="/" class="flex items-center gap-3 py-2 -ml-2 px-2 rounded-lg hover:bg-gray-50 transition-colors">
						<TonishLogo variant="icon" size="md" />
						<span class="hidden sm:inline-block text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">Tonish</span>
					</a>
					<div class="hidden md:flex items-center space-x-6">
						{#each navLinks as link}
							<a
								href={link.href}
								class={`touch-target inline-flex items-center gap-2 px-4 py-3 rounded-full text-sm font-medium transition touch-feedback ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
									? 'bg-blue-50 text-blue-700'
									: 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}`}
							>
								<link.icon size={18} />
								<span>{link.label}</span>
							</a>
						{/each}
					</div>
					<div class="flex items-center gap-4">
						<div class="hidden md:flex flex-col items-end">
							<span class="text-sm text-gray-500">Tonish v1.0</span>
							<span class="text-xs text-gray-400">Sailendra</span>
						</div>
						<button
							type="button"
							class="md:hidden touch-target-large inline-flex items-center justify-center rounded-full p-3 text-gray-600 hover:bg-gray-100 touch-feedback focus-ring"
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
								class={`touch-target flex items-center gap-3 rounded-xl px-4 py-3.5 text-base font-medium touch-feedback no-select ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
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
		<div class="md:hidden fixed bottom-0 inset-x-0 bg-white border-t shadow-lg z-40 safe-bottom">
			<nav class="flex justify-around items-stretch py-2">
				{#each navLinks as link}
					<a
						href={link.href}
						class={`touch-target-large flex flex-col items-center justify-center text-xs font-medium gap-1 px-3 touch-feedback no-select ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
							? 'text-blue-600'
							: 'text-gray-500'}`}
					>
						<div class="text-lg" aria-hidden="true">
							<link.icon size={22} />
						</div>
						<span class="text-[10px]">{link.label}</span>
					</a>
				{/each}
			</nav>
		</div>
	{/if}
</div>
