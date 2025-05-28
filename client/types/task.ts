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
export type FormValues = {
  title: string;
  description?: string;
  priority: Priority;
  status: Status;
  dueDate: Date | null;
  tags: string[];
};

export type TaskFormValues = Omit<Task, 'id' | 'createdAt'>;

export interface Tag {
  id: string;
  name: string;
  color: string;
}

