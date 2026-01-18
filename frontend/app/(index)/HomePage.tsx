import CTA from "./CTA"
import Demo from "./Demo"
import FAQ from "./FAQ"
import Features from "./Features"
import Footer from "./Footer"
import Header from "./Header"
import Hero from "./Hero"

export default function HomePage() {
  return (
    <div className="bg-base">
      <Header />
      <main className="isolate">
        <Hero />
        <Demo />
        <Features />
        <FAQ />
        <CTA />
      </main>
      <Footer />
    </div>
  )
}
