import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { GlobalErrorBoundary } from "@/components/global-error-boundary";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "tany.ai · AI Client Chat Assistant",
  description:
    "Prototipe tany.ai yang menggabungkan Next.js dan backend Golang untuk menjawab calon klien secara otomatis.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id">
      <body
        className={`${geistSans.variable} ${geistMono.variable} bg-slate-950 font-sans antialiased`}
      >
        <GlobalErrorBoundary>{children}</GlobalErrorBoundary>
      </body>
    </html>
  );
}
