"use client";
import { useForm } from "react-hook-form";
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
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { CalendarIcon, Clock } from "lucide-react";
import { format } from "date-fns";
import { cn } from "@/lib/utils";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

const formSchema = z
  .object({
    gameTitle: z
      .string()
      .min(1, "Game title is required")
      .max(50, "Game title must be 50 characters or less"),
    deadline: z.date(),
    deadlineTime: z
      .string()
      .regex(/^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$/, "Invalid time format"),
    songCount: z.number().min(1).max(5),
  })
  .refine(
    (data) => {
      const now = new Date();
      const deadline = new Date(data.deadline);
      const [hours, minutes] = data.deadlineTime.split(":");
      deadline.setHours(Number.parseInt(hours), Number.parseInt(minutes));
      return deadline > now;
    },
    {
      message: "Deadline must be in the future",
      path: ["deadline"],
    },
  );

type FormData = z.infer<typeof formSchema>;

export default function CreateGame() {
  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      gameTitle: "",
      deadline: new Date(),
      deadlineTime: format(new Date(), "HH:mm"),
      songCount: 3,
    },
  });

  const onSubmit = (data: FormData) => {
    const combinedDeadline = new Date(data.deadline);
    const [hours, minutes] = data.deadlineTime.split(":");
    combinedDeadline.setHours(Number.parseInt(hours), Number.parseInt(minutes));
    const unixTimestamp = Math.floor(combinedDeadline.getTime() / 1000);
    console.log({
      gameTitle: data.gameTitle,
      deadline: unixTimestamp,
      songCount: data.songCount,
    });
    console.log(unixTimestamp);
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <Card className="w-full max-w-md bg-gray-800 text-white">
        <CardHeader>
          <CardTitle className="text-3xl font-bold text-center text-purple-400">
            Song Sleuths
          </CardTitle>
          <CardDescription className="text-center text-gray-400">
            Create a new game
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="gameTitle"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Game Title</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Enter game title"
                        {...field}
                        className="bg-gray-700 border-gray-600 text-white"
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="deadline"
                render={({ field }) => (
                  <FormItem className="flex flex-col">
                    <FormLabel>Submission Deadline</FormLabel>
                    <Popover>
                      <PopoverTrigger asChild>
                        <FormControl>
                          <Button
                            variant={"outline"}
                            className={cn(
                              "w-full justify-start text-left font-normal bg-gray-700 border-gray-600 text-white",
                              !field.value && "text-gray-400",
                            )}
                          >
                            <CalendarIcon className="mr-2 h-4 w-4" />
                            {field.value ? (
                              format(field.value, "PPP")
                            ) : (
                              <span>Pick a date</span>
                            )}
                          </Button>
                        </FormControl>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0 bg-white">
                        <Calendar
                          mode="single"
                          selected={field.value}
                          onSelect={field.onChange}
                          initialFocus
                          className="bg-white text-gray-900"
                        />
                      </PopoverContent>
                    </Popover>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="deadlineTime"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Deadline Time</FormLabel>
                    <FormControl>
                      <div className="flex items-center bg-gray-700 border border-gray-600 rounded-md">
                        <Clock className="ml-2 h-4 w-4 text-gray-400" />
                        <Input
                          type="time"
                          {...field}
                          className="bg-transparent border-none text-white focus:outline-none focus:ring-0"
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="songCount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Number of Songs</FormLabel>
                    <Select
                      onValueChange={(value) =>
                        field.onChange(Number.parseInt(value))
                      }
                      defaultValue={field.value.toString()}
                    >
                      <FormControl>
                        <SelectTrigger className="bg-gray-700 border-gray-600 text-white">
                          <SelectValue placeholder="Select number of songs" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent className="bg-gray-700 border-gray-600 text-white">
                        {[1, 2, 3, 4, 5].map((num) => (
                          <SelectItem key={num} value={num.toString()}>
                            {num}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button
                type="submit"
                className="w-full bg-purple-600 hover:bg-purple-700 text-white"
              >
                Create Game
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
