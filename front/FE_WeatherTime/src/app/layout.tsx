import Header from "../components/header";
import Footer from "../components/footer";
import "./globals.css";

export const metadata = {
  title: "Alexander's weather-time city",
  description: "Лучшие погода и время в мире",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased">
        <Header />
        <main>{children}</main>
        <Footer />
      </body>
    </html>
  );
}