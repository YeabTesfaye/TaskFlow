import { Priority } from '@/types';
import { cn } from '@/lib/utils';
import { AlertTriangle, ArrowUp, ArrowRight, ArrowDown } from 'lucide-react';

interface PriorityBadgeProps {
  priority: Priority;
  className?: string;
  showLabel?: boolean;
  size?: 'sm' | 'default';
}

export function PriorityBadge({
  priority,
  className,
  showLabel = true,
  size = 'default',
}: PriorityBadgeProps) {
  const iconSize = size === 'sm' ? 12 : 14;

  const config = {
    Urgent: {
      icon: <AlertTriangle size={iconSize} />,
      label: 'Urgent',
      className:
        'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400 border-red-200 dark:border-red-800',
    },
    High: {
      icon: <ArrowUp size={iconSize} />,
      label: 'High',
      className:
        'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400 border-orange-200 dark:border-orange-800',
    },
    Medium: {
      icon: <ArrowRight size={iconSize} />,
      label: 'Medium',
      className:
        'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 border-blue-200 dark:border-blue-800',
    },
    Low: {
      icon: <ArrowDown size={iconSize} />,
      label: 'Low',
      className:
        'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 border-green-200 dark:border-green-800',
    },
  };

  const { icon, label, className: priorityClass } = config[priority];

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded border px-2 py-0.5 text-xs font-medium',
        priorityClass,
        size === 'sm' ? 'px-1.5 py-0' : 'px-2 py-0.5',
        className,
      )}
    >
      {icon}
      {showLabel && <span>{label}</span>}
    </span>
  );
}
