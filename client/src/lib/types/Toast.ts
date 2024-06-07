export enum ToastType {
  Success = 'success',
  Error = 'error',
  Warning = 'warning',
  Info = 'info',
}

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
  timeout?: number;
}