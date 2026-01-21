import Image from "next/image";

type Feature = {
  title: string;
  subtitle: string;
  description: string;
  imageSrc: string;
  imageAlt: string;
};

const features: Feature[] = [
  {
    title: "Media Processing",
    subtitle: "Built-in Optimization",
    description:
      "Upload your audio and images, and TuneSlap automatically processes them for the best quality and smallest file size. Audio files are normalized and converted to optimal formats. Images are compressed and resized without losing clarity.",
    imageSrc: "/defaultKey.png", // Placeholder - replace with actual screenshot
    imageAlt: "Media processing interface showing file optimization",
  },
  {
    title: "Collaboration",
    subtitle: "Work Together",
    description:
      "Share your soundboards with teammates and set permissions for who can view or edit. Perfect for podcast teams, streaming groups, or anyone who needs to coordinate sounds across multiple people.",
    imageSrc: "/defaultBoard.png", // Placeholder - replace with actual screenshot
    imageAlt: "Collaboration interface showing shared board with team members",
  },
  {
    title: "Self-Hosted",
    subtitle: "Own Your Data",
    description:
      "Run TuneSlap on your own infrastructure with Docker. Your sounds, your servers, your rules. No vendor lock-in, no unexpected costs, and complete control over your data and privacy.",
    imageSrc: "/logo.png", // Placeholder - replace with actual screenshot
    imageAlt: "Docker deployment showing self-hosted setup",
  },
];

function FeatureRow({
  feature,
  imageOnLeft,
}: {
  feature: Feature;
  imageOnLeft: boolean;
}) {
  const imageBlock = (
    <div className="relative aspect-video overflow-hidden rounded-xl bg-card shadow-lg ring-1 ring-border">
      <Image
        src={feature.imageSrc}
        alt={feature.imageAlt}
        fill
        className="object-cover"
      />
    </div>
  );

  const textBlock = (
    <div className="flex flex-col justify-center">
      <p className="text-base/7 font-semibold text-primary">{feature.subtitle}</p>
      <h3 className="mt-2 text-2xl font-semibold tracking-tight text-highlight sm:text-3xl">
        {feature.title}
      </h3>
      <p className="mt-4 text-base/7 text-muted-foreground">
        {feature.description}
      </p>
    </div>
  );

  return (
    <div className="grid grid-cols-1 gap-8 lg:grid-cols-2 lg:gap-16 items-center">
      {imageOnLeft ? (
        <>
          <div className="lg:order-1">{imageBlock}</div>
          <div className="lg:order-2">{textBlock}</div>
        </>
      ) : (
        <>
          <div className="lg:order-2">{imageBlock}</div>
          <div className="lg:order-1">{textBlock}</div>
        </>
      )}
    </div>
  );
}

export default function Features() {
  return (
    <section
      id="features"
      className="mx-auto max-w-7xl px-6 lg:px-8 mt-24 sm:mt-32"
    >
      <div className="mx-auto max-w-2xl lg:text-center mb-16">
        <h2 className="text-base/7 font-semibold text-primary">Built for Creators</h2>
        <p className="mt-2 text-4xl font-semibold tracking-tight text-pretty text-highlight sm:text-5xl lg:text-balance">
          Everything you need for professional soundboards
        </p>
        <p className="mt-6 text-lg/8 text-pretty text-muted-foreground">
          Modern tools that keep up with your creativity. Organize, trigger, and
          share sounds effortlessly with zero-lag performance and beautiful
          design.
        </p>
      </div>

      <div className="space-y-24 lg:space-y-32">
        {features.map((feature, index) => (
          <FeatureRow
            key={feature.title}
            feature={feature}
            imageOnLeft={index % 2 === 0}
          />
        ))}
      </div>
    </section>
  );
}
