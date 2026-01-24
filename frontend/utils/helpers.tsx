export function classNames(...classes: unknown[]) {
  return classes.filter(Boolean).join(' ')
}

export function formatBytes(bytes: number): string {
  // -1 indicates unlimited storage
  if (bytes === -1) return 'Unlimited';
  if (bytes === 0) return '0 MB';
  const k = 1024;
  const sizeUnits = new Map<number, string>([
    [0, 'Bytes'],
    [1, 'KB'],
    [2, 'MB'],
    [3, 'GB'],
    [4, 'TB'],
  ]);
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  const sizeIndex = Math.min(Math.max(i, 0), 4);
  const unit = sizeUnits.get(sizeIndex) ?? 'Bytes';
  return parseFloat((bytes / Math.pow(k, sizeIndex)).toFixed(2)) + ' ' + unit;
}