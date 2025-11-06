import {EventContainer} from "@/containers/events/event.container.tsx";
import {EventExplorerContainer} from "@/containers/events/event-explorer.container.tsx";

export const EventPage = () => {
  return (
    <EventContainer>
      <EventExplorerContainer />
    </EventContainer>
  )
}


