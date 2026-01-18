import Link from "next/link"

const footerNavigation = {
  solutions: [
    { name: "For Streamers", href: "#" },
    { name: "For Podcasters", href: "#" },
    { name: "For Content Creators", href: "#" },
    { name: "For Teams", href: "#" },
    { name: "Sound Packs", href: "#" },
  ],
  support: [
    { name: "Help Center", href: "#" },
    { name: "Documentation", href: "#" },
    { name: "Community", href: "#" },
    { name: "Contact Support", href: "#" },
  ],
  company: [
    { name: "About", href: "#" },
    { name: "Blog", href: "#" },
    { name: "Careers", href: "#" },
    { name: "Press Kit", href: "#" },
  ],
  legal: [
    { name: "Terms of Service", href: "#" },
    { name: "Privacy Policy", href: "#" },
    { name: "Cookie Policy", href: "#" },
  ],
}

export default function Footer() {
  return (
    <footer className="relative mx-auto mt-32 max-w-7xl px-6 lg:px-8">
      <div className="border-t border-muted py-16 sm:py-24 lg:py-32">
        <div className="xl:grid xl:grid-cols-3 xl:gap-8">
          <img alt="Company name" src="/logo.png" className="h-9" />
          <div className="mt-16 grid grid-cols-2 gap-8 xl:col-span-2 xl:mt-0">
            <div className="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 className="text-sm/6 font-semibold text-highlight">Solutions</h3>
                <ul className="mt-6 space-y-4">
                  {footerNavigation.solutions.map((item) => (
                    <li key={item.name}>
                      <Link href={item.href} className="text-sm/6 text-base hover:text-highlight">
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
              <div className="mt-10 md:mt-0">
                <h3 className="text-sm/6 font-semibold text-highlight">Support</h3>
                <ul className="mt-6 space-y-4">
                  {footerNavigation.support.map((item) => (
                    <li key={item.name}>
                      <Link href={item.href} className="text-sm/6 text-base hover:text-highlight">
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
            <div className="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 className="text-sm/6 font-semibold text-highlight">Company</h3>
                <ul className="mt-6 space-y-4">
                  {footerNavigation.company.map((item) => (
                    <li key={item.name}>
                      <Link href={item.href} className="text-sm/6 text-base hover:text-highlight">
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
              <div className="mt-10 md:mt-0">
                <h3 className="text-sm/6 font-semibold text-highlight">Legal</h3>
                <ul className="mt-6 space-y-4">
                  {footerNavigation.legal.map((item) => (
                    <li key={item.name}>
                      <Link href={item.href} className="text-sm/6 text-base hover:text-highlight">
                        {item.name}
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}
