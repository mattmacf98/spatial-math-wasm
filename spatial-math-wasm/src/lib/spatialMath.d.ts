/**
 * TypeScript declarations for Spatial Math WASM functions
 * Generated from Go WASM module
 */

/**
 * Represents an orientation as a vector with rotation.
 * Corresponds to the Go OrientationVector type.
 */
export interface OrientationVector {
	/** Rotation angle in radians around the orientation vector */
	theta: number;
	/** X component of the orientation vector */
	ox: number;
	/** Y component of the orientation vector */
	oy: number;
	/** Z component of the orientation vector */
	oz: number;
}

/**
 * Global declarations for WASM functions
 */
declare global {
	/**
	 * Converts a quaternion to an Orientation Vector.
	 * 
	 * A quaternion represents a rotation in 3D space using four components (w, x, y, z).
	 * This function converts it to an orientation vector representation.
	 * 
	 * @param w - The real (scalar) component of the quaternion
	 * @param x - The i (first imaginary) component of the quaternion
	 * @param y - The j (second imaginary) component of the quaternion
	 * @param z - The k (third imaginary) component of the quaternion
	 * @returns An OrientationVector object containing theta, ox, oy, oz
	 * 
	 * @example
	 * // Identity quaternion (no rotation)
	 * const result = quatToOv(1, 0, 0, 0);
	 * console.log(result); // { theta: 0, ox: 0, oy: 0, oz: 1 }
	 * 
	 * @example
	 * // 90 degree rotation around X axis
	 * const result = quatToOv(0.707, 0.707, 0, 0);
	 */
	function quatToOv(w: number, x: number, y: number, z: number): OrientationVector;

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

