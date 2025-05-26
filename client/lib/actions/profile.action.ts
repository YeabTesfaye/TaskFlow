import { user } from "@/api";
import { isRedirectError } from "next/dist/client/components/redirect";
import { formatError } from "../utils";
import { passwordUpdateSchema, profileUpdateSchema } from "../validator";


// Get user profile
export async function getProfile() {
    try {
        const profile = await user.getProfile()
        return {success : true, data : profile}
    } catch (error) {
        if(isRedirectError(error)) throw error
        return {
            success: false,
            message: formatError(error),
          };
    }
}

// Update user profile
export async function updateProfile(prevState: unknown, formData: FormData) {
    try {
      const data = profileUpdateSchema.parse({
        name: formData.get("name"),
      });
  
      await user.updateProfile(data);
      return { success: true, message: "Profile updated successfully" };
    } catch (error: any) {
      if (isRedirectError(error)) throw error;
      return {
        success: false,
        message: formatError(error),
      };
    }
  }
  

  // Update password
export async function updatePassword(prevState: unknown, formData: FormData) {
    try {
      const data = passwordUpdateSchema.parse({
        currentPassword: formData.get("currentPassword"),
        newPassword: formData.get("newPassword"),
      });
  
      await user.changePassword(data.currentPassword, data.newPassword);
      return { success: true, message: "Password updated successfully" };
    } catch (error: any) {
      if (isRedirectError(error)) throw error;
      return {
        success: false,
        message: formatError(error),
      };
    }
  }
  
  // Delete account
export async function deleteAccount() {
    try {
      await user.deleteAccount();
      return { success: true, message: "Account deleted successfully" };
    } catch (error: any) {
      if (isRedirectError(error)) throw error;
      return {
        success: false,
        message: formatError(error),
      };
    }
  }