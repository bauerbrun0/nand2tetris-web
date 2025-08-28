export {};

declare global {
  interface Window {
    WASM: {
      // exported JS functions (called *from Go*)
      setProgressWASM: (str: string) => void;

      // exported Go functions (called *from JS*)
      startComputing: (n: number, delayNS: number) => void;
    };
  }
}
