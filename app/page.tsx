"use client";

import CreateGame from "@/components/CreateGame";
import { useState, useEffect } from "react";
import { auth } from "@/components/auth";
import { useRouter } from "next/navigation";

export default function Home() {
  const [loggedIn, setLoggedIn] = useState(false);

  const router = useRouter();
  useEffect(() => {
    auth.onAuthStateChanged((user) => {
      if (!user) {
        router.push("/login");
      } else {
        setLoggedIn(true);
      }
    });
  }, []);
  return (
    <div className="flex flex-col gap-8 py-8">
      {loggedIn ? <CreateGame /> : "loading..."}
    </div>
  );
}
