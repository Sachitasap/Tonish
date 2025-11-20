<script lang="ts">
	interface Props {
		size?: 'sm' | 'md' | 'lg';
		showVersion?: boolean;
	}

	let { size = 'md', showVersion = false }: Props = $props();

	const sizeClasses = {
		sm: 'h-8 w-8',
		md: 'h-10 w-10',
		lg: 'h-12 w-12'
	};

	const textSizeClasses = {
		sm: 'text-lg',
		md: 'text-2xl',
		lg: 'text-3xl'
	};

	const dotSizeClasses = {
		sm: 'w-1.5 h-1.5',
		md: 'w-2 h-2',
		lg: 'w-2.5 h-2.5'
	};
</script>

<div class="flex items-center gap-1">
	<!-- Tonish Text with Handwriting Style -->
	<div class="flex items-center gap-0.5">
		<span 
			class="{textSizeClasses[size]} font-bold tracking-tight"
			style="font-family: 'Brush Script MT', cursive; letter-spacing: -0.5px; background: linear-gradient(135deg, #3b82f6 0%, #1e40af 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;"
		>
			Tonish
		</span>
		
		<!-- Three Dots: Start → Continuous → Finish -->
		<div class="flex items-center gap-1 ml-1">
			<div class="{dotSizeClasses[size]} rounded-full bg-gradient-to-r from-green-400 to-green-500 shadow-sm"></div>
			<div class="{dotSizeClasses[size]} rounded-full bg-gradient-to-r from-blue-400 to-blue-500 shadow-sm animate-pulse"></div>
			<div class="{dotSizeClasses[size]} rounded-full bg-gradient-to-r from-purple-400 to-purple-500 shadow-sm"></div>
		</div>
	</div>

	<!-- Optional Version Badge -->
	{#if showVersion}
		<span class="text-xs font-medium text-gray-500 ml-2">v1.0</span>
	{/if}
</div>

<style>
	:global {
		@import url('https://fonts.googleapis.com/css2?family=Kalam:wght@700&display=swap');
	}

	span {
		font-family: 'Kalam', cursive !important;
		font-weight: 700;
	}
</style>
