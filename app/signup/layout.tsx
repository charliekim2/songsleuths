import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Song Sleuths - Sign up",
  description: "Join Song Sleuths to create and play games",
};

export default function SignupLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return children;
}
