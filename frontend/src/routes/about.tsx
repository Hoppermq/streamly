import { createFileRoute} from "@tanstack/react-router";

const About = () => {
  return (
    <div>
      <h1>About Page</h1>
    </div>
  )
}

export const Route = createFileRoute('/about')({
  component: About,
})
