import { Navigate } from "@solidjs/router";
import { useUser } from "./UserContext.tsx";
import { type ParentProps, Show } from "solid-js";

function ProtectedRoute(props: ParentProps) {
  const { user, isLoading } = useUser()

  return (
    <Show when={!isLoading()} fallback={<div>Loading...</div>}>
      <Show when={user()} fallback={<Navigate href="/login" />}>
        {props.children}
      </Show>
    </Show>
  )
}

export default ProtectedRoute;