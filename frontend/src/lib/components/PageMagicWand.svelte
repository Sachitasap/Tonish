<script lang="ts">
	import { aiAPI, pageAPI } from '$lib/api';
	
	export let page: any;
	export let onPageUpdated: (updatedPage: any) => void;
	
	let activeFeature: 'summary' | 'tags' | 'enhance' | null = null;
	let isProcessing = false;
	let error = '';
	
	// AI Results
	let summaryData: any = null;
	let tagsData: any = null;
	let enhanceData: any = null;
	
	async function summarizePage() {
		isProcessing = true;
		error = '';
		activeFeature = 'summary';
		
		try {
			const result = await aiAPI.summarizePage(page.id);
			summaryData = result;
		} catch (e: any) {
			error = e.message || 'Failed to summarize page';
			console.error('Summarize error:', e);
		} finally {
			isProcessing = false;
		}
	}
	
	async function generateTags() {
		isProcessing = true;
		error = '';
		activeFeature = 'tags';
		
		try {
			const result = await aiAPI.generateTags(page.id);
			tagsData = result;
		} catch (e: any) {
			error = e.message || 'Failed to generate tags';
			console.error('Tags error:', e);
		} finally {
			isProcessing = false;
		}
	}
	
	async function enhancePage() {
		isProcessing = true;
		error = '';
		activeFeature = 'enhance';
		
		try {
			const result = await aiAPI.enhancePage(page.id);
			enhanceData = result;
		} catch (e: any) {
			error = e.message || 'Failed to enhance page';
			console.error('Enhance error:', e);
		} finally {
			isProcessing = false;
		}
	}
	
	async function applyTags() {
		if (!tagsData || !tagsData.result) return;
		
		try {
			// Convert existing tags string to array or start with empty array
			const existingTags = page.tags ? page.tags.split(',').map((t: string) => t.trim()).filter((t: string) => t) : [];
			// Merge and deduplicate
			const newTags = [...new Set([...existingTags, ...tagsData.result])];
			
			const updated = {
				...page,
				tags: newTags.join(', ')
			};
			
			await pageAPI.update(page.id, updated);
			onPageUpdated(updated);
			tagsData = null;
			activeFeature = null;
		} catch (e: any) {
			error = e.message || 'Failed to apply tags';
		}
	}
	
	function closePanel() {
		activeFeature = null;
		summaryData = null;
		tagsData = null;
		enhanceData = null;
		error = '';
	}
</script>

<div class="page-magic-wand">
	<div class="wand-header">
		<h3 class="wand-title">
			<span class="icon">📖✨</span>
			Page AI Assistant
		</h3>
		<p class="wand-subtitle">Powered by Qwen AI</p>
	</div>
	
	<div class="wand-actions">
		<button 
			class="wand-btn btn-summary" 
			on:click={summarizePage}
			disabled={isProcessing}
		>
			{#if isProcessing && activeFeature === 'summary'}
				<span class="spinner"></span>
				Analyzing...
			{:else}
				📝 Summarize Content
			{/if}
		</button>
		
		<button 
			class="wand-btn btn-tags" 
			on:click={generateTags}
			disabled={isProcessing}
		>
			{#if isProcessing && activeFeature === 'tags'}
				<span class="spinner"></span>
				Generating...
			{:else}
				🏷️ Generate Tags
			{/if}
		</button>
		
		<button 
			class="wand-btn btn-enhance" 
			on:click={enhancePage}
			disabled={isProcessing}
		>
			{#if isProcessing && activeFeature === 'enhance'}
				<span class="spinner"></span>
				Analyzing...
			{:else}
				💡 Enhance Page
			{/if}
		</button>
	</div>
	
	{#if error}
		<div class="error-box">
			⚠️ {error}
		</div>
	{/if}
	
	<!-- Summary Panel -->
	{#if summaryData}
		<div class="result-panel panel-summary">
			<div class="panel-header">
				<h4>📝 Content Summary</h4>
				<button class="btn-close" on:click={closePanel}>✕</button>
			</div>
			
			<div class="summary-content">
				{#if summaryData.result}
					<p class="summary-text">{summaryData.result}</p>
				{/if}
				
				{#if summaryData.tags && summaryData.tags.length > 0}
					<div class="suggested-tags">
						<h5>🏷️ Suggested Tags</h5>
						<div class="tag-list">
							{#each summaryData.tags as tag}
								<span class="tag">{tag}</span>
							{/each}
						</div>
					</div>
				{/if}
				
				{#if summaryData.priority_level}
					<div class="priority-badge">
						<span class="badge">Priority Level: {summaryData.priority_level}</span>
					</div>
				{/if}
			</div>
		</div>
	{/if}
	
	<!-- Tags Panel -->
	{#if tagsData}
		<div class="result-panel panel-tags">
			<div class="panel-header">
				<h4>🏷️ Suggested Tags</h4>
				<button class="btn-close" on:click={closePanel}>✕</button>
			</div>
			
			<div class="tags-content">
				{#if tagsData.result && tagsData.result.length > 0}
					<div class="tag-list">
						{#each tagsData.result as tag}
							<span class="tag">{tag}</span>
						{/each}
					</div>
					
					{#if tagsData.message}
						<div class="reasoning">
							<p>{tagsData.message}</p>
						</div>
					{/if}
					
					<button class="btn-apply" on:click={applyTags}>
						✓ Apply These Tags
					</button>
				{:else}
					<p>No tags suggested</p>
				{/if}
			</div>
		</div>
	{/if}
	
	<!-- Enhancement Panel -->
	{#if enhanceData}
		<div class="result-panel panel-enhance">
			<div class="panel-header">
				<h4>💡 Content Enhancement Suggestions</h4>
				<button class="btn-close" on:click={closePanel}>✕</button>
			</div>
			
			<div class="enhance-content">
				{#if enhanceData.message}
					<div class="suggestion-box">
						<h5>📋 Overview</h5>
						<p>{enhanceData.message}</p>
					</div>
				{/if}
				
				{#if enhanceData.result}
					{#if enhanceData.result.sectionsToAdd && enhanceData.result.sectionsToAdd.length > 0}
						<div class="sections-box">
							<h5>➕ Sections to Add</h5>
							<ul>
								{#each enhanceData.result.sectionsToAdd as section}
									<li>{section}</li>
								{/each}
							</ul>
						</div>
					{/if}
					
					{#if enhanceData.result.areasForDetail && enhanceData.result.areasForDetail.length > 0}
						<div class="details-box">
							<h5>🔍 Areas Needing More Detail</h5>
							<ul>
								{#each enhanceData.result.areasForDetail as area}
									<li>{area}</li>
								{/each}
							</ul>
						</div>
					{/if}
					
					{#if enhanceData.result.betterOrganization && enhanceData.result.betterOrganization.length > 0}
						<div class="organization-box">
							<h5>🎯 Organization Tips</h5>
							<ul>
								{#each enhanceData.result.betterOrganization as tip}
									<li>{tip}</li>
								{/each}
							</ul>
						</div>
					{/if}
					
					{#if enhanceData.result.relatedResources && enhanceData.result.relatedResources.length > 0}
						<div class="resources-box">
							<h5>📚 Related Resources</h5>
							<ul>
								{#each enhanceData.result.relatedResources as resource}
									<li>{resource}</li>
								{/each}
							</ul>
						</div>
					{/if}
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.page-magic-wand {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 12px;
		padding: 1.5rem;
		color: white;
		box-shadow: 0 8px 24px rgba(102, 126, 234, 0.25);
		margin-bottom: 1.5rem;
	}
	
	.wand-header {
		margin-bottom: 1rem;
	}
	
	.wand-title {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin: 0 0 0.25rem 0;
		font-size: 1.25rem;
		font-weight: 600;
	}
	
	.icon {
		font-size: 1.5rem;
		animation: float 3s ease-in-out infinite;
	}
	
	@keyframes float {
		0%, 100% { transform: translateY(0px); }
		50% { transform: translateY(-5px); }
	}
	
	.wand-subtitle {
		margin: 0;
		font-size: 0.875rem;
		opacity: 0.9;
	}
	
	.wand-actions {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
		gap: 0.75rem;
		margin-bottom: 1rem;
	}
	
	.wand-btn {
		padding: 0.75rem 1rem;
		background: rgba(255, 255, 255, 0.2);
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-radius: 8px;
		color: white;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.3s ease;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		font-size: 0.875rem;
	}
	
	.wand-btn:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.3);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	}
	
	.wand-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
	
	.spinner {
		width: 14px;
		height: 14px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
	
	.error-box {
		background: rgba(239, 68, 68, 0.2);
		border: 1px solid rgba(239, 68, 68, 0.5);
		border-radius: 6px;
		padding: 0.75rem;
		margin-top: 1rem;
		font-size: 0.875rem;
	}
	
	.result-panel {
		background: rgba(255, 255, 255, 0.98);
		color: #1f2937;
		border-radius: 8px;
		padding: 1.5rem;
		margin-top: 1rem;
		animation: slideIn 0.3s ease-out;
	}
	
	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
	
	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
		padding-bottom: 0.75rem;
		border-bottom: 2px solid #e5e7eb;
	}
	
	.panel-header h4 {
		margin: 0;
		color: #667eea;
		font-size: 1.125rem;
	}
	
	.btn-close {
		background: none;
		border: none;
		color: #9ca3af;
		font-size: 1.5rem;
		cursor: pointer;
		padding: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 4px;
		transition: all 0.2s ease;
	}
	
	.btn-close:hover {
		background: #f3f4f6;
		color: #374151;
	}
	
	.summary-content {
		line-height: 1.6;
	}
	
	.summary-text {
		background: #f0fdf4;
		padding: 1rem;
		border-radius: 6px;
		border-left: 4px solid #10b981;
		margin-bottom: 1rem;
		font-size: 0.9375rem;
	}
	
	.key-points, .action-items {
		margin: 1rem 0;
	}
	
	.key-points h5, .action-items h5 {
		margin: 0 0 0.5rem 0;
		color: #374151;
		font-size: 0.9375rem;
	}
	
	.key-points ul, .action-items ul {
		margin: 0;
		padding-left: 1.5rem;
	}
	
	.key-points li, .action-items li {
		margin: 0.5rem 0;
		color: #4b5563;
		font-size: 0.875rem;
	}
	
	.priority-badge {
		margin-top: 1rem;
	}
	
	.badge {
		display: inline-block;
		padding: 0.375rem 0.875rem;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
	}
	
	.badge-high {
		background: #fee2e2;
		color: #991b1b;
	}
	
	.badge-medium {
		background: #fed7aa;
		color: #9a3412;
	}
	
	.badge-low {
		background: #dbeafe;
		color: #1e40af;
	}
	
	.tags-content {
		line-height: 1.6;
	}
	
	.tag-list {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}
	
	.tag {
		background: #e0e7ff;
		color: #3730a3;
		padding: 0.375rem 0.875rem;
		border-radius: 12px;
		font-size: 0.875rem;
		font-weight: 500;
	}
	
	.reasoning {
		background: #fef3c7;
		padding: 1rem;
		border-radius: 6px;
		border-left: 4px solid #f59e0b;
		margin: 1rem 0;
		font-size: 0.875rem;
	}
	
	.reasoning p {
		margin: 0.5rem 0;
	}
	
	.btn-apply {
		background: #10b981;
		color: white;
		border: none;
		padding: 0.75rem 1.5rem;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		width: 100%;
		margin-top: 1rem;
	}
	
	.btn-apply:hover {
		background: #059669;
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
	}
	
	.enhance-content {
		line-height: 1.6;
	}
	
	.suggestion-box, .sections-box, .details-box, .organization-box, .resources-box {
		background: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		padding: 1rem;
		margin-bottom: 1rem;
	}
	
	.suggestion-box h5, .sections-box h5, .details-box h5, .organization-box h5, .resources-box h5 {
		margin: 0 0 0.75rem 0;
		color: #374151;
		font-size: 0.9375rem;
	}
	
	.suggestion-box p, .organization-box p {
		margin: 0;
		color: #4b5563;
		font-size: 0.875rem;
	}
	
	.sections-box ul, .details-box ul, .resources-box ul {
		margin: 0;
		padding-left: 1.5rem;
	}
	
	.sections-box li, .details-box li, .resources-box li {
		margin: 0.5rem 0;
		color: #4b5563;
		font-size: 0.875rem;
	}
</style>
