import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Song Sleuths",
  description: "Submit, guess, and rank songs with friends",
};

export default function GameLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return <main>{children}</main>;
}
