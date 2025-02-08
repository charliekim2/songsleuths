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
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
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
import { Slider } from "@/components/ui/slider";
import { useRef, useState } from "react";
import ReactCanvasDraw from "react-canvas-draw";

const youtubeUrlRegex = /^.*\?v=([A-Za-z0-9_\-]{11})&.*$/;

const formSchema = z.object({
  nickname: z
    .string()
    .min(3, "Nickname must be at least 3 characters")
    .max(10, "Nickname must be 10 characters or less"),
  songs: z
    .array(
      z.object({
        url: z.string().regex(youtubeUrlRegex, "Must be a valid YouTube URL"),
      }),
    )
    .length(3, "You must submit exactly 3 songs"), // set to songNum from backend
});

type FormData = z.infer<typeof formSchema>;

function YouTubeEmbed({ url }: { url: string }) {
  const videoId = url.split("v=")[1];
  if (!videoId) {
    return <Card className="w-full aspect-video bg-black border-black"></Card>;
  }
  return (
    <div className="w-full aspect-video">
      <iframe
        src={`https://www.youtube.com/embed/${videoId}`}
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowFullScreen
        className="w-full h-full"
      />
    </div>
  );
}

export default function SubmitSongs() {
  const canvasRef = useRef<ReactCanvasDraw>(null);
  const [brushColor, setBrushColor] = useState("#000000");
  const [brushRadius, setBrushRadius] = useState(2);
  let submitDrawing = false;

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      nickname: "",
      songs: [{ url: "" }, { url: "" }, { url: "" }], // set to songNum from backend
    },
  });

  const { fields } = useFieldArray({
    name: "songs",
    control: form.control,
  });

  const onSubmit = (data: FormData) => {
    // @ts-expect-error - getDataURL is not in the types
    const drawingDataUrl = canvasRef.current?.getDataURL(
      "png", // Export canvas data as PNG
      false, // Export canvas data without background image
      "#FFFFFF", // Background color
    );
    console.log(submitDrawing);
    console.log({
      ...data,
      drawing: drawingDataUrl,
    });
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <Card className="w-full max-w-2xl bg-gray-800 text-white">
        <CardHeader>
          <CardTitle className="text-3xl font-bold text-center text-purple-400">
            Song Sleuths
          </CardTitle>
          <CardDescription className="text-center text-gray-400">
            Submit your songs
          </CardDescription>
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
              {fields.map((field, index) => (
                <FormField
                  key={field.id}
                  control={form.control}
                  name={`songs.${index}.url`}
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Song {index + 1} YouTube URL</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Enter YouTube URL"
                          {...field}
                          className="bg-gray-700 border-gray-600 text-white"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              ))}
              <div className="mt-4 mx-10">
                <h3 className="text-lg font-semibold mb-2">Song Previews</h3>
                <Carousel>
                  <CarouselContent>
                    {fields.map((_, index) => (
                      <CarouselItem key={index} className="max-w-60 mx-auto">
                        <YouTubeEmbed url={form.watch(`songs.${index}.url`)} />
                      </CarouselItem>
                    ))}
                  </CarouselContent>
                  <CarouselPrevious />
                  <CarouselNext />
                </Carousel>
              </div>
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
                              ref={canvasRef}
                              brushColor={brushColor}
                              brushRadius={brushRadius}
                              backgroundColor="#FFFFFF"
                              canvasWidth={300}
                              canvasHeight={300}
                              lazyRadius={0}
                              className="border border-gray-600 rounded-md"
                              onChange={() => (submitDrawing = true)}
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
                                  canvasRef.current?.clear();
                                  submitDrawing = false;
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
              <Button
                type="submit"
                className="w-full bg-purple-600 hover:bg-purple-700 text-white"
              >
                Submit Songs
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
