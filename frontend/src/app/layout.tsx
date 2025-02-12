import type { Metadata } from "next";
import "./globals.css";
import { Inter } from 'next/font/google'
import { Providers } from './providers'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: "Lesyr",
  description: "AI Powered Feed Reader",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
    <body className={`${inter.className} bg-white text-slate-900 dark:bg-slate-900 dark:text-white`}>
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
