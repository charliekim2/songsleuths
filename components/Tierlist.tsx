"use client";

import { useState, useMemo } from "react";
import { DragDropContext, Droppable, Draggable } from "@hello-pangea/dnd";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import {
  TooltipContent,
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
} from "./ui/tooltip";

interface TierListState {
  [tierId: string]: string[]; // Maps tier IDs to arrays of item IDs
}

interface TierListProps {
  data: Tierlist;
  items: Song[];
  initialState?: TierListState;
}

interface DragResult {
  [tierId: string]: string[]; // Maps tier IDs to arrays of item IDs
}

function DraggableItem({ item, index }: { item: Song; index: number }) {
  return (
    <Draggable draggableId={item.id.toString()} index={index}>
      {(provided, snapshot) => (
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger>
              <div
                ref={provided.innerRef}
                {...provided.draggableProps}
                {...provided.dragHandleProps}
                className={`group relative w-20 h-20 rounded-lg overflow-hidden ${
                  snapshot.isDragging ? "ring-2 ring-purple-500" : ""
                }`}
              >
                <Image
                  src={item.album_art || "/placeholder.svg"}
                  alt={item.name}
                  fill
                  className="object-cover"
                />
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p>{item.name}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      )}
    </Draggable>
  );
}

export default function TierList({
  data,
  items,
  initialState = {},
}: TierListProps) {
  // Create a Set of valid item IDs and tier IDs for quick lookup
  const validItemIds = useMemo(
    () => new Set(items.map((item) => item.id.toString())),
    [items],
  );
  const validTierIds = useMemo(
    () => new Set(data.tiers.map((tier) => tier.id.toString())),
    [data.tiers],
  );

  // Initialize tierItems state with validated initial state
  const [tierItems, setTierItems] = useState<{ [key: string]: string[] }>(
    () => {
      // Start with empty arrays for each tier
      const initialTierItems: { [key: string]: string[] } = {
        unranked: [],
        ...data.tiers.reduce((acc, tier) => ({ ...acc, [tier.id]: [] }), {}),
      };

      // Create a Set to track which items have been placed
      const placedItems = new Set<string>();

      // Process initial state, validating tier IDs and item IDs
      Object.entries(initialState).forEach(([tierId, itemIds]) => {
        if (validTierIds.has(tierId)) {
          // Filter out invalid item IDs and add valid ones to the tier
          const validItems = itemIds.filter((id) => validItemIds.has(id));
          initialTierItems[tierId] = validItems;
          validItems.forEach((id) => placedItems.add(id));
        }
      });

      // Add remaining unplaced valid items to unranked section
      items.forEach((item) => {
        if (!placedItems.has(item.id.toString())) {
          initialTierItems.unranked.push(item.id.toString());
        }
      });

      return initialTierItems;
    },
  );

  // eslint-disable-next-line
  const onDragEnd = (result: any) => {
    const { source, destination } = result;

    if (!destination) return;

    if (
      source.droppableId === destination.droppableId &&
      source.index === destination.index
    )
      return;

    const sourceList = Array.from(tierItems[source.droppableId]);
    const destList =
      source.droppableId === destination.droppableId
        ? sourceList
        : Array.from(tierItems[destination.droppableId]);

    const [removed] = sourceList.splice(source.index, 1);

    if (source.droppableId === destination.droppableId) {
      sourceList.splice(destination.index, 0, removed);
    } else {
      destList.splice(destination.index, 0, removed);
    }

    setTierItems({
      ...tierItems,
      [source.droppableId]: sourceList,
      ...(source.droppableId === destination.droppableId
        ? {}
        : { [destination.droppableId]: destList }),
    });
  };

  const handleSubmit = () => {
    const result = Object.entries(tierItems).reduce((acc, [tierId, items]) => {
      if (tierId !== "unranked" && items.length > 0) {
        acc[tierId] = items;
      }
      return acc;
    }, {} as DragResult);

    console.log(JSON.stringify(result, null, 2));
  };

  const getItemById = (id: string) =>
    items.find((item) => item.id.toString() === id);

  return (
    <div className="w-full max-w-4xl mx-auto bg-gray-800 p-6 rounded-lg">
      <h2 className="text-2xl font-bold text-white mb-6">
        {data.type == "ranking"
          ? "Rank Songs by Preference"
          : "Match Songs to the Player"}
      </h2>

      <DragDropContext onDragEnd={onDragEnd}>
        {/* Tiers */}
        {data.tiers
          .sort((a, b) => a.rank - b.rank)
          .map((tier) => (
            <div key={tier.id} className="flex mb-2">
              <div className="w-24 bg-gray-700 flex items-center justify-center p-4 rounded-l-lg">
                {tier.drawing ? (
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger>
                        <img src={tier.drawing} alt={tier.name}></img>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>{tier.name}</p>
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                ) : (
                  <span className="text-xl font-bold text-white">
                    {tier.name}
                  </span>
                )}
              </div>
              <Droppable
                droppableId={tier.id.toString()}
                direction="horizontal"
              >
                {(provided, snapshot) => (
                  <div
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                    className={`flex-1 flex flex-wrap gap-2 p-2 min-h-[100px] rounded-r-lg ${
                      snapshot.isDraggingOver ? "bg-gray-600" : "bg-gray-700"
                    }`}
                  >
                    {tierItems[tier.id.toString()].map((itemId, index) => {
                      const item = getItemById(itemId);
                      if (!item) return null;

                      return (
                        <DraggableItem key={itemId} item={item} index={index} />
                      );
                    })}
                    {provided.placeholder}
                  </div>
                )}
              </Droppable>
            </div>
          ))}

        {/* Unranked section */}
        <div className="mt-8">
          <h3 className="text-xl font-bold text-white mb-2">Unranked Items</h3>
          <Droppable droppableId="unranked" direction="horizontal">
            {(provided, snapshot) => (
              <div
                ref={provided.innerRef}
                {...provided.droppableProps}
                className={`flex flex-wrap gap-2 p-4 min-h-[120px] rounded-lg ${
                  snapshot.isDraggingOver ? "bg-gray-600" : "bg-gray-700"
                }`}
              >
                {tierItems.unranked.map((itemId, index) => {
                  const item = getItemById(itemId);
                  if (!item) return null;

                  return (
                    <DraggableItem key={itemId} item={item} index={index} />
                  );
                })}
                {provided.placeholder}
              </div>
            )}
          </Droppable>
        </div>

        <Button
          onClick={handleSubmit}
          className="mt-6 w-full bg-purple-600 hover:bg-purple-700 text-white"
        >
          Submit Tierlist
        </Button>
      </DragDropContext>
    </div>
  );
}
