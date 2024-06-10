import { ToastType } from '$lib/enums/ToastType';
import { showToast } from '$lib/store/toasts';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const handleError = <WrappedFunction extends (...args: any[]) => Promise<any>>(callback: WrappedFunction) => {
  return async (...args: Parameters<WrappedFunction>) : Promise<ReturnType<WrappedFunction> | null> => callback(...args).catch(error => {
    showToast({ message: error.message, type: ToastType.Error });
    return null;
  });
}