import { ArrowPathIcon, CloudArrowUpIcon, FingerPrintIcon, LockClosedIcon } from "@heroicons/react/24/outline"

const features = [
  {
    name: "Lightning-Fast Audio",
    description:
      "Trigger sounds instantly with zero lag. Built for real-time performance with modern web technology and optimized audio processing.",
    icon: CloudArrowUpIcon,
  },
  {
    name: "Cloud Sync & Storage",
    description:
      "Your soundboards sync across all devices with secure S3 cloud storage. Access your sounds anywhere, anytime.",
    icon: LockClosedIcon,
  },
  {
    name: "Full Customization",
    description:
      "Optimize your audio with built-in editing tools and personalize your soundboards with custom themes, layouts, and hotkeys. Make it truly yours.",
    icon: ArrowPathIcon,
  },
  {
    name: "Team Collaboration",
    description:
      "Share boards with your team, set permissions, and collaborate on sound libraries. Perfect for podcasts and streaming teams.",
    icon: FingerPrintIcon,
  },
]

export default function Features() {
  return (
    <div id="features" className="mx-auto mt-24 max-w-7xl px-6 sm:mt-24 lg:px-8">
      <div className="mx-auto max-w-2xl lg:text-center">
        <h2 className="text-base/7 font-semibold text-accent">Built for Creators</h2>
        <p className="mt-2 text-4xl font-semibold tracking-tight text-pretty text-highlight sm:text-5xl lg:text-balance">
          Everything you need for professional soundboards
        </p>
        <p className="mt-6 text-lg/8 text-pretty text-base">
          Modern tools that keep up with your creativity. Organize, trigger, and share sounds effortlessly with zero-lag
          performance and beautiful design.
        </p>
      </div>
      <div className="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-4xl">
        <dl className="grid max-w-xl grid-cols-1 gap-x-8 gap-y-10 lg:max-w-none lg:grid-cols-2 lg:gap-y-16">
          {features.map((feature) => (
            <div key={feature.name} className="relative pl-16">
              <dt className="text-base/7 font-semibold text-highlight">
                <div className="absolute top-0 left-0 flex size-10 items-center justify-center rounded-lg bg-accent">
                  <feature.icon aria-hidden="true" className="size-6 text-highlight" />
                </div>
                {feature.name}
              </dt>
              <dd className="mt-2 text-base/7 text-base">{feature.description}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  )
}
