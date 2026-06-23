import Errors from "./Errors.tsx";
import TextInput from "./TextInput";
import { type ErrorKey } from "./Errors.tsx";
import { Show, createSignal } from "solid-js";
import TimeRange from "./TimeRange.tsx";
import CreatorSelection from "./CreatorSelection.tsx";

function CreateReportModule() {
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  
  const [name, setName] = createSignal("");

  return (
    <div class="module">
      <form class="form-module">
        <div class="form-title">
          <h2>New report</h2>
        </div>
        <div class="form-name">
          <TextInput 
            errors={errors()}
            errorsSetFn={setErrors}
            value={name()}
            valueSetFn={setName}
            autofocus={false}
            name={"Name"}
            bgColor={"var(--pale-green)"}/>
          
          <Show when={errors().length > 0} >
            <Errors errors={errors()}/>
          </Show>
        </div>
        <div class="form-timerange"
          style={{"margin-top":"1rem"}}>
          <h3>Time</h3>
          <TimeRange />
        </div>
        <div class="form-users">
          <h3>Creators</h3>
          <CreatorSelection />
        </div>
        <div class="form-redes">

        </div>
      </form>
    </div>
  )
}

export default CreateReportModule;