import type { ToastType } from '$lib/enums/ToastType';

export interface Toast {
  id: number;
  message: string;
  type: ToastType;
  timeout?: number;
}