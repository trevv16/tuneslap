import { ExclamationTriangleIcon } from '@heroicons/react/20/solid';

type DemoBannerProps = {
  message: string;
}

const isDemoMode = process.env.NEXT_PUBLIC_DEMO_MODE === 'true';

export default function DemoBanner({ message }: DemoBannerProps) {
  if (!isDemoMode) return null;

  return (
    <div className="rounded-md bg-warning/10 border border-warning p-4 mb-4">
      <div className="flex">
        <div className="shrink-0">
          <ExclamationTriangleIcon aria-hidden="true" className="size-5 text-warning" />
        </div>
        <div className="ml-3">
          <p className="text-sm text-warning">{message}</p>
        </div>
      </div>
    </div>
  );
}
