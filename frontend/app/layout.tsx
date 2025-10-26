import type { Metadata } from "next";
import { Inter, Plus_Jakarta_Sans } from "next/font/google";
import { GlobalErrorBoundary } from "@/components/global-error-boundary";
import { headers } from "next/headers";
import "./globals.css";

const display = Plus_Jakarta_Sans({
  variable: "--font-display",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
});

const body = Inter({
  variable: "--font-body",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "tany.ai · AI Client Chat Assistant",
  description:
    "Prototipe tany.ai yang menggabungkan Next.js dan backend Golang untuk menjawab calon klien secara otomatis.",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const headersList = await headers();
  const nonce = headersList.get("x-nonce") || "";
  
  return (
    <html lang="id">
      <head>
        {nonce && <meta name="x-nonce" content={nonce} />}
      </head>
      <body
        className={`${display.variable} ${body.variable} bg-[var(--bg)] font-sans antialiased text-[var(--text)]`}
      >
        <GlobalErrorBoundary>{children}</GlobalErrorBoundary>
      </body>
    </html>
  );
}
