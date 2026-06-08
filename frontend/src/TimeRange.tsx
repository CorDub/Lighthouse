import { createSignal } from "solid-js";
import Switch from "./Switch.tsx"

function TimeRange() {
  const [indefinite, setIndefinite] = createSignal(true)

  return (
    <div class="timerange">
      <Switch 
        onText="Recurrring"
        offText="One-off"
        status={indefinite()}
        setStatus={setIndefinite} />
      {indefinite()
        ? <div>
            <input type="date"/>
          </div>
        : <div>
            <input type="date"/>
            <input type="date"/>
          </div>
      }
    </div>
  )
}

export default TimeRange;