interface Game {
  id: string;
  name: string;
  deadline: number;
  n_songs: number;

  submission?: Submission;
  player_list?: Submission[];

  guess_list?: Tierlist;
  ranking_list?: Tierlist;
  playlist?: string;
  songs?: string[];
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
}

interface Submission {
  songs?: string[];
  nickname: string;
  drawing: string;
}
