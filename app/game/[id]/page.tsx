"use client";

import SubmitSongs from "@/components/SubmitSongs";
import { useState, useEffect } from "react";
import { auth } from "@/components/auth";
import { useParams, useRouter } from "next/navigation";
import { Dispatch, SetStateAction } from "react";
import { User } from "firebase/auth";

async function GetGame(
  gameId: string,
  setHook: Dispatch<SetStateAction<Game | null>>,
  user: User,
) {
  try {
    const token = await user.getIdToken();
    const response = await fetch(`/api/games/${gameId}`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });
    if (response.ok) {
      const data = await response.json();
      console.log(data);
      setHook(data);
    } else {
      throw new Error("Failed to fetch game data");
    }
  } catch (error) {
    console.error(error);
  }
}

export default function Home() {
  const [loggedIn, setLoggedIn] = useState(false);
  const router = useRouter();
  const { id } = useParams<{ id: string }>();
  const [game, setGame] = useState<Game | null>(null);
  useEffect(() => {
    auth.onAuthStateChanged((user) => {
      if (!user) {
        router.push("/login?redirect=/game/<id>"); // fetch id from route and use to return here after login
      } else {
        setLoggedIn(true);
        GetGame(id, setGame, user);
      }
    });
  }, []);
  return (
    <div className="flex flex-col gap-8 py-8">
      {loggedIn && game ? (
        <SubmitSongs
          gameId={game.id}
          title={game.name}
          deadline={game.deadline}
          numSongs={game.n_songs}
          nickname={game.submission?.nickname}
          songs={game.submission?.songs}
          drawing={game.submission?.drawing}
        />
      ) : (
        "loading..."
      )}
    </div>
  );
}
