import { Portal } from "solid-js/web"
import { onMount } from "solid-js";
import Text from "./Text";
import { useDefaults } from "./DefaultsContext";
import alertText from "./translations/Alert.json";
import "./styles/Alert.css";

type AlertType = "confirmation" | "error" | "warning"

export type AlertKeys = keyof typeof alertText

type AlertProps = {
  message: AlertKeys,
  type: AlertType,
  alertOpen: boolean,
  setAlertOpen: (status: boolean) => void
}

function Alert(props: AlertProps) {
  const { defaults } = useDefaults();

  onMount(() => {
    setTimeout(() => props.setAlertOpen(false), 6000)
  })

  return (
    <Portal>
      <div class={`alert ${props.alertOpen ? "open" : ""} ${props.type} clickable`}
        onClick={() => props.setAlertOpen(false)}>
        <Text 
          value={alertText[props.message]}
          lang={defaults().lang}/>
      </div>
    </Portal>
  )
}

export default Alert