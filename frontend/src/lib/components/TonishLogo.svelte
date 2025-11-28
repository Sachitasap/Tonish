<script lang="ts">
	export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
	export let showText: boolean = true;
	export let variant: 'full' | 'icon' = 'full';

	const logoSizes = {
		sm: { width: 80, height: 56 },
		md: { width: 110, height: 77 },
		lg: { width: 140, height: 98 },
		xl: { width: 180, height: 126 }
	};

	const iconSizes = {
		sm: 32,
		md: 40,
		lg: 48,
		xl: 56
	};
</script>

{#if variant === 'full' && showText}
	<div class="tonish-logo">
		<svg 
			width={logoSizes[size].width} 
			height={logoSizes[size].height} 
			viewBox="0 0 200 140" 
			fill="none" 
			xmlns="http://www.w3.org/2000/svg"
		>
			<!-- Handwritten "Tonish" text -->
			<g id="tonish-text">
				<path d="M20 35 Q25 25 35 30 T55 35 M30 35 L30 50 Q30 55 35 55 T45 50" stroke="#6366f1" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
				<path d="M60 40 Q65 35 70 40 Q75 50 70 55 Q65 60 60 55 Q55 50 60 40 Z" fill="#6366f1"/>
				<path d="M85 35 L85 55 M85 35 Q90 30 95 35 L95 55" stroke="#8b5cf6" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
				<path d="M105 35 L105 55 M105 42 Q110 38 115 42" stroke="#a855f7" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
				<path d="M125 40 Q128 35 133 38 Q136 42 133 46 Q130 50 125 48 Q122 46 125 40 Z" fill="#c026d3"/>
				<path d="M145 35 L145 55 M145 35 Q150 30 155 35 L155 50 Q155 55 160 55 T170 50" stroke="#d946ef" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
			</g>
			
			<!-- Double-ended arrow -->
			<g id="arrow">
				<path d="M30 80 L170 80" stroke="url(#arrowGradient)" stroke-width="3" stroke-linecap="round"/>
				<!-- Left arrow head -->
				<path d="M30 80 L40 75 M30 80 L40 85" stroke="url(#arrowGradient)" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
				<!-- Right arrow head -->
				<path d="M170 80 L160 75 M170 80 L160 85" stroke="url(#arrowGradient)" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
			</g>
			
			<!-- Three circles: Start, Progress, Finish -->
			<g id="status-circles">
				<!-- Start circle (blue) -->
				<circle cx="50" cy="110" r="12" fill="#3b82f6" opacity="0.9"/>
				<circle cx="50" cy="110" r="8" fill="#60a5fa" opacity="0.6"/>
				
				<!-- Progress circle (purple) -->
				<circle cx="100" cy="110" r="12" fill="#8b5cf6" opacity="0.9"/>
				<circle cx="100" cy="110" r="8" fill="#a78bfa" opacity="0.6"/>
				
				<!-- Finish circle (green) -->
				<circle cx="150" cy="110" r="12" fill="#10b981" opacity="0.9"/>
				<circle cx="150" cy="110" r="8" fill="#34d399" opacity="0.6"/>
			</g>
			
			<!-- Gradients -->
			<defs>
				<linearGradient id="arrowGradient" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" style="stop-color:#3b82f6;stop-opacity:1" />
					<stop offset="50%" style="stop-color:#8b5cf6;stop-opacity:1" />
					<stop offset="100%" style="stop-color:#10b981;stop-opacity:1" />
				</linearGradient>
			</defs>
		</svg>
	</div>
{:else}
	<!-- Icon variant (simplified) -->
	<svg 
		width={iconSizes[size]} 
		height={iconSizes[size]} 
		viewBox="0 0 60 60" 
		fill="none" 
		xmlns="http://www.w3.org/2000/svg"
	>
		<rect width="60" height="60" rx="12" fill="url(#bgGradient)"/>
		
		<!-- Stylized "T" -->
		<path d="M15 18 Q18 15 25 18 T35 18 M25 18 L25 42" stroke="white" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
		
		<!-- Three status dots -->
		<circle cx="20" cy="48" r="3" fill="#60a5fa"/>
		<circle cx="30" cy="48" r="3" fill="#a78bfa"/>
		<circle cx="40" cy="48" r="3" fill="#34d399"/>
		
		<defs>
			<linearGradient id="bgGradient" x1="0%" y1="0%" x2="100%" y2="100%">
				<stop offset="0%" style="stop-color:#6366f1;stop-opacity:1" />
				<stop offset="50%" style="stop-color:#8b5cf6;stop-opacity:1" />
				<stop offset="100%" style="stop-color:#a855f7;stop-opacity:1" />
			</linearGradient>
		</defs>
	</svg>
{/if}

<style>
	.tonish-logo {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		transition: opacity 0.2s ease;
	}
	
	.tonish-logo:hover {
		opacity: 0.85;
	}
	
	.tonish-logo svg {
		display: block;
	}
</style>
