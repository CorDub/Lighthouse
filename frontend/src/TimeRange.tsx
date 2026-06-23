import { createSignal } from "solid-js";
import Switch from "./Switch.tsx";
import "./styles/TimeRange.css";

function TimeRange() {
  const [indefinite, setIndefinite] = createSignal(true)

  return (
    <div class="timerange">
      <Switch 
        onText="Recurrring"
        offText="One-off"
        status={indefinite()}
        setStatus={setIndefinite} />
      <div class="tr-dates">
      {indefinite()
        ? <div style={{"margin-top":"0.5rem"}}
            class="tr-recurring">
            <p>Start date</p>
            <input type="date"/>
          </div>
        : <div style={{"margin-top":"0.5rem"}} 
            class="tr-oneoff">
            <div class="tr-startdate">
              <p>Start Date</p>
              <input type="date"/>
            </div>
            <div class="tr-enddate">
              <p>End Date</p>
              <input type="date"/>
            </div>
          </div>
      }
      </div>
    </div>
  )
}

export default TimeRange;