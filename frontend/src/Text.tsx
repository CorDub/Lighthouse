import { useUser } from "./UserContext"

export type TextValues = {
  "en": string,
  "es": string
}

type TextProps = {
  value: TextValues,
  strReturn?: boolean
}

function Text(props: TextProps) {
  const { user } = useUser()
  const userResolved = user()
  const lang: keyof TextValues = userResolved?.language ?? "en"

  if (props.strReturn) {
    return props.value[lang]
  }

  return(
    <p>{props.value[lang]}</p>
  )
}

export default Text
