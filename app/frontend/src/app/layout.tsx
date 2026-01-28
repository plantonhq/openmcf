import { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { Layout } from '@/components';
import { Providers } from '@/components/providers';
import { getAllCookiesParsed, getCookieThemeMode, getCookieNavbarOpen } from '@/lib/server/cookies';

export const metadata: Metadata = {
  title: 'OpenMCF Web App',
  description: 'OpenMCF Web Application',
};

const inter = Inter({
  weight: ['400', '500', '600', '700'],
  subsets: ['latin'],
  display: 'swap',
});

export default function RootLayout({ children }: { children: React.ReactNode }) {
  // Read from API_URL (server-side, runtime) or fall back to NEXT_PUBLIC_API_URL (build-time)
  // Server components can read regular env vars at request time
  const connectHost = process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:50051';

  // Get cookies for SSR
  const allCookies = getAllCookiesParsed();
  const cookieThemeMode = getCookieThemeMode(allCookies);
  const cookieNavbarOpen = getCookieNavbarOpen(allCookies);

  return (
    <html lang="en" className={inter.className}>
      <body>
        <Providers
          connectHost={connectHost}
          cookieThemeMode={cookieThemeMode}
          cookieNavbarOpen={cookieNavbarOpen}
        >
          <Layout>{children}</Layout>
        </Providers>
      </body>
    </html>
  );
}
