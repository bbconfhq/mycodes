export function debounce<Args extends unknown[]>(func: (...args: Args) => void, wait: number) {
  let timeout: ReturnType<typeof setTimeout> | undefined;
  const debounced = (...args: Args) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      func(...args);
    }, wait);
  };
  debounced.cancel = () => clearTimeout(timeout);
  return debounced;
}
