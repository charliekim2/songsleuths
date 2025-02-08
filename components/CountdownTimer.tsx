"use client";

import { useState, useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";

interface CountdownTimerProps {
  targetTimestamp: number;
}

interface TimeLeft {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
}

export default function CountdownTimer({
  targetTimestamp,
}: CountdownTimerProps) {
  const [timeLeft, setTimeLeft] = useState(calculateTimeLeft());

  function calculateTimeLeft() {
    const difference = targetTimestamp * 1000 - Date.now();
    let timeLeft: TimeLeft = {
      days: 0,
      hours: 0,
      minutes: 0,
      seconds: 0,
    };

    if (difference > 0) {
      timeLeft = {
        days: Math.floor(difference / (1000 * 60 * 60 * 24)),
        hours: Math.floor((difference / (1000 * 60 * 60)) % 24),
        minutes: Math.floor((difference / 1000 / 60) % 60),
        seconds: Math.floor((difference / 1000) % 60),
      };
    }

    return timeLeft;
  }

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft(calculateTimeLeft());
    }, 1000);

    return () => clearInterval(timer);
  }, [targetTimestamp, calculateTimeLeft]); // Added calculateTimeLeft to dependencies

  const timeComponents = Object.keys(timeLeft).map((interval) => {
    if (!timeLeft[interval as keyof TimeLeft]) {
      return null;
    }

    return (
      <div key={interval} className="flex flex-col items-center">
        <span className="text-4xl font-bold">
          {timeLeft[interval as keyof TimeLeft]}
        </span>
        <span className="text-sm text-gray-400">{interval}</span>
      </div>
    );
  });

  return (
    <Card className="bg-gray-800 text-white">
      <CardContent className="p-6">
        <h2 className="text-xl font-semibold mb-4 text-center">
          Time Remaining
        </h2>
        <div className="flex justify-around">
          {timeComponents.length ? timeComponents : <span>Time is up!</span>}
        </div>
      </CardContent>
    </Card>
  );
}
