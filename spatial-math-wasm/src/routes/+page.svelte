<script lang="ts">
	import type { OrientationVector } from '$lib/spatialMath.js';
	import { initWasm, quatToOv } from '$lib/wasm.js';
	import { onMount } from 'svelte';

	let w = $state(0);
	let x = $state(0);
	let y = $state(0);
	let z = $state(0);

	let ov = $state<OrientationVector>({theta: 0, ox: 0, oy: 0, oz: 0});
	let wasmReady = $state(false);
	let error = $state<string | null>(null);

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

	function convertQuatToOv() {
        const convertedOv = quatToOv(w, x, y, z);
        if (convertedOv) {
            ov = convertedOv;
        }
	}
</script>

<div class="max-w-md mx-auto mt-10 bg-white shadow-lg rounded-lg p-8">
	<h1 class="text-2xl font-semibold mb-6 text-center text-gray-800">Quaternion to Orientation Vector</h1>
	
	{#if error}
		<div class="mb-4 p-4 bg-red-50 border border-red-200 text-red-700 rounded">
			<strong>Error:</strong> {error}
		</div>
	{/if}
	
	{#if !wasmReady}
		<div class="mb-4 p-4 bg-blue-50 border border-blue-200 text-blue-700 rounded">
			Loading WASM module...
		</div>
	{/if}
	
	<div class="space-y-4">
		<div class="flex items-center space-x-3">
			<label for="w" class="w-12 font-medium text-gray-700">W:</label>
			<input
				id="w"
				type="number"
				bind:value={w}
				step="any"
				oninput={convertQuatToOv}
				class="flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
			/>
		</div>
		<div class="flex items-center space-x-3">
			<label for="x" class="w-12 font-medium text-gray-700">X:</label>
			<input
				id="x"
				type="number"
				bind:value={x}
				step="any"
				oninput={convertQuatToOv}
				class="flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
			/>
		</div>
		<div class="flex items-center space-x-3">
			<label for="y" class="w-12 font-medium text-gray-700">Y:</label>
			<input
				id="y"
				type="number"
				bind:value={y}
				step="any"
				oninput={convertQuatToOv}
				class="flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
			/>
		</div>
		<div class="flex items-center space-x-3">
			<label for="z" class="w-12 font-medium text-gray-700">Z:</label>
			<input
				id="z"
				type="number"
				bind:value={z}
				step="any"
				oninput={convertQuatToOv}
				class="flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400 transition"
			/>
		</div>
	</div>
	{#if ov}
		<div class="mt-8 bg-blue-50 p-6 rounded-lg">
			<h2 class="text-lg font-semibold mb-4 text-blue-700">Orientation Vector Result</h2>
			<div class="grid grid-cols-2 gap-y-2 gap-x-8 text-gray-700">
				<span class="font-medium">Theta:</span>
				<span>{ov.theta}</span>
				<span class="font-medium">OX:</span>
				<span>{ov.ox}</span>
				<span class="font-medium">OY:</span>
				<span>{ov.oy}</span>
				<span class="font-medium">OZ:</span>
				<span>{ov.oz}</span>
			</div>
		</div>
	{/if}
</div>

