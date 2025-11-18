<script lang="ts">
	import { initWasm, getPosesFromTrajectory } from '$lib/wasm.js';
	import { onMount } from 'svelte';

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

	function testWasm() {
		const frameSystem = `{
			"frames": {
				"frame1": {
					"x": 0,
					"y": 0,
					"z": 0
				}
			}
		}`;
		const trajectory = {
			frame1: [0, 0, 0]
		};
		const frameName = "frame1";
		const poses = getPosesFromTrajectory(frameSystem, JSON.stringify(trajectory), frameName);
		console.log(poses);
	}

	// function convertQuatToOv() {
    //     const convertedOv = quatToOv(w, x, y, z);
    //     if (convertedOv) {
    //         ov = convertedOv;
    //     }
	// }
</script>

<button onclick={testWasm}>Get Poses</button>

