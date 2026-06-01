import { useCallback, useEffect, useRef, useState } from 'react'

type Poller = (signal: AbortSignal) => void | Promise<void>

type VisibilityPollingOptions = {
  enabled?: boolean
  intervalMs?: number
  runImmediately?: boolean
}

type VisibilityPollingResult = {
  error: unknown
  isPolling: boolean
  isVisible: boolean
  lastPolledAt: Date | null
  pollNow: () => Promise<void>
}

const defaultIntervalMs = 3000

function isDocumentVisible() {
  return typeof document === 'undefined' || !document.hidden
}

export function useVisibilityPolling(
  poller: Poller,
  {
    enabled = true,
    intervalMs = defaultIntervalMs,
    runImmediately = true,
  }: VisibilityPollingOptions = {},
): VisibilityPollingResult {
  const pollerRef = useRef(poller)
  const intervalRef = useRef<number | undefined>(undefined)
  const activeRequestRef = useRef<AbortController | null>(null)

  const [error, setError] = useState<unknown>(null)
  const [isPolling, setIsPolling] = useState(false)
  const [isVisible, setIsVisible] = useState(isDocumentVisible)
  const [lastPolledAt, setLastPolledAt] = useState<Date | null>(null)

  useEffect(() => {
    pollerRef.current = poller
  }, [poller])

  const stopActiveRequest = useCallback(() => {
    activeRequestRef.current?.abort()
    activeRequestRef.current = null
    setIsPolling(false)
  }, [])

  const pollNow = useCallback(async () => {
    if (!enabled || !isDocumentVisible() || activeRequestRef.current) {
      return
    }

    const controller = new AbortController()
    activeRequestRef.current = controller
    setIsPolling(true)

    try {
      await pollerRef.current(controller.signal)
      setLastPolledAt(new Date())
      setError(null)
    } catch (pollError) {
      if (!controller.signal.aborted) {
        setError(pollError)
      }
    } finally {
      if (activeRequestRef.current === controller) {
        activeRequestRef.current = null
        setIsPolling(false)
      }
    }
  }, [enabled])

  useEffect(() => {
    function clearPollingInterval() {
      if (intervalRef.current === undefined) {
        return
      }

      window.clearInterval(intervalRef.current)
      intervalRef.current = undefined
    }

    function startPollingInterval() {
      if (!enabled || !isDocumentVisible() || intervalRef.current !== undefined) {
        return
      }

      if (runImmediately) {
        void pollNow()
      }

      intervalRef.current = window.setInterval(() => {
        void pollNow()
      }, intervalMs)
    }

    function syncPollingWithVisibility() {
      const visible = isDocumentVisible()
      setIsVisible(visible)

      if (visible) {
        startPollingInterval()
      } else {
        clearPollingInterval()
        stopActiveRequest()
      }
    }

    syncPollingWithVisibility()
    document.addEventListener('visibilitychange', syncPollingWithVisibility)

    return () => {
      clearPollingInterval()
      document.removeEventListener('visibilitychange', syncPollingWithVisibility)
      stopActiveRequest()
    }
  }, [enabled, intervalMs, pollNow, runImmediately, stopActiveRequest])

  return {
    error,
    isPolling,
    isVisible,
    lastPolledAt,
    pollNow,
  }
}
