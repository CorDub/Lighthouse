import { createContext, createSignal, useContext} from "solid-js";
import { type ParentProps } from "solid-js";
import { type User } from "./schemas/user.ts";

type UserContext = {
  user: () => User;
  setUser: (u: User) => void;
}

export const UserContext = createContext<UserContext>();

export function UserProvider(props: ParentProps) {
  const [user, setUser] = createSignal<User>(null);

  return (
    <UserContext.Provider value ={{ user, setUser }}>
      {props.children}
    </UserContext.Provider>
  );
}

export function useUser() {
  const ctx = useContext(UserContext);
  if (!ctx) throw new Error("useUser must be used within a UserProvider");
  return ctx;
}
