interface Game {
  id: string;
  name: string;
  deadline: number;
  n_songs: number;

  submission?: Submission;
  // player_list?: Submission[];

  guess_list?: Tierlist;
  ranking_list?: Tierlist;
  playlist?: string;
  songs?: Song[];
}

interface Song {
  id: number;
  spotify: string;
  album_art: string;
  name: string;
}

interface Tierlist {
  id: string;
  type: string;
  tiers: Tier[];
}

interface Tier {
  id: string;
  name: string;
  rank: number;
  drawing?: string;
}

interface Submission {
  songs?: string[];
  nickname: string;
  drawing: string;
}
