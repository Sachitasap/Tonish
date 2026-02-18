<script lang="ts">
	type Size = 'sm' | 'md' | 'lg' | 'xl';
	type Variant = 'full' | 'icon';

	interface Props {
		size?: Size;
		variant?: Variant;
		class?: string;
	}

	let { size = 'md', variant = 'full', class: className = '' }: Props = $props();

	const sizeClasses = {
		sm: 'h-8 w-auto',
		md: 'h-12 w-auto',
		lg: 'h-16 w-auto',
		xl: 'h-24 w-auto'
	};
</script>

<img 
	src="/tonish-logo.svg" 
	alt="Tonish Logo" 
	class="{sizeClasses[size]} {className}"
/>
