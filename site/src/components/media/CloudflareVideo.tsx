'use client';

import React from 'react';

export interface CloudflareVideoProps {
  url: string;
  title?: string;
  className?: string;
}

export function getEmbedInfoFromUrl(
  inputUrl: string
): { kind: 'iframe'; src: string } | { kind: 'video'; src: string; type?: string } | null {
  try {
    const url = new URL(inputUrl);
    const host = url.hostname.toLowerCase();
    const path = url.pathname.replace(/\/+$/, '');
    const segments = path.split('/').filter(Boolean);
    const lastSegment = segments[segments.length - 1] || '';
    const cleanHref = inputUrl.split('#')[0];

    // Direct file links
    if (/\.mp4(\?.*)?$/i.test(cleanHref) || cleanHref.includes('/downloads/')) {
      return { kind: 'video', src: inputUrl, type: 'video/mp4' };
    }
    if (/\.m3u8(\?.*)?$/i.test(cleanHref) || cleanHref.includes('/manifest/')) {
      return { kind: 'video', src: inputUrl, type: 'application/x-mpegURL' };
    }

    // Already an iframe embed URL
    if (host === 'iframe.videodelivery.net' || host === 'iframe.cloudflarestream.com') {
      return { kind: 'iframe', src: inputUrl };
    }

    // Common Cloudflare Stream hostnames -> derive iframe embed from id
    if (
      host.endsWith('cloudflarestream.com') ||
      host.includes('watch.cloudflarestream.com') ||
      host.endsWith('videodelivery.net')
    ) {
      // Handle customer subdomain style: /:id/iframe
      let id = '';
      if (segments.length >= 2 && lastSegment === 'iframe') {
        id = segments[segments.length - 2];
      } else if (segments.length >= 1) {
        id = segments[0];
      }

      if (id) {
        const query = url.search ? url.search : '';
        return { kind: 'iframe', src: `https://iframe.videodelivery.net/${id}${query}` };
      }
    }

    return null;
  } catch {
    // If it's not a URL but looks like a Stream video id, try to embed
    if (/^[a-zA-Z0-9_-]{20,}$/.test(inputUrl)) {
      return { kind: 'iframe', src: `https://iframe.videodelivery.net/${inputUrl}` };
    }
    return null;
  }
}

const CloudflareVideo: React.FC<CloudflareVideoProps> = ({ url, title, className }) => {
  const info = getEmbedInfoFromUrl(url);
  if (!info) return null;

  return (
    <div className={`relative w-full h-0 pb-[56.25%] rounded-lg overflow-hidden shadow-lg ${className || ''}`}>
      {info.kind === 'iframe' ? (
        <iframe
          src={info.src}
          title={title || 'Video'}
          className="absolute top-0 left-0 w-full h-full"
          frameBorder={0}
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; fullscreen; web-share"
          allowFullScreen
          referrerPolicy="no-referrer"
        />
      ) : (
        <video
          className="absolute top-0 left-0 w-full h-full"
          src={info.src}
          controls
          playsInline
          preload="metadata"
        />
      )}
    </div>
  );
};

export default CloudflareVideo;


