// src/hooks/useDeleteConfirmation.ts
import { useState, useCallback } from 'react';

export function useDeleteConfirmation(onConfirm: () => void) {
  const [isConfirming, setIsConfirming] = useState(false);

  const handleClick = useCallback(() => {
    if (isConfirming) {
      onConfirm();
    } else {
      setIsConfirming(true);
      // Auto-reset after 3 seconds if not clicked again
      setTimeout(() => setIsConfirming(false), 3000);
    }
  }, [isConfirming, onConfirm]);

  const cancel = useCallback(() => {
    setIsConfirming(false);
  }, []);

  return { isConfirming, handleClick, cancel };
}
