import Text from "./Text";
import { useDefaults } from "./DefaultsContext";
import subbutText from "./translations/SubmitButton.json"

type SubmitButtonProps = {
  clickFn: (...args: any[]) => any,
}

function SubmitButton(props: SubmitButtonProps) {
  const {defaults} = useDefaults();

  return(
    <button
      class="green-button clickable"
      onClick={props.clickFn}>
        <Text 
          value={subbutText.submit}
          lang={defaults().lang}/>
    </button>
  )
}

export default SubmitButton