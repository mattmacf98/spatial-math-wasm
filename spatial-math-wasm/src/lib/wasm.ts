/**
 * WASM Module Initialization
 * 
 * This module handles loading and initializing the Go WASM module.
 * It ensures the module is only loaded once and provides a clean API.
 */

import type { OrientationVector } from './spatialMath';

// Global reference to track initialization
let wasmInitialized = false;
let wasmInitPromise: Promise<void> | null = null;

/**
 * Initialize the WASM module
 * @returns Promise that resolves when WASM is ready to use
 */
export async function initWasm(): Promise<void> {
	// Return existing promise if already initializing
	if (wasmInitPromise) {
		return wasmInitPromise;
	}

	// Return immediately if already initialized
	if (wasmInitialized) {
		return Promise.resolve();
	}

	// Create and store the initialization promise
	wasmInitPromise = (async () => {
		try {
			// Dynamically import the Go WASM runtime
			// @ts-expect-error - wasm_exec.js doesn't have TypeScript definitions
			await import('../wasm_exec.js');

			// Create Go instance and load WASM
			const go = new Go();
			const response = await fetch('/main.wasm');
			
			if (!response.ok) {
				throw new Error(`Failed to fetch WASM module: ${response.statusText}`);
			}

			const result = await WebAssembly.instantiateStreaming(response, go.importObject);
			
			// Run the Go program (non-blocking, runs in background)
			go.run(result.instance);
			
			wasmInitialized = true;
		} catch (error) {
			// Reset promise so initialization can be retried
			wasmInitPromise = null;
			throw new Error(`WASM initialization failed: ${error}`);
		}
	})();

	return wasmInitPromise;
}


export function isWasmReady(): boolean {
	return wasmInitialized;
}

// Wrapper functions go here (they should checck if WASM is ready and return undefined if not)
export function quatToOv(
	w: number,
	x: number,
	y: number,
	z: number
): OrientationVector | undefined {

    if (!isWasmReady()) {
        console.warn('WASM is not ready');
        return undefined;
    }

    // we have to be very careful about not passing things that may cause functions to panic through the WASM or else the WASM will exit
    if (w === null || x === null || y === null || z === null) {
        console.warn('w, x, y, z are not all numbers', w, x, y, z);
        return undefined;
    }
	
	if (typeof globalThis.quatToOv !== 'function') {
		console.error('quatToOv function not available in WASM module');
        return undefined;
	}
	
	return globalThis.quatToOv(w, x, y, z);
}

