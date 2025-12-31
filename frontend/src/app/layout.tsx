import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Tau-Tau Run 5K | Event Registration',
  description: 'Register for the Tau-Tau Run Fun Run 5K event',
  keywords: ['fun run', '5k', 'tau-tau run', 'event registration'],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="antialiased bg-gray-50">
        {children}
      </body>
    </html>
  )
}
