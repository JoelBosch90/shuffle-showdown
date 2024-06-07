import { writable } from 'svelte/store';
import { type Toast, ToastType } from '$lib/types/Toast';

const DEFAULT_TOAST: Toast = {
  id: 0,
  message: '',
  type: ToastType.Error,
  timeout: 3500,
};


export const toasts = writable<Toast[]>([]);
let latestId = 0;

export const showToast = (toast: Omit<Toast, 'id'>) => {
  const newToast = { ...DEFAULT_TOAST, ...toast, id: latestId++ };
  toasts.update(all => [...all, newToast]);

  if (newToast.timeout) {
    setTimeout(() => {
      dismissToast(newToast.id);
    }, newToast.timeout);
  }
};

export const dismissToast = (id: number) =>
  toasts.update(all => all.filter(toast => toast.id !== id));