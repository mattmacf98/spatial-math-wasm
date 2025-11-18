<script lang="ts">
	import { initWasm, getPosesFromTrajectory } from '$lib/wasm.js';
	import { onMount } from 'svelte';

	let wasmReady = $state(false);
	let error = $state<string | null>(null);
	let frameSystemJson = $state<string | null>(null);
	let frameSystemFileName = $state<string | null>(null);
	let plansJson = $state<string | null>(null);
	let plansFileName = $state<string | null>(null);
	let frameSystemUploadError = $state<string | null>(null);
	let plansUploadError = $state<string | null>(null);

    // TODO: maybe I can wasm init into a hook?
	onMount(async () => {
		try {
			await initWasm();
			wasmReady = true;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to initialize WASM';
			console.error('WASM initialization error:', err);
		}
	});

	async function handleFrameSystemUpload(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		
		if (!file) {
			return;
		}

		// Reset errors
		frameSystemUploadError = null;

		// Check if file is JSON
		if (!file.name.endsWith('.json')) {
			frameSystemUploadError = 'Please upload a JSON file';
			frameSystemJson = null;
			frameSystemFileName = null;
			return;
		}

		try {
			const text = await file.text();
			// Validate JSON
			JSON.parse(text);
			frameSystemJson = text;
			frameSystemFileName = file.name;
		} catch (err) {
			frameSystemUploadError = err instanceof Error ? err.message : 'Invalid JSON file';
			frameSystemJson = null;
			frameSystemFileName = null;
		}
	}

	async function handleTrajectoryUpload(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		
		if (!file) {
			return;
		}

		// Reset errors
		plansUploadError = null;

		// Check if file is JSON
		if (!file.name.endsWith('.json')) {
			plansUploadError = 'Please upload a JSON file';
			plansJson = null;
			plansFileName = null;
			return;
		}

		try {
			const text = await file.text();
			// Validate JSON
			JSON.parse(text);
			plansJson = text;
			plansFileName = file.name;
		} catch (err) {
			plansUploadError = err instanceof Error ? err.message : 'Invalid JSON file';
			plansJson = null;
			plansFileName = null;
		}
	}

	function testWasm() {
		if (!frameSystemJson) {
			frameSystemUploadError = 'Please upload a frame system JSON first';
			return;
		}

		if (!plansJson) {
			plansUploadError = 'Please upload a trajectory JSON first';
			return;
		}

		const plans = JSON.parse(plansJson);
		const trajectories = plans.trajectories;

		console.log(`${trajectories.length} trajectories`);

		
		const trajectory = trajectories[0];
		console.log(trajectory);

		// DUMMY TRAJECTORY
		// const trajectory = [{
		// 	compliance: [0, 0, 0]
		// }];

		const frameName = "compliance";
		const poses = getPosesFromTrajectory(frameSystemJson, JSON.stringify(trajectory[0]), frameName);
		console.log(poses);
	}
</script>

<div class="container">
	<h1>Spatial Math WASM</h1>
	
	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if wasmReady}
		<div class="upload-section">
			<h2>Frame System</h2>
			<label for="frameSystemUpload" class="file-label">
				<input
					id="frameSystemUpload"
					type="file"
					accept=".json"
					onchange={handleFrameSystemUpload}
				/>
				<span class="file-button">Choose JSON File</span>
			</label>
			
			{#if frameSystemFileName}
				<div class="file-status success">
					✓ Loaded: {frameSystemFileName}
				</div>
			{/if}
			
			{#if frameSystemUploadError}
				<div class="file-status error">{frameSystemUploadError}</div>
			{/if}
		</div>

		<div class="upload-section">
			<h2>Trajectory</h2>
			<label for="trajectoryUpload" class="file-label">
				<input
					id="trajectoryUpload"
					type="file"
					accept=".json"
					onchange={handleTrajectoryUpload}
				/>
				<span class="file-button">Choose JSON File</span>
			</label>
			
			{#if plansFileName}
				<div class="file-status success">
					✓ Loaded: {plansFileName}
				</div>
			{/if}
			
			{#if plansUploadError}
				<div class="file-status error">{plansUploadError}</div>
			{/if}
		</div>

		<button onclick={testWasm} disabled={!frameSystemJson || !plansJson}>Get Poses</button>
	{:else}
		<p>Loading WASM...</p>
	{/if}
</div>

<style>
	.container {
		max-width: 800px;
		margin: 2rem auto;
		padding: 2rem;
		font-family: system-ui, -apple-system, sans-serif;
	}

	h1 {
		margin-bottom: 2rem;
	}

	h2 {
		font-size: 1.2rem;
		margin-bottom: 1rem;
	}

	.upload-section {
		margin-bottom: 2rem;
		padding: 1.5rem;
		border: 2px dashed #ccc;
		border-radius: 8px;
		background: #f9f9f9;
	}

	.file-label {
		display: inline-block;
		cursor: pointer;
	}

	.file-label input[type="file"] {
		display: none;
	}

	.file-button {
		display: inline-block;
		padding: 0.75rem 1.5rem;
		background: #4CAF50;
		color: white;
		border-radius: 4px;
		font-weight: 500;
		transition: background 0.2s;
	}

	.file-button:hover {
		background: #45a049;
	}

	.file-status {
		margin-top: 1rem;
		padding: 0.75rem;
		border-radius: 4px;
		font-size: 0.9rem;
	}

	.file-status.success {
		background: #d4edda;
		color: #155724;
		border: 1px solid #c3e6cb;
	}

	.file-status.error {
		background: #f8d7da;
		color: #721c24;
		border: 1px solid #f5c6cb;
	}

	.error {
		padding: 1rem;
		margin-bottom: 1rem;
		background: #f8d7da;
		color: #721c24;
		border: 1px solid #f5c6cb;
		border-radius: 4px;
	}

	button {
		padding: 0.75rem 2rem;
		font-size: 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-weight: 500;
		transition: background 0.2s;
	}

	button:hover:not(:disabled) {
		background: #0056b3;
	}

	button:disabled {
		background: #ccc;
		cursor: not-allowed;
	}
</style>

