import { Navigate } from "@solidjs/router";
import { useUser } from "./UserContext.tsx";
import type { ParentProps } from "solid-js";

function ProtectedRoute(props: ParentProps) {
  const { user } = useUser()

  if (!user()) {
    return (
      <Navigate href="/"/>
    )
  } else {
    return (
      <>
        {props.children}
      </>
    )
  }
}

export default ProtectedRoute;