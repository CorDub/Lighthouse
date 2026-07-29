import { Show } from "solid-js";
import "./styles/Switch.css";

type SwitchProps = {
  onText: string,
  offText: string,
  status: boolean,
  setStatus: (value: boolean) => void,
}

function Switch(props: SwitchProps) {
  return (
    <div class="switch">
      <Show when={props.status}>
        <div>
          <p style={{"margin-bottom":"0.25rem"}}>Recurring / One-off</p>
          <button
            type="button"
            class="switch-button clickable"
            onClick={() => props.setStatus(false)}>
            {props.onText}
          </button>
        </div>
      </Show>

      <Show when={!props.status}>
        <div>
          <p style={{"margin-bottom":"0.25rem"}}>Recurring / One-off</p>
          <button
            class="switch-button clickable"
            onClick={() => props.setStatus(true)}>
            {props.offText}
          </button>
        </div>
      </Show>
    </div>
  )
}

export default Switch;