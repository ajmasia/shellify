import { useEffect } from 'react'

const FAVICON_LIGHT = '/shellify-light.svg'
const FAVICON_DARK = '/shellify-dark.svg'

function setFavicon(href: string): void {
  const link = document.querySelector<HTMLLinkElement>("link[rel*='icon']")
  if (link) {
    link.href = href
  }
}

export function useFaviconTheme(): void {
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')

    const updateFavicon = (isDark: boolean): void => {
      // Dark system theme → light favicon (for visibility)
      // Light system theme → dark favicon (for visibility)
      setFavicon(isDark ? FAVICON_LIGHT : FAVICON_DARK)
    }

    // Set initial favicon
    updateFavicon(mediaQuery.matches)

    // Listen for theme changes
    const handleChange = (e: MediaQueryListEvent): void => {
      updateFavicon(e.matches)
    }

    mediaQuery.addEventListener('change', handleChange)

    return () => {
      mediaQuery.removeEventListener('change', handleChange)
    }
  }, [])
}
