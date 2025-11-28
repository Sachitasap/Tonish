<script lang="ts">
	import { aiAPI, taskAPI } from '$lib/api';
	
	export let task: any;
	export let onTaskUpdated: (updatedTask: any) => void;
	
	let isEnhancing = false;
	let isBreakingDown = false;
	let showEnhancedPreview = false;
	let enhancedData: any = null;
	let subtasks: any[] = [];
	let error = '';
	
	async function enhanceTask() {
		isEnhancing = true;
		error = '';
		
		try {
			const result = await aiAPI.enhanceTask(task);
			enhancedData = result.enhanced;
			showEnhancedPreview = true;
		} catch (e: any) {
			error = e.message || 'AI enhancement failed';
			console.error('Enhancement error:', e);
		} finally {
			isEnhancing = false;
		}
	}
	
	async function applyEnhancement() {
		if (!enhancedData) return;
		
		try {
			// Update task with enhanced data
			const updated = {
				...task,
				title: enhancedData.title || task.title,
				description: enhancedData.description || task.description,
				priority: enhancedData.priority || task.priority,
				quadrant: enhancedData.quadrant || task.quadrant,
				tags: JSON.stringify(enhancedData.tags || [])
			};
			
			await taskAPI.update(task.id, updated);
			onTaskUpdated(updated);
			showEnhancedPreview = false;
		} catch (e: any) {
			error = e.message || 'Failed to apply enhancement';
		}
	}
	
	async function breakdownTask() {
		isBreakingDown = true;
		error = '';
		
		try {
			const result = await aiAPI.breakdownTask(task);
			subtasks = result.subtasks || [];
		} catch (e: any) {
			error = e.message || 'Task breakdown failed';
			console.error('Breakdown error:', e);
		} finally {
			isBreakingDown = false;
		}
	}
	
	async function createSubtask(subtask: any) {
		try {
			await taskAPI.create({
				...subtask,
				tags: JSON.stringify(subtask.tags || [])
			});
			// Remove from list after creation
			subtasks = subtasks.filter(s => s !== subtask);
		} catch (e: any) {
			error = e.message || 'Failed to create subtask';
		}
	}
</script>

<div class="magic-wand-panel">
	<h3 class="panel-title">
		<span class="icon">✨</span>
		Magic Wand AI
	</h3>
	
	<div class="actions">
		<button 
			class="btn-magic" 
			on:click={enhanceTask}
			disabled={isEnhancing}
		>
			{#if isEnhancing}
				<span class="spinner"></span>
				Enhancing...
			{:else}
				🎯 Enhance Task
			{/if}
		</button>
		
		<button 
			class="btn-magic" 
			on:click={breakdownTask}
			disabled={isBreakingDown}
		>
			{#if isBreakingDown}
				<span class="spinner"></span>
				Analyzing...
			{:else}
				📋 Break Down Task
			{/if}
		</button>
	</div>
	
	{#if error}
		<div class="error-message">
			⚠️ {error}
		</div>
	{/if}
	
	{#if showEnhancedPreview && enhancedData}
		<div class="preview-panel">
			<h4>✨ AI Enhancement Preview</h4>
			
			<div class="preview-field">
				<label>Title:</label>
				<div class="value">{enhancedData.title}</div>
			</div>
			
			<div class="preview-field">
				<label>Description:</label>
				<div class="value">{enhancedData.description}</div>
			</div>
			
			<div class="preview-field">
				<label>Priority:</label>
				<span class="badge badge-{enhancedData.priority}">
					{enhancedData.priority}
				</span>
			</div>
			
			<div class="preview-field">
				<label>Quadrant:</label>
				<div class="value">{enhancedData.quadrant}</div>
			</div>
			
			{#if enhancedData.tags && enhancedData.tags.length > 0}
				<div class="preview-field">
					<label>Tags:</label>
					<div class="tags">
						{#each enhancedData.tags as tag}
							<span class="tag">{tag}</span>
						{/each}
					</div>
				</div>
			{/if}
			
			{#if enhancedData.suggestions}
				<div class="preview-field">
					<label>💡 Suggestions:</label>
					<div class="value suggestions">{enhancedData.suggestions}</div>
				</div>
			{/if}
			
			<div class="preview-actions">
				<button class="btn-apply" on:click={applyEnhancement}>
					✓ Apply Enhancement
				</button>
				<button class="btn-cancel" on:click={() => showEnhancedPreview = false}>
					✗ Cancel
				</button>
			</div>
		</div>
	{/if}
	
	{#if subtasks.length > 0}
		<div class="subtasks-panel">
			<h4>📋 Suggested Subtasks</h4>
			
			{#each subtasks as subtask, i}
				<div class="subtask-card">
					<div class="subtask-header">
						<strong>{i + 1}. {subtask.title}</strong>
						<span class="badge badge-{subtask.priority}">
							{subtask.priority}
						</span>
					</div>
					<p class="subtask-description">{subtask.description}</p>
					<button 
						class="btn-create-subtask"
						on:click={() => createSubtask(subtask)}
					>
						+ Create Task
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.magic-wand-panel {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 12px;
		padding: 1.5rem;
		color: white;
		box-shadow: 0 8px 24px rgba(102, 126, 234, 0.25);
	}
	
	.panel-title {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin-bottom: 1rem;
		font-size: 1.25rem;
		font-weight: 600;
	}
	
	.icon {
		font-size: 1.5rem;
		animation: sparkle 2s infinite;
	}
	
	@keyframes sparkle {
		0%, 100% { transform: scale(1) rotate(0deg); }
		50% { transform: scale(1.2) rotate(180deg); }
	}
	
	.actions {
		display: flex;
		gap: 1rem;
		margin-bottom: 1rem;
	}
	
	.btn-magic {
		flex: 1;
		padding: 0.75rem 1.5rem;
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
	}
	
	.btn-magic:hover:not(:disabled) {
		background: rgba(255, 255, 255, 0.3);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	}
	
	.btn-magic:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}
	
	.spinner {
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
	
	.error-message {
		background: rgba(239, 68, 68, 0.2);
		border: 1px solid rgba(239, 68, 68, 0.5);
		border-radius: 6px;
		padding: 0.75rem;
		margin-top: 1rem;
		font-size: 0.875rem;
	}
	
	.preview-panel, .subtasks-panel {
		background: rgba(255, 255, 255, 0.95);
		color: #1f2937;
		border-radius: 8px;
		padding: 1.5rem;
		margin-top: 1rem;
	}
	
	.preview-panel h4, .subtasks-panel h4 {
		margin-top: 0;
		margin-bottom: 1rem;
		color: #667eea;
	}
	
	.preview-field {
		margin-bottom: 1rem;
	}
	
	.preview-field label {
		display: block;
		font-weight: 600;
		color: #4b5563;
		margin-bottom: 0.25rem;
		font-size: 0.875rem;
	}
	
	.preview-field .value {
		color: #1f2937;
		line-height: 1.5;
	}
	
	.value.suggestions {
		background: #fef3c7;
		padding: 0.75rem;
		border-radius: 6px;
		border-left: 4px solid #f59e0b;
		font-style: italic;
	}
	
	.badge {
		display: inline-block;
		padding: 0.25rem 0.75rem;
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
	
	.tags {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}
	
	.tag {
		background: #e0e7ff;
		color: #3730a3;
		padding: 0.25rem 0.75rem;
		border-radius: 12px;
		font-size: 0.75rem;
		font-weight: 500;
	}
	
	.preview-actions {
		display: flex;
		gap: 1rem;
		margin-top: 1.5rem;
	}
	
	.btn-apply, .btn-cancel {
		flex: 1;
		padding: 0.75rem;
		border: none;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
	}
	
	.btn-apply {
		background: #10b981;
		color: white;
	}
	
	.btn-apply:hover {
		background: #059669;
	}
	
	.btn-cancel {
		background: #e5e7eb;
		color: #374151;
	}
	
	.btn-cancel:hover {
		background: #d1d5db;
	}
	
	.subtask-card {
		background: #f9fafb;
		border: 1px solid #e5e7eb;
		border-radius: 6px;
		padding: 1rem;
		margin-bottom: 1rem;
	}
	
	.subtask-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.5rem;
	}
	
	.subtask-description {
		color: #6b7280;
		margin: 0.5rem 0;
		font-size: 0.875rem;
	}
	
	.btn-create-subtask {
		background: #667eea;
		color: white;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 6px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		margin-top: 0.5rem;
	}
	
	.btn-create-subtask:hover {
		background: #5568d3;
		transform: translateY(-1px);
	}
</style>
