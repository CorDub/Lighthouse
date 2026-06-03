import { useUser } from "./UserContext.tsx";
import { Navigate } from "@solidjs/router";
import { type ParentProps } from "solid-js";

function AgencyProtectedRoute(props: ParentProps) {
  const { user } = useUser()

  if (user()?.role === "agency") {
    return (
      <>
        {props.children}
      </>
    )
  } else {
    return (
      <Navigate href="/" />
    )
  }
}

export default AgencyProtectedRoute;