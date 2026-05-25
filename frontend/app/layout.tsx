import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AI \u6570\u5b57\u4eba\u89c6\u9891\u5de5\u5382",
  description: "AI \u89c6\u9891\u751f\u4ea7\u540e\u53f0",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  );
}
