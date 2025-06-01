import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';
import {
  Pencil,
  Trash2,
  CheckCircle,
  RotateCw,
  Clock,
  MoreHorizontal,
} from 'lucide-react';
import Link from 'next/link';
import { CardHeader } from '@/components/ui/card';
import { Task } from '@/types';

interface Props {
  task: Task;
  onDelete(): void;
  onStatusChange(status: Task['status']): void;
}

export function CardHeaderSection({ task, onDelete, onStatusChange }: Props) {
  return (
    <CardHeader className="flex flex-row items-start justify-between gap-2 p-4 pb-2">
      <div className="space-y-1.5">
        <Link
          href={`/task/${task.id}`}
          className="line-clamp-2 text-base font-semibold leading-tight text-foreground hover:underline"
        >
          {task.title}
        </Link>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <MoreHorizontal className="h-4 w-4" />
            <span className="sr-only">Open menu</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-[180px]">
          <DropdownMenuItem asChild>
            <Link href={`/task/${task.id}`} className="flex items-center">
              <Pencil className="mr-2 h-4 w-4" /> Edit task
            </Link>
          </DropdownMenuItem>
          {['Completed', 'In Progress', 'Pending'].map((status) => (
            <DropdownMenuItem
              key={status}
              onClick={() => onStatusChange(status as Task['status'])}
              className="flex items-center"
            >
              {status === 'Completed' && (
                <CheckCircle className="mr-2 h-4 w-4" />
              )}
              {status === 'In Progress' && (
                <RotateCw className="mr-2 h-4 w-4" />
              )}
              {status === 'Pending' && <Clock className="mr-2 h-4 w-4" />}
              Mark {status.toLowerCase()}
            </DropdownMenuItem>
          ))}
          <DropdownMenuItem
            onClick={onDelete}
            className="flex items-center text-destructive focus:text-destructive"
          >
            <Trash2 className="mr-2 h-4 w-4" /> Delete task
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </CardHeader>
  );
}
