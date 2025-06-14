'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { format } from 'date-fns';
import { CalendarIcon } from 'lucide-react';
import { cn } from '@/lib/utils';
import { PRIORITY_VALUES, STATUS_VALUES, Tag, Task } from '@/types';
import { useToast } from '@/hooks/use-toast';
import { useRouter } from 'next/navigation';
import { Checkbox } from '@/components/ui/checkbox';
import { tasks } from '@/api';
import { formSchema } from '@/lib/validator';
import { useTags } from '@/hooks/useTags';
import { CollaboratorManager } from './collaborator-manager';

type FormValues = z.infer<typeof formSchema>;

interface TaskFormProps {
  task?: Task;
  mode: 'create' | 'edit';
}

export function TaskForm({ task, mode }: TaskFormProps) {
  const { toast } = useToast();
  const router = useRouter();
  const { tagList } = useTags() as { tagList: Tag[] };

  const defaultValues: FormValues = task
    ? {
        title: task.title,
        description: task.description,
        priority: task.priority,
        status: task.status,
        dueDate: task.dueDate ? new Date(task.dueDate) : null,
        tags: task.tags,
        collaborators: task.collaborators || [],
      }
    : {
        title: '',
        description: '',
        priority: 'Medium',
        status: 'Pending',
        dueDate: new Date(),
        tags: [],
        collaborators: [],
      };

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues,
  });

  async function onSubmit(values: FormValues) {
    try {
      if (mode === 'create') {
        await tasks.create({
          ...values,
          due_date: values.dueDate,
        });
        toast({
          title: 'Task Created',
          description: 'Your new task has been created successfully.',
        });
      } else if (task) {
        await tasks.update(task.id, {
          ...values,
          due_date: values.dueDate,
        });
        toast({
          title: 'Task Updated',
          description: 'Your task has been updated successfully.',
        });
      }
      router.push('/');
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error?.response?.data?.message || 'Something went wrong',
        variant: 'destructive',
      });
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        {/* Title Field */}
        <FormField
          control={form.control}
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Task Title</FormLabel>
              <FormControl>
                <Input
                  placeholder="Enter a clear, specific task title"
                  {...field}
                  autoFocus
                />
              </FormControl>
              <FormDescription>Keep it short and descriptive</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Description Field */}
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="Add details about your task..."
                  className="min-h-[120px] resize-none"
                  {...field}
                  value={field.value || ''}
                />
              </FormControl>
              <FormDescription>
                 Include any important details or context
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Priority & Status */}
        <div className="grid gap-6 md:grid-cols-2">
          <FormField
            control={form.control}
            name="priority"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Priority</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select priority" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {PRIORITY_VALUES.map((p) => (
                      <SelectItem key={p} value={p}>
                        {p}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Status</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select status" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {STATUS_VALUES.map((s) => (
                      <SelectItem key={s} value={s}>
                        {s}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        {/* Due Date */}
        <FormField
          control={form.control}
          name="dueDate"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Due Date</FormLabel>
              <Popover>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant={'outline'}
                      className={cn(
                        'w-full pl-3 text-left font-normal',
                        !field.value && 'text-muted-foreground',
                      )}
                    >
                      {field.value ? (
                        format(field.value, 'PPP')
                      ) : (
                        <span>Pick a date</span>
                      )}
                      <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value || undefined}
                    onSelect={field.onChange}
                    disabled={(date) =>
                      date < new Date(new Date().setHours(0, 0, 0, 0))
                    }
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
              <FormDescription>
                Optional: Set a deadline for this task
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Tags Field */}
        <FormField
          control={form.control}
          name="tags"
          render={({ field }) => (
            <FormItem>
              <div className="mb-4">
                <FormLabel className="text-base">Tags</FormLabel>
                <FormDescription>
                  Select relevant tags for this task
                </FormDescription>
              </div>
              <div className="grid grid-cols-2 gap-2 md:grid-cols-4">
                {tagList.map((tag) => (
                  <div
                    key={tag.id}
                    className="flex flex-row items-start space-x-2 space-y-0 rounded-md border p-2"
                    style={{ borderColor: `${tag.color}50` }}
                  >
                    <Checkbox
                      id={tag.id}
                      checked={field.value?.includes(tag.id)}
                      onCheckedChange={(checked) => {
                        const currentTags = field.value || [];
                        const newTags = checked
                          ? [...currentTags, tag.id]
                          : currentTags.filter((id) => id !== tag.id);
                        field.onChange(newTags);
                      }}
                    />
                    <FormLabel
                      htmlFor={`tag-${tag.id}`}
                      className="font-normal"
                      style={{ color: tag.color }}
                    >
                      {tag.name}
                    </FormLabel>
                  </div>
                ))}
              </div>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Collaborator Manager - only in edit mode */}
        {mode === 'edit' && task && (
          <FormField
            control={form.control}
            name="collaborators"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Collaborators</FormLabel>
                <FormControl>
                  <CollaboratorManager
                    taskId={task.id}
                    collaborators={field.value}
                    onUpdate={() => {
                      router.refresh();
                    }}
                  />
                </FormControl>
                <FormDescription>
                  Add or remove collaborators for this task
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        )}

        {/* Submit / Cancel Buttons */}
        <div className="flex justify-end space-x-2">
          <Button variant="outline" type="button" onClick={() => router.back()}>
            Cancel
          </Button>
          <Button type="submit">
            {mode === 'create' ? 'Create Task' : 'Update Task'}
          </Button>
        </div>
      </form>
    </Form>
  );
}
