import { useState } from 'react';
import { formatDistanceToNow } from 'date-fns';
import { Pencil, Trash2, X } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';
import { Card, CardContent, CardFooter } from './ui/card';
import { Textarea } from './ui/textarea';
import { Button } from './ui/button';
import { Comment } from '@/types';

interface CommentProps {
  comment: Comment;
  onDelete: () => Promise<void>;
  onUpdate: (content: string) => Promise<void>;
}

export function CommentCard({ comment, onDelete, onUpdate }: CommentProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [content, setContent] = useState(comment?.content);
  const [isUpdating, setIsUpdating] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const { toast } = useToast();

  const handleUpdate = async () => {
    if (!content.trim()) {
      toast({
        title: 'Invalid input',
        description: 'Comment content cannot be empty.',
        variant: 'destructive',
        duration: 1000,
      });
      return;
    }
    setIsUpdating(true);
    try {
      await onUpdate(content.trim());
      setIsEditing(false);
      toast({
        title: 'Comment updated',
        description: 'Your comment has been updated successfully.',
        duration: 1000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to update comment.',
        variant: 'destructive',
        duration: 1000,
      });
    } finally {
      setIsUpdating(false);
    }
  };

  const handleDelete = async () => {
    setIsDeleting(true);
    try {
      await onDelete();
      toast({
        title: 'Comment deleted',
        description: 'The comment has been removed.',
        duration: 1000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to delete comment.',
        variant: 'destructive',
        duration: 1000,
      });
    } finally {
      setIsDeleting(false);
    }
  };

  const cancelEdit = () => {
    setContent(comment.content);
    setIsEditing(false);
  };

  const isValidDate =
    comment.createdAt instanceof Date && !isNaN(comment.createdAt.getTime());

  return (
    <Card className="mb-4">
      <CardContent className="pt-4">
        {isEditing ? (
          <Textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="min-h-[100px]"
            disabled={isUpdating || isDeleting}
          />
        ) : (
          <p className="text-sm text-gray-700">{comment.content}</p>
        )}
      </CardContent>
      <CardFooter className="flex justify-between text-xs text-gray-500">
        <span>
          {isValidDate
            ? `${formatDistanceToNow(comment.createdAt)} ago`
            : 'Unknown time'}
        </span>
        <div className="flex gap-2">
          {isEditing ? (
            <>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleUpdate}
                disabled={isUpdating || !content.trim() || isDeleting}
              >
                Save
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={cancelEdit}
                disabled={isUpdating || isDeleting}
              >
                <X className="h-4 w-4" />
              </Button>
            </>
          ) : (
            <>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setIsEditing(true)}
                disabled={isDeleting || isUpdating}
              >
                <Pencil className="h-4 w-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleDelete}
                disabled={isDeleting || isUpdating}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            </>
          )}
        </div>
      </CardFooter>
    </Card>
  );
}
