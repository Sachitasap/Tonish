<script lang="ts">
	import "../app.css";
	import { page } from '$app/stores';
	import { LayoutGrid, Kanban, BookOpen, Archive, Calendar } from 'lucide-svelte';
	import { onMount, onDestroy } from 'svelte';
	import { wsService } from '$lib/websocket';

	let { children } = $props();
	let isKioskMode = $state(false);
	let idleTimer: number | null = null;

	const navLinks = [
		{ href: '/', label: 'Dashboard', icon: LayoutGrid },
		{ href: '/myflow', label: 'MyFlow', icon: Kanban },
		{ href: '/myflowbook', label: 'MyFlowBook', icon: BookOpen },
		{ href: '/calendar', label: 'Calendar', icon: Calendar },
		{ href: '/lookback', label: 'LookBack', icon: Archive }
	];

	$effect(() => {
		// Track route changes
		$page.url.pathname;
	});

	// Kiosk Mode functionality
	onMount(() => {
		// Initialize WebSocket connection
		wsService.connect();

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
			wsService.disconnect();
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

<div class="min-h-screen bg-gray-950">
	{#if !$page.url.pathname.startsWith('/login') && !$page.url.pathname.startsWith('/register')}
		<nav class="bg-gray-900 shadow-sm border-b border-gray-800 sticky top-0 z-30 touch-pan-y">
			<div class="app-shell">
				<div class="flex items-center justify-between h-12">
					<a href="/" class="flex items-center py-1 -ml-1 px-1.5 rounded-lg hover:bg-gray-800 transition-colors touch-manipulation">
						<img src="/tonish-logo.svg" alt="Tonish" class="h-7 w-auto" />
					</a>
					<div class="hidden md:flex items-center space-x-4">
						{#each navLinks as link}
							<a
								href={link.href}
								class={`touch-target inline-flex items-center gap-1.5 px-3 py-2 rounded-full text-sm font-medium transition touch-feedback ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
									? 'bg-blue-900 text-blue-300'
									: 'text-gray-400 hover:text-gray-200 hover:bg-gray-800'}`}
							>
								<link.icon size={16} />
								<span>{link.label}</span>
							</a>
						{/each}
					</div>

				</div>
			</div>
		</nav>
	{/if}

	<main class="app-shell py-3 pb-24 md:pb-10">
		{@render children()}
	</main>

	{#if !$page.url.pathname.startsWith('/login') && !$page.url.pathname.startsWith('/register')}
		<div class="md:hidden fixed bottom-0 inset-x-0 bg-gray-900 border-t border-gray-800 shadow-lg z-40 safe-bottom">
			<nav class="flex justify-around items-stretch py-1.5">
				{#each navLinks as link}
					<a
						href={link.href}
						class={`touch-target-large flex flex-col items-center justify-center text-xs font-medium gap-0.5 px-2 touch-feedback no-select ${$page.url.pathname === link.href || ($page.url.pathname.startsWith(link.href) && link.href !== '/')
							? 'text-blue-400'
							: 'text-gray-400'}`}
					>
						<div class="text-lg" aria-hidden="true">
							<link.icon size={20} />
						</div>
						<span class="text-[10px]">{link.label}</span>
					</a>
				{/each}
			</nav>
		</div>
	{/if}
</div>
