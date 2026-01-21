import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Providers from "./providers";

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: {
    default: "TuneSlap - The Soundboard Platform",
    template: "%s | TuneSlap"
  },
  description: "Create, collaborate, and share your soundboards with TuneSlap. Upload audio clips, manage soundboards, and connect with other audio professionals in our creative community.",
  generator: "TuneSlap",
  alternates: {
    canonical: "https://tuneslap.com",
  },
  openGraph: {
    title: "TuneSlap - The Soundboard Platform",
    description: "Create, collaborate, and share your soundboards with TuneSlap. Upload audio clips, manage soundboards, and connect with other audio professionals in our creative community.",
    url: "https://tuneslap.com",
    siteName: "TuneSlap",
    images: [
      {
        url: "https://tuneslap.com/opengraph-image.png",
        width: 1200,
        height: 630,
        alt: "TuneSlap - The Soundboard Platform"
      }
    ],
    locale: "en_US",
    type: "website"
  },
  twitter: {
    card: "summary_large_image",
    title: "TuneSlap - The Soundboard Platform",
    description: "Create, collaborate, and share your soundboards with TuneSlap. Upload audio clips, manage soundboards, and connect with other audio professionals in our creative community.",
    creator: "@tuneslap",
    site: "@tuneslap",
    images: ["https://tuneslap.com/opengraph-image.png"]
  },
  icons: {
    icon: "/favicon.ico"
  }
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full" suppressHydrationWarning>
      <body
        className={`${inter.className} h-full`}
      >
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
