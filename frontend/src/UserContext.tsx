import { 
  createContext, 
  createSignal, 
  useContext,
  onMount,
  type ParentProps
} from "solid-js";
import { type User, UserSchema } from "./schemas/user.ts";
import { BASE_URL } from "./helpers/config.ts";

type UserContext = {
  user: () => User;
  setUser: (u: User) => void;
}

const UserContext = createContext<UserContext>();

export function UserProvider(props: ParentProps) {
  const [user, setUser] = createSignal<User>(null);

  onMount(
    async () => {
      try {
        const response = await fetch(`${BASE_URL}/api/checkAuth`,  {
          method: "GET",
          credentials: "include",
        })

        if (response.ok) {
          const data = await response.json()
          const parsedUser = UserSchema.parse(data.user);
          setUser(parsedUser);
        }
      } catch(error) {
        console.error(error)
      }
    }
  )

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
