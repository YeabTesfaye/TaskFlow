/* File: hooks/useProfile.ts */
import { useState, useEffect } from 'react';
import { user } from '@/api';
import { useToast } from '@/hooks/use-toast';

export function useProfile() {
  const [profile, setProfile] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const { toast } = useToast();

  const loadProfile = async () => {
    try {
      const data = await user.getProfile();
      setProfile(data);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to load profile',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProfile();
  }, []);

  return { profile, setProfile, loading, loadProfile };
}
