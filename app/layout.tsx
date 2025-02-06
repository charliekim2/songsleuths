import "./globals.css";
import { Metadata } from "next";
import { Inter } from "next/font/google";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Song Sleuths - Create Game",
  description: "Create a new Song Sleuths game",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <style>{`
          .canvas-draw {
            position: relative;
            width: 100%;
            height: 100%;
          }
        `}</style>
      </head>
      <body className={`${inter.className} bg-gray-900`}>{children}</body>
    </html>
  );
}
