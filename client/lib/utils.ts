import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

// Tailwind + clsx merge
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Format various errors (Zod or generic)
export function formatError(error: any): string {
  if (!error) return 'An unknown error occurred';

  // Handle Zod validation errors
  if (error.name === 'ZodError' && Array.isArray(error.errors)) {
    const messages = error.errors.map((err: any) => err.message);
    return messages.join('. ');
  }

  // Handle other errors
  return typeof error.message === 'string'
    ? error.message
    : JSON.stringify(error.message ?? 'An unknown error occurred');
}


