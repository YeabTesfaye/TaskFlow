import { useState, useEffect } from 'react';
import { Button } from './ui/button';
import { Textarea } from './ui/textarea';
import { useToast } from '@/hooks/use-toast';
import { comments } from '@/api';
import { Comment } from '@/types/task';
import { CommentCard } from './comment-card';
import { commentSchema } from '@/lib/validator';

interface CommentSectionProps {
  taskId: string;
}

export function CommentSection({ taskId }: CommentSectionProps) {
  const [commentList, setCommentList] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState('');
  const { toast } = useToast();

  // Fetch comments when component mounts
  useEffect(() => {
    const fetchComments = async () => {
      try {
        const result = await comments.getAll(taskId);
        setCommentList(Array.isArray(result) ? result : []);
      } catch (error: any) {
        toast({
          title: 'Error',
          description:
            error.response?.data?.message ||
            'Something went wrong while loading comments.',
          variant: 'destructive',
        });
      }
    };

    fetchComments();
  }, [taskId, toast]);

  const handleAddComment = async () => {
    const parsed = commentSchema.safeParse({ content: newComment });
    if (!parsed.success) {
      toast({
        title: 'Validation Error',
        description: parsed.error.errors[0].message,
        variant: 'destructive',
      });
      return;
    }
    try {
      const comment = await comments.create(taskId, { content: newComment });
      setCommentList((prev) => [...prev, comment]);
      setNewComment('');
      toast({
        title: 'Comment added',
        description: 'Your comment has been added successfully',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description: error.response?.data?.message || 'Failed to add comment',
        variant: 'destructive',
      });
    }
  };

  const handleDeleteComment = async (commentId: string) => {
    try {
      await comments.delete(taskId, commentId);
      setCommentList((prev) => prev.filter((c) => c.id !== commentId));
      toast({
        title: 'Comment deleted',
        description: 'The comment has been removed',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to delete comment',
        variant: 'destructive',
      });
    }
  };

  const handleUpdateComment = async (commentId: string, content: string) => {
    try {
      const updated = await comments.update(taskId, commentId, { content });
      setCommentList((prev) =>
        prev.map((c) => (c.id === commentId ? updated : c)),
      );
      toast({
        title: 'Comment updated',
        description: 'Your comment has been updated successfully',
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to update comment',
        variant: 'destructive',
      });
    }
  };

  return (
    <div className="mt-6">
      <h3 className="text-lg font-semibold mb-4">Comments</h3>
      <div className="space-y-4">
        <div className="flex flex-col gap-2">
          <Textarea
            placeholder="Add a comment..."
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            className="min-h-[100px]"
          />
          <Button
            onClick={handleAddComment}
            disabled={!newComment.trim()}
            className="self-end"
          >
            Add Comment
          </Button>
        </div>
        <div className="space-y-4">
          {commentList.map((comment) => (
            <CommentCard
              key={comment.id}
              comment={comment}
              onDelete={() => handleDeleteComment(comment.id)}
              onUpdate={(content) => handleUpdateComment(comment.id, content)}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
