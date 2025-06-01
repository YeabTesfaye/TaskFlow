'use client';

import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { updateProfile } from '@/lib/actions/profile.action';
import { useToast } from '@/hooks/use-toast';

interface ProfileFormProps {
  profile: any;
  setProfile: (data: any) => void;
  reload: () => void;
}

export function ProfileForm({ profile, setProfile, reload }: ProfileFormProps) {
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    const formData = new FormData(e.currentTarget as HTMLFormElement);
    const res = await updateProfile(null, formData);

    if (res.success) {
      toast({ title: 'Success', description: res.message, duration: 1000 });
      reload();
    } else {
      toast({
        title: 'Error',
        description: res.message,
        variant: 'destructive',
        duration: 1000,
      });
    }

    setIsLoading(false);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <label className="text-sm font-medium">Name</label>
        <Input
          name="name"
          value={profile.name}
          onChange={(e) => setProfile({ ...profile, name: e.target.value })}
        />
      </div>
      <div className="space-y-2">
        <label className="text-sm font-medium">Email</label>
        <Input name="email" value={profile.email} disabled />
      </div>
      <Button type="submit" className="w-full" disabled={isLoading}>
        {isLoading ? 'Saving...' : 'Save Changes'}
      </Button>
    </form>
  );
}
