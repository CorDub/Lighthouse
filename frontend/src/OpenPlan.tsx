import { Dynamic } from "solid-js/web";
import { type Component } from "solid-js";
import "./styles/OpenPlan.css";

type OpenPlanProps = {
  comp: Component
}

function OpenPlan(props: OpenPlanProps) {
  return(
    <div class="open-plan">
      <Dynamic component={props.comp} />
    </div>
  )
}

export default OpenPlan;