export type LanguageCode = "en" | "es"

export type TextValues = Record<LanguageCode, string>

type TextProps = {
  value: TextValues,
  lang: LanguageCode
}

function Text(props: TextProps) {
  return(
    <>{props.value[props.lang]}</>
  )
}

export default Text
