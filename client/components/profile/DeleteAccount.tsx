'use client';

import { Button } from '@/components/ui/button';
import { user } from '@/api';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useToast } from '@/hooks/use-toast';

export function DeleteAccount() {
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const { toast } = useToast();

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete your account?'))
      return;
    setLoading(true);
    try {
      await user.deleteAccount();
      localStorage.removeItem('token');
      router.push('/login');
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to delete account',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <Button
      variant="destructive"
      onClick={handleDelete}
      disabled={loading}
      className="w-full"
    >
      {loading ? 'Deleting...' : 'Delete Account'}
    </Button>
  );
}
