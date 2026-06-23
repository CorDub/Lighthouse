import { createSignal, For } from "solid-js";
import TextInput from "./TextInput.tsx";
import type { ErrorKey } from "./Errors.tsx";
import AddButton from "./AddButton";
import "./styles/CreatorSelection.css";

function CreatorSelection() {
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const [creators, setCreators] = createSignal([]);
  const [creator, setCreator] = createSignal("");

  return(
    <div class="creator-selection">
      <div class="cs-list">
        <For each={creators()}>{(creator, i) => 
          <p>{creator}</p>
        }</For>
      </div>
      <div class="cs-input">
        <div class="cs-text-input">
          <TextInput 
            errors={errors()}
            errorsSetFn={setErrors}
            value={creator()}
            valueSetFn={setCreator}
            autofocus={true}
            name={"Add a creator"}
            bgColor={"var(--pale-green)"}/>
        </div>
        <div class="cs-add-button">
          <AddButton />
        </div>
      </div>
    </div>

  )
}

export default CreatorSelection;