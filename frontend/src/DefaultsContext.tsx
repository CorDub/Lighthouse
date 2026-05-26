import { 
  createContext, 
  type ParentProps, 
  useContext, 
  createSignal,
  type Setter,
  type Accessor,
} from "solid-js";
import { type Defaults, DefaultsSchema } from "./schemas/defaults.ts";

type DefaultsContext = {
  defaults : Accessor<Defaults>;
  setDefaults: Setter<Defaults>;
}

const DefaultsContext = createContext<DefaultsContext>();

export function DefaultsProvider(props: ParentProps) {
  const [defaults, setDefaults] = createSignal<Defaults>(
    DefaultsSchema.parse({})
  );

  return (
    <DefaultsContext.Provider value ={{ defaults, setDefaults }}>
      {props.children}
    </DefaultsContext.Provider>
  );
}

export function useDefaults() {
  const ctx = useContext(DefaultsContext);
  if (!ctx) throw new Error("useDefaults must be used within a DefaultsProvider");
  return ctx;
}

