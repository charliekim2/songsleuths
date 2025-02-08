import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Song Sleuths - Log in",
  description: "Log in to play Song Sleuths",
};

export default function LoginLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return children;
}
