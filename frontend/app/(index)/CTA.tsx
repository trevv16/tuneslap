import Link from "next/link"

export default function CTA() {
  return (
    <div className="relative -z-10 mt-32 px-6 lg:px-8">
      <div
        aria-hidden="true"
        className="absolute inset-x-0 top-1/2 -z-10 flex -translate-y-1/2 transform-gpu justify-center overflow-hidden blur-3xl sm:top-auto sm:right-[calc(50%-6rem)] sm:bottom-0 sm:translate-y-0 sm:transform-gpu sm:justify-end"
      >
        <div
          style={{
            clipPath:
              "polygon(73.6% 48.6%, 91.7% 88.5%, 100% 53.9%, 97.4% 18.1%, 92.5% 15.4%, 75.7% 36.3%, 55.3% 52.8%, 46.5% 50.9%, 45% 37.4%, 50.3% 13.1%, 21.3% 36.2%, 0.1% 0.1%, 5.4% 49.1%, 21.4% 36.4%, 58.9% 100%, 73.6% 48.6%)",
          }}
          className="aspect-1108/632 w-277 flex-none bg-linear-to-r from-[#F9B716] to-[#FF930F] opacity-25"
        />
      </div>
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-4xl font-semibold tracking-tight text-balance text-highlight sm:text-5xl">
          Ready to level up your sound game?
        </h2>
        <p className="mx-auto mt-6 max-w-xl text-lg/8 text-pretty text-base">
          Join thousands of creators who&apos;ve upgraded their workflow with TuneSlap. Fast, beautiful, and built for
          professionals.
        </p>
        <div className="mt-10 flex items-center justify-center gap-x-6">
          <Link
            href="/auth/signup"
            className="rounded-md bg-accent px-3.5 py-2.5 text-sm font-semibold text-inverted shadow-xs hover:bg-highlight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            Get started
          </Link>
          <Link href="#demo" className="text-sm/6 font-semibold text-base">
            Demo <span aria-hidden="true">→</span>
          </Link>
        </div>
      </div>
      <div
        aria-hidden="true"
        className="absolute top-full right-0 left-1/2 -z-10 hidden -translate-y-1/2 transform-gpu overflow-hidden blur-3xl sm:block"
      >
        <div
          style={{
            clipPath:
              "polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)",
          }}
          className="aspect-1155/678 w-288.75 bg-linear-to-tr from-[#F9B716] to-[#FF930F] opacity-30"
        />
      </div>
    </div>
  )
}
