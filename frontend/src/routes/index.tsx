import { createFileRoute} from "@tanstack/react-router";
import App from "../App.tsx";

const Index = () => {
  return (
    <div>
      <App />
    </div>
  )
}

export const Route = createFileRoute('/')({
  component: Index,
})
