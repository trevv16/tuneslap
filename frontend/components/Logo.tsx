import Image from "next/image";

export function Logo() {
  return (
    <Image
      alt="TuneSlap"
      src="/logo.png"
      width={804}
      height={169}
      className="h-8 w-auto"
      priority
    />
  );
}
