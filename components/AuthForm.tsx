"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  type Auth,
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
} from "firebase/auth";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { toast } from "@/hooks/use-toast";
import { useRouter } from "next/navigation";

const formSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(6, "Password must be at least 6 characters"),
});

type FormData = z.infer<typeof formSchema>;

interface AuthFormProps {
  auth: Auth;
  mode: "login" | "signup";
}

export default function AuthForm({ auth, mode }: AuthFormProps) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = async (data: FormData) => {
    setIsLoading(true);
    try {
      if (mode === "login") {
        await signInWithEmailAndPassword(auth, data.email, data.password);
        toast({
          title: "Logged in successfully",
          description: "Welcome back!",
        });
        router.push("/");
      } else {
        await createUserWithEmailAndPassword(auth, data.email, data.password);
        const token = await auth.currentUser?.getIdToken(true);
        const res = await fetch("/api/signup", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({}),
        });
        if (res.ok) {
          router.push("/");
        } else {
          throw new Error("Failed to create account");
        }
        toast({
          title: "Account created successfully",
          description: "Welcome to Song Sleuths!",
        });
      }
    } catch (error) {
      console.error(error);
      toast({
        title: "Authentication error",
        description: "There was a problem with your request. Please try again.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-md bg-gray-800 text-white">
      <CardHeader>
        <CardTitle className="text-2xl font-bold text-center text-purple-400">
          {mode === "login" ? "Log In" : "Sign Up"}
        </CardTitle>
        <CardDescription className="text-center text-gray-400">
          {mode === "login"
            ? "Enter your credentials to access your account"
            : "Create a new account to join Song Sleuths"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Enter your email"
                      {...field}
                      type="email"
                      className="bg-gray-700 border-gray-600 text-white"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Enter your password"
                      {...field}
                      type="password"
                      className="bg-gray-700 border-gray-600 text-white"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button
              type="submit"
              className="w-full bg-purple-600 hover:bg-purple-700 text-white"
              disabled={isLoading}
            >
              {isLoading
                ? "Processing..."
                : mode === "login"
                  ? "Log In"
                  : "Sign Up"}
            </Button>
          </form>
        </Form>
      </CardContent>
      <CardFooter className="flex justify-center">
        <p className="text-sm text-gray-400">
          {mode === "login"
            ? "Don't have an account? "
            : "Already have an account? "}
          <a
            href={mode === "login" ? "/signup" : "/login"}
            className="text-purple-400 hover:underline"
          >
            {mode === "login" ? "Sign up" : "Log in"}
          </a>
        </p>
      </CardFooter>
    </Card>
  );
}
