'use client';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Priority, PRIORITY_VALUES, Status, STATUS_VALUES } from '@/types';
import { useState } from 'react';
import { Filter, SortAsc, SortDesc } from 'lucide-react';
import { motion } from 'framer-motion';

interface TaskFiltersProps {
  onFilterChange(filters: FilterState): void;

  tags: { id: string; name: string; color: string }[];
}

export interface FilterState {
  statuses: Status[];
  priorities: Priority[];
  tagIds: string[];
  sortBy: 'dueDate' | 'priority' | 'createdAt';
  sortDirection: 'asc' | 'desc';
  view: 'grid' | 'list';
}

const initialFilters: FilterState = {
  statuses: ['Pending', 'In Progress', 'Completed'],
  priorities: ['Low', 'Medium', 'High', 'Urgent'],
  tagIds: [],
  sortBy: 'createdAt',
  sortDirection: 'desc',
  view: 'grid',
};

export function TaskFilters({ onFilterChange, tags }: TaskFiltersProps) {
  const [filters, setFilters] = useState<FilterState>(initialFilters);

  const updateFilters = (newFilters: Partial<FilterState>) => {
    const updatedFilters = { ...filters, ...newFilters };
    setFilters(updatedFilters);
    onFilterChange(updatedFilters);
  };

  const toggleStatus = (status: Status) => {
    const statuses = filters.statuses.includes(status)
      ? filters.statuses.filter((s) => s !== status)
      : [...filters.statuses, status];
    updateFilters({ statuses });
  };

  const togglePriority = (priority: Priority) => {
    const priorities = filters.priorities.includes(priority)
      ? filters.priorities.filter((p) => p !== priority)
      : [...filters.priorities, priority];
    updateFilters({ priorities });
  };

  const toggleTag = (tagId: string) => {
    const tagIds = filters.tagIds.includes(tagId)
      ? filters.tagIds.filter((id) => id !== tagId)
      : [...filters.tagIds, tagId];
    updateFilters({ tagIds });
  };

  const sortOptions = [
    { value: 'createdAt', label: 'Creation Date' },
    { value: 'dueDate', label: 'Due Date' },
    { value: 'priority', label: 'Priority' },
  ] as const;

  return (
    <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex flex-wrap gap-2">
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                <Filter className="h-4 w-4" />
                <span>Status</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              <DropdownMenuLabel>Filter by Status</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuCheckboxItem
                checked={filters.statuses.includes(STATUS_VALUES[0])}
                onCheckedChange={() => toggleStatus(STATUS_VALUES[0])}
              >
                Pending
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.statuses.includes(STATUS_VALUES[1])}
                onCheckedChange={() => toggleStatus(STATUS_VALUES[1])}
              >
                In Progress
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.statuses.includes(STATUS_VALUES[2])}
                onCheckedChange={() => toggleStatus(STATUS_VALUES[2])}
              >
                Completed
              </DropdownMenuCheckboxItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.1 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                <Filter className="h-4 w-4" />
                <span>Priority</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              <DropdownMenuLabel>Filter by Priority</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes(PRIORITY_VALUES[0])}
                onCheckedChange={() => togglePriority(PRIORITY_VALUES[0])}
              >
                Low
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes(PRIORITY_VALUES[1])}
                onCheckedChange={() => togglePriority(PRIORITY_VALUES[1])}
              >
                Medium
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes(PRIORITY_VALUES[2])}
                onCheckedChange={() => togglePriority(PRIORITY_VALUES[2])}
              >
                High
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={filters.priorities.includes(PRIORITY_VALUES[3])}
                onCheckedChange={() => togglePriority(PRIORITY_VALUES[3])}
              >
                Urgent
              </DropdownMenuCheckboxItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.2 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                <Filter className="h-4 w-4" />
                <span>Tags</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-48">
              <DropdownMenuLabel>Filter by Tags</DropdownMenuLabel>
              <DropdownMenuSeparator />
              {tags.length > 0 ? (
                tags.map((tag) => (
                  <DropdownMenuCheckboxItem
                    key={tag.id}
                    checked={filters.tagIds.includes(tag.id)}
                    onCheckedChange={() => toggleTag(tag.id)}
                  >
                    <span className="inline-flex items-center gap-1.5">
                      <span
                        className="h-2 w-2 rounded-full"
                        style={{ backgroundColor: tag.color }}
                      ></span>
                      {tag.name}
                    </span>
                  </DropdownMenuCheckboxItem>
                ))
              ) : (
                <div className="px-2 py-1.5 text-sm text-muted-foreground">
                  No tags available
                </div>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>
      </div>

      <div className="flex items-center gap-2">
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: 0.3 }}
        >
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="gap-1">
                {filters.sortDirection === 'asc' ? (
                  <SortAsc className="h-4 w-4" />
                ) : (
                  <SortDesc className="h-4 w-4" />
                )}
                <span>Sort</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <DropdownMenuLabel>Sort by</DropdownMenuLabel>
              <DropdownMenuSeparator />
              {sortOptions.map((option) => (
                <DropdownMenuCheckboxItem
                  key={option.value}
                  checked={filters.sortBy === option.value}
                  onCheckedChange={() =>
                    updateFilters({ sortBy: option.value })
                  }
                >
                  {option.label}
                </DropdownMenuCheckboxItem>
              ))}
              <DropdownMenuSeparator />
              <DropdownMenuCheckboxItem
                checked={filters.sortDirection === 'asc'}
                onCheckedChange={() =>
                  updateFilters({
                    sortDirection:
                      filters.sortDirection === 'asc' ? 'desc' : 'asc',
                  })
                }
              >
                Ascending Order
              </DropdownMenuCheckboxItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </motion.div>
      </div>
    </div>
  );
}
