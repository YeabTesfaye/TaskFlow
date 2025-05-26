"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { motion } from "framer-motion";
import Link from "next/link";
import SignUpForm from "./signup-form";

export default function SignupPage() {
  return (
    <div className="container flex min-h-screen items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="w-full max-w-sm"
      >
        <Card className="w-full">
          <CardHeader>
            <CardTitle className="text-2xl font-bold text-center">Sign Up</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <SignUpForm />
            <div className="mt-4 text-center text-sm">
              Already have an account?{" "}
              <Link href="/login" className="text-primary hover:underline">
                Login
              </Link>
            </div>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  );
}
