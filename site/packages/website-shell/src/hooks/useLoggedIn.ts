'use client';

import { useSyncExternalStore } from 'react';

const noop = () => () => {};

/**
 * On localhost the website runs standalone without the console backend.
 * A stale `planton_logged_in` cookie from a console dev session on the
 * same host would cause "Dashboard" to show instead of "Sign in / Sign up",
 * linking to a route that doesn't exist on the website dev server.
 */
const isLocalhost = () =>
  typeof window !== 'undefined' &&
  (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1');

const getCookieSnapshot = () =>
  !isLocalhost() && /planton_logged_in=/.test(document.cookie);

const getServerSnapshot = () => false;

/**
 * Reads the `planton_logged_in` cookie to determine whether the current
 * user has an active session. The cookie is a non-HttpOnly marker set by
 * NextAuth — its presence means the user is authenticated.
 *
 * Returns `false` on localhost regardless of cookie state, because the
 * website dev server has no console backend to dashboard into.
 *
 * Uses `useSyncExternalStore` so the component re-renders if the cookie
 * changes (e.g. after login/logout in another tab). The server snapshot
 * always returns `false` (SSR cannot read browser cookies).
 */
export function useLoggedIn(): boolean {
  return useSyncExternalStore(noop, getCookieSnapshot, getServerSnapshot);
}
