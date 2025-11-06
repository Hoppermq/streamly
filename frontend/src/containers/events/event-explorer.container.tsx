import {cn} from "@/lib/utils.ts";
import {Tooltip, TooltipContent, TooltipTrigger} from "@/components/ui/tooltip.tsx";
import {ChevronDown, Clock} from "lucide-react";
import {useState} from "react";
import {Popover, PopoverContent, PopoverTrigger} from "@/components/ui/popover.tsx";

export const EventExplorerContainer = () => {
  const [isInputOpen, _setIsInputOpen] = useState<boolean>(false);
  return (
      <div className={cn(
        "flex items-center gap-2  transition-[color,box-shadow] outline-none",
        "ring-ring/50 ring-[0.5px] border p-1 w-fit",
        "border-main shadow-sm",
        "rounded-sm")
      }>
        <span>
          <Tooltip>
            <TooltipTrigger asChild>
              <Clock className={`h-4 w-4 ${isInputOpen ? "text-foreground" : "text-muted-foreground"}`}/>
            </TooltipTrigger>
            <TooltipContent side={'top'}>{"here goes the content"}</TooltipContent>
          </Tooltip>
        </span>
        <span className={'relative'}>
          <input
            type={'text'}
            placeholder={'Enter Time serie'}
            value={''}
            className={cn(
              'bg-transparent border-0 outline-none text-sm',
              'w-40 px-2 py-1',
              'placeholder-muted-foreground focus:placeholder-transparent'
            )}
          />
        </span>

        <span>
          <Popover>
            <PopoverTrigger asChild>
              <ChevronDown className={'h-4 w-4 text-muted-foreground'}/>
            </PopoverTrigger>
            <PopoverContent>
              HELOO
            </PopoverContent>
          </Popover>
        </span>
      </div>
  )
}
