import { Status } from '@/types/task';
import { cn } from '@/lib/utils';
import { CheckCircle2, Clock, RotateCw } from 'lucide-react';

interface StatusBadgeProps {
  status: Status;
  className?: string;
  showLabel?: boolean;
  size?: 'sm' | 'default';
}

export function StatusBadge({
  status,
  className,
  showLabel = true,
  size = 'default',
}: StatusBadgeProps) {
  const iconSize = size === 'sm' ? 12 : 14;

  const config = {
    Completed: {
      icon: <CheckCircle2 size={iconSize} />,
      label: 'Completed',
      className:
        'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 border-green-200 dark:border-green-800',
    },
    'In Progress': {
      icon: <RotateCw size={iconSize} />,
      label: 'In Progress',
      className:
        'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 border-blue-200 dark:border-blue-800',
    },
    Pending: {
      icon: <Clock size={iconSize} />,
      label: 'Pending',
      className:
        'bg-gray-100 text-gray-700 dark:bg-gray-800/50 dark:text-gray-400 border-gray-200 dark:border-gray-700',
    },
  };

  const { icon, label, className: statusClass } = config[status];

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded border px-2 py-0.5 text-xs font-medium',
        statusClass,
        size === 'sm' ? 'px-1.5 py-0' : 'px-2 py-0.5',
        className,
      )}
    >
      {icon}
      {showLabel && <span>{label}</span>}
    </span>
  );
}
