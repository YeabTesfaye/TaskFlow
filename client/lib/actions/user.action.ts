import { signInFormSchema, signUpFormSchema } from "../validator";
import { auth } from "@/api";
import { isRedirectError } from "next/dist/client/components/redirect";
import { formatError } from "../utils";


// Sign in the user with credentials
export async function signInWithCredentials(prevState: unknown, formData: FormData) {
  try {
    const email = formData.get('email') as string;
    const password = formData.get('password') as string;

    const validation = signInFormSchema.safeParse({ email, password });

    if (!validation.success) {
      return {
        success: false,
        message: "Validation failed",
        errors: validation.error.flatten().fieldErrors,
      };
    }

    const user =   await auth.login(email, password);
    // if (user.token) {
    //   cookies().set("token", user.token); // secure, httpOnly = false if you want JS access
    // }
    return { success: true, message: "Signed in successfully", user };
  } catch (error: any) {
    if (isRedirectError(error)) throw error;

    return {
      success: false,
      message: formatError(error),
    };
  }
}

// Register a new user
export async function signUp(prevState: unknown, formData: FormData) {
  try {
    const user = signUpFormSchema.parse({
      name: formData.get("name"),
      email: formData.get("email"),
      password: formData.get("password"),
      confirmPassword: formData.get("confirmPassword"),
    });

    const { email, password, name } = user;
    await auth.signup(email, password, name);

    return { success: true, message: "User created successfully" };
  } catch (error) {
    if (isRedirectError(error)) throw error;

    return {
      success: false,
      message: formatError(error),
    };
  }
}
