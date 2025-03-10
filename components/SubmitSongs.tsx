"use client";

import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Card,
  CardContent,
  CardDescription,
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Slider } from "@/components/ui/slider";
import { useRef, useState } from "react";
import ReactCanvasDraw from "react-canvas-draw";
import { Search, Trash } from "lucide-react";
import CountdownTimer from "./CountdownTimer";
import { auth } from "@/components/auth";

interface SearchResults {
  id: string;
  name: string;
  album: string;
  artists: string[];
  image: string;
}

function SearchDropdown({ onSelect }: { onSelect: (url: string) => void }) {
  const [searchTerm, setSearchTerm] = useState("");
  const [searchResults, setSearchResults] = useState<SearchResults[]>([]);

  const handleSearch = async () => {
    try {
      const res = await fetch("/api/spotify/search?q=" + searchTerm);
      if (!res.ok) {
        throw new Error("Failed to fetch search results");
      }

      const data = await res.json();
      setSearchResults(data);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="p-4 bg-gray-800 rounded-md w-full">
      <div className="flex mb-4">
        <Input
          type="text"
          placeholder="Search songs..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="flex-grow mr-2 bg-gray-700 text-white"
        />
        <Button
          onClick={handleSearch}
          className="bg-purple-600 hover:bg-purple-700"
        >
          Search
        </Button>
      </div>
      <ul className="space-y-2 max-h-60 overflow-y-auto">
        {searchResults.map((result) => (
          <DropdownMenuItem
            key={result.id}
            onSelect={() =>
              onSelect(`https://open.spotify.com/track/${result.id}`)
            }
            className="cursor-pointer hover:bg-gray-700 focus:bg-gray-700"
          >
            <div className="font-semibold text-white">{result.name}</div>
            <div className="text-sm text-gray-400">{result.artists.join()}</div>
          </DropdownMenuItem>
        ))}
      </ul>
    </div>
  );
}

const songUrlRegex =
  /(?:https:\/\/open\.spotify\.com\/track\/|spotify:track:)([a-zA-Z0-9]{22})/;

function SongEmbed({ id }: { id: string }) {
  return (
    <div className="w-full h-fit">
      <iframe
        src={`https://open.spotify.com/embed/track/${id}?utm_source=generator`}
        width="100%"
        height="80"
        allowFullScreen={false}
        allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
        loading="lazy"
      ></iframe>
    </div>
  );
}

export default function SubmitSongs({
  gameId,
  title,
  deadline,
  numSongs,
  songs,
  nickname,
  drawing,
}: {
  gameId: string;
  title: string;
  deadline: number;
  numSongs: number;
  nickname: string | undefined;
  songs: string[] | undefined;
  drawing: string | undefined;
}) {
  const formSchema = z.object({
    nickname: z
      .string()
      .min(3, "Nickname must be at least 3 characters")
      .max(10, "Nickname must be 10 characters or less")
      .regex(/^[a-z0-9]+$/i, "Nickname must contain only letters and numbers"),
    songs: z
      .array(
        z.object({
          url: z.string().regex(songUrlRegex, "Must be a valid Spotify URL"),
        }),
      )
      .length(numSongs, `You must submit exactly ${numSongs} songs`), // set to songNum from backend
  });

  type FormData = z.infer<typeof formSchema>;

  const canvasRef = useRef<ReactCanvasDraw>(null);
  const [brushColor, setBrushColor] = useState("#000000");
  const [brushRadius, setBrushRadius] = useState(2);

  const [img, setImg] = useState(drawing);
  const [reset, setReset] = useState(0);
  const [drew, setDrew] = useState(drawing ? true : false);

  const [loading, setLoading] = useState(false);

  const songURLs = Array(numSongs)
    .fill({})
    .map(() => ({ url: "" }));
  if (songs) {
    songs.forEach((song, idx) => {
      songURLs[idx].url = "https://open.spotify.com/track/" + song;
    });
  }
  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      nickname: nickname ?? "",
      songs: songURLs,
    },
  });

  const { fields } = useFieldArray({
    name: "songs",
    control: form.control,
  });

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    const drawingDataUrl = drew
      ? // @ts-expect-error - getDataURL is not in the types
        canvasRef.current?.getDataURL(
          "png", // Export canvas data as PNG
          img ? true : false, // Export canvas data with background image
          "#FFFFFF", // Background color
        )
      : "";
    const submitData = {
      nickname: data.nickname,
      songs: data.songs.map((song) => {
        const matches = songUrlRegex.exec(song.url);
        return matches ? matches[1] : "";
      }),
      drawing: drawingDataUrl,
    };
    try {
      const token = await auth.currentUser?.getIdToken(true);
      const res = await fetch(`/api/submit/${gameId}`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(submitData),
      });
      if (!res.ok) {
        throw new Error("Failed to submit songs");
      }
      console.log("Submitted songs successfully");
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
  };

  const onDelete = async () => {
    setLoading(true);
    try {
      const token = await auth.currentUser?.getIdToken(true);
      const res = await fetch(`/api/submit/${gameId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (!res.ok) {
        throw new Error("Failed to delete songs");
      }
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
    window.location.reload();
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <Card className="w-full max-w-2xl bg-gray-800 text-white">
        <CardHeader>
          <CardTitle className="text-3xl font-bold text-center text-purple-400">
            {title}
          </CardTitle>
          <CardDescription className="text-center text-gray-400">
            Submit your songs
          </CardDescription>
          <br />
          <CountdownTimer targetTimestamp={deadline} />
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="nickname"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nickname</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Enter your nickname"
                        {...field}
                        className="bg-gray-700 border-gray-600 text-white"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              {fields.map((field, index) => {
                const songId =
                  form.watch(`songs.${index}.url`).match(songUrlRegex)?.[1] ??
                  "";
                return (
                  <div key={field.id} className="space-y-2">
                    <FormField
                      key={field.id}
                      control={form.control}
                      name={`songs.${index}.url`}
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Song {index + 1} Spotify URL</FormLabel>
                          <FormControl>
                            <div className="flex items-center bg-gray-700 border border-gray-600 rounded-md relative">
                              <Input
                                placeholder="Enter Spotify URL"
                                {...field}
                                className="bg-transparent border-none text-white flex-grow pr-10"
                              />
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="absolute right-0 top-0 h-full aspect-square rounded-l-none"
                                  >
                                    <Search className="h-4 w-4" />
                                    <span className="sr-only">
                                      Search songs
                                    </span>
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent
                                  className="w-[350px] bg-gray-800 border-gray-700"
                                  align="end"
                                >
                                  <SearchDropdown
                                    onSelect={(url) => {
                                      form.setValue(`songs.${index}.url`, url);
                                    }}
                                  />
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </div>
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    {songId && <SongEmbed id={songId} />}
                  </div>
                );
              })}
              <FormField
                name="drawing"
                render={() => (
                  <FormItem>
                    <FormLabel>Draw yourself! (optional)</FormLabel>
                    <FormControl>
                      <div className="space-y-2">
                        <div className="flex flex-col sm:flex-row space-x-2">
                          <div className="flex-grow">
                            <ReactCanvasDraw
                              key={reset}
                              ref={canvasRef}
                              brushColor={brushColor}
                              brushRadius={brushRadius}
                              backgroundColor="#FFFFFF"
                              imgSrc={img ?? ""}
                              canvasWidth={300}
                              canvasHeight={300}
                              lazyRadius={0}
                              className="border border-gray-600 rounded-md"
                              onChange={() => {
                                if (!drew) {
                                  setDrew(true);
                                }
                              }}
                            />
                          </div>
                          <div className="w-72 space-y-4">
                            <div className="flex space-x-4 items-end">
                              <div>
                                <label className="block text-sm font-medium text-gray-400 mb-1">
                                  Color
                                </label>
                                <Input
                                  type="color"
                                  value={brushColor}
                                  onChange={(e) =>
                                    setBrushColor(e.target.value)
                                  }
                                  className="w-full h-8 p-1 bg-gray-700 border-gray-600"
                                />
                              </div>
                              <Button
                                type="button"
                                variant="outline"
                                onClick={() => {
                                  const curr = canvasRef.current;
                                  if (!curr) return;
                                  curr.clear();
                                  setImg("");
                                  setDrew(false);
                                  setReset(reset + 1);
                                }}
                                className="w-full bg-gray-700 text-white hover:bg-gray-600"
                              >
                                Clear
                              </Button>
                            </div>
                            <div>
                              <label className="block text-sm font-medium text-gray-400 mb-1">
                                Thickness
                              </label>
                              <Slider
                                value={[brushRadius]}
                                onValueChange={(value) =>
                                  setBrushRadius(value[0])
                                }
                                min={1}
                                max={20}
                                step={1}
                                className="w-full"
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                    </FormControl>
                  </FormItem>
                )}
              />
              <br />
              <div className="flex space-x-2">
                <Button
                  type="submit"
                  className="w-full bg-purple-600 hover:bg-purple-700 text-white"
                  disabled={loading}
                >
                  Submit Songs
                </Button>
                <Button
                  className="bg-red-500"
                  disabled={loading}
                  onClick={onDelete}
                >
                  <Trash className="h-4 w-4" />
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
