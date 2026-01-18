import { PlusIcon } from '@heroicons/react/20/solid';

type EmptyStateProps = {
  title: string;
  description: string;
  buttonText: string;
  buttonOnClick: () => void;
}

export default function EmptyState({
  title,
  description,
  buttonText,
  buttonOnClick
}: EmptyStateProps) {
  return (
    <div className="text-center">
      <svg
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        aria-hidden="true"
        className="mx-auto size-12 text-highlight"
      >
        <path
          d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"
          strokeWidth={2}
          vectorEffect="non-scaling-stroke"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
      <h3 className="mt-2 text-sm font-semibold text-base">{title}</h3>
      <p className="mt-1 text-sm text-highlight">{description}</p>
      <div className="mt-6">
        <button
          type="button"
          onClick={buttonOnClick}
          className="inline-flex items-center rounded-md bg-accent px-3 py-2 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          <PlusIcon aria-hidden="true" className="mr-1.5 -ml-0.5 size-5 text-inverted group-hover:text-inverted" />
          {buttonText}
        </button>
      </div>
    </div>
  )
}
