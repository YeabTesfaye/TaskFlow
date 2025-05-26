'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useToast } from '@/hooks/use-toast';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from '@/components/ui/card';
import { user } from '@/api';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { Trash2, User, Lock } from 'lucide-react';
import { updatePassword, updateProfile } from '@/lib/actions/profile.action';

export default function ProfilePage() {
  const [isLoading, setIsLoading] = useState(false);
  const [profile, setProfile] = useState<any>(null);
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [isProfileLoading, setIsProfileLoading] = useState(false);
  const [isPasswordLoading, setIsPasswordLoading] = useState(false)
  const { toast } = useToast();
  const router = useRouter();

  useEffect(() => {
    loadProfile();
  }, []);

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
    }
  };


  const handleUpdateProfile = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsProfileLoading(true);
    const formData = new FormData(e.currentTarget);
    const result = await updateProfile(null, formData);
    
    if (result.success) {
      toast({
        title: 'Success',
        description: result.message,
      });
      loadProfile();
    } else {
      toast({
        title: 'Error',
        description: result.message,
        variant: 'destructive',
      });
    }
    setIsProfileLoading(false);
  };


  

  const handlePasswordUpdate = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsPasswordLoading(true);
    const formData = new FormData(e.currentTarget);
    const result = await updatePassword(null, formData);

    if (result.success) {
      toast({
        title: 'Success',
        description: result.message,
      });
      setCurrentPassword('');
      setNewPassword('');
    } else {
      toast({
        title: 'Error',
        description: result.message,
        variant: 'destructive',
      });
    }
    setIsPasswordLoading(false);
  };

  const handleDeleteAccount = async () => {
    if (!window.confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
      return;
    }
    setIsLoading(true);
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
      setIsLoading(false);
    }
  };

  if (!profile) return (
    <div className="flex items-center justify-center min-h-[50vh]">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
    </div>
  );

  return (
    <div className="max-w-4xl mx-auto py-8 px-4 space-y-8">
      <div className="flex items-center space-x-6 pb-6">
        <Avatar className="h-24 w-24">
          <AvatarImage src={profile.avatar || ''} />
          <AvatarFallback className="text-2xl">{profile.name?.[0]?.toUpperCase()}</AvatarFallback>
        </Avatar>
        <div>
          <h1 className="text-2xl font-bold">{profile.name}</h1>
          <p className="text-muted-foreground">{profile.email}</p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center space-x-4 pb-2">
            <User className="w-5 h-5 text-primary" />
            <div>
              <CardTitle>Profile Information</CardTitle>
              <CardDescription>Update your personal details</CardDescription>
            </div>
          </CardHeader>
          <Separator className="mb-6" />
          <CardContent>
            <form onSubmit={handleUpdateProfile} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Name</label>
                <Input
                name='name'
                  value={profile.name || ''}
                  onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                  disabled={isLoading}
                  placeholder="Your name"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Email</label>
                <Input
                name='email'
                  value={profile.email || ''}
                  // onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                  disabled
                  placeholder="your.email@example.com"
                />
              </div>
            <Button type="submit" className="w-full" disabled={isProfileLoading}>
              {isProfileLoading ? 'Saving...' : 'Save Changes'}
          </Button>

            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center space-x-4 pb-2">
            <Lock className="w-5 h-5 text-primary" />
            <div>
              <CardTitle>Security</CardTitle>
              <CardDescription>Manage your password</CardDescription>
            </div>
          </CardHeader>
          <Separator className="mb-6" />
          <CardContent>
            <form onSubmit={handlePasswordUpdate} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Current Password</label>
                <Input
             name="currentPassword"
            type="password"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
              disabled={isLoading}
            placeholder="••••••••"
            />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">New Password</label>
                <Input
             name="newPassword"
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              disabled={isLoading}
             placeholder="••••••••"
         />
              </div>
              <Button type="submit" className="w-full" disabled={isPasswordLoading}>
              {isPasswordLoading ? 'Updating...' : 'Update Password'}
                </Button>
            </form>
          </CardContent>
        </Card>
      </div>

      <Card className="border-destructive/50">
        <CardHeader className="flex flex-row items-center space-x-4 pb-2">
          <Trash2 className="w-5 h-5 text-destructive" />
          <div>
            <CardTitle>Delete Account</CardTitle>
            <CardDescription>Permanently remove your account and all data</CardDescription>
          </div>
        </CardHeader>
        <Separator className="mb-6" />
        <CardContent>
          <p className="text-sm text-muted-foreground mb-4">
            Once you delete your account, there is no going back. Please be certain.
          </p>
          <Button
            variant="destructive"
            onClick={handleDeleteAccount}
            disabled={isLoading}
            className="w-full"
          >
            {isLoading ? 'Deleting...' : 'Delete Account'}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}