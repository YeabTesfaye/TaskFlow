export const PRIORITY_OPTIONS = ['High', 'Medium', 'Low', 'Urgent'] as const;
export const PRIORITY_VALUES = ['High', 'Medium', 'Low', 'Urgent'] as const;
export type Priority = (typeof PRIORITY_VALUES)[number];

export const STATUS_VALUES = ['Pending', 'In Progress', 'Completed'] as const;
export type Status = (typeof STATUS_VALUES)[number];

export interface Task {
  id: string;
  title: string;
  description: string;
  priority: Priority;
  status: Status;
  dueDate: Date | null;
  createdAt: Date;
  tags: string[];
}

export type TaskFormValues = Omit<Task, 'id' | 'createdAt'>;

export interface Tag {
  id: string;
  name: string;
  color: string;
}

export interface Comment {
  id: string;
  taskId: string;
  userId: string;
  content: string;
  createdAt: Date;
  updatedAt: Date;
}
export interface Stats {
  id: string;
  user_id: string;
  total_tasks: number;
  completed_tasks: number;
  pending_tasks: number;
  overdue_tasks: number;
  completion_rate: number;
  average_completion: number;
  by_priority: Record<string, number>;
  by_category: Record<string, number> | null;
  updated_at: string;
}

export interface FilterState {
  statuses: Status[];
  priorities: Priority[];
  tagIds: string[];
  sortBy: 'dueDate' | 'priority' | 'createdAt';
  sortDirection: 'asc' | 'desc';
  view: 'grid' | 'list';
}
export const sortOptions = [
  { value: 'createdAt', label: 'Creation Date' },
  { value: 'dueDate', label: 'Due Date' },
  { value: 'priority', label: 'Priority' },
] as const;
