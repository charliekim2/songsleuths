"use client";

import AuthForm from "@/components/AuthForm";
import { auth } from "@/components/auth";

export default function LoginPage() {
  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <AuthForm auth={auth} mode="login" />
    </div>
  );
}
