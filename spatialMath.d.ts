export interface OrientationVector {
    theta: number;
    ox: number;
    oy: number;
    oz: number;
}

/**
 * Converts a quaternion (w, x, y, z) to an OrientationVector.
 * Example: quatToOv(1, 0, 0, 0)
 */
export function quatToOv(w: number, x: number, y: number, z: number): OrientationVector;
