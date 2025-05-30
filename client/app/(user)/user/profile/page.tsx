'use client';

import { AvatarUpload } from '@/components/profile/AvatarUpload';
import { ProfileForm } from '@/components/profile/ProfileForm';
import { PasswordForm } from '@/components/profile/PasswordForm';
import { DeleteAccount } from '@/components/profile/DeleteAccount';
import { useAuth } from '@/hooks/use-auth';
import { useProfile } from '@/hooks/useProfile';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Lock, Trash2, User } from 'lucide-react';

export default function ProfilePage() {
  const { isAuthenticated, loading: authLoading } = useAuth();
  const { profile, setProfile, loading, loadProfile } = useProfile();
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !isAuthenticated) router.push('/login');
  }, [authLoading, isAuthenticated, router]);

  if (loading || !profile)
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );

  return (
    <div className="max-w-4xl mx-auto py-8 px-4 space-y-8">
      <div className="flex items-center space-x-6 pb-6">
        <AvatarUpload profile={profile} reload={loadProfile} />
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
            <ProfileForm
              profile={profile}
              setProfile={setProfile}
              reload={loadProfile}
            />
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
            <PasswordForm />
          </CardContent>
        </Card>
      </div>

      <Card className="border-destructive/50">
        <CardHeader className="flex flex-row items-center space-x-4 pb-2">
          <Trash2 className="w-5 h-5 text-destructive" />
          <div>
            <CardTitle>Delete Account</CardTitle>
            <CardDescription>Permanently remove your account</CardDescription>
          </div>
        </CardHeader>
        <Separator className="mb-6" />
        <CardContent>
          <DeleteAccount />
        </CardContent>
      </Card>
    </div>
  );
}
