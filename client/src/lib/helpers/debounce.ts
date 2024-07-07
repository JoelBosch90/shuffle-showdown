type FunctionToDebounce<ReturnType> = (...args: unknown[]) => ReturnType

export const debounce = <ReturnType>(callback: FunctionToDebounce<ReturnType>, milliseconds: number = 1000) : FunctionToDebounce<Promise<ReturnType>> => {
  let timer: number | null

  return function(this: unknown, ...args: unknown[]) : Promise<ReturnType> {
    return new Promise<ReturnType>((resolve: CallableFunction) => {
      if (timer) clearTimeout(timer)

      timer = setTimeout(() => {
        timer = null
        
        resolve(callback.apply(this, args))
      }, milliseconds)
    })
  }
}