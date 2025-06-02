import { FilterState } from '@/types';

export const notificationSettings = [
  {
    label: 'Email Notifications',
    description: 'Receive notifications about your tasks via email',
    key: 'emailNotifications',
  },
  {
    label: 'Push Notifications',
    description: 'Receive notifications in your browser',
    key: 'pushNotifications',
  },
  {
    label: 'Daily Digest',
    description: 'Receive a daily summary of your tasks',
    key: 'dailyDigest',
  },
] as const;

// Initial filter
export const initialFilters: FilterState = {
  statuses: ['Pending', 'In Progress', 'Completed'],
  priorities: ['Low', 'Medium', 'High', 'Urgent'],
  tagIds: [],
  sortBy: 'createdAt',
  sortDirection: 'desc',
  view: 'grid',
};
