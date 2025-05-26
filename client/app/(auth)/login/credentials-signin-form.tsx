"use client"
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import { signInWithCredentials } from "@/lib/actions/user.action";
import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";


const CredentialsSignInForm = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const toast = useToast();
    const router = useRouter();
    const searchParams = useSearchParams();
  
    const handleSubmit = async (e: React.FormEvent) => {
      e.preventDefault();
      setIsLoading(true);
  
      try {
        const formData = new FormData();
        formData.append("email", email);
        formData.append("password", password);
  
        const result = await signInWithCredentials(null, formData);
  
        if (result.success) {
          // Save token if your signInWithCredentials returns it (optional)
          document.cookie = `token=${result.user.token}; path=/; max-age=${60 * 60 * 24 * 7}`;
          localStorage.setItem("token", result.user?.token)
  
          toast.toast({
            title: "Success",
            description: result.message,
          });
  
          // Redirect to home or next param
          const redirectTo = searchParams.get("next") || "/";
          router.push(redirectTo);
        } else {
          toast.toast({
            title: "Error",
            description: result.message || "Failed to login",
            variant: "destructive",
          });
        }
      } catch (error) {
        toast.toast({
          title: "Error",
          description: "An unexpected error occurred",
          variant: "destructive",
        });
      } finally {
        setIsLoading(false);
      }
    };
  
    return (
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-2">
          <Input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            disabled={isLoading}
          />
        </div>
        <div className="space-y-2">
          <Input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            disabled={isLoading}
          />
        </div>
        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ? "Logging in..." : "Log in"}
        </Button>
      </form>
    );
  };

  export default CredentialsSignInForm