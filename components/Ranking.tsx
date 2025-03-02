import Link from "next/link";
import TierList from "./Tierlist";
import { Button } from "./ui/button";

interface RankingProps {
  guess?: Tierlist;
  ranking?: Tierlist;
  playlist: string;
  songs?: Song[];
}

export default function Ranking({
  guess,
  ranking,
  playlist,
  songs,
}: RankingProps) {
  if (!guess || !ranking || !songs) {
    return <div className="text-white">Nothing...</div>;
  }
  const initState = JSON.parse(localStorage.getItem("list_state") ?? "{}");
  return (
    <div>
      <Button asChild>
        <Link href={`https://open.spotify.com/playlist/${playlist}`}>
          Spotify Playlist
        </Link>
      </Button>
      <TierList data={guess} items={songs} />
      <TierList data={ranking} items={songs} initialState={initState} />
    </div>
  );
}
