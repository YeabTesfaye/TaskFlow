'use client';

import { UploadCloud } from 'lucide-react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useToast } from '@/hooks/use-toast';
import { user } from '@/api';

interface AvatarUploadProps {
  profile: any;
  reload: () => void;
}

export function AvatarUpload({ profile, reload }: AvatarUploadProps) {
  const { toast } = useToast();
  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith('image/')) {
      toast({
        title: 'Error',
        description: 'Upload an image',
        variant: 'destructive',
        duration: 1000,
      });
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      toast({
        title: 'Error',
        description: 'Max size 5MB',
        variant: 'destructive',
        duration: 1000,
      });
      return;
    }

    try {
      await user.uploadProfilePicture(file);
      toast({
        title: 'Success',
        description: 'Profile picture updated',
        duration: 1000,
      });
      reload();
    } catch (err) {
      toast({
        title: 'Error',
        description: 'Upload failed',
        variant: 'destructive',
        duration: 1000,
      });
    } finally {
      e.target.value = '';
    }
  };

  return (
    <div className="relative">
      <Avatar className="h-24 w-24">
        <AvatarImage
          src={
            profile?.profile_picture?.url &&
            `http://localhost:8080${profile.profile_picture.url}`
          }
        />
        <AvatarFallback className="text-2xl">
          {profile.name?.[0]}
        </AvatarFallback>
      </Avatar>
      <label
        htmlFor="upload-avatar"
        className="absolute bottom-0 right-0 p-2 bg-primary hover:bg-primary/90 rounded-full cursor-pointer"
      >
        <UploadCloud className="h-4 w-4 text-white" />
      </label>
      <input
        type="file"
        id="upload-avatar"
        accept="image/*"
        className="hidden"
        onChange={handleUpload}
      />
    </div>
  );
}
