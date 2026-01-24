export function classNames(...classes: unknown[]) {
  return classes.filter(Boolean).join(' ')
}

export function formatBytes(bytes: number): string {
  // -1 indicates unlimited storage
  if (bytes === -1) return 'Unlimited';
  if (bytes === 0) return '0 MB';
  const k = 1024;
  const sizes: readonly string[] = ['Bytes', 'KB', 'MB', 'GB', 'TB'] as const;
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const sizeIndex = Math.min(i, sizes.length - 1);
  return parseFloat((bytes / Math.pow(k, sizeIndex)).toFixed(2)) + ' ' + sizes[sizeIndex];
}