import { createMemo } from "solid-js";

export type LanguageCode = "en" | "es"

export type TextValues = Record<LanguageCode, string>

type TextProps = {
  value: TextValues,
  lang: LanguageCode,
  var?: string[]
}

function Text(props: TextProps) {
  const text = createMemo(() => {
    let text = props.value[props.lang]

    if (props.var) {
      for (const variable of props.var) {
        text = text.replace("{}", variable)
      }
    }

    return text
  })
  

  return(
    <>{text}</>
  )
}

export default Text
