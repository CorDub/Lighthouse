import { Show } from "solid-js";

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
        <button
          class="clickable"
          onClick={() => props.setStatus(false)}>
          {props.onText}
        </button>
      </Show>

      <Show when={!props.status}>
        <button
          class="clickable"
          onClick={() => props.setStatus(true)}>
          {props.offText}
        </button>
      </Show>
    </div>
  )
}

export default Switch;