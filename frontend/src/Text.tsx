export type TextValues = {
  "en": string,
  "es": string
}

type TextProps = {
  value: TextValues,
  lang: "en" | "es"
}

function Text(props: TextProps) {
  return(
    <>{props.value[props.lang]}</>
  )
}

export default Text
