import { createFileRoute, Outlet } from "@tanstack/react-router";

// Layout route for a single task. The detail view lives in the index child
// (app.tasks.$id.index.tsx); nested children like `/applications` render here
// through the <Outlet/>. Without this Outlet, navigating to a child route
// changes the URL but renders nothing.
export const Route = createFileRoute("/app/tasks/$id")({
  component: () => <Outlet />,
});
