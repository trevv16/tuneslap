const faqs = [
  {
    id: 1,
    question: "How fast is the audio playback?",
    answer:
      "TuneSlap is built for real-time performance with zero-lag audio triggering. Our modern web architecture ensures your sounds play instantly when you need them.",
  },
  {
    id: 2,
    question: "Can I use TuneSlap for live streaming?",
    answer:
      "TuneSlap is perfect for streamers, podcasters, and live content creators. Use hotkeys to trigger sounds without interrupting your flow.",
  },
  {
    id: 3,
    question: "Do I own my uploaded sounds?",
    answer:
      "Yes, you own all your content. You can export your sounds anytime and your data is securely stored in the cloud with full ownership rights.",
  },
  {
    id: 4,
    question: "Can I collaborate with my team?",
    answer:
      "Yes! TuneSlap includes team collaboration features. Share soundboards, set permissions, and work together on audio libraries.",
  },
  {
    id: 5,
    question: "Is TuneSlap open source?",
    answer:
      "Yes, TuneSlap is fully open source and self-hostable. You can deploy it on your own infrastructure and have complete control over your data. TuneSlap will remain open source and self-hostable forever.",
  },
]

export default function FAQ() {
  return (
    <div className="mx-auto max-w-2xl px-6 pb-8 sm:pt-12 sm:pb-24 lg:max-w-7xl lg:px-8 lg:pb-32">
      <h2 className="text-4xl font-semibold tracking-tight text-highlight sm:text-5xl">Frequently asked questions</h2>
      <dl className="mt-20 divide-y divide-gray-900/10">
        {faqs.map((faq) => (
          <div key={faq.id} className="py-8 first:pt-0 last:pb-0 lg:grid lg:grid-cols-12 lg:gap-8">
            <dt className="text-base/7 font-semibold text-highlight lg:col-span-5">{faq.question}</dt>
            <dd className="mt-4 lg:col-span-7 lg:mt-0">
              <p className="text-base/7 text-base">{faq.answer}</p>
            </dd>
          </div>
        ))}
      </dl>
    </div>
  )
}
