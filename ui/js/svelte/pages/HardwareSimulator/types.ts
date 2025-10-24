export type HardwareSimulatorError = {
  line?: number;
  column?: number;
  message: string;
};

export type SimulationSpeed = {
  text: string;
  delayMs: number;
};
