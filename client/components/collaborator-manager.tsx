import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useToast } from '@/hooks/use-toast';
import { tasks } from '@/api';

interface CollaboratorManagerProps {
  taskId: string;
  collaborators: string[];
  onUpdate: () => void;
}

export function CollaboratorManager({
  taskId,
  collaborators,
  onUpdate,
}: CollaboratorManagerProps) {
  const [newCollaborator, setNewCollaborator] = useState('');
  const { toast } = useToast();

  const handleAddCollaborator = async () => {
    try {
      await tasks.addCollaborator(taskId, newCollaborator);
      setNewCollaborator('');
      onUpdate();
      toast({
        title: 'Success',
        description: 'Collaborator added successfully',
        duration: 1000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to add collaborator',
        variant: 'destructive',
        duration: 1000,
      });
    }
  };

  const handleRemoveCollaborator = async (collaboratorId: string) => {
    try {
      await tasks.removeCollaborator(taskId, collaboratorId);
      onUpdate();
      toast({
        title: 'Success',
        description: 'Collaborator removed successfully',
        duration: 1000,
      });
    } catch (error: any) {
      toast({
        title: 'Error',
        description:
          error.response?.data?.message || 'Failed to remove collaborator',
        variant: 'destructive',
        duration: 1000,
      });
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <Input
          placeholder="Enter collaborator ID"
          value={newCollaborator}
          onChange={(e) => setNewCollaborator(e.target.value)}
        />
        <Button onClick={handleAddCollaborator}>Add</Button>
      </div>

      <div className="space-y-2">
        {collaborators.map((collaborator) => (
          <div key={collaborator} className="flex items-center justify-between">
            <span>{collaborator}</span>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => handleRemoveCollaborator(collaborator)}
            >
              Remove
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
}
