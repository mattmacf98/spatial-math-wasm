/**
 * Global declarations for WASM functions
 */

declare global {
	
	function getPosesFromTrajectory(frameSystem: string, trajectory: string, frameName: string): unknown[];
	/**
	 * Go WASM runtime required for executing Go WebAssembly modules
	 */
	class Go {
		constructor();
		/**
		 * Import object for WebAssembly instantiation
		 */
		importObject: WebAssembly.Imports;
		/**
		 * Run the Go WASM instance
		 * @param instance - The WebAssembly instance to run
		 */
		run(instance: WebAssembly.Instance): Promise<void>;
	}
}

export {};

