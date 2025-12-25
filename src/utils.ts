export function getElement<T extends HTMLElement>(
  selector: string,
  parent?: DocumentFragment,
) {
  const el = (parent ? parent : document).querySelector<T>(selector);
  if (!el || !(el instanceof HTMLElement)) {
    throw new Error(`failed to query ${selector} selector`);
  }
  return el;
}

const SIZES = ["k", "M", "G", "T", "P", "E"] as const;
export type ReadableSize = `${string} ${(typeof SIZES)[number] | ""}B`;

export function toReadableSize(size: number): ReadableSize {
  const unit = 1000;

  if (size < unit) {
    return `${size} B`;
  }

  let exponent = 0;
  let division = unit;
  for (let i = size / unit; i >= unit; i /= unit) {
    division *= unit;
    exponent++;
  }

  return `${(size / division).toFixed(2)} ${SIZES[exponent]}B`;
}

export type FileProgress = {
  bytes?: number;
  time?: number;
  percentage: number;
  samples?: number[];
};

export type ReadableSpeed = `${ReadableSize}/s`;

export function calculateProgress(
  progress: FileProgress,
  current: number,
  total: number,
) {
  const now = Date.now();
  const timeDiff = !progress.time ? 0 : (now - progress.time) / 1000;
  const byteDiff = current - (progress.bytes || 0);
  const samples = progress.samples || [];
  const currentSpeed = byteDiff / timeDiff;
  if (progress.time) samples.push(currentSpeed);
  if (samples.length > 5) samples.shift();
  const avgSpeed = samples.reduce((t, c) => t + c, 0) / samples.length || 0;
  const speed: ReadableSpeed = `${toReadableSize(avgSpeed)}/s`;

  const newProgress: FileProgress = {
    samples,
    time: now,
    bytes: current,
    percentage: Math.round((current / total) * 100),
  };

  return { speed, newProgress };
}
