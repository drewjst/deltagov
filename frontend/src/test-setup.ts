// Polyfill ResizeObserver for tests
if (typeof ResizeObserver === 'undefined') {
  // @ts-expect-error - testing mock
  global.ResizeObserver = class ResizeObserver {
    observe(): void {
        // mock
    }
    unobserve(): void {
        // mock
    }
    disconnect(): void {
        // mock
    }
  };
}
