import TierList from "./Tierlist";

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
    <div className="space-y-4">
      <div className="w-full max-w-4xl mx-auto bg-gray-800 p-6 rounded-lg">
        <iframe
          src={`https://open.spotify.com/embed/playlist/${playlist}?utm_source=generator`}
          width="100%"
          height="352"
          allowFullScreen={false}
          allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
          loading="lazy"
        ></iframe>
      </div>
      <TierList data={guess} items={songs} />
      <TierList data={ranking} items={songs} initialState={initState} />
    </div>
  );
}
