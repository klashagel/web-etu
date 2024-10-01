import './globals.css'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Andritz AB Controller Management',
  description: 'Manage and monitor controllers',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
